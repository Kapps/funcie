package publisher_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWebsocketPublisher_Publish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	srv := mocks.NewServer(t)
	pub := publisher.NewWebsocketPublisher(srv)

	msg := funcie.NewMessage("app", messages.MessageKindForwardRequest, json.RawMessage(`{"data":"test"}`))

	t.Run("successful publish", func(t *testing.T) {
		srv.EXPECT().SendMessage(ctx, msg).
			Return(&funcie.Response{}, nil).Once()

		_, err := pub.Publish(ctx, msg)
		require.NoError(t, err)
	})

	t.Run("error publishing", func(t *testing.T) {
		expectedErr := errors.New("error publishing")
		srv.EXPECT().SendMessage(ctx, msg).
			Return(nil, expectedErr).Once()

		_, err := pub.Publish(ctx, msg)
		require.ErrorIs(t, err, expectedErr)
	})
}
