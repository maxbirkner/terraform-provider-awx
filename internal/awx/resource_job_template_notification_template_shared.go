package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagJobTemplateNotificationTitle = "Job Template - Notification Template"

func getResourceJobTemplateNotificationTemplateAssociateFuncForType(client *awx.JobTemplateNotificationTemplatesService, typ string) func(jobTemplateID int, notificationTemplateID int) (*awx.NotificationTemplate, error) {
	switch typ {
	case "error":
		return client.AssociateJobTemplateNotificationTemplatesError
	case "success":
		return client.AssociateJobTemplateNotificationTemplatesSuccess
	case "started":
		return client.AssociateJobTemplateNotificationTemplatesStarted
	}
	return nil
}

func getResourceJobTemplateNotificationTemplateDisassociateFuncForType(client *awx.JobTemplateNotificationTemplatesService, typ string) func(jobTemplateID int, notificationTemplateID int) (*awx.NotificationTemplate, error) {
	switch typ {
	case "error":
		return client.DisassociateJobTemplateNotificationTemplatesError
	case "success":
		return client.DisassociateJobTemplateNotificationTemplatesSuccess
	case "started":
		return client.DisassociateJobTemplateNotificationTemplatesStarted
	}
	return nil
}

func resourceJobTemplateNotificationTemplateCreateForType(typ string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		jobTemplateID := d.Get("job_template_id").(int)
		if _, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string)); err != nil {
			return utils.DiagNotFound(diagJobTemplateNotificationTitle, jobTemplateID, err)
		}

		notificationTemplateID := d.Get("notification_template_id").(int)
		associationFunc := getResourceJobTemplateNotificationTemplateAssociateFuncForType(client.JobTemplateNotificationTemplatesService, typ)
		if associationFunc == nil {
			return utils.Diagf(
				"Create: JobTemplate not AssociateJobTemplateNotificationTemplates",
				"Fail to find association function for notification_template type %s", typ,
			)
		}

		result, err := associationFunc(jobTemplateID, notificationTemplateID)
		if err != nil {
			return utils.Diagf(
				"Create: JobTemplate not AssociateJobTemplateNotificationTemplates",
				"Fail to associate notification_template credentials with ID %v, for job_template ID %v, got error: %s",
				notificationTemplateID, jobTemplateID, err,
			)
		}

		d.SetId(strconv.Itoa(result.ID))
		return nil
	}
}

func resourceJobTemplateNotificationTemplateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceJobTemplateNotificationTemplateDeleteForType(typ string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		jobTemplateID := d.Get("job_template_id").(int)
		if _, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string)); err != nil {
			return utils.DiagNotFound(diagJobTemplateNotificationTitle, jobTemplateID, err)
		}

		notificationTemplateID := d.Get("notification_template_id").(int)
		disassociationFunc := getResourceJobTemplateNotificationTemplateDisassociateFuncForType(client.JobTemplateNotificationTemplatesService, typ)
		if disassociationFunc == nil {
			return utils.Diagf(
				"Create: JobTemplate not DisassociateJobTemplateNotificationTemplates",
				"Fail to find disassociation function for notification_template type %s", typ,
			)
		}

		if _, err := disassociationFunc(jobTemplateID, notificationTemplateID); err != nil {
			return utils.Diagf(
				"Create: JobTemplate not DisassociateJobTemplateNotificationTemplates",
				"Fail to associate notification_template credentials with ID %v, for job_template ID %v, got error: %s",
				notificationTemplateID, jobTemplateID, err,
			)
		}

		d.SetId("")
		return nil
	}
}
