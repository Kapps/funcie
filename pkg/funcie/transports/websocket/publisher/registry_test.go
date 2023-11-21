package publisher_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"log/slog"
	ws "nhooyr.io/websocket"
	"testing"
)

func TestRegistry_EndToEnd(t *testing.T) {
	t.Parallel()

	logger := slog.Default()
	ctx := context.Background()
	appId := faker.Word()
	conn := mocks.NewClientConnection(t)
	secondConn := mocks.NewClientConnection(t)
	reg := publisher.NewRegistry(logger)

	conn.EXPECT().ApplicationId().Return(appId)
	secondConn.EXPECT().ApplicationId().Return(appId)

	t.Run("acquiring a never registered connection", func(t *testing.T) {
		conn, err := reg.AcquireExclusive(ctx, appId)
		require.NoError(t, err)
		require.Nil(t, conn)
	})

	t.Run("registering a connection", func(t *testing.T) {
		err := reg.Register(ctx, conn)
		require.NoError(t, err)
	})

	t.Run("acquiring and releasing a registered connection", func(t *testing.T) {
		conn, err := reg.AcquireExclusive(ctx, appId)
		require.NoError(t, err)
		require.NotNil(t, conn)

		err = reg.ReleaseExclusive(ctx, appId, conn)
		require.NoError(t, err)
	})

	t.Run("overwriting a registered connection should close the old one", func(t *testing.T) {
		conn.EXPECT().Close(ws.StatusGoingAway, "replaced by new connection").Return(nil).Once()
		err := reg.Register(ctx, secondConn)
		require.NoError(t, err)
	})

	t.Run("acquiring should return the new connection", func(t *testing.T) {
		conn, err := reg.AcquireExclusive(ctx, appId)
		require.NoError(t, err)
		require.Equal(t, secondConn, conn)

		err = reg.ReleaseExclusive(ctx, appId, conn)
		require.NoError(t, err)
	})

	t.Run("unregistering a connection should close it", func(t *testing.T) {
		secondConn.EXPECT().Close(ws.StatusNormalClosure, "unregistered").Return(nil).Once()
		err := reg.Unregister(ctx, appId)
		require.NoError(t, err)
	})

	t.Run("acquiring a unregistered connection should return nil", func(t *testing.T) {
		conn, err := reg.AcquireExclusive(ctx, appId)
		require.NoError(t, err)
		require.Nil(t, conn)
	})
}

func TestRegistry_AcquireBetweenLifecycles(t *testing.T) {
	// Test scenario:
	// 1. Register a connection
	// 2. Acquire the connection
	// 3. Register a new connection, replacing the old one
	// 4. Release the old connection

	t.Parallel()

	logger := slog.Default()
	ctx := context.Background()
	appId := faker.Word()
	conn := mocks.NewClientConnection(t)
	secondConn := mocks.NewClientConnection(t)
	overwritten := make(chan struct{})

	reg := publisher.NewRegistry(logger)

	conn.EXPECT().ApplicationId().Return(appId)
	secondConn.EXPECT().ApplicationId().Return(appId)

	require.NoError(t, reg.Register(ctx, conn))

	acquiredConn, err := reg.AcquireExclusive(ctx, appId)
	require.NoError(t, err)
	require.Equal(t, conn, acquiredConn)

	go func() {
		conn.EXPECT().Close(ws.StatusGoingAway, "replaced by new connection").Return(nil).Once()
		require.NoError(t, reg.Register(ctx, secondConn))
		overwritten <- struct{}{}
	}()

	// This should block until the old connection is released.
	require.NoError(t, reg.ReleaseExclusive(ctx, appId, conn))

	<-overwritten

	// This should not block once the old connection is released.
	acquiredConn, err = reg.AcquireExclusive(ctx, appId)
	require.NoError(t, err)
	require.Equal(t, secondConn, acquiredConn)

	require.NoError(t, reg.ReleaseExclusive(ctx, appId, acquiredConn))
}
