variable "subnet_ids" {
  description = "List of (private) subnet IDs for the Lambda function"
  type        = list(string)
}

variable "vpc_id" {
  description = "ID of the VPC to deploy into"
  type        = string
}

data "archive_file" "zip" {
  type        = "zip"
  source_dir  = "."
  output_path = "funciejs.zip"
  excludes    = [".terraform", "terraform.tfstate*", "*.tfvars", "deploy.tf", "README.md", "funciejs.zip", ".terraform.lock.hcl"]
}

resource "aws_security_group" "funcie_js_egress" {
  name        = "funcie-js-egress"
  description = "funcie-js-egress"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lambda_function" "funcie_js" {
  function_name    = "FuncieJsLambdaUrlSample"
  filename         = "funciejs.zip"
  handler          = "src/index.handler"
  source_code_hash = data.archive_file.zip.output_base64sha256
  role             = aws_iam_role.iam_for_lambda.arn
  runtime          = "nodejs18.x"
  memory_size      = 128
  timeout          = 30

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = aws_security_group.funcie_js_egress[*].id
  }

  environment {
    variables = {
      FUNCIE_LOG_LEVEL = "debug"
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

// Grant access to SSM and such:

resource "aws_iam_policy" "ssm_policy" {
  name        = "ssm_policy"
  description = "Allow access to SSM"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ssm_policy_attachment" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.ssm_policy.arn
}

resource "aws_lambda_function_url" "funcie_js" {
  function_name      = aws_lambda_function.funcie_js.function_name
  authorization_type = "NONE"
}

output "lambda" {
  value = aws_lambda_function.funcie_js.arn
}

output "lambda_url" {
  value = aws_lambda_function_url.funcie_js.function_url
}
