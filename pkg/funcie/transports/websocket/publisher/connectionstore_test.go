package publisher_test

import (
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryConnectionStore(t *testing.T) {
	t.Parallel()

	store := publisher.NewMemoryConnectionStore()

	t.Run("register and get connection", func(t *testing.T) {
		t.Parallel()
		appID := "test-app"
		mockConn := mocks.NewClientConnection(t)

		store.RegisterConnection(appID, mockConn)
		retrievedConn := store.GetConnection(appID)

		require.NotNil(t, retrievedConn, "Expected non-nil connection")
		require.Equal(t, mockConn, retrievedConn, "Expected retrieved connection to match registered connection")
	})

	t.Run("unregister connection", func(t *testing.T) {
		t.Parallel()
		appID := "test-app"
		mockConn := mocks.NewClientConnection(t)

		store.RegisterConnection(appID, mockConn)
		unregisteredConn := store.UnregisterConnection(appID)

		require.NotNil(t, unregisteredConn, "Expected non-nil connection on unregister")
		require.Equal(t, mockConn, unregisteredConn, "Expected unregistered connection to match the original")
		require.Nil(t, store.GetConnection(appID), "Expected nil connection after unregistering")
	})

	t.Run("get connection for unregistered app", func(t *testing.T) {
		t.Parallel()
		appID := "non-existent-app"

		retrievedConn := store.GetConnection(appID)

		require.Nil(t, retrievedConn, "Expected nil connection for unregistered app")
	})
}
