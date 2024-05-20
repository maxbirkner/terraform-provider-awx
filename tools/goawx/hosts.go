package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// HostService implements awx Hosts apis.
type HostService struct {
	client *Client
}

// AssociateGroup implement the awx group association request.
type AssociateGroup struct {
	ID        int  `json:"id"`
	Associate bool `json:"associate"`
}

// ListHostsResponse represents `ListHosts` endpoint response.
type ListHostsResponse struct {
	Pagination
	Results []*Host `json:"results"`
}

const hostsAPIEndpoint = "/api/v2/hosts/"

// GetHostByID shows the details of a awx inventroy sources.
func (h *HostService) GetHostByID(id int, params map[string]string) (*Host, error) {
	result := new(Host)
	endpoint := fmt.Sprintf("%s%d/", hostsAPIEndpoint, id)
	resp, err := h.client.Requester.GetJSON(endpoint, result, params)
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

// ListHosts shows list of awx Hosts.
func (h *HostService) ListHosts(params map[string]string) ([]*Host, *ListHostsResponse, error) {
	result := new(ListHostsResponse)
	resp, err := h.client.Requester.GetJSON(hostsAPIEndpoint, result, params)
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

// CreateHost creates an awx Host.
func (h *HostService) CreateHost(data map[string]interface{}, params map[string]string) (*Host, error) {
	mandatoryFields = []string{"name", "inventory"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Host)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if Host exists and return proper error

	resp, err := h.client.Requester.PostJSON(hostsAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateHost update an awx Host.
func (h *HostService) UpdateHost(id int, data map[string]interface{}, _ map[string]string) (*Host, error) {
	result := new(Host)
	endpoint := fmt.Sprintf("%s%d", hostsAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// AssociateGroup update an awx Host.
func (h *HostService) AssociateGroup(id int, data map[string]interface{}, _ map[string]string) (*Host, error) {
	result := new(Host)
	endpoint := fmt.Sprintf("%s%d/groups/", hostsAPIEndpoint, id)
	data["associate"] = true
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// DisAssociateGroup update an awx Host.
func (h *HostService) DisAssociateGroup(id int, data map[string]interface{}, _ map[string]string) (*Host, error) {
	result := new(Host)
	endpoint := fmt.Sprintf("%s%d/groups/", hostsAPIEndpoint, id)
	data["disassociate"] = true
	mandatoryFields = []string{"id"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// DeleteHost delete an awx Host.
func (h *HostService) DeleteHost(id int) (*Host, error) {
	result := new(Host)
	endpoint := fmt.Sprintf("%s%d", hostsAPIEndpoint, id)

	resp, err := h.client.Requester.Delete(endpoint, result, nil)
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
