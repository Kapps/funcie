package websocket

import (
	"context"
	"fmt"
	ws "nhooyr.io/websocket"
)

// Client allows connecting to a websocket server.
type Client interface {
	// Dial connects to the given URL.
	Dial(ctx context.Context, url string, opts *ws.DialOptions) (Connection, error)
}

type client struct {
}

// NewClient creates a new websocket client.
func NewClient() Client {
	return &client{}
}

func (c *client) Dial(ctx context.Context, url string, opts *ws.DialOptions) (Connection, error) {
	socket, _, err := ws.Dial(ctx, url, opts)
	if err != nil {
		return nil, fmt.Errorf("dialing %v: %w", url, err)
	}

	conn := NewConnection(socket)
	return conn, nil
}
