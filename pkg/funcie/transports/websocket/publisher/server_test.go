package publisher_test

import (
	"context"
	"errors"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log"
	"log/slog"
	"net/http"
	ws "nhooyr.io/websocket"
	"testing"
	"time"
)

func TestServer_Listen_HappyPath(t *testing.T) {
	ctx := context.Background()
	acceptor := mocks.NewAcceptor(t)
	reg := mocks.NewRegistry(t)

	srv := publisher.NewServer(acceptor, reg, slog.Default())
	addr := "localhost:8086"
	endpoint := "http://" + addr

	go func() {
		err := srv.Listen(ctx, addr)
		require.NoError(t, err)
	}()

	t.Cleanup(func() {
		require.NoError(t, srv.Close())
	})

	time.Sleep(100 * time.Millisecond)

	var clientConn publisher.ClientConnection
	acceptor.EXPECT().Accept(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) (publisher.ClientConnection, error) {
			log.Println("Accepting connection")
			conn, err := ws.Accept(rw, req, nil)
			require.NoError(t, err)

			go func() {
				time.Sleep(100 * time.Millisecond)
				require.NoError(t, conn.Close(ws.StatusNormalClosure, "closing"))
			}()

			wsConn := websocket.NewConnection(conn)
			clientConn = publisher.NewClientConnection(wsConn, "app")
			return clientConn, nil
		}).Once()
	reg.EXPECT().Register(mock.Anything, mock.Anything).
		Return(nil).Once()

	conn, resp, err := ws.Dial(ctx, endpoint, nil)
	require.NoError(t, err)
	require.Equal(t, 101, resp.StatusCode)

	err = conn.Close(ws.StatusNormalClosure, "closing")
	require.NoError(t, err)
}

func TestServer_Listen_AcceptFailed(t *testing.T) {
	ctx := context.Background()
	acceptor := mocks.NewAcceptor(t)
	reg := mocks.NewRegistry(t)

	srv := publisher.NewServer(acceptor, reg, slog.Default())
	addr := "localhost:8086"
	endpoint := "http://" + addr

	go func() {
		err := srv.Listen(ctx, addr)
		require.NoError(t, err)
	}()

	t.Cleanup(func() {
		require.NoError(t, srv.Close())
	})

	time.Sleep(100 * time.Millisecond)

	acceptor.EXPECT().Accept(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("failed to accept connection")).Once()

	conn, resp, err := ws.Dial(ctx, endpoint, nil)
	require.Error(t, err)
	require.Equal(t, 400, resp.StatusCode)
	require.Nil(t, conn)
}
