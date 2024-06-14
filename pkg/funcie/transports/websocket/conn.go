package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	ws "nhooyr.io/websocket"
)

// Connection represents a websocket connection, either from a client or a server.
type Connection interface {
	// Close closes the connection with the given code and reason.
	// The code represents an RFC 6455 close status.
	Close(code ws.StatusCode, reason string) error
	// Read reads an entire message from the connection.
	// The details of deserialization are up to the implementation.
	Read(ctx context.Context, message interface{}) error
	// Write writes an entire message to the connection.
	// The details of serialization are up to the implementation.
	Write(ctx context.Context, message interface{}) error
}

type connection struct {
	ws *ws.Conn
}

// NewConnection wraps an existing websocket connection.
// Messages sent and received through this connection will be in JSON text format.
func NewConnection(ws *ws.Conn) Connection {
	return &connection{ws}
}

func (c *connection) Close(code ws.StatusCode, reason string) error {
	err := c.ws.Close(code, reason)
	if err != nil {
		return fmt.Errorf("closing connection: %w", err)
	}

	return nil
}

func (c *connection) Read(ctx context.Context, message interface{}) error {
	kind, content, err := c.ws.Read(ctx)
	if err != nil {
		return fmt.Errorf("reading message: %w", err)
	}

	if kind != ws.MessageText {
		return fmt.Errorf("invalid message type: %v", kind)
	}

	if err := json.Unmarshal(content, message); err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	return nil
}

func (c *connection) Write(ctx context.Context, message interface{}) error {
	content, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling message: %w", err)
	}

	if err := c.ws.Write(ctx, ws.MessageText, content); err != nil {
		return fmt.Errorf("writing message: %w", err)
	}

	return nil
}
