package aws_test

import (
	"context"
	. "github.com/Kapps/funcie/cmd/funcie/funcli/aws"
	"github.com/Kapps/funcie/cmd/funcie/funcli/aws/mocks"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/stretchr/testify/require"
)

func TestListVpcs(t *testing.T) {
	ctx := context.Background()
	mockEC2 := mocks.NewEC2Client(t)
	mockElastiCache := mocks.NewElastiCacheClient(t)
	r := NewAwsResourceLister(mockEC2, mockElastiCache)

	vpcsOutput := &ec2.DescribeVpcsOutput{
		Vpcs: []types.Vpc{
			{
				VpcId: aws.String("vpc-2"),
				Tags:  []types.Tag{{Key: aws.String("Name"), Value: aws.String("beta-vpc")}},
			},
			{
				VpcId: aws.String("vpc-1"),
				Tags:  []types.Tag{{Key: aws.String("Name"), Value: aws.String("alpha-vpc")}},
			},
		},
	}

	mockEC2.EXPECT().DescribeVpcs(ctx, &ec2.DescribeVpcsInput{}).
		Return(vpcsOutput, nil).
		Once()

	vpcs, err := r.ListVpcs(ctx)

	require.NoError(t, err)
	require.Len(t, vpcs, 2)
	require.Equal(t, "vpc-1", vpcs[0].Id)
	require.Equal(t, "alpha-vpc", vpcs[0].Name)
	require.Equal(t, "vpc-2", vpcs[1].Id)
	require.Equal(t, "beta-vpc", vpcs[1].Name)
}

func TestListSubnets(t *testing.T) {
	ctx := context.Background()
	mockEC2 := mocks.NewEC2Client(t)
	mockElastiCache := mocks.NewElastiCacheClient(t)
	r := NewAwsResourceLister(mockEC2, mockElastiCache)

	subnetsOutput := &ec2.DescribeSubnetsOutput{
		Subnets: []types.Subnet{
			{
				SubnetId: aws.String("subnet-2"),
				VpcId:    aws.String("vpc-1"),
				Tags:     []types.Tag{{Key: aws.String("Name"), Value: aws.String("beta-subnet")}},
			},
			{
				SubnetId: aws.String("subnet-1"),
				VpcId:    aws.String("vpc-1"),
				Tags:     []types.Tag{{Key: aws.String("Name"), Value: aws.String("alpha-subnet")}},
			},
		},
	}

	mockEC2.EXPECT().DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{}).
		Return(subnetsOutput, nil).
		Once()

	routeTablesOutput := &ec2.DescribeRouteTablesOutput{
		RouteTables: []types.RouteTable{
			{
				Associations: []types.RouteTableAssociation{
					{SubnetId: aws.String("subnet-1")},
				},
				Routes: []types.Route{
					{GatewayId: aws.String("igw-1")},
				},
			},
			{
				Associations: []types.RouteTableAssociation{
					{SubnetId: aws.String("subnet-2")},
				},
				Routes: []types.Route{
					{GatewayId: aws.String("igw-2")},
				},
			},
		},
	}

	mockEC2.EXPECT().DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []string{"subnet-2", "subnet-1"},
			},
		},
	}).Return(routeTablesOutput, nil).Once()

	subnets, err := r.ListSubnets(ctx)

	require.NoError(t, err)
	require.Len(t, subnets, 2)
	require.Equal(t, "subnet-1", subnets[0].Id)
	require.Equal(t, "vpc-1", subnets[0].VpcId)
	require.Equal(t, "alpha-subnet", subnets[0].Name)
	require.True(t, subnets[0].Public)
	require.Equal(t, "subnet-2", subnets[1].Id)
	require.Equal(t, "vpc-1", subnets[1].VpcId)
	require.Equal(t, "beta-subnet", subnets[1].Name)
	require.True(t, subnets[1].Public)

	mockEC2.AssertExpectations(t)
}

func TestListElastiCacheClusters(t *testing.T) {
	ctx := context.Background()
	mockEC2 := mocks.NewEC2Client(t)
	mockElastiCache := mocks.NewElastiCacheClient(t)
	r := NewAwsResourceLister(mockEC2, mockElastiCache)

	clustersOutput := &elasticache.DescribeCacheClustersOutput{
		CacheClusters: []ecTypes.CacheCluster{
			{
				ARN:            aws.String("arn:aws:elasticache:cluster:beta-cluster"),
				CacheClusterId: aws.String("beta-cluster"),
				CacheNodes: []ecTypes.CacheNode{
					{Endpoint: &ecTypes.Endpoint{Address: aws.String("beta-endpoint")}},
				},
				ConfigurationEndpoint: &ecTypes.Endpoint{Address: aws.String("beta-config-endpoint")},
			},
			{
				ARN:            aws.String("arn:aws:elasticache:cluster:alpha-cluster"),
				CacheClusterId: aws.String("alpha-cluster"),
				CacheNodes: []ecTypes.CacheNode{
					{Endpoint: &ecTypes.Endpoint{Address: aws.String("alpha-endpoint")}},
				},
				ConfigurationEndpoint: &ecTypes.Endpoint{Address: aws.String("alpha-config-endpoint")},
			},
		},
	}

	mockElastiCache.EXPECT().DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		ShowCacheNodeInfo: aws.Bool(true),
	}).Return(clustersOutput, nil).Once()

	clusters, err := r.ListElastiCacheClusters(ctx)

	require.NoError(t, err)
	require.Len(t, clusters, 2)
	require.Equal(t, "arn:aws:elasticache:cluster:alpha-cluster", clusters[0].Arn)
	require.Equal(t, "alpha-cluster", clusters[0].Name)
	require.Equal(t, "alpha-config-endpoint", clusters[0].PrimaryEndpoint)
	require.Equal(t, "arn:aws:elasticache:cluster:beta-cluster", clusters[1].Arn)
	require.Equal(t, "beta-cluster", clusters[1].Name)
	require.Equal(t, "beta-config-endpoint", clusters[1].PrimaryEndpoint)

	mockElastiCache.AssertExpectations(t)
}
