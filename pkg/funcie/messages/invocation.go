package messages

import "encoding/json"

// MessageKindForwardRequest is the kind of a message that is used to forward an application request.
const MessageKindForwardRequest = "FORWARD_REQUEST"

// ForwardRequestPayload is the payload for an invocation message.
type ForwardRequestPayload struct {
	Body json.RawMessage `json:"body"`
}

// ForwardRequestResponsePayload is the payload for an invocation response.
type ForwardRequestResponsePayload struct {
	Body json.RawMessage `json:"body"`
}
