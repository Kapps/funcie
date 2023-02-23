package tunnel

import "time"

// Response represents a response to a message sent to a tunnel.
type Response struct {
	// ID is the unique identifier for this message.
	ID string
	// Data is the actual message payload.
	Data []byte
	// Created is the time the message was created.
	Created time.Time
	// Received is the time the response was received.
	Received time.Time
}

type Consumer interface {
	// Consume consumes a message from the tunnel, synchronously waiting for a response from the other side.

}
