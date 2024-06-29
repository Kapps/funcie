package funcietunnel

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/clients/go/funcietunnel/internal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"net/url"
	"os"
)

// FuncieConfig is the basic configuration for both the local and Lambda versions of the Funcie tunnel.
type FuncieConfig struct {
	// ClientBastionEndpoint is the endpoint of the bastion running on the client machine.
	ClientBastionEndpoint url.URL `json:"clientBastionEndpoint"`
	// ServerBastionEndpoint is the endpoint of the bastion running in the cloud, often private.
	ServerBastionEndpoint url.URL `json:"serverBastionEndpoint"`
	// ListenAddress is the address that the local server will listen on.
	ListenAddress string `json:"listenAddress"`
	// ApplicationId is the ID of the application that the tunnel is for.
	ApplicationId string `json:"applicationId"`
}

// SsmParameterStoreClient is a minimal interface for the SSM client.
type SsmParameterStoreClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(options *ssm.Options)) (*ssm.GetParameterOutput, error)
}

// NewConfigFromEnvironment creates a new Config from environment variables.
// The following environment variables are used:
//
//	FUNCIE_APPLICATION_ID (required)
//	FUNCIE_CLIENT_BASTION_ENDPOINT (optional; for client, defaults to port 24193 on localhost)
//	FUNCIE_SERVER_BASTION_ENDPOINT (required for server)
//	FUNCIE_LISTEN_ADDRESS (optional; defaults to localhost on a random port)
func NewConfigFromEnvironment() *FuncieConfig {
	return &FuncieConfig{
		ClientBastionEndpoint: internal.OptionalUrlEnv("FUNCIE_CLIENT_BASTION_ENDPOINT", "http://127.0.0.1:24193"),
		ServerBastionEndpoint: internal.RequireUrlEnv("FUNCIE_SERVER_BASTION_ENDPOINT", internal.ConfigPurposeServer),
		ApplicationId:         internal.RequiredEnv("FUNCIE_APPLICATION_ID", internal.ConfigPurposeAny),
		ListenAddress:         internal.OptionalEnv("FUNCIE_LISTEN_ADDRESS", "0.0.0.0:0"),
	}
}

// NewConfig creates a new Config from a combination of environment variables and SSM parameters.
// The following variables are used:
//
//	FUNCIE_CLIENT_BASTION_ENDPOINT (optional; for client, defaults to port 24193 on localhost)
//	FUNCIE_SERVER_BASTION_ENDPOINT -> /funcie/<env>/bastion_host (required)
//	FUNCIE_LISTEN_ADDRESS (optional; defaults to localhost on a random port)
func NewConfig(ctx context.Context, applicationId string, ssmClient *ssm.Client) *FuncieConfig {
	serverEndpoint := os.Getenv("FUNCIE_SERVER_BASTION_ENDPOINT")
	if serverEndpoint == "" {
		// Do this after the env check to allow avoiding SSM calls entirely.
		serverEndpoint = loadSSMParameter(ctx, ssmClient, "default", "bastion_host")
	}

	return &FuncieConfig{
		ClientBastionEndpoint: internal.OptionalUrlEnv("FUNCIE_CLIENT_BASTION_ENDPOINT", "http://127.0.0.1:24193"),
		ServerBastionEndpoint: internal.OptionalUrlEnv("FUNCIE_SERVER_BASTION_ENDPOINT", serverEndpoint),
		ApplicationId:         applicationId,
		ListenAddress:         internal.OptionalEnv("FUNCIE_LISTEN_ADDRESS", "0.0.0.0:0"),
	}
}

func loadSSMParameter(ctx context.Context, ssmClient SsmParameterStoreClient, env string, name string) string {
	path := fmt.Sprintf("/funcie/%s/%s", env, name)
	req := ssm.GetParameterInput{
		Name: aws.String(path),
	}
	resp, err := ssmClient.GetParameter(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("failed to load SSM parameter %s: %s", path, err))
	}

	return *resp.Parameter.Value
}
