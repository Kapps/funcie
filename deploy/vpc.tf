# TODO: Move this into a module.

locals {
  public_subnet_cidrs  = var.vpc_id != "" ? var.public_subnet_cidrs : aws_subnet.funcie_public_subnets[*].cidr_block
  public_subnet_ids    = var.vpc_id != "" ? var.public_subnet_ids : aws_subnet.funcie_public_subnets[*].id
  private_subnet_cidrs = var.vpc_id != "" ? var.private_subnet_cidrs : aws_subnet.funcie_private_subnets[*].cidr_block
  private_subnet_ids   = var.vpc_id != "" ? var.private_subnet_ids : aws_subnet.funcie_private_subnets[*].id

  vpc_id   = var.vpc_id != "" ? var.vpc_id : aws_vpc.funcie_vpc[0].id
  vpc_cidr = data.aws_vpc.funcie_vpc.cidr_block
}

data "aws_vpc" "funcie_vpc" {
  id = local.vpc_id
}

resource "aws_vpc" "funcie_vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr

  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "funcie-vpc"
  }
}

resource "aws_subnet" "funcie_public_subnets" {
  count                   = var.vpc_id == "" ? length(var.public_subnet_cidrs) : 0
  vpc_id                  = aws_vpc.funcie_vpc[0].id
  cidr_block              = var.public_subnet_cidrs[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "funcie-public-subnet-${count.index}"
  }
}

resource "aws_subnet" "funcie_private_subnets" {
  count                   = var.vpc_id == "" ? length(var.private_subnet_cidrs) : 0
  vpc_id                  = aws_vpc.funcie_vpc[0].id
  cidr_block              = var.private_subnet_cidrs[count.index]
  map_public_ip_on_launch = false

  tags = {
    Name = "funcie-private-subnet-${count.index}"
  }
}

resource "aws_internet_gateway" "funcie_igw" {
  count  = var.vpc_id == "" ? 1 : 0
  vpc_id = aws_vpc.funcie_vpc[0].id

  tags = {
    Name = "funcie-igw"
  }
}

resource "aws_route_table" "funcie_igw_route_table" {
  count  = var.vpc_id == "" ? 1 : 0
  vpc_id = aws_vpc.funcie_vpc[0].id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.funcie_igw[0].id
  }

  tags = {
    Name = "funcie-igw-route-table"
  }
}

resource "aws_route_table" "funcie_nat_route_table" {
  count  = var.vpc_id == "" ? 1 : 0
  vpc_id = aws_vpc.funcie_vpc[0].id

  // The route here is created by the userdata script in the ASG.
  // This is since the NAT instance may change if the instance becomes unhealthy or terminates.

  tags = {
    Name = "funcie-nat-route-table"
  }
}

resource "aws_route_table_association" "funcie_public_subnet_associations" {
  count          = var.vpc_id == "" ? length(var.public_subnet_cidrs) : 0
  subnet_id      = aws_subnet.funcie_public_subnets[count.index].id
  route_table_id = aws_route_table.funcie_igw_route_table[0].id
}

resource "aws_route_table_association" "funcie_private_subnet_associations" {
  count          = var.vpc_id == "" ? length(var.private_subnet_cidrs) : 0
  subnet_id      = aws_subnet.funcie_private_subnets[count.index].id
  route_table_id = aws_route_table.funcie_nat_route_table[0].id
}

output "vpc_id" {
  value = var.vpc_id == "" ? aws_vpc.funcie_vpc[0].id : var.vpc_id
}

output "public_subnet_ids" {
  value = local.public_subnet_ids
}

output "public_subnet_cidrs" {
  value = local.public_subnet_cidrs
}

output "private_subnet_ids" {
  value = local.private_subnet_ids
}

output "private_subnet_cidrs" {
  value = local.private_subnet_cidrs
}
