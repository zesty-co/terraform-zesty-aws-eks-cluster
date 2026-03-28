# Terraform Simple Example

Links your AWS account to Zesty and deploys the Kompass Helm chart into a **single EKS cluster**. Everything lives in one directory with one state file.

## Prerequisites

- Terraform >= 1.1
- AWS credentials configured
- Zesty API token — set the `token` value in `provider.tf`
- `kubectl` access to the target EKS cluster

## Usage

```bash
terraform init
terraform apply -var="cluster_name=my-eks-cluster"
```

## Files

| File | Description |
|------|-------------|
| `provider.tf` | AWS, Zesty, and Helm provider configuration |
| `main.tf` | Module call + `helm_release` |
| `variables.tf` | `cluster_name` |
| `versions.tf` | Required providers (AWS, Zesty, Helm) |
