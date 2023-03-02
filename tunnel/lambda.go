package tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/exp/slog"
)

// Start replaces the built-in lambda.Start function.
// If the current execution environment is Lambda, it will proxy requests to the
// tunnel instead of directly to the Lambda API (if a consumer is active).
// If the current execution environment is not Lambda, it will wait for proxied requests from the tunnel.
func Start(handler interface{}) {
	// First, check if we're running inside a Lambda.
	// If we are, then we need to wrap the handler in a LambdaHandler.
	var wrappedHandler lambda.Handler
	if IsRunningWithLambda() {
		slog.Info("running in Lambda")
		wrappedHandler = lambdaHandler(handler)
	} else {
		// TODO: Wrapped receiver
		slog.Info("running outside of Lambda")
		wrappedHandler = lambda.NewHandler(handler)
	}

	slog.Info("starting lambda")
	lambda.Start(wrappedHandler)
}

// lambdaHandler is the Lambda-level wrapper around the handler.
// It is responsible for publishing the message to the tunnel, and waiting for a response.
func lambdaHandler(userHandler interface{}) lambda.Handler {
	slog.Info("creating lambda handler")
	awsHandler := lambda.NewHandler(userHandler)
	slog.Info("created lambda handler")
	wrapper := func(ctx context.Context, message json.RawMessage) (json.RawMessage, error) {
		// TODO: publish message to tunnel
		slog.Info("publishing message to tunnel", "payload", message)
		bytes, err := message.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
		res, err := awsHandler.Invoke(ctx, bytes)
		slog.Info("received response from tunnel", "response", string(res), "err", err)
		if err != nil {
			return nil, fmt.Errorf("failed to invoke handler: %w", err)
		}

		rawResp := json.RawMessage{}
		err = rawResp.UnmarshalJSON(res)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return rawResp, nil
	}
	slog.Info("returning lambda handler")
	return lambda.NewHandler(wrapper)
}
