package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Enum of job statuses.
const (
	JobStatusNew        = "new"
	JobStatusPending    = "pending"
	JobStatusWaiting    = "waiting"
	JobStatusRunning    = "running"
	JobStatusSuccessful = "successful"
	JobStatusFailed     = "failed"
	JobStatusError      = "error"
	JobStatusCanceled   = "canceled"
)

// JobService implements awx job apis.
type JobService struct {
	client *Client
}

// HostSummariesResponse represents `JobHostSummaries` endpoint response.
type HostSummariesResponse struct {
	Pagination
	Results []HostSummary `json:"results"`
}

// JobEventsResponse represents `JobEvents` endpoint response.
type JobEventsResponse struct {
	Pagination
	Results []JobEvent `json:"results"`
}

// CancelJobResponse represents `CancelJob` endpoint response.
type CancelJobResponse struct {
	Detail string `json:"detail"`
}

const jobAPIEndpoint = "/api/v2/jobs/"

// GetJob shows the details of a job.
func (j *JobService) GetJob(id int, params map[string]string) (*Job, error) {
	result := new(Job)
	endpoint := fmt.Sprintf("%s%d/", jobAPIEndpoint, id)
	resp, err := j.client.Requester.GetJSON(endpoint, result, params)
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

// CancelJob cancels a job.
func (j *JobService) CancelJob(id int, data map[string]interface{}, params map[string]string) (*CancelJobResponse, error) {
	result := new(CancelJobResponse)
	endpoint := fmt.Sprintf("%s%d/cancel/", jobAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := j.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, params)
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

// RelaunchJob relaunch a job.
func (j *JobService) RelaunchJob(id int, data map[string]interface{}, params map[string]string) (*JobLaunch, error) {
	result := new(JobLaunch)
	endpoint := fmt.Sprintf("%s%d/relaunch/", jobAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := j.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, params)
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

// GetHostSummaries get a job hosts summaries.
func (j *JobService) GetHostSummaries(id int, params map[string]string) ([]HostSummary, *HostSummariesResponse, error) {
	result := new(HostSummariesResponse)
	endpoint := fmt.Sprintf("%s%d/job_host_summaries/", jobAPIEndpoint, id)
	resp, err := j.client.Requester.GetJSON(endpoint, result, params)
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

// GetJobEvents get a list of job events.
func (j *JobService) GetJobEvents(id int, params map[string]string) ([]JobEvent, *JobEventsResponse, error) {
	result := new(JobEventsResponse)
	endpoint := fmt.Sprintf("%s%d/job_events/", jobAPIEndpoint, id)
	resp, err := j.client.Requester.GetJSON(endpoint, result, params)
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
