package publisher

import (
	"context"
	"net/http"
	ws "nhooyr.io/websocket"
)

type Websocket interface {
	Close(code ws.StatusCode, reason string) error
	Read(ctx context.Context) (ws.MessageType, []byte, error)
	Write(ctx context.Context, typ ws.MessageType, p []byte) error
	Subprotocol() string
}

type WebsocketServer interface {
	Accept(w http.ResponseWriter, r *http.Request, opts *ws.AcceptOptions) (Websocket, error)
}

type WebsocketServerWrapper struct{}

func (w *WebsocketServerWrapper) Accept(rw http.ResponseWriter, r *http.Request, opts *ws.AcceptOptions) (Websocket, error) {
	return ws.Accept(rw, r, opts)
}
