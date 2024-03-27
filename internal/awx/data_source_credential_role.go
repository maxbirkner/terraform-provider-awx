package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceCredentialMachineRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialMachineRoleRead,
		Description: "Use this data source to get the role of a machine credential",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the credential role",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the credential role",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the credential to fetch the role from",
			},
		},
	}
}

func dataSourceCredentialMachineRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	credID := d.Get("credential_id").(int)
	credentialID, err := client.CredentialsService.GetCredentialsByID(credID, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Credential",
			"Fail to find the credential, got: %s",
			err.Error(),
		)
	}

	rolesList := []*awx.ApplyRole{
		credentialID.SummaryFields.ObjectRoles.UseRole,
		credentialID.SummaryFields.ObjectRoles.AdminRole,
		credentialID.SummaryFields.ObjectRoles.AdhocRole,
		credentialID.SummaryFields.ObjectRoles.UpdateRole,
		credentialID.SummaryFields.ObjectRoles.ReadRole,
		credentialID.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch machine credential role - Not Found",
		"The credential role was not found",
	)
}

func setCredentialRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
