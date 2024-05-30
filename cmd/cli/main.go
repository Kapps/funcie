package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli"
	"github.com/alexflint/go-arg"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"go.uber.org/fx"
	"os"
	"strings"
)

//go:embed version.txt
var version string

// Runnable is an interface for commands that can be run.
type Runnable interface {
	Run(ctx context.Context) error
}

func main() {
	cli, err := makeCli()
	if err != nil {
		fmt.Println("Failed to initialize CLI:", err)
		return
	}

	argConfig := arg.Config{
		Program: "funcie",
	}

	parser, err := arg.NewParser(argConfig, cli)
	if err != nil {
		fmt.Println("Failed to initialize parser:", err)
		return
	}

	parser.MustParse(os.Args[1:])

	cmd := parser.Subcommand()
	if cmd == nil {
		parser.WriteHelp(os.Stdout)
		return
	}

	ctx := context.Background()
	if err := cmd.(Runnable).Run(ctx); err != nil {
		fmt.Println(err)
		return
	}
}

func makeCli() (*funcli.Cli, error) {
	var res *funcli.Cli

	app := fx.New(
		fx.Provide(
			func() context.Context {
				return context.Background()
			},
			config.LoadDefaultConfig,
			ssm.NewFromConfig,
			func(ssmClient *ssm.Client) funcli.ConfigStore {
				return funcli.NewConfigStore("default", ssmClient)
			},
			func(configStore funcli.ConfigStore, connectClient *ssm.Client, tunneller funcli.Tunneller) *funcli.ConnectCommand {
				return funcli.NewConnectCommand(configStore, connectClient, tunneller)
			},
			funcli.NewCli,
			funcli.NewWebhookTunneller,
		),
		fx.NopLogger,
		fx.Populate(&res),
	)

	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("failed to initialize CLI: %w", err)
	}

	res.VersionString = strings.TrimSpace(version) // Feels icky

	return res, nil
}
