package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Host",
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		DeleteContext: resourceHostDelete,
		UpdateContext: resourceHostUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the host",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the host",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The inventory id of the host",
			},
			"group_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "The group ids of the host",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     "",
				Description: "The enabled status of the host",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The instance id of the host",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   normalizeJsonYaml,
				Description: "The variables of the host",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.HostService

	result, err := awxService.CreateHost(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(int),
		"enabled":     d.Get("enabled").(bool),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementHostTitle, err)
	}

	hostID := result.ID
	if d.IsNewResource() {
		rawGroups := d.Get("group_ids").([]interface{})
		for _, v := range rawGroups {

			_, err := awxService.AssociateGroup(hostID, map[string]interface{}{
				"id": v.(int),
			}, map[string]string{})
			if err != nil {
				return buildDiagnosticsMessage(
					diagElementHostTitle,
					"Assign Group Id %v to hostid %v fail, got  %s",
					v, hostID, err.Error(),
				)
			}
		}
	}
	d.SetId(strconv.Itoa(result.ID))
	return resourceHostRead(ctx, d, m)
}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	_, err := awxService.UpdateHost(id, map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(int),
		"enabled":     d.Get("enabled").(bool),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementHostTitle, id, err)
	}

	if d.HasChange("group_ids") {
		// TODO Check whats happen with removin groups ....
		rawGroups := d.Get("group_ids").([]interface{})
		for _, v := range rawGroups {
			_, err := awxService.AssociateGroup(id, map[string]interface{}{
				"id": v.(int),
			}, map[string]string{})
			if err != nil {
				return buildDiagnosticsMessage(
					diagElementHostTitle,
					"Assign Group Id %v to hostid %v fail, got  %s",
					v, id, err.Error(),
				)
			}
		}
	}
	return resourceHostRead(ctx, d, m)

}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}
	res, err := awxService.GetHostByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementHostTitle, id, err)
	}
	d = setHostResourceData(d, res)
	return nil
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteHost(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setHostResourceData(d *schema.ResourceData, r *awx.Host) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		return d
	}
	if err := d.Set("description", r.Description); err != nil {
		return d
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		return d
	}
	if err := d.Set("enabled", r.Enabled); err != nil {
		return d
	}
	if err := d.Set("instance_id", r.InstanceID); err != nil {
		return d
	}
	if err := d.Set("variables", normalizeJsonYaml(r.Variables)); err != nil {
		return d
	}
	if err := d.Set("group_ids", d.Get("group_ids").([]interface{})); err != nil {
		return d
	}
	return d
}
