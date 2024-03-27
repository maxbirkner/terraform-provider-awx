package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// SettingService implements awx settings apis.
type SettingService struct {
	client *Client
}

// ListSettingsResponse represents `ListSettings` endpoint response.
type ListSettingsResponse struct {
	Pagination
	Results []*SettingSummary `json:"results"`
}

const settingsAPIEndpoint = "/api/v2/settings/"

// ListSettings shows list of awx settings.
func (p *SettingService) ListSettings(params map[string]string) ([]*SettingSummary, *ListSettingsResponse, error) {
	result := new(ListSettingsResponse)
	resp, err := p.client.Requester.GetJSON(settingsAPIEndpoint, result, params)
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

// GetSettingsBySlug shows the details of a setting.
func (p *SettingService) GetSettingsBySlug(slug string, params map[string]string) (*Setting, error) {
	result := new(Setting)
	endpoint := fmt.Sprintf("%s%s/", settingsAPIEndpoint, slug)
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

// UpdateSettings update an awx Setting.
func (p *SettingService) UpdateSettings(slug string, data map[string]interface{}, params map[string]string) (*Setting, error) {
	result := new(Setting)
	endpoint := fmt.Sprintf("%s%s", settingsAPIEndpoint, slug)
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

// DeleteSettings delete an awx Setting.
func (p *SettingService) DeleteSettings(slug string) (*Setting, error) {
	result := new(Setting)
	endpoint := fmt.Sprintf("%s%s", settingsAPIEndpoint, slug)

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
