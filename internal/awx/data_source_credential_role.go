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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the credential role",
				ExactlyOneOf: []string{"id", "name"},
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the credential to fetch the role from",
			},
		},
	}
}

func dataSourceCredentialMachineRoleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	credID := d.Get("credential_id").(int)
	credentialID, err := client.CredentialsService.GetCredentialsByID(credID, params)
	if err != nil {
		return utils.DiagFetch("Machine Credential Role", credID, err)
	}

	rolesList := []*awx.ApplyRole{
		credentialID.SummaryFields.ObjectRoles.UseRole,
		credentialID.SummaryFields.ObjectRoles.AdminRole,
		credentialID.SummaryFields.ObjectRoles.AdhocRole,
		credentialID.SummaryFields.ObjectRoles.UpdateRole,
		credentialID.SummaryFields.ObjectRoles.ReadRole,
		credentialID.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, ok := d.GetOk("id"); ok {
		id := roleID.(int)
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, ok := d.GetOk("name"); ok {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	return utils.DiagNotFound("Machine Credential Role", credID, nil)
}

func setCredentialRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
