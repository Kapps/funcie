package funcie_tunnel

import (
	"fmt"
	"github.com/Kapps/funcie/clients/go/funcie-tunnel/internal"
	"net/url"
	"os"
)

type FuncieConfig struct {
	ClientBastionEndpoint url.URL `json:"clientBastionEndpoint"`
	ServerBastionEndpoint url.URL `json:"serverBastionEndpoint"`
	ListenAddress         string  `json:"listenAddress"`
	ApplicationId         string  `json:"applicationId"`
}

// NewConfigFromEnvironment creates a new Config from environment variables.
// The following environment variables are used:
//
//	FUNCIE_APPLICATION_ID (required)
//	FUNCIE_CLIENT_BASTION_ENDPOINT (required for client)
//	FUNCIE_SERVER_BASTION_ENDPOINT (required for server)
//	FUNCIE_LISTEN_ADDRESS (optional; defaults to localhost on a random port)
func NewConfigFromEnvironment() *FuncieConfig {
	return &FuncieConfig{
		ClientBastionEndpoint: requireUrlEnv("FUNCIE_CLIENT_BASTION_ENDPOINT", internal.ConfigPurposeClient),
		ServerBastionEndpoint: requireUrlEnv("FUNCIE_SERVER_BASTION_ENDPOINT", internal.ConfigPurposeServer),
		ApplicationId:         requiredEnv("FUNCIE_APPLICATION_ID", internal.ConfigPurposeAny),
		ListenAddress:         optionalEnv("FUNCIE_LISTEN_ADDRESS", "0.0.0.0:0"),
	}
}

func requireUrlEnv(name string, purpose internal.ConfigPurpose) url.URL {
	value := requiredEnv(name, purpose)
	parsedUrl, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %s %s: %s", name, value, err))
	}
	return *parsedUrl
}

func requiredEnv(name string, purpose internal.ConfigPurpose) string {
	value := os.Getenv(name)
	if value == "" {
		currPurpose := internal.GetConfigPurpose()
		purposeMatches := purpose == internal.ConfigPurposeAny || currPurpose == purpose
		if purposeMatches {
			panic(fmt.Sprintf("required environment variable %s not set", name))
		}
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
