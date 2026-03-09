###############################################################################
# Zesty Kompass — single-cluster Helm deployment
###############################################################################

dependency "account" {
  config_path = "../account"
}

locals {
  cluster_name = "my-eks-cluster" # replace with your EKS cluster name
}

generate "backend" {
  path      = "backend.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
terraform {
  backend "s3" {
    bucket  = "my-terraform-state"
    key     = "zesty/kompass/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}
EOF
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
provider "aws" {
  region = "us-east-1"
}

data "aws_eks_cluster" "cluster" {
  name = "${local.cluster_name}"
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", "${local.cluster_name}"]
    }
  }
}
EOF
}

generate "main" {
  path      = "main.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
variable "kompass_values_yaml" {
  type      = string
  sensitive = true
}

resource "helm_release" "kompass" {
  name             = "kompass"
  repository       = "https://zesty-co.github.io/kompass"
  chart            = "kompass"
  namespace        = "zesty-system"
  cleanup_on_fail  = true
  create_namespace = true

  values = [var.kompass_values_yaml]
}
EOF
}

generate "versions" {
  path      = "versions.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
terraform {
  required_version = ">= 1.1"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }
}
EOF
}

inputs = {
  kompass_values_yaml = dependency.account.outputs.kompass_values_yaml
}
