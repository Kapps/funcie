package main

import (
	"fmt"
	"os"
)

type Config struct {
	// RedisAddress is the address of the Redis server.
	RedisAddress string `json:"redisAddress"`
	// ListenAddress is the address to listen on for client requests.
	ListenAddress string `json:"listenAddress"`
}

// NewConfig creates a new Config with no values set.
func NewConfig() *Config {
	return &Config{}
}

// NewConfigFromEnvironment creates a new Config from environment variables.
// The following environment variables are used:
//
//	FUNCIE_REDIS_ADDRESS (required)
//	FUNCIE_LISTEN_ADDRESS (required)
func NewConfigFromEnvironment() *Config {
	return &Config{
		RedisAddress:  requiredEnv("FUNCIE_REDIS_ADDRESS"),
		ListenAddress: requiredEnv("FUNCIE_LISTEN_ADDRESS"),
	}
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s not set", name))
	}
	return value
}
