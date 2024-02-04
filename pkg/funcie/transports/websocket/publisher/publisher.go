package publisher

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"log/slog"
)

type websocketPublisher struct {
	connStore ConnectionStore
	exchange  websocket.Exchange
	logger    *slog.Logger
}

// NewWebsocketPublisher creates a publisher that communicates over websockets.
func NewWebsocketPublisher(connStore ConnectionStore, exchange websocket.Exchange, logger *slog.Logger) funcie.Publisher {
	return &websocketPublisher{
		connStore: connStore,
		exchange:  exchange,
		logger:    logger,
	}
}

func (w *websocketPublisher) Publish(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	w.logger.DebugContext(ctx, "Publishing message", "messageId", message.ID, "application", message.Application)

	conn, err := w.connStore.GetConnection(message.Application)
	if err != nil {
		return nil, fmt.Errorf("getting connection for app %v: %w", message.Application, err)
	}

	resp, err := w.exchange.Send(ctx, conn, message)
	if err != nil {
		return nil, fmt.Errorf("sending message %v to application %v: %w", message.ID, message.Application, err)
	}

	return resp, nil
}
