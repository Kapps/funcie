package aws

import "fmt"

// Vpc provides some basic information about a VPC.
type Vpc struct {
	// ID is the AWS representation of the VPC ID (not the ARN).
	Id string
	// Name is a human-readable name for the VPC.
	Name string
}

func (v Vpc) String() string {
	return fmt.Sprintf("%s (%s)", v.Name, v.Id)
}

// Subnet provides some basic information about a Subnet.
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

func (s Subnet) String() string {
	pubStr := ""
	if s.Public {
		pubStr = "Public"
	} else {
		pubStr = "Private"
	}
	return fmt.Sprintf("%s: %s (%s)", s.Name, s.Id, pubStr)
}

// ElastiCacheCluster represents an AWS ElastiCache cluster and its nodes.
type ElastiCacheCluster struct {
	// Arn is the ARN of the ElastiCache cluster, such as `arn:aws:elasticache:ca-central-1:123:cluster:foo`.
	Arn string
	// Name is the human-readable name of the ElastiCache cluster.
	Name string
	/*// Nodes is a list of nodes in the ElastiCache cluster.
	Nodes []ElastiCacheNode*/
	// PrimaryEndpoint is the primary endpoint for the ElastiCache cluster.
	// If the cluster is not configured as a cluster, this is the endpoint for the single node.
	PrimaryEndpoint string
}

func (e ElastiCacheCluster) String() string {
	return fmt.Sprintf("%s (%s)", e.Name, e.PrimaryEndpoint)
}

// ElastiCacheNode represents a node in an AWS ElastiCache cluster.
type ElastiCacheNode struct {
	// Name is the human-readable name of the node.
	Name string
	// Endpoint is the endpoint of the node.
	Endpoint string
}
