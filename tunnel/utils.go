package tunnel

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slog"
)

// GetResponseKeyForMessage returns the Redis key for the response of a message
func GetResponseKeyForMessage(messageId string) string {
	if messageId == "" {
		panic("messageId cannot be empty")
	}
	return fmt.Sprintf("funcie:resp:%v", messageId)
}

// MustSerialize serializes the given value to JSON, or panics if it fails.
func MustSerialize(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

type Closable interface {
	Close() error
}

// CloseOrLog closes the given closable, logging any errors (but continuing execution).
func CloseOrLog(name string, c Closable) {
	if err := c.Close(); err != nil {
		slog.Error("error closing", err)
	} else {
		slog.Debug("closed resource", "resource", name)
	}
}
