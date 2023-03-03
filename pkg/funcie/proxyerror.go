package funcie

import "encoding/json"

type proxyErrorJsonWrapper ProxyError

// ProxyError is an error that can be sent through a tunnel.
type ProxyError struct {
	// Message is the error message.
	Message string `json:"message,omitempty"`
}

// NewProxyError creates a new ProxyError with the given message.
func NewProxyError(message string) *ProxyError {
	return &ProxyError{
		Message: message,
	}
}

// NewProxyErrorFromError creates a new ProxyError from the given error.
// If the error is nil, nil is returned.
func NewProxyErrorFromError(err error) *ProxyError {
	if err == nil {
		return nil
	}
	return NewProxyError(err.Error())
}

// Error returns the error message.
func (e *ProxyError) Error() string {
	return e.Message
}

// MarshalJSON implements json.Marshaler.
func (e *ProxyError) MarshalJSON() ([]byte, error) {
	wrapper := proxyErrorJsonWrapper(*e)
	return json.Marshal(wrapper)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ProxyError) UnmarshalJSON(data []byte) error {
	var wrapper proxyErrorJsonWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	*e = ProxyError(wrapper)
	return nil
}
