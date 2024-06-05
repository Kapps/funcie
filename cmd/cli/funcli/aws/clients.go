package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SsmClient is a minimal interface for the AWS SSM client with only the methods used by the CLI.
type SsmClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)
	TerminateSession(ctx context.Context, params *ssm.TerminateSessionInput, optFns ...func(*ssm.Options)) (*ssm.TerminateSessionOutput, error)
}

// EC2Client is a minimal interface for the EC2 client with only the methods used by the CLI.
type EC2Client interface {
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, opts ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, opts ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
	DescribeRouteTables(ctx context.Context, params *ec2.DescribeRouteTablesInput, opts ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)
}

// ElastiCacheClient is a minimal interface for the ElastiCache client with only the methods used by the CLI.
type ElastiCacheClient interface {
	DescribeCacheClusters(ctx context.Context, params *elasticache.DescribeCacheClustersInput, optFns ...func(*elasticache.Options)) (*elasticache.DescribeCacheClustersOutput, error)
}
