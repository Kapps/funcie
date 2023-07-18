package provider_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/examples/go-lambda-url-bastion/provider"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLambdaBastionReceiver_Integration(t *testing.T) {
	listenerAddress := registerServer(t)

	t.Run("should forward requests to the handler", func(t *testing.T) {
		ev := events.LambdaFunctionURLRequest{
			QueryStringParameters: map[string]string{
				"name": "world",
			},
		}
		forwardRequestPayload := messages.NewForwardRequestPayload(funcie.MustSerialize(ev))
		forwardMessage := funcie.NewMessageWithPayload("app", messages.MessageKindForwardRequest, &forwardRequestPayload)
		requestBytes := funcie.MustSerialize(forwardMessage)

		resp, err := http.Post(listenerAddress.String(), "application/json", bytes.NewReader(requestBytes))
		require.NoError(t, err)

		respBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var responseMessage funcie.Message
		require.NoError(t, json.Unmarshal(respBytes, &responseMessage))

		require.Equal(t, forwardMessage.ID, responseMessage.ID)
		require.Equal(t, "Hello world", responseMessage.Payload)
	})

}

func registerServer(t *testing.T) funcie.Endpoint {
	applicationId := "app"
	registrationChannel := make(chan funcie.Endpoint)

	bastionRegistrationStubHandler := func(w http.ResponseWriter, r *http.Request) {
		req, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var message messages.RegistrationMessage
		require.NoError(t, json.Unmarshal(req, &message))

		require.Equal(t, messages.MessageKindRegister, message.Kind)

		respPayload := messages.NewRegistrationResponsePayload(uuid.New())
		resp := funcie.NewResponseWithPayload(message.ID, &respPayload, err)

		respBytes, err := json.Marshal(resp)
		require.NoError(t, err)

		_, err = w.Write(respBytes)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)

		registrationChannel <- message.Payload.Endpoint
	}

	bastionServer := httptest.NewServer(http.HandlerFunc(bastionRegistrationStubHandler))
	t.Cleanup(bastionServer.Close)

	handler := func(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       fmt.Sprintf("Hello %s", request.QueryStringParameters["name"]),
		}, nil
	}

	bastionUrl, err := url.Parse(fmt.Sprintf("http://%s", bastionServer.Listener.Addr().String()))
	require.NoError(t, err)

	receiver := provider.NewLambdaBastionReceiver(applicationId, "localhost:0", *bastionUrl, handler)
	t.Cleanup(receiver.Stop)

	go receiver.Start()

	return <-registrationChannel
}
