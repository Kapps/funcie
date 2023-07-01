package bastion_test

import (
	"context"
	. "github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPApplicationClient(t *testing.T) {
	t.Parallel()
	t.Run("invalid protocol", func(t *testing.T) {
		t.Parallel()
		require.Panics(t, func() {
			NewHTTPApplicationClient("invalid", http.DefaultClient)
		})
	})

	t.Run("http", func(t *testing.T) {
		t.Parallel()
		client := NewHTTPApplicationClient("http", http.DefaultClient)
		require.NotNil(t, client)
	})

	t.Run("https", func(t *testing.T) {
		t.Parallel()
		client := NewHTTPApplicationClient("https", http.DefaultClient)
		require.NotNil(t, client)
	})
}

func TestHttpApplicationClient_ProcessRequest(t *testing.T) {
	ctx := context.Background()
	resp := funcie.NewResponse("id", []byte("hello"), nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.Equal(t, "foo", body)
		_, err = w.Write(funcie.MustSerialize(resp))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPApplicationClient("http", http.DefaultClient)
	endpoint, err := funcie.NewEndpointFromAddress(server.URL)
	require.NoError(t, err)

	app := funcie.Application{
		Name:     "test-app",
		Endpoint: endpoint,
	}

	payload := "foo"

	req := funcie.NewMessage("test-app", funcie.MessageKindDispatch, []byte(payload), time.Minute)

	returned, err := client.ProcessRequest(ctx, app, req)
	require.NoError(t, err)

	require.Equal(t, returned, resp)

}
