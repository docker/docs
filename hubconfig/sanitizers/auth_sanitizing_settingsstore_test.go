package sanitizers

import (
	"reflect"
	"testing"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/memory"

	garantconfig "github.com/docker/garant/config"
)

var authTestCases = []struct {
	domainName             string
	unsanitized, sanitized garantconfig.Configuration
}{
	// Empty unsanitized config
	{
		sanitized: garantconfig.Configuration{
			Version:    defaultconfigs.DefaultAuthConfig.Version,
			Auth:       defaultconfigs.DefaultAuthConfig.Auth,
			HTTP:       defaultconfigs.DefaultAuthConfig.HTTP,
			Logging:    defaultconfigs.DefaultAuthConfig.Logging,
			SigningKey: defaultconfigs.DefaultAuthConfig.SigningKey,
			Issuer:     defaultconfigs.DefaultAuthConfig.Issuer,
		},
	},
	{
		unsanitized: garantconfig.Configuration{
			Auth: garantconfig.Auth{
				BackendName: "notarealbackend",
			},
		},
		sanitized: garantconfig.Configuration{
			Version:    defaultconfigs.DefaultAuthConfig.Version,
			Auth:       defaultconfigs.DefaultAuthConfig.Auth,
			HTTP:       defaultconfigs.DefaultAuthConfig.HTTP,
			Logging:    defaultconfigs.DefaultAuthConfig.Logging,
			SigningKey: defaultconfigs.DefaultAuthConfig.SigningKey,
			Issuer:     defaultconfigs.DefaultAuthConfig.Issuer,
		},
	},
	{
		unsanitized: garantconfig.Configuration{
			Version: "0.45",
		},
		sanitized: garantconfig.Configuration{
			Version:    defaultconfigs.DefaultAuthConfig.Version,
			Auth:       defaultconfigs.DefaultAuthConfig.Auth,
			HTTP:       defaultconfigs.DefaultAuthConfig.HTTP,
			Logging:    defaultconfigs.DefaultAuthConfig.Logging,
			SigningKey: defaultconfigs.DefaultAuthConfig.SigningKey,
			Issuer:     defaultconfigs.DefaultAuthConfig.Issuer,
		},
	},
	{
		unsanitized: garantconfig.Configuration{
			Version: defaultconfigs.DefaultAuthConfig.Version,
		},
		sanitized: garantconfig.Configuration{
			Version:    defaultconfigs.DefaultAuthConfig.Version,
			Auth:       defaultconfigs.DefaultAuthConfig.Auth,
			HTTP:       defaultconfigs.DefaultAuthConfig.HTTP,
			Logging:    defaultconfigs.DefaultAuthConfig.Logging,
			SigningKey: defaultconfigs.DefaultAuthConfig.SigningKey,
			Issuer:     defaultconfigs.DefaultAuthConfig.Issuer,
		},
	},
}

func TestSanitizeAuthConfig(t *testing.T) {
	for _, testCase := range authTestCases {
		ss := AuthSanitizingSettingsStore{memory.NewSettingsStore()}

		err := ss.SetUserHubConfig(&hubconfig.UserHubConfig{
			DTRHost: testCase.domainName,
		})
		if err != nil {
			t.Fatal(err)
		}

		err = ss.SetAuthConfig(&testCase.unsanitized)
		if err != nil {
			t.Fatal(err)
		}
		retrieved, err := ss.AuthConfig()
		if err != nil {
			t.Fatal(err)
		}
		if retrieved == nil {
			t.Fatal("Retrieved nil auth config")
		}

		assertAuthConfigEquals(t, *retrieved, testCase.sanitized)
	}
}

func TestHubConfigSanitizesAuthConfig(t *testing.T) {
	for _, testCase := range authTestCases {
		if testCase.domainName == "" {
			continue
		}
		ss := AuthSanitizingSettingsStore{memory.NewSettingsStore()}

		err := ss.SetUserHubConfig(&hubconfig.UserHubConfig{
			DTRHost: "some-domain.com",
		})
		if err != nil {
			t.Fatal(err)
		}

		err = ss.SetAuthConfig(&testCase.unsanitized)
		if err != nil {
			t.Fatal(err)
		}

		err = ss.SetUserHubConfig(&hubconfig.UserHubConfig{
			DTRHost: testCase.domainName,
		})
		if err != nil {
			t.Fatal(err)
		}

		retrieved, err := ss.AuthConfig()
		if err != nil {
			t.Fatal(err)
		}
		if retrieved == nil {
			t.Fatal("Retrieved nil auth config")
		}

		assertAuthConfigEquals(t, *retrieved, testCase.sanitized)
	}
}

func assertAuthConfigEquals(t *testing.T, received, expected garantconfig.Configuration) {
	if expected.Logging.Fields == nil {
		expected.Logging.Fields = map[string]string{}
	}
	if received.Logging.Fields == nil {
		received.Logging.Fields = map[string]string{}
	}

	if received.Version != expected.Version {
		t.Fatalf("Versions not equal!\nReceived: %#v\nExpected: %#v", received.Version, expected.Version)
	}

	if !reflect.DeepEqual(received.Logging, expected.Logging) {
		t.Fatalf("Logging not equal!\nReceived: %#v\nExpected: %#v", received.Logging, expected.Logging)
	}

	if received.SigningKey != expected.SigningKey {
		t.Fatalf("SigningKeys not equal!\nReceived: %#v\nExpected: %#v", received.SigningKey, expected.SigningKey)
	}

	if received.Issuer != expected.Issuer {
		t.Fatalf("Issuers not equal!\nReceived: %#v\nExpected: %#v", received.Issuer, expected.Issuer)
	}

	if !reflect.DeepEqual(received.Reporting, expected.Reporting) {
		t.Fatalf("Reporting not equal!\nReceived: %#v\nExpected: %#v", received.Reporting, expected.Reporting)
	}

	if !reflect.DeepEqual(received.HTTP, expected.HTTP) {
		t.Fatalf("HTTP not equal!\nReceived: %#v\nExpected: %#v", received.HTTP, expected.HTTP)
	}

	if !reflect.DeepEqual(received.Auth, expected.Auth) {
		t.Fatalf("Auth not equal!\nReceived: %#v\nExpected: %#v", received.Auth, expected.Auth)
	}
}
