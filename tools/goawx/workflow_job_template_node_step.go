package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// WorkflowJobTemplateNodeStepService implements awx job template nodes apis.
type WorkflowJobTemplateNodeStepService struct {
	endpoint string
	client   *Client
}

// ListWorkflowJobTemplateNodes shows a list of job templates nodes.
func (jt *WorkflowJobTemplateNodeStepService) ListWorkflowJobTemplateNodes(id int, params map[string]string) ([]*WorkflowJobTemplateNode, *ListWorkflowJobTemplateNodesResponse, error) {
	workflowJobTemplateNodesActionEndpoint := fmt.Sprintf(jt.endpoint, id)
	return fetchWorkflowJobTemplateNode(jt.client, params, workflowJobTemplateNodesActionEndpoint)
}

func fetchWorkflowJobTemplateNode(client *Client, params map[string]string, workflowJobTemplateNodesActionEndpoint string) ([]*WorkflowJobTemplateNode, *ListWorkflowJobTemplateNodesResponse, error) {
	result := new(ListWorkflowJobTemplateNodesResponse)
	resp, err := client.Requester.GetJSON(workflowJobTemplateNodesActionEndpoint, result, params)
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

func createWorkflowJobTemplateNode(client *Client, data map[string]interface{}, params map[string]string, workflowJobTemplateNodesActionEndpoint string) (*WorkflowJobTemplateNode, error) {
	result := new(WorkflowJobTemplateNode)
	mandatoryFields = []string{"unified_job_template", "identifier"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := client.Requester.PostJSON(workflowJobTemplateNodesActionEndpoint, bytes.NewReader(payload), result, params)
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
	log.Printf("Created ID %v", result.ID)
	return result, nil
}

// CreateWorkflowJobTemplateNodeStep will be creating a template node for an existing node.
func (jt *WorkflowJobTemplateNodeStepService) CreateWorkflowJobTemplateNodeStep(id int, data map[string]interface{}, params map[string]string) (*WorkflowJobTemplateNode, error) {
	workflowJobTemplateNodesActionEndpoint := fmt.Sprintf(jt.endpoint, id)
	return createWorkflowJobTemplateNode(jt.client, data, params, workflowJobTemplateNodesActionEndpoint)
}
