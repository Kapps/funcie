package messages

import (
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"time"
)

// MessageKindForwardRequest is the kind of a message that is used to forward an application request.
const MessageKindForwardRequest = "FORWARD_REQUEST"

// ForwardRequestMessage is a message containing an invocation request.
type ForwardRequestMessage funcie.MessageBase[ForwardRequestPayload]

// ForwardRequestResponse is a message containing an invocation response.
type ForwardRequestResponse funcie.ResponseBase[ForwardRequestResponsePayload]

// NewForwardRequestMessage creates a new ForwardRequestMessage with the given application name and payload.
func NewForwardRequestMessage(application string, payload ForwardRequestPayload, ttl time.Duration) *ForwardRequestMessage {
	return (*ForwardRequestMessage)(funcie.NewMessageWithPayload(application, MessageKindForwardRequest, payload, ttl))
}

// NewForwardRequestResponse creates a new ForwardRequestResponse with the given id, payload and error.
func NewForwardRequestResponse(id string, payload ForwardRequestResponsePayload, error *funcie.ProxyError) *ForwardRequestResponse {
	return (*ForwardRequestResponse)(funcie.NewResponseWithPayload(id, payload, error))
}

// ForwardRequestPayload is the payload for an invocation message.
type ForwardRequestPayload struct {
	Body json.RawMessage `json:"body"`
}

// ForwardRequestResponsePayload is the payload for an invocation response.
type ForwardRequestResponsePayload struct {
	Body json.RawMessage `json:"body"`
}
