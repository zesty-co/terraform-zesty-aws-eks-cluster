provider "aws" {
  region = var.region
}

# ---------------------------------------------------------------------------
# EKS cluster data — used by the Helm provider
# ---------------------------------------------------------------------------
locals {
  cluster_name = "kompass-eks-data"
}

data "aws_eks_cluster" "example" {
  name = local.cluster_name
}

data "aws_eks_cluster_auth" "example" {
  name = local.cluster_name
}

# Configure the Helm provider
provider "helm" {
  kubernetes = {
    host                   = data.aws_eks_cluster.example.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.example.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.example.token
  }
}
