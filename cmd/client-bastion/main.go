package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/receiver"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

func newRedisClient(conf *bastion.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       conf.RedisAddress,
		ClientName: "funcie-client-bastion",
	})
}

func newApplicationRegistry(publisher funcie.Publisher) funcie.ApplicationRegistry {
	underlying := receiver.NewMemoryApplicationRegistry()
	return bastion.NewForwardingApplicationRegistry(underlying, publisher)
}

func newPublisher(redisClient *redis.Client, conf *bastion.Config) funcie.Publisher {
	return r.NewPublisher(redisClient, conf.BaseChannelName)
}

func newHost(conf *bastion.Config, messageProcessor bastion.MessageProcessor) bastion.Host {
	return bastion.NewHost(conf.ListenAddress, messageProcessor)
}

func main() {
	ctx := context.Background()
	fx.New(
		fx.Provide(
			func() context.Context { return ctx },
			func() *http.Client { return http.DefaultClient },
			bastion.NewConfigFromEnvironment,
			newRedisClient,
			utils.NewClientHandlerRouter,
			bastion.NewMessageProcessor,
			newApplicationRegistry,
			newPublisher,
			newHost,

			bastion.NewHTTPApplicationClient,
			bastion.NewHandler,
		),
		fx.Invoke(func(lc fx.Lifecycle, consumer funcie.Consumer, host bastion.Host) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return Start(ctx, consumer, host)
				},
				OnStop: func(ctx context.Context) error {
					return host.Close(ctx)
				},
			})
		}),
	).Run()
}

func Start(ctx context.Context, consumer funcie.Consumer, host bastion.Host) error {
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
	}

	return nil
}
