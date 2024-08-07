package main

import (
	"context"
	"errors"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/receiver"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func newRedisClient(conf *bastion.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       conf.RedisAddress,
		ClientName: "funcie-client-bastion",
	})
}

func newApplicationRegistry(redis *redis.Client) funcie.ApplicationRegistry {
	return receiver.NewRedisApplicationRegistry(redis)
}

func newPublisher(redisClient *redis.Client, conf *bastion.Config) funcie.Publisher {
	return r.NewPublisher(redisClient, conf.BaseChannelName)
}

func newHost(conf *bastion.Config, messageProcessor transports.MessageProcessor) transports.Host {
	return transports.NewHost(conf.ListenAddress, messageProcessor)
}

func newConsumer(redisClient *redis.Client, conf *bastion.Config, router utils.ClientHandlerRouter) funcie.Consumer {
	return r.NewConsumer(redisClient, conf.BaseChannelName, router)
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
			utils.NewClientHandlerRouter,
			transports.NewMessageProcessor,
			newApplicationRegistry,
			newPublisher,
			newHost,
			newConsumer,
			bastion.NewHTTPApplicationClient,
			bastion.NewHandler,
			bastion.NewDockerHostTranslator,
		),
		fx.StartTimeout(time.Hour*24*365*100), // Effectively infinite timeout to allow launching without starting Redis tunnel
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
	for {
		err := consumer.Connect(ctx)
		if err == nil {
			break
		}

		slog.WarnContext(ctx, "failed to connect to Redis; trying again in 10 seconds", "error", err.Error())
		time.Sleep(10 * time.Second)
	}

	go func() {
		// Goroutine for host requests -- a socket for receiving messages from other clients.
		err := host.Listen(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "host closed", err)
			os.Exit(1)
		}
		slog.WarnContext(ctx, "host closed", "error", err.Error())
	}()

	go func() {
		// Goroutine for incoming messages -- registers on the consumer and starts listening.
		err := consumer.Consume(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "consume", err)
			os.Exit(1)
		}
		slog.WarnContext(ctx, "consume", "error", err.Error())
	}()

	return nil
}
