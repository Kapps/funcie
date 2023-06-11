package main

import (
	"github.com/Kapps/funcie/cmd/server-bastion/bastion"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	//panic("not implemented")
	config := bastion.NewConfigFromEnvironment()
	redisClient := redis.NewClient(&redis.Options{
		Addr:       config.RedisAddress,
		ClientName: "Funcie Server Bastion",
	})

	publisher := r.NewPublisher(redisClient, config.RequestChannel)
	handler := bastion.NewRequestHandler(publisher, config.RequestTtl)
	server := bastion.NewServer(config.ListenAddress, handler)

	err := server.Listen()
	if err != nil {
		panic("stopped listening on server")
	}
}
