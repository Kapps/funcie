package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateTestServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(srv.Close)

	return srv
}
