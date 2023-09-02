package bastion

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports"
	"golang.org/x/exp/slog"
)

type requestHandler struct {
	publisher funcie.Publisher
}

// NewRequestHandler creates a new RequestHandler.
func NewRequestHandler(publisher funcie.Publisher) transports.MessageHandler {
	return &requestHandler{
		publisher: publisher,
	}
}

func (r *requestHandler) Register(ctx context.Context, message messages.RegistrationMessage) (*messages.RegistrationResponse, error) {
	return nil, fmt.Errorf("register unsupported")
}

func (r *requestHandler) Deregister(ctx context.Context, message messages.DeregistrationMessage) (*messages.DeregistrationResponse, error) {
	return nil, fmt.Errorf("deregister unsupported")
}

func (r *requestHandler) ForwardRequest(ctx context.Context, message messages.ForwardRequestMessage) (*messages.ForwardRequestResponse, error) {
	slog.DebugCtx(ctx, "forwarding request", "message", &message)

	marshaled, err := funcie.MarshalMessagePayload(message)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := r.publisher.Publish(ctx, marshaled)
	if err != nil {
		return nil, fmt.Errorf("publish request: %w", err)
	}

	unmarshaledResp, err := funcie.UnmarshalResponsePayload[messages.ForwardRequestResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return unmarshaledResp, nil
}
