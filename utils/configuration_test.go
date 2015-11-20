package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const envPrefix = "NOTARY_TESTING_ENV_PREFIX"

// initializes a viper object with test configuration
func configure(jsonConfig string) *viper.Viper {
	config := viper.New()
	SetupViper(config, envPrefix)
	config.SetConfigType("json")
	config.ReadConfig(bytes.NewBuffer([]byte(jsonConfig)))
	return config
}

// Sets the environment variables in the given map, prefixed by envPrefix.
func setupEnvironmentVariables(t *testing.T, vars map[string]string) {
	for k, v := range vars {
		err := os.Setenv(fmt.Sprintf("%s_%s", envPrefix, k), v)
		assert.NoError(t, err)
	}
}

// Unsets whatever environment variables were set with this map
func cleanupEnvironmentVariables(t *testing.T, vars map[string]string) {
	for k := range vars {
		err := os.Unsetenv(fmt.Sprintf("%s_%s", envPrefix, k))
		assert.NoError(t, err)
	}

}

// An error is returned if the log level is not parsable
func TestParseInvalidLogLevel(t *testing.T) {
	_, err := ParseLogLevel(configure(`{"logging": {"level": "horatio"}}`),
		logrus.DebugLevel)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid logrus Level")
}

// If there is no logging level configured it is set to the default level
func TestParseNoLogLevel(t *testing.T) {
	empties := []string{`{}`, `{"logging": {}}`}
	for _, configJSON := range empties {
		lvl, err := ParseLogLevel(configure(configJSON), logrus.DebugLevel)
		assert.NoError(t, err)
		assert.Equal(t, logrus.DebugLevel, lvl)
	}
}

// If there is logging level configured, it is set to the configured one
func TestParseLogLevel(t *testing.T) {
	lvl, err := ParseLogLevel(configure(`{"logging": {"level": "error"}}`),
		logrus.DebugLevel)
	assert.NoError(t, err)
	assert.Equal(t, logrus.ErrorLevel, lvl)
}

func TestParseLogLevelWithEnvironmentVariables(t *testing.T) {
	vars := map[string]string{"LOGGING_LEVEL": "error"}
	setupEnvironmentVariables(t, vars)
	defer cleanupEnvironmentVariables(t, vars)

	lvl, err := ParseLogLevel(configure(`{}`),
		logrus.DebugLevel)
	assert.NoError(t, err)
	assert.Equal(t, logrus.ErrorLevel, lvl)
}

