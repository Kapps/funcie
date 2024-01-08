package publisher_test

import (
	wsMocks "github.com/Kapps/funcie/pkg/funcie/transports/websocket/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryConnectionStore(t *testing.T) {
	t.Parallel()

	store := publisher.NewMemoryConnectionStore()

	t.Run("register and get connection", func(t *testing.T) {
		t.Parallel()
		appID := "test-app"
		mockConn := wsMocks.NewConnection(t)

		store.RegisterConnection(appID, mockConn)
		retrievedConn, err := store.GetConnection(appID)

		require.NotNil(t, retrievedConn)
		require.NoError(t, err)
		require.Equal(t, mockConn, retrievedConn)
	})

	t.Run("unregister connection", func(t *testing.T) {
		t.Parallel()
		appID := "test-app"
		mockConn := wsMocks.NewConnection(t)

		store.RegisterConnection(appID, mockConn)
		unregisteredConn, err := store.UnregisterConnection(appID)

		require.NotNil(t, unregisteredConn)
		require.Equal(t, mockConn, unregisteredConn)

		reretrieved, err := store.GetConnection(appID)
		require.Nil(t, reretrieved)
		require.ErrorIs(t, err, publisher.ErrNoConnection)
	})

	t.Run("get connection for unregistered app", func(t *testing.T) {
		t.Parallel()
		appID := "non-existent-app"

		retrievedConn, err := store.GetConnection(appID)

		require.Nil(t, retrievedConn)
		require.ErrorIs(t, err, publisher.ErrNoConnection)
	})
}
