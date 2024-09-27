package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
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
func (cs *CredentialTypeService) ListCredentialTypes(params map[string]string) ([]*CredentialType, error) {

	results, err := cs.getAllPages(credentialTypesAPIEndpoint, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (cs *CredentialTypeService) getAllPages(firstURL string, params map[string]string) ([]*CredentialType, error) {
	results := make([]*CredentialType, 0)
	nextURL := firstURL
	for {
		nextURLParsed, err := url.Parse(nextURL)
		if err != nil {
			return nil, err
		}

		nextURLQueryParams := make(map[string]string)
		for paramName, paramValues := range nextURLParsed.Query() {
			if len(paramValues) > 0 {
				nextURLQueryParams[paramName] = paramValues[0]
			}
		}

		for paramName, paramValue := range params {
			nextURLQueryParams[paramName] = paramValue
		}

		result := new(ListCredentialTypeResponse)
		resp, err := cs.client.Requester.GetJSON(nextURLParsed.Path, result, nextURLQueryParams)
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

		results = append(results, result.Results...)

		if result.Next == nil || result.Next.(string) == "" {
			break
		}
		nextURL = result.Next.(string)
	}
	return results, nil
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

// GetCredentialTypeByName : Fetches a credential type by Name.
func (cs *CredentialTypeService) GetCredentialTypeByName(name string, params map[string]string) (*CredentialType, error) {
	credentialTypes, err := cs.ListCredentialTypes(params)
	if err != nil {
		return nil, err
	}

	for _, credentialType := range credentialTypes {
		if credentialType.Name == name {
			return credentialType, nil
		}
	}
	return nil, fmt.Errorf("could not find credential type with name %s", name)
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
