package sanitizers

import (
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"

	garantconfig "github.com/docker/garant/config"
)

var permissiveAccessSet = map[string][]string{"*": {"pull", "push"}}

type AuthSanitizingSettingsStore struct {
	hubconfig.SettingsStore
}

func (s AuthSanitizingSettingsStore) setSanitizedAuthConfig(userHubConfig *hubconfig.UserHubConfig, authConfig *garantconfig.Configuration) error {
	var (
		err        error
		domainName string
	)
	if userHubConfig == nil {
		userHubConfig, err = s.UserHubConfig()
		if err != nil {
			return err
		}
	}
	if userHubConfig != nil {
		domainName = userHubConfig.DTRHost
	}
	if authConfig == nil {
		authConfig, err = s.AuthConfig()
		if err != nil {
			return err
		} else if authConfig == nil {
			authConfig = &defaultconfigs.DefaultAuthConfig
		}
	}

	authConfig.Auth = defaultconfigs.DefaultAuthConfig.Auth
	authConfig.HTTP = defaultconfigs.DefaultAuthConfig.HTTP
	authConfig.Logging = defaultconfigs.DefaultAuthConfig.Logging
	authConfig.Version = defaultconfigs.DefaultAuthConfig.Version

	authConfig.SigningKey = defaultconfigs.DefaultAuthConfig.SigningKey
	if domainName != "" {
		authConfig.Issuer = domainName
	} else {
		authConfig.Issuer = defaultconfigs.DefaultAuthConfig.Issuer
	}

	return s.SettingsStore.SetAuthConfig(authConfig)
}

func (s AuthSanitizingSettingsStore) SetUserHubConfig(userHubConfig *hubconfig.UserHubConfig) error {
	if err := s.setSanitizedAuthConfig(userHubConfig, nil); err != nil {
		return err
	}
	// TODO it's awkward if SetUserHubConfig fails...
	return s.SettingsStore.SetUserHubConfig(userHubConfig)
}

func (s AuthSanitizingSettingsStore) SetAuthConfig(authConfig *garantconfig.Configuration) error {
	return s.setSanitizedAuthConfig(nil, authConfig)
}
