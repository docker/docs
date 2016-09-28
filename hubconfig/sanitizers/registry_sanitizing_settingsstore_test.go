package sanitizers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/memory"
	"github.com/docker/dhe-deploy/hubconfig/util"

	"github.com/docker/distribution/configuration"
)

type testCaseType struct {
	domainName             string
	loadBalancerHTTPSPort  uint16
	unsanitized, sanitized configuration.Configuration
}

var registryTestCases = []testCaseType{
	{
		// Empty unsanitized config
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Loglevel = "warn"
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Log.Level = "warn"
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			// Log.Level has higher precendence than the deprecated Loglevel
			config.Loglevel = "warn"
			config.Log.Level = "debug"
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Log.Level = "debug"
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = fmt.Sprintf("localhost:%d", deploy.AdminPort)
			config.HTTP.Secret = "placeholder"
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Storage = configuration.Storage{"inmemory": nil}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = configuration.Storage{}
			for k, v := range defaultconfigs.DefaultRegistryConfig.Storage {
				config.Storage[k] = v
			}
			delete(config.Storage, "filesystem")
			config.Storage["inmemory"] = nil
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Storage = configuration.Storage{
				"inmemory": nil,
				"delete": configuration.Parameters{
					"enabled": false,
					"herp":    "derp",
				},
				"maintenance": configuration.Parameters{
					"bogus": "stuff",
				},
			}
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = configuration.Storage{}
			for k, v := range defaultconfigs.DefaultRegistryConfig.Storage {
				config.Storage[k] = v
			}
			delete(config.Storage, "filesystem")
			config.Storage["inmemory"] = nil
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Storage = configuration.Storage{
				"inmemory": nil,
				"maintenance": configuration.Parameters{
					"readonly": map[string]interface{}{
						"enabled": true,
					},
				},
			}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = configuration.Storage{}
			for k, v := range defaultconfigs.DefaultRegistryConfig.Storage {
				config.Storage[k] = v
			}
			delete(config.Storage, "filesystem")
			config.Storage["inmemory"] = nil
			config.Storage["maintenance"] = configuration.Parameters{
				"readonly": map[string]interface{}{
					"enabled": true,
				},
			}
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Storage = configuration.Storage{"filesystem": nil}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Auth = configuration.Auth{
				"silly": nil,
			}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Notifications.Endpoints = []configuration.Endpoint{
				{
					Name:    "some-endpoint",
					URL:     "https://some-url/asdf",
					Headers: http.Header{"X-Header-Name": []string{"some value"}},
				},
			}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications.Endpoints = append(defaultconfigs.DefaultRegistryConfig.Notifications.Endpoints,
				configuration.Endpoint{
					Name:    "some-endpoint",
					URL:     "https://some-url/asdf",
					Headers: http.Header{"X-Header-Name": []string{"some value"}},
				},
			)
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName: "docker.com",
		// Empty unsanitized config
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://docker.com/auth/token",
					"issuer":         "docker.com",
					"service":        "docker.com",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName: "my.nested.domain.org",
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = fmt.Sprintf("localhost:%d", deploy.AdminPort)
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://my.nested.domain.org/auth/token",
					"issuer":         "my.nested.domain.org",
					"service":        "my.nested.domain.org",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName: "vagrant.local",
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Auth = configuration.Auth{
				"silly": nil,
			}
			config.Notifications.Endpoints = []configuration.Endpoint{
				{
					Name:    "some-endpoint",
					URL:     "https://some-url/asdf",
					Headers: http.Header{"X-Header-Name": []string{"some value"}},
				},
			}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://vagrant.local/auth/token",
					"issuer":         "vagrant.local",
					"service":        "vagrant.local",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications.Endpoints = append(defaultconfigs.DefaultRegistryConfig.Notifications.Endpoints,
				configuration.Endpoint{
					Name:    "some-endpoint",
					URL:     "https://some-url/asdf",
					Headers: http.Header{"X-Header-Name": []string{"some value"}},
				},
			)
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName: "vagrant.local",
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Auth = configuration.Auth{
				"silly": nil,
			}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://vagrant.local/auth/token",
					"issuer":         "vagrant.local",
					"service":        "vagrant.local",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName: "vagrant.local",
		unsanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Auth = configuration.Auth{
				"silly": nil,
			}
			return config
		}(),
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://vagrant.local/auth/token",
					"issuer":         "vagrant.local",
					"service":        "vagrant.local",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName:            "docker.com",
		loadBalancerHTTPSPort: deploy.AdminTlsPort,
		// Empty unsanitized config
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://docker.com/auth/token",
					"issuer":         "docker.com",
					"service":        "docker.com",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
	{
		domainName:            "docker.com",
		loadBalancerHTTPSPort: 9999,
		// Empty unsanitized config
		sanitized: func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.HTTP.Secret = "placeholder"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "https://docker.com:9999/auth/token",
					"issuer":         "docker.com:9999",
					"service":        "docker.com:9999",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			config.Notifications = defaultconfigs.DefaultRegistryConfig.Notifications
			config.Compatibility.Schema1.DisableSignatureStore = true
			config.Middleware = defaultconfigs.DefaultMiddlewareConfig
			return config
		}(),
	},
}

