package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const workflowJobPrefix = "workflow_"
const surveySpecAPIEndpoint = "/api/v2/%sjob_templates/%d/survey_spec/"

// SurveySpecService implements awx job template nodes apis.
type SurveySpecService struct {
	client *Client
}

func computeEndpoint(isWorkflow bool, jobTemplateID int) string {
	if isWorkflow {
		return fmt.Sprintf(surveySpecAPIEndpoint, workflowJobPrefix, jobTemplateID)
	}
	return fmt.Sprintf(surveySpecAPIEndpoint, "", jobTemplateID)
}

func (jts *SurveySpecService) GetSurveySpec(isWorkflow bool, jobTemplateID int, params map[string]string) (*SurveySpec, error) {
	endpoint := computeEndpoint(isWorkflow, jobTemplateID)
	result := new(SurveySpec)
	resp, err := jts.client.Requester.GetJSON(endpoint, result, params)
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

func (jts *SurveySpecService) CreateSurveySpec(isWorkflow bool, jobTemplateID int, data map[string]interface{}) (*SurveySpec, error) {
	endpoint := computeEndpoint(isWorkflow, jobTemplateID)
	result := new(SurveySpec)
	mandatoryFields = []string{"name", "description", "spec"}
	validate, status := ValidateParams(data, mandatoryFields)
	if !status {
		err := fmt.Errorf("mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := jts.client.Requester.PostJSON(endpoint, bytes.NewReader(payload), result, map[string]string{})
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

func (jts *SurveySpecService) DeleteSurveySpec(isWorkflow bool, jobTemplateID int) error {
	endpoint := computeEndpoint(isWorkflow, jobTemplateID)
	resp, err := jts.client.Requester.Delete(endpoint, nil, map[string]string{})
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
