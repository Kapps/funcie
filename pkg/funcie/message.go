package funcie

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// MessageKind is the type of message that is being sent.
type MessageKind string

type Message MessageBase[json.RawMessage]

// MessageBase represents a message to be sent through a tunnel, typed to a specific payload kind.
// The generic message is aliased as Message, to allow for untyped use.
type MessageBase[T any] struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Kind is the type of message that is being sent.
	Kind MessageKind `json:"kind"`
	// Application is the name of the application that this message is for.
	Application string `json:"application"`
	// Payload is the actual message payload.
	Payload T `json:"payload"`
	// Created is the time the message was created.
	Created time.Time `json:"created"`
	// Ttl is the time to live for this message.
	// If the message is not processed by Created+Ttl, it should be discarded.
	Ttl time.Duration `json:"ttl"`
}

// NewMessage creates a new message with the given payload.
func NewMessage(application string, kind MessageKind, payload []byte, ttl time.Duration) *Message {
	serialized := json.RawMessage(payload)
	return &Message{
		ID:          uuid.New().String(),
		Application: application,
		Kind:        kind,
		Payload:     serialized,
		Created:     time.Now().Truncate(time.Millisecond),
		Ttl:         ttl,
	}
}

// NewMessageWithPayload creates a new message with the given payload, which is serialized using funcie.MustSerialize.
func NewMessageWithPayload[T any](application string, kind MessageKind, payload T, ttl time.Duration) *MessageBase[T] {
	return &MessageBase[T]{
		ID:          uuid.New().String(),
		Application: application,
		Kind:        kind,
		Payload:     payload,
		Created:     time.Now().Truncate(time.Millisecond),
		Ttl:         ttl,
	}
}
