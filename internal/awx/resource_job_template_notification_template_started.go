package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJobTemplateNotificationTemplateStarted() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource for creating a notification template for a job template that will be sent when the job template is started.",
		CreateContext: resourceJobTemplateNotificationTemplateCreateForType("started"),
		DeleteContext: resourceJobTemplateNotificationTemplateDeleteForType("started"),
		ReadContext:   resourceJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The job template to associate the notification template with.",
			},
			"notification_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The notification template to associate with the job template.",
			},
		},
	}
}
