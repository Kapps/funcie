package websocket

import "encoding/json"

// PayloadKind is the kind of payload that is encompassed by an Envelope.
type PayloadKind int

const (
	// PayloadKindUnknown is an unknown payload kind. This should not intentionally be used.
	PayloadKindUnknown PayloadKind = iota
	// PayloadKindRequest is a request or message payload kind.
	PayloadKindRequest
	// PayloadKindResponse is a response payload kind.
	PayloadKindResponse
)

// Envelope is a wrapper around a payload that is sent over the websocket.
type Envelope struct {
	// Kind is the kind of payload that is encompassed by the envelope.
	Kind PayloadKind
	// Data is the raw data of the payload.
	Data *json.RawMessage
}

// NewEnvelope creates a new envelope with the given kind and data.
func NewEnvelope(kind PayloadKind, data *json.RawMessage) *Envelope {
	return &Envelope{
		Kind: kind,
		Data: data,
	}
}
