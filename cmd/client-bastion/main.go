package main

import (
	"context"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	conf := NewConfigFromEnvironment()
	redisClient := redis.NewClient(&redis.Options{
		Addr:       conf.RedisAddress,
		ClientName: "Funcie Client Bastion",
	})

	handlerRouter := utils.NewClientHandlerRouter()

	consumer := r.NewConsumer(redisClient, conf.BaseChannelName, handlerRouter)

	err := consumer.Consume(ctx)
	if err != nil {
		panic("stopped consuming")
	}
}
