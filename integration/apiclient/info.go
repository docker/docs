package apiclient

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const infoSubroute = "/api/v0/admin/info"

func (c *apiClient) Version() (string, error) {
	infoResp, err := c.getInfoResponse()
	if err != nil {
		return "", err
	}

	defer infoResp.Body.Close()

	versionContainer := struct {
		DTR struct {
			Version string
		}
	}{}
	if json.NewDecoder(infoResp.Body).Decode(&versionContainer); err != nil {
		return "", err
	}
	return versionContainer.DTR.Version, nil
}

func (c *apiClient) getInfoResponse() (*http.Response, error) {
	return c.makeRequest("GET", url.URL{Path: infoSubroute}, nil)
}
