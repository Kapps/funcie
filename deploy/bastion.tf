variable "subnet_ids" {
  type = set(string)
}

variable "redis_host" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "public_subnet_ids" {
  type = set(string)
}

resource "aws_ecs_cluster" "funcie_cluster" {
  name = "funcie-cluster"
}

resource "aws_ecs_task_definition" "server_bastion_task" {
  family                   = "funcie-server-bastion"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = <<DEFINITION
  [
    {
      "name": "server-bastion-container",
      "image": "public.ecr.aws/w1h1o7p8/funcie-server-bastion:8e4df269d18a3bff7994c911a3df0caee3cab317",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8082,
          "hostPort": 8082
        }
      ],
      "environment" : [
        { "name" : "FUNCIE_REDIS_ADDRESS", "value" : "${var.redis_host}" },
        { "name" : "FUNCIE_LISTEN_ADDRESS", "value" : "0.0.0.0:8082" },
        { "name" : "FUNCIE_LOG_LEVEL", "value" : "debug" }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/funcie-server-bastion",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
  DEFINITION
}

resource "aws_cloudwatch_log_group" "funcie_server_bastion_lg" {
  name = "/ecs/funcie-server-bastion"
}

resource "aws_security_group" "server_bastion_sg" {
  name        = "funcie-server-bastion-sg"
  description = "funcie-server-bastion-sg"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 8082
    to_port     = 8082
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_service" "server_bastion_service" {
  name            = "funcie-server-bastion-service"
  cluster         = aws_ecs_cluster.funcie_cluster.id
  task_definition = aws_ecs_task_definition.server_bastion_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    # TODO: This should not be public in a deploy but rather behind a NAT Gateway with no public IP address.
    # It should be in a private subnet, but that's beyond the scope of the example.
    assign_public_ip = true
    subnets          = var.public_subnet_ids
    security_groups  = [aws_security_group.server_bastion_sg.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.funcie_server_bastion_tg.arn
    container_name   = "server-bastion-container"
    container_port   = 8082
  }
}

resource "aws_iam_role" "ecs_execution_role" {
  name = "ecs_execution_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ecs_logging" {
  name = "ecs_logging"
  role = aws_iam_role.ecs_execution_role.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_lb" "funcie_lb" {
  name               = "funcie-lb"
  internal           = true
  load_balancer_type = "network"
  subnets            = var.subnet_ids
}

resource "aws_lb_target_group" "funcie_server_bastion_tg" {
  name     = "funcie-server-bastion-tg"
  port     = 8082
  protocol = "TCP"
  vpc_id   = var.vpc_id

  target_type = "ip"
}

resource "aws_lb_listener" "funcie_server_bastion_listener" {
  load_balancer_arn = aws_lb.funcie_lb.arn
  port              = 8082
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.funcie_server_bastion_tg.arn
  }
}
