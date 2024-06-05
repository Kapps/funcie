package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go/aws"
	"sort"
	"strings"
)

// ResourceLister allows retrieving a list of AWS resources.
type ResourceLister interface {
	// ListVpcs returns a list of VPCs within the region configured in the client.
	ListVpcs(ctx context.Context) ([]Vpc, error)
	// ListSubnets returns a list of all Subnets within the region configured in the client, regardless of VPC.
	ListSubnets(ctx context.Context) ([]Subnet, error)
	// ListElastiCacheClusters returns a list of all ElastiCache clusters within the region configured in the client.
	ListElastiCacheClusters(ctx context.Context) ([]ElastiCacheCluster, error)
}

type awsResourceLister struct {
	ec2Client         EC2Client
	elasticacheClient ElastiCacheClient
}

// NewAwsResourceLister creates a new ResourceLister that uses the provided EC2 client.
func NewAwsResourceLister(client EC2Client, elasticacheClient ElastiCacheClient) ResourceLister {
	return &awsResourceLister{
		ec2Client:         client,
		elasticacheClient: elasticacheClient,
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

	sort.Slice(vpcs, func(i, j int) bool {
		return vpcs[i].Name < vpcs[j].Name
	})

	return vpcs, nil
}

func (r *awsResourceLister) ListSubnets(ctx context.Context) ([]Subnet, error) {
	input := &ec2.DescribeSubnetsInput{}

	result, err := r.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %w", err)
	}

	subnetIds := make([]string, 0, len(result.Subnets))
	var subnets []Subnet
	for _, subnet := range result.Subnets {
		subnets = append(subnets, Subnet{
			Id:    *subnet.SubnetId,
			VpcId: *subnet.VpcId,
			Name:  *subnet.Tags[0].Value,
		})
		subnetIds = append(subnetIds, *subnet.SubnetId)
	}

	routeTables, err := r.ec2Client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: subnetIds,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list route tables: %w", err)
	}

	// Ick.
	for i, subnet := range subnets {
	rt:
		for _, routeTable := range routeTables.RouteTables {
			for _, association := range routeTable.Associations {
				if *association.SubnetId == subnet.Id {
					for _, route := range routeTable.Routes {
						if strings.HasPrefix(*route.GatewayId, "igw-") {
							subnets[i].Public = true
							break rt
						}
					}
				}
			}
		}
	}

	sort.Slice(subnets, func(i, j int) bool {
		return subnets[i].Id < subnets[j].Id
	})

	return subnets, nil
}

func (r *awsResourceLister) ListElastiCacheClusters(ctx context.Context) ([]ElastiCacheCluster, error) {
	input := &elasticache.DescribeCacheClustersInput{
		ShowCacheNodeInfo: aws.Bool(true),
	}

	result, err := r.elasticacheClient.DescribeCacheClusters(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list ElastiCache clusters: %w", err)
	}

	var clusters []ElastiCacheCluster
	for _, cluster := range result.CacheClusters {
		/*var nodes []ElastiCacheNode
		for _, node := range cluster.CacheNodes {
			nodes = append(nodes, ElastiCacheNode{
				Name:     *node.CacheNodeId,
				Endpoint: *node.Endpoint.Address,
			})
		}*/

		c := ElastiCacheCluster{
			Arn:  *cluster.ARN,
			Name: *cluster.CacheClusterId,
		}
		if cluster.ConfigurationEndpoint != nil && cluster.ConfigurationEndpoint.Address != nil {
			c.PrimaryEndpoint = *cluster.ConfigurationEndpoint.Address
		} else {
			c.PrimaryEndpoint = *cluster.CacheNodes[0].Endpoint.Address
		}

		clusters = append(clusters, c)
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name < clusters[j].Name
	})

	return clusters, nil
}
