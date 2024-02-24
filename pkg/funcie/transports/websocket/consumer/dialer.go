package consumer

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"net/http"
	ws "nhooyr.io/websocket"
)

// Dialer allows connecting to a websocket server.
type Dialer interface {
	// Dial connects to the given URL.
	Dial(ctx context.Context, url string) (websocket.Connection, error)
}

type dialer struct {
	opts DialerOptions
}

// DialerOptions are options for the websocket client.
type DialerOptions struct {
	// AuthToken is the authentication token to use when connecting.
	AuthToken string
}

// NewDialer creates a new websocket client.
func NewDialer(opts DialerOptions) Dialer {
	return &dialer{
		opts: opts,
	}
}

func (c *dialer) Dial(ctx context.Context, url string) (websocket.Connection, error) {
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
