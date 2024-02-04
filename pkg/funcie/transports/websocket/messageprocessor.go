package websocket

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
)

// MessageProcessor is responsible for processing messages from a websocket endpoint.
type MessageProcessor interface {
	// ProcessMessage processes an untyped message and returns an untyped response.
	ProcessMessage(ctx context.Context, conn Connection, msg *funcie.Message) (*funcie.Response, error)
}
