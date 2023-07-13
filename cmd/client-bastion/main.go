package main

import (
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/Kapps/funcie/pkg/receiver"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"net/http"
)

func newRedisClient(conf *bastion.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       conf.RedisAddress,
		ClientName: "Funcie Client Bastion",
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
		fx.Invoke(func(lc fx.Lifecycle, host bastion.Host) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return host.Listen(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return host.Close(ctx)
				},
			})
		}),
	).Run()
}
