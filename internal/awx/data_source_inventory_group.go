package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceInventoryGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventoryGroupRead,
		Description: "Use this data source to get the details of an existing Inventory Group.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Inventory Group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Inventory Group.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the Inventory.",
			},
		},
	}
}

func dataSourceInventoryGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or id)",
		)
	}
	inventoryID := d.Get("inventory_id").(int)
	groups, _, err := client.InventoryGroupService.ListInventoryGroups(inventoryID, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Inventory Group",
			"Fail to find the group got: %s",
			err.Error(),
		)
	}
	if len(groups) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one Inventory Group, %d",
			len(groups),
		)
	}
	if len(groups) == 0 {
		return buildDiagnosticsMessage(
			"Get: Inventory Group does not exist",
			"The Query Returns no Inventory Group matching filter %v",
			params,
		)
	}

	group := groups[0]
	d = setInventoryGroupResourceData(d, group)
	return diags
}
