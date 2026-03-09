package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiClusterAccountExample runs a full E2E test of the account layer.
//
// Required environment variables:
//   - ZESTY_API_TOKEN: Zesty API token
//
// Run:
//
//	go test -v -run TestMultiClusterAccountExample -timeout 30m
func TestMultiClusterAccountExample(t *testing.T) {
	t.Parallel()

	exampleDir := "../../examples/multi_clusters/terraform/account"
	tmpDir := files.CopyTerraformFolderToTemp(t, exampleDir, t.Name())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tmpDir,
		NoColor:      true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify the account module outputs kompass_values_yaml
	kompassValues := terraform.Output(t, terraformOptions, "kompass_values_yaml")
	assert.NotEmpty(t, kompassValues, "kompass_values_yaml output should not be empty")
}

// TestMultiClusterKompassExample runs a full E2E test of the kompass layer.
// Requires the account layer to already be applied (remote state must exist).
//
// Required environment variables:
//   - CLUSTER_NAME: EKS cluster name to deploy into
//
// Run:
//
//	CLUSTER_NAME=eks-prod go test -v -run TestMultiClusterKompassExample -timeout 30m
func TestMultiClusterKompassExample(t *testing.T) {
	t.Parallel()

	clusterName := os.Getenv("CLUSTER_NAME")
	require.NotEmpty(t, clusterName, "CLUSTER_NAME environment variable must be set")

	exampleDir := "../../examples/multi_clusters/terraform/kompass"
	tmpDir := files.CopyTerraformFolderToTemp(t, exampleDir, t.Name())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tmpDir,
		Vars: map[string]interface{}{
			"cluster_name": clusterName,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}

// TestMultiClusterFullE2E runs account + kompass in sequence.
// This is the closest to a real customer workflow.
//
// Required environment variables:
//   - CLUSTER_NAME:    EKS cluster name
//   - ZESTY_API_TOKEN: Zesty API token
//
// Run:
//
//	CLUSTER_NAME=eks-prod go test -v -run TestMultiClusterFullE2E -timeout 45m
func TestMultiClusterFullE2E(t *testing.T) {
	clusterName := os.Getenv("CLUSTER_NAME")
	require.NotEmpty(t, clusterName, "CLUSTER_NAME environment variable must be set")

	// Step 1: Apply account layer
	accountDir := "../../examples/multi_clusters/terraform/account"
	accountTmpDir := files.CopyTerraformFolderToTemp(t, accountDir, t.Name()+"-account")

	accountOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: accountTmpDir,
		NoColor:      true,
	})

	defer terraform.Destroy(t, accountOpts)
	terraform.InitAndApply(t, accountOpts)

	kompassValues := terraform.Output(t, accountOpts, "kompass_values_yaml")
	assert.NotEmpty(t, kompassValues, "account layer should output kompass_values_yaml")

	// Step 2: Apply kompass layer
	kompassDir := "../../examples/multi_clusters/terraform/kompass"
	kompassTmpDir := files.CopyTerraformFolderToTemp(t, kompassDir, t.Name()+"-kompass")

	kompassOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: kompassTmpDir,
		Vars: map[string]interface{}{
			"cluster_name": clusterName,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, kompassOpts)
	terraform.InitAndApply(t, kompassOpts)
}
