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
	Connect(ctx context.Context) error
	Consume(ctx context.Context) error
	Subscribe(ctx context.Context, applicationId string, handler Handler) error
	Unsubscribe(ctx context.Context, applicationId string) error
}
