package funcli

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// ConfigStore retrieves dynamic configuration values on-demand.
type ConfigStore interface {
	// GetConfigValue retrieves the value of a configuration key.
	GetConfigValue(ctx context.Context, key string) (string, error)
}

type SsmParameterStore interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type configStore struct {
	environment string
	ssmClient   SsmParameterStore
}

// NewConfigStore creates a new ConfigStore.
func NewConfigStore(environment string, ssmClient SsmParameterStore) ConfigStore {
	return &configStore{
		environment: environment,
		ssmClient:   ssmClient,
	}
}

func (c *configStore) GetConfigValue(ctx context.Context, key string) (string, error) {
	param, err := c.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/funcie/%v/%v", c.environment, key)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get config key %v: %w", key, err)
	}

	return *param.Parameter.Value, nil
}
