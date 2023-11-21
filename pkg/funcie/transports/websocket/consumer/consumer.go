package consumer

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"log/slog"
	ws "nhooyr.io/websocket"
)

type consumer struct {
	serverUrl string
	router    utils.ClientHandlerRouter
	client    Client
	logger    *slog.Logger

	conn      websocket.Connection
	connected bool
}

// NewConsumer creates a new consumer that consumes messages from the given URL.
func NewConsumer(serverUrl string, router utils.ClientHandlerRouter, client Client, logger *slog.Logger) funcie.Consumer {
	return &consumer{
		serverUrl: serverUrl,
		router:    router,
		client:    client,
		logger:    logger,
	}
}

func (c *consumer) Connect(ctx context.Context) error {
	conn, err := c.client.Dial(ctx, c.serverUrl)
	if err != nil {
		return fmt.Errorf("dialing %v: %w", c.serverUrl, err)
	}

	c.conn = conn
	c.connected = true
	return nil
}

func (c *consumer) Consume(ctx context.Context) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	defer func() {
		c.connected = false
		if err := c.conn.Close(ws.StatusAbnormalClosure, "exiting consumer"); err != nil {
			c.logger.WarnContext(ctx, "error closing Websocket", "error", err)
		}
	}()

	for ctx.Err() == nil {
		var msg funcie.Message
		err := c.conn.Read(ctx, &msg)
		if err != nil {
			return fmt.Errorf("reading message: %w", err)
		}

		c.logger.DebugContext(ctx, "received message", "message", msg)

		resp, err := c.router.Handle(ctx, &msg)
		if err != nil {
			return fmt.Errorf("handling message: %w", err)
		}

		if resp != nil {
			err = c.conn.Write(ctx, resp)
			if err != nil {
				return fmt.Errorf("writing response: %w", err)
			}
		}
	}

	return nil
}

func (c *consumer) Subscribe(ctx context.Context, applicationId string, handler funcie.Handler) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	if err := c.router.AddClientHandler(applicationId, handler); err != nil {
		return fmt.Errorf("adding client handler: %w", err)
	}

	return nil
}

func (c *consumer) Unsubscribe(ctx context.Context, applicationId string) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	if err := c.router.RemoveClientHandler(applicationId); err != nil {
		return fmt.Errorf("removing client handler: %w", err)
	}

	return nil
}
