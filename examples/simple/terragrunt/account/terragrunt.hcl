###############################################################################
# Zesty AWS Account — applied ONCE per AWS account
###############################################################################

terraform {
  source = "tfr:///zesty-co/aws-eks-cluster/zesty?version=0.2.0"
}

generate "backend" {
  path      = "backend.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
terraform {
  backend "s3" {
    bucket  = "my-terraform-state"
    key     = "zesty/account/terraform.tfstate"
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

provider "zesty" {
  token = "your-zesty-api-token" # replace with your actual token
}
EOF
}

inputs = {
  role_name                = "ZestyIamRole"
  policy_name              = "ZestyPolicy"
  create_values_local_file = false
}
