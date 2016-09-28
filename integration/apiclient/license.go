package apiclient

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/docker/dhe-deploy/adminserver"
	"github.com/docker/dhe-deploy/hubconfig"
)

const licenseSubroute = "/api/v0/admin/settings/license"

func (c *apiClient) GetLicenseSettings() (*adminserver.LicenseSettings, error) {
	resp, err := c.GetLicenseSettingsResponse()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var licenseSettings adminserver.LicenseSettings
	if err := json.NewDecoder(resp.Body).Decode(&licenseSettings); err != nil {
		return nil, err
	}

	return &licenseSettings, err
}

func (c *apiClient) GetLicenseSettingsResponse() (*http.Response, error) {
	return c.makeRequest("GET", url.URL{Path: licenseSubroute}, nil)
}

func (c *apiClient) SetLicenseSettings(licenseConfig *hubconfig.LicenseConfig) (*adminserver.LicenseSettings, error) {
	resp, err := c.SetLicenceSetttingsResponse(licenseConfig)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusAccepted); err != nil {
		return nil, err
	}
	var licenseSettings adminserver.LicenseSettings
	if err := json.NewDecoder(resp.Body).Decode(&licenseSettings); err != nil {
		return nil, err
	}

	return &licenseSettings, err
}

func (c *apiClient) SetLicenceSetttingsResponse(licenseConfig *hubconfig.LicenseConfig) (*http.Response, error) {
	return c.makeRequest("PUT", url.URL{Path: licenseSubroute}, licenseConfig)
}
