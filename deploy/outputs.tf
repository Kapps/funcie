
output "server_bastion_lb_host" {
  value = aws_lb.funcie_lb.dns_name
}

output "server_bastion_lb_arn" {
  value = aws_lb.funcie_lb.arn
}

output "funcie_cluster_arn" {
  value = aws_ecs_cluster.funcie_cluster.arn
}

output "server_bastion_sg_id" {
  value = aws_security_group.server_bastion_sg.id
}

output "server_bastion_host" {
  value = aws_instance.server_bastion.public_ip
}