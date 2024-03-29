package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagWorkflowJobTemplateTitle = "Workflow Job Template"

func dataSourceWorkflowJobTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowJobTemplateRead,
		Description: "Use this data source to get the details of a workflow job template.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the workflow job template.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the workflow job template.",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceWorkflowJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	workflowJobTemplate, _, err := client.WorkflowJobTemplateService.ListWorkflowJobTemplates(params)
	if err != nil {
		return utils.DiagFetch(diagWorkflowJobTemplateTitle, params, err)
	}
	if groupName, okName := d.GetOk("name"); okName {
		for _, template := range workflowJobTemplate {
			if template.Name == groupName {
				d = setWorkflowJobTemplateResourceData(d, template)
				return diags
			}
		}
	}
	if _, okGroupID := d.GetOk("id"); okGroupID {
		if len(workflowJobTemplate) > 1 {
			return utils.Diagf(
				"Get: find more than one Element",
				"The Query Returns more than one Group, %d",
				len(workflowJobTemplate),
			)
		}
		if len(workflowJobTemplate) == 0 {
			return utils.Diagf(
				"Get: Workflow template does not exist",
				"The Query Returns no Workflow template matching filter %v",
				params,
			)
		}
		d = setWorkflowJobTemplateResourceData(d, workflowJobTemplate[0])
		return diags
	}
	return utils.Diagf(
		"Get: find more than one Element",
		"The Query Returns more than one Group, %d",
		len(workflowJobTemplate),
	)
}
