package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagEETitle = "Execution Environment"

func dataSourceExecutionEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExecutionEnvironmentsRead,
		Description: "Use this data source to get the details of an execution environment.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the execution environment.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the execution environment.",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceExecutionEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	executionEnvironments, _, err := client.ExecutionEnvironmentsService.ListExecutionEnvironments(params)
	if err != nil {
		return utils.DiagFetch(diagEETitle, params, err)
	}
	if len(executionEnvironments) > 1 {
		return utils.Diagf(
			"Get: find more than one element",
			"The query returns more than one execution environment, %d",
			len(executionEnvironments),
		)
	}
	if len(executionEnvironments) == 0 {
		return utils.Diagf(
			"Get: Execution Environment does not exist",
			"The Query Returns no Execution Environment matching filter %v",
			params,
		)
	}

	ee := executionEnvironments[0]
	d = setExecutionEnvironmentsResourceData(d, ee)
	return diags
}
