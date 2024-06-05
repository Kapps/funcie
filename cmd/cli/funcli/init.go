package funcli

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli/aws"
	"github.com/charmbracelet/huh"
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

	subnet, err := c.promptSubnet(ctx, vpc.Id)
	if err != nil {
		return fmt.Errorf("failed to prompt for subnet: %w", err)
	}

	fmt.Printf("Selected Subnet: %v\n", subnet)

	return nil
}

func (c *InitCommand) promptVpc(ctx context.Context) (aws.Vpc, error) {
	vpcs, err := c.resourceList.ListVpcs(ctx)
	if err != nil {
		return aws.Vpc{}, fmt.Errorf("failed to list vpcs: %w", err)
	}

	var selected aws.Vpc
	err = huh.NewSelect[aws.Vpc]().
		Title("Which VPC would you like to use?").
		Options(huh.NewOptions[aws.Vpc](vpcs...)...).
		Value(&selected).
		Run()
	if err != nil {
		return aws.Vpc{}, fmt.Errorf("failed to select VPC: %w", err)
	}

	return selected, nil
}

func (c *InitCommand) promptSubnet(ctx context.Context, vpcId string) (aws.Subnet, error) {
	subnets, err := c.resourceList.ListSubnets(ctx)
	if err != nil {
		return aws.Subnet{}, fmt.Errorf("failed to list subnets: %w", err)
	}

	filtered := make([]aws.Subnet, 0, 4)
	for _, subnet := range subnets {
		if subnet.VpcId == vpcId {
			filtered = append(filtered, subnet)
		}
	}

	var selected aws.Subnet
	err = huh.NewSelect[aws.Subnet]().
		Title("Which subnet would you like to use?").
		Options(huh.NewOptions[aws.Subnet](filtered...)...).
		Value(&selected).
		Run()
	if err != nil {
		return aws.Subnet{}, fmt.Errorf("failed to select subnet: %w", err)
	}

	return selected, nil
}
