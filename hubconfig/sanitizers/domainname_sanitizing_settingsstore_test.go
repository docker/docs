package sanitizers

import (
	"reflect"
	"testing"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/memory"
)

func TestSanitizedDomainName(t *testing.T) {
	testCases := []struct {
		unsanitized, sanitized hubconfig.UserHubConfig
		expectSetError         bool
	}{
		// Empty config
		{},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com/",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com/some/path",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "http://docker.com",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com/asdf#fragment",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com/asdf?queryarg=stuff&otherarg=otherstuff",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "https://docker.com/some/path",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "docker.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "vagrant.local",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "vagrant.local",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "lots.of.sub-domains.com",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "lots.of.sub-domains.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "https://sub.my-domain.com/some/path#andfragment",
			},
			sanitized: hubconfig.UserHubConfig{
				DTRHost: "sub.my-domain.com",
			},
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "/just/a/path",
			},
			expectSetError: true,
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "http://",
			},
			expectSetError: true,
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "?queryparams=only",
			},
			expectSetError: true,
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "unix:///var/run/docker.sock",
			},
			expectSetError: true,
		},
		{
			unsanitized: hubconfig.UserHubConfig{
				DTRHost: "bad%20scheme://docker.com",
			},
			expectSetError: true,
		},
	}
	for _, testCase := range testCases {
		ss := DomainNameSanitizingSettingsStore{memory.NewSettingsStore()}

		err := ss.SetUserHubConfig(&testCase.unsanitized)
		if testCase.expectSetError {
			if err == nil {
				t.Fatalf("Expected setting of config to fail: %#v", testCase.unsanitized)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			retrieved, err := ss.UserHubConfig()
			if err != nil {
				t.Fatal(err)
			}
			if retrieved == nil {
				t.Fatal("Retrieved nil auth config")
			}

			assertUserHubConfigEquals(t, *retrieved, testCase.sanitized)

		}
	}
}

func assertUserHubConfigEquals(t *testing.T, received, expected hubconfig.UserHubConfig) {
	if !reflect.DeepEqual(received, expected) {
		t.Fatalf("Configs not equal!\nReceived: %#v\nExpected: %#v", received, expected)
	}
}
