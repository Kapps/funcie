package consumer_test

import (
	"context"
	"errors"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/consumer"
	conMocks "github.com/Kapps/funcie/pkg/funcie/transports/websocket/consumer/mocks"
	wsMocks "github.com/Kapps/funcie/pkg/funcie/transports/websocket/mocks"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestConsumer_EndToEnd(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.Default()
	exchange := wsMocks.NewExchange(t)
	router := mocks.NewClientHandlerRouter(t)
	dialer := conMocks.NewDialer(t)
	serverUrl := "wss://example.com/ws"

	c := consumer.NewConsumer(serverUrl, exchange, router, dialer, logger)

	t.Run("should not allow consume/sub/unsub if not connected", func(t *testing.T) {
		ops := []func(context.Context) error{
			func(ctx context.Context) error {
				return c.Consume(ctx)
			},
			func(ctx context.Context) error {
				return c.Subscribe(ctx, "app-123", nil)
			},
			func(ctx context.Context) error {
				return c.Unsubscribe(ctx, "app-123")
			},
		}

		for _, op := range ops {
			err := op(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "not connected")
		}
	})

	t.Run("should connect successfully", func(t *testing.T) {
		conn := wsMocks.NewConnection(t)
		dialer.EXPECT().Dial(ctx, serverUrl).Return(conn, nil)

		require.NoError(t, c.Connect(ctx))
	})

	t.Run("should be able to subscribe to an app", func(t *testing.T) {

	}

	require.NoError(t, c.Connect(ctx), "Connect should not return an error")
	require.True(t, c.Connected(), "Consumer should be connected after Connect")
}

func TestConsumer_Connect_Fail(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New()
	exchange := mocks.NewExchange(t)
	router := mocks.NewClientHandlerRouter(t)
	dialer := mocks.NewDialer(t)
	serverUrl := "wss://invalid-url"

	c := consumer.NewConsumer(serverUrl, exchange, router, dialer, logger).(*consumer.Consumer)

	dialer.EXPECT().Dial(ctx, serverUrl).Return(nil, errors.New("dial error"))

	require.Error(t, c.Connect(ctx), "Connect should return an error")
	require.False(t, c.Connected(), "Consumer should not be connected after failed Connect")
}

func TestConsumer_Subscribe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New()
	exchange := mocks.NewExchange(t)
	router := mocks.NewClientHandlerRouter(t)
	dialer := mocks.NewDialer(t)
	serverUrl := "wss://example.com/ws"

	c := consumer.NewConsumer(serverUrl, exchange, router, dialer, logger).(*consumer.Consumer)
	c.SetConnected(true) // Simulate successful connection

	applicationId := "app-123"
	handler := func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
		return nil, nil
	}

	router.EXPECT().AddClientHandler(applicationId, handler).Return(nil)

	require.NoError(t, c.Subscribe(ctx, applicationId, handler), "Subscribe should not return an error")
}

func TestConsumer_Unsubscribe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New()
	exchange := mocks.NewExchange(t)
	router := mocks.NewClientHandlerRouter(t)
	dialer := mocks.NewDialer(t)
	serverUrl := "wss://example.com/ws"

	c := consumer.NewConsumer(serverUrl, exchange, router, dialer, logger).(*consumer.Consumer)
	c.SetConnected(true) // Simulate successful connection

	applicationId := "app-123"

	router.EXPECT().RemoveClientHandler(applicationId).Return(nil)

	require.NoError(t, c.Unsubscribe(ctx, applicationId), "Unsubscribe should not return an error")
}
