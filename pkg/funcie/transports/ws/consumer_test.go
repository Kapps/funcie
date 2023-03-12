package ws_test

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConsumer_Connect(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns a valid connection", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)

		conn, err := consumer.Connect(ctx)
		require.NoError(t, err)
		require.Equal(t, mockSocket, conn)
	})

	t.Run("returns an error if the connection fails", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(nil, nil, fmt.Errorf("error"))

		conn, err := consumer.Connect(ctx)
		require.Error(t, err)
		require.Nil(t, conn)
	})
}
