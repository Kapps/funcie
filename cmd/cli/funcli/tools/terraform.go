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
	// Destroy destroys a Terraform project.
	Destroy(directory string, varFilePath string) error
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

func (t *terraformCliClient) run(cmd string, args ...string) (string, error) {
	return t.runner.RunWithOptions(cmd, RunnerOpts{
		Args: args,
		Env:  map[string]string{"TF_IN_AUTOMATION": "true"},
	})
}

func (t *terraformCliClient) Init(directory string) error {
	_, err := t.run("terraform", "-chdir="+directory, "init", "-input=false")
	if err != nil {
		return fmt.Errorf("failed to initialize Terraform project in directory %v: %w", directory, err)
	}

	return nil
}

func (t *terraformCliClient) Apply(directory string, varFilePath string) error {
	_, err := t.run("terraform", "-chdir="+directory, "apply", "-input=false", "-auto-approve", "-var-file="+varFilePath)
	if err != nil {
		return fmt.Errorf("failed to apply Terraform configuration in directory %v: %w", directory, err)
	}

	return nil
}

func (t *terraformCliClient) Destroy(directory string, varFilePath string) error {
	_, err := t.run("terraform", "-chdir="+directory, "destroy", "-input=false", "-auto-approve", "-var-file="+varFilePath)
	if err != nil {
		return fmt.Errorf("failed to destroy Terraform project in directory %v: %w", directory, err)
	}

	return nil
}
