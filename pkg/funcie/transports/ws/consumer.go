package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"golang.org/x/exp/slog"
	"net/url"
	ws "nhooyr.io/websocket"
)

func (c *Consumer) Subscribe(ctx context.Context, url *url.URL, channel string) error {
	conn, _, err := ws.Dial(ctx, url.String(), nil)
	if err != nil {
		return fmt.Errorf("error dialing websocket: %w", err)
	}
	defer conn.Close(ws.StatusInternalError, "unexpected error, closing WebSocket")

	c.ws = conn

	r := ClientToServerMessage{
		Channel:     channel,
		RequestType: ClientToServerMessageRequestTypeSubscribe,
	}

	return c.writeJson(ctx, r)
}

func (c *Consumer) writeJson(ctx context.Context, v interface{}) (err error) {
	w, err := c.ws.Writer(ctx, ws.MessageText)
	if err != nil {
		return err
	}

	// json.Marshal cannot reuse buffers between calls as it has to return
	// a copy of the byte slice but Encoder does as it directly writes to w.
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return w.Close()
}

// Consumer represents a consumer that consumes messages from a Redis channel.
type Consumer struct {
	channelName string
	URL         *url.URL
	ws          websocket
}

// NewConsumer creates a new Websocket consumer that consumes messages from the given URL.
func NewConsumer(url *url.URL, channelName string) funcie.Consumer {
	return &Consumer{
		URL:         url,
		channelName: channelName,
	}
}

// NewConsumerWithWS creates a new Websocket consumer that consumes messages from the given URL, with a given websocket.
func NewConsumerWithWS(ws websocket, url *url.URL, channelName string) funcie.Consumer {
	return &Consumer{
		ws:          ws,
		channelName: channelName,
		URL:         url,
	}
}

func (c *Consumer) Close() error {
	return c.ws.Close(ws.StatusNormalClosure, "closing websocket")
}

// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
func (c *Consumer) Consume(ctx context.Context, handler funcie.Handler) error {
	err := c.Subscribe(ctx, c.URL, c.channelName)
	if err != nil {
		return fmt.Errorf("error subscribing to channel: %w", err)
	}
	defer funcie.CloseOrLog(fmt.Sprintf("pubsub channel %v", c.channelName), c)

	for {
		select {
		case <-ctx.Done():
			slog.Warn("context cancelled", "err", ctx.Err())
			return ctx.Err()
		default:
			message, err := readMessage(ctx, c.ws)
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

			if err := c.ws.Write(ctx, ws.MessageText, []byte(responseData)); err != nil {
				return fmt.Errorf("error writing message: %w", err)
			}
		}
	}
}

func readMessage(ctx context.Context, conn websocket) (*funcie.Message, error) {
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
