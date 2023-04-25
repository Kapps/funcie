package funcie

import (
	"encoding/json"
	"golang.org/x/exp/slog"
	"os"
)

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

// IsRunningWithLambda returns true if the current process is running in AWS Lambda.
func IsRunningWithLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
