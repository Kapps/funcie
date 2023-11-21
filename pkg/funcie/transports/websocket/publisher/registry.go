package publisher

import (
	"context"
	"log/slog"
	ws "nhooyr.io/websocket"
	"sync"
)

// Registry is a registry for websocket connections.
type Registry interface {
	// Register registers a websocket connection for the given application.
	Register(ctx context.Context, conn ClientConnection) error
	// Unregister unregisters a websocket connection.
	Unregister(ctx context.Context, appId string) error
	// AcquireExclusive locks and returns a connection for the given application ID.
	// The connection must be released using ReleaseExclusive.
	// If the connection is already locked, this method blocks until the connection is released or the context is canceled.
	// If the application is not registered, this method returns nil.
	AcquireExclusive(ctx context.Context, appId string) (ClientConnection, error)
	// ReleaseExclusive releases the exclusive lock for the given application ID.
	// If the connection is not locked, a panic will occur.
	ReleaseExclusive(ctx context.Context, appId string, conn ClientConnection) error
}

type connectionWrapper struct {
	conn ClientConnection
	lock sync.Mutex
}

type registry struct {
	connections sync.Map
	logger      *slog.Logger
}

// NewRegistry creates a new websocket registry.
func NewRegistry(logger *slog.Logger) Registry {
	return &registry{
		logger: logger,
	}
}

func (r *registry) Register(ctx context.Context, conn ClientConnection) error {
	wrapper := &connectionWrapper{
		conn: conn,
		lock: sync.Mutex{},
	}

	appId := conn.ApplicationId()

	if existing, loaded := r.connections.Swap(appId, wrapper); loaded {
		// If a connection is already registered for this application, we'll close the old connection.
		// The new connection will be registered instead.
		// Example scenario would be a client reconnecting after restarting the application.

		existingWrapper := existing.(*connectionWrapper)
		existingWrapper.lock.Lock()
		defer existingWrapper.lock.Unlock()
		err := existingWrapper.conn.Close(ws.StatusGoingAway, "replaced by new connection")
		if err != nil {
			r.logger.Error("Failed to close old connection", err, "appId", appId)
		}
	}

	r.logger.Info("Registered connection", "appId", appId)

	return nil
}

func (r *registry) Unregister(ctx context.Context, appId string) error {
	if existing, loaded := r.connections.LoadAndDelete(appId); loaded {
		wrapper := existing.(*connectionWrapper)
		wrapper.lock.Lock()
		defer wrapper.lock.Unlock()
		err := wrapper.conn.Close(ws.StatusNormalClosure, "unregistered")
		if err != nil {
			r.logger.Error("Failed to close connection", err, "appId", appId)
		}
	}
	return nil
}

func (r *registry) AcquireExclusive(ctx context.Context, appId string) (ClientConnection, error) {
	if existing, loaded := r.connections.Load(appId); loaded {
		wrapper := existing.(*connectionWrapper)
		wrapper.lock.Lock()
		return wrapper.conn, nil
	}
	return nil, nil
}

func (r *registry) ReleaseExclusive(ctx context.Context, appId string, conn ClientConnection) error {
	if existing, loaded := r.connections.Load(appId); loaded {
		if conn != existing.(*connectionWrapper).conn {
			// This connection is not the one we locked, so we'll just ignore it.
			r.logger.Error("Attempted to release exclusive lock for an old connection", "appId", appId)
		} else {
			wrapper := existing.(*connectionWrapper)
			wrapper.lock.Unlock()
			return nil
		}
	}
	return nil
}
