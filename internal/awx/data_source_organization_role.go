package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceOrganizationRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRolesRead,
		Description: "Data source for an organization role",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the organization role",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the organization role",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the organization",
			},
		},
	}
}

func dataSourceOrganizationRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	org_id := d.Get("organization_id").(int)

	organization, err := client.OrganizationsService.GetOrganizationsByID(org_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch organization role",
			"Fail to find the organization role got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		organization.SummaryFields.ObjectRoles.AdhocRole,
		organization.SummaryFields.ObjectRoles.AdminRole,
		organization.SummaryFields.ObjectRoles.ApprovalRole,
		organization.SummaryFields.ObjectRoles.AuditorRole,
		organization.SummaryFields.ObjectRoles.CredentialAdminRole,
		organization.SummaryFields.ObjectRoles.ExecuteRole,
		organization.SummaryFields.ObjectRoles.InventoryAdminRole,
		organization.SummaryFields.ObjectRoles.JobTemplateAdminRole,
		organization.SummaryFields.ObjectRoles.MemberRole,
		organization.SummaryFields.ObjectRoles.NotificationAdminRole,
		organization.SummaryFields.ObjectRoles.ProjectAdminRole,
		organization.SummaryFields.ObjectRoles.ReadRole,
		organization.SummaryFields.ObjectRoles.UpdateRole,
		organization.SummaryFields.ObjectRoles.UseRole,
		organization.SummaryFields.ObjectRoles.WorkflowAdminRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch organization role - Not Found",
		"The organization role was not found",
	)
}

func setOrganizationRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		return d
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
