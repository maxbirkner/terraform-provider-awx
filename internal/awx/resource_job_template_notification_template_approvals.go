package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJobTemplateNotificationTemplateApprovals() *schema.Resource {
	return &schema.Resource{
		Description:   "A notification template for a job template that is triggered on request for approvals.",
		CreateContext: resourceJobTemplateNotificationTemplateCreateForType("approvals"),
		DeleteContext: resourceJobTemplateNotificationTemplateDeleteForType("approvals"),
		ReadContext:   resourceJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The job template ID that the notification template is associated with.",
			},
			"notification_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The notification template ID that the notification template is associated with.",
			},
		},
	}
}
