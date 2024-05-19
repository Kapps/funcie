variable "redis_host" {
  description = "Address of the Redis host, including trailing :6379 port."
  type        = string
}

variable "vpc_id" {
  description = "ID of the VPC to deploy into."
  type        = string
}

variable "public_subnet_ids" {
  description = "IDs of the public subnets to deploy into."
  type        = set(string)
}

variable "private_subnet_ids" {
  description = "IDs of the private subnets to deploy into."
  type        = set(string)
}

variable "bastion_public_key_path" {
  description = "Path to the SSH key to use for the bastion host."
  type        = string
  default     = "~/.ssh/funcie_rsa.pub"
}

variable "has_vpn" {
  description = "Whether the VPC has a VPN connection and therefore should use only private networking and not include SSH tunneling. Not yet implemented."
  type        = bool
  default     = false
}
