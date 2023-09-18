data "local_file" "version" {
  filename = "${path.module}/../VERSION"
}

locals {
  version = trim(data.local_file.version.content, "\n")
}

output "version" {
  value = local.version
}
