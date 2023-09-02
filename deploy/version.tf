data "local_file" "version" {
  filename = "../VERSION"
}

output "version" {
  value = data.local_file.version.content
}
