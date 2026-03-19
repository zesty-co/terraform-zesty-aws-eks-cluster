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

### Terraform

- [Simple single-cluster example](./examples/simple/terraform/)
- [Multi-cluster example](./examples/multi_clusters/terraform/)

### Terragrunt

- [Simple single-cluster example](./examples/simple/terragrunt/)
- [Multi-cluster example](./examples/multi_clusters/terragrunt/)


## Kompass Helm Values Reference

The module outputs a `kompass_values_yaml` string containing the credentials and
metadata needed to connect your cluster to Zesty. It is passed directly to the
`helm_release` resource via the `values` argument.

For the full list of configurable chart values, see the
[Kompass `values.yaml`](https://github.com/zesty-co/kompass-insights/blob/main/charts/zesty/values.yaml).
