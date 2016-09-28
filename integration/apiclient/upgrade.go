package apiclient

import (
	"net/http"
	"net/url"
)

const upgradeSubroute = "/api/v0/admin/upgrade"

func (c *apiClient) Upgrade(toVersion string) error {
	resp, err := c.makeRequest("POST", url.URL{Path: upgradeSubroute}, struct {
		Version string
	}{
		Version: toVersion,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusAccepted); err != nil {
		return err
	}

	return nil
}
