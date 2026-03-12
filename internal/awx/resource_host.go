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

const diagHostTitle = "Host"

func resourceHost() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Host",
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		DeleteContext: resourceHostDelete,
		UpdateContext: resourceHostUpdate,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    resourceHostV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceHostStateUpgradeV0,
			},
		},

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
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
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
				StateFunc:   utils.Normalize,
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
		return utils.DiagCreate(diagHostTitle, err)
	}

	hostID := result.ID
	if d.IsNewResource() {
		rawGroups := d.Get("group_ids").(*schema.Set).List()
		for _, v := range rawGroups {

			if _, err := awxService.AssociateGroup(hostID, map[string]interface{}{
				"id": v.(int),
			}, map[string]string{}); err != nil {
				return utils.Diagf(diagHostTitle, "Assign Group Id %v to hostid %v fail, got  %s", v, hostID, err)
			}
		}
	}
	d.SetId(strconv.Itoa(result.ID))
	return resourceHostRead(ctx, d, m)
}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.HostService.UpdateHost(id, map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(int),
		"enabled":     d.Get("enabled").(bool),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}, nil); err != nil {
		return utils.DiagUpdate(diagHostTitle, id, err)
	}

	if d.HasChange("group_ids") {
		oldVal, newVal := d.GetChange("group_ids")
		oldSet := oldVal.(*schema.Set)
		newSet := newVal.(*schema.Set)

		for _, v := range newSet.Difference(oldSet).List() {
			gid := v.(int)
			if _, err := client.HostService.AssociateGroup(id, map[string]interface{}{
				"id": gid,
			}, map[string]string{}); err != nil {
				return utils.Diagf(diagHostTitle, "Associate Group Id %v to hostid %v fail, got  %s", gid, id, err)
			}
		}

		for _, v := range oldSet.Difference(newSet).List() {
			gid := v.(int)
			if _, err := client.HostService.DisAssociateGroup(id, map[string]interface{}{
				"id": gid,
			}, map[string]string{}); err != nil {
				return utils.Diagf(diagHostTitle, "Disassociate Group Id %v from hostid %v fail, got  %s", gid, id, err)
			}
		}
	}
	return resourceHostRead(ctx, d, m)
}

func resourceHostRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagHostTitle, d)
	if diags.HasError() {
		return diags
	}
	res, err := client.HostService.GetHostByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagHostTitle, id, err)
	}
	d = setHostResourceData(d, res)

	// Fetch actual group memberships from AWX API
	groups, err := client.HostService.ListHostGroups(id, make(map[string]string))
	if err != nil {
		return utils.Diagf(diagHostTitle, "Failed to list groups for host %v: %s", id, err)
	}
	groupIDs := make([]int, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}
	if err := d.Set("group_ids", groupIDs); err != nil {
		return utils.Diagf(diagHostTitle, "Error setting group_ids for host %v: %s", id, err)
	}
	return nil
}

func resourceHostDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.HostService.DeleteHost(id); err != nil {
		return utils.DiagDelete(diagHostTitle, id, err)
	}
	d.SetId("")
	return nil
}

func setHostResourceData(d *schema.ResourceData, r *awx.Host) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		fmt.Println("Error setting inventory_id", err)
	}
	if err := d.Set("enabled", r.Enabled); err != nil {
		fmt.Println("Error setting enabled", err)
	}
	if err := d.Set("instance_id", r.InstanceID); err != nil {
		fmt.Println("Error setting instance_id", err)
	}
	if err := d.Set("variables", utils.Normalize(r.Variables)); err != nil {
		fmt.Println("Error setting variables", err)
	}
	// group_ids are set by resourceHostRead after fetching from the API
	return d
}

// resourceHostV0 returns the v0 schema (before TypeList→TypeSet migration)
// used by the state upgrader to interpret old state.
func resourceHostV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":         {Type: schema.TypeString},
			"description":  {Type: schema.TypeString},
			"inventory_id": {Type: schema.TypeInt},
			"group_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{Type: schema.TypeInt},
			},
			"enabled":     {Type: schema.TypeBool},
			"instance_id": {Type: schema.TypeString},
			"variables":   {Type: schema.TypeString},
		},
	}
}

// resourceHostStateUpgradeV0 migrates state from v0 (TypeList) to v1 (TypeSet).
// The JSON representation is identical (array of ints), so the state is passed through unchanged.
func resourceHostStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	return rawState, nil
}
