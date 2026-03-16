# Terraform module to onboard an AWS account and connect AWS EKS clusters to Zesty Kompass

Terraform module to onboard an AWS account and connect AWS EKS clusters to Zesty Kompass

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
  config = {
    bucket = "my-tf-state"
    key    = "zesty/account/terraform.tfstate"
    region = "us-east-1"
  }
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

### Simple (single cluster, two states)

> [Full example](./examples/simple/terragrunt/)

Uses the same `live/` directory hierarchy as Multi-Cluster, with a single `kompass/` stack.

Account (`zesty/account/terragrunt.hcl`):

```hcl
include "datacenter" {
  path = find_in_parent_folders("datacenter.hcl")
}

terraform {
  source = "tfr:///zesty-co/aws-eks-cluster/zesty?version=0.2.0"
}

generate "zesty_provider" {
  path      = "zesty_provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
provider "zesty" {
  token = "your-zesty-api-token" # replace with your actual token
}
EOF
}

inputs = {
  create_values_local_file = false
}
```

Kompass (`zesty/kompass/terragrunt.hcl`):

```hcl
include "datacenter" {
  path = find_in_parent_folders("datacenter.hcl")
}

dependency "account" {
  config_path = find_in_parent_folders("account/terragrunt.hcl")
}

locals {
  cluster_name = "my-eks-cluster"
}

terraform {
  source = "${get_repo_root()}/examples/simple/terragrunt/modules/kompass"
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
