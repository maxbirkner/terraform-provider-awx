package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GroupService implements awx Groups apis.
type GroupService struct {
	client *Client
}

// ListGroupsResponse represents `ListGroups` endpoint response.
type ListGroupsResponse struct {
	Pagination
	Results []*Group `json:"results"`
}

const groupsAPIEndpoint = "/api/v2/groups/"

// GetGroupByID shows the details of a awx group.
func (g *GroupService) GetGroupByID(id int, params map[string]string) (*Group, error) {
	result := new(Group)
	endpoint := fmt.Sprintf("%s%d/", groupsAPIEndpoint, id)
	resp, err := g.client.Requester.GetJSON(endpoint, result, params)
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

// ListGroups shows list of awx Groups.
func (g *GroupService) ListGroups(params map[string]string) ([]*Group, *ListGroupsResponse, error) {
	result := new(ListGroupsResponse)
	resp, err := g.client.Requester.GetJSON(groupsAPIEndpoint, result, params)
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

// CreateGroup creates an awx Group.
func (g *GroupService) CreateGroup(data map[string]interface{}, params map[string]string) (*Group, error) {
	mandatoryFields = []string{"name", "inventory"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Group)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if Group exists and return proper error

	resp, err := g.client.Requester.PostJSON(groupsAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateGroup update an awx group.
func (g *GroupService) UpdateGroup(id int, data map[string]interface{}, _ map[string]string) (*Group, error) {
	result := new(Group)
	endpoint := fmt.Sprintf("%s%d", groupsAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// DeleteGroup delete an awx Group.
func (g *GroupService) DeleteGroup(id int) (*Group, error) {
	result := new(Group)
	endpoint := fmt.Sprintf("%s%d", groupsAPIEndpoint, id)

	resp, err := g.client.Requester.Delete(endpoint, result, nil)
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
