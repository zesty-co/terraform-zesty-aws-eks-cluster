# Terragrunt Simple Example

Links your AWS account to Zesty and deploys the Kompass Helm chart into a **single EKS cluster** using Terragrunt with a production-grade `live/` directory hierarchy.

## Directory Structure

```
simple/terragrunt/
├── modules/kompass/                            # reusable Helm module
│   ├── main.tf
│   ├── variables.tf
│   └── versions.tf
└── live/
    ├── root.hcl                                # org + project
    └── prod/
        ├── environment.hcl                     # => "prod" (from dirname)
        └── aws/
            ├── datacenter.hcl                  # S3 backend + AWS provider
            └── us-east-1/
                ├── region.hcl                  # => "us-east-1" (from dirname)
                └── my-account/
                    ├── account.hcl             # AWS account ID + profile
                    └── zesty/
                        ├── account/            # IAM + Zesty registration (apply first)
                        └── kompass/            # Helm release (depends on account)
```

## Prerequisites

- Terraform >= 1.1
- Terragrunt >= 0.45
- AWS credentials configured
- Zesty API token — set the `token` value in `account/terragrunt.hcl`
- `kubectl` access to the target EKS cluster

## Usage

### Apply both at once

```bash
cd live/prod/aws/us-east-1/my-account/zesty
terragrunt run-all apply
```

### Apply individually

```bash
cd live/prod/aws/us-east-1/my-account/zesty

# 1. Account first
cd account && terragrunt apply

# 2. Kompass
cd ../kompass && terragrunt apply
```

## Configuration

- **Cluster name**: change `cluster_name` in `kompass/terragrunt.hcl` locals
- **Region**: rename the `us-east-1/` directory — value is derived from the folder name
- **Environment**: rename the `prod/` directory — value is derived from the folder name
- **AWS profile**: change `profile` in `account.hcl`
- **S3 backend**: change bucket pattern in `datacenter.hcl`
