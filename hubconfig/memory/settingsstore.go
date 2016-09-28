package memory

import (
	"crypto/x509"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/libtrust"

	distributionconfig "github.com/docker/distribution/configuration"
	garantconfig "github.com/docker/garant/config"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/samalba/dockerclient"
)

type settingsStore struct {
	storage map[string]interface{}
}

const (
	userHubConfigKey       = "hubConfig"
	authConfigKey          = "authConfig"
	registryIndexConfigKey = "registryIndexConfig"
	licenseConfigKey       = "licenseConfig"
	haConfigKey            = "haConfig"
	registryKey            = "registry"
	hubCredentialsKey      = "hubCredentials"
	garantSigningKey       = "garantSigningKey"
	enziSigningKey         = "enziSigningKey"
	enziService            = "enziService"
	notaryServer           = "notaryServer"
	notarySigner           = "notarySigner"
)

func NewSettingsStore() hubconfig.SettingsStore {
	return settingsStore{storage: make(map[string]interface{})}
}

func (s settingsStore) UserHubConfig() (*hubconfig.UserHubConfig, error) {
	userHubConfig, ok := s.storage[userHubConfigKey]
	if !ok {
		return nil, nil
	}
	return userHubConfig.(*hubconfig.UserHubConfig), nil
}

func (s settingsStore) SetUserHubConfig(userHubConfig *hubconfig.UserHubConfig) error {
	s.storage[userHubConfigKey] = userHubConfig
	return nil
}

func (s settingsStore) AuthConfig() (*garantconfig.Configuration, error) {
	authConfig, ok := s.storage[authConfigKey]
	if !ok {
		return nil, nil
	}
	return authConfig.(*garantconfig.Configuration), nil
}

func (s settingsStore) SetAuthConfig(authConfig *garantconfig.Configuration) error {
	s.storage[authConfigKey] = authConfig
	return nil
}

func (s settingsStore) LicenseConfig() (*hubconfig.LicenseConfig, error) {
	licenseConfig, ok := s.storage[licenseConfigKey]
	if !ok {
		return nil, nil
	}
	return licenseConfig.(*hubconfig.LicenseConfig), nil
}

func (s settingsStore) SetLicenseConfig(licenseConfig *hubconfig.LicenseConfig) error {
	s.storage[licenseConfigKey] = licenseConfig
	return nil
}

func (s settingsStore) SetHAConfig(haConfig *hubconfig.HAConfig) error {
	s.storage[haConfigKey] = haConfig
	return nil
}

func (s settingsStore) HAConfig() (*hubconfig.HAConfig, error) {
	haConfig, ok := s.storage[haConfigKey]
	if !ok {
		return nil, nil
	}
	return haConfig.(*hubconfig.HAConfig), nil
}

func (s settingsStore) SetLicenseStatus(isValid bool) error {
	return nil
}

func (s settingsStore) RegistryConfig() (*distributionconfig.Configuration, error) {
	registryConfig, ok := s.storage[registryKey]
	if !ok {
		return nil, nil
	}
	return registryConfig.(*distributionconfig.Configuration), nil
}

func (s settingsStore) SetRegistryConfig(registryConfig *distributionconfig.Configuration) error {
	s.storage[registryKey] = registryConfig
	return nil
}

func (s settingsStore) HubCredentials() (*dockerclient.AuthConfig, error) {
	hubCredentials, ok := s.storage[hubCredentialsKey]
	if !ok {
		return nil, nil
	}
	return hubCredentials.(*dockerclient.AuthConfig), nil
}

func (s settingsStore) SetHubCredentials(hubCredentials *dockerclient.AuthConfig) error {
	s.storage[hubCredentialsKey] = hubCredentials
	return nil
}

func (s settingsStore) GarantRootCert() (string, error) {
	return "", nil
}

func (s settingsStore) SetGarantRootCert(*x509.Certificate) error {
	return nil
}

func (s settingsStore) GarantSigningKey() (libtrust.PrivateKey, error) {
	key, ok := s.storage[garantSigningKey]
	if !ok {
		return nil, nil
	}
	return key.(libtrust.PrivateKey), nil
}

func (s settingsStore) SetGarantSigningKey(pk libtrust.PrivateKey) error {
	s.storage[garantSigningKey] = pk
	return nil
}

func (s settingsStore) RawEnziSigningKey() (string, error) {
	return "", nil
}

func (s settingsStore) EnziSigningKey() (libtrust.PrivateKey, error) {
	key, ok := s.storage[enziSigningKey]
	if !ok {
		return nil, nil
	}
	return key.(libtrust.PrivateKey), nil
}

func (s settingsStore) SetEnziSigningKey(pk libtrust.PrivateKey) error {
	s.storage[enziSigningKey] = pk
	return nil
}

func (s settingsStore) EnziService() (*enziresponses.Service, error) {
	service, ok := s.storage[enziService]
	if !ok {
		return nil, nil
	}
	return service.(*enziresponses.Service), nil
}

func (s settingsStore) SetEnziService(service *enziresponses.Service) error {
	s.storage[enziService] = service
	return nil
}

func (s settingsStore) Ping() error {
	return nil
}
