package bastion

import (
	"fmt"
	"os"
	"time"
)

// Config allows the configuration of the Bastion.
type Config struct {
	// RedisAddress is the address of the Redis server.
	RedisAddress string `json:"redisAddress"`
	// ListenAddress is the address to listen on.
	ListenAddress string `json:"listenAddress"`
	// RequestTtl indicates the time to live for a request.
	RequestTtl time.Duration `json:"requestTtl"`
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
//	FUNCIE_REQUEST_TTL (optional; defaults to 15 minutes; values are parsed using time.ParseDuration)
func NewConfigFromEnvironment() *Config {
	return &Config{
		RedisAddress:  requiredEnv("FUNCIE_REDIS_ADDRESS"),
		ListenAddress: requiredEnv("FUNCIE_LISTEN_ADDRESS"),
		RequestTtl:    parseTimeDuration(optionalEnv("FUNCIE_REQUEST_TTL", "15m")),
	}
}

func parseTimeDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("failed to parse duration %s: %v", value, err))
	}
	return duration
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
