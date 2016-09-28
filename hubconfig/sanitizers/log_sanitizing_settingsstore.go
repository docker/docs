package sanitizers

import (
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/distribution/configuration"
)

type LogSanitizingSettingsStore struct {
	hubconfig.SettingsStore
}

func (s LogSanitizingSettingsStore) SetHAConfig(haConfig *hubconfig.HAConfig) error {
	err := haConfig.LoggingValid()
	if err != nil {
		return err
	}
	_, _, err = haConfig.SyslogLogrusLevels()
	if err != nil {
		return err
	}

	return s.SettingsStore.SetHAConfig(haConfig)
}

func (s LogSanitizingSettingsStore) SetRegistryConfig(registryConfig *configuration.Configuration) error {
	if registryConfig.Log.Formatter == "" {
		registryConfig.Log.Formatter = "json"
	}

	return s.SettingsStore.SetRegistryConfig(registryConfig)
}
