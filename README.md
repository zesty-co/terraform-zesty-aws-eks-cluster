# Terraform module for connecting an AWS EKS cluster to Zesty Kompass

This module allows onboarding to Zesty in order to connect Zesty Kompass
to an AWS EKS cluster

## Prerequisites

- [Terraform](https://developer.hashicorp.com/terraform/install) 0.13+

## Providers

- aws >= 5.0
- random >= 3.7.2
- local >= 2.5.3

## Optional Provider

If you wish to onboard the cluster itself via Terraform, you can include:

- helm >= 3

## Example Usage

```terraform

module "zesty" {
  source              = "zesty-co/eks-cluster/zesty-co"
  aws_region          = "{{your aws region}}"
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
