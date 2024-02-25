package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"log/slog"
	nhws "nhooyr.io/websocket"
)

// An Exchange is a service that can be used to send and receive messages for many connections.
type Exchange interface {
	// RegisterConnection registers a connection with the exchange.
	// This method transfers ownership of the connection to the exchange.
	RegisterConnection(ctx context.Context, conn Connection) error
	// Send sends a message to the given connection that is registered in this exchange.
	Send(ctx context.Context, conn Connection, message *funcie.Message) (*funcie.Response, error)
}

type exchange struct {
	processor        MessageProcessor
	responseNotifier ResponseNotifier
	logger           *slog.Logger
}

// NewExchange returns a default in-memory implementation of an Exchange.
func NewExchange(
	responseNotifier ResponseNotifier,
	processor MessageProcessor,
	logger *slog.Logger,
) Exchange {
	return &exchange{
		responseNotifier: responseNotifier,
		processor:        processor,
		logger:           logger,
	}
}

func (x *exchange) RegisterConnection(ctx context.Context, conn Connection) error {
	go x.readLoop(ctx, conn)

	return nil
}

func (x *exchange) Send(ctx context.Context, conn Connection, message *funcie.Message) (*funcie.Response, error) {

	// TODO: Validate that the connection is registered with this exchange.

	payload, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("marshalling message: %w", err)
	}

	jsonPayload := json.RawMessage(payload)

	x.logger.Debug("Sending message", "messageId", message.ID)

	err = conn.Write(ctx, &Envelope{
		Kind: PayloadKindRequest,
		Data: &jsonPayload,
	})
	if err != nil {
		return nil, fmt.Errorf("writing message: %w", err)
	}

	// TODO: If the connection closes, we should stop waiting early.
	// TODO: Should we have a timeout at this level, or parent?

	resp, err := x.responseNotifier.WaitForResponse(ctx, message.ID)
	if err != nil {
		return nil, fmt.Errorf("waiting for response: %w", err)
	}

	x.logger.Debug("Received response", "messageId", message.ID)

	return resp, nil
}

func (x *exchange) readLoop(ctx context.Context, conn Connection) {
	closeConn := func(reason string) {
		if err := conn.Close(nhws.StatusNormalClosure, reason); err != nil {
			x.logger.ErrorContext(ctx, "Error closing connection", "error", err)
		}
	}

	for {
		x.logger.DebugContext(ctx, "Waiting for next message")

		envelope, err := x.readNextMessage(ctx, conn)
		if err != nil {
			x.logger.ErrorContext(ctx, "Error reading message; closing connection", "error", err)
			closeConn("error reading message")
			break
		}

		err = x.processMessage(ctx, envelope, conn)
		if err != nil {
			// TODO: Should we send an error response?
			// TODO: Should we close the connection just because processing fails?
			x.logger.ErrorContext(ctx, "Error processing message; closing connection", "error", err)
			closeConn("error processing message")
			break
		}
	}
}

func (x *exchange) readNextMessage(ctx context.Context, conn Connection) (*Envelope, error) {
	var msg Envelope
	err := conn.Read(ctx, &msg)
	if err != nil {
		return nil, fmt.Errorf("reading message: %w", err)
	}

	return &msg, nil
}

func (x *exchange) processMessage(ctx context.Context, envelope *Envelope, conn Connection) error {
	switch envelope.Kind {
	case PayloadKindRequest:
		return x.processRequest(ctx, envelope, conn)
	case PayloadKindResponse:
		return x.processResponse(ctx, envelope)
	default:
		return fmt.Errorf("invalid message type: %v", envelope.Kind)
	}

	return nil
}

func (x *exchange) processRequest(ctx context.Context, envelope *Envelope, conn Connection) error {
	var msg *funcie.Message
	err := json.Unmarshal(*envelope.Data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	x.logger.DebugContext(ctx, "Received request", "message", msg)

	resp, err := x.processor.ProcessMessage(ctx, conn, msg)
	if err != nil {
		return fmt.Errorf("handling request: %w", err)
	}

	responsePayload, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshalling response: %w", err)
	}

	jsonPayload := json.RawMessage(responsePayload)

	err = conn.Write(ctx, &Envelope{
		Kind: PayloadKindResponse,
		Data: &jsonPayload,
	})
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

func (x *exchange) processResponse(ctx context.Context, envelope *Envelope) error {
	var msg *funcie.Response
	err := json.Unmarshal(*envelope.Data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	x.logger.DebugContext(ctx, "Received response", "message", msg)

	x.responseNotifier.Notify(ctx, msg)

	return nil
}
