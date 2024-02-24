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
	dialer    Dialer
	exchange  websocket.Exchange
	logger    *slog.Logger
	router    utils.ClientHandlerRouter

	conn      websocket.Connection
	connected bool
}

// NewConsumer creates a new consumer that consumes messages from the given URL.
func NewConsumer(serverUrl string, exchange websocket.Exchange, router utils.ClientHandlerRouter, dialer Dialer, logger *slog.Logger) funcie.Consumer {
	return &consumer{
		serverUrl: serverUrl,
		dialer:    dialer,
		logger:    logger,
		router:    router,
		exchange:  exchange,
	}
}

func (c *consumer) Connect(ctx context.Context) error {
	conn, err := c.dialer.Dial(ctx, c.serverUrl)
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

	if err := c.exchange.RegisterConnection(ctx, c.conn); err != nil {
		// If we fail to register the connection, we should close it.
		if err := c.conn.Close(ws.StatusAbnormalClosure, "failed to register connection with exchange"); err != nil {
			c.logger.WarnContext(ctx, "error closing Websocket", "error", err)
		}
		return fmt.Errorf("registering connection with excahnge: %w", err)
	}

	defer func() {
		c.connected = false
		if err := c.conn.Close(ws.StatusAbnormalClosure, "exiting consumer"); err != nil {
			c.logger.WarnContext(ctx, "error closing Websocket", "error", err)
		}
	}()

	<-ctx.Done()

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
