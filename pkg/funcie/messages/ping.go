package messages

import "github.com/Kapps/funcie/pkg/funcie"

// MessageKindPing is a request to validate whether a client is still alive.
// This can be used for both clients and bastions.
const MessageKindPing funcie.MessageKind = "PING"

// PingMessage is a message containing a ping request.
type PingMessage = funcie.MessageBase[PingRequestPayload]

// PingResponse is a message containing a ping response.
type PingResponse = funcie.ResponseBase[PingResponsePayload]

// PingRequestPayload is a ping request.
type PingRequestPayload struct{}

// NewPingRequestPayload creates a new PingRequestPayload.
func NewPingRequestPayload() *PingRequestPayload {
	return &PingRequestPayload{}
}

// PingResponsePayload is a response to a ping request.
type PingResponsePayload struct{}

// NewPingResponsePayload creates a new PingResponsePayload.
func NewPingResponsePayload() *PingResponsePayload {
	return &PingResponsePayload{}
}
