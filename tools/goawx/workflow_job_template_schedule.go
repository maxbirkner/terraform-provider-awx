package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const workflowJobTemplateSchedulesAPIEndpoint = "/api/v2/workflow_job_templates/%d/schedules/"

// WorkflowJobTemplateScheduleService implements awx job template nodes apis.
type WorkflowJobTemplateScheduleService struct {
	client *Client
}

// ListWorkflowJobTemplateSchedules shows a list of schedules for a given workflow_job_template.
func (jt *WorkflowJobTemplateScheduleService) ListWorkflowJobTemplateSchedules(id int, params map[string]string) ([]*Schedule, *ListSchedulesResponse, error) {
	result := new(ListSchedulesResponse)
	resp, err := jt.client.Requester.GetJSON(
		fmt.Sprintf(workflowJobTemplateSchedulesAPIEndpoint, id),
		result, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return nil, result, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, result, err
	}

	return result.Results, result, nil
}

// CreateWorkflowJobTemplateSchedule will create a schedule for an existing workflow_job_template.
func (jt *WorkflowJobTemplateScheduleService) CreateWorkflowJobTemplateSchedule(id int, data map[string]interface{}, params map[string]string) (*Schedule, error) {
	mandatoryFields = []string{"name", "rrule"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Schedule)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := jt.client.Requester.PostJSON(
		fmt.Sprintf(workflowJobTemplateSchedulesAPIEndpoint, id),
		bytes.NewReader(payload), result, params)
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
