package funcli

import (
	"fmt"
)

// Cli is the top-level command-line arguments.
type Cli struct {
	ConnectArgs *ConnectCommand `arg:"subcommand:connect" help:"Connect to a funcie deployment to allow local development."`
	Environment string          `arg:"--env" help:"Funcie environment used if multiple deployments are present." default:"default"`

	VersionString string `arg:"-"`
}

func (c *Cli) Version() string {
	return fmt.Sprintf("funcie v%v", c.VersionString)
}

func NewCli(connectArgs *ConnectCommand) *Cli {
	return &Cli{
		ConnectArgs: connectArgs,
	}
}
