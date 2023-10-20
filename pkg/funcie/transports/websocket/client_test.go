package websocket_test

import (
	"context"
	. "github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/stretchr/testify/require"
	"net/http"
	ws "nhooyr.io/websocket"
	"testing"
)

func TestClient_Dial(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := NewClient()

	srv := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		socket, err := ws.Accept(w, r, nil)
		require.NoError(t, err)
		require.NotNil(t, socket)
		require.NoError(t, socket.Close(ws.StatusNormalClosure, ""))
	})

	t.Run("should dial the server", func(t *testing.T) {
		socket, err := client.Dial(ctx, srv.URL, nil)
		require.NoError(t, err)
		require.NotNil(t, socket)
		require.NoError(t, socket.Close(ws.StatusNormalClosure, ""))
	})
}
