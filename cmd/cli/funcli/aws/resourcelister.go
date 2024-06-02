package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Vpc provides some basic information about a VPC.
type Vpc struct {
	// ID is the AWS representation of the VPC ID (not the ARN).
	Id string
	// Name is a human-readable name for the VPC.
	Name string
}

type Subnet struct {
	// ID is the AWS representation of the Subnet ID (not the ARN).
	Id string
	// Name is a human-readable name for the Subnet.
	Name string
	// Public indicates whether the subnet is public or private.
	Public bool
	// VpcId is the ID of the VPC that the subnet is associated with.
	VpcId string
}

// ResourceLister allows retrieving a list of AWS resources.
type ResourceLister interface {
	ListVpcs(ctx context.Context) ([]Vpc, error)
	ListSubnets(ctx context.Context) ([]Subnet, error)
}

// Ec2VpcClient is a minimal interface for the EC2 client.
type Ec2VpcClient interface {
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeRouteTables(ctx context.Context, params *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error)
}

type awsResourceLister struct {
	ec2Client Ec2VpcClient
}

// NewAwsResourceLister creates a new ResourceLister that uses the provided EC2 client.
func NewAwsResourceLister(client Ec2VpcClient) ResourceLister {
	return &awsResourceLister{
		ec2Client: client,
	}
}

func (r *awsResourceLister) ListVpcs(ctx context.Context) ([]Vpc, error) {
	input := &ec2.DescribeVpcsInput{}

	result, err := r.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list VPCs: %w", err)
	}

	var vpcs []Vpc
	for _, vpc := range result.Vpcs {
		vpcs = append(vpcs, Vpc{
			Id:   *vpc.VpcId,
			Name: *vpc.Tags[0].Value,
		})
	}

	return vpcs, nil
}

func (r *awsResourceLister) ListSubnets(ctx context.Context) ([]Subnet, error) {
	input := &ec2.DescribeSubnetsInput{}

	result, err := r.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %w", err)
	}

	subnetIds := make([]*string, 0, len(result.Subnets))
	var subnets []Subnet
	for _, subnet := range result.Subnets {
		subnets = append(subnets, Subnet{
			Id:    *subnet.SubnetId,
			Name:  *subnet.Tags[0].Value,
			VpcId: *subnet.VpcId,
		})
		subnetIds = append(subnetIds, subnet.SubnetId)
	}

	routeTables, err := r.ec2Client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: subnetIds,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list route tables: %w", err)
	}

	for i, subnet := range subnets {
		for _, routeTable := range routeTables.RouteTables {
			for _, association := range routeTable.Associations {
				if *association.SubnetId == subnet.Id {
					subnets[i].Public = true
					break
				}
			}
		}
	}

	return subnets, nil
}
