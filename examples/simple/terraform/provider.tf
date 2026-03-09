provider "aws" {
  region = "us-east-1"
}

provider "zesty" {
  token = "your-zesty-api-token" # replace with your actual token
}

# ---------------------------------------------------------------------------
# EKS cluster data — used by the Helm provider
# ---------------------------------------------------------------------------
data "aws_eks_cluster" "cluster" {
  name = var.cluster_name
}

provider "helm" {
  kubernetes = {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
    exec = {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
    }
  }
}
