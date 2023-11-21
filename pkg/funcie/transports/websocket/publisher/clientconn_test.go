package publisher_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/testutils"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/stretchr/testify/require"
	"net/http"
	ws "nhooyr.io/websocket"
	"testing"
)

func TestClientConn_Send(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := funcie.NewMessage("app", messages.MessageKindForwardRequest, []byte("\"input\""))
	resp := funcie.NewResponse(req.ID, []byte("\"output\""), nil)

	srv := testutils.CreateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		socket, err := ws.Accept(w, r, nil)
		require.NoError(t, err)

		conn := websocket.NewConnection(socket)
		clientConn := publisher.NewClientConnection(conn, "app")

		received, err := clientConn.Send(ctx, req)
		require.NoError(t, err)

		require.Equal(t, resp.ID, received.ID)

		require.NoError(t, clientConn.Close(ws.StatusNormalClosure, ""))
	})

	client, _, err := ws.Dial(ctx, srv.URL, nil)
	require.NoError(t, err)

	conn := websocket.NewConnection(client)

	var incoming funcie.Message
	err = conn.Read(ctx, &incoming)
	require.NoError(t, err)

	require.Equal(t, req.ID, incoming.ID)

	err = conn.Write(ctx, resp)
	require.NoError(t, err)

	require.NoError(t, conn.Close(ws.StatusNormalClosure, ""))
}
