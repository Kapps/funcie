package publisher

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"net/http"
	ws "nhooyr.io/websocket"
	"strings"
)

// AuthorizationHandler is a function that handles authorization for a websocket connection.
// If the function returns an error, the connection will be closed and no requests accepted.
type AuthorizationHandler = func(ctx context.Context, req *http.Request) error

// Acceptor allows accepting websocket connections from an existing HTTP request.
type Acceptor interface {
	// Accept accepts a websocket connection from the given HTTP request.
	// It is the responsibility of the caller to close the connection.
	// If the connection is not accepted, the acceptor will write an error to the response writer.
	Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (ClientConnection, error)
}

type acceptor struct {
	opts AcceptorOptions
}

// AcceptorOptions are parameters for how to create a new websocket acceptor.
// The zero value is a valid acceptor, but it is strongly recommended to provide an authorization handler.
type AcceptorOptions struct {
	// AuthorizationHandler is invoked on a new request to authorize the connection.
	AuthorizationHandler AuthorizationHandler
	// AcceptOptions are options for accepting the websocket connection, passed to the underlying provider.
	// At a minimum, this must include a subprotocol of "funcie".
	AcceptOptions *ws.AcceptOptions
}

// NewAcceptor creates a new websocket acceptor with the given options.
// If no authorization handler is provided, all connections will be accepted.
// It is strongly recommended to provide an authorization handler.
func NewAcceptor(opts AcceptorOptions) Acceptor {
	if opts.AuthorizationHandler == nil {
		opts.AuthorizationHandler = func(context.Context, *http.Request) error { return nil }
	}
	if opts.AcceptOptions == nil {
		opts.AcceptOptions = &ws.AcceptOptions{
			Subprotocols: []string{"funcie"},
		}
	}
	return &acceptor{
		opts: opts,
	}
}

func (acc *acceptor) Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (conn ClientConnection, err error) {
	defer func() {
		if err != nil {
			rw.Header().Set("Connection", "close")
			req.Close = true
			_ = req.Body.Close()
		}
	}()
	if err := acc.opts.AuthorizationHandler(ctx, req); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return nil, fmt.Errorf("authorizing connection: %w", err)
	}

	app := req.Header.Get("X-Funcie-App")
	if app == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("missing X-Funcie-App header")
	}

	socket, err := ws.Accept(rw, req, acc.opts.AcceptOptions)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return nil, fmt.Errorf("accepting connection: %w", err)
	}

	wsConn := websocket.NewConnection(socket)
	conn = NewClientConnection(wsConn, app)

	return conn, nil
}

// BearerAuthorizationHandler returns an authorization handler for a bearer token.
// The token is the expected value of the Authorization header, with the kind "Bearer".
func BearerAuthorizationHandler(token string) AuthorizationHandler {
	return func(ctx context.Context, req *http.Request) error {
		auth := req.Header.Get("Authorization")
		kind, value, found := strings.Cut(auth, " ")
		if !found {
			return fmt.Errorf("authorization header required")
		}

		if kind != "Bearer" {
			return fmt.Errorf("invalid authorization header")
		}

		if subtle.ConstantTimeCompare([]byte(value), []byte(token)) != 1 {
			return fmt.Errorf("invalid authorization token")
		}

		return nil
	}
}