// An error is returned if there's no API key
func TestParseInvalidBugsnag(t *testing.T) {
	_, err := ParseBugsnag(configure(
		`{"reporting": {"bugsnag": {"endpoint": "http://12345"}}}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must provide an API key")
}

// If there's no bugsnag, a nil pointer is returned
func TestParseNoBugsnag(t *testing.T) {
	empties := []string{`{}`, `{"reporting": {}}`}
	for _, configJSON := range empties {
		bugconf, err := ParseBugsnag(configure(configJSON))
		assert.NoError(t, err)
		assert.Nil(t, bugconf)
	}
}

func TestParseBugsnag(t *testing.T) {
	config := configure(`{
		"reporting": {
			"bugsnag": {
				"api_key": "12345",
				"release_stage": "production",
				"endpoint": "http://1234.com"
			}
		}
	}`)

	expected := bugsnag.Configuration{
		APIKey:       "12345",
		ReleaseStage: "production",
		Endpoint:     "http://1234.com",
	}

	bugconf, err := ParseBugsnag(config)
	assert.NoError(t, err)
	assert.Equal(t, expected, *bugconf)
}

func TestParseBugsnagWithEnvironmentVariables(t *testing.T) {
	config := configure(`{
		"reporting": {
			"bugsnag": {
				"api_key": "12345",
				"release_stage": "staging"
			}
		}
	}`)

	vars := map[string]string{
		"REPORTING_BUGSNAG_RELEASE_STAGE": "production",
		"REPORTING_BUGSNAG_ENDPOINT":      "http://1234.com",
	}
	setupEnvironmentVariables(t, vars)
	defer cleanupEnvironmentVariables(t, vars)

	expected := bugsnag.Configuration{
		APIKey:       "12345",
		ReleaseStage: "production",
		Endpoint:     "http://1234.com",
	}

	bugconf, err := ParseBugsnag(config)
	assert.NoError(t, err)
	assert.Equal(t, expected, *bugconf)
}

// If the storage parameters are invalid, an error is returned
func TestParseInvalidStorage(t *testing.T) {
	invalids := []string{
		`{"storage": {"backend": "postgres", "db_url": "1234"}}`,
		`{"storage": {"db_url": "12345"}}`,
		`{"storage": {"backend": "mysql"}}`,
		`{"storage": {"backend": "sqlite3", "db_url": ""}}`,
	}
	for _, configJSON := range invalids {
		_, err := ParseStorage(configure(configJSON), []string{"mysql", "sqlite3"})
		assert.Error(t, err, fmt.Sprintf("'%s' should be an error", configJSON))
		if strings.Contains(configJSON, "mysql") || strings.Contains(configJSON, "sqlite3") {
			assert.Contains(t, err.Error(),
				"must provide a non-empty database source")
		} else {
			assert.Contains(t, err.Error(),
				"must specify one of these supported backends: mysql, sqlite3")
		}
	}
}

// If there is no storage, a nil pointer is returned
func TestParseNoStorage(t *testing.T) {
	empties := []string{`{}`, `{"storage": {}}`}
	for _, configJSON := range empties {
		store, err := ParseStorage(configure(configJSON), []string{"mysql"})
		assert.NoError(t, err)
		assert.Nil(t, store)
	}
}

func TestParseStorage(t *testing.T) {
	config := configure(`{
		"storage": {
			"backend": "MySQL",
			"db_url": "username:passord@tcp(hostname:1234)/dbname"
		}
	}`)

	expected := Storage{
		Backend: "mysql",
		Source:  "username:passord@tcp(hostname:1234)/dbname",
	}

	store, err := ParseStorage(config, []string{"mysql"})
	assert.NoError(t, err)
	assert.Equal(t, expected, *store)
}

func TestParseStorageWithEnvironmentVariables(t *testing.T) {
	config := configure(`{
		"storage": {
			"db_url": "username:passord@tcp(hostname:1234)/dbname"
		}
	}`)

	vars := map[string]string{"STORAGE_BACKEND": "MySQL"}
	setupEnvironmentVariables(t, vars)
	defer cleanupEnvironmentVariables(t, vars)

	expected := Storage{
		Backend: "mysql",
		Source:  "username:passord@tcp(hostname:1234)/dbname",
	}

	store, err := ParseStorage(config, []string{"mysql"})
	assert.NoError(t, err)
	assert.Equal(t, expected, *store)
}

// If TLS is required and the parameters are missing, an error is returned
func TestParseTLSNoTLSWhenRequired(t *testing.T) {
	invalids := []string{
		`{"server": {"tls_cert_file": "path/to/cert"}}`,
		`{"server": {"tls_key_file": "path/to/key"}}`,
	}
	for _, configJSON := range invalids {
		_, err := ParseServerTLS(configure(configJSON), true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(),
			"both the TLS certificate and key are mandatory")
	}
}

// If TLS is not and the cert/key are partially provided, an error is returned
func TestParseTLSPartialTLS(t *testing.T) {
	invalids := []string{
		`{"server": {"tls_cert_file": "path/to/cert"}}`,
		`{"server": {"tls_key_file": "path/to/key"}}`,
	}
	for _, configJSON := range invalids {
		_, err := ParseServerTLS(configure(configJSON), false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(),
			"either include both a cert and key file, or neither to disable TLS")
	}
}

func TestParseTLSNoTLSNotRequired(t *testing.T) {
	config := configure(`{
		"server": {}
	}`)

	tlsOpts, err := ParseServerTLS(config, false)
	assert.NoError(t, err)
	assert.Nil(t, tlsOpts)
}

func TestParseTLSWithTLS(t *testing.T) {
	config := configure(`{
		"server": {
			"tls_cert_file": "path/to/cert",
			"tls_key_file": "path/to/key",
			"client_ca_file": "path/to/clientca"
		}
	}`)

	expected := ServerTLSOpts{
		ServerCertFile: "path/to/cert",
		ServerKeyFile:  "path/to/key",
		ClientCAFile:   "path/to/clientca",
	}

	tlsOpts, err := ParseServerTLS(config, false)
	assert.NoError(t, err)
	assert.Equal(t, expected, *tlsOpts)
}

func TestParseTLSWithTLSRelativeToConfigFile(t *testing.T) {
	config := configure(`{
		"server": {
			"tls_cert_file": "path/to/cert",
			"tls_key_file": "/abspath/to/key",
			"client_ca_file": ""
		}
	}`)
	config.SetConfigFile("/opt/me.json")

	expected := ServerTLSOpts{
		ServerCertFile: "/opt/path/to/cert",
		ServerKeyFile:  "/abspath/to/key",
		ClientCAFile:   "",
	}

	tlsOpts, err := ParseServerTLS(config, false)
	assert.NoError(t, err)
	assert.Equal(t, expected, *tlsOpts)
}

func TestParseTLSWithEnvironmentVariables(t *testing.T) {
	config := configure(`{
		"server": {
			"tls_cert_file": "path/to/cert",
			"client_ca_file": "nosuchfile"
		}
	}`)

	vars := map[string]string{
		"SERVER_TLS_KEY_FILE":   "path/to/key",
		"SERVER_CLIENT_CA_FILE": "path/to/clientca",
	}
	setupEnvironmentVariables(t, vars)
	defer cleanupEnvironmentVariables(t, vars)

	expected := ServerTLSOpts{
		ServerCertFile: "path/to/cert",
		ServerKeyFile:  "path/to/key",
		ClientCAFile:   "path/to/clientca",
	}

	tlsOpts, err := ParseServerTLS(config, true)
	assert.NoError(t, err)
	assert.Equal(t, expected, *tlsOpts)
}
