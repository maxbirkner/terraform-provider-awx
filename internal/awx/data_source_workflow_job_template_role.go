package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagWorkflowJobTemplateRole = "Workflow Job Template Role"

func dataSourceWorkflowJobTemplateRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowJobTemplateRoleRead,
		Description: "Data source for AWX Workflow Job Template Role",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the role",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the role",
				ExactlyOneOf: []string{"id", "name"},
			},
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the workflow job template",
			},
		},
	}
}

func dataSourceWorkflowJobTemplateRoleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	templateID := d.Get("workflow_job_template_id").(int)
	workflowJobTemplate, err := client.WorkflowJobTemplateService.GetWorkflowJobTemplateByID(templateID, params)
	if err != nil {
		return utils.DiagFetch(diagWorkflowJobTemplateRole, templateID, err)
	}

	rolesList := []*awx.ApplyRole{
		workflowJobTemplate.SummaryFields.ObjectRoles.UseRole,
		workflowJobTemplate.SummaryFields.ObjectRoles.AdminRole,
		workflowJobTemplate.SummaryFields.ObjectRoles.AdhocRole,
		workflowJobTemplate.SummaryFields.ObjectRoles.UpdateRole,
		workflowJobTemplate.SummaryFields.ObjectRoles.ReadRole,
		workflowJobTemplate.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setWorkflowJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setWorkflowJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	return utils.DiagNotFound(diagWorkflowJobTemplateRole, templateID, nil)
}

func setWorkflowJobTemplateRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
