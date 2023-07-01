package bastion_test

import (
	. "github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	client := NewHTTPApplicationClient("http", http.DefaultClient)

	response, err := client.ProcessRequest(nil, funcie.Application{
		Endpoint: server.URL,
	}, nil)

	server.
}
