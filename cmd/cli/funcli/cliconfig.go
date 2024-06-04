package funcli

import "fmt"

// CliConfig provides user input parameters specific to the CLI tool.
type CliConfig struct {
	ConnectConfig *ConnectConfig `arg:"subcommand:connect" help:"Connect to a funcie deployment to allow local development."`

	Environment string `arg:"--env" help:"Funcie environment used if multiple deployments are present." default:"default"`
	Region      string `arg:"env:AWS_REGION" help:"AWS region to use for deployments; otherwise uses the default AWS CLI region."`

	versionString string `arg:"-"`
}

func (c *CliConfig) Version() string {
	return fmt.Sprintf("funcie v%v", c.versionString)
}

func NewCliConfig(version string) *CliConfig {
	return &CliConfig{
		versionString: version,
	}
}
