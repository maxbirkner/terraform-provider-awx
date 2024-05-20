package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func getResourceWorkflowJobTemplateNotificationTemplateAssociateFuncForType(client *awx.WorkflowJobTemplateNotificationTemplatesService, typ string) func(workflowJobTemplateID int, notificationTemplateID int) (*awx.NotificationTemplate, error) {
	switch typ {
	case "error":
		return client.AssociateWorkflowJobTemplateNotificationTemplatesError
	case "success":
		return client.AssociateWorkflowJobTemplateNotificationTemplatesSuccess
	case "started":
		return client.AssociateWorkflowJobTemplateNotificationTemplatesStarted
	}
	return nil
}

func getResourceWorkflowJobTemplateNotificationTemplateDisassociateFuncForType(client *awx.WorkflowJobTemplateNotificationTemplatesService, typ string) func(workflowJobTemplateID int, notificationTemplateID int) (*awx.NotificationTemplate, error) {
	switch typ {
	case "error":
		return client.DisassociateWorkflowJobTemplateNotificationTemplatesError
	case "success":
		return client.DisassociateWorkflowJobTemplateNotificationTemplatesSuccess
	case "started":
		return client.DisassociateWorkflowJobTemplateNotificationTemplatesStarted
	}
	return nil
}

func resourceWorkflowJobTemplateNotificationTemplateCreateForType(typ string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		wjtID := d.Get("workflow_job_template_id").(int)
		if _, err := client.WorkflowJobTemplateService.GetWorkflowJobTemplateByID(wjtID, make(map[string]string)); err != nil {
			return utils.DiagNotFound("Workflow Job Template", wjtID, err)
		}

		ntID := d.Get("notification_template_id").(int)
		associationFunc := getResourceWorkflowJobTemplateNotificationTemplateAssociateFuncForType(client.WorkflowJobTemplateNotificationTemplatesService, typ)
		if associationFunc == nil {
			return utils.Diagf(
				"Create: WorkflowJobTemplate not AssociateWorkflowJobTemplateNotificationTemplates",
				"Fail to find association function for notification_template type %s", typ,
			)
		}

		result, err := associationFunc(wjtID, ntID)
		if err != nil {
			return utils.Diagf(
				"Create: WorkflowJobTemplate not AssociateWorkflowJobTemplateNotificationTemplates",
				"Fail to associate notification_template credentials with ID %v, for workflow_job_template ID %v, got error: %s",
				ntID, wjtID, err,
			)
		}

		d.SetId(strconv.Itoa(result.ID))
		return nil
	}
}

func resourceWorkflowJobTemplateNotificationTemplateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowJobTemplateNotificationTemplateDeleteForType(typ string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		wjtID := d.Get("workflow_job_template_id").(int)
		if _, err := client.WorkflowJobTemplateService.GetWorkflowJobTemplateByID(wjtID, make(map[string]string)); err != nil {
			return utils.DiagNotFound("workflow job template", wjtID, err)
		}

		ntID := d.Get("notification_template_id").(int)
		disassociationFunc := getResourceWorkflowJobTemplateNotificationTemplateDisassociateFuncForType(client.WorkflowJobTemplateNotificationTemplatesService, typ)
		if disassociationFunc == nil {
			return utils.Diagf(
				"Create: WorkflowJobTemplate not DisassociateWorkflowJobTemplateNotificationTemplates",
				"Fail to find disassociation function for notification_template type %s", typ,
			)
		}

		if _, err := disassociationFunc(wjtID, ntID); err != nil {
			return utils.Diagf(
				"Create: WorkflowJobTemplate not DisassociateWorkflowJobTemplateNotificationTemplates",
				"Fail to associate notification_template credentials with ID %v, for job_template ID %v, got error: %s",
				ntID, wjtID, err,
			)
		}

		d.SetId("")
		return nil
	}
}
