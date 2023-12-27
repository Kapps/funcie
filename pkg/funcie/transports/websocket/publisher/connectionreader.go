package publisher

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"log/slog"
)

type ResponseCallback = func(ctx context.Context, resp *funcie.Response)

// ConnectionReader reads messages from a websocket connection and routes them to the appropriate handler.
type ConnectionReader interface {
	// Register registers a new connection to be read from.
	// From this point on, all messages received from the connection will be routed to the appropriate handler.
	// This method is non-blocking.
	Register(ctx context.Context, conn ClientConnection) error
	// RegisterResponseCallback registers a callback for a specific message ID.
	// The callback will be called when a response with the same message ID is received.
	// This method is non-blocking.
	RegisterResponseCallback(ctx context.Context, messageId string, callback ResponseCallback) error
}

type connectionReader struct {
	connections       []ClientConnection
	connectionStore   ConnectionStore
	responseCallbacks map[string]ResponseCallback
	logger            *slog.Logger
}

func NewConnectionReader(connectionStore ConnectionStore, logger *slog.Logger) ConnectionReader {
	return &connectionReader{
		responseCallbacks: make(map[string]ResponseCallback),
		connectionStore:   connectionStore,
		logger:            logger,
	}
}

func (c *connectionReader) Register(ctx context.Context, conn ClientConnection) error {
	c.connections = append(c.connections, conn)

	go c.readLoop(ctx, conn)

	return nil
}
