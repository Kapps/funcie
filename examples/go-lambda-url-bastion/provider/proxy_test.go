package provider

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLambdaProxy_Start(t *testing.T) {
	app := "app"
	handler := func(ctx context.Context, payload events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       "Hello world",
		}, nil
	}

	proxy := NewLambdaFunctionProxy(app, nil, handler)
	require.NotNil(t, proxy)
}
