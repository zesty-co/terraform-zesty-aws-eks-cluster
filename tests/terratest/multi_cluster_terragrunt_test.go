package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestMultiClusterTerragruntAccount runs the multi_clusters/terragrunt account layer.
//
// Required environment variables:
//   - ZESTY_API_TOKEN: Zesty API token
//
// Run:
//
//	go test -v -run TestMultiClusterTerragruntAccount -timeout 30m
func TestMultiClusterTerragruntAccount(t *testing.T) {
	t.Parallel()

	exampleDir := "../../examples/multi_clusters/terragrunt/live/prod/aws/us-east-1/my-account/zesty/account"
	tmpDir, _ := files.CopyTerraformFolderToTemp(t, exampleDir, t.Name())

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

// TestMultiClusterTerragruntKompassProd runs kompass for the eks-prod cluster.
// Requires the account layer to already be applied.
//
// Required environment variables:
//   - CLUSTER_NAME: overridden to "eks-prod" by the terragrunt.hcl locals
//
// Run:
//
//	go test -v -run TestMultiClusterTerragruntKompassProd -timeout 30m
func TestMultiClusterTerragruntKompassProd(t *testing.T) {
	t.Parallel()

	exampleDir := "../../examples/multi_clusters/terragrunt/live/prod/aws/us-east-1/my-account/zesty/kompass-eks-prod"
	tmpDir, _ := files.CopyTerraformFolderToTemp(t, exampleDir, t.Name())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    tmpDir,
		TerraformBinary: "terragrunt",
		NoColor:         true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}

// TestMultiClusterTerragruntRunAll runs `terragrunt run-all apply` from the zesty/ root.
// This is the closest to a real multi-cluster customer workflow:
// account applies first, then kompass-eks-prod and kompass-eks-staging in parallel.
//
// Required environment variables:
//   - ZESTY_API_TOKEN: Zesty API token
//   - EKS clusters "eks-prod" and "eks-staging" must exist
//
// Run:
//
//	go test -v -run TestMultiClusterTerragruntRunAll -timeout 45m
func TestMultiClusterTerragruntRunAll(t *testing.T) {
	_ = os.Getenv("ZESTY_API_TOKEN")

	exampleDir := "../../examples/multi_clusters/terragrunt/live/prod/aws/us-east-1/my-account/zesty"
	tmpDir := files.CopyTerraformFolderToTemp(t, exampleDir, t.Name())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    tmpDir,
		TerraformBinary: "terragrunt",
		NoColor:         true,
	})

	defer terraform.RunTerraformCommand(t, terraformOptions, "run-all", "destroy", "--terragrunt-non-interactive")

	terraform.RunTerraformCommand(t, terraformOptions, "run-all", "apply", "--terragrunt-non-interactive")
}
