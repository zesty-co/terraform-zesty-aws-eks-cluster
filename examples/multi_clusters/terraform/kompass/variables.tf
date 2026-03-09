variable "cluster_name" {
  description = "Name of the EKS cluster to deploy Kompass into"
  type        = string
}

variable "region" {
  description = "AWS region where the EKS cluster lives"
  type        = string
  default     = "us-east-1"
}
