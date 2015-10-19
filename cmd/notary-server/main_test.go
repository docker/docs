package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	Cert = "../../fixtures/notary-server.crt"
	Key  = "../../fixtures/notary-server.key"
	Root = "../../fixtures/root-ca.crt"
)

// initializes a viper object with test configuration
func configure(jsonConfig []byte) *viper.Viper {
	config := viper.New()
	config.SetConfigType("json")
	config.ReadConfig(bytes.NewBuffer(jsonConfig))
	return config
}

// If neither the cert nor the key are provided, a nil tls config is returned.
func TestServerTLSMissingCertAndKey(t *testing.T) {
	tlsConfig, err := serverTLS(configure([]byte(`{"server": {}}`)))
	assert.NoError(t, err)
	assert.Nil(t, tlsConfig)
}

func TestServerTLSMissingCertAndOrKey(t *testing.T) {
	configs := []string{
		fmt.Sprintf(`{"tls_cert_file": "%s"}`, Cert),
		fmt.Sprintf(`{"tls_key_file": "%s"}`, Key),
	}
	for _, serverConfig := range configs {
		config := configure(
			[]byte(fmt.Sprintf(`{"server": %s}`, serverConfig)))
		tlsConfig, err := serverTLS(config)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.True(t,
			strings.Contains(err.Error(), "Partial TLS configuration found."))
	}
}

// The rest of the functionality of serverTLS depends upon
// utils.ConfigureServerTLS, so this test just asserts that if successful,
// the correct tls.Config is returned based on all the configuration parameters
func TestServerTLSSuccess(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(Cert, Key)
	assert.NoError(t, err, "Unable to load cert and key for testing")

	config := fmt.Sprintf(
		`{"server": {"tls_cert_file": "%s", "tls_key_file": "%s"}}`,
		Cert, Key)
	tlsConfig, err := serverTLS(configure([]byte(config)))
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
}

// The rest of the functionality of singerTLS depends upon
// utils.ConfigureServerTLS, so this test just asserts that if it fails,
// the error is propogated.
func TestServerTLSFailure(t *testing.T) {
	config := fmt.Sprintf(
		`{"server": {"tls_cert_file": "non-exist", "tls_key_file": "%s"}}`,
		Key)
	tlsConfig, err := serverTLS(configure([]byte(config)))
	assert.Error(t, err)
	assert.Nil(t, tlsConfig)
	assert.True(t, strings.Contains(err.Error(), "Unable to set up TLS"))
}
