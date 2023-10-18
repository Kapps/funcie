package publisher

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"sync"
)

// Registry is a registry for websocket connections.
type Registry interface {
	// Register registers a websocket connection for the given application.
	Register(ctx context.Context, appId string, conn websocket.Connection) error
	// Unregister unregisters a websocket connection.
	Unregister(ctx context.Context, appId string) error
	// AcquireExclusive locks and returns a connection for the given application ID.
	// The connection must be released using ReleaseExclusive.
	// If the connection is already locked, this method blocks until the connection is released or the context is canceled.
	AcquireExclusive(ctx context.Context, appId string) (websocket.Connection, error)
	// ReleaseExclusive releases the exclusive lock for the given application ID.
	// If the connection is not locked, a panic will occur.
	ReleaseExclusive(ctx context.Context, appId string) error
}

type connectionWrapper struct {
	conn websocket.Connection
	lock sync.Mutex
}

type registry struct {
	connections sync.Map
}

// NewRegistry creates a new websocket registry.
func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) Register(ctx context.Context, appId string, conn websocket.Connection) error {
	if _, loaded := r.connections.LoadOrStore(appId, &connectionWrapper{conn: conn}); loaded {
		// If a connection is already registered for this application, we'll close the old connection.
		// The new connection will be registered instead.
		// Example scenario would be a client reconnecting after restarting the application.
		return ErrAlreadyRegistered
	}
	return nil
}
