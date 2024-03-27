package awx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		Description: "This resource configure generic AWX settings.\n" +
			"Please note that resource deletion only delete object from terraform state and do not reset setting to his initial value.\n\n" +
			"See available settings list here: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Settings/Settings_settings_update",
		CreateContext: resourceSettingUpdate,
		ReadContext:   resourceSettingRead,
		DeleteContext: resourceSettingDelete,
		UpdateContext: resourceSettingUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of setting to modify",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value to be modified for given setting.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

type setting map[string]string

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.SettingService

	_, err := awxService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch setting, got: %s", err.Error(),
		)
	}

	var (
		mapDecoded     map[string]interface{}
		arrayDecoded   []interface{}
		formattedValue interface{} // nolint:typecheck
	)

	name := d.Get("name").(string)
	value := d.Get("value").(string)

	// Attempt to unmarshall string into a map
	err = json.Unmarshal([]byte(value), &mapDecoded)

	if err != nil {
		// Attempt to unmarshall string into an array
		err = json.Unmarshal([]byte(value), &arrayDecoded)

		if err != nil {
			formattedValue = value
		} else {
			formattedValue = arrayDecoded
		}
	} else {
		formattedValue = mapDecoded
	}

	payload := map[string]interface{}{
		name: formattedValue,
	}

	_, err = awxService.UpdateSettings("all", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: setting not created",
			"failed to save setting data, got: %s, %s", err.Error(), value,
		)
	}

	d.SetId(name)
	return resourceSettingRead(ctx, d, m)
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.SettingService

	_, err := awxService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to fetch settings",
			"Unable to load settings with slug all: got %s", err.Error(),
		)
	}

	d.Set("name", d.Id())
	d.Set("value", d.Get("value").(string))
	return diags
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
