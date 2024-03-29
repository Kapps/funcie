package transports

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
)

// TODO: Refactor -- MessageHandler vs MessageProcessor is confusing.

// MessageProcessor dispatches messages to the appropriate handler.
type MessageProcessor interface {
	// ProcessMessage processes the given message and returns a response.
	ProcessMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error)
}

var ErrUnknownMessageKind = fmt.Errorf("unknown message kind")

type messageProcessor struct {
	handler MessageHandler
}

// NewMessageProcessor creates a new MessageProcessor.
func NewMessageProcessor(handler MessageHandler) MessageProcessor {
	return &messageProcessor{
		handler: handler,
	}
}

func (p *messageProcessor) ProcessMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	switch message.Kind {
	case messages.MessageKindForwardRequest:
		// Usually comes from consumer
		return p.forwardRequest(ctx, message)
	case messages.MessageKindRegister:
		// Usually comes from host
		return p.register(ctx, message)
	case messages.MessageKindDeregister:
		// Usually comes from host
		return p.deregister(ctx, message)
	default:
		return nil, ErrUnknownMessageKind
	}
}

func (p *messageProcessor) forwardRequest(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	forwardMessage, err := funcie.UnmarshalMessagePayload[messages.ForwardRequestMessage](message)
	if err != nil {
		return nil, fmt.Errorf("unmarshal payload %v: %w", message.Payload, err)
	}

	resp, err := p.handler.ForwardRequest(ctx, *forwardMessage)
	if err != nil {
		return nil, fmt.Errorf("forward request %v: %w", forwardMessage.ID, err)
	}
	serializedResponse, err := funcie.MarshalResponsePayload(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response %v: %w", resp, err)
	}
	return serializedResponse, nil
}

func (p *messageProcessor) register(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	registerMessage, err := funcie.UnmarshalMessagePayload[messages.RegistrationMessage](message)
	if err != nil {
		return nil, fmt.Errorf("unmarshal payload %v: %w", message.Payload, err)
	}
	resp, err := p.handler.Register(ctx, *registerMessage)
	if err != nil {
		return nil, fmt.Errorf("register application %v: %w", message.Application, err)
	}
	serializedResponse, err := funcie.MarshalResponsePayload(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response %v: %w", resp, err)
	}
	return serializedResponse, nil
}

func (p *messageProcessor) deregister(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	deregisterMessage, err := funcie.UnmarshalMessagePayload[messages.DeregistrationMessage](message)
	if err != nil {
		return nil, fmt.Errorf("unmarshal payload %v: %w", message.Payload, err)
	}
	resp, err := p.handler.Deregister(ctx, *deregisterMessage)
	if err != nil {
		return nil, fmt.Errorf("deregister application %v: %w", resp.Data, err)
	}
	serializedResponse, err := funcie.MarshalResponsePayload(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response %v: %w", resp, err)
	}
	return serializedResponse, nil
}
