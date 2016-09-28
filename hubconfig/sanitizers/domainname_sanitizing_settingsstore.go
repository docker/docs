package sanitizers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/dhe-deploy/hubconfig"
)

type DomainNameSanitizingSettingsStore struct {
	hubconfig.SettingsStore
}

func (s DomainNameSanitizingSettingsStore) SetUserHubConfig(userHubConfig *hubconfig.UserHubConfig) error {
	host := userHubConfig.DTRHost
	if host != "" {
		if !strings.Contains(host, "://") {
			host = "scheme://" + host
		}
		if url, err := url.Parse(host); err != nil {
			return err
		} else if url.Host == "" {
			return fmt.Errorf("Invalid domain name: %q", host)
		} else {
			parts := strings.Split(url.Host, ":")
			if len(parts) > 1 && parts[1] == "443" {
				host = parts[0]
			} else {
				host = url.Host
			}
		}
	}
	userHubConfig.DTRHost = host
	return s.SettingsStore.SetUserHubConfig(userHubConfig)
}
func (s DomainNameSanitizingSettingsStore) UserHubConfig() (*hubconfig.UserHubConfig, error) {
	userHubConfig, err := s.SettingsStore.UserHubConfig()
	if err != nil {
		return nil, err
	}
	host := userHubConfig.DTRHost
	parts := strings.Split(host, ":")
	if len(parts) > 1 && parts[1] == "443" {
		host = parts[0]
	}
	userHubConfig.DTRHost = host
	return userHubConfig, nil
}
