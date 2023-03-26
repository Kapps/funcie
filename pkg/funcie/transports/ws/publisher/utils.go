package publisher

import (
	"context"
	ws "nhooyr.io/websocket"
)

type Websocket interface {
	Close(code ws.StatusCode, reason string) error
	Read(ctx context.Context) (ws.MessageType, []byte, error)
	Write(ctx context.Context, typ ws.MessageType, p []byte) error
}
