package utils

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// initializes a viper object with test configuration
func configure(jsonConfig string) *viper.Viper {
	config := viper.New()
	config.SetConfigType("json")
	config.ReadConfig(bytes.NewBuffer([]byte(jsonConfig)))
	return config
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

// If the storage parameters are invalid, an error is returned
func TestParseInvalidStorage(t *testing.T) {
	invalids := []string{
		`{"storage": {"backend": "memow", "db_url": "1234"}}`,
		`{"storage": {"db_url": "12345"}}`,
		`{"storage": {"backend": "mysql"}}`,
		`{"storage": {"backend": "mysql", "db_url": ""}}`,
	}
	for _, configJSON := range invalids {
		_, err := ParseStorage(configure(configJSON))
		assert.Error(t, err, fmt.Sprintf("'%s' should be an error", configJSON))
		if strings.Contains(configJSON, "mysql") {
			assert.Contains(t, err.Error(),
				"must provide a non-empty database URL")
		} else {
			assert.Contains(t, err.Error(),
				"must specify one of these supported backends: mysql")
		}
	}
}

// If there is no storage, a nil pointer is returned
func TestParseNoStorage(t *testing.T) {
	empties := []string{`{}`, `{"storage": {}}`}
	for _, configJSON := range empties {
		store, err := ParseStorage(configure(configJSON))
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
		URL:     "username:passord@tcp(hostname:1234)/dbname",
	}

	store, err := ParseStorage(config)
	assert.NoError(t, err)
	assert.Equal(t, expected, *store)
}

// If the server section is missing or missing HTTP/GRPC addresses, an error is
// returned
func TestParseInvalidOrNoServer(t *testing.T) {
	invalids := []string{`{}`, `{"server": {}}`}
	for _, configJSON := range invalids {
		_, err := ParseServer(configure(configJSON), false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have an HTTP and/or GRPC address")
	}
}

// If TLS is required and the parameters are missing, an error is returned
func TestParseInvalidServerNoTLSWhenRequired(t *testing.T) {
	invalids := []string{
		`{"server": {"http_addr": ":443", "tls_cert_file": "path/to/cert"}}`,
		`{"server": {"http_addr": ":443", "tls_key_file": "path/to/key"}}`,
	}
	for _, configJSON := range invalids {
		_, err := ParseServer(configure(configJSON), true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(),
			"both the TLS certificate and key are mandatory")
	}
}

// If TLS is not and the cert/key are partially provided, an error is returned
func TestParseInvalidServerPartialTLS(t *testing.T) {
	invalids := []string{
		`{"server": {"http_addr": ":443", "tls_cert_file": "path/to/cert"}}`,
		`{"server": {"http_addr": ":443", "tls_key_file": "path/to/key"}}`,
	}
	for _, configJSON := range invalids {
		_, err := ParseServer(configure(configJSON), false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(),
			"either include both a cert and key file, or neither to disable TLS")
	}
}

func TestParseServerNoTLS(t *testing.T) {
	config := configure(`{
		"server": {
			"http_addr": ":4443",
			"grpc_addr": ":7899"
		}
	}`)

	expected := Server{
		HTTPAddr:      ":4443",
		GRPCAddr:      ":7899",
		ServerTLSOpts: nil,
	}

	server, err := ParseServer(config, false)
	assert.NoError(t, err)
	assert.Equal(t, expected, *server)
}

func TestUnmarshalConfigServerWithTLS(t *testing.T) {
	config := configure(`{
		"server": {
			"http_addr": ":4443",
			"grpc_addr": ":7899",
			"tls_cert_file": "path/to/cert",
			"tls_key_file": "path/to/key",
			"client_ca_file": "path/to/clientca"
		}
	}`)

	expected := Server{
		HTTPAddr: ":4443",
		GRPCAddr: ":7899",
		ServerTLSOpts: &ServerTLSOpts{
			ServerCertFile: "path/to/cert",
			ServerKeyFile:  "path/to/key",
			ClientCAFile:   "path/to/clientca",
		},
	}

	server, err := ParseServer(config, false)
	assert.NoError(t, err)
	assert.Equal(t, expected, *server)
}
