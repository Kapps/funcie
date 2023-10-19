package websocket_test

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAcceptor_Accept(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	token := newAuthToken()
	opts := []websocket.AcceptorOpt{
		websocket.WithBearerAuthorizationHandler(token),
	}
	acceptor := websocket.NewAcceptor(opts...)

	t.Run("bearer auth (valid token)", func(t *testing.T) {
		srv := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			conn, err := acceptor.Accept(ctx, w, r)
			require.NoError(t, err)
			require.NotNil(t, conn)
		})

		request := createRequest(t, srv, token)
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		require.NoError(t, resp.Body.Close())
	})

	t.Run("bearer auth (invalid token)", func(t *testing.T) {
		srv := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			conn, err := acceptor.Accept(ctx, w, r)
			require.Error(t, err)
			require.Nil(t, conn)
		})

		request := createRequest(t, srv, "invalid")
		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		require.NoError(t, resp.Body.Close())
	})
}

func createTestServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(srv.Close)

	return srv
}

func createRequest(t *testing.T, srv *httptest.Server, token string) *http.Request {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Protocol", "funcie")

	return request
}
