# TODO: Move this into a module.
# Gets a bit tricky with needing the VPC for creating the instance but needing the instance for routing.

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
  vpc_id                  = aws_vpc.funcie_vpc.id
  cidr_block              = var.public_subnet_cidrs[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "funcie-public-subnet-${count.index}"
  }
}

resource "aws_internet_gateway" "funcie_igw" {
  count = var.vpc_id == "" ? 1 : 0
  vpc_id = aws_vpc.funcie_vpc.id

  tags = {
    Name = "funcie-igw"
  }
}

resource "aws_subnet" "funcie_private_subnets" {
  count                   = var.vpc_id == "" ? length(var.private_subnet_cidrs) : 0
  vpc_id                  = aws_vpc.funcie_vpc.id
  cidr_block              = var.private_subnet_cidrs[count.index]
  map_public_ip_on_launch = false

  tags = {
    Name = "funcie-private-subnet-${count.index}"
  }
}

resource "aws_route_table" "funcie_igw_route_table" {
  count  = var.vpc_id == "" ? 1 : 0
  vpc_id = aws_vpc.funcie_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.funcie_igw.id
  }
}

resource "aws_route_table" "funcie_nat_route_table" {
  count  = var.vpc_id == "" ? 1 : 0
  vpc_id = aws_vpc.funcie_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    instance_id = aws_instance.funcie_nat.id
