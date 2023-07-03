package bastion

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
)

// MessageProcessor dispatches messages to the appropriate handler.
type MessageProcessor interface {
	// ProcessMessage processes the given message and returns a response.
	ProcessMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error)
}

type messageProcessor struct {
	handler Handler
}

// NewMessageProcessor creates a new MessageProcessor.
func NewMessageProcessor(handler Handler) MessageProcessor {
	return &messageProcessor{
		handler: handler,
	}
}

func (p *messageProcessor) ProcessMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	switch message.Kind {
	case funcie.MessageKindDispatch:
		return p.handler.ForwardRequest(ctx, message)
	case messages.MessageKindRegister:
		return p.register(ctx, message)
	case funcie.MessageKindResponse:
		return p.handler.Response(message)
	default:
		return nil, ErrUnknownMessageKind
	}
}

func (p *messageProcessor) register(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	var payload messages.RegistrationRequestPayload
	err := json.Unmarshal(message.Payload, &payload)
	if err != nil {
		return nil, fmt.Errorf("unmarshal payload %v: %w", message.Payload, err)
	}

	err = p.handler.Register(ctx, funcie.NewApplication(payload.Name, payload.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("register application %v: %w", payload.Name, err)
	}

	responsePayload := messages.NewRegistrationResponsePayload()
	return funcie.NewResponseWithPayload(message.ID, nil, nil)
}
