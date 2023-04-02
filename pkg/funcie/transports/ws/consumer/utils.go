package consumer

import (
	"context"
	"net/http"
	ws "nhooyr.io/websocket"
)

type Websocket interface {
	Close(code ws.StatusCode, reason string) error
	Read(ctx context.Context) (ws.MessageType, []byte, error)
	Write(ctx context.Context, typ ws.MessageType, p []byte) error
}

type WebsocketClient interface {
	Dial(ctx context.Context, u string, opts *ws.DialOptions) (Websocket, *http.Response, error)
}

type WebsocketClientWrapper struct{}

func (w *WebsocketClientWrapper) Dial(ctx context.Context, u string, opts *ws.DialOptions) (Websocket, *http.Response, error) {
	return ws.Dial(ctx, u, opts)
}
