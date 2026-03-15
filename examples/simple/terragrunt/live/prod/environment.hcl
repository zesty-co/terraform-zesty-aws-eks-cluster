# Derives the environment name from the directory name (e.g. "prod").
# Directory structure is the single source of truth — no hardcoded duplicates
# to keep in sync. Rename the folder to change the environment everywhere at once.
locals {
  environment = basename(get_terragrunt_dir()) # => "prod"
}
