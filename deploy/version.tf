data "local_file" "version" {
  filename = "../VERSION"
}

locals {
  version = trim(data.local_file.version.content, "\n")
}

output "version" {
  value = local.version
}
