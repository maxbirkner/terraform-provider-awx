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

const diagOrganizationRole = "Organization Role"

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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the organization role",
				ExactlyOneOf: []string{"id", "name"},
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the organization",
			},
		},
	}
}

func dataSourceOrganizationRolesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	orgID := d.Get("organization_id").(int)

	organization, err := client.OrganizationsService.GetOrganizationsByID(orgID, params)
	if err != nil {
		return utils.DiagFetch(diagOrganizationRole, orgID, err)
	}

	rolesList := []*awx.ApplyRole{
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
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	return utils.DiagNotFound(diagOrganizationRole, orgID, nil)
}

func setOrganizationRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
