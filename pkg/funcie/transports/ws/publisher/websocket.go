package publisher

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"nhooyr.io/websocket"
)

type WebsocketClient struct {
	conn Websocket
}

func NewWebsocketClient(conn Websocket) *WebsocketClient {
	return &WebsocketClient{conn: conn}
}

func (c *WebsocketClient) HandleMessage(ctx context.Context, msg funcie.Message) error {
	//todo
	return nil
}

func (c *WebsocketClient) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "closing")
}
