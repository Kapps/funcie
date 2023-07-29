package provider

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttpBastionClient_SendRequest(t *testing.T) {
	ctx := context.Background()
	requestPayload := messages.NewForwardRequestPayload([]byte("\"input\""))
	responsePayload := messages.NewForwardRequestResponsePayload([]byte("\"Hello world\""))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBytes, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req messages.ForwardRequestMessage
		require.NoError(t, json.Unmarshal(reqBytes, &req))

		require.Equal(t, "app", req.Application)
		require.Equal(t, *requestPayload, req.Payload)

		resp := funcie.NewResponseWithPayload("id", responsePayload, nil)
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(funcie.MustSerialize(resp))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	parsedUrl, err := url.Parse(server.URL)
	require.NoError(t, err)

	client := NewHTTPBastionClient(*parsedUrl)

	t.Run("should send requests to the bastion", func(t *testing.T) {
		req := funcie.NewMessage("app", messages.MessageKindForwardRequest, funcie.MustSerialize(requestPayload))
		resp, err := client.SendRequest(ctx, req)
		require.NoError(t, err)

		require.Equal(t, funcie.MustSerialize(responsePayload), []byte(*resp.Data))
	})
}
