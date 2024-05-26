locals {
  redis_host = var.redis_host != "" ? var.redis_host : "${aws_elasticache_cluster.redis[0].cache_nodes.0.address}:${aws_elasticache_cluster.redis[0].port}"
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id = "funcie-redis"

  engine               = "redis"
  engine_version       = "7.1"
  node_type            = "cache.t4g.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.redis.name
  security_group_ids   = [aws_security_group.redis.id]

  tags = {
    Name = "funcie-redis"
  }

  count = var.redis_host == "" ? 1 : 0
}

resource "aws_elasticache_subnet_group" "redis" {
  name       = "funcie-redis"
  subnet_ids = var.private_subnet_ids
}

resource "aws_security_group" "redis" {
  name        = "funcie-redis"
  description = "Security group for the Redis cluster"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8"] # TODO: Should depend on the VPC CIDR
  }
}

resource "aws_ssm_parameter" "redis_host" {
  name  = "/funcie/${var.funcie_env}/redis_host"
  type  = "String"
  value = local.redis_host
}
