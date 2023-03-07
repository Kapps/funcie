package main

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/cmd/funcie-lambda-ext/lambdaext"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	extensionClient = lambdaext.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	extensionName   = "funcie-lambda-ext"
	redisAddr       = os.Getenv("FUNCIE_REDIS_ADDR")
)

func main() {
	if redisAddr == "" {
		panic("FUNCIE_REDIS_ADDR not set")
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		slog.Info("Received signal; exiting", "signal", s)
		cancel()
	}()

	res, err := extensionClient.Register(ctx, extensionName)
	if err != nil {
		panic(err)
	}

	slog.Info("Register response:", "response", res)

	// Will block until shutdown event is received or cancelled via the context.
	if err := processEvents(ctx); err != nil {
		slog.Info("Error; exiting", "error", err)
		panic(err)
	}
}

// Method to process events
func processEvents(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			slog.Debug("Waiting for event...")
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				return fmt.Errorf("error getting next event: %w", err)
			}

			// Exit if we receive a SHUTDOWN event
			if res.EventType == lambdaext.Shutdown {
				slog.Info("Received SHUTDOWN event; exiting")
				return nil
			}

			if res.EventType == lambdaext.Invoke {
				slog.Debug("Received INVOKE event; invoking consumer")

			}
		}
	}
}
