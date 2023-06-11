package main

import (
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	conf := bastion.NewConfigFromEnvironment()
	redisClient := redis.NewClient(&redis.Options{
		Addr:       conf.RedisAddress,
		ClientName: "Funcie Client Bastion",
	})

	handlerRouter := utils.NewClientHandlerRouter()

	consumer := r.NewConsumer(redisClient, conf.BaseChannelName, handlerRouter)
	registry := r.
	err := consumer.Consume(ctx)
	if err != nil {
		panic("stopped consuming")
	}
}
