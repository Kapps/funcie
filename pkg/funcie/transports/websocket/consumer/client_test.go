package consumer_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie/testutils"
	. "github.com/Kapps/funcie/pkg/funcie/transports/websocket/consumer"
	"github.com/stretchr/testify/require"
	"net/http"
	ws "nhooyr.io/websocket"
	"testing"
)

func TestClient_Dial(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := NewClient(ClientOptions{
		AuthToken: "token",
	})

	srv := testutils.CreateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		socket, err := ws.Accept(w, r, nil)
		require.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		require.NoError(t, err)
		require.NotNil(t, socket)
		require.NoError(t, socket.Close(ws.StatusNormalClosure, ""))
	})

	t.Run("should dial the server", func(t *testing.T) {
		socket, err := client.Dial(ctx, srv.URL)
		require.NoError(t, err)
		require.NotNil(t, socket)
		require.NoError(t, socket.Close(ws.StatusNormalClosure, ""))
	})
}
