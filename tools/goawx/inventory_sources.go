package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// InventorySourcesService implements awx inventory sources apis.
type InventorySourcesService struct {
	client *Client
}

// ListInventorySourcesResponse represents `ListInventorySources` endpoint response.
type ListInventorySourcesResponse struct {
	Pagination
	Results []*InventorySource `json:"results"`
}

const inventorySourcesAPIEndpoint = "/api/v2/inventory_sources/"

// GetInventorySourceByID shows the details of a awx inventory sources.
func (i *InventorySourcesService) GetInventorySourceByID(id int, params map[string]string) (*InventorySource, error) {
	result := new(InventorySource)
	endpoint := fmt.Sprintf("%s%d/", inventorySourcesAPIEndpoint, id)
	resp, err := i.client.Requester.GetJSON(endpoint, result, params)
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

// ListInventorySources shows list of awx inventories.
func (i *InventorySourcesService) ListInventorySources(params map[string]string) ([]*InventorySource, *ListInventorySourcesResponse, error) {
	result := new(ListInventorySourcesResponse)
	resp, err := i.client.Requester.GetJSON(inventorySourcesAPIEndpoint, result, params)
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

// CreateInventorySource creates an awx InventorySource.
func (i *InventorySourcesService) CreateInventorySource(data map[string]interface{}, params map[string]string) (*InventorySource, error) {
	mandatoryFields = []string{"name", "inventory"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(InventorySource)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if InventorySource exists and return proper error

	resp, err := i.client.Requester.PostJSON(inventorySourcesAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateInventorySource update an awx InventorySource.
func (i *InventorySourcesService) UpdateInventorySource(id int, data map[string]interface{}, _ map[string]string) (*InventorySource, error) {
	result := new(InventorySource)
	endpoint := fmt.Sprintf("%s%d", inventorySourcesAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := i.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// GetInventorySource retrieves the InventorySource information from its ID or Name.
func (i *InventorySourcesService) GetInventorySource(id int, _ map[string]string) (*InventorySource, error) {
	endpoint := fmt.Sprintf("%s%d", inventorySourcesAPIEndpoint, id)
	result := new(InventorySource)
	resp, err := i.client.Requester.GetJSON(endpoint, result, map[string]string{})
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
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

// DeleteInventorySource delete an InventorySource from AWX.
func (i *InventorySourcesService) DeleteInventorySource(id int) (*InventorySource, error) {
	result := new(InventorySource)
	endpoint := fmt.Sprintf("%s%d", inventorySourcesAPIEndpoint, id)

	resp, err := i.client.Requester.Delete(endpoint, result, nil)
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
