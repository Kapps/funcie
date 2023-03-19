package funcie

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

// ErrNoActiveConsumer is returned when a consumer is not active on a tunnel.
var ErrNoActiveConsumer = errors.New("no consumer is active on this tunnel")

// Message represents a message to be sent through a tunnel.
type Message struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Application is the name of the application that this message is for.
	Application string `json:"application"`
	// Data is the actual message payload.
	Data []byte `json:"data"`
	// Created is the time the message was created.
	Created time.Time `json:"created"`
	// Ttl is the time to live for this message.
	// If the message is not processed by Created+Ttl, it should be discarded.
	Ttl time.Duration `json:"ttl"`
}

func NewMessage(application string, data []byte, ttl time.Duration) *Message {
	return &Message{
		ID:          newId(),
		Application: application,
		Data:        data,
		Created:     time.Now().Truncate(time.Millisecond),
		Ttl:         ttl,
	}
}

func newId() string {
	return uuid.New().String()
}

// Publisher represents the publishing a synchronous tunnel that can be used to send messages to a consumer and wait for a response.
type Publisher interface {
	// Publish publishes a message to the tunnel, synchronously waiting for a response from the other side.
	// If no consumer is active, ErrNoConsumerActive is returned.
	Publish(ctx context.Context, message *Message) (*Response, error)
}
