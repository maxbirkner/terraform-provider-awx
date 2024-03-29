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
				StateFunc:   utils.Normalize,
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
		return utils.DiagCreate(diagInventoryGroupTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryGroupRead(ctx, d, m)

}

func resourceInventoryGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.GroupService.UpdateGroup(id, map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(string),
		"variables":   d.Get("variables").(string),
	}, nil); err != nil {
		return utils.DiagUpdate(diagInventoryGroupTitle, id, err)
	}

	return resourceInventoryGroupRead(ctx, d, m)

}

func resourceInventoryGroupDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.GroupService.DeleteGroup(id); err != nil {
		return utils.DiagDelete(diagInventoryGroupTitle, id, err)
	}
	d.SetId("")
	return nil
}

func resourceInventoryGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := client.GroupService.GetGroupByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagFetch(diagInventoryGroupTitle, id, err)
	}
	d = setInventoryGroupResourceData(d, res)
	return nil
}

func setInventoryGroupResourceData(d *schema.ResourceData, r *awx.Group) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		fmt.Println("Error setting inventory_id", err)
	}
	if err := d.Set("variables", utils.Normalize(r.Variables)); err != nil {
		fmt.Println("Error setting variables", err)
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
