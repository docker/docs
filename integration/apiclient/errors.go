package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
)

// APIError describes an API Error.
type APIError struct {
	HTTPStatusCode int               `json:"httpStatusCode"`
	Errors         []errors.APIError `json:"errors"`
}

// Error implements the builtin error interface for APIError.
func (apiErr *APIError) Error() string {
	errJSON, err := json.MarshalIndent(apiErr, "", "  ")
	if err != nil {
		return fmt.Sprintf("APIError %d: %v", apiErr.HTTPStatusCode, apiErr.Errors)
	}

	return string(errJSON)
}

// validateStatusCode is used as a helper function by various client methods
// to parse API Errors from a response.
func validateStatusCode(resp *http.Response, expectedStatusCode int) error {
	requestURL := resp.Request.URL.String()
	requestMethod := resp.Request.Method

	if resp.StatusCode == expectedStatusCode {
		return nil
	}

	var apiErr APIError
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error at %s %s: unable to read body: %s", requestMethod, requestURL, err)
	}
	err = json.Unmarshal(body, &apiErr)
	if err != nil {
		return fmt.Errorf("error at %s %s: unable to decode API error %q: %s", requestMethod, requestURL, string(body), err)
	}
	if len(apiErr.Errors) == 0 {
		return fmt.Errorf("error at %s %s: no errors found in json response: %s", requestMethod, requestURL, string(body))
	}
	apiErr.HTTPStatusCode = resp.StatusCode
	return &apiErr
}
