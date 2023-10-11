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
type AuthorizationHandler = func(r *http.Request) error

// ServerOpt is a function that modifies a websocket server.
type ServerOpt = func(*server)

// Server allows accepting websocket connections.
type Server interface {
	// Listen listens for websocket connections on the given address.
	Listen(ctx context.Context, addr string) error
}

type server struct {
	authHandler AuthorizationHandler
}

// NewServer creates a new websocket server with the given options.
// If no authorization handler is provided, all connections will be accepted.
// It is strongly recommended to provide an authorization handler.
func NewServer(opts ...ServerOpt) Server {
	svr := &server{}
	for _, opt := range opts {
		opt(svr)
	}
	if svr.authHandler == nil {
		svr.authHandler = func(r *http.Request) error { return nil }
	}
	return svr
}

func (s *server) Listen(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		}),
	}
}

func (s *server) accept(rw http.ResponseWriter, r *http.Request) (Connection, error) {
	if err := s.authHandler(r); err != nil {
		// TODO: Is this the best way to close the connection?
		r.Close = true
		r.Header.Set("Connection", "close")
		_ = r.Body.Close()
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized"))
		return nil, fmt.Errorf("authorizing connection: %w", err)
	}

	socket, err := ws.Accept(rw, r, opts)
	if err != nil {
		return nil, fmt.Errorf("accepting connection: %w", err)
	}

	conn := NewConnection(socket)
	return conn, nil
}

// WithAuthorizationHandler sets the authorization handler for the server.
func WithAuthorizationHandler(handler AuthorizationHandler) ServerOpt {
	if handler == nil {
		panic("authorization handler cannot be nil")
	}
	return func(s *server) {
		s.authHandler = handler
	}
}

// WithBasicAuthorizationHandler sets the authorization handler for the server to a basic authorization handler.
// The token is the expected value of the Authorization header, with the kind "Basic".
func WithBasicAuthorizationHandler(token string) ServerOpt {
	return WithAuthorizationHandler(func(r *http.Request) error {
		auth := r.Header.Get("Authorization")
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
