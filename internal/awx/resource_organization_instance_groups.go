package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagOrganizationInstanceGroup = "Organization Instance Group"

func resourceOrganizationsInstanceGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationsInstanceGroupsCreate,
		DeleteContext: resourceOrganizationsInstanceGroupsDelete,
		ReadContext:   resourceOrganizationsInstanceGroupsRead,

		Schema: map[string]*schema.Schema{

			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_groups_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOrganizationsInstanceGroupsCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	OrganizationID := d.Get("organization_id").(int)
	_, err := awxService.GetOrganizationsByID(OrganizationID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("organization", OrganizationID, err)
	}

	result, err := awxService.AssociateInstanceGroups(OrganizationID, map[string]interface{}{
		"id": d.Get("instance_groups_id").(int),
	}, map[string]string{})

	if err != nil {
		return utils.DiagCreate(diagOrganizationInstanceGroup, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return diags
}

func resourceOrganizationsInstanceGroupsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceOrganizationsInstanceGroupsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	OrganizationID := d.Get("organization_id").(int)
	res, err := awxService.GetOrganizationsByID(OrganizationID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("organization", OrganizationID, err)
	}

	_, err = awxService.DisAssociateInstanceGroups(res.ID, map[string]interface{}{
		"id": d.Get("instance_groups_id").(int),
	}, map[string]string{})
	if err != nil {
		return utils.DiagDelete("Organization DisAssociateInstanceGroups", res.ID, err)
	}

	d.SetId("")
	return diags
}
