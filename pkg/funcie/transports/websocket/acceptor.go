package websocket

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	ws "nhooyr.io/websocket"
	"strings"
)

// AuthorizationHandler is a function that handles authorization for a websocket connection.
// If the function returns an error, the connection will be closed and no requests accepted.
type AuthorizationHandler = func(ctx context.Context, req *http.Request) error

// AcceptorOpt is a function that modifies a websocket acceptor.
type AcceptorOpt = func(*acceptor)

// Acceptor allows accepting websocket connections from an existing HTTP request.
type Acceptor interface {
	// Accept accepts a websocket connection from the given HTTP request.
	// In case of an error, it is the responsibility of the caller to close the connection.
	// If no error, the returned connection must be closed by the caller.
	// The context is used for the lifetime of the connection.
	Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (Connection, error)
}

type acceptor struct {
	authHandler AuthorizationHandler
	acceptOpts  *ws.AcceptOptions
}

// NewAcceptor creates a new websocket acceptor with the given options.
// If no authorization handler is provided, all connections will be accepted.
// It is strongly recommended to provide an authorization handler.
func NewAcceptor(opts ...AcceptorOpt) Acceptor {
	acc := &acceptor{}
	for _, opt := range opts {
		opt(acc)
	}
	if acc.authHandler == nil {
		acc.authHandler = func(context.Context, *http.Request) error { return nil }
	}
	if acc.acceptOpts == nil {
		acc.acceptOpts = &ws.AcceptOptions{
			Subprotocols: []string{"funcie"},
		}
	}
	return acc
}

func (acc *acceptor) Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (Connection, error) {
	if err := acc.authHandler(ctx, req); err != nil {
		return nil, fmt.Errorf("authorizing connection: %w", err)
	}

	socket, err := ws.Accept(rw, req, acc.acceptOpts)
	if err != nil {
		return nil, fmt.Errorf("accepting connection: %w", err)
	}

	conn := NewConnection(socket)
	return conn, nil
}

// WithAuthorizationHandler sets the authorization handler for the acceptor.
func WithAuthorizationHandler(handler AuthorizationHandler) AcceptorOpt {
	if handler == nil {
		panic("authorization handler cannot be nil")
	}
	return func(s *acceptor) {
		s.authHandler = handler
	}
}

// WithBasicAuthorizationHandler sets the authorization handler for the server to a basic authorization handler.
// The token is the expected value of the Authorization header, with the kind "Basic".
func WithBasicAuthorizationHandler(token string) AcceptorOpt {
	return WithAuthorizationHandler(func(ctx context.Context, req *http.Request) error {
		auth := req.Header.Get("Authorization")
		kind, value, found := strings.Cut(auth, " ")
		if !found {
			return fmt.Errorf("authorization header required")
		}

		if kind != "Basic" {
			return fmt.Errorf("invalid authorization header")
		}

		if subtle.ConstantTimeCompare([]byte(value), []byte(token)) != 1 {
			return fmt.Errorf("invalid authorization token")
		}

		return nil
	})
}

// WithAcceptOptions sets the accept options for the server.
// These must include a subprotocol of "funcie".
func WithAcceptOptions(opts *ws.AcceptOptions) AcceptorOpt {
	if opts == nil {
		panic("accept options cannot be nil")
	}
	return func(s *acceptor) {
		s.acceptOpts = opts
	}
}
