package publisher

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"nhooyr.io/websocket"
)

type WebsocketClientConnection struct {
	conn Websocket
}

func NewWebsocketClient(conn Websocket) *WebsocketClientConnection {
	return &WebsocketClientConnection{conn: conn}
}

func (c *WebsocketClientConnection) HandleMessage(ctx context.Context, msg messages.Message) error {
	//todo
	return nil
}

func (c *WebsocketClientConnection) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "closing")
}
