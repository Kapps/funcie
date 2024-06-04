package funcli

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli/aws"
	"github.com/c-bata/go-prompt"
)

type InitConfig struct {
	OutputFile string `arg:"--output-file,-o" help:"Output file to write the generated terraform configuration to."`
}

type InitCommand struct {
	cliConfig    *CliConfig
	resourceList aws.ResourceLister
}

func NewInitCommand(cliConfig *CliConfig, resourceList aws.ResourceLister) *InitCommand {
	return &InitCommand{
		cliConfig:    cliConfig,
		resourceList: resourceList,
	}
}

func (c *InitCommand) Run(ctx context.Context) error {
	vpc, err := c.promptVpc(ctx)
	if err != nil {
		return fmt.Errorf("failed to prompt for VPC: %w", err)
	}

	fmt.Printf("Selected VPC: %s\n", vpc)

	return nil
}

func (c *InitCommand) promptVpc(ctx context.Context) (string, error) {
	vpcs, err := c.resourceList.ListVpcs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list vpcs: %w", err)
	}

	vpcNames := make([]prompt.Suggest, 0, len(vpcs))
	for _, vpc := range vpcs {
		vpcNames = append(vpcNames, prompt.Suggest{Text: vpc.Name, Description: vpc.Id})
	}

	completer := func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterFuzzy(vpcNames, d.GetWordBeforeCursor(), true)
	}

	promp
	p := prompt.New(
		func(in string) {},
		completer,
		prompt.OptionPrefix("Select a VPC: "),
		prompt.OptionShowCompletionAtStart(),
	)

	return p.Input(), nil
}
