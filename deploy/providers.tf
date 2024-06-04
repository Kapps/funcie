terraform {
  required_providers {
    aws = {
      version = "> 2.0"
    }
    archive = {
      version = "> 1.3.0"
    }
    tls = {
      version = "> 4.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
}

provider "tls" {

}

provider "null" {

}
