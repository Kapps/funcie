############################################
# Variables that are always required to deploy funcie.
############################################

variable "region" {
  description = "The AWS region to deploy to, such as us-east-1 or ca-central-1."
  type        = string
}

variable "redis_host" {
  description = "Address of the Redis host, including trailing :6379 port. If empty, an Elasticache instance will be created and used."
  type        = string
  default     = ""
}

variable "vpc_id" {
  description = "ID of the VPC to deploy into, or an empty string to create a new VPC (in which case )"
  type        = string
  default     = ""
}

############################################
# Variables that are required to deploy funcie into an existing VPC.
############################################

variable "public_subnet_ids" {
  description = "IDs of the public subnets to deploy into. Must be set if deploying into an existing VPC; ignored otherwise."
  type        = list(string)
  default     = []

  // Terraform 1.9 will allow this, but it's still too new to require for such a minor benefit.
  /*validation {
    condition     = var.vpc_id == "" || length(var.public_subnet_ids) > 0
    error_message = "At least one public subnet ID must be provided if a vpc_id is set."
  }*/
}

variable "private_subnet_ids" {
  description = "IDs of the private subnets to deploy into. Must be set if deploying into an existing VPC; ignored otherwise."
  type        = list(string)
  default     = []

  // Terraform 1.9 will allow this, but it's still too new to require for such a minor benefit.
  /*validation {
    condition     = var.vpc_id == "" || length(var.private_subnet_ids) > 0
    error_message = "At least one private subnet ID must be provided if a vpc_id is set."
  }*/
}

############################################
# Variables that are used to deploy funcie into a new VPC (all of which have valid defaults).
############################################

variable "public_subnet_cidrs" {
  description = "If having funcie create a VPC, the CIDR blocks for the public subnets."
  type        = list(string)
  default = [
    "10.0.1.0/24",
    "10.0.2.0/24",
  ]
}

variable "private_subnet_cidrs" {
  description = "If having funcie create a VPC, the CIDR blocks for the private subnets."
  type        = list(string)
  default = [
    "10.0.128.0/24",
    "10.0.129.0/24",
  ]
}

variable "vpc_cidr" {
  description = "If having funcie create a VPC, the CIDR block for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

############################################
# Other optional variables.
############################################

variable "funcie_env" {
  description = "An environment name to allow differentiating deployments and configuration. If you don't need multiple environments, you can use the default value. This is not yet implemented and must remain default for now."
  type        = string
  default     = "default"
}

variable "redis_instance_type" {
  description = "The instance type to use for the Elasticache instance. If redis_host is set, this is ignored."
  type        = string
  default     = "cache.t4g.micro"
}

variable "bastion_instance_type" {
  description = "The instance type to use for the bastion host."
  type        = string
  default     = "t3.micro"
}