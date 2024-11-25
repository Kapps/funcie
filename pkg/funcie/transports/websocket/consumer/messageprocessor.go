package consumer

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
)

type messageProcessor struct {
	router utils.ClientHandlerRouter
}

// NewMessageProcessor creates a new message processor for consumer messages.
func NewMessageProcessor(router utils.ClientHandlerRouter) websocket.MessageProcessor {
	return &messageProcessor{
		router: router,
	}
}

func (m *messageProcessor) ProcessMessage(ctx context.Context, conn websocket.Connection, msg *funcie.Message) (*funcie.Response, error) {
	if msg.Kind != messages.MessageKindForwardRequest {
		return nil, fmt.Errorf("unsupported message kind: %v", msg.Kind)
	}

	resp, err := m.router.Handle(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("handling message: %w", err)
	}

	return resp, nil
}
