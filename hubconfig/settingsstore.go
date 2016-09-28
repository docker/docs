package hubconfig

import (
	"crypto/x509"
	"time"

	"github.com/docker/libtrust"

	distributionconfig "github.com/docker/distribution/configuration"
	garantconfig "github.com/docker/garant/config"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/samalba/dockerclient"
)

type LicenseChecker interface {
	Initialize() error
	IsValid() bool
	BeginLicenseSyncing()
	IsExpired() bool
	LicensingEnforced() bool
	LicenseTier() string
	LicenseType() string
	GetLicenseID() string
	Expiration() time.Time
	ToggleAutoRefresh(autoRefresh bool) error
	ChangeLicenseFromId(keyID, privateKey string) error
	LoadLicenseFromConfig(config *LicenseConfig, newLicense bool) error
}

type SettingsStore interface {
	SettingsReader
	SettingsWriter
	Ping() error
}

type SettingsReader interface {
	UserHubConfig() (*UserHubConfig, error)
	AuthConfig() (*garantconfig.Configuration, error)
	LicenseConfig() (*LicenseConfig, error)
	RegistryConfig() (*distributionconfig.Configuration, error)
	HubCredentials() (*dockerclient.AuthConfig, error)
	GarantSigningKey() (libtrust.PrivateKey, error)
	GarantRootCert() (string, error)
	EnziSigningKey() (libtrust.PrivateKey, error)
	RawEnziSigningKey() (string, error)
	EnziService() (*enziresponses.Service, error)
	HAConfig() (*HAConfig, error)
}

type SettingsWriter interface {
	SetLicenseStatus(bool) error
	SetUserHubConfig(*UserHubConfig) error
	SetAuthConfig(*garantconfig.Configuration) error
	SetLicenseConfig(*LicenseConfig) error
	SetRegistryConfig(*distributionconfig.Configuration) error
	SetHubCredentials(*dockerclient.AuthConfig) error
	SetGarantRootCert(*x509.Certificate) error
	SetGarantSigningKey(libtrust.PrivateKey) error
	SetEnziSigningKey(libtrust.PrivateKey) error
	SetEnziService(*enziresponses.Service) error
	SetHAConfig(*HAConfig) error
}