func TestSanitizedRegistryConfig(t *testing.T) {
	for _, testCase := range registryTestCases {
		ss := RegistrySanitizingSettingsStore{memory.NewSettingsStore()}

		err := ss.SetUserHubConfig(&hubconfig.UserHubConfig{
			DTRHost: testCase.formatDomain(),
		})
		if err != nil {
			t.Fatal(err)
		}

		err = ss.SetRegistryConfig(&testCase.unsanitized)
		if err != nil {
			t.Fatal(err)
		}
		retrieved, err := ss.RegistryConfig()
		if err != nil {
			t.Fatal(err)
		}
		if retrieved == nil {
			t.Fatal("Retrieved nil auth config")
		}

		assertRegistryConfigEquals(t, *retrieved, testCase.sanitized)
	}
}

func (testCase testCaseType) formatDomain() string {
	if testCase.domainName != "" {
		if testCase.loadBalancerHTTPSPort != 0 {
			return fmt.Sprintf("%s:%d", testCase.domainName, testCase.loadBalancerHTTPSPort)
		} else {
			return fmt.Sprintf("%s", testCase.domainName)
		}
	}
	return ""
}

func TestHubConfigSanitizesRegistryConfig(t *testing.T) {
	for _, testCase := range registryTestCases {
		ss := RegistrySanitizingSettingsStore{memory.NewSettingsStore()}

		err := ss.SetUserHubConfig(&hubconfig.UserHubConfig{
			DTRHost: fmt.Sprintf("%s:%d", "some-domain.com", 123),
		})
		if err != nil {
			t.Fatal(err)
		}

		err = ss.SetRegistryConfig(&testCase.unsanitized)
		if err != nil {
			t.Fatal(err)
		}

		err = ss.SetUserHubConfig(&hubconfig.UserHubConfig{
			DTRHost: testCase.formatDomain(),
		})
		if err != nil {
			t.Fatal(err)
		}

		retrieved, err := ss.RegistryConfig()
		if err != nil {
			t.Fatal(err)
		}
		if retrieved == nil {
			t.Fatal("Retrieved nil auth config")
		}

		assertRegistryConfigEquals(t, *retrieved, testCase.sanitized)
	}
}

func assertRegistryConfigEquals(t *testing.T, received, expected configuration.Configuration) {
	util.SetReadonlyModeJSON(&received.Storage, util.GetReadonlyMode(&received.Storage))
	util.SetReadonlyModeJSON(&expected.Storage, util.GetReadonlyMode(&expected.Storage))
	if !reflect.DeepEqual(received, expected) {
		t.Errorf("Configs not equal!\n=> Received: %s\n=> Expected: %s", spew.Sdump(received), spew.Sdump(expected))
	}
}
