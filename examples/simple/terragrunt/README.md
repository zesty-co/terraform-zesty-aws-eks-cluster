# Terragrunt Simple Example

Links your AWS account to Zesty and deploys the Kompass Helm chart into a **single EKS cluster** using Terragrunt.

## Directory Structure

```
account/    → IAM role + policy + zesty_account  (applied first)
kompass/    → Helm release                       (depends on account)
```

## Prerequisites

- Terraform >= 1.1
- Terragrunt >= 0.45
- AWS credentials configured
- Zesty API token — set the `token` value in the `generate "provider"` block inside `account/terragrunt.hcl`
- `kubectl` access to the target EKS cluster

## Usage

### Apply both at once

```bash
terragrunt run-all apply
```

### Apply individually

```bash
# 1. Account first
cd account && terragrunt apply

# 2. Kompass
cd ../kompass && terragrunt apply
```

## Configuration

- **Cluster name**: change `cluster_name` in `kompass/terragrunt.hcl` locals
- **S3 backend**: change bucket/key in the `generate "backend"` blocks
- **Region**: change region in the `generate "provider"` blocks
