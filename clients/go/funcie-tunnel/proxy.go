package funcie_tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/exp/slog"
)

var lambdaStart = lambda.Start

// FunctionProxy represents a proxy that can be used to wrap the invocation of a function, such as a Lambda.
type FunctionProxy interface {
	// Start starts the tunnel. This function never returns unless Stop is called by another goroutine.
	Start()
}

type lambdaProxy struct {
	applicationId string
	client        BastionClient
	handler       lambda.Handler
}

// NewLambdaFunctionProxy creates a new FunctionProxy for AWS Lambda operations.
// The handler is the handler that will be invoked when a request is received.
// It is subject to the same restrictions as the handler for the underlying serverless function provider (such as lambda.Start).
func NewLambdaFunctionProxy(
	applicationId string,
	client BastionClient,
	handler interface{},
) FunctionProxy {
	lambdaHandler := lambda.NewHandler(handler)

	return &lambdaProxy{
		applicationId: applicationId,
		client:        client,
		handler:       lambdaHandler,
	}
}

func (p *lambdaProxy) Start() {
	handler := p.lambdaHandler()
	slog.Info("starting lambda proxy")
	lambdaStart(handler)
}

// lambdaHandler is the Lambda-level wrapper around the handler.
// It is responsible for publishing the message to the tunnel, and waiting for a response.
func (p *lambdaProxy) lambdaHandler() lambda.Handler {
	wrapper := func(ctx context.Context, payload *json.RawMessage) (*json.RawMessage, error) {
		slog.Debug("publishing message to tunnel", "message", payload)

		// Raw constant to avoid cycles -- this needs to be moved.
		forwardPayload := messages.NewForwardRequestPayload(*payload)
		message := funcie.NewMessageWithPayload(p.applicationId, "FORWARD_REQUEST", forwardPayload)

		marshaled, err := funcie.MarshalMessagePayload(*message)
		if err != nil {
			return nil, fmt.Errorf("marshalling message payload: %w", err)
		}

		resp, err := p.client.SendRequest(ctx, marshaled)
		if err != nil {
			// If we can't reach the bastion, we should just handle the request directly.
			slog.ErrorCtx(ctx, "failed to send request to bastion", err, "message", message)
			return p.handleDirect(ctx, payload)
		}

		forwardResponse, err := funcie.UnmarshalResponsePayload[messages.ForwardRequestResponse](resp)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling response from bastion: %w", err)
		}

		if forwardResponse.Error != nil {
			if forwardResponse.Error.Error() == funcie.ErrNoActiveConsumer.Error() {
				slog.InfoCtx(ctx, "no active consumer for request", "message", message)
				return p.handleDirect(ctx, payload)
			}
			slog.DebugCtx(ctx, "received error from bastion", "error", forwardResponse.Error)
			return nil, fmt.Errorf("received error from proxied implementation: %w", forwardResponse.Error)
		}

		slog.DebugCtx(ctx, "received response from bastion", "response", string(forwardResponse.Data.Body))

		return &forwardResponse.Data.Body, nil
	}
	return lambda.NewHandler(wrapper)
}

func (p *lambdaProxy) handleDirect(ctx context.Context, payload *json.RawMessage) (*json.RawMessage, error) {
	res, err := p.handler.Invoke(ctx, *payload)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke handler: %w", err)
	}

	raw := json.RawMessage(res)
	return &raw, nil
}
