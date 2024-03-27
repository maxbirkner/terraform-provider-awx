package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowJobTemplateNotificationTemplateSuccess() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource for creating a notification template for a workflow job template that will be triggered on success.",
		CreateContext: resourceWorkflowJobTemplateNotificationTemplateCreateForType("success"),
		DeleteContext: resourceWorkflowJobTemplateNotificationTemplateDeleteForType("success"),
		ReadContext:   resourceWorkflowJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the workflow job template to associate the notification template with.",
			},
			"notification_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the notification template to associate with the workflow job template.",
			},
		},
	}
}
