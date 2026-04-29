package awx

import (
	"fmt"
	"net/url"
)

func getAllLabelPages(requester *Requester, firstURL string, params map[string]string) ([]*Label, error) {
	results := make([]*Label, 0)
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

		result := new(ListLabelsResponse)
		resp, err := requester.GetJSON(nextURLParsed.Path, result, nextURLQueryParams)
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

		nextURLValue, ok := result.Next.(string)
		if !ok || nextURLValue == "" {
			break
		}

		nextURL = nextURLValue
	}

	return results, nil
}
