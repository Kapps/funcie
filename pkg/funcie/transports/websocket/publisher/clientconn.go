package publisher

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
)

// While this is a simple wrapper, it makes testing quite a bit easier.

// ClientConnection wraps a websocket connection for easier request-response functionality.
type ClientConnection interface {
	websocket.Connection

	// Send sends a message to the client, and waits for a response.
	Send(ctx context.Context, message *funcie.Message) (*funcie.Response, error)

	// ApplicationId returns the application ID of the client.
	ApplicationId() string
}

type clientConn struct {
	websocket.Connection
	applicationId string
}

// NewClientConnection creates a new ClientConnection from a websocket connection.
func NewClientConnection(conn websocket.Connection, applicationId string) ClientConnection {
	return &clientConn{conn, applicationId}
}

func (c *clientConn) Send(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	err := c.Write(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("writing message: %w", err)
	}

	var resp *funcie.Response
	err = c.Read(ctx, &resp)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return resp, nil
}

func (c *clientConn) ApplicationId() string {
	return c.applicationId
}
