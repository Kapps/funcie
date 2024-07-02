
resource "aws_ecs_task_definition" "server_bastion_task" {
  family                   = "funcie-server-bastion"
  network_mode             = "bridge"
  requires_compatibilities = ["EC2"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = <<DEFINITION
  [
    {
      "name": "server-bastion-container",
      "image": "public.ecr.aws/w1h1o7p8/funcie-server-bastion:v${local.version}",
      "essential": true,

      "portMappings": [
        {
          "containerPort": 8082,
          "hostPort": 8082
        }
      ],
      "environment" : [
        { "name" : "FUNCIE_REDIS_ADDRESS", "value" : "${local.redis_host}" },
        { "name" : "FUNCIE_LISTEN_ADDRESS", "value" : "0.0.0.0:8082" },
        { "name" : "FUNCIE_LOG_LEVEL", "value" : "debug" },
        { "name" : "FUNCIE_VERSION", "value" : "${local.version}" },
        { "name" : "FUNCIE_ENV", "value" : "${var.funcie_env}" }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/funcie-server-bastion",
          "awslogs-region": "${var.region}",
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
  vpc_id      = local.vpc_id

  ingress {
    from_port   = 8082
    to_port     = 8082
    protocol    = "tcp"
    cidr_blocks = [local.vpc_cidr]
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
  launch_type     = "EC2"
}

resource "aws_iam_role" "ecs_execution_role" {
  name = "ecs_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "ecs_logging" {
  name = "ecs_logging"
  role = aws_iam_role.ecs_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "*"
      }
    ]
  })
}
