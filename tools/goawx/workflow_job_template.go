package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// WorkflowJobTemplateService implements awx workflow job template apis.
type WorkflowJobTemplateService struct {
	client *Client
}

// ListWorkflowJobTemplatesResponse represents `ListWorkflowJobTemplate` endpoint response.
type ListWorkflowJobTemplatesResponse struct {
	Pagination
	Results []*WorkflowJobTemplate `json:"results"`
}

const workflowJobTemplateAPIEndpoint = "/api/v2/workflow_job_templates/"

// GetWorkflowJobTemplateByID shows the details of a workflow job template.
func (jt *WorkflowJobTemplateService) GetWorkflowJobTemplateByID(id int, params map[string]string) (*WorkflowJobTemplate, error) {
	result := new(WorkflowJobTemplate)
	endpoint := fmt.Sprintf("%s%d/", workflowJobTemplateAPIEndpoint, id)
	resp, err := jt.client.Requester.GetJSON(endpoint, result, params)
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

// ListWorkflowJobTemplates shows a list of workflow job templates.
func (jt *WorkflowJobTemplateService) ListWorkflowJobTemplates(params map[string]string) ([]*WorkflowJobTemplate, *ListWorkflowJobTemplatesResponse, error) {
	result := new(ListWorkflowJobTemplatesResponse)
	resp, err := jt.client.Requester.GetJSON(workflowJobTemplateAPIEndpoint, result, params)
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

// CreateWorkflowJobTemplate creates a workflow job template.
func (jt *WorkflowJobTemplateService) CreateWorkflowJobTemplate(data map[string]interface{}, params map[string]string) (*WorkflowJobTemplate, error) {
	result := new(WorkflowJobTemplate)
	mandatoryFields = []string{"name"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := jt.client.Requester.PostJSON(workflowJobTemplateAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateWorkflowJobTemplate updates a workflow job template.
func (jt *WorkflowJobTemplateService) UpdateWorkflowJobTemplate(id int, data map[string]interface{}, params map[string]string) (*WorkflowJobTemplate, error) {
	result := new(WorkflowJobTemplate)
	endpoint := fmt.Sprintf("%s%d", workflowJobTemplateAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := jt.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, params)
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

// DeleteWorkflowJobTemplate deletes a workflow job template.
func (jt *WorkflowJobTemplateService) DeleteWorkflowJobTemplate(id int) (*WorkflowJobTemplate, error) {
	result := new(WorkflowJobTemplate)
	endpoint := fmt.Sprintf("%s%d", workflowJobTemplateAPIEndpoint, id)

	resp, err := jt.client.Requester.Delete(endpoint, result, nil)
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

// Launch a job with the workflow job template.
func (jt *WorkflowJobTemplateService) Launch(id int, data map[string]interface{}, params map[string]string) (*JobLaunch, error) {
	result := new(JobLaunch)
	endpoint := fmt.Sprintf("%s%d/launch/", workflowJobTemplateAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := jt.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, params)
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
