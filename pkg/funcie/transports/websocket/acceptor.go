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
	// It is the responsibility of the caller to close the connection.
	// If the connection is not accepted, the acceptor will write an error to the response writer.
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

func (acc *acceptor) Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (conn Connection, err error) {
	defer func() {
		if err != nil {
			rw.Header().Set("Connection", "close")
			req.Close = true
			_ = req.Body.Close()
		}
	}()
	if err := acc.authHandler(ctx, req); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return nil, fmt.Errorf("authorizing connection: %w", err)
	}

	socket, err := ws.Accept(rw, req, acc.acceptOpts)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return nil, fmt.Errorf("accepting connection: %w", err)
	}

	conn = NewConnection(socket)
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

// WithBearerAuthorizationHandler sets the authorization handler for the server to a bearer authorization handler.
// The token is the expected value of the Authorization header, with the kind "Bearer".
func WithBearerAuthorizationHandler(token string) AcceptorOpt {
	return WithAuthorizationHandler(func(ctx context.Context, req *http.Request) error {
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