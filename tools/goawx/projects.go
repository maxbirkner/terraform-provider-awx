package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ProjectService implements awx projects apis.
type ProjectService struct {
	client *Client
}

// ListProjectsResponse represents `ListProjects` endpoint response.
type ListProjectsResponse struct {
	Pagination
	Results []*Project `json:"results"`
}

const projectsAPIEndpoint = "/api/v2/projects/"

// ListProjects shows list of awx projects.
func (p *ProjectService) ListProjects(params map[string]string) ([]*Project, *ListProjectsResponse, error) {
	result := new(ListProjectsResponse)
	resp, err := p.client.Requester.GetJSON(projectsAPIEndpoint, result, params)
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

// GetProjectByID shows the details of a project.
func (p *ProjectService) GetProjectByID(id int, params map[string]string) (*Project, error) {
	result := new(Project)
	endpoint := fmt.Sprintf("%s%d/", projectsAPIEndpoint, id)
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

// CreateProject creates an awx project.
func (p *ProjectService) CreateProject(data map[string]interface{}, params map[string]string) (*Project, error) {
	mandatoryFields = []string{"name", "organization", "scm_type"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Project)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if project exists and return proper error

	resp, err := p.client.Requester.PostJSON(projectsAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateProject update an awx Project.
func (p *ProjectService) UpdateProject(id int, data map[string]interface{}, params map[string]string) (*Project, error) {
	result := new(Project)
	endpoint := fmt.Sprintf("%s%d", projectsAPIEndpoint, id)
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

// DeleteProject delete an awx Project.
func (p *ProjectService) DeleteProject(id int) (*Project, error) {
	result := new(Project)
	endpoint := fmt.Sprintf("%s%d", projectsAPIEndpoint, id)

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
