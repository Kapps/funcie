package publisher_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestClientConn_ReadLoop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req := funcie.NewMessage("app", messages.MessageKindForwardRequest, []byte("\"input\""))
	resp := funcie.NewResponse(req.ID, []byte("\"output\""), nil)

	var activeHandler func(ctx context.Context, conn publisher.ClientConnection, msg *funcie.Message) (*funcie.Response, error)
	mockConn := mocks.NewConnection(t)
	mockNotifier := publisher.NewResponseNotifier() // TODO: Mock

	mockHandler := publisher.RequestHandler(func(ctx context.Context, conn publisher.ClientConnection, msg *funcie.Message) (*funcie.Response, error) {
		return activeHandler(ctx, conn, msg)
	})
	conn := publisher.NewClientConnection(
		ctx,
		mockConn,
		mockHandler,

	)

	messages := make(chan *websocket.Envelope)

	mockConn.EXPECT().Read(ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, payload interface{}) error {
			msg := payload.(*websocket.Envelope)
			*msg = *<- messages
			return nil
		})

	t.Run("request with nil response", func(t *testing.T) {
		mockConn.EXPECT().Write(ctx, req).Return(nil)

		c := websocket.NewClientConnection(ctx, mockConn, nil, nil, nil)

		go c.ReadLoop(ctx)

		messages <- req

		resp, err := c.Send(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	}
}
