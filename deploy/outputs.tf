output "funcie_cluster_arn" {
  value = aws_ecs_cluster.funcie_cluster.arn
}

output "server_bastion_sg_id" {
  value = aws_security_group.server_bastion_sg.id
}

output "redis_host" {
  value = local.redis_host
}

output "bastion_host" {
  value = "${aws_service_discovery_service.server_bastion.name}.${aws_service_discovery_private_dns_namespace.funcie_local.name}"
}
