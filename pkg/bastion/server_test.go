package bastion_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/pkg/bastion"
	bastionMocks "github.com/Kapps/funcie/pkg/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

const testAddress = "http://127.0.0.1:8080/dispatch"

func TestServer_Listen_Shutdown(t *testing.T) {
	_ = makeServerAndListen(t)
}

func TestServer_Request_InvalidMethod(t *testing.T) {
	_ = makeServerAndListen(t)

	req, err := http.NewRequest("GET", testAddress, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestServer_Request_InvalidPath(t *testing.T) {
	_ = makeServerAndListen(t)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/invalid", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestServer_Request_InternalError(t *testing.T) {
	// Currently, an invalid request gives an internal error, so can test via that.
	_ = makeServerAndListen(t)

	req, err := http.NewRequest("POST", testAddress, bytes.NewBufferString("invalid body"))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestServer_Request_Valid(t *testing.T) {
	stubs := makeServerAndListen(t)

	payload := json.RawMessage(`{"foo":"bar"}`)
	request := bastion.NewRequest("foo", &payload, map[string]string{"param": "one"})
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)

	handlerResponse := funcie.NewResponse("resp", []byte(`{"foo":"bar"}`), nil)

	stubs.handler.EXPECT().Dispatch(mock.Anything, request).
		Return(handlerResponse, nil).Once()

	req, err := http.NewRequest("POST", testAddress, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response bastion.Response
	err = json.Unmarshal(respBytes, &response)
	require.NoError(t, err)
}

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

	s := &stubs{
		ctx:        ctx,
		handler:    handler,
		server:     server,
		httpServer: httpServer,
		done:       done,
	}

	t.Cleanup(func() { s.shutdown(t) })

	return s
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
