package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// UserService implements awx Users apis.
type UserService struct {
	client *Client
}

// ListUsersResponse represents `ListUsers` endpoint response.
type ListUsersResponse struct {
	Pagination
	Results []*User `json:"results"`
}

// ListUsersEntitlementsResponse represents a list of ApplyRole for a User.
type ListUsersEntitlementsResponse struct {
	Pagination
	Results []*ApplyRole `json:"results"`
}

const usersAPIEndpoint = "/api/v2/users/"

// ListUsers shows list of awx Users.
func (u *UserService) ListUsers(params map[string]string) ([]*User, *ListUsersResponse, error) {
	result := new(ListUsersResponse)
	resp, err := u.client.Requester.GetJSON(usersAPIEndpoint, result, params)
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

// CreateUser creates an awx User.
func (u *UserService) CreateUser(data map[string]interface{}, params map[string]string) (*User, error) {
	mandatoryFields = []string{"username", "password", "first_name", "last_name", "email"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(User)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if User exists and return proper error

	resp, err := u.client.Requester.PostJSON(usersAPIEndpoint, bytes.NewReader(payload), result, params)
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

// UpdateUser update an awx user.
func (u *UserService) UpdateUser(id int, data map[string]interface{}, params map[string]string) (*User, error) {
	result := new(User)
	endpoint := fmt.Sprintf("%s%d", usersAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Requester.PutJSON(endpoint, bytes.NewReader(payload), result, nil)
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

// DeleteUser delete an awx User.
func (u *UserService) DeleteUser(id int) (*User, error) {
	result := new(User)
	endpoint := fmt.Sprintf("%s%d", usersAPIEndpoint, id)

	resp, err := u.client.Requester.Delete(endpoint, result, nil)
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

// GetUserByID read an awx User.
func (u *UserService) GetUserByID(id int, params map[string]string) (*User, error) {
	result := new(User)
	endpoint := fmt.Sprintf("%s%d", usersAPIEndpoint, id)
	resp, err := u.client.Requester.GetJSON(endpoint, result, params)
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

// ListUserRoleEntitlements shows list of awx User Role Entitlements.
func (u *UserService) ListUserRoleEntitlements(id int, params map[string]string) ([]*ApplyRole, *ListUsersEntitlementsResponse, error) {
	result := new(ListUsersEntitlementsResponse)
	endpoint := fmt.Sprintf("%s%d/roles/", usersAPIEndpoint, id)
	resp, err := u.client.Requester.GetJSON(endpoint, result, params)
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

// UpdateUserRoleEntitlement updates an awx user role entitlement.
func (u *UserService) UpdateUserRoleEntitlement(id int, data map[string]interface{}, params map[string]string) (interface{}, error) {
	result := new(interface{})
	endpoint := fmt.Sprintf("%s%d/roles/", usersAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, nil)
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
