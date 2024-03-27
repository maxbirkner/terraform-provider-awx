package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ApplicationService implements awx application api endpoints.
type ApplicationService struct {
	client *Client
}

// ListApplicationResponse represents `ListApplication` endpoint response.
type ListApplicationResponse struct {
	Pagination
	Results []*Application `json:"results"`
}

const applicationAPIEndpoint = "/api/v2/applications/"

// ListApplication shows list of awx authentication applications.
func (c *ApplicationService) ListApplication(params map[string]string) ([]*Application, *ListApplicationResponse, error) {
	result := new(ListApplicationResponse)
	resp, err := c.client.Requester.GetJSON(applicationAPIEndpoint, result, params)
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

// GetApplicationByID shows an of awx application by its ID.
func (c *ApplicationService) GetApplicationByID(id int, params map[string]string) (*Application, error) {
	result := new(Application)
	endpoint := fmt.Sprintf("%s%d", applicationAPIEndpoint, id)
	resp, err := c.client.Requester.GetJSON(endpoint, result, params)
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

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateApplication creates an awx authentication application.
func (c *ApplicationService) CreateApplication(data map[string]interface{}, params map[string]string) (*Application, error) {
	mandatoryFields = []string{"name", "client_type", "authorization_grant_type", "organization"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Application)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if Application exists and return proper error

	resp, err := c.client.Requester.PostJSON(applicationAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateApplication update an awx application.
func (c *ApplicationService) UpdateApplication(id int, data map[string]interface{}, _ map[string]string) (*Application, error) {
	result := new(Application)
	endpoint := fmt.Sprintf("%s%d", applicationAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Requester.PutJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// DeleteApplication delete an awx application.
func (c *ApplicationService) DeleteApplication(id int) (*Application, error) {
	result := new(Application)
	endpoint := fmt.Sprintf("%s%d", applicationAPIEndpoint, id)

	resp, err := c.client.Requester.Delete(endpoint, result, nil)
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
