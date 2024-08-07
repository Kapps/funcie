package funcietunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/aws/aws-lambda-go/lambda"
	"log/slog"
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
	logger        *slog.Logger
}

// NewLambdaFunctionProxy creates a new FunctionProxy for AWS Lambda operations.
// The handler is the handler that will be invoked when a request is received.
// It is subject to the same restrictions as the handler for the underlying serverless function provider (such as lambda.Start).
func NewLambdaFunctionProxy(
	applicationId string,
	client BastionClient,
	handler interface{},
	logger *slog.Logger,
) FunctionProxy {
	lambdaHandler := lambda.NewHandler(handler)

	return &lambdaProxy{
		applicationId: applicationId,
		client:        client,
		handler:       lambdaHandler,
		logger:        logger,
	}
}

func (p *lambdaProxy) Start() {
	handler := p.lambdaHandler()
	p.logger.Info("starting lambda proxy")
	lambdaStart(handler)
}

// lambdaHandler is the Lambda-level wrapper around the handler.
// It is responsible for publishing the message to the tunnel, and waiting for a response.
func (p *lambdaProxy) lambdaHandler() lambda.Handler {
	wrapper := func(ctx context.Context, payload *json.RawMessage) (*json.RawMessage, error) {
		p.logger.DebugContext(ctx, "publishing message to tunnel", "message", string(*payload))

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
			p.logger.WarnContext(ctx, "failed to send request to bastion", "error", err, "messageId", message.ID)
			p.logger.DebugContext(ctx, "failed delivery details", "message", message)
			return p.handleDirect(ctx, payload)
		}

		forwardResponse, err := funcie.UnmarshalResponsePayload[messages.ForwardRequestResponse](resp)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling response from bastion: %w", err)
		}

		if forwardResponse.Error != nil {
			// This is a bit of a gross way to check this, but... it is what it is.
			// We need to add error codes in the future and make this less gross.
			if isExpectedProxyError(forwardResponse.Error) {
				// If there is no active consumer, we should just handle the request directly.
				p.logger.DebugContext(ctx, "no active consumer for request", "message", message)
				return p.handleDirect(ctx, payload)
			}
			// In this case though, the request was handled and the handling returned an error.
			// So we should forward that error back to the Lambda.
			p.logger.DebugContext(ctx, "received error from bastion", "error", forwardResponse.Error)
			return nil, fmt.Errorf("received error from proxied implementation: %w", forwardResponse.Error)
		}

		p.logger.DebugContext(ctx, "received response from bastion", "response", string(forwardResponse.Data.Body))

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

func isExpectedProxyError(err *funcie.ProxyError) bool {
	str := err.Error()
	return str == funcie.ErrNoActiveConsumer.Error() || str == funcie.ErrApplicationNotFound.Error()
}
