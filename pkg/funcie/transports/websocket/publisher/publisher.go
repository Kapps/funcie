package publisher

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
)

type websocketPublisher struct {
	server Server
}

// NewWebsocketPublisher creates a publisher that communicates over websockets.
func NewWebsocketPublisher(server Server) funcie.Publisher {
	return &websocketPublisher{
		server: server,
	}
}

func (w *websocketPublisher) Publish(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	resp, err := w.server.SendMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("sending message: %w", err)
	}

	return resp, nil
}
