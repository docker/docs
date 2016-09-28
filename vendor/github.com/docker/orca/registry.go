package orca

import (
	"github.com/docker/distribution/reference"
	"net/http"
	"net/url"
)

type (
	RegistryClient struct {
		URL        *url.URL
		HttpClient *http.Client
	}

	RegistryConfig struct {
		Type     string `json:"type"`
		ID       string `json:"id,omitempty"`
		Name     string `json:"name"`
		URL      string `json:"url"`
		Insecure bool   `json:"insecure"`
		CACert   string `json:"ca_cert"`
	}

	// This should be replaced by docker/engine-api/types/AuthConfig
	AuthConfig struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Auth     string `json:"auth"`
		Email    string `json:"email"`
		Hostname string `json:"serveraddress"`
	}

	Registry interface {
		GetConfig() *RegistryConfig
		GetAuthToken(username string, accessType string, named reference.Named) (string, error)
	}
)
