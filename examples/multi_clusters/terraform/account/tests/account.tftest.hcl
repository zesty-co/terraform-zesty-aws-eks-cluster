# ─── Mock providers (no real credentials needed) ─────────────────────────────

mock_provider "aws" {}

mock_provider "zesty" {}

# ─── Test: module creates resources and exposes output ────────────────────────

run "account_module_creates_resources" {
  command = apply

  assert {
    condition     = module.zesty.kompass_values_yaml != null
    error_message = "module should output kompass_values_yaml"
  }
}
