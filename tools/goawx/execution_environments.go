package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ExecutionEnvironmentsService implements awx execution environments apis.
type ExecutionEnvironmentsService struct {
	client *Client
}

// ListExecutionEnvironmentsResponse represents `ListExecutionEnvironments` endpoint response.
type ListExecutionEnvironmentsResponse struct {
	Pagination
	Results []*ExecutionEnvironment `json:"results"`
}

const executionEnvironmentsAPIEndpoint = "/api/v2/execution_environments/"

// ListExecutionEnvironments shows list of awx execution environments.
func (p *ExecutionEnvironmentsService) ListExecutionEnvironments(params map[string]string) ([]*ExecutionEnvironment, *ListExecutionEnvironmentsResponse, error) {
	result := new(ListExecutionEnvironmentsResponse)
	resp, err := p.client.Requester.GetJSON(executionEnvironmentsAPIEndpoint, result, params)
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

// GetExecutionEnvironmentByID shows the details of a ExecutionEnvironment.
func (p *ExecutionEnvironmentsService) GetExecutionEnvironmentByID(id int, params map[string]string) (*ExecutionEnvironment, error) {
	result := new(ExecutionEnvironment)
	endpoint := fmt.Sprintf("%s%d/", executionEnvironmentsAPIEndpoint, id)
	resp, err := p.client.Requester.GetJSON(endpoint, result, params)
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

// CreateExecutionEnvironment creates an awx ExecutionEnvironment.
func (p *ExecutionEnvironmentsService) CreateExecutionEnvironment(data map[string]interface{}, params map[string]string) (*ExecutionEnvironment, error) {
	mandatoryFields = []string{"name", "image"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(ExecutionEnvironment)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Requester.PostJSON(executionEnvironmentsAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateExecutionEnvironment update an awx ExecutionEnvironment.
func (p *ExecutionEnvironmentsService) UpdateExecutionEnvironment(id int, data map[string]interface{}, params map[string]string) (*ExecutionEnvironment, error) {
	result := new(ExecutionEnvironment)
	endpoint := fmt.Sprintf("%s%d", executionEnvironmentsAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// DeleteExecutionEnvironment delete an awx ExecutionEnvironment.
func (p *ExecutionEnvironmentsService) DeleteExecutionEnvironment(id int) (*ExecutionEnvironment, error) {
	result := new(ExecutionEnvironment)
	endpoint := fmt.Sprintf("%s%d", executionEnvironmentsAPIEndpoint, id)

	resp, err := p.client.Requester.Delete(endpoint, result, nil)
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
