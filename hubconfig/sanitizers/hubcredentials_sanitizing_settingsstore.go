package sanitizers

import (
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"

	"github.com/samalba/dockerclient"
)

type HubCredentialsSanitizingSettingsStore struct {
	hubconfig.SettingsStore
}

func (s HubCredentialsSanitizingSettingsStore) SetHubCredentials(hubCredentials *dockerclient.AuthConfig) error {
	if hubCredentials != nil {
		hubCredentials.Email = deploy.DummyHubUserEmail
	}
	return s.SettingsStore.SetHubCredentials(hubCredentials)
}
