//nolint:testpackage
package test

import (
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform" //nolint:depguard
	"github.com/stretchr/testify/assert"                  //nolint:depguard
)

func TestTerraformMinimalAccExample(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/k8s/base",
	}
	terraformDataOptions := &terraform.Options{
		TerraformDir: "../examples/k8s/data",
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
	//nolint:ineffassign
	output := terraform.Output(t, terraformOptions, "inventory_id")
	nr, err := strconv.Atoi(output)
	if err != nil {
		t.Logf("Inventory id is not a number")
		t.Fail()
	}
	assert.Greater(t, nr, 0)
	terraform.InitAndApply(t, terraformDataOptions)
	defer terraform.Destroy(t, terraformDataOptions)
	output = terraform.Output(t, terraformDataOptions, "job")
}
