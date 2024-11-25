package publisher

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/google/uuid"
	"log/slog"
)

type messageHandler struct {
	connStore ConnectionStore
	logger    *slog.Logger
}

// NewMessageProcessor returns a new MessageHandler for server-side connections.
func NewMessageProcessor(connStore ConnectionStore, logger *slog.Logger) websocket.MessageProcessor {
	return &messageHandler{
		connStore: connStore,
		logger:    logger,
	}
}

func (m *messageHandler) ProcessMessage(ctx context.Context, conn websocket.Connection, msg *funcie.Message) (*funcie.Response, error) {
	switch msg.Kind {
	case messages.MessageKindRegister:
		registerMessage, err := funcie.UnmarshalMessagePayload[messages.RegistrationMessage](msg)
		if err != nil {
			return nil, fmt.Errorf("unmarshal payload %v: %w", msg.Payload, err)
		}

		m.connStore.RegisterConnection(registerMessage.Application, conn)
		m.logger.InfoContext(ctx, "Registered connection", "application", msg.Application)

		registerPayload := messages.NewRegistrationResponsePayload(uuid.New())
		resp := funcie.NewResponseWithPayload(msg.ID, registerPayload, nil)
		untyped, err := funcie.MarshalResponsePayload(resp)
		if err != nil {
			return nil, fmt.Errorf("marshalling registratoin response payload: %w", err)
		}

		return untyped, nil
	case messages.MessageKindDeregister:
		_, err := m.connStore.UnregisterConnection(msg.Application)
		if err != nil {
			return nil, fmt.Errorf("unregistering connection: %w", err)
		}

		deregisterPayload := messages.NewDeregistrationResponsePayload()
		resp := funcie.NewResponseWithPayload(msg.ID, deregisterPayload, nil)
		untyped, err := funcie.MarshalResponsePayload(resp)
		if err != nil {
			return nil, fmt.Errorf("marshalling deregistration response payload: %w", err)
		}

		m.logger.InfoContext(ctx, "Unregistered connection", "application", msg.Application)
		return untyped, nil
	default:
		return nil, fmt.Errorf("invalid client to server message type: %v", msg.Kind)
	}
}
