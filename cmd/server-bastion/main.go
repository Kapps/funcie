package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/cmd/server-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

func newRedisClient(config *bastion.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       config.RedisAddress,
		ClientName: "funcie-server-bastion",
	})
}

func newPublisher(redisClient *redis.Client, config *bastion.Config) funcie.Publisher {
	return r.NewPublisher(redisClient, config.RequestChannel)
}

func newConsumer(redisClient *redis.Client, config *bastion.Config, router utils.ClientHandlerRouter) funcie.Consumer {
	return r.NewConsumer(redisClient, config.RequestChannel, router)
}

func newHost(config *bastion.Config, processor transports.MessageProcessor) transports.Host {
	return transports.NewHost(config.ListenAddress, processor)
}

func main() {
	ctx := context.Background()

	funcie.ConfigureLogging()

	fx.New(
		fx.Provide(
			func() context.Context { return ctx },
			func() *http.Client { return http.DefaultClient },
			bastion.NewConfigFromEnvironment,
			newRedisClient,
			newPublisher,
			bastion.NewRequestHandler,
			newHost,
			transports.NewMessageProcessor,
			utils.NewClientHandlerRouter,
			newConsumer,
		),
		fx.Invoke(func(lc fx.Lifecycle, consumer funcie.Consumer, host transports.Host) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return Start(ctx, consumer, host)
				},
				OnStop: func(_ context.Context) error {
					return host.Close(ctx)
				},
			})
		}),
	).Run()
}

func Start(ctx context.Context, consumer funcie.Consumer, host transports.Host) error {
	err := consumer.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect to consumer: %w", err)
	}

	go func() {
		// Goroutine for host requests -- a socket for receiving messages from other clients.
		err := host.Listen(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorCtx(ctx, "host closed", err)
			os.Exit(1)
		}
		slog.WarnCtx(ctx, "host closed", "error", err.Error())
	}()

	go func() {
		// Goroutine for incoming messages -- registers on the consumer and starts listening.
		err := consumer.Consume(ctx)
		if err != nil {
			slog.ErrorCtx(ctx, "consume", err)
			os.Exit(1)
		}
		slog.WarnCtx(ctx, "consume", "error", err.Error())
	}()

	return nil
}
