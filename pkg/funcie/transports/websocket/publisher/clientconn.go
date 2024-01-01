package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"log/slog"
)

// While this is a simple wrapper, it makes testing quite a bit easier.

type RequestHandler = func(ctx context.Context, conn ClientConnection, msg *funcie.Message) (*funcie.Response, error)

// ClientConnection wraps a websocket connection and manages sending and receiving data.
// This type is aware of the funcie protocol and will automatically handle responses.
type ClientConnection interface {
	websocket.Connection

	// Send sends a message to the client, and waits for a response.
	Send(ctx context.Context, message *funcie.Message) (*funcie.Response, error)
}

type clientConn struct {
	websocket.Connection
	logger           *slog.Logger
	requestHandler   RequestHandler
	responseNotifier ResponseNotifier
}

// NewClientConnection creates a new ClientConnection from a websocket connection.
// The context will be used for all operations on the connection.
func NewClientConnection(
	ctx context.Context,
	conn websocket.Connection,
	requestHandler RequestHandler,
	responseNotifier ResponseNotifier,
	logger *slog.Logger,
) ClientConnection {
	c := &clientConn{
		Connection:       conn,
		logger:           logger,
		requestHandler:   requestHandler,
		responseNotifier: responseNotifier,
	}

	go c.readLoop(ctx)

	return c
}

func (c *clientConn) Send(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	err := c.Write(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("writing message: %w", err)
	}

	c.logger.DebugContext(ctx, "Sent request", "message", message.ID)

	resp, err := c.responseNotifier.WaitForResponse(ctx, message.ID)
	if err != nil {
		return nil, fmt.Errorf("waiting for response: %w", err)
	}

	return resp, nil
}

func (c *clientConn) readLoop(ctx context.Context) {
	for {
		c.logger.DebugContext(ctx, "Waiting for next message")

		envelope, err := c.readNextMessage(ctx)
		if err != nil {
			c.logger.ErrorContext(ctx, "Error reading message", "error", err)
			continue
		}

		err = c.processMessage(ctx, envelope)
		if err != nil {
			c.logger.ErrorContext(ctx, "Error processing message", "error", err)
			continue
		}
	}
}

func (c *clientConn) readNextMessage(ctx context.Context) (*websocket.Envelope, error) {
	var msg websocket.Envelope
	err := c.Read(ctx, &msg)
	if err != nil {
		return nil, fmt.Errorf("reading message: %w", err)
	}

	return &msg, nil
}

func (c *clientConn) processMessage(ctx context.Context, envelope *websocket.Envelope) error {
	switch envelope.Kind {
	case websocket.PayloadKindRequest:
		return c.processRequest(ctx, envelope)
	case websocket.PayloadKindResponse:
		return c.processResponse(ctx, envelope)
	default:
		return fmt.Errorf("invalid message type: %v", envelope.Kind)
	}

	return nil
}

func (c *clientConn) processRequest(ctx context.Context, envelope *websocket.Envelope) error {
	var msg *funcie.Message
	err := json.Unmarshal(*envelope.Data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	c.logger.DebugContext(ctx, "Received request", "message", msg)

	resp, err := c.requestHandler(ctx, c, msg)
	if err != nil {
		// TODO: Should we send an error response?
		return fmt.Errorf("handling request: %w", err)
	}

	if resp == nil {
		return nil
	}

	responsePayload, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshalling response: %w", err)
	}

	jsonPayload := json.RawMessage(responsePayload)

	err = c.Write(ctx, &websocket.Envelope{
		Kind: websocket.PayloadKindResponse,
		Data: &jsonPayload,
	})
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

func (c *clientConn) processResponse(ctx context.Context, envelope *websocket.Envelope) error {
	var msg *funcie.Response
	err := json.Unmarshal(*envelope.Data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	c.logger.DebugContext(ctx, "Received response", "message", msg)

	c.responseNotifier.Notify(ctx, msg)

	return nil
}
