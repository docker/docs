package sanitizers

import (
	"fmt"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/util"

	"github.com/docker/distribution/configuration"
)

type RegistrySanitizingSettingsStore struct {
	hubconfig.SettingsStore
}

// XXX: wtf is going on in this function?
func (s RegistrySanitizingSettingsStore) setSanitizedRegistryConfig(userHubConfig *hubconfig.UserHubConfig, registryConfig *configuration.Configuration) error {
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
	if registryConfig == nil {
		registryConfig, err = s.RegistryConfig()
		if err != nil {
			return err
		} else if registryConfig == nil {
			registryConfig = &defaultconfigs.DefaultRegistryConfig
		}
	}

	emptyConfig := configuration.Configuration{}
	registryConfig.HTTP = emptyConfig.HTTP
	registryConfig.HTTP.Addr = fmt.Sprintf(":%d", deploy.StorageContainerPort)
	// TODO(andrewnguyen): replace this in the future
	registryConfig.HTTP.Secret = "placeholder"
	registryConfig.Auth = emptyConfig.Auth

	if registryConfig.Storage.Type() == "" {
		registryConfig.Storage = defaultconfigs.DefaultRegistryConfig.Storage
	}

	if registryConfig.Storage.Type() == defaultconfigs.DefaultRegistryConfig.Storage.Type() {
		if rootDir, _ := registryConfig.Storage.Parameters()["rootdirectory"].(string); rootDir == "" {
			if registryConfig.Storage.Parameters() == nil {
				registryConfig.Storage = map[string]configuration.Parameters{registryConfig.Storage.Type(): make(configuration.Parameters)}
			}
			registryConfig.Storage.Parameters()["rootdirectory"] = defaultconfigs.DefaultRegistryConfig.Storage.Parameters()["rootdirectory"]
		}
	}

	registryConfig.Storage["delete"] = configuration.Parameters{"enabled": true}

	// Clear the maintenance section of the storage except for the current readonly mode
	storageReadOnly := util.GetReadonlyMode(&registryConfig.Storage)
	delete(registryConfig.Storage, "maintenance")
	util.SetReadonlyMode(&registryConfig.Storage, storageReadOnly)

	if domainName != "" {
		config := util.GetRegistryAuthConfig(domainName)
		registryConfig.Auth = configuration.Auth{
			"token": configuration.Parameters{
				"realm":          config.Realm,
				"issuer":         config.Issuer,
				"service":        config.Service,
				"rootcertbundle": config.CertBundle,
			},
		}
	}

	if registryConfig.Loglevel != "" {
		if registryConfig.Log.Level == "" {
			registryConfig.Log.Level = registryConfig.Loglevel
		}
		registryConfig.Loglevel = ""
	}

	numDefaultEndpoints := len(defaultconfigs.DefaultRegistryConfig.Notifications.Endpoints)
	endpoints := make([]configuration.Endpoint, numDefaultEndpoints)
	for i := 0; i < numDefaultEndpoints; i++ {
		endpoints[i] = defaultconfigs.DefaultRegistryConfig.Notifications.Endpoints[i]
	}

	for _, endpoint := range registryConfig.Notifications.Endpoints {
		isDefaultEndpoint := false
		for _, defaultEndpoint := range defaultconfigs.DefaultRegistryConfig.Notifications.Endpoints {
			if endpoint.Name == defaultEndpoint.Name || endpoint.URL == defaultEndpoint.URL {
				isDefaultEndpoint = true
				break
			}
		}
		if !isDefaultEndpoint {
			endpoints = append(endpoints, endpoint)
		}
	}
	registryConfig.Notifications.Endpoints = endpoints

	// Ensure we disable saving V1 signatures within the registry. Each
	// registry container will create their own signing key and attach
	// signatures on the fly. They'll be different depending on the pod
	// you communicate with, but they're never verified and no-one cares
	// about them what-so-fucking-ever so it's a good riddance.
	registryConfig.Compatibility.Schema1.DisableSignatureStore = true

	// Also set up our metadata middleware if this isn't set as a pull-through
	// cache (pull-through cache registries should store manifests and tags
	// locally)
	if registryConfig.Proxy.RemoteURL == "" {
		registryConfig.Middleware = defaultconfigs.DefaultMiddlewareConfig
	}

	// Ensure we disable saving V1 signatures within the registry. Each
	// registry container will create their own signing key and attach
	// signatures on the fly. They'll be different depending on the pod
	// you communicate with, but they're never verified and no-one cares
	// about them what-so-fucking-ever so it's a good riddance.
	registryConfig.Compatibility.Schema1.DisableSignatureStore = true

	return s.SettingsStore.SetRegistryConfig(registryConfig)
}

func (s RegistrySanitizingSettingsStore) SetUserHubConfig(userHubConfig *hubconfig.UserHubConfig) error {
	if err := s.setSanitizedRegistryConfig(userHubConfig, nil); err != nil {
		return err
	}
	// TODO it's awkward if SetUserHubConfig fails...
	return s.SettingsStore.SetUserHubConfig(userHubConfig)
}

func (s RegistrySanitizingSettingsStore) SetRegistryConfig(config *configuration.Configuration) error {
	return s.setSanitizedRegistryConfig(nil, config)
}
