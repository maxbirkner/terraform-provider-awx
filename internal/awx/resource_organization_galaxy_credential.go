package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagOrganizationGalaxyCredentialTitle = "Organization Galaxy Credential"

func resourceOrganizationsGalaxyCredentials() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource OrganizationsGalaxyCredentials manages the association of Galaxy credentials to an organization.",
		CreateContext: resourceOrganizationsGalaxyCredentialsCreate,
		DeleteContext: resourceOrganizationsGalaxyCredentialsDelete,
		ReadContext:   resourceOrganizationsGalaxyCredentialsRead,

		Schema: map[string]*schema.Schema{

			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The organization ID to associate the Galaxy credentials.",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The Galaxy credential ID to associate with the organization.",
			},
		},
	}
}

func resourceOrganizationsGalaxyCredentialsCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	orgID := d.Get("organization_id").(int)
	if _, err := client.OrganizationsService.GetOrganizationsByID(orgID, make(map[string]string)); err != nil {
		return utils.DiagNotFound(diagOrganizationGalaxyCredentialTitle, orgID, err)
	}

	result, err := client.OrganizationsService.AssociateGalaxyCredentials(orgID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{})

	if err != nil {
		return utils.DiagCreate(diagOrganizationGalaxyCredentialTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return nil
}

func resourceOrganizationsGalaxyCredentialsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOrganizationsGalaxyCredentialsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	orgID := d.Get("organization_id").(int)
	res, err := client.OrganizationsService.GetOrganizationsByID(orgID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagOrganizationGalaxyCredentialTitle, orgID, err)
	}

	if _, err = client.OrganizationsService.DisAssociateGalaxyCredentials(res.ID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{}); err != nil {
		return utils.DiagDelete(diagOrganizationGalaxyCredentialTitle, orgID, err)
	}

	d.SetId("")
	return nil
}
