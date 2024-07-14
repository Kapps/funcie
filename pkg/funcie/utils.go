package funcie

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
)

// MustSerialize serializes the given value to JSON, or panics if it fails.
func MustSerialize(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// MustDeserialize deserializes the given JSON to the given type, or panics if it fails.
func MustDeserialize[T any](b []byte) T {
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	return v
}

type Closable interface {
	Close() error
}

// CloseOrLog closes the given closable, logging any errors (but continuing execution).
func CloseOrLog(name string, c Closable) {
	if err := c.Close(); err != nil {
		slog.Error("error closing", "error", err)
	} else {
		slog.Debug("closed resource", "resource", name)
	}
}

// IsRunningWithLambda returns true if the current process is running in AWS Lambda.
func IsRunningWithLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

// ConfigureLogging configures slog to log to stdout at the given level.
// If FUNCIE_LOG_LEVEL is set, it will be used as the log level.
// Otherwise, the default is Info.
func ConfigureLogging() {
	programLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{
		//AddSource: true,
		Level: logLevelFromEnv(),
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
	programLevel.Set(slog.LevelDebug)
}

func logLevelFromEnv() slog.Level {
	logLevel := os.Getenv("FUNCIE_LOG_LEVEL")
	if logLevel == "" {
		return slog.LevelInfo
	}

	switch strings.ToLower(logLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Warn("unknown log level", "level", logLevel)
		return slog.LevelInfo
	}
}
