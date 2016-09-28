package sanitizers

import (
	"reflect"
	"testing"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig/memory"

	"github.com/samalba/dockerclient"
)

func TestSanitizedHubCredentials(t *testing.T) {
	testCases := []struct {
		unsanitized, sanitized dockerclient.AuthConfig
	}{
		// Empty unsanitized config
		{
			sanitized: dockerclient.AuthConfig{
				Email: deploy.DummyHubUserEmail,
			},
		},
		{
			unsanitized: dockerclient.AuthConfig{
				Email: "someguy@domain.com",
			},
			sanitized: dockerclient.AuthConfig{
				Email: deploy.DummyHubUserEmail,
			},
		},
		{
			unsanitized: dockerclient.AuthConfig{
				Email: "brianbland@docker.com",
			},
			sanitized: dockerclient.AuthConfig{
				Email: deploy.DummyHubUserEmail,
			},
		},
		{
			unsanitized: dockerclient.AuthConfig{
				Username: "someguy",
				Password: "somepassword",
				Email:    "someguy@domain.com",
			},
			sanitized: dockerclient.AuthConfig{
				Username: "someguy",
				Password: "somepassword",
				Email:    deploy.DummyHubUserEmail,
			},
		},
		{
			unsanitized: dockerclient.AuthConfig{
				Username: "someguy",
				Password: "P4$$w0rd!",
				Email:    "herp derp",
			},
			sanitized: dockerclient.AuthConfig{
				Username: "someguy",
				Password: "P4$$w0rd!",
				Email:    deploy.DummyHubUserEmail,
			},
		},
	}
	for _, testCase := range testCases {
		ss := HubCredentialsSanitizingSettingsStore{memory.NewSettingsStore()}

		err := ss.SetHubCredentials(&testCase.unsanitized)
		if err != nil {
			t.Fatal(err)
		}
		retrieved, err := ss.HubCredentials()
		if err != nil {
			t.Fatal(err)
		}
		if retrieved == nil {
			t.Fatal("Retrieved nil auth config")
		}

		assertHubCredentialsEquals(t, *retrieved, testCase.sanitized)
	}
}

func assertHubCredentialsEquals(t *testing.T, received, expected dockerclient.AuthConfig) {
	if !reflect.DeepEqual(received, expected) {
		t.Fatalf("Configs not equal!\nReceived: %#v\nExpected: %#v", received, expected)
	}
}
