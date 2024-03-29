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

const diagJobTemplateRole = "Job Template Role"

func dataSourceJobTemplateRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobTemplateRoleRead,
		Description: "Data source for AWX Job Template Role",
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
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the job template",
			},
		},
	}
}

func dataSourceJobTemplateRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	templateID := d.Get("job_template_id").(int)
	jobTemplate, err := client.JobTemplateService.GetJobTemplateByID(templateID, params)
	if err != nil {
		return utils.DiagFetch(diagJobTemplateRole, templateID, err)
	}

	rolesList := []*awx.ApplyRole{
		jobTemplate.SummaryFields.ObjectRoles.UseRole,
		jobTemplate.SummaryFields.ObjectRoles.AdminRole,
		jobTemplate.SummaryFields.ObjectRoles.AdhocRole,
		jobTemplate.SummaryFields.ObjectRoles.UpdateRole,
		jobTemplate.SummaryFields.ObjectRoles.ReadRole,
		jobTemplate.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	return utils.DiagNotFound(diagJobTemplateRole, templateID, nil)
}

func setJobTemplateRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
