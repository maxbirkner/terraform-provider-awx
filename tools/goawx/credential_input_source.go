// Package awx provides a client for using the Ansible Tower / AWX REST API.
package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CredentialInputSourceService implements awx credential input source api endpoints.
type CredentialInputSourceService struct {
	client *Client
}

// ListCredentialInputSourceResponse represents `ListCredentialInputSource` endpoint response.
type ListCredentialInputSourceResponse struct {
	Pagination
	Results []*CredentialInputSource `json:"results"`
}

const credentialInputSourceAPIEndpoint = "/api/v2/credential_input_sources/" //nolint:gosec

// ListCredentialInputSources shows list of awx credential input sources.
func (cs *CredentialInputSourceService) ListCredentialInputSources(params map[string]string) ([]*CredentialInputSource,
	*ListCredentialInputSourceResponse,
	error) {
	result := new(ListCredentialInputSourceResponse)
	resp, err := cs.client.Requester.GetJSON(credentialInputSourceAPIEndpoint, result, params)
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

// CreateCredentialInputSource creates an awx credential input source.
func (cs *CredentialInputSourceService) CreateCredentialInputSource(data map[string]interface{}, params map[string]string) (*CredentialInputSource, error) {
	result := new(CredentialInputSource)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PostJSON(credentialInputSourceAPIEndpoint, bytes.NewReader(payload), result, params)
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

// GetCredentialInputSourceByID : Gets a specific input source by ID.
func (cs *CredentialInputSourceService) GetCredentialInputSourceByID(id int, params map[string]string) (*CredentialInputSource, error) {
	result := new(CredentialInputSource)
	endpoint := fmt.Sprintf("%s%d", credentialInputSourceAPIEndpoint, id)
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

// UpdateCredentialInputSourceByID : Updates an input source by ID.
func (cs *CredentialInputSourceService) UpdateCredentialInputSourceByID(id int, data map[string]interface{},
	params map[string]string) (*CredentialInputSource, error) {
	result := new(CredentialInputSource)
	endpoint := fmt.Sprintf("%s%d", credentialInputSourceAPIEndpoint, id)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, params)
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

// DeleteCredentialInputSourceByID : Deletes an input source by ID.
func (cs *CredentialInputSourceService) DeleteCredentialInputSourceByID(id int, params map[string]string) error {
	endpoint := fmt.Sprintf("%s%d", credentialInputSourceAPIEndpoint, id)
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
