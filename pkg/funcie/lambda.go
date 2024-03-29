package funcie

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log/slog"
)

type Tunnel interface {
	// Start starts the tunnel. This function never returns.
	// The handler is the handler that will be invoked when a request is received.
	// It is subject to the same restrictions as the handler for the serverless function provider (such as lambda.Start).
	Start()
}

type lambdaTunnel struct {
	applicationId string
	handler       interface{}
	publisher     Publisher
	consumer      Consumer
}

func NewLambdaTunnel(applicationId string, handler interface{}, publisher Publisher, consumer Consumer) Tunnel {
	return &lambdaTunnel{
		applicationId: applicationId,
		handler:       handler,
		publisher:     publisher,
		consumer:      consumer,
	}
}

// Start replaces the built-in lambda.Start function.
// If the current execution environment is Lambda, it will proxy requests to the
// tunnel instead of directly to the Lambda API (if a consumer is active).
// If the current execution environment is not Lambda, it will wait for proxied requests from the tunnel.
func (t *lambdaTunnel) Start() {
	// First, check if we're running inside a Lambda.
	// If we are, then we need to wrap the handler in a LambdaHandler.
	var wrappedHandler lambda.Handler
	if IsRunningWithLambda() {
		slog.Info("running in Lambda", "applicationId", t.applicationId)
		wrappedHandler = t.lambdaHandler()
		slog.Info("starting lambda")
		lambda.Start(wrappedHandler)
	} else {
		slog.Info("running outside of Lambda", "applicationId", t.applicationId)
		t.beginProxyConsume()

		wrappedHandler = lambda.NewHandler(t.handler)
	}
}

// proxyHandler is the consumer-level wrapper around the handler.
// It is responsible for receiving messages from the tunnel, invoking the handler, and returning the response.
func (t *lambdaTunnel) beginProxyConsume() {
	//TODO -- FIX ME -- Doesn't yet follow the new Consume/Subscribe pattern.

	//slog.Info("creating proxy handler")
	//localHandler := lambda.NewHandler(t.handler)
	//slog.Info("created proxy handler")
	//
	//err := t.consumer.Consume(context.Background(), func(ctx context.Context, message *Message) (*Response, error) {
	//	slog.Debug("received message from tunnel", "message", string(message.Payload), "applicationId", t.applicationId)
	//
	//	// Invoke the handler.
	//	resp, err := localHandler.Invoke(ctx, message.Payload)
	//	response := NewResponse(message.ID, resp, err)
	//
	//	// Publish the response to the tunnel.
	//	slog.Debug("returning response to tunnel", "response", response, "err", err)
	//
	//	return response, nil
	//})
	//if err != nil {
	//	slog.Error("failed to consume from tunnel", err)
	//	// TODO: Figure out how to handle this error.
	//	os.Exit(1)
	//}
}

// lambdaHandler is the Lambda-level wrapper around the handler.
// It is responsible for publishing the message to the tunnel, and waiting for a response.
func (t *lambdaTunnel) lambdaHandler() lambda.Handler {
	slog.Info("creating lambda handler")
	awsHandler := lambda.NewHandler(t.handler)
	slog.Info("created lambda handler")

	wrapper := func(ctx context.Context, message *json.RawMessage) (*json.RawMessage, error) {
		slog.Debug("publishing message to tunnel", "payload", message)

		bytes, err := message.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}

		var rawResp *json.RawMessage

		// Raw constant to avoid cycles -- this needs to be moved.
		msg := NewMessage(t.applicationId, "FORWARD_REQUEST", bytes)
		res, err := t.publisher.Publish(ctx, msg)
		if err == nil {
			// If we got a response, then we can return it immediately.
			slog.Debug("received response from tunnel", "response", res.Data, "err", err)
			if res.Error != nil {
				return nil, fmt.Errorf("received error from proxied implementation: %w", res.Error)
			}

			if res.Data != nil {
				rawResp = res.Data
				if err != nil {
					return nil, fmt.Errorf("failed to marshal response from proxied implementation: %w", err)
				}
			}

			return rawResp, nil
		}
		if err != ErrNoActiveConsumer {
			return nil, fmt.Errorf("failed to publish message: %w", err)
		}

		implRes, err := awsHandler.Invoke(ctx, bytes)
		slog.Debug("received response from direct implementation", "response", string(implRes), "err", err)
		if err != nil {
			return nil, fmt.Errorf("failed to invoke handler: %w", err)
		}

		if len(implRes) > 0 {
			rawResp = &json.RawMessage{}
			err = rawResp.UnmarshalJSON(implRes)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response from direct implementation: %w", err)
			}
		}

		return rawResp, nil
	}
	slog.Info("returning lambda handler")
	return lambda.NewHandler(wrapper)
}
