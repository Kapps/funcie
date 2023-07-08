package messages

import (
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
)

// MessageKindForwardRequest is the kind of a message that is used to forward an application request.
const MessageKindForwardRequest = "FORWARD_REQUEST"

// ForwardRequestMessage is a message containing an invocation request.
type ForwardRequestMessage funcie.MessageBase[ForwardRequestPayload]

// ForwardRequestResponse is a message containing an invocation response.
type ForwardRequestResponse funcie.ResponseBase[ForwardRequestResponsePayload]

// ForwardRequestPayload is the payload for an invocation message.
type ForwardRequestPayload struct {
	Body json.RawMessage `json:"body"`
}

// ForwardRequestResponsePayload is the payload for an invocation response.
type ForwardRequestResponsePayload struct {
	Body json.RawMessage `json:"body"`
}
