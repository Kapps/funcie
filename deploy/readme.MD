# Terraform Module Overview:

## Variables:

- `subnet_ids`: IDs for the private subnets where load balancer resources will be created.
- `redis_host`: Address for the Redis host that both the server and client bastion can access, including port.
- `vpc_id`: ID for the VPC where resources will be created.
- `public_subnet_ids`: IDs for public subnets where the EC2 instance will be created. Can be null if using a private IP.
- `assign_public_ip`: Whether to use a public IP and public subnets for the server bastion; may require a NAT Gateway otherwise.

## Resources:

## High Level Overview:

- Creates an ECS cluster named `funcie-cluster` that uses a Fargate service for launching the server bastion, logging to CloudWatch.
  - For production usage, the service should be within a private subnet with no public IP address and behind a NAT Gateway.
  - By default, it is a public subnet with a public IP address, though it is only accessible from the VPC CIDR block.
- Creates an internal load balancer that listens on port 8082 for directing traffic to the Fargate instance.

### AWS ECS:

- Creates an ECS cluster named `funcie-cluster`.
- Creates an ECS task definition for the "funcie-server-bastion" with FARGATE compatibility.
  - Uses AWS logs for logging to CloudWatch (see below section).
- Creates an ECS service named "funcie-server-bastion-service" with FARGATE launch type.
  - This service is associated with the previously defined task and uses a public IP by default with a security group for allowing internal access.
  - Note: There is a TODO for making this private behind a NAT Gateway.

### AWS CloudWatch:

- Sets up a log group for the ECS task to log its output: "/ecs/funcie-server-bastion".

### AWS Security Group:

- Creates a security group allowing TCP traffic on port 8082 from the CIDR block "10.0.0.0/8".
- Allows all outbound traffic.

### AWS IAM:

- Sets up an IAM role for the task that allows creating log stream and putting log events.

### AWS Load Balancer:

- Sets up a network load balancer (internal type).
- Configures a target group for the bastion service.
- Creates a listener for directing incoming traffic to the target group on port 8082.

## Caveats and Warnings:

- The ECS service is configured by default to have a public IP. This is intended for demonstration purposes and should be changed for production usage.
  - The ECS service should ideally be in a private subnet, and optionally behind a NAT Gateway, with no public IP address.
  - That being said, the ECS service is configured to only allow traffic from the VPC CIDR block via the security group.
- Security groups are allowing specific traffic, ensure that they match the security needs of your environment.
- Ensure that IAM roles and permissions are reviewed to match least privilege principles.
