package provider

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/examples/go-lambda-url-bastion/provider/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLambdaProxy_Start(t *testing.T) {
	app := "app"
	ctx := context.Background()

	rawHandler := func(ctx context.Context, payload events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       "Hello world",
		}, nil
	}

	client := mocks.NewBastionClient(t)

	proxy := NewLambdaFunctionProxy(app, client, rawHandler)

	var handler lambda.Handler
	lambdaProxy := func(wrappedHandler interface{}) {
		handler = wrappedHandler.(lambda.Handler)
	}

	lambdaStart = lambdaProxy
	t.Cleanup(func() {
		lambdaStart = lambda.Start
	})

	proxy.Start()

	require.NotNil(t, handler)

	t.Run("successful request", func(t *testing.T) {
		req := events.LambdaFunctionURLRequest{}
		reqBytes := funcie.MustSerialize(req)
		urlResp := &events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       "Hello world",
		}

		respPayload := messages.NewForwardRequestResponsePayload(funcie.MustSerialize(urlResp))
		resp := funcie.NewResponse("id", funcie.MustSerialize(respPayload), nil)
		client.EXPECT().SendRequest(ctx, mock.Anything).Return(resp, nil)

		responseBytes, err := handler.Invoke(ctx, reqBytes)
		require.NoError(t, err)

		var response events.LambdaFunctionURLResponse
		require.NoError(t, json.Unmarshal(responseBytes, &response))

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Hello world", response.Body)
	})
}
