package apiclient

import (
	"encoding/json"
	"net/http"
	"net/url"

	distributionconfig "github.com/docker/distribution/configuration"
	"gopkg.in/yaml.v2"
)

const registrySettingsSubroute = "/settings/registry"

func (c *apiClient) GetRegistrySettings() (*distributionconfig.Configuration, error) {
	if resp, err := c.makeRequest("GET", url.URL{Path: registrySettingsSubroute}, nil); err != nil {
		return nil, err
	} else {
		resp.Body.Close()
		if err := validateStatusCode(resp, http.StatusOK); err != nil {
			return nil, err
		}

		var config struct {
			Storage distributionconfig.Configuration
		}
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			return nil, err
		}
		return &config.Storage, nil
	}
}

func (c *apiClient) SetRegistrySettings(config *distributionconfig.Configuration) error {
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	resp, err := c.makeRequest("PUT", url.URL{Path: registrySettingsSubroute}, map[string]interface{}{
		"config": &[]string{string(configBytes)}[0],
	})
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
