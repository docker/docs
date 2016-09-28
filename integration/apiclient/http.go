package apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"

	"github.com/gorilla/websocket"
)

const httpSubroute = "/api/v0/meta/settings"

func (c *apiClient) GetCA() ([]byte, error) {
	resp, err := c.makeRequest("GET", url.URL{Path: "/ca"}, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bts, err
}

func (c *apiClient) GetClusterStatus() (*responses.ClusterStatus, error) {
	resp, err := c.makeRequest("GET", url.URL{Path: "/api/v0/meta/cluster_status"}, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var clusterStatus responses.ClusterStatus
	if err := json.NewDecoder(resp.Body).Decode(&clusterStatus); err != nil {
		return nil, err
	}

	return &clusterStatus, err
}

func (c *apiClient) GetEvents() (*responses.Events, *url.Values, error) {
	return c.GetEventsWithParams(url.Values{})
}

func (c *apiClient) GetEventsWithParams(params url.Values) (*responses.Events, *url.Values, error) {
	var nextPageQuery *url.Values
	resp, err := c.makeRequest("GET", url.URL{Path: "/api/v0/events", RawQuery: params.Encode()}, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	pagHeader := resp.Header.Get("Link")
	if pagHeader != "" {
		re := regexp.MustCompile(`<(.*)\?(.*)>;\s?rel="next"`)
		match := re.FindStringSubmatch(pagHeader)
		if len(match) > 1 {
			nextPageObj, err := url.ParseQuery(match[2])
			if err != nil {
				return nil, nil, err
			}
			nextPageQuery = &nextPageObj
		}
	}

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, nil, err
	}

	var events responses.Events
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, nil, err
	}

	return &events, nextPageQuery, nil
}

func (c *apiClient) EventsWebsocket() (*websocket.Conn, error) {
	return c.openUnauthenticatedWebsocket("/ws/events")
}

func (c *apiClient) EventsWebsocketAuthd() (*websocket.Conn, error) {
	return c.openAuthenticatedWebsocket("/ws/events")
}

func (c *apiClient) GetHTTPSettings() (*responses.Settings, error) {
	resp, err := c.GetHTTPSettingsResponse()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var httpSettings responses.Settings
	if err := json.NewDecoder(resp.Body).Decode(&httpSettings); err != nil {
		return nil, err
	}

	return &httpSettings, err
}

func (c *apiClient) GetHTTPSettingsResponse() (*http.Response, error) {
	return c.makeRequest("GET", url.URL{Path: httpSubroute}, nil)
}

func (c *apiClient) SetHTTPSettings(settings *forms.Settings) error {
	resp, err := c.SetHTTPSettingsResponse(settings)
	if err != nil {
		return err
	}

	if err := validateStatusCode(resp, http.StatusAccepted); err != nil {
		return err
	}

	return nil
}

func (c *apiClient) SetHTTPSettingsResponse(settings *forms.Settings) (*http.Response, error) {
	return c.makeRequest("POST", url.URL{Path: httpSubroute}, settings)
}
