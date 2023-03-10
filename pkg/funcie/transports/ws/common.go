package ws

import (
	"context"
	"io"
	ws "nhooyr.io/websocket"
)

type websocket interface {
	Close(code ws.StatusCode, reason string) error
	Read(ctx context.Context) (ws.MessageType, []byte, error)
	Write(ctx context.Context, typ ws.MessageType, p []byte) error
	Writer(ctx context.Context, typ ws.MessageType) (io.WriteCloser, error)
	Reader(ctx context.Context) (ws.MessageType, io.Reader, error)
}

type ClientToServerMessage struct {
	Channel     string `json:"channel"`
	RequestType string `json:"request_type"`
}

type ServerToClientMessage struct {
	Channel     string `json:"channel"`
	Payload     string `json:"payload"`
	RequestType string `json:"request_type"`
}

const ClientToServerMessageRequestTypeSubscribe = "subscribe"
const ServerToClientMessageRequestTypeInvoke = "invoke"
