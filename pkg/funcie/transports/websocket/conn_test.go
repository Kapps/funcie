package websocket_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie/testutils"
	. "github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/stretchr/testify/require"
	"net/http"
	ws "nhooyr.io/websocket"
	"testing"
)

func TestConnection_EndToEnd(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	srv := testutils.CreateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		socket, err := ws.Accept(w, r, nil)
		require.NoError(t, err)

		conn := NewConnection(socket)

		err = conn.Write(ctx, "hello")
		require.NoError(t, err)

		var response string
		err = conn.Read(ctx, &response)
		require.NoError(t, err)
		require.Equal(t, "world", response)

		require.NoError(t, conn.Close(ws.StatusNormalClosure, ""))
	})

	client, _, err := ws.Dial(ctx, srv.URL, nil)
	require.NoError(t, err)

	conn := NewConnection(client)

	t.Run("should be able to read the initial message", func(t *testing.T) {
		var message string
		err = conn.Read(ctx, &message)
		require.NoError(t, err)

		require.Equal(t, "hello", message)
	})

	t.Run("should be able to write a message", func(t *testing.T) {
		err = conn.Write(ctx, "world")
		require.NoError(t, err)
	})

	t.Run("should be able to close the connection", func(t *testing.T) {
		require.NoError(t, conn.Close(ws.StatusNormalClosure, ""))
	})
}
