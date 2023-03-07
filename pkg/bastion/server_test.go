package bastion_test

import (
	"bytes"
	"context"
	"github.com/Kapps/funcie/pkg/bastion"
	bastionMocks "github.com/Kapps/funcie/pkg/bastion/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const testAddress = "127.0.0.1:8080/dispatch"

func TestServer_Listen_Shutdown(t *testing.T) {
	stubs := makeServerAndListen(t)
	stubs.shutdown(t)
}

func TestServer_Request_InvalidMethod(t *testing.T) {
	stubs := makeServerAndListen(t)

	req, err := http.NewRequest("GET", testAddress, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	stubs.shutdown(t)
}

func TestServer_Request_InvalidPath(t *testing.T) {
	stubs := makeServerAndListen(t)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/invalid", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	stubs.shutdown(t)
}

func TestServer_Request_InternalError(t *testing.T) {
	// Currently, an invalid request gives an internal error, so can test via that.
	stubs := makeServerAndListen(t)

	req, err := http.NewRequest("POST", testAddress, bytes.NewBufferString("invalid body"))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	stubs.shutdown(t)
}

/*func TestServer_Request_Valid(t *testing.T) {
	stubs := makeServerAndListen(t)
	payload := json.RawMessage(`{"foo":"bar"}`)
	request := bastion.NewRequest("foo", &payload, map[string]string{"param": "one"})
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", testAddress, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

}*/

func makeServerAndListen(t *testing.T) *stubs {
	ctx := context.Background()
	handler := bastionMocks.NewRequestHandler(t)
	httpServer := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	server := bastion.NewServerWithHTTPServer(httpServer, handler)
	done := make(chan struct{})

	go func() {
		err := server.Listen()
		require.NoError(t, err)
		close(done)
	}()

	return &stubs{
		ctx:        ctx,
		handler:    handler,
		server:     server,
		httpServer: httpServer,
		done:       done,
	}
}

type stubs struct {
	ctx        context.Context
	handler    *bastionMocks.RequestHandler
	server     bastion.Server
	httpServer *http.Server
	done       chan struct{}
}

func (s *stubs) shutdown(tb testing.TB) {
	err := s.httpServer.Shutdown(nil)
	require.NoError(tb, err)

	<-s.done
}
