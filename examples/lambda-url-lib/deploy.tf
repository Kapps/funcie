data "archive_file" "zip" {
  type        = "zip"
  source_file = "bin/main"
  output_path = "funciego.zip"
}

resource "aws_lambda_function" "funcie_go" {
  function_name    = "HandleRequest"
  filename         = "funciego.zip"
  handler          = "main"
  source_code_hash = "data.archive_file.zip.output_base64sha256"
  role             = aws_iam_role.iam_for_lambda.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 10
  vpc_config {
    subnet_ids         = subnet_ids
    security_group_ids = security_group_ids
  }
  environment {
    variables = {
      FUNCIE_REDIS_ADDR = redis_host
    }
  }
}

variable "subnet_ids" {
  type = set(string)
  default = []
}

variable "security_group_ids" {
  type = set(string)
  default = []
}

variable "redis_host" {
  type = string
  default = "localhost:6379"
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

output "lambda" {
  value = aws_lambda_function.funcie_go.arn
}