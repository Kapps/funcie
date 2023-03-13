package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"golang.org/x/exp/slog"
	ws "nhooyr.io/websocket"
)

func (c *Consumer) Connect(ctx context.Context) (Websocket, error) {
	conn, _, err := c.wsClient.Dial(ctx, c.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing Websocket: %w", err)
	}

	return conn, nil
}

func (c *Consumer) Subscribe(ctx context.Context, conn Websocket, channel string) error {
	r := ClientToServerMessage{
		Channel:     channel,
		RequestType: ClientToServerMessageRequestTypeSubscribe,
	}

	jsonValue, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = conn.Write(ctx, ws.MessageText, jsonValue)
	if err != nil {
		return fmt.Errorf("error writing to Websocket: %w", err)
	}
	return nil
}

// Consumer represents a consumer that consumes messages from a Redis channel.
type Consumer struct {
	channelName string
	URL         string
	wsClient    WebsocketClient
}

// NewConsumer creates a new Websocket consumer that consumes messages from the given URL.
func NewConsumer(url string, channelName string) funcie.Consumer {
	return &Consumer{
		URL:         url,
		channelName: channelName,
		wsClient:    &WebsocketClientWrapper{},
	}
}

// NewConsumerWithWS creates a new Websocket consumer that consumes messages from the given URL, with a given Websocket.
func NewConsumerWithWS(wsClient WebsocketClient, url string, channelName string) *Consumer {
	return &Consumer{
		wsClient:    wsClient,
		channelName: channelName,
		URL:         url,
	}
}

// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
func (c *Consumer) Consume(ctx context.Context, handler funcie.Handler) error {
	conn, err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to Websocket: %w", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "exiting consumer")

	err = c.Subscribe(ctx, conn, c.channelName)
	if err != nil {
		return fmt.Errorf("error subscribing to channel: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Warn("context cancelled", "err", ctx.Err())
			return ctx.Err()
		default:
			message, err := readMessage(ctx, conn)
			if err != nil {
				return fmt.Errorf("error reading message: %w", err)
			}

			response, err := handler(ctx, message)
			if err != nil {
				return fmt.Errorf("error handling message: %w", err)
			}

			responseData, err := formatResponse(response)
			if err != nil {
				return fmt.Errorf("error formatting response: %w", err)
			}

			if err := conn.Write(ctx, ws.MessageText, []byte(responseData)); err != nil {
				return fmt.Errorf("error writing message: %w", err)
			}
		}
	}
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
	var msg funcie.Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func formatResponse(response *funcie.Response) (string, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
