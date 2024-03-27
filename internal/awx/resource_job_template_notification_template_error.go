package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJobTemplateNotificationTemplateError() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource for creating a job template notification template error.",
		CreateContext: resourceJobTemplateNotificationTemplateCreateForType("error"),
		DeleteContext: resourceJobTemplateNotificationTemplateDeleteForType("error"),
		ReadContext:   resourceJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The job template ID to associate with the notification template.",
			},
			"notification_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The notification template ID to associate with the job template.",
			},
		},
	}
}
