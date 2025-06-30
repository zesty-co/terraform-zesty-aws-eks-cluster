variable "role_name" {
  description = "Name of the IAM role"
  type        = string
  default     = "ZestyIamRole"
}

variable "policy_name" {
  description = "Name of the IAM policy"
  type        = string
  default     = "ZestyPolicy"
}

variable "max_session_duration" {
  description = "Assumed role max session duration in seconds"
  type        = number
  default     = 43200
}

variable "cloud_provider" {
  description = "Name of the cloud provider"
  type        = string
  default     = "AWS"
}

variable "products" {
  description = "List of all enabled products"
  type        = list(map(any))
  default = [{
    name   = "Kompass"
    active = true
  }]
}

variable "trusted_principal" {
  description = "Trusted AWS principal allowed to assume the role"
  type        = string
  default     = "arn:aws:iam::672188301118:root"
}

variable "values_yaml_filename" {
  description = "Location of the Kompass values.yaml output"
  type        = string
  default     = "values.yaml"
}
