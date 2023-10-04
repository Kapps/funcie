package main

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/cmd/funcie-lambda-ext/lambdaext"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	extensionName = "funcie-lambda-ext"
)

var (
	runtimeApi = os.Getenv("AWS_LAMBDA_RUNTIME_API")
)

func main() {
	isLocal := runtimeApi == ""
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		slog.Info("Received signal; exiting", "signal", s)
	}()

	extensionClient := lambdaext.NewClient(runtimeApi)
	if !isLocal {
		registerExtension(ctx, extensionClient)
	} else {
		slog.Info("Running locally; not registering extension")
	}

	/*bastionServer := createBastion()
	go func() {
		if err := bastionServer.Listen(); err != nil {
			slog.Info("Error; exiting", "error", err)
			panic(err)
		}
	}()*/

	if !isLocal {
		// Will block until shutdown event is received or cancelled via the context.
		if err := processEvents(ctx, extensionClient); err != nil {
			slog.Info("Error; exiting", "error", err)
			panic(err)
		}
	} else {
		slog.Info("Running locally; not processing events")
		<-ctx.Done()
	}
}

/*func createBastion() bastion2.Server {
	config := bastion2.NewConfigFromEnvironment()
	redisClient := redis.NewClient(&redis.Options{
		Addr:       config.RedisAddress,
		ClientName: "Funcie Bastion",
	})
	publisher := r.NewPublisher(redisClient, config.RequestChannel)
	handler := bastion2.NewRequestHandler(publisher, config.RequestTtl)
	server := bastion2.NewServer(config.ListenAddress, handler)
	return server
}*/

func registerExtension(ctx context.Context, extensionClient *lambdaext.Client) {
	res, err := extensionClient.Register(ctx, extensionName)
	if err != nil {
		panic(err)
	}

	slog.Info("Register response:", "response", res)
}

// Method to process events
func processEvents(ctx context.Context, extensionClient *lambdaext.Client) error {
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
				continue
			}
		}
	}
}
