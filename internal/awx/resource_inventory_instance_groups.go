package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceInventoryInstanceGroups() *schema.Resource {
	return &schema.Resource{
		Description:   "Associates an instance group to a inventory",
		CreateContext: resourceInventoryInstanceGroupsCreate,
		DeleteContext: resourceInventoryInstanceGroupsDelete,
		ReadContext:   resourceInventoryInstanceGroupsRead,

		Schema: map[string]*schema.Schema{

			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the inventory to associate the instance group with",
			},
			"instance_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the instance group to associate with the inventory",
			},
		},
	}
}

func resourceInventoryInstanceGroupsCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	inventoryId := d.Get("inventory_id").(int)
	if _, err := client.InventoriesService.GetInventoryByID(inventoryId, make(map[string]string)); err != nil {
		return utils.DiagNotFound("Inventory InstanceGroup", inventoryId, err)
	}

	result, err := client.InventoriesService.AssociateInstanceGroups(inventoryId, map[string]interface{}{
		"id": d.Get("instance_group_id").(int),
	}, map[string]string{})

	if err != nil {
		return utils.DiagCreate("Inventory AssociateInstanceGroups", err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return nil
}

func resourceInventoryInstanceGroupsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInventoryInstanceGroupsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	inventoryId := d.Get("inventory_id").(int)
	res, err := client.InventoriesService.GetInventoryByID(inventoryId, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("Inventory InstanceGroup", inventoryId, err)
	}

	if _, err = client.InventoriesService.DisAssociateInstanceGroups(res.ID, map[string]interface{}{
		"id": d.Get("instance_group_id").(int),
	}, map[string]string{}); err != nil {
		return utils.DiagDelete("Inventory DisAssociateInstanceGroups", inventoryId, err)
	}

	d.SetId("")
	return nil
}
