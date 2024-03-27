package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const jobTemplateNotificationTemplatesAPIEndpoint = "/api/v2/job_templates/%d/notification_templates_%s/"

// JobTemplateNotificationTemplatesService implements awx job template nodes apis.
type JobTemplateNotificationTemplatesService struct {
	client *Client
}

func (jt *JobTemplateNotificationTemplatesService) associateJobTemplateNotificationTemplatesForType(jobTemplateID int, notificationTemplateID int, typ string) (*NotificationTemplate, error) {
	result := new(NotificationTemplate)

	data := map[string]interface{}{
		"id": notificationTemplateID,
	}
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	endpoint := fmt.Sprintf(jobTemplateNotificationTemplatesAPIEndpoint, jobTemplateID, typ)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// AssociateJobTemplateNotificationTemplatesError will associate an error notification_template for a job_template.
func (jt *JobTemplateNotificationTemplatesService) AssociateJobTemplateNotificationTemplatesError(jobTemplateID int, notificationTemplateID int) (*NotificationTemplate, error) {
	return jt.associateJobTemplateNotificationTemplatesForType(jobTemplateID, notificationTemplateID, "error")
}

// AssociateJobTemplateNotificationTemplatesSuccess will associate a success notification_template for a job_template.
func (jt *JobTemplateNotificationTemplatesService) AssociateJobTemplateNotificationTemplatesSuccess(jobTemplateID int, notificationTemplateID int) (*NotificationTemplate, error) {
	return jt.associateJobTemplateNotificationTemplatesForType(jobTemplateID, notificationTemplateID, "success")
}

// AssociateJobTemplateNotificationTemplatesStarted will associate a started notification_template for a job_template.
func (jt *JobTemplateNotificationTemplatesService) AssociateJobTemplateNotificationTemplatesStarted(jobTemplateID int, notificationTemplateID int) (*NotificationTemplate, error) {
	return jt.associateJobTemplateNotificationTemplatesForType(jobTemplateID, notificationTemplateID, "started")
}

func (jt *JobTemplateNotificationTemplatesService) disassociateJobTemplateNotificationTemplatesForType(jobTemplateID int, notificationTemplateID int, typ string) (*NotificationTemplate, error) {
	result := new(NotificationTemplate)

	data := map[string]interface{}{
		"id":           notificationTemplateID,
		"disassociate": true,
	}
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	endpoint := fmt.Sprintf(jobTemplateNotificationTemplatesAPIEndpoint, jobTemplateID, typ)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// DisassociateJobTemplateNotificationTemplatesError will disassociate an error notification_template for a job_template.
func (jt *JobTemplateNotificationTemplatesService) DisassociateJobTemplateNotificationTemplatesError(jobTemplateID int, notificationTemplateID int) (*NotificationTemplate, error) {
	return jt.disassociateJobTemplateNotificationTemplatesForType(jobTemplateID, notificationTemplateID, "error")
}

// DisassociateJobTemplateNotificationTemplatesSuccess will disassociate a success notification_template for a job_template.
func (jt *JobTemplateNotificationTemplatesService) DisassociateJobTemplateNotificationTemplatesSuccess(jobTemplateID int, notificationTemplateID int) (*NotificationTemplate, error) {
	return jt.disassociateJobTemplateNotificationTemplatesForType(jobTemplateID, notificationTemplateID, "success")
}

// DisassociateJobTemplateNotificationTemplatesStarted will disassociate a started notification_template for a job_template.
func (jt *JobTemplateNotificationTemplatesService) DisassociateJobTemplateNotificationTemplatesStarted(jobTemplateID int, notificationTemplateID int) (*NotificationTemplate, error) {
	return jt.disassociateJobTemplateNotificationTemplatesForType(jobTemplateID, notificationTemplateID, "started")
}
