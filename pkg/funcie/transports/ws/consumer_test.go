package ws_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	wsl "nhooyr.io/websocket"
	"testing"
	"time"
)

func TestConsumer_Subscribe(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("writes a subscribe message to the connection", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		jsonValue, err := json.Marshal(ws.ClientToServerMessage{
			Channel:     "channelName",
			RequestType: ws.ClientToServerMessageRequestTypeSubscribe,
		})

		mockSocket.EXPECT().Write(ctx, mock.Anything, jsonValue).Return(nil)

		err = consumer.Subscribe(ctx, mockSocket, "channelName")
		require.NoError(t, err)
	})

	t.Run("returns an error if the write fails", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		mockSocket.EXPECT().Write(ctx, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))

		err := consumer.Subscribe(ctx, mockSocket, "channelName")
		require.Error(t, err)
	})
}

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

func TestConsumer_Consume(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns an error if the connection fails", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(nil, nil, fmt.Errorf("error"))

		err := consumer.Consume(ctx, func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			return nil, nil
		})

		require.Error(t, err)
	})

	t.Run("subscribes to the channel, consumes and responds to a message", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)

		subscriberJsonValue, err := json.Marshal(ws.ClientToServerMessage{
			Channel:     "channelName",
			RequestType: ws.ClientToServerMessageRequestTypeSubscribe,
		})
		require.NoError(t, err)

		serverToClient := &funcie.Message{
			ID:      "S2C",
			Data:    []byte("DataS2C"),
			Created: time.Now(),
			Ttl:     600,
		}
		serverToClientJson, err := json.Marshal(serverToClient)
		require.NoError(t, err)

		clientToServer := &funcie.Response{
			ID:       "C2S",
			Data:     []byte("DataC2S"),
			Error:    nil,
			Received: time.Now(),
		}
		clientToServerJson, err := json.Marshal(clientToServer)
		require.NoError(t, err)

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)

		mockSocket.EXPECT().Write(ctx, wsl.MessageText, subscriberJsonValue).Return(nil)
		mockSocket.EXPECT().Read(ctx).Return(wsl.MessageText, serverToClientJson, nil)
		mockSocket.EXPECT().Write(ctx, wsl.MessageText, clientToServerJson).Return(nil)
		mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil)

		_ = consumer.Consume(ctx, func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			cancel()
			return clientToServer, nil
		})
	})

	t.Run("errors if can't subscribe", func(t *testing.T) {
		t.Parallel()

		subscriberJsonValue, err := json.Marshal(ws.ClientToServerMessage{
			Channel:     "channelName",
			RequestType: ws.ClientToServerMessageRequestTypeSubscribe,
		})
		require.NoError(t, err)

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)

		mockSocket.EXPECT().Write(ctx, wsl.MessageText, subscriberJsonValue).Return(fmt.Errorf("error"))
		mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil)

		err = consumer.Consume(ctx, func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			return nil, nil
		})

		require.Error(t, err)
	})

	t.Run("errors if can't read message", func(t *testing.T) {
		t.Parallel()

		subscriberJsonValue, err := json.Marshal(ws.ClientToServerMessage{
			Channel:     "channelName",
			RequestType: ws.ClientToServerMessageRequestTypeSubscribe,
		})
		require.NoError(t, err)

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)

		mockSocket.EXPECT().Write(ctx, wsl.MessageText, subscriberJsonValue).Return(nil)
		mockSocket.EXPECT().Read(ctx).Return(0, nil, fmt.Errorf("error123"))
		mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil)

		err = consumer.Consume(ctx, func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			return nil, nil
		})

		require.Errorf(t, err, "error123")
	})

	t.Run("errors if cant write response", func(t *testing.T) {
		t.Parallel()

		subscriberJsonValue, err := json.Marshal(ws.ClientToServerMessage{
			Channel:     "channelName",
			RequestType: ws.ClientToServerMessageRequestTypeSubscribe,
		})
		require.NoError(t, err)

		serverToClient := &funcie.Message{
			ID:      "S2C",
			Data:    []byte("DataS2C"),
			Created: time.Now(),
			Ttl:     600,
		}
		serverToClientJson, err := json.Marshal(serverToClient)
		require.NoError(t, err)

		clientToServer := &funcie.Response{
			ID:       "C2S",
			Data:     []byte("DataC2S"),
			Error:    nil,
			Received: time.Now(),
		}
		clientToServerJson, err := json.Marshal(clientToServer)
		require.NoError(t, err)

		wsClient := mocks.NewWebsocketClient(t)
		consumer := ws.NewConsumerWithWS(wsClient, "ws://localhost:8080", "channelName")
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)

		mockSocket.EXPECT().Write(ctx, wsl.MessageText, subscriberJsonValue).Return(nil)
		mockSocket.EXPECT().Read(ctx).Return(wsl.MessageText, serverToClientJson, nil)
		mockSocket.EXPECT().Write(ctx, wsl.MessageText, clientToServerJson).Return(fmt.Errorf("error123"))
		mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil)

		err = consumer.Consume(ctx, func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			return clientToServer, nil
		})

		require.Errorf(t, err, "error123")
	})

}
