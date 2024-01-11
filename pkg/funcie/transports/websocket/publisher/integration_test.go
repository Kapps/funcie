package publisher_test

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	. "github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/stretchr/testify/require"
	"log/slog"
	ws "nhooyr.io/websocket"
	"syscall"
	"testing"
	"time"
)

func TestIntegration_ServerAccept(t *testing.T) {
	ctx := context.Background()

	_ = createIntegrationServer(t, ctx)

	t.Run("successful connection", func(t *testing.T) {
		_ = createIntegrationClient(t, ctx)
	})

	t.Run("unauthorized connection: no header", func(t *testing.T) {
		conn, resp, err := ws.Dial(ctx, "ws://localhost:8086", &ws.DialOptions{
			Subprotocols: []string{"funcie"},
		})
		require.Error(t, err)
		require.Nil(t, conn)

		require.Equal(t, 401, resp.StatusCode)
	})

	t.Run("unauthorized connection: incorrect bearer value", func(t *testing.T) {
		conn, resp, err := ws.Dial(ctx, "ws://localhost:8086", &ws.DialOptions{
			HTTPHeader: map[string][]string{
				"Authorization": {"Bearer bar"},
			},
		})
		require.Error(t, err)
		require.Nil(t, conn)

		require.Equal(t, 401, resp.StatusCode)
	})

	t.Run("unauthorized connection: incorrect header format", func(t *testing.T) {
		conn, resp, err := ws.Dial(ctx, "ws://localhost:8086", &ws.DialOptions{
			HTTPHeader: map[string][]string{
				"Authorization": {"foo"},
			},
		})
		require.Error(t, err)
		require.Nil(t, conn)

		require.Equal(t, 401, resp.StatusCode)
	})
}

func TestIntegration_ServerClose(t *testing.T) {
	ctx := context.Background()

	scaffold := createIntegrationServer(t, ctx)

	require.NoError(t, scaffold.close())
	scaffold.close = func() error { return nil }

	time.Sleep(100 * time.Millisecond)

	conn, resp, err := ws.Dial(ctx, "ws://localhost:8086", &ws.DialOptions{
		HTTPHeader: map[string][]string{
			"Authorization": {"Bearer foo"},
		},
	})
	require.Nil(t, resp)
	require.Nil(t, conn)

	require.ErrorIs(t, err, syscall.ECONNREFUSED)
}

func TestIntegration_Register(t *testing.T) {
	ctx := context.Background()

	scaffold := createIntegrationServer(t, ctx)

	t.Run("successful registration", func(t *testing.T) {
		conn := createIntegrationClient(t, ctx)
		regMessage := funcie.NewMessage("foo", messages.MessageKindRegister, nil)
		regPayload := json.RawMessage(funcie.MustSerialize(regMessage))

		err := conn.Write(ctx, ws.MessageText, funcie.MustSerialize(&websocket.Envelope{
			Kind: websocket.PayloadKindRequest,
			Data: &regPayload,
		}))
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		kind, respBytes, err := conn.Read(ctx)
		require.NoError(t, err)

		require.Equal(t, ws.MessageText, kind)

		var resp websocket.Envelope
		err = json.Unmarshal(respBytes, &resp)
		require.NoError(t, err)

		require.Equal(t, websocket.PayloadKindResponse, resp.Kind)

		var regResp funcie.Response
		err = json.Unmarshal(*resp.Data, &regResp)
		require.NoError(t, err)

		require.Equal(t, regResp.ID, regMessage.ID)

		storedConn, err := scaffold.connStore.GetConnection(regMessage.Application)
		require.NoError(t, err)
		require.NotNil(t, storedConn)
	})
}

func createIntegrationClient(t *testing.T, ctx context.Context) *ws.Conn {
	conn, resp, err := ws.Dial(ctx, "ws://localhost:8086", &ws.DialOptions{
		HTTPHeader: map[string][]string{
			"Authorization": {"Bearer foo"},
		},
		Subprotocols: []string{"funcie"},
	})
	require.NoError(t, err)
	require.NotNil(t, conn)
	require.NotNil(t, resp)

	t.Cleanup(func() {
		require.NoError(t, conn.Close(ws.StatusNormalClosure, "closing"))
	})

	return conn
}

func createIntegrationServer(t *testing.T, ctx context.Context) *serverScaffold {
	connStore := NewMemoryConnectionStore()
	responseNotifier := websocket.NewResponseNotifier()
	acceptor := NewAcceptor(AcceptorOptions{
		AuthorizationHandler: BearerAuthorizationHandler("foo"),
		UpgradeHandler:       DefaultUpgradeHandler(),
	})
	logger := slog.Default()

	srv := NewServer(connStore, responseNotifier, acceptor, logger)

	go func() {
		err := srv.Listen(ctx, "localhost:8086")
		require.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	scaffold := &serverScaffold{
		connStore:        connStore,
		responseNotifier: responseNotifier,
		acceptor:         acceptor,
		close:            srv.Close,
		logger:           logger,
		server:           srv,
	}

	t.Cleanup(func() {
		require.NoError(t, scaffold.close())
	})

	return scaffold
}

type serverScaffold struct {
	connStore        ConnectionStore
	responseNotifier websocket.ResponseNotifier
	acceptor         Acceptor
	close            func() error
	logger           *slog.Logger
	server           Server
}
