package provider

import (
	"fmt"
	"net/url"
	"os"
)

type FuncieConfig struct {
	BastionEndpoint url.URL `json:"bastionEndpoint"`
	ApplicationId   string  `json:"applicationId"`
}

// NewConfigFromEnvironment creates a new Config from environment variables.
// The following environment variables are used:
//
//	FUNCIE_BASTION_ENDPOINT (required)
//	FUNCIE_LISTEN_ADDRESS (required)
//	FUNCIE_BASE_CHANNEL_NAME (optional)
func NewConfigFromEnvironment() *FuncieConfig {
	bastionEndpoint := requiredEnv("FUNCIE_BASTION_ENDPOINT")
	bastionEndpointUrl, err := url.Parse(bastionEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed to parse bastion endpoint %s: %s", bastionEndpoint, err))
	}
	return &FuncieConfig{
		BastionEndpoint: *bastionEndpointUrl,
		ApplicationId:   requiredEnv("FUNCIE_APPLICATION_ID"),
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
