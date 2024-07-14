package funcie

import (
	"context"
	"fmt"
)

var ErrPubSubChannelClosed = fmt.Errorf("pubsub channel closed")

// Handler is a function that handles a message from a tunnel.
type Handler func(ctx context.Context, message *Message) (*Response, error)

// Consumer represents a consumer of a synchronous tunnel that can be used to receive messages from a publisher and send a response.
type Consumer interface {
	// Connect opens the connection to the tunnel, allowing subscribing or consuming messages.
	// This method must be called before subscribing or consuming messages.
	Connect(ctx context.Context) error
	// Consume begins listening for messages from the tunnel.
	// This method blocks until the connection is closed or an error occurs.
	Consume(ctx context.Context) error
	// Subscribe subscribes to messages for a specific application from the tunnel.
	// The connection must be open before calling this method.
	Subscribe(ctx context.Context, applicationId string, handler Handler) error
	// Unsubscribe unsubscribes from messages for a specific application from the tunnel.
	// The connection must be open before calling this method.
	Unsubscribe(ctx context.Context, applicationId string) error
}
