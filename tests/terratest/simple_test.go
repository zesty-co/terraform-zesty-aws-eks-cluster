package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleExample runs a full E2E test of the simple/terraform example.
//
// Required environment variables:
//   - CLUSTER_NAME:    EKS cluster name to deploy into
//   - ZESTY_API_TOKEN: Zesty API token (set in provider.tf or override)
//
// Run:
//
//	CLUSTER_NAME=my-cluster go test -v -run TestSimpleExample -timeout 30m
func TestSimpleExample(t *testing.T) {
	t.Parallel()

	clusterName := os.Getenv("CLUSTER_NAME")
	require.NotEmpty(t, clusterName, "CLUSTER_NAME environment variable must be set")

	// Copy to temp dir so parallel tests don't collide on .terraform/
	exampleDir := "../../examples/simple/terraform"
	tmpDir, err := files.CopyTerraformFolderToTemp(exampleDir, t.Name())
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tmpDir,
		Vars: map[string]interface{}{
			"cluster_name": clusterName,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify the module produced a non-empty kompass_values_yaml output
	kompassValues := terraform.Output(t, terraformOptions, "kompass_values_yaml")
	assert.NotEmpty(t, kompassValues, "kompass_values_yaml output should not be empty")
}

// TestSimpleExamplePlanOnly runs terraform plan without applying.
// Useful as a smoke test when you don't want to create real resources.
//
// Required environment variables:
//   - CLUSTER_NAME: EKS cluster name
func TestSimpleExamplePlanOnly(t *testing.T) {
	t.Parallel()

	clusterName := os.Getenv("CLUSTER_NAME")
	require.NotEmpty(t, clusterName, "CLUSTER_NAME environment variable must be set")

	exampleDir := "../../examples/simple/terraform"
	tmpDir, err := files.CopyTerraformFolderToTemp(exampleDir, t.Name())
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tmpDir,
		Vars: map[string]interface{}{
			"cluster_name": clusterName,
		},
		NoColor: true,
	})

	terraform.InitAndPlan(t, terraformOptions)
}
