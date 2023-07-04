package consumer_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	utilMocks "github.com/Kapps/funcie/pkg/funcie/transports/utils/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/common"
	c "github.com/Kapps/funcie/pkg/funcie/transports/ws/consumer"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/consumer/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	wsl "nhooyr.io/websocket"
	"testing"
	"time"
)

// This is a test helper function that returns a connected c, a mock websocket client, and a mock websocket.
func getConnectedConsumer(t *testing.T, ctx context.Context) (*c.Consumer, *mocks.WebsocketClient, *mocks.Websocket) {
	wsClient := mocks.NewWebsocketClient(t)
	consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", utils.NewClientHandlerRouter())
	mockSocket := mocks.NewWebsocket(t)
	mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil).Maybe()

	wsClient.On("Dial", ctx, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)
	err := consumer.Connect(ctx)
	require.NoError(t, err)

	return consumer, wsClient, mockSocket
}

func nilHandler(_ context.Context, _ *funcie.Message) (*funcie.Response, error) {
	return nil, nil
}

func TestConsumer_Subscribe(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns an error if the connection is not established", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", utils.NewClientHandlerRouter())

		err := consumer.Subscribe(ctx, "channelName", nilHandler)
		require.Errorf(t, err, "not connected")
	})

	t.Run("writes a subscribe message to the connection", func(t *testing.T) {
		t.Parallel()

		consumer, _, mockSocket := getConnectedConsumer(t, ctx)

		jsonValue, err := json.Marshal(common.ClientToServerMessage{
			Application: "channelName",
			RequestType: common.ClientToServerMessageRequestTypeSubscribe,
		})

		mockSocket.EXPECT().Write(ctx, mock.Anything, jsonValue).Return(nil)

		err = consumer.Subscribe(ctx, "channelName", nilHandler)
		require.NoError(t, err)
	})

	t.Run("returns an error if the write fails", func(t *testing.T) {
		t.Parallel()

		consumer, _, mockSocket := getConnectedConsumer(t, ctx)

		mockSocket.EXPECT().Write(ctx, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))

		err := consumer.Subscribe(ctx, "channelName", nilHandler)
		require.Error(t, err)
	})
}

func TestConsumer_Unsubscribe(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns an error if the connection is not established", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", utils.NewClientHandlerRouter())

		err := consumer.Unsubscribe(ctx, "channelName")
		require.Errorf(t, err, "not connected")
	})

	t.Run("writes an unsubscribe message to the connection", func(t *testing.T) {
		t.Parallel()

		mockRouter := utilMocks.NewClientHandlerRouter(t)
		wsClient := mocks.NewWebsocketClient(t)
		consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", mockRouter)
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", ctx, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)
		err := consumer.Connect(ctx)
		require.NoError(t, err)

		jsonValue, err := json.Marshal(common.ClientToServerMessage{
			Application: "channelName",
			RequestType: common.ClientToServerMessageRequestTypeUnsubscribe,
		})

		mockSocket.EXPECT().Write(ctx, mock.Anything, jsonValue).Return(nil)
		mockRouter.EXPECT().RemoveClientHandler("channelName").Return(nil)
		err = consumer.Unsubscribe(ctx, "channelName")
		require.NoError(t, err)
	})

	t.Run("returns an error if the write fails", func(t *testing.T) {
		t.Parallel()

		consumer, _, mockSocket := getConnectedConsumer(t, ctx)

		mockSocket.EXPECT().Write(ctx, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))

		err := consumer.Unsubscribe(ctx, "channelName")
		require.Error(t, err)
	})
}

