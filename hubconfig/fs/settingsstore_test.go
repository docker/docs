// TODO test fs KeyValueStore and SettingsStore independently of each other!
// this is dumb.
package fs

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/shared/containers"

	"github.com/docker/distribution/configuration"
	garantconfig "github.com/docker/garant/config"
	"github.com/samalba/dockerclient"
)

// WARNING: through this test we are creating a fake settings store that assumes that
// we are running an unlicensed version of DTR. This doesn't affect this test right
// now, but it might in the future if we change the effects of not being licensed

func TestUserHubConfig(t *testing.T) {
	hubConfig, err := util.DefaultUserHubConfig("")
	if err != nil {
		t.Fatal(err)
	}
	hubConfig1 := *hubConfig
	hubConfig2 := *hubConfig
	hubConfig2.DTRHost = "docker.com"
	testCases := []hubconfig.UserHubConfig{
		hubConfig1,
		hubConfig2,
	}
	for _, testCase := range testCases {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		ss := settingsstore.New(NewKeyValueStore(tempDir))
		if testCase.DTRHost != "" {
			cert, err := util.GenTLSCert(testCase.DTRHost)
			if err != nil {
				t.Fatal(err)
			}

			err = util.SetTLSCertificateInHubConfig(&testCase, cert, cert)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = ss.SetUserHubConfig(&testCase)
		if err != nil {
			t.Fatal(err)
		}
		received, err := ss.UserHubConfig()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(*received, testCase) {
			t.Fatalf("Configs not equal!\nReceived: %#v\nExpected: %#v", *received, testCase)
		}
	}
}

func TestAuthConfig(t *testing.T) {
	testCases := []garantconfig.Configuration{
		{
			Auth: garantconfig.Auth{
				BackendName: "dtr",
			},
		},
		{
			Version: "0.1",
			Logging: garantconfig.Logging{
				Level: "debug",
			},
			SigningKey: "path/to/a/signing/key.json",
			Issuer:     "some.guy",
			Auth: garantconfig.Auth{
				BackendName: "dtr",
			},
			HTTP: garantconfig.HTTP{
				Addr: ":80",
			},
		},
	}
	for _, testCase := range testCases {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		ss := settingsstore.New(NewKeyValueStore(tempDir))
		err = ss.SetAuthConfig(&testCase)
		if err != nil {
			t.Fatal(err)
		}
		received, err := ss.AuthConfig()
		if err != nil {
			t.Fatal(err)
		}

		if testCase.Logging.Fields == nil {
			testCase.Logging.Fields = map[string]string{}
		}
		if received.Logging.Fields == nil {
			received.Logging.Fields = map[string]string{}
		}

		if testCase.Auth.Parameters == nil {
			testCase.Auth.Parameters = make(garantconfig.Parameters)
		}
		if received.Auth.Parameters == nil {
			received.Auth.Parameters = make(garantconfig.Parameters)
		}

		if received.Version != testCase.Version {
			t.Fatalf("Versions not equal!\nReceived: %#v\nExpected: %#v", received.Version, testCase.Version)
		}

		if !reflect.DeepEqual(received.Logging, testCase.Logging) {
			t.Fatalf("Logging not equal!\nReceived: %#v\nExpected: %#v", received.Logging, testCase.Logging)
		}

		if received.SigningKey != testCase.SigningKey {
			t.Fatalf("SigningKeys not equal!\nReceived: %#v\nExpected: %#v", received.SigningKey, testCase.SigningKey)
		}

		if received.Issuer != testCase.Issuer {
			t.Fatalf("Issuers not equal!\nReceived: %#v\nExpected: %#v", received.Issuer, testCase.Issuer)
		}

		if !reflect.DeepEqual(received.Reporting, testCase.Reporting) {
			t.Fatalf("Reporting not equal!\nReceived: %#v\nExpected: %#v", received.Reporting, testCase.Reporting)
		}

		if !reflect.DeepEqual(received.HTTP, testCase.HTTP) {
			t.Fatalf("HTTP not equal!\nReceived: %#v\nExpected: %#v", received.HTTP, testCase.HTTP)
		}

		if !reflect.DeepEqual(received.Auth, testCase.Auth) {
			t.Fatalf("Auth not equal!\nReceived: %#v\nExpected: %#v", received.Auth, testCase.Auth)
		}
	}
}

func TestTLSCertificate(t *testing.T) {
	testCases := [][]byte{
		// Key then Cert
		[]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA2YY+nYx5JIYORwLyvnqzsQ8u9jV5kTyGMitZg6Ye9XPCocc/
5RaOtFHHdiB0QMYRJMf5X5PFwkyCO+7FOqBR9dgNr/MJbn3jGTkTxb4cS5Px4dAq
NB9/vLaYuh4ti3I+nEc+e7tRGmqWDnqBgH1q8IL5S1rn0Eu3sTwhmIqLlyL7pBwr
YyLUuDbmN7FbFBShumq8UOruV4L6sOBN3/eZpJjKKV7YGmPbteNEs2X0v3ueG5VB
cufCtF1ykGpGcejH9AcYyVRkc40C8lvxu9sbV17Bw8TZW8QDid0wCW0LgJiKWbRt
rKqqDMi+2Aw24JaYQfbRtt1I9a5oQ+OB1344pwIDAQABAoIBAAjYbalocfCrRt9a
XlailLYJgQZgDE58oJRfsOcqlS20BGEhwhZlwW3RuYOzNCcCJuZQ/3IEh1EsXUtk
nT6SEGMn4v6ZqnOEyPxYltY+sdXc8UQWe/DcqwwYfMNswLtf7O3b882CQ30Igoua
LbP05alcGhkXXD/bJGAfqtoSIDUIEuYUhXk70x7JrmVUAL0UJ/wxTL5gcDCAcFAx
Xd4N3YhIJVqK/XPXIFUpyFnOEbxlHn74Hc/XJiVj/GdoYDbKk49n22uJHqUcIl5T
z33MZjtb/YhCETUsr33Dt5CfYH0Pmo3+gtYryCxVpnZfsYfo4lLi5Dd548dwRbyD
T/EJ4KECgYEA/q/p0cctm9dZ7oUvX2HVLYaiY2117h+5FoGAYq4IT40gKIZA5r5I
zmCKx2/Il9+5PQR4HmeGTe1HtIuSOp4ZKrBnx9tRmxysFj9VLYnokhUbmLr2QvDt
5fG4IqlE7+/I0wABFkwwYWrfvrMLOEA+qBrUaZEd/AgXNa7xfwfsf5cCgYEA2qVK
gQRxTDlB6LYGDzBJ6xdSS6nNtLhy4PRUGPmrUZH6iiHiGQyDKh+jFB40/ViFS6kQ
9VLy/Qj4LO4x5g9K01kA7gvVbcWoNAcP0D18L/I5FZs9ry122yxxAQwcFumodu7H
h7nBDAHkHvo+AH1WeYpwlv1LtzwprB5OeR+CMXECgYB7J/+uNF2mrWVMhNTaj8lx
IVinMchEJMzwyCCYF0XnifvR/NGngr5cJa0WMcTTRBkkG2Qmd7MnPaVd9dv5Qngy
/2i/6Rs2IZBJlciPo891zIUyvr3UswKnHdMH18iBKfd3qNnduWvvv6mAYr37Ln3d
9lNe1RClzhfDGtymq+M98wKBgEqUxVQ7AragdVX++RQnQZ+ahezfUBbMMAuB7EnU
qFabt910b15iAT/WKNeM8kBU+Kr1UZ6NG+uqKpWQ7p5uKvTq0EFi8fuOx7BvlDpR
LtJgCo7PNHxTws5CW42i5tX+AAQJsTAf8bS51GmorIuYNW4iBgPLBQ/Myt42PEhe
shdRAoGAcshEBWe76THTPwZI2eE8fVkeXUElvpF+aNcqjERIQNemzAzvtvi0dYDe
A4qpEylTMG0BiHnkAvWokB37JKSMqBmPY6o2FsJ1b0IE5zNJgkdvehdbASTYJSia
82uhAOidQPoVKavm4SHYVN4+oUwP4JYdek8RApeL1u3ZIQ+fszQ=
-----END RSA PRIVATE KEY-----
-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAIQLUVbo2fP6MA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTAzMzExODQzMjZaFw0yNTAzMjgxODQz
MjZaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBANmGPp2MeSSGDkcC8r56s7EPLvY1eZE8hjIrWYOmHvVz
wqHHP+UWjrRRx3YgdEDGESTH+V+TxcJMgjvuxTqgUfXYDa/zCW594xk5E8W+HEuT
8eHQKjQff7y2mLoeLYtyPpxHPnu7URpqlg56gYB9avCC+Uta59BLt7E8IZiKi5ci
+6QcK2Mi1Lg25jexWxQUobpqvFDq7leC+rDgTd/3maSYyile2Bpj27XjRLNl9L97
nhuVQXLnwrRdcpBqRnHox/QHGMlUZHONAvJb8bvbG1dewcPE2VvEA4ndMAltC4CY
ilm0bayqqgzIvtgMNuCWmEH20bbdSPWuaEPjgdd+OKcCAwEAAaNQME4wHQYDVR0O
BBYEFGQUtOCedlcJwHvG3mBRM2Y+2cZ+MB8GA1UdIwQYMBaAFGQUtOCedlcJwHvG
3mBRM2Y+2cZ+MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAEEeMs0q
NSZkR4UxObWrSZ3LmrQTwLvhhECeXJHBIZXGfbgXOkbrxr1wHgAcYWQAUa507q91
+GGx+KaubzKa4qmQywaHnvR0cEzMozSI53eGxIiOKyI4sj4oUrJuVSfn+U8g6l4I
KDO/ztTASFKrqQS0gkCBUSRGCJhlBW6hS06lAQyzvn1rASsjZinMSEiuMx1qj9yM
TSCUwCIYy1ADZOFPt5LQvS2ytAypeDbEPTqwn83w06927Qg7jw3vsy3bHPfanpw4
zpwBMEzr5oxSQyEDLLQRVBvzMF86KKVKFSzmpcwtgju+WNtvl+hKG1gsuCuuX4P8
SBlNCMvHjCd+2vw=
-----END CERTIFICATE-----
`),
		// Cert then Key
		[]byte(`-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAIQLUVbo2fP6MA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTAzMzExODQzMjZaFw0yNTAzMjgxODQz
MjZaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBANmGPp2MeSSGDkcC8r56s7EPLvY1eZE8hjIrWYOmHvVz
wqHHP+UWjrRRx3YgdEDGESTH+V+TxcJMgjvuxTqgUfXYDa/zCW594xk5E8W+HEuT
8eHQKjQff7y2mLoeLYtyPpxHPnu7URpqlg56gYB9avCC+Uta59BLt7E8IZiKi5ci
+6QcK2Mi1Lg25jexWxQUobpqvFDq7leC+rDgTd/3maSYyile2Bpj27XjRLNl9L97
nhuVQXLnwrRdcpBqRnHox/QHGMlUZHONAvJb8bvbG1dewcPE2VvEA4ndMAltC4CY
ilm0bayqqgzIvtgMNuCWmEH20bbdSPWuaEPjgdd+OKcCAwEAAaNQME4wHQYDVR0O
BBYEFGQUtOCedlcJwHvG3mBRM2Y+2cZ+MB8GA1UdIwQYMBaAFGQUtOCedlcJwHvG
3mBRM2Y+2cZ+MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAEEeMs0q
NSZkR4UxObWrSZ3LmrQTwLvhhECeXJHBIZXGfbgXOkbrxr1wHgAcYWQAUa507q91
+GGx+KaubzKa4qmQywaHnvR0cEzMozSI53eGxIiOKyI4sj4oUrJuVSfn+U8g6l4I
KDO/ztTASFKrqQS0gkCBUSRGCJhlBW6hS06lAQyzvn1rASsjZinMSEiuMx1qj9yM
TSCUwCIYy1ADZOFPt5LQvS2ytAypeDbEPTqwn83w06927Qg7jw3vsy3bHPfanpw4
zpwBMEzr5oxSQyEDLLQRVBvzMF86KKVKFSzmpcwtgju+WNtvl+hKG1gsuCuuX4P8
SBlNCMvHjCd+2vw=
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA2YY+nYx5JIYORwLyvnqzsQ8u9jV5kTyGMitZg6Ye9XPCocc/
5RaOtFHHdiB0QMYRJMf5X5PFwkyCO+7FOqBR9dgNr/MJbn3jGTkTxb4cS5Px4dAq
NB9/vLaYuh4ti3I+nEc+e7tRGmqWDnqBgH1q8IL5S1rn0Eu3sTwhmIqLlyL7pBwr
YyLUuDbmN7FbFBShumq8UOruV4L6sOBN3/eZpJjKKV7YGmPbteNEs2X0v3ueG5VB
cufCtF1ykGpGcejH9AcYyVRkc40C8lvxu9sbV17Bw8TZW8QDid0wCW0LgJiKWbRt
rKqqDMi+2Aw24JaYQfbRtt1I9a5oQ+OB1344pwIDAQABAoIBAAjYbalocfCrRt9a
XlailLYJgQZgDE58oJRfsOcqlS20BGEhwhZlwW3RuYOzNCcCJuZQ/3IEh1EsXUtk
nT6SEGMn4v6ZqnOEyPxYltY+sdXc8UQWe/DcqwwYfMNswLtf7O3b882CQ30Igoua
LbP05alcGhkXXD/bJGAfqtoSIDUIEuYUhXk70x7JrmVUAL0UJ/wxTL5gcDCAcFAx
Xd4N3YhIJVqK/XPXIFUpyFnOEbxlHn74Hc/XJiVj/GdoYDbKk49n22uJHqUcIl5T
z33MZjtb/YhCETUsr33Dt5CfYH0Pmo3+gtYryCxVpnZfsYfo4lLi5Dd548dwRbyD
T/EJ4KECgYEA/q/p0cctm9dZ7oUvX2HVLYaiY2117h+5FoGAYq4IT40gKIZA5r5I
zmCKx2/Il9+5PQR4HmeGTe1HtIuSOp4ZKrBnx9tRmxysFj9VLYnokhUbmLr2QvDt
5fG4IqlE7+/I0wABFkwwYWrfvrMLOEA+qBrUaZEd/AgXNa7xfwfsf5cCgYEA2qVK
gQRxTDlB6LYGDzBJ6xdSS6nNtLhy4PRUGPmrUZH6iiHiGQyDKh+jFB40/ViFS6kQ
9VLy/Qj4LO4x5g9K01kA7gvVbcWoNAcP0D18L/I5FZs9ry122yxxAQwcFumodu7H
h7nBDAHkHvo+AH1WeYpwlv1LtzwprB5OeR+CMXECgYB7J/+uNF2mrWVMhNTaj8lx
IVinMchEJMzwyCCYF0XnifvR/NGngr5cJa0WMcTTRBkkG2Qmd7MnPaVd9dv5Qngy
/2i/6Rs2IZBJlciPo891zIUyvr3UswKnHdMH18iBKfd3qNnduWvvv6mAYr37Ln3d
9lNe1RClzhfDGtymq+M98wKBgEqUxVQ7AragdVX++RQnQZ+ahezfUBbMMAuB7EnU
qFabt910b15iAT/WKNeM8kBU+Kr1UZ6NG+uqKpWQ7p5uKvTq0EFi8fuOx7BvlDpR
LtJgCo7PNHxTws5CW42i5tX+AAQJsTAf8bS51GmorIuYNW4iBgPLBQ/Myt42PEhe
shdRAoGAcshEBWe76THTPwZI2eE8fVkeXUElvpF+aNcqjERIQNemzAzvtvi0dYDe
A4qpEylTMG0BiHnkAvWokB37JKSMqBmPY6o2FsJ1b0IE5zNJgkdvehdbASTYJSia
82uhAOidQPoVKavm4SHYVN4+oUwP4JYdek8RApeL1u3ZIQ+fszQ=
-----END RSA PRIVATE KEY-----
`),
	}
	for _, testCase := range testCases {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		ss := settingsstore.New(NewKeyValueStore(tempDir))
		expectedCert, err := tls.X509KeyPair(testCase, testCase)
		if err != nil {
			t.Fatal(err)
		}
		certBytes, err := util.CertificateToCertPEM(&expectedCert)
		if err != nil {
			t.Fatal(err)
		}
		keyBytes, err := util.CertificateToKeyPEM(&expectedCert)
		if err != nil {
			t.Fatal(err)
		}
		hubConfig, err := util.DefaultUserHubConfig("")
		if err != nil {
			t.Fatal(err)
		}
		hubConfig.WebTLSCert = string(certBytes)
		hubConfig.WebTLSKey = string(keyBytes)
		hubConfig.WebTLSCA = string(certBytes)
		err = ss.SetUserHubConfig(hubConfig)
		if err != nil {
			t.Fatal(err)
		}
		hubConfig2, err := ss.UserHubConfig()
		if err != nil {
			t.Fatal(err)
		}
		if hubConfig2.WebTLSCert != hubConfig.WebTLSCert {
			t.Fatal("webTLSCert mismatch")
		}
		if hubConfig2.WebTLSKey != hubConfig.WebTLSKey {
			t.Fatal("webTLSCert mismatch")
		}
		if hubConfig2.WebTLSCA != hubConfig.WebTLSCA {
			t.Fatal("webTLSCA mismatch")
		}
	}
}

func TestRegistryConfig(t *testing.T) {
	var eventsNotificationEndpoint = configuration.Endpoint{
		Name:    "audit-logs",
		URL:     fmt.Sprintf("http://%s/%s/", containers.APIServer.BridgeNameLocalReplica(), deploy.EventsEndpointSubroute),
		Headers: http.Header{"X-Registry-Events": []string{"true"}},
	}
	var testCases = []configuration.Configuration{
		func() configuration.Configuration {
			config := configuration.Configuration{}
			config.Notifications.Endpoints = []configuration.Endpoint{
				{
					Name:    "some-endpoint",
					URL:     "https://some-url/asdf",
					Headers: http.Header{"X-Header-Name": []string{"some value"}},
				},
			}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			return config
		}(),
		func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "//docker.com/auth/v2/token/",
					"issuer":         "docker.com",
					"service":        "docker.com",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Notifications.Endpoints = []configuration.Endpoint{eventsNotificationEndpoint}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			return config
		}(),
		func() configuration.Configuration {
			config := configuration.Configuration{}
			config.HTTP.Addr = ":5000"
			config.Auth = configuration.Auth{
				"token": configuration.Parameters{
					"realm":          "//my.nested.domain.org/auth/v2/token/",
					"issuer":         "my.nested.domain.org",
					"service":        "my.nested.domain.org",
					"rootcertbundle": filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
				},
			}
			config.Notifications.Endpoints = []configuration.Endpoint{eventsNotificationEndpoint}
			config.Storage = defaultconfigs.DefaultRegistryConfig.Storage
			return config
		}(),
	}
	for _, testCase := range testCases {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		ss := settingsstore.New(NewKeyValueStore(tempDir))
		err = ss.SetRegistryConfig(&testCase)
		if err != nil {
			t.Fatal(err)
		}
		received, err := ss.RegistryConfig()
		if err != nil {
			t.Fatal(err)
		}

		if testCase.Log.Fields == nil {
			testCase.Log.Fields = map[string]interface{}{}
		}
		if received.Log.Fields == nil {
			received.Log.Fields = map[string]interface{}{}
		}

		// don't compare the storage config because we modify the path
		received.Storage = nil
		testCase.Storage = nil
		if !reflect.DeepEqual(*received, testCase) {
			t.Fatalf("Configs not equal!\nReceived: %#v\nExpected: %#v", *received, testCase)
		}
	}
}

func TestHubCredentials(t *testing.T) {
	var testCases = []dockerclient.AuthConfig{
		// Empty config
		{},
		{
			Username: "brianbland",
			Password: "P4$$w0rd!",
			Email:    "testemail@docker.com",
		},
	}
	for _, testCase := range testCases {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		ss := settingsstore.New(NewKeyValueStore(tempDir))
		err = ss.SetHubCredentials(&testCase)
		if err != nil {
			t.Fatal(err)
		}
		received, err := ss.HubCredentials()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(*received, testCase) {
			t.Fatalf("Configs not equal!\nReceived: %#v\nExpected: %#v", *received, testCase)
		}
	}
}
