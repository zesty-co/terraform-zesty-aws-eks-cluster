# Derives the AWS region from the directory name (e.g. "us-east-1").
# Directory structure is the single source of truth — no hardcoded duplicates
# to keep in sync. Rename the folder to change the region everywhere at once.
locals {
  region = basename(get_terragrunt_dir()) # => "us-east-1"
}