func TestConsumer_Connect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("returns no error on valid connection", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", utils.NewClientHandlerRouter())
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)

		err := consumer.Connect(ctx)
		require.NoError(t, err)
	})

	t.Run("disconnects when the context is cancelled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)

		wsClient := mocks.NewWebsocketClient(t)
		consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", utils.NewClientHandlerRouter())
		mockSocket := mocks.NewWebsocket(t)

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(mockSocket, nil, nil)
		mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil)

		err := consumer.Connect(ctx)
		require.NoError(t, err)

		cancel()

		//todo -- not sure how to test this without a sleep
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("returns an error if the connection fails", func(t *testing.T) {
		t.Parallel()

		wsClient := mocks.NewWebsocketClient(t)
		consumer := c.NewConsumerWithWS(wsClient, "ws://localhost:8080", utils.NewClientHandlerRouter())

		wsClient.On("Dial", mock.Anything, "ws://localhost:8080", mock.Anything).Return(nil, nil, fmt.Errorf("error"))

		err := consumer.Connect(ctx)
		require.Error(t, err)
	})
}

func TestConsumer_Consume(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("consumes and responds to a message", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)

		subscriberJsonValue, err := json.Marshal(common.ClientToServerMessage{
			Application: "app",
			RequestType: common.ClientToServerMessageRequestTypeSubscribe,
		})
		require.NoError(t, err)

		serverToClient := &funcie.Message{
			Application: "app",
			ID:          "S2C",
			Payload:     []byte("DataS2C"),
			Created:     time.Now().Truncate(0),
			Ttl:         600,
		}
		serverToClientJson, err := json.Marshal(serverToClient)
		require.NoError(t, err)

		clientToServer := &funcie.Response{
			ID:       "C2S",
			Data:     []byte("DataC2S"),
			Error:    nil,
			Received: time.Now().Truncate(0),
		}
		clientToServerJson, err := json.Marshal(clientToServer)
		require.NoError(t, err)

		consumer, _, mockSocket := getConnectedConsumer(t, ctx)

		mockSocket.EXPECT().Write(ctx, wsl.MessageText, subscriberJsonValue).Return(nil)
		mockSocket.EXPECT().Read(ctx).Return(wsl.MessageText, serverToClientJson, nil)
		mockSocket.EXPECT().Write(ctx, wsl.MessageText, clientToServerJson).Return(nil)
		//mockSocket.EXPECT().Close(wsl.StatusNormalClosure, mock.Anything).Return(nil)

		_ = consumer.Subscribe(ctx, "app", func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			require.Equal(t, serverToClient, message)
			cancel()
			return clientToServer, nil
		})

		_ = consumer.Consume(ctx)
	})

	t.Run("errors if can't read message", func(t *testing.T) {
		t.Parallel()

		consumer, _, mockSocket := getConnectedConsumer(t, ctx)

		mockSocket.EXPECT().Read(ctx).Return(0, nil, fmt.Errorf("error123"))

		err := consumer.Consume(ctx)

		require.Errorf(t, err, "error123")
	})

	t.Run("errors if can't write response", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(ctx)

		consumer, _, mockSocket := getConnectedConsumer(t, ctx)

		subscriberJsonValue, err := json.Marshal(common.ClientToServerMessage{
			Application: "app",
			RequestType: common.ClientToServerMessageRequestTypeSubscribe,
		})
		require.NoError(t, err)

		serverToClient := &funcie.Message{
			Application: "app",
			ID:          "S2C",
			Payload:     []byte("DataS2C"),
			Created:     time.Now().Truncate(0),
			Ttl:         600,
		}
		serverToClientJson, err := json.Marshal(serverToClient)
		require.NoError(t, err)

		clientToServer := &funcie.Response{
			ID:       "C2S",
			Data:     []byte("DataC2S"),
			Error:    nil,
			Received: time.Now().Truncate(0),
		}
		clientToServerJson, err := json.Marshal(clientToServer)
		require.NoError(t, err)

		mockSocket.EXPECT().Write(ctx, wsl.MessageText, subscriberJsonValue).Return(nil)
		mockSocket.EXPECT().Read(ctx).Return(wsl.MessageText, serverToClientJson, nil)
		mockSocket.EXPECT().Write(ctx, wsl.MessageText, clientToServerJson).Return(fmt.Errorf("error123"))

		_ = consumer.Subscribe(ctx, "app", func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
			require.Equal(t, serverToClient, message)
			cancel()
			return clientToServer, nil
		})

		err = consumer.Consume(ctx)
		require.Errorf(t, err, "error123")
	})
}
