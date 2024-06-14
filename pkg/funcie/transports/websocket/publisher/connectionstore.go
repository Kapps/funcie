package publisher

import (
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"sync"
)

var ErrNoConnection = fmt.Errorf("no connection found for this application")

// ConnectionStore is a store for websocket connections to allow multiple applications to share the same connection.
type ConnectionStore interface {
	// GetConnection returns the connection for the given application.
	// If no connection is found, returns nil.
	GetConnection(app string) (websocket.Connection, error)
	// RegisterConnection registers the given connection for the given application, or associates an existing one.
	RegisterConnection(app string, conn websocket.Connection)
	// UnregisterConnection unregisters the connection for the given application.
	// Returns the connection that was unregistered, or nil if no connection was found.
	UnregisterConnection(app string) (websocket.Connection, error)
}

type connectionWrapper struct {
	conn websocket.Connection
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

func (c *connections) GetConnection(appId string) (websocket.Connection, error) {
	wrapper := c.getWrapperForConnection(appId, true)
	if wrapper == nil {
		return nil, ErrNoConnection
	}

	return wrapper.conn, nil
}

func (c *connections) RegisterConnection(appId string, conn websocket.Connection) {
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

func (c *connections) UnregisterConnection(appId string) (websocket.Connection, error) {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	wrapper := c.getWrapperForConnection(appId, false)
	if wrapper == nil {
		return nil, ErrNoConnection
	}

	_, loaded := wrapper.apps.LoadAndDelete(appId)
	if !loaded {
		panic("connection wrapper did not contain app ID")
	}

	conn := wrapper.conn
	return conn, nil
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
