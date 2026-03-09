###############################################################################
# Zesty AWS Account
#
# This configuration calls the current zesty-co/aws-eks-cluster/zesty module
# to create the IAM role, policy, and register the AWS account with Zesty.
#
# Applied ONCE per AWS account.
# The kompass_values_yaml output is consumed by the kompass/ example via
# terraform_remote_state.
###############################################################################

terraform {
  backend "s3" {
    bucket  = "my-terraform-state"
    key     = "zesty/account/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}

module "zesty" {
  source = "../../../../"

  role_name                = "ZestyIamRole"
  policy_name              = "ZestyPolicy"
  create_values_local_file = false
}
