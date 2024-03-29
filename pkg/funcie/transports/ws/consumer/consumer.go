package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/common"
	"log"
	"log/slog"
	ws "nhooyr.io/websocket"
)

type wsConsumer struct {
	URL       string
	wsClient  WebsocketClient
	websocket Websocket
	connected bool
	router    utils.ClientHandlerRouter
}

// NewConsumer creates a new Websocket consumer that consumes messages from the given URL.
func NewConsumer(url string) funcie.Consumer {
	return &wsConsumer{
		URL:      url,
		wsClient: &WebsocketClientWrapper{},
		router:   utils.NewClientHandlerRouter(),
	}
}

// NewConsumerWithWS creates a new Websocket consumer that consumes messages from the given URL, with a given Websocket.
func NewConsumerWithWS(wsClient WebsocketClient, url string, router utils.ClientHandlerRouter) funcie.Consumer {
	return &wsConsumer{
		wsClient: wsClient,
		URL:      url,
		router:   router,
	}
}

func (c *wsConsumer) Connect(ctx context.Context) error {
	var err error
	c.websocket, err = c.connectSocket(ctx)
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-ctx.Done():
			err := c.websocket.Close(ws.StatusNormalClosure, "exiting consumer")
			c.connected = false
			if err != nil {
				log.Fatalf("error closing Websocket, was probably shutting down anyhow: %v", err)
			}
		}
	}()

	c.connected = true
	return nil
}

func (c *wsConsumer) connectSocket(ctx context.Context) (Websocket, error) {
	conn, _, err := c.wsClient.Dial(ctx, c.URL, &ws.DialOptions{
		Subprotocols: []string{"funcie"},
	})

	if err != nil {
		return nil, fmt.Errorf("error dialing Websocket: %w", err)
	}

	return conn, nil
}

func (c *wsConsumer) Subscribe(ctx context.Context, application string, handler funcie.Handler) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	r := common.ClientToServerMessage{
		Application: application,
		RequestType: common.ClientToServerMessageRequestTypeSubscribe,
	}

	jsonValue, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = c.websocket.Write(ctx, ws.MessageText, jsonValue)
	if err != nil {
		return fmt.Errorf("error writing to Websocket: %w", err)
	}

	err = c.router.AddClientHandler(application, handler)
	if err != nil {
		return fmt.Errorf("error adding handler: %w", err)
	}

	return nil
}

func (c *wsConsumer) Unsubscribe(ctx context.Context, channel string) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	r := common.ClientToServerMessage{
		Application: channel,
		RequestType: common.ClientToServerMessageRequestTypeUnsubscribe,
	}

	jsonValue, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = c.websocket.Write(ctx, ws.MessageText, jsonValue)
	if err != nil {
		return fmt.Errorf("error writing to Websocket: %w", err)
	}

	err = c.router.RemoveClientHandler(channel)
	if err != nil {
		return fmt.Errorf("error removing handler: %w", err)
	}

	return nil
}

// Consume starts the consume loop, reading from the Websocket and passing it to the router for handling.
func (c *wsConsumer) Consume(ctx context.Context) error {
	messageChannel := make(chan *funcie.Message, 10)

	var readError error = nil

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.WarnContext(ctx, "context cancelled", "err", ctx.Err())
				close(messageChannel)
				return
			default:
				message, err := readMessage(ctx, c.websocket)
				if err != nil {
					slog.ErrorContext(ctx, "error reading message", err)
					readError = err
					close(messageChannel)
					return
				}

				messageChannel <- message
			}
		}
	}()

	for message := range messageChannel {
		response, err := c.router.Handle(ctx, message)
		if err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}

		responseData, err := formatResponse(response)
		if err != nil {
			return fmt.Errorf("error formatting response: %w", err)
		}

		if err := c.websocket.Write(ctx, ws.MessageText, []byte(responseData)); err != nil {
			return fmt.Errorf("error writing message: %w", err)
		}
	}

	return readError
}

func readMessage(ctx context.Context, conn Websocket) (*funcie.Message, error) {
	messageType, message, err := conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	if messageType != ws.MessageText {
		return nil, fmt.Errorf("invalid message type: %v", messageType)
	}

	msg, err := parseMessage(string(message))
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func parseMessage(message string) (*funcie.Message, error) {
	var msg common.ServerToClientMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, err
	}

	if msg.RequestType == common.ServerToClientMessageRequestTypeRequest {
		return msg.Message, nil
	}

	return nil, fmt.Errorf("unsupported server to client message type: %v", msg.RequestType)
}

func formatResponse(response *funcie.Response) (string, error) {
	msg := &common.ClientToServerMessage{
		RequestType: common.ClientToServerMessageRequestTypeResponse,
		Response:    response,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
