package funcie_tunnel

import (
	"fmt"
	"net/url"
	"os"
)

type FuncieConfig struct {
	ClientBastionEndpoint url.URL `json:"bastionEndpoint"`
	ServerBastionEndpoint url.URL `json:"bastionEndpoint"`
	ListenAddress         string  `json:"listenAddress"`
	ApplicationId         string  `json:"applicationId"`
}

// NewConfigFromEnvironment creates a new Config from environment variables.
// The following environment variables are used:
//
//	FUNCIE_APPLICATION_ID (required)
//	FUNCIE_CLIENT_BASTION_ENDPOINT (required for client)
//	FUNCIE_SERVER_BASTION_ENDPOINT (required for server)
//	FUNCIE_LISTEN_ADDRESS (required for client)
func NewConfigFromEnvironment() *FuncieConfig {
	return &FuncieConfig{
		ClientBastionEndpoint: requireUrlEnv("FUNCIE_CLIENT_BASTION_ENDPOINT"),
		ServerBastionEndpoint: requireUrlEnv("FUNCIE_SERVER_BASTION_ENDPOINT"),
		ListenAddress:         requiredEnv("FUNCIE_LISTEN_ADDRESS"),
		ApplicationId:         requiredEnv("FUNCIE_APPLICATION_ID"),
	}
}

func requireUrlEnv(name string) url.URL {
	value := requiredEnv(name)
	url, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %s %s: %s", name, value, err))
	}
	return *url
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
