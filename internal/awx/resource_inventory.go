package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceInventory() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Inventory is used to define an inventory in AWX",
		CreateContext: resourceInventoryCreate,
		ReadContext:   resourceInventoryRead,
		DeleteContext: resourceInventoryDelete,
		UpdateContext: resourceInventoryUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the inventory",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the inventory",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id of the inventory",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The kind of the inventory",
			},
			"host_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The host filter of the inventory",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   normalizeJsonYaml,
				Description: "The variables of the inventory",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService

	result, err := awxService.CreateInventory(map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(string),
		"description":  d.Get("description").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInventoryTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryRead(ctx, d, m)

}

func resourceInventoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, diags := convertStateIDToNummeric(diagElementInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	_, err := awxService.UpdateInventory(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(string),
		"description":  d.Get("description").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInventoryTitle, id, err)
	}

	return resourceInventoryRead(ctx, d, m)

}

func resourceInventoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, err := strconv.Atoi(d.Id())
	id, diags := convertStateIDToNummeric(diagElementInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	r, err := awxService.GetInventory(id, map[string]string{})
	if err != nil {
		return buildDiagNotFoundFail(diagElementInventoryTitle, id, err)
	}
	d = setInventoryResourceData(d, r)
	return nil
}

func resourceInventoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, diags := convertStateIDToNummeric(diagElementInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := awxService.DeleteInventory(id); err != nil {
		return buildDiagDeleteFail(
			diagElementInventoryTitle,
			fmt.Sprintf(
				"%s %v, got %s ",
				diagElementInventoryTitle, id, err.Error(),
			),
		)
	}
	d.SetId("")
	return nil
}

func setInventoryResourceData(d *schema.ResourceData, r *awx.Inventory) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		return d
	}
	if err := d.Set("organization_id", strconv.Itoa(r.Organization)); err != nil {
		return d
	}
	if err := d.Set("description", r.Description); err != nil {
		return d
	}
	if err := d.Set("kind", r.Kind); err != nil {
		return d
	}
	if err := d.Set("host_filter", r.HostFilter); err != nil {
		return d
	}
	if err := d.Set("variables", normalizeJsonYaml(r.Variables)); err != nil {
		return d
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
