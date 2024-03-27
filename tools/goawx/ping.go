package awx

import "fmt"

// PingService implements awx ping apis.
type PingService struct {
	client *Client
}

const pingAPIEndpoint = "/api/v2/ping/"

// Ping do ping with awx servers.
func (p *PingService) Ping() (*Ping, error) {
	result := new(Ping)
	resp, err := p.client.Requester.GetJSON(pingAPIEndpoint, result, map[string]string{})
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
