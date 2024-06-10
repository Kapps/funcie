package tools

import (
	"fmt"
)

// TerraformClient provides an interface for interacting with Terraform.
type TerraformClient interface {
	// Init initializes a Terraform project.
	Init(directory string) error

	// Apply applies a Terraform configuration with variables from the specified file.
	Apply(directory string, varFilePath string) error
}

type terraformCliClient struct {
	runner ProcessRunner
}

// NewTerraformCliClient creates a new TerraformClient that wraps the Terraform CLI tool.
func NewTerraformCliClient(runner ProcessRunner) TerraformClient {
	return &terraformCliClient{
		runner: runner,
	}
}

func (t *terraformCliClient) Init(directory string) error {
	_, err := t.runner.Run("terraform", directory, "init")
	if err != nil {
		return fmt.Errorf("failed to initialize Terraform project in directory %v: %w", directory, err)
	}

	return nil
}

func (t *terraformCliClient) Apply(directory string, varFilePath string) error {
	_, err := t.runner.Run("terraform", directory, "apply", "-auto-approve", "-var-file="+varFilePath, directory)
	if err != nil {
		return fmt.Errorf("failed to apply Terraform configuration in directory %v: %w", directory, err)
	}

	return nil
}
