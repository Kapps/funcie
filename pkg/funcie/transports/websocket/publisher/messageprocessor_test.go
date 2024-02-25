package publisher_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	pubMocks "github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestMessageHandler_ProcessMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connStore := pubMocks.NewConnectionStore(t)
	logger := slog.Default()
	handler := publisher.NewMessageProcessor(connStore, logger)

	t.Run("registration message", func(t *testing.T) {
		t.Parallel()

		conn := mocks.NewConnection(t)
		appID := "app"
		payload := messages.NewRegistrationRequestPayload(appID, funcie.MustNewEndpointFromAddress("http://localhost:8080"))
		message := funcie.NewMessage(appID, messages.MessageKindRegister, funcie.MustSerialize(payload))

		connStore.EXPECT().RegisterConnection(appID, conn).Once()

		resp, err := handler.ProcessMessage(ctx, conn, message)
		require.NoError(t, err)

		require.Equal(t, message.ID, resp.ID)
	})

	t.Run("deregistration message", func(t *testing.T) {
		t.Parallel()

		conn := mocks.NewConnection(t)
		appID := "app"
		message := funcie.NewMessage(appID, messages.MessageKindDeregister, nil)

		connStore.EXPECT().UnregisterConnection(appID).Return(conn, nil).Once()

		resp, err := handler.ProcessMessage(ctx, conn, message)
		require.NoError(t, err)

		require.NotNil(t, resp)
		require.Equal(t, message.ID, resp.ID)
	})
}
