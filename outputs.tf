output "kompass_values_yaml" {
  value       = local_file.kompass_values.content
  description = "The contents of the values.yaml file used to onboard Kompass"
}
