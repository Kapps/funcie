variable "subnet_ids" {
  description = "List of (private) subnet IDs for the Lambda function"
  type        = list(string)
}

variable "security_group_ids" {
  description = "List of security group IDs for the Lambda function"
  type        = list(string)
}

variable "redis_host" {
  description = "Address of the Redis host"
  type        = string
}

variable "bastion_lb_arn" {
  description = "The ARN for the load balancer created to be used for the bastion (generally listens on port 8082)."
  type        = string
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "bin/main"
  output_path = "funciego.zip"
}

data "aws_lb" "bastion_lb" {
  arn = var.bastion_lb_arn
}

resource "aws_lambda_function" "funcie_go" {
  function_name    = "FuncieGoLambdaUrlSample"
  filename         = "funciego.zip"
  handler          = "main"
  source_code_hash = data.archive_file.zip.output_base64sha256
  role             = aws_iam_role.iam_for_lambda.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 30
  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.security_group_ids
  }
  environment {
    variables = {
      FUNCIE_REDIS_ADDR              = var.redis_host,
      FUNCIE_SERVER_BASTION_ENDPOINT = "http://${data.aws_lb.bastion_lb.dns_name}:8082/dispatch",
      FUNCIE_APPLICATION_ID          = "url"
      FUNCIE_LOG_LEVEL               = "debug"
    }
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "terraform_lambda_policy" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_lambda_function_url" "funcie_go" {
  function_name      = aws_lambda_function.funcie_go.function_name
  authorization_type = "NONE"
}

output "lambda" {
  value = aws_lambda_function.funcie_go.arn
}

output "lambda_url" {
  value = aws_lambda_function_url.funcie_go.function_url
}
