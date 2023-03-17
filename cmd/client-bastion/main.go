package main

import "github.com/redis/go-redis/v9"

func main() {
	conf := NewConfigFromEnvironment()
	redisClient := redis.NewClient(&redis.Options{
		Addr: conf.RedisAddress,
	})

}
