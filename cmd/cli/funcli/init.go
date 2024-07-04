package funcli

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli/aws"
	"github.com/Kapps/funcie/cmd/cli/funcli/internal"
	"github.com/Kapps/funcie/cmd/cli/funcli/tools"
	"github.com/charmbracelet/huh"
	"os"
	"strings"
)

const tfModuleRepo = "git@github.com:Kapps/terraform-aws-funcie.git"

type InitConfig struct {
}

type InitCommand struct {
	cliConfig       *CliConfig
	resourceList    aws.ResourceLister
	gitClient       tools.GitClient
	terraformClient tools.TerraformClient
}

type TerraformVars struct {
	VpcId          string   `yaml:"vpc_id"`
	PrivateSubnets []string `yaml:"private_subnet_ids"`
	PublicSubnets  []string `yaml:"public_subnet_ids"`
	RedisHost      string   `yaml:"redis_host"`
	Region         string   `yaml:"region"`
}

func NewInitCommand(
	cliConfig *CliConfig,
	resourceList aws.ResourceLister,
	gitClient tools.GitClient,
	terraformClient tools.TerraformClient,
) *InitCommand {
	return &InitCommand{
		cliConfig:       cliConfig,
		resourceList:    resourceList,
		gitClient:       gitClient,
		terraformClient: terraformClient,
	}
}

func (c *InitCommand) Run(ctx context.Context) error {
	vpc, err := c.promptVpc(ctx)
	if err != nil {
		return fmt.Errorf("failed to prompt for VPC: %w", err)
	}

	var privSubnets, pubSubnets []aws.Subnet
	if vpc.Id != "" {
		privSubnets, pubSubnets, err = c.promptSubnet(ctx, vpc.Id)
		if err != nil {
			return fmt.Errorf("failed to prompt for subnet: %w", err)
		}
	}

	elastiCache, err := c.promptElasticache(ctx)
	if err != nil {
		return fmt.Errorf("failed to prompt for ElastiCache cluster: %w", err)
	}

	vars := TerraformVars{
		VpcId:          vpc.Id,
		PrivateSubnets: make([]string, 0, len(privSubnets)),
		PublicSubnets:  make([]string, 0, len(pubSubnets)),
		Region:         c.cliConfig.Region,
	}

	if elastiCache != nil {
		vars.RedisHost = elastiCache.PrimaryEndpoint
	} else {
		vars.RedisHost = ""
	}

	for _, subnet := range privSubnets {
		vars.PrivateSubnets = append(vars.PrivateSubnets, subnet.Id)
	}
	for _, subnet := range pubSubnets {
		vars.PublicSubnets = append(vars.PublicSubnets, subnet.Id)
	}

	if err := c.writeTerraformVars(vars); err != nil {
		return fmt.Errorf("failed to write terraform vars: %w", err)
	}

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
		Description("Or default to have funcie provision / reuse a VPC.").
		Options(append(
			[]huh.Option[aws.Vpc]{huh.NewOption[aws.Vpc]("<create funcie-managed VPC with NAT instance config>", aws.Vpc{})},
			huh.NewOptions[aws.Vpc](vpcs...)...,
		)...).
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

	err = huh.NewMultiSelect[aws.Subnet]().
		Title("Which private subnets would you like to use?").
		Options(huh.NewOptions[aws.Subnet](filtered...)...).
		Description("This will be used for resources such as the Elasticache instance.").
		Value(&privSubnets).
		Run()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to select private subnets: %w", err)
	}

	err = huh.NewMultiSelect[aws.Subnet]().
		Title("Which public subnets would you like to use?").
		Options(huh.NewOptions[aws.Subnet](filtered...)...).
		Description("This will be used for resources such as the bastion host.").
		Value(&pubSubnets).
		Run()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to select public subnets: %w", err)
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
		Description("Or default to have funcie provision / reuse a cluster.").
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

func (c *InitCommand) writeTerraformVars(vars TerraformVars) error {
	varFileContents := marshalVariables(vars)

	baseDir := internal.GetFuncieBaseDir()
	varFile := internal.GetTerraformVarsPath()

	if err := os.Mkdir(baseDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory for terraform vars: %w", err)
	}

	if err := os.WriteFile(varFile, []byte(varFileContents), 0644); err != nil {
		return fmt.Errorf("failed to write terraform vars to %s: %w", varFile, err)
	}

	branch := "v" + strings.TrimSpace(c.cliConfig.versionString)
	fmt.Printf("Cloning funcie terraform module with tag %v...\n", branch)

	tfModuleDir := internal.GetTerraformDir()
	if err := c.gitClient.Checkout(tfModuleRepo, tfModuleDir, branch); err != nil {
		return fmt.Errorf("failed to clone terraform module: %w", err)
	}

	fmt.Printf("Cloned terraform module to %s; running init.\n", tfModuleDir)

	if err := c.terraformClient.Init(tfModuleDir); err != nil {
		return fmt.Errorf("failed to initialize terraform module: %w", err)
	}

	fmt.Println("Terraform module initialized; will deploy module with the following parameters:")
	fmt.Println(varFileContents)

	confirmed := false
	err := huh.NewConfirm().
		Title("Would you like to proceed with the deployment? Remember, some resources are not part of the AWS free tier and therefore charges will apply.").
		Affirmative("I understand, deploy module").
		Value(&confirmed).
		Run()
	if err != nil {
		return fmt.Errorf("failed to confirm deployment: %w", err)
	}

	if !confirmed {
		return fmt.Errorf("deployment cancelled")
	}

	fmt.Println("Applying terraform configuration...")

	if err := c.terraformClient.Apply(tfModuleDir, varFile); err != nil {
		return fmt.Errorf("failed to apply terraform configuration: %w", err)
	}

	fmt.Println("Terraform configuration applied successfully.")

	return nil
}

func marshalVariables(vars TerraformVars) string {
	return fmt.Sprintf(`
vpc_id             = "%s"
private_subnet_ids = %s
public_subnet_ids  = %s
redis_host         = "%s"
region             = "%s"
`, vars.VpcId, internal.MarshalArray(vars.PrivateSubnets), internal.MarshalArray(vars.PublicSubnets), vars.RedisHost, vars.Region)
}
