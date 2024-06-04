resource "aws_service_discovery_private_dns_namespace" "funcie_local" {
  name        = "funcie.local"
  description = "Service discovery namespace for the funcie to allow discovering of bastions and other services"
  vpc         = var.vpc_id
}

resource "aws_service_discovery_service" "server_bastion" {
  name          = "server-bastion"
  description   = "Service discovery service for connecting to the funcie server bastion"
  namespace_id  = aws_service_discovery_private_dns_namespace.funcie_local.id
  force_destroy = true

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.funcie_local.id
    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }

  #   depends_on = [null_resource.server_bastion_deregisterer]
}

# resource "null_resource" "server_bastion_deregisterer" {
#   triggers = {
#     service_id = aws_service_discovery_service.server_bastion.id
#     region     = var.region
#   }

#   provisioner "local-exec" {
#     when    = destroy
#     command = "AWS_PROFILE=${self.triggers.region} ${path.module}/deregister_service_map.sh ${self.triggers.service_id}"
#   }
# }
