package funcie

import (
	"context"
	"fmt"
	"time"
)

var ErrPubSubChannelClosed = fmt.Errorf("pubsub channel closed")

// Response represents a response to a message sent to a tunnel.
type Response struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Data is the actual message payload, or nil if an error occurred.
	// Exactly one of Data or Error are not nil.
	Data []byte `json:"data,omitempty"`
	// Error is the error that occurred, or nil if no error occurred.
	// Exactly one of Data or Error are not nil.
	Error *ProxyError `json:"error,omitempty"`
	// Received is the time the response was received.
	Received time.Time `json:"received"`
}

// NewResponse creates a new response with the given data and the current time as the received time.
func NewResponse(id string, data []byte, error error) *Response {
	return &Response{
		ID:       id,
		Data:     data,
		Received: time.Now().Truncate(time.Millisecond),
		Error:    NewProxyErrorFromError(error),
	}
}

// Handler is a function that handles a message from a tunnel.
type Handler func(ctx context.Context, message *Message) (*Response, error)

// Consumer represents a consumer of a synchronous tunnel that can be used to receive messages from a publisher and send a response.
type Consumer interface {
	// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
	Consume(ctx context.Context, handler Handler) error
}
