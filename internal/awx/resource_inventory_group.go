package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceInventoryGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Inventory Group is used to manage the group in the AWX",
		CreateContext: resourceInventoryGroupCreate,
		ReadContext:   resourceInventoryGroupRead,
		UpdateContext: resourceInventoryGroupUpdate,
		DeleteContext: resourceInventoryGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the group",
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The inventory id of the group",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   normalizeJsonYaml,
				Description: `The variables of the group. This can be in JSON or YAML format. For example:  {"key": "value"}`,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventoryGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.GroupService

	result, err := awxService.CreateGroup(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInventoryGroupTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryGroupRead(ctx, d, m)

}

func resourceInventoryGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.GroupService
	id, diags := convertStateIDToNummeric(diagElementInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	_, err := awxService.UpdateGroup(id, map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(string),
		"variables":   d.Get("variables").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInventoryGroupTitle, id, err)
	}

	return resourceInventoryGroupRead(ctx, d, m)

}

func resourceInventoryGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.GroupService

	id, diags := convertStateIDToNummeric(diagElementInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteGroup(id); err != nil {
		return buildDiagDeleteFail(
			diagElementInventoryGroupTitle,
			fmt.Sprintf("ID: %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func resourceInventoryGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.GroupService

	id, diags := convertStateIDToNummeric(diagElementInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetGroupByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementInventoryGroupTitle, id, err)
	}
	d = setInventoryGroupResourceData(d, res)
	return diags
}

func setInventoryGroupResourceData(d *schema.ResourceData, r *awx.Group) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("variables", normalizeJsonYaml(r.Variables))

	d.SetId(strconv.Itoa(r.ID))
	return d
}
