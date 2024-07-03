package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli"
	funcliAws "github.com/Kapps/funcie/cmd/cli/funcli/aws"
	"github.com/Kapps/funcie/cmd/cli/funcli/tools"
	"github.com/alexflint/go-arg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"go.uber.org/fx"
	"os"
)

//go:embed version.txt
var version string

// Runnable is an interface for commands that can be run.
type Runnable interface {
	Run(ctx context.Context) error
}

func main() {
	cliConfig := funcli.NewCliConfig(version)

	argConfig := arg.Config{
		Program: "funcie",
	}

	parser, err := arg.NewParser(argConfig, cliConfig)
	if err != nil {
		fmt.Println("Failed to initialize parser:", err)
		return
	}

	parser.MustParse(os.Args[1:])

	cli, err := makeCli(cliConfig)
	if err != nil {
		fmt.Println("Failed to initialize CLI:", err)
		return
	}

	cmd := parser.Subcommand()
	if cmd == nil {
		parser.WriteHelp(os.Stdout)
		return
	}

	cmdInstance, ok := cli.commands[cmd]
	if !ok {
		panic("command not found")
	}

	ctx := context.Background()
	if err := cmdInstance.Run(ctx); err != nil {
		fmt.Println(err)
		return
	}
}

func makeCli(cliConfig *funcli.CliConfig) (*cli, error) {
	var res *cli

	app := fx.New(
		fx.Supply(cliConfig),
		fx.Supply(fx.Annotate(context.Background(), fx.As(new(context.Context)))),
		fx.Provide(
			fx.Annotate(ssm.NewFromConfig, fx.As(new(funcliAws.SsmClient))),
			fx.Annotate(ec2.NewFromConfig, fx.As(new(funcliAws.EC2Client))),
			fx.Annotate(elasticache.NewFromConfig, fx.As(new(funcliAws.ElastiCacheClient))),
			loadAwsConfig,
			funcli.NewConfigStore,
			funcli.NewConnectCommand,
			funcli.NewInitCommand,
			newCli,
			funcli.NewSsmTunneller,
			funcli.NewHttpConnectivityService,
			funcliAws.NewAwsResourceLister,
			tools.NewProcessRunner,
			tools.NewGitCliClient,
			tools.NewTerraformCliClient,
			funcli.NewDestroyCommand,
		),
		fx.NopLogger,
		fx.Populate(&res),
	)

	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("failed to initialize CLI: %w", err)
	}

	return res, nil
}

func loadAwsConfig(cliConfig *funcli.CliConfig) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error
	if cliConfig.Region != "" {
		opts = append(opts, config.WithRegion(cliConfig.Region))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if cliConfig.Region == "" {
		cliConfig.Region = cfg.Region
	}

	return cfg, nil
}

func newCli(
	conf *funcli.CliConfig,
	connectCmd *funcli.ConnectCommand,
	initCmd *funcli.InitCommand,
	destroyCmd *funcli.DestroyCommand,
) *cli {
	inst := &cli{
		commands: make(map[interface{}]Runnable),
	}
	inst.RegisterCommand(conf.ConnectConfig, connectCmd)
	inst.RegisterCommand(conf.InitConfig, initCmd)
	inst.RegisterCommand(conf.DestroyConfig, destroyCmd)

	return inst
}

type cli struct {
	commands map[interface{}]Runnable
}

func (c *cli) RegisterCommand(instance interface{}, cmd Runnable) {
	c.commands[instance] = cmd
}
