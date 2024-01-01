package publisher

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
)

// ConnectionManager is responsible for managing all connections and communicating with them.
type ConnectionManager interface {
	// RegisterConnection registers a connection with the manager.
	RegisterConnection(ctx context.Context, conn ClientConnection) error
	// SendMessage sends a message to the application it is intended for.
	// If no connection is found for the given application, ErrNoConnection is returned.
	SendMessage(ctx context.Context, message *funcie.Message) error
}

type messageSender struct {
	connStore ConnectionStore
}

// NewMessageSender creates a new MessageSender.
func NewMessageSender(connStore ConnectionStore) MessageSender {
	return &messageSender{
		connStore: connStore,
	}
}
