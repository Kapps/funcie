package funcie

import (
	"encoding/json"
	"fmt"
	"time"
)

// Response is a response to a message sent to a tunnel, with either a generic JSON payload or an error.
type Response = ResponseBase[json.RawMessage]

// ResponseBase represents a response to a message sent to a tunnel, with either a generic payload or an error.
type ResponseBase[T any] struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Data is the actual message payload, or nil if an error occurred.
	// Exactly one of Data or Error are not nil.
	Data *T `json:"data,omitempty"`
	// Error is the error that occurred, or nil if no error occurred.
	// Exactly one of Data or Error are not nil.
	Error *ProxyError `json:"error,omitempty"`
	// Received is the time the response was received.
	Received time.Time `json:"received"`
}

// NewResponse creates a new response with the given data and the current time as the received time.
func NewResponse(id string, data []byte, error error) *Response {
	var formatted *json.RawMessage
	if len(data) > 0 {
		formatted = (*json.RawMessage)(&data)
	}
	return &Response{
		ID:       id,
		Data:     formatted,
		Received: time.Now().UTC().Truncate(time.Millisecond),
		Error:    NewProxyErrorFromError(error),
	}
}

// NewResponseWithPayload creates a new response with the given payload, which is serialized using funcie.MustSerialize, and the current time as the received time.
func NewResponseWithPayload[T any](id string, payload *T, error error) *ResponseBase[T] {
	return &ResponseBase[T]{
		ID:       id,
		Data:     payload,
		Received: time.Now().UTC().Truncate(time.Millisecond),
		Error:    NewProxyErrorFromError(error),
	}
}

// UnmarshalResponsePayload unmarshals the payload of the given response into the given payload type.
func UnmarshalResponsePayload[ResponseType ResponseBase[T], T any](response *Response) (*ResponseType, error) {
	var data T
	if response.Data != nil {
		err := json.Unmarshal(*response.Data, &data)
		if err != nil {
			return nil, err
		}
	}
	return &ResponseType{
		ID:       response.ID,
		Data:     &data,
		Received: response.Received,
		Error:    response.Error,
	}, nil
}

// MarshalResponsePayload marshals the payload of the given response into a JSON byte array.
func MarshalResponsePayload[T any](response *ResponseBase[T]) (*Response, error) {
	var raw *json.RawMessage

	if response.Data != nil {
		serialized, err := json.Marshal(response.Data)
		if err != nil {
			return nil, fmt.Errorf("marshalling response data: %w", err)
		}
		raw = (*json.RawMessage)(&serialized)
	}

	return &Response{
		ID:       response.ID,
		Data:     raw,
		Error:    response.Error,
		Received: response.Received,
	}, nil
}
