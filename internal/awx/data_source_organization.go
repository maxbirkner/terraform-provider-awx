package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagOrganizationTitle = "Organization"

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRead,
		Description: "Data source for an AWX organization",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the organization",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the organization",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	organizations, err := client.OrganizationsService.ListOrganizations(params)
	if err != nil {
		return utils.DiagFetch(diagOrganizationTitle, params, err)
	}
	if len(organizations) > 1 {
		return utils.Diagf(
			"Get: find more than one Element",
			"The Query Returns more than one organization, %d",
			len(organizations),
		)
	}
	if len(organizations) == 0 {
		return utils.Diagf(
			"Get: Organization does not exist",
			"The Query Returns no Organization matching filter %v",
			params,
		)
	}

	organization := organizations[0]
	d = setOrganizationsResourceData(d, organization)
	return diags
}
