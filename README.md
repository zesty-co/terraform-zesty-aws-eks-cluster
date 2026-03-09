# Terraform module to connect an AWS EKS cluster to Zesty Kompass

This module onboards an AWS account to Zesty, and you can connect an
AWS EKS cluster to Zesty Kompass

## Prerequisites

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.1
- `token` for Zesty Provider - API token for Zesty platform, provided by a Zesty representative

  ```terraform
  provider "zesty" {
    token = "your-zesty-api-token" # replace with your actual token
  }
  ```

## Optional Provider

To connect a cluster (optional) as part of the installation, include:

- helm >= 3

## Examples

<details>
<summary><strong>Terraform</strong></summary>

### Simple (single cluster, one state file)

> [Full example](./examples/simple/terraform/)

```terraform
module "zesty" {
  source  = "zesty-co/aws-eks-cluster/zesty"
  version = "~> 0.2"

  create_values_local_file = false
}

resource "helm_release" "kompass" {
  name             = "kompass"
  repository       = "https://zesty-co.github.io/kompass"
  chart            = "kompass"
  namespace        = "zesty-system"
  cleanup_on_fail  = true
  create_namespace = true

  values = [module.zesty.kompass_values_yaml]
}
```

### Multi-Cluster (one account, N clusters, separate states)

> [Full example](./examples/multi_clusters/terraform/)

Account onboarding (`account/`):

```terraform
module "zesty" {
  source  = "zesty-co/aws-eks-cluster/zesty"
  version = "~> 0.2"

  create_values_local_file = false
}

output "kompass_values_yaml" {
  value     = module.zesty.kompass_values_yaml
  sensitive = true
}
```

Per-cluster Helm deploy (`kompass/`):

```terraform
data "terraform_remote_state" "account" {
  backend = "s3"
  config  = { bucket = "my-tf-state", key = "zesty/account/terraform.tfstate", region = "us-east-1" }
}

resource "helm_release" "kompass" {
  name             = "kompass"
  repository       = "https://zesty-co.github.io/kompass"
  chart            = "kompass"
  namespace        = "zesty-system"
  cleanup_on_fail  = true
  create_namespace = true

  values = [data.terraform_remote_state.account.outputs.kompass_values_yaml]
}
```

</details>

<details>
<summary><strong>Terragrunt</strong></summary>

### Simple (single cluster, one state file)

> [Full example](./examples/simple/terragrunt/)

Account (`account/terragrunt.hcl`):

```hcl
terraform {
  source = "tfr:///zesty-co/aws-eks-cluster/zesty?version=0.2.0"
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
    provider "aws" {
      region = "us-east-1"
    }
    provider "zesty" {
      token = "your-zesty-api-token"
    }
  EOF
}

inputs = {
  create_values_local_file = false
}
```

Kompass (`kompass/terragrunt.hcl`):

```hcl
dependency "account" {
  config_path = find_in_parent_folders("account/terragrunt.hcl")
}

inputs = {
  kompass_values_yaml = dependency.account.outputs.kompass_values_yaml
}
```

### Multi-Cluster (production-grade live hierarchy)

> [Full example](./examples/multi_clusters/terragrunt/)

Uses a `live/` directory layout: environment / datacenter / region / account.
One `account/` stack is applied once, then each EKS cluster gets its own `kompass-<cluster>/` stack that reads the account output via `dependency`.

Per-cluster Kompass (`kompass-eks-prod/terragrunt.hcl`):

```hcl
include "datacenter" {
  path = find_in_parent_folders("datacenter.hcl")
}

dependency "account" {
  config_path = find_in_parent_folders("account/terragrunt.hcl")
}

locals {
  cluster_name = "eks-prod"
}

terraform {
  source = "${get_repo_root()}/examples/multi_clusters/terragrunt/modules/kompass"
}

inputs = {
  kompass_values_yaml = dependency.account.outputs.kompass_values_yaml
}
```

</details>

## Kompass Helm Values Reference

The module outputs a `kompass_values_yaml` string containing the credentials and
metadata needed to connect your cluster to Zesty. It is passed directly to the
`helm_release` resource via the `values` argument.

For the full list of configurable chart values, see the
[Kompass `values.yaml`](https://github.com/zesty-co/kompass-insights/blob/main/charts/zesty/values.yaml).


<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.1 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 6.0 |
| <a name="requirement_local"></a> [local](#requirement\_local) | ~> 2.5.3 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 3.7.2 |
| <a name="requirement_zesty"></a> [zesty](#requirement\_zesty) | ~> 0.3.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.35.1 |
| <a name="provider_local"></a> [local](#provider\_local) | 2.5.3 |
| <a name="provider_random"></a> [random](#provider\_random) | 3.7.2 |
| <a name="provider_zesty"></a> [zesty](#provider\_zesty) | 0.3.0 |

## Resources

| Name | Type |
|------|------|
| [aws_iam_role.zesty_iam_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.zesty_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [local_file.kompass_values](https://registry.terraform.io/providers/hashicorp/local/latest/docs/resources/file) | resource |
| [random_uuid.zesty_external_id](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/uuid) | resource |
| [zesty_account.result](https://registry.terraform.io/providers/zesty-co/zesty/latest/docs/resources/account) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_create_values_local_file"></a> [create\_values\_local\_file](#input\_create\_values\_local\_file) | Enables the creation of a local values.yaml file | `bool` | `true` | no |
| <a name="input_max_session_duration"></a> [max\_session\_duration](#input\_max\_session\_duration) | Maximum session duration of the assumed role (in seconds) | `number` | `43200` | no |
| <a name="input_policy_name"></a> [policy\_name](#input\_policy\_name) | IAM policy name | `string` | `"ZestyPolicy"` | no |
| <a name="input_products"></a> [products](#input\_products) | List of all products to enable | `list(map(any))` | <pre>[<br/>  {<br/>    "active": true,<br/>    "name": "Kompass"<br/>  }<br/>]</pre> | no |
| <a name="input_region"></a> [region](#input\_region) | AWS region | `string` | `""` | no |
| <a name="input_role_name"></a> [role\_name](#input\_role\_name) | IAM role name | `string` | `"ZestyIamRole"` | no |
| <a name="input_trusted_principal"></a> [trusted\_principal](#input\_trusted\_principal) | Trusted AWS principal allowed to assume the role | `string` | `"arn:aws:iam::672188301118:root"` | no |
| <a name="input_values_yaml_filename"></a> [values\_yaml\_filename](#input\_values\_yaml\_filename) | Path of values.yaml (default is the current working directory) | `string` | `"values.yaml"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_kompass_values_yaml"></a> [kompass\_values\_yaml](#output\_kompass\_values\_yaml) | The contents of the values.yaml file used to onboard Kompass |
<!-- END_TF_DOCS -->