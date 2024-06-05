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

	privSubnets, pubSubnets, err := c.promptSubnet(ctx, vpc.Id)
	if err != nil {
		return fmt.Errorf("failed to prompt for subnet: %w", err)
	}

	elastiCache, err := c.promptElasticache(ctx)
	if err != nil {
		return fmt.Errorf("failed to prompt for ElastiCache cluster: %w", err)
	}

	fmt.Println("Selected VPC:", vpc.Name)
	fmt.Printf("Selected private subnets: %v\n", privSubnets)
	fmt.Printf("Selected public subnets: %v\n", pubSubnets)
	fmt.Println("Selected ElastiCache cluster:", elastiCache)

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

func (c *InitCommand) promptSubnet(ctx context.Context, vpcId string) ([]aws.Subnet, []aws.Subnet, error) {
	subnets, err := c.resourceList.ListSubnets(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list subnets: %w", err)
	}

	filtered := make([]aws.Subnet, 0, 4)
	for _, subnet := range subnets {
		if subnet.VpcId == vpcId {
			filtered = append(filtered, subnet)
		}
	}

	pubSubnets := make([]aws.Subnet, 0, 4)
	privSubnets := make([]aws.Subnet, 0, 4)
	for _, subnet := range filtered {
		if subnet.Public {
			pubSubnets = append(pubSubnets, subnet)
		} else {
			privSubnets = append(privSubnets, subnet)
		}
	}

	//var selected []aws.Subnet
	err = huh.NewMultiSelect[aws.Subnet]().
		Title("Which public subnets would you like to use?").
		Options(huh.NewOptions[aws.Subnet](filtered...)...).
		Description("This will be used for resources such as the bastion host.").
		Value(&pubSubnets).
		Run()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to select public subnets: %w", err)
	}

	err = huh.NewMultiSelect[aws.Subnet]().
		Title("Which private subnets would you like to use?").
		Options(huh.NewOptions[aws.Subnet](filtered...)...).
		Description("This will be used for resources such as the Elasticache instance.").
		Value(&privSubnets).
		Run()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to select private subnets: %w", err)
	}

	return privSubnets, pubSubnets, nil
}

func (c *InitCommand) promptElasticache(ctx context.Context) (*aws.ElastiCacheCluster, error) {
	clusters, err := c.resourceList.ListElastiCacheClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list ElastiCache clusters: %w", err)
	}

	var selected aws.ElastiCacheCluster
	err = huh.NewSelect[aws.ElastiCacheCluster]().
		Title("Which ElastiCache cluster would you like to use?").
		Description("Leave blank to have funcie provision a new single-node cluster.").
		Options(append(
			[]huh.Option[aws.ElastiCacheCluster]{huh.NewOption[aws.ElastiCacheCluster]("<create new cluster>", aws.ElastiCacheCluster{})},
			huh.NewOptions[aws.ElastiCacheCluster](clusters...)...,
		)...).
		Value(&selected).
		Run()
	if err != nil {
		return nil, fmt.Errorf("failed to select ElastiCache cluster: %w", err)
	}

	if selected.Name == "" {
		return nil, nil
	}

	return &selected, nil
}
