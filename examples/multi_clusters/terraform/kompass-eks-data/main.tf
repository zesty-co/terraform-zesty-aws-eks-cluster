###############################################################################
# Zesty Kompass — per-cluster Helm deployment
#
# Deploy this ONCE PER EKS CLUSTER.
# Each cluster gets its own directory copy with its own state file and its own
# var.cluster_name value.
#
# Example directory layout for multi-cluster:
#
#   infrastructure/
#   ├── account/               # applied once per AWS account
#   ├── kompass-eks-prod/      # copy of this directory — cluster_name = "eks-prod"
#   ├── kompass-eks-staging/   # copy of this directory — cluster_name = "eks-staging"
#   └── kompass-eks-data/      # copy of this directory — cluster_name = "eks-data"
#
# Each copy only differs in backend key and cluster_name variable.
###############################################################################

terraform {
  backend "s3" {
    bucket  = "my-terraform-state"
    key     = "zesty/kompass-eks-prod/terraform.tfstate" # change per cluster
    region  = "us-east-1"
    encrypt = true
  }
}

# ---------------------------------------------------------------------------
# Read Kompass values from the account state
# ---------------------------------------------------------------------------
data "terraform_remote_state" "account" {
  backend = "s3"
  config = {
    bucket = "my-terraform-state"
    key    = "zesty/account/terraform.tfstate"
    region = "us-east-1"
  }
}

# ---------------------------------------------------------------------------
# Kompass Helm release
# ---------------------------------------------------------------------------
resource "helm_release" "kompass" {
  name             = "kompass"
  repository       = "https://zesty-co.github.io/kompass"
  chart            = "kompass"
  namespace        = "zesty-system"
  cleanup_on_fail  = true
  create_namespace = true

  values = [data.terraform_remote_state.account.outputs.kompass_values_yaml]
}