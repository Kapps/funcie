package websocket_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/consumer"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestWebsocket_EndToEnd(t *testing.T) {
	logger := slog.Default()
	authToken := faker.Jwt()
	acceptor := publisher.NewAcceptor(publisher.AcceptorOptions{
		AuthorizationHandler: publisher.BearerAuthorizationHandler(authToken),
	})
	registry := publisher.NewRegistry(logger)
	srv := publisher.NewListener(acceptor, registry, logger)
	ctx := context.Background()
	host := "localhost:8086"
	endpoint := "http://" + host
	appId := faker.Word()

	go func() {
		err := srv.Listen(ctx, host)
		require.NoError(t, err)
	}()

	t.Cleanup(func() {
		require.NoError(t, srv.Close())
	})

	consumerClient := consumer.NewClient(consumer.ClientOptions{
		AuthToken: authToken,
	})
	router := utils.NewClientHandlerRouter()
	cons := consumer.NewConsumer(endpoint, router, consumerClient, logger)

	require.NoError(t, cons.Connect(ctx))

	// TODO: Subscribe doesn't really make sense given that the connection is for an application.

	consumerMessages := make(chan *funcie.Message, 1)
	require.NoError(t, cons.Subscribe(ctx, appId, func(ctx context.Context, msg *funcie.Message) (*funcie.Response, error) {
		consumerMessages <- msg
		return funcie.NewResponse(msg.ID, []byte("\"foo\""), nil), nil
	}))

	require.NoError(t, cons.Consume(ctx))

	conn, err := registry.AcquireExclusive(ctx, appId)
	require.NoError(t, err)
	require.NotNil(t, conn)

	t.Cleanup(func() {
		require.NoError(t, registry.ReleaseExclusive(ctx, appId, conn))
	})

	message := funcie.NewMessage(appId, messages.MessageKindRegister, []byte("{}"))
	resp, err := conn.Send(ctx, message)
	require.NoError(t, err)

	require.Equal(t, "\"foo\"", resp.Data)

	select {
	case msg := <-consumerMessages:
		require.Equal(t, message, msg)
	default:
		require.Fail(t, "expected message to be received")
	}

	require.NoError(t, cons.Unsubscribe(ctx, appId))

	select {
	case <-consumerMessages:
		require.Fail(t, "expected no message to be received")
	default:
	}
}
