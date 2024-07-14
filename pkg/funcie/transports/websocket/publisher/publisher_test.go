package publisher_test

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	wsMocks "github.com/Kapps/funcie/pkg/funcie/transports/websocket/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestWebsocketPublisher_Publish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.Default()
	connStore := mocks.NewConnectionStore(t)
	exchange := wsMocks.NewExchange(t)

	pub := publisher.NewWebsocketPublisher(connStore, exchange, logger)

	msg := funcie.NewMessage("app", messages.MessageKindForwardRequest, json.RawMessage(`{"data":"test"}`))

	t.Run("successful publish", func(t *testing.T) {
		conn := wsMocks.NewConnection(t)
		connStore.EXPECT().GetConnection("app").
			Return(conn, nil).Once()
		exchange.EXPECT().Send(ctx, conn, msg).
			Return(&funcie.Response{}, nil).Once()

		_, err := pub.Publish(ctx, msg)
		require.NoError(t, err)
	})
}
