package consumer

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"net/http"
	ws "nhooyr.io/websocket"
)

// Client allows connecting to a websocket server.
type Client interface {
	// Dial connects to the given URL.
	Dial(ctx context.Context, url string) (websocket.Connection, error)
}

type client struct {
	opts ClientOptions
}

// ClientOptions are options for the websocket client.
type ClientOptions struct {
	// AuthToken is the authentication token to use when connecting.
	AuthToken string
}

// NewClient creates a new websocket client.
func NewClient(opts ClientOptions) Client {
	return &client{
		opts: opts,
	}
}

func (c *client) Dial(ctx context.Context, url string) (websocket.Connection, error) {
	headers := http.Header{
		"Sec-Websocket-Protocol": []string{"funcie"},
	}
	if c.opts.AuthToken != "" {
		headers["Authorization"] = []string{fmt.Sprintf("Bearer %v", c.opts.AuthToken)}
	}

	opts := &ws.DialOptions{
		Subprotocols: []string{"funcie"},
		HTTPHeader:   headers,
	}

	socket, _, err := ws.Dial(ctx, url, opts)
	if err != nil {
		return nil, fmt.Errorf("dialing %v: %w", url, err)
	}

	conn := websocket.NewConnection(socket)
	return conn, nil
}
