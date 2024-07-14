package websocket_test

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	wsMocks "github.com/Kapps/funcie/pkg/funcie/transports/websocket/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestExchange_Send(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	notifier := mocks.NewResponseNotifier(t)
	processor := wsMocks.NewMessageProcessor(t)
	conn := wsMocks.NewConnection(t)
	logger := slog.Default()
	exchange := websocket.NewExchange(notifier, processor, logger)

	req := funcie.NewMessage("app", messages.MessageKindForwardRequest, []byte("\"hello\""))
	reqBytes := json.RawMessage(funcie.MustSerialize(req))
	envelope := &websocket.Envelope{
		Kind: websocket.PayloadKindRequest,
		Data: &reqBytes,
	}

	conn.EXPECT().Write(ctx, envelope).Return(nil)

	resp := funcie.NewResponse(req.ID, []byte("\"world\""), nil)
	notifier.EXPECT().WaitForResponse(ctx, req.ID).Return(resp, nil)

	actualResp, err := exchange.Send(ctx, conn, req)
	require.NoError(t, err)
	require.Equal(t, resp, actualResp)

}
