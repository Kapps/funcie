package publisher

import (
	"sync"
)

// ConnectionStore is a store for websocket connections to allow multiple applications to share the same connection.
type ConnectionStore interface {
	// GetConnection returns the connection for the given application.
	// If no connection is found, returns nil.
	GetConnection(app string) *ClientConnection
	// RegisterConnection registers the given connection for the given application, or associates an existing one.
	RegisterConnection(app string, conn ClientConnection)
	// UnregisterConnection unregisters the connection for the given application.
	// Returns the connection that was unregistered, or nil if no connection was found.
	UnregisterConnection(app string) *ClientConnection
}

type connectionWrapper struct {
	conn ClientConnection
	apps sync.Map
}

type connections struct {
	connections    []*connectionWrapper // array to pointer, as we can modify elements during array modification.
	connectionLock sync.Mutex
}

// NewMemoryConnectionStore creates a new in-memory connection store.
func NewMemoryConnectionStore() ConnectionStore {
	return &connections{}
}

func (c *connections) GetConnection(appId string) *ClientConnection {
	wrapper := c.getWrapperForConnection(appId, true)
	if wrapper == nil {
		return nil
	}

	return &wrapper.conn
}

func (c *connections) RegisterConnection(appId string, conn ClientConnection) {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	wrapper := c.getWrapperForConnection(appId, false)
	if wrapper == nil {
		wrapper = &connectionWrapper{
			conn: conn,
			apps: sync.Map{},
		}
		c.connections = append(c.connections, wrapper)
	}

	wrapper.apps.Store(appId, struct{}{})
}

func (c *connections) UnregisterConnection(appId string) *ClientConnection {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	wrapper := c.getWrapperForConnection(appId, false)
	if wrapper == nil {
		return nil
	}

	_, loaded := wrapper.apps.LoadAndDelete(appId)
	if !loaded {
		panic("connection wrapper did not contain app ID")
	}

	conn := wrapper.conn
	return &conn
}

func (c *connections) getWrapperForConnection(appId string, acquireLock bool) *connectionWrapper {
	if acquireLock {
		c.connectionLock.Lock()
		defer c.connectionLock.Unlock()
	}

	for _, v := range c.connections {
		if _, found := v.apps.Load(appId); found {
			return v
		}
	}

	return nil
}
