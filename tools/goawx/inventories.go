package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// InventoriesService implements awx inventories apis.
type InventoriesService struct {
	client *Client
}

// ListInventoriesResponse represents `ListInventories` endpoint response.
type ListInventoriesResponse struct {
	Pagination
	Results []*Inventory `json:"results"`
}

const inventoriesAPIEndpoint = "/api/v2/inventories/"

// GetInventoryByID shows the details of a awx inventroy sources.
func (i *InventoriesService) GetInventoryByID(id int, params map[string]string) (*Inventory, error) {
	result := new(Inventory)
	endpoint := fmt.Sprintf("%s%d/", inventoriesAPIEndpoint, id)
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

// ListInventories shows list of awx inventories.
func (i *InventoriesService) ListInventories(params map[string]string) ([]*Inventory, *ListInventoriesResponse, error) {
	result := new(ListInventoriesResponse)
	resp, err := i.client.Requester.GetJSON(inventoriesAPIEndpoint, result, params)
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

// CreateInventory creates an awx inventory.
func (i *InventoriesService) CreateInventory(data map[string]interface{}, params map[string]string) (*Inventory, error) {
	mandatoryFields = []string{"name", "organization"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Inventory)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if inventory exists and return proper error

	resp, err := i.client.Requester.PostJSON(inventoriesAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateInventory update an awx inventory.
func (i *InventoriesService) UpdateInventory(id int, data map[string]interface{}, _ map[string]string) (*Inventory, error) {
	result := new(Inventory)
	endpoint := fmt.Sprintf("%s%d", inventoriesAPIEndpoint, id)
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

// GetInventory retrieves the inventory information from its ID or Name.
func (i *InventoriesService) GetInventory(id int, _ map[string]string) (*Inventory, error) {
	endpoint := fmt.Sprintf("%s%d", inventoriesAPIEndpoint, id)
	result := new(Inventory)
	resp, err := i.client.Requester.GetJSON(endpoint, result, map[string]string{})
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

// DeleteInventory delete an inventory from AWX.
func (i *InventoriesService) DeleteInventory(id int) (*Inventory, error) {
	result := new(Inventory)
	endpoint := fmt.Sprintf("%s%d", inventoriesAPIEndpoint, id)

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

// DisAssociateInstanceGroups remove InstanceGroup from an awx Inventory.
func (i *InventoriesService) DisAssociateInstanceGroups(id int, data map[string]interface{}, _ map[string]string) (*Inventory, error) {
    result := new(Inventory)
    endpoint := fmt.Sprintf("%s%d/instance_groups/", inventoriesAPIEndpoint, id)
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
    resp, err := i.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// AssociateInstanceGroups  adding InstanceGroup to Inventory.
func (i *InventoriesService) AssociateInstanceGroups(id int, data map[string]interface{}, _ map[string]string) (*Inventory, error) {
    result := new(Inventory)

    endpoint := fmt.Sprintf("%s%d/instance_groups/", inventoriesAPIEndpoint, id)
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
    resp, err := i.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
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
