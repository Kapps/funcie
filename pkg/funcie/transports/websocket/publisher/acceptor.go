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

// UpgradeHandler is a function that handles upgrading an HTTP request to a websocket connection.
type UpgradeHandler = func(ctx context.Context, rw http.ResponseWriter, req *http.Request, opts *ws.AcceptOptions) (websocket.Connection, error)

// Acceptor allows accepting websocket connections from an existing HTTP request.
type Acceptor interface {
	// Accept accepts a websocket connection from the given HTTP request.
	// It is the responsibility of the caller to close the connection.
	// If the connection is not accepted, the acceptor will write an error to the response writer.
	Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (websocket.Connection, error)
}

type acceptor struct {
	opts AcceptorOptions
}

// AcceptorOptions are parameters for how to create a new websocket acceptor.
// The zero value is not a valid acceptor as an authorization handler is required.
type AcceptorOptions struct {
	// AuthorizationHandler is invoked on a new request to authorize the connection.
	// If the function returns an error, the connection will be closed and no requests accepted.
	// The zero value is not valid and will cause NewAcceptor to panic.
	AuthorizationHandler AuthorizationHandler
	// AcceptOptions are options for accepting the websocket connection, passed to the underlying provider.
	// At a minimum, this must include a subprotocol of "funcie".
	// The zero value is valid and will be replaced with a default set of options.
	AcceptOptions *ws.AcceptOptions
	// UpgradeHandler is invoked to upgrade the HTTP request to a websocket connection.
	// The zero value is valid and will be replaced with the default upgrade handler.
	UpgradeHandler UpgradeHandler
}

// NewAcceptor creates a new websocket acceptor with the given options.
// While for most fields the zero value is valid, an authorization handler is required.
// If no authorization handler is provided, this method will panic.
func NewAcceptor(opts AcceptorOptions) Acceptor {
	if opts.AuthorizationHandler == nil {
		panic("acceptor authorization handler required")
	}
	if opts.AcceptOptions == nil {
		opts.AcceptOptions = &ws.AcceptOptions{
			Subprotocols: []string{"funcie"},
		}
	}
	if opts.UpgradeHandler == nil {
		opts.UpgradeHandler = DefaultUpgradeHandler()
	}
	return &acceptor{
		opts: opts,
	}
}

func (acc *acceptor) Accept(ctx context.Context, rw http.ResponseWriter, req *http.Request) (conn websocket.Connection, err error) {
	defer func() {
		if err != nil {
			rw.Header().Set("Connection", "close")
			req.Close = true
			_ = req.Body.Close()
		}
	}()
	if err := acc.opts.AuthorizationHandler(ctx, req); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return nil, fmt.Errorf("accepting connection: %w", err)
	}

	conn, err = acc.opts.UpgradeHandler(ctx, rw, req, acc.opts.AcceptOptions)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return nil, fmt.Errorf("accepting connection: %w", err)
	}

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

// DefaultUpgradeHandler is the default upgrade handler for a websocket connection.
func DefaultUpgradeHandler() UpgradeHandler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, opts *ws.AcceptOptions) (websocket.Connection, error) {
		socket, err := ws.Accept(rw, req, opts)
		if err != nil {
			return nil, fmt.Errorf("accepting websocket connection: %w", err)
		}

		wsConn := websocket.NewConnection(socket)
		return wsConn, nil
	}
}
