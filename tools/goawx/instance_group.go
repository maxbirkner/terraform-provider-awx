package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// InstanceGroupsService implements awx execution environments apis.
type InstanceGroupsService struct {
	client *Client
}

// ListInstanceGroupsResponse represents `ListInstanceGroups` endpoint response.
type ListInstanceGroupsResponse struct {
	Pagination
	Results []*InstanceGroup `json:"results"`
}

// InstanceGroupsAPIEndpoint represents the endpoint for the InstanceGroups.
const InstanceGroupsAPIEndpoint = "/api/v2/instance_groups/"

// ListInstanceGroups shows list of awx execution environments.
func (p *InstanceGroupsService) ListInstanceGroups(params map[string]string) ([]*InstanceGroup, *ListInstanceGroupsResponse, error) {
	result := new(ListInstanceGroupsResponse)
	resp, err := p.client.Requester.GetJSON(InstanceGroupsAPIEndpoint, result, params)
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

// GetInstanceGroupByID shows the details of a InstanceGroup.
func (p *InstanceGroupsService) GetInstanceGroupByID(id int, params map[string]string) (*InstanceGroup, error) {
	result := new(InstanceGroup)
	endpoint := fmt.Sprintf("%s%d/", InstanceGroupsAPIEndpoint, id)
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

// CreateInstanceGroup creates an awx InstanceGroup.
func (p *InstanceGroupsService) CreateInstanceGroup(data map[string]interface{}, params map[string]string) (*InstanceGroup, error) {
	mandatoryFields = []string{"name"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(InstanceGroup)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Requester.PostJSON(InstanceGroupsAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateInstanceGroup update an awx InstanceGroup.
func (p *InstanceGroupsService) UpdateInstanceGroup(id int, data map[string]interface{}, _ map[string]string) (*InstanceGroup, error) {
	result := new(InstanceGroup)
	endpoint := fmt.Sprintf("%s%d", InstanceGroupsAPIEndpoint, id)
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

// DeleteInstanceGroup delete an awx InstanceGroup.
func (p *InstanceGroupsService) DeleteInstanceGroup(id int) (*InstanceGroup, error) {
	result := new(InstanceGroup)
	endpoint := fmt.Sprintf("%s%d", InstanceGroupsAPIEndpoint, id)

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
