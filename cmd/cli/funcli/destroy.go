package funcli

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli/internal"
	"github.com/Kapps/funcie/cmd/cli/funcli/tools"
	"os"
)

type DestroyConfig struct {
}

type DestroyCommand struct {
	cliConfig       *CliConfig
	terraformClient tools.TerraformClient
}

func NewDestroyCommand(
	cliConfig *CliConfig,
	terraformClient tools.TerraformClient,
) *DestroyCommand {
	return &DestroyCommand{
		cliConfig:       cliConfig,
		terraformClient: terraformClient,
	}
}

func (c *DestroyCommand) Run(ctx context.Context) error {
	tfDir := internal.GetTerraformDir()
	varsFile := internal.GetTerraformVarsPath()

	_, err := os.Stat(varsFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("could not find an instance of funcie to destroy")
	}

	if err := c.terraformClient.Destroy(tfDir, varsFile); err != nil {
		return fmt.Errorf("failed to destroy funcie instance: %w", err)
	}

	return nil
}
