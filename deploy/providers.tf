variable "region" {
  description = "The AWS region to deploy to"
  type        = string
}

terraform {
  required_providers {
    aws = {
      version = "> 2.0"
    }
    archive = {
      version = "> 1.3.0"
    }
  }
}

provider "aws" {
  region = var.region
}
