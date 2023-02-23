package tunnel

import (
	"context"
	"errors"
	"time"
)

// ErrNoConsumerActive is returned when a consumer is not active on a tunnel.
var ErrNoConsumerActive = errors.New("no consumer is active on this tunnel")

// Message represents a message to be sent through a tunnel.
type Message struct {
	// ID is the unique identifier for this message.
	ID string
	// Data is the actual message payload.
	Data []byte
	// Created is the time the message was created.
	Created time.Time
	// Ttl is the time to live for this message.
	// If the message is not processed by Created+Ttl, it should be discarded.
	Ttl time.Duration
}

// Publisher represents the publishing a synchronous tunnel that can be used to send messages to a consumer and wait for a response.
type Publisher interface {
	// Publish publishes a message to the tunnel, synchronously waiting for a response from the other side.
	Publish(ctx context.Context, message Message) (error, Response)
}

type RedisPublisher struct {
}
