package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CredentialTypeService implements awx CredentialType apis.
type CredentialTypeService struct {
	client *Client
}

// ListCredentialTypeResponse represents `ListCredentialTypes` endpoint response.
type ListCredentialTypeResponse struct {
	Pagination
	Results []*CredentialType `json:"results"`
}

const credentialTypesAPIEndpoint = "/api/v2/credential_types/" //nolint:gosec

// ListCredentialTypes shows list of awx CredentialTypes.
func (cs *CredentialTypeService) ListCredentialTypes(params map[string]string) ([]*CredentialType,
	*ListCredentialTypeResponse, error) {
	result := new(ListCredentialTypeResponse)
	resp, err := cs.client.Requester.GetJSON(credentialTypesAPIEndpoint, result, params)
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

	err = CheckResponse(resp)
	if err != nil {
		return nil, result, err
	}

	return result.Results, result, nil
}

// CreateCredentialType : Creates a new credential type in AWX.
func (cs *CredentialTypeService) CreateCredentialType(data map[string]interface{}, params map[string]string) (*CredentialType, error) {
	result := new(CredentialType)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PostJSON(credentialTypesAPIEndpoint, bytes.NewReader(payload), result, params)
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

// GetCredentialTypeByID : Fetches a credential type by ID.
func (cs *CredentialTypeService) GetCredentialTypeByID(id int, params map[string]string) (*CredentialType, error) {
	result := new(CredentialType)
	endpoint := fmt.Sprintf("%s%d", credentialTypesAPIEndpoint, id)
	resp, err := cs.client.Requester.GetJSON(endpoint, result, params)
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

// UpdateCredentialTypeByID : Updates a credential type by ID.
func (cs *CredentialTypeService) UpdateCredentialTypeByID(id int, data map[string]interface{}, params map[string]string) (*CredentialType, error) {
	result := new(CredentialType)
	endpoint := fmt.Sprintf("%s%d", credentialTypesAPIEndpoint, id)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PutJSON(endpoint, bytes.NewReader(payload), result, params)
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

// DeleteCredentialTypeByID : Deletes a credential type by ID.
func (cs *CredentialTypeService) DeleteCredentialTypeByID(id int, params map[string]string) error {
	endpoint := fmt.Sprintf("%s%d", credentialTypesAPIEndpoint, id)
	resp, err := cs.client.Requester.Delete(endpoint, nil, params)
	if resp != nil {
		func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err != nil {
		return err
	}

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	return nil
}
