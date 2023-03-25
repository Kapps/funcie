package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/orcaman/concurrent-map/v2"
	"log"
	"net"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"os/signal"
	"time"
)

func Listen(port int32) error {
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	s := &http.Server{
		Handler:      NewWebsocketClientManager(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errorChan:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}

type ClientManager interface {
	AddClient(id string, conn Websocket)
	RemoveClient(id string)
	GetClient(id string) Websocket
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type WebsocketClientManager struct {
	clientMap cmap.ConcurrentMap[string, Websocket]
	logf      func(f string, v ...interface{})
}

func NewWebsocketClientManager() *WebsocketClientManager {
	return &WebsocketClientManager{
		clientMap: cmap.New[Websocket](),
		logf:      log.Printf,
	}
}

func (c *WebsocketClientManager) AddClient(id string, conn Websocket) {
	c.clientMap.Set(id, conn)
}

func (c *WebsocketClientManager) RemoveClient(id string) {
	c.clientMap.Remove(id)
}

func (c *WebsocketClientManager) GetClient(id string) (Websocket, error) {
	v, ok := c.clientMap.Get(id)

	if !ok {
		return nil, errors.New("client not found")
	}
	return v.(Websocket), nil
}

func (c *WebsocketClientManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"funcie"},
	})
	if err != nil {
		c.logf("%v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	if conn.Subprotocol() != "funcie" {
		conn.Close(websocket.StatusPolicyViolation, "client must speak the echo sub-protocol")
		return
	}

	err = c.RegisterClient(ctx, conn)
	if err != nil {
		c.logf("%v", err)
		return
	}
}

// RegisterClient waits for a subscribe message from the client (required) and then registers a client to the client manager.
func (c *WebsocketClientManager) RegisterClient(ctx context.Context, conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	for {
		typ, p, err := conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}
		if typ != websocket.MessageText {
			return fmt.Errorf("unexpected message type: %v", typ)
		}

		var message ClientToServerMessage
		if err := json.Unmarshal(p, &message); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		switch message.RequestType {
		case ClientToServerMessageRequestTypeSubscribe:
			c.AddClient(message.Channel, conn)
		}
	}
}

//type PublishClient interface {
//	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
//}

type Publisher struct {
	clientManager WebsocketClientManager
}

// NewPublisher creates a new RedisPublisher that publishes messages to the given channel.
func NewPublisher(redisClient PublishClient, channelName string) *Publisher {
	return &Publisher{redisClient: redisClient, channelName: channelName}
}

func (p *Publisher) Publish(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	// Publish the message to the channel.
	messageContents, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	pub := p.redisClient.Publish(ctx, p.channelName, messageContents)
	if err := pub.Err(); err != nil {
		return nil, fmt.Errorf("failed to publish message to channel %s: %w", p.channelName, err)
	}

	consumers, err := pub.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get result of publish: %w", err)
	}

	if consumers == 0 {
		return nil, funcie.ErrNoActiveConsumer
	}

	// Wait for a response from the consumer.
	responseKey := funcie.GetResponseKeyForMessage(message.ID)
	resp, err := p.redisClient.BRPop(ctx, message.Ttl, responseKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get response within ttl: %w", err)
	}

	// The response should be a list of two items, the first being the key and the second being the value.
	if len(resp) == 0 {
		return nil, fmt.Errorf("no response received within ttl of %s", message.Ttl)
	}

	if len(resp) != 2 {
		panic(fmt.Sprintf("expected response to be a list of two items, got %d", len(resp)))
	}

	// First entry is key; value should be the serialized response data.
	data := []byte(resp[1])

	var response *funcie.Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
