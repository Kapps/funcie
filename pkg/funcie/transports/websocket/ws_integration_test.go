package websocket_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/consumer"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"time"
)

func TestWebsocket_EndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	t.Cleanup(cancel)

	logger := slog.Default()
	authToken := faker.Jwt()

	connStore := publisher.NewMemoryConnectionStore()
	responseNotifier := websocket.NewResponseNotifier()
	messageProcessor := publisher.NewMessageProcessor(connStore, logger)
	exchange := websocket.NewExchange(responseNotifier, messageProcessor, logger)
	acceptor := publisher.NewAcceptor(publisher.AcceptorOptions{
		AuthorizationHandler: publisher.BearerAuthorizationHandler(authToken),
	})

	srv := publisher.NewServer(connStore, exchange, acceptor, logger)
	pub := publisher.NewWebsocketPublisher(connStore, exchange, logger)

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

	consumerClient := consumer.NewDialer(consumer.DialerOptions{
		AuthToken: authToken,
	})
	router := utils.NewClientHandlerRouter()
	cons := consumer.NewConsumer(endpoint, router, consumerClient, logger)

	require.NoError(t, cons.Connect(ctx))

	consumerMessages := make(chan *funcie.Message, 1)
	require.NoError(t, cons.Subscribe(ctx, appId, func(ctx context.Context, msg *funcie.Message) (*funcie.Response, error) {
		consumerMessages <- msg
		return funcie.NewResponse(msg.ID, []byte("\"foo\""), nil), nil
	}))

	go func() {
		err := cons.Consume(ctx)
		require.ErrorIs(t, err, context.Canceled)
	}()

	time.Sleep(100 * time.Millisecond)

	message := funcie.NewMessage(appId, messages.MessageKindRegister, []byte("{}"))
	resp, err := pub.Publish(ctx, message)
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
