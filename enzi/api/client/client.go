package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
)

// RequestAuthenticator is any type which can set authentication credentials
// on an HTTP request.
type RequestAuthenticator interface {
	AuthenticateRequest(r *http.Request)
}

// BasicAuthenticator implements RequestAuthenticator using a username and
// password.
type BasicAuthenticator struct {
	Username, Password string
}

// AuthenticateRequest sets a basic authorization header on the given request.
func (a *BasicAuthenticator) AuthenticateRequest(r *http.Request) {
	r.SetBasicAuth(a.Username, a.Password)
}

var (
	_ RequestAuthenticator = (*BasicAuthenticator)(nil)
	_ RequestAuthenticator = (*responses.LoginSession)(nil)
)

// Session is used to access eNZi API endpoints.
type Session struct {
	httpClient    *http.Client
	serverAddr    string
	pathPrefix    string
	authenticator RequestAuthenticator
}

// New returns a new eNZi API Session making requests to the given serverAddr,
// authenticating these requests using the given authenticator (which may be
// nil to use no authentication) and performs those requests using the given
// httpClient.
func New(httpClient *http.Client, serverAddr, pathPrefix string, authenticator RequestAuthenticator) *Session {
	return &Session{
		httpClient:    httpClient,
		serverAddr:    serverAddr,
		pathPrefix:    pathPrefix,
		authenticator: authenticator,
	}
}

func (s *Session) buildURL(routePath string, queryParams url.Values) *url.URL {
	endpoint := &url.URL{
		Scheme: "https",
		Host:   s.serverAddr,
		Path:   path.Join("/", s.pathPrefix, routePath),
	}

	if len(queryParams) > 0 {
		endpoint.RawQuery = queryParams.Encode()
	}

	return endpoint
}

func (s *Session) performRequestRawResponse(method string, endpoint *url.URL, payload interface{}) (*http.Response, error) {
	var jsonPayload io.Reader
	if payload != nil {
		payloadJSONBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("unable to encode payload into JSON: %s", err)
		}

		jsonPayload = bytes.NewReader(payloadJSONBytes)
	}

	req, err := http.NewRequest(method, endpoint.String(), jsonPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare request: %s", err)
	}

	if jsonPayload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if s.authenticator != nil {
		s.authenticator.AuthenticateRequest(req)
	}

	return s.httpClient.Do(req)
}

func (s *Session) performRequest(method string, endpoint *url.URL, payload interface{}, expectedStatus int, responseVal interface{}, nextPageStart *string) error {
	resp, err := s.performRequestRawResponse(method, endpoint, payload)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return s.handleUnexpectedResponse(expectedStatus, resp)
	}

	if nextPageStart != nil {
		*nextPageStart = resp.Header.Get("X-Next-Page-Start")
	}

	if responseVal == nil {
		// No response body expected.
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(responseVal); err != nil {
		return fmt.Errorf("unable to decode %T JSON response: %s", responseVal, err)
	}

	return nil
}

// handleUnexpectedResponse attempts to decode the unexpected response body
// into an API Error response envelope. If this is successful, the returned
// response will have the type *(github.com/docker/orca/enzi/api/errors).APIErrors
// otherwise it will be an error string describing the unexpected response.
func (s *Session) handleUnexpectedResponse(expectedStatus int, resp *http.Response) error {
	// Read the full response.
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read unexpected response: expected %d, got %d", expectedStatus, resp.StatusCode)
	}

	// Attempt to parse the response as APIErrors.
	var errResponse errors.APIErrors
	if err := json.Unmarshal(respData, &errResponse); err != nil {
		return fmt.Errorf("unable to decode unexpected response: expected %d, got %d - %s", expectedStatus, resp.StatusCode, string(respData))
	}

	// This field is not set by the JSON decoder.
	errResponse.HTTPStatusCode = resp.StatusCode

	return &errResponse
}
