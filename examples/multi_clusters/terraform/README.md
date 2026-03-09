# Terraform Multi-Cluster Example

This example deploys Zesty across **one AWS account** and **multiple EKS clusters** using plain Terraform with S3 remote state.

## Architecture

```
account/       → applied ONCE per AWS account  (IAM role + policy + zesty_account)
kompass/       → applied ONCE PER EKS cluster  (Helm release)
```

Each cluster gets its own copy of the `kompass/` directory with a unique backend key and `cluster_name`.

```
infrastructure/
├── account/                   # one per AWS account — own state
├── kompass-eks-prod/          # cluster: eks-prod   — own state
├── kompass-eks-staging/       # cluster: eks-staging — own state
└── kompass-eks-data/          # cluster: eks-data   — own state
```

## Prerequisites

- Terraform >= 1.1
- AWS credentials configured
- Zesty API token — set the `token` value in `account/provider.tf`
- `kubectl` access to the target EKS clusters
- S3 bucket for remote state

## Usage

### 1. Apply the account layer (once per AWS account)

```bash
cd account
terraform init
terraform apply
```

### 2. Deploy Kompass per cluster

Copy the `kompass/` directory for each EKS cluster. Change the backend `key` and set `cluster_name`:

```bash
cd kompass-eks-prod
terraform init
terraform apply -var="cluster_name=eks-prod"
```

Repeat for each cluster.

## Files

### `account/`

| File | Description |
|------|-------------|
| `provider.tf` | AWS and Zesty provider configuration |
| `main.tf` | S3 backend + module call with `create_values_local_file = false` |
| `outputs.tf` | Exposes `kompass_values_yaml` for downstream consumers |
| — | — |
| `versions.tf` | Required providers (AWS, Zesty) |

### `kompass/`

| File | Description |
|------|-------------|
| `provider.tf` | AWS + Helm provider (authenticates to EKS via data sources) |
| `main.tf` | S3 backend + `terraform_remote_state` + `helm_release` |
| `variables.tf` | `cluster_name`, `region` |
| `versions.tf` | Required providers (AWS, Helm) |
