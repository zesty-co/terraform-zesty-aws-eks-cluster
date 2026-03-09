package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleTerragruntAccount runs the simple/terragrunt/account example.
//
// Required environment variables:
//   - ZESTY_API_TOKEN: Zesty API token (set in the generate block)
//
// Run:
//
//	go test -v -run TestSimpleTerragruntAccount -timeout 30m
func TestSimpleTerragruntAccount(t *testing.T) {
	t.Parallel()

	exampleDir := "../../examples/simple/terragrunt/account"
	tmpDir, err := files.CopyTerraformFolderToTemp(exampleDir, t.Name())
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    tmpDir,
		TerraformBinary: "terragrunt",
		NoColor:         true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	kompassValues := terraform.Output(t, terraformOptions, "kompass_values_yaml")
	assert.NotEmpty(t, kompassValues, "kompass_values_yaml output should not be empty")
}

// TestSimpleTerragruntKompass runs the simple/terragrunt/kompass example.
// Requires the account layer to already be applied.
//
// Required environment variables:
//   - CLUSTER_NAME: EKS cluster name
//
// Run:
//
//	CLUSTER_NAME=my-cluster go test -v -run TestSimpleTerragruntKompass -timeout 30m
func TestSimpleTerragruntKompass(t *testing.T) {
	t.Parallel()

	clusterName := os.Getenv("CLUSTER_NAME")
	require.NotEmpty(t, clusterName, "CLUSTER_NAME environment variable must be set")

	exampleDir := "../../examples/simple/terragrunt/kompass"
	tmpDir, err := files.CopyTerraformFolderToTemp(exampleDir, t.Name())
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    tmpDir,
		TerraformBinary: "terragrunt",
		NoColor:         true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}

// TestSimpleTerragruntFullE2E runs account + kompass sequentially via terragrunt.
//
// Required environment variables:
//   - CLUSTER_NAME:    EKS cluster name
//   - ZESTY_API_TOKEN: Zesty API token
//
// Run:
//
//	CLUSTER_NAME=my-cluster go test -v -run TestSimpleTerragruntFullE2E -timeout 45m
func TestSimpleTerragruntFullE2E(t *testing.T) {
	clusterName := os.Getenv("CLUSTER_NAME")
	require.NotEmpty(t, clusterName, "CLUSTER_NAME environment variable must be set")

	// Step 1: Apply account
	accountDir := "../../examples/simple/terragrunt/account"
	accountTmpDir, err := files.CopyTerraformFolderToTemp(accountDir, t.Name()+"-account")
	require.NoError(t, err)

	accountOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    accountTmpDir,
		TerraformBinary: "terragrunt",
		NoColor:         true,
	})

	defer terraform.Destroy(t, accountOpts)
	terraform.InitAndApply(t, accountOpts)

	kompassValues := terraform.Output(t, accountOpts, "kompass_values_yaml")
	assert.NotEmpty(t, kompassValues, "account layer should output kompass_values_yaml")

	// Step 2: Apply kompass
	kompassDir := "../../examples/simple/terragrunt/kompass"
	kompassTmpDir, err := files.CopyTerraformFolderToTemp(kompassDir, t.Name()+"-kompass")
	require.NoError(t, err)

	kompassOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    kompassTmpDir,
		TerraformBinary: "terragrunt",
		NoColor:         true,
	})

	defer terraform.Destroy(t, kompassOpts)
	terraform.InitAndApply(t, kompassOpts)
}
