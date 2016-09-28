package sanitizers

import (
	"github.com/docker/dhe-deploy/hubconfig"
)

func Wrap(settingsStore hubconfig.SettingsStore) hubconfig.SettingsStore {
	return HubCredentialsSanitizingSettingsStore{
		DomainNameSanitizingSettingsStore{
			LogSanitizingSettingsStore{
				AuthSanitizingSettingsStore{
					RegistrySanitizingSettingsStore{
						PortSanitizingSettingsStore{
							settingsStore,
						},
					},
				},
			},
		},
	}
}
