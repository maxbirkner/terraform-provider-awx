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

const diagNotificationTemplateTitle = "Notification Template"

func resourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_notification_template` manages notification templates within an AWX organization.",
		CreateContext: resourceNotificationTemplateCreate,
		ReadContext:   resourceNotificationTemplateRead,
		UpdateContext: resourceNotificationTemplateUpdate,
		DeleteContext: resourceNotificationTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the notification template.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization ID to associate with the notification template.",
			},
			"notification_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of notification template.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the notification template.",
			},
			"notification_configuration": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     nil,
				Description: "Build custom message responses for the notification template.",
			},
			"messages": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The description of the notification template. Options are `started`, `success`, `error`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"started": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The message to send when the job starts.",
						},
						"success": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The message to send when the job starts.",
						},
						"error": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The message to send when the job starts.",
						},
					},
				},
			},
		},
	}
}

type notifyMessage struct {
	Started map[string]interface{} `json:"started,omitempty"`
	Success map[string]interface{} `json:"success,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

func resourceNotificationTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	payload := map[string]interface{}{
		"name":                       d.Get("name").(string),
		"description":                d.Get("description").(string),
		"organization":               d.Get("organization_id").(string),
		"notification_type":          d.Get("notification_type").(string),
		"notification_configuration": parseNotifyConfig(d.Get("notification_configuration").(map[string]interface{})),
	}
	messages := d.Get("messages").(*schema.Set).List()
	if len(messages) != 0 {
		payload["messages"] = messages[0].(map[string]interface{})
	}
	result, err := client.NotificationTemplatesService.Create(payload, map[string]string{})
	if err != nil {
		return utils.DiagCreate("NotificationTemplate", err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update NotificationTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.NotificationTemplatesService.GetByID(id, params); err != nil {
		return utils.DiagNotFound(diagNotificationTemplateTitle, id, err)
	}
	payload := map[string]interface{}{
		"name":                       d.Get("name").(string),
		"description":                d.Get("description").(string),
		"organization":               d.Get("organization_id").(string),
		"notification_type":          d.Get("notification_type").(string),
		"notification_configuration": parseNotifyConfig(d.Get("notification_configuration").(map[string]interface{})),
	}
	messages := d.Get("messages").(*schema.Set).List()
	if len(messages) != 0 {
		payload["messages"] = messages[0].(map[string]interface{})
	}
	if _, err := client.NotificationTemplatesService.Update(id, payload, map[string]string{}); err != nil {
		return utils.DiagUpdate(diagNotificationTemplateTitle, id, err)
	}

	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read notification_template", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.NotificationTemplatesService.GetByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagNotificationTemplateTitle, id, err)

	}
	d = setNotificationTemplateResourceData(d, res)
	return nil
}

func resourceNotificationTemplateDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagNotificationTemplateTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.NotificationTemplatesService.Delete(id); err != nil {
		return utils.DiagDelete(diagNotificationTemplateTitle, id, err)
	}
	d.SetId("")
	return nil
}

func setNotificationTemplateResourceData(d *schema.ResourceData, r *awx.NotificationTemplate) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("organization", r.Organization); err != nil {
		fmt.Println("Error setting organization", err)
	}
	if err := d.Set("notification_type", r.NotificationType); err != nil {
		fmt.Println("Error setting notification_type", err)
	}
	if err := d.Set("notification_configuration", r.NotificationConfiguration); err != nil {
		fmt.Println("Error setting notification_configuration", err)
	}

	if err := d.Set("messages", schema.NewSet(func(i interface{}) int { return len(i.(map[string]interface{})) }, []interface{}{r.Messages})); err != nil {
		fmt.Println("Error setting messages", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}

func parseNotifyConfig(n map[string]interface{}) map[string]interface{} {
	for key, value := range n {
		if value == "" {
			delete(n, key)
			continue
		}
		switch key {
		case "channels":
			n[key] = []interface{}{value}
		}
	}
	return n
}
