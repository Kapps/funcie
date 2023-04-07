package funcie

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

// ErrNoActiveConsumer is returned when a consumer is not active on a tunnel.
var ErrNoActiveConsumer = errors.New("no consumer is active on this tunnel")

func newId() string {
	return uuid.New().String()
}

// Publisher represents the publishing a synchronous tunnel that can be used to send messages to a consumer and wait for a response.
type Publisher interface {
	// Publish publishes a message to the tunnel, synchronously waiting for a response from the other side.
	// If no consumer is active, ErrNoConsumerActive is returned.
	Publish(ctx context.Context, message *Message) (*Response, error)
}
