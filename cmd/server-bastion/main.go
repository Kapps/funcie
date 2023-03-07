package main

import (
	"github.com/Kapps/funcie/pkg/bastion"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	config := NewConfigFromEnvironment()
	redisClient := &redis.Client{}
	publisher := r.NewPublisher(redisClient, config.RequestChannel)
	handler := bastion.NewRequestHandler(publisher, config.RequestTtl)
	server := bastion.NewServer(config.ListenAddress, handler)

	err := server.Listen()
	if err != nil {
		panic("stopped listening on server")
	}
}
