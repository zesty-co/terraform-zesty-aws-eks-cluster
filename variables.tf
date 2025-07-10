variable "role_name" {
  description = "IAM role name"
  type        = string
  default     = "ZestyIamRole"
}

variable "policy_name" {
  description = "IAM policy name"
  type        = string
  default     = "ZestyPolicy"
}

variable "max_session_duration" {
  description = "Maximum session duration of the assumed role (in seconds)"
  type        = number
  default     = 43200
}

variable "products" {
  description = "List of all products to enable"
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
  description = "Path of values.yaml (default is the current working directory)"
  type        = string
  default     = "values.yaml"
}
