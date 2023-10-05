package bastion

import (
	"fmt"
	"os"
)

type Config struct {
	// RedisAddress is the address of the Redis server.
	RedisAddress string `json:"redisAddress"`
	// ListenAddress is the address to listen on for client requests.
	ListenAddress string `json:"listenAddress"`
	// BaseChannelName is the base name of the Redis channel keys to use.
	BaseChannelName string `json:"baseChannelName"`
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
//	FUNCIE_BASE_CHANNEL_NAME (optional)
func NewConfigFromEnvironment() *Config {
	return &Config{
		RedisAddress:    requiredEnv("FUNCIE_REDIS_ADDRESS"),
		ListenAddress:   optionalEnv("FUNCIE_LISTEN_ADDRESS", "127.0.0.1:24193"),
		BaseChannelName: optionalEnv("FUNCIE_BASE_CHANNEL_NAME", "funcie:requests"),
	}
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s not set", name))
	}
	return value
}

func optionalEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
