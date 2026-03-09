SHELL := /bin/bash

SIMPLE_TF     := ../examples/simple/terraform
MC_ACCOUNT_TF := ../examples/multi_clusters/terraform/account
MC_KOMPASS_TF := ../examples/multi_clusters/terraform/kompass

SIMPLE_TG_ACCOUNT := ../examples/simple/terragrunt/account
SIMPLE_TG_KOMPASS := ../examples/simple/terragrunt/kompass
MC_TG_ZESTY       := ../examples/multi_clusters/terragrunt/live/prod/aws/us-east-1/my-account/zesty

# ─── Layer 1: terraform validate (no credentials needed) ─────────────────────

.PHONY: validate validate-simple validate-mc-account validate-mc-kompass

validate: validate-simple validate-mc-account validate-mc-kompass
	@echo "All examples validated successfully."

validate-simple:
	@echo "── validate: simple/terraform ──"
	cd $(SIMPLE_TF) && terraform init -backend=false && terraform validate

validate-mc-account:
	@echo "── validate: multi_clusters/terraform/account ──"
	cd $(MC_ACCOUNT_TF) && terraform init -backend=false && terraform validate

validate-mc-kompass:
	@echo "── validate: multi_clusters/terraform/kompass ──"
	cd $(MC_KOMPASS_TF) && terraform init -backend=false && terraform validate

# ─── Layer 2: terraform plan (needs AWS_* + ZESTY_API_TOKEN) ─────────────────

.PHONY: plan plan-simple plan-mc-account plan-mc-kompass

plan: plan-simple plan-mc-account plan-mc-kompass
	@echo "All plans succeeded."

plan-simple:
	@echo "── plan: simple/terraform ──"
	cd $(SIMPLE_TF) && terraform init -backend=false && terraform plan -var="cluster_name=$(CLUSTER_NAME)"

plan-mc-account:
	@echo "── plan: multi_clusters/terraform/account ──"
	cd $(MC_ACCOUNT_TF) && terraform init -backend=false && terraform plan

plan-mc-kompass:
	@echo "── plan: multi_clusters/terraform/kompass ──"
	cd $(MC_KOMPASS_TF) && terraform init -backend=false && terraform plan -var="cluster_name=$(CLUSTER_NAME)"

# ─── Layer 3: terraform test (mock providers, TF >= 1.7) ─────────────────────

.PHONY: test-tf test-tf-simple test-tf-mc-account test-tf-mc-kompass

test-tf: test-tf-simple test-tf-mc-account test-tf-mc-kompass
	@echo "All terraform tests passed."

test-tf-simple:
	@echo "── terraform test: simple/terraform ──"
	cd $(SIMPLE_TF) && terraform test

test-tf-mc-account:
	@echo "── terraform test: multi_clusters/terraform/account ──"
	cd $(MC_ACCOUNT_TF) && terraform test

test-tf-mc-kompass:
	@echo "── terraform test: multi_clusters/terraform/kompass ──"
	cd $(MC_KOMPASS_TF) && terraform test

# ─── Layer 4: Terratest E2E — Terraform (needs real infra) ───────────────────

.PHONY: test-e2e-tf test-e2e-tf-simple test-e2e-tf-mc

test-e2e-tf: test-e2e-tf-simple test-e2e-tf-mc
	@echo "All Terraform E2E tests passed."

test-e2e-tf-simple:
	@echo "── Terratest E2E: simple/terraform ──"
	cd terratest && CLUSTER_NAME=$(CLUSTER_NAME) go test -v -run TestSimpleExample -timeout 30m

test-e2e-tf-mc:
	@echo "── Terratest E2E: multi_clusters/terraform ──"
	cd terratest && CLUSTER_NAME=$(CLUSTER_NAME) go test -v -run TestMultiClusterFullE2E -timeout 45m

# ─── Layer 5: Terratest E2E — Terragrunt (needs real infra) ──────────────────

.PHONY: test-e2e-tg test-e2e-tg-simple test-e2e-tg-mc

test-e2e-tg: test-e2e-tg-simple test-e2e-tg-mc
	@echo "All Terragrunt E2E tests passed."

test-e2e-tg-simple:
	@echo "── Terratest E2E: simple/terragrunt ──"
	cd terratest && CLUSTER_NAME=$(CLUSTER_NAME) go test -v -run TestSimpleTerragruntFullE2E -timeout 45m

test-e2e-tg-mc:
	@echo "── Terratest E2E: multi_clusters/terragrunt (run-all) ──"
	cd terratest && go test -v -run TestMultiClusterTerragruntRunAll -timeout 45m

# ─── Run everything ──────────────────────────────────────────────────────────

.PHONY: test-e2e-all

test-e2e-all: test-e2e-tf test-e2e-tg
	@echo "All E2E tests passed."

# ─── Cleanup ─────────────────────────────────────────────────────────────────

.PHONY: clean

clean:
	find ../examples -name ".terraform" -type d -exec rm -rf {} + 2>/dev/null || true
	find ../examples -name ".terragrunt-cache" -type d -exec rm -rf {} + 2>/dev/null || true
	find ../examples -name "*.tfstate*" -type f -delete 2>/dev/null || true
	find ../examples -name ".terraform.lock.hcl" -type f -delete 2>/dev/null || true
	@echo "Cleaned up."
