package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagInventoryTitle = "Inventory"

func dataSourceInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventoriesRead,
		Description: "Use this data source to get the details of an existing Inventory.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the inventory.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the inventory.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the organization to which the inventory belongs.",
			},
		},
	}
}

func dataSourceInventoriesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if organizationID, okIOrgID := d.GetOk("organization_id"); okIOrgID {
		params["organization"] = strconv.Itoa(organizationID.(int))
	}
	inventories, _, err := client.InventoriesService.ListInventories(params)
	if err != nil {
		return utils.DiagFetch(diagInventoryTitle, params, err)
	}
	if len(inventories) > 1 {
		return utils.Diagf(
			"Get: find more than one Element",
			"The Query Returns more than one Inventory, %d",
			len(inventories),
		)
	}

	if len(inventories) == 0 {
		return utils.Diagf(
			"Get: Inventory does not exist",
			"The Query Returns no Inventory matching filter %v",
			params,
		)
	}

	inventory := inventories[0]
	d = setInventoryResourceData(d, inventory)
	return diags
}
