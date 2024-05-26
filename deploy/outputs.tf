output "funcie_cluster_arn" {
  value = aws_ecs_cluster.funcie_cluster.arn
}

output "server_bastion_sg_id" {
  value = aws_security_group.server_bastion_sg.id
}

output "redis_host" {
  value = local.redis_host
}

output "bastion_private_key" {
  value     = tls_private_key.bastion_key.private_key_pem
  sensitive = true
}

output "bastion_public_key" {
  value = tls_private_key.bastion_key.public_key_openssh
}
