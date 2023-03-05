package bastion

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
)

// RequestHandler allows the handling of incoming valid Bastion requests.
type RequestHandler interface {
	// Dispatch dispatches the given request to the proxy handler.
	// The response is sent back to the caller.
	Dispatch(ctx context.Context, request *Request) (*funcie.Response, error)
}

type requestHandler struct {
	config    *Config
	publisher funcie.Publisher
}

// NewRequestHandler creates a new RequestHandler.
func NewRequestHandler(config *Config, publisher funcie.Publisher) RequestHandler {
	return &requestHandler{
		config:    config,
		publisher: publisher,
	}
}

func (h *requestHandler) Dispatch(ctx context.Context, request *Request) (*funcie.Response, error) {
	contents, err := request.Payload.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	message := funcie.NewMessage(contents, h.config.RequestTtl)

	resp, err := h.publisher.Publish(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	return resp, nil
}
