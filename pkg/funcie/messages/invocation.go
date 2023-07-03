package messages

import "encoding/json"

// MessageKindForwardRequest is the kind of a message that is used to forward an application request.
const MessageKindForwardRequest = "FORWARD_REQUEST"

// InvokeRequestPayload is the payload for an invoke message.
type InvokeRequestPayload struct {
	Body json.RawMessage `json:"body"`
}
