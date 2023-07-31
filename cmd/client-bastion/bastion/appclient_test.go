package bastion_test

import (
	"context"
	"encoding/json"
	. "github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpApplicationClient_ProcessRequest(t *testing.T) {
	ctx := context.Background()
	resp := funcie.NewResponse("id", []byte("\"hello\""), nil)
	payload := "\"foo\""
	req := funcie.NewMessage("test-app", messages.MessageKindForwardRequest, json.RawMessage(payload))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.Equal(t, funcie.MustSerialize(req), body)
		_, err = w.Write(funcie.MustSerialize(resp))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := NewHTTPApplicationClient(http.DefaultClient)
	endpoint := funcie.MustNewEndpointFromAddress(server.URL)

	app := funcie.Application{
		Name:     "test-app",
		Endpoint: endpoint,
	}

	returned, err := client.ProcessRequest(ctx, app, req)
	require.NoError(t, err)

	require.Equal(t, returned, resp)
}
