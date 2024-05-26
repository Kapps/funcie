variable "region" {
  description = "The AWS region to deploy to"
  type        = string
}

variable "redis_host" {
  description = "Address of the Redis host, including trailing :6379 port. If empty, an Elasticache instance will be created and used."
  type        = string
  default     = ""
}

variable "vpc_id" {
  description = "ID of the VPC to deploy into."
  type        = string
}

variable "public_subnet_ids" {
  description = "IDs of the public subnets to deploy into."
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "IDs of the private subnets to deploy into."
  type        = list(string)
}

variable "has_vpn" {
  description = "Whether the VPC has a VPN connection and therefore should use only private networking and not include SSH tunneling. Not yet implemented."
  type        = bool
  default     = false
}

variable "funcie_env" {
  description = "An environment name to allow differentiating deployments and configuration. If you don't need multiple environments, you can use the default value."
  type        = string
  default     = "default"
}
