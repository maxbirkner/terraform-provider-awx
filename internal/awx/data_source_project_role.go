package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceProjectRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRolesRead,
		Description: "Use this data source to get the details of a project role in AWX.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the project role.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the project role.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the project.",
			},
		},
	}
}

func dataSourceProjectRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	proj_id := d.Get("project_id").(int)

	Project, err := client.ProjectService.GetProjectByID(proj_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Project",
			"Fail to find the project got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		Project.SummaryFields.ObjectRoles.UseRole,
		Project.SummaryFields.ObjectRoles.AdminRole,
		Project.SummaryFields.ObjectRoles.UpdateRole,
		Project.SummaryFields.ObjectRoles.ReadRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setProjectRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setProjectRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch project role - Not Found",
		"The project role was not found",
	)
}

func setProjectRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
