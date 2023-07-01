package funcie

import (
	"github.com/google/uuid"
	"time"
)

// MessageKind is the type of message that is being sent.
type MessageKind int

const (
	// MessageKindUnknown is an unknown message kind; generally indicates uninitialized data.
	MessageKindUnknown MessageKind = 0
	// MessageKindDispatch is a request to forward a message to a consumer.
	MessageKindDispatch MessageKind = 1
)

// Message represents a message to be sent through a tunnel.
type Message struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Kind is the type of message that is being sent.
	Kind MessageKind `json:"kind"`
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

// NewMessage creates a new message with the given data.
func NewMessage(application string, kind MessageKind, data []byte, ttl time.Duration) *Message {
	return &Message{
		ID:          uuid.New().String(),
		Application: application,
		Kind:        kind,
		Data:        data,
		Created:     time.Now().Truncate(time.Millisecond),
		Ttl:         ttl,
	}
}

// NewMessageWithPayload creates a new message with the given payload, which is serialized using funcie.MustSerialize.
func NewMessageWithPayload[T any](application string, kind MessageKind, payload T, ttl time.Duration) *Message {
	return NewMessage(application, kind, MustSerialize(payload), ttl)
}
