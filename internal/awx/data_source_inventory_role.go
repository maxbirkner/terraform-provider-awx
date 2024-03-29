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

const diagInventoryRole = "Inventory Role"

func dataSourceInventoryRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventoryRoleRead,
		Description: "Data source for inventory role",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the inventory role",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the inventory role",
				ExactlyOneOf: []string{"id", "name"},
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the inventory",
			},
		},
	}
}

func dataSourceInventoryRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	invID := d.Get("inventory_id").(int)
	inventory, err := client.InventoriesService.GetInventoryByID(invID, params)
	if err != nil {
		return utils.DiagFetch(diagInventoryRole, invID, err)
	}

	rolesList := []*awx.ApplyRole{
		inventory.SummaryFields.ObjectRoles.UseRole,
		inventory.SummaryFields.ObjectRoles.AdminRole,
		inventory.SummaryFields.ObjectRoles.AdhocRole,
		inventory.SummaryFields.ObjectRoles.UpdateRole,
		inventory.SummaryFields.ObjectRoles.ReadRole,
		inventory.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setInventoryRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setInventoryRoleData(d, v)
				return diags
			}
		}
	}

	return utils.DiagNotFound(diagInventoryRole, invID, nil)
}

func setInventoryRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
