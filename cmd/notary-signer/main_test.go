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
	Cert = "../../fixtures/notary-signer.crt"
	Key  = "../../fixtures/notary-signer.key"
	Root = "../../fixtures/root-ca.crt"
)

// initializes a viper object with test configuration
func configure(jsonConfig []byte) *viper.Viper {
	config := viper.New()
	config.SetConfigType("json")
	config.ReadConfig(bytes.NewBuffer(jsonConfig))
	return config
}

func TestSignerTLSMissingCertAndOrKey(t *testing.T) {
	configs := []string{
		"{}",
		fmt.Sprintf(`{"cert_file": "%s"}`, Cert),
		fmt.Sprintf(`{"key_file": "%s"}`, Key),
	}
	for _, serverConfig := range configs {
		config := configure(
			[]byte(fmt.Sprintf(`{"server": %s}`, serverConfig)))
		tlsConfig, err := signerTLS(config, false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Equal(t, "Certificate and key are mandatory", err.Error())
	}
}

// The rest of the functionality of singerTLS depends upon
// utils.ConfigureServerTLS, so this test just asserts that if successful,
// the correct tls.Config is returned based on all the configuration parameters
func TestSignerTLSSuccess(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(Cert, Key)
	assert.NoError(t, err, "Unable to load cert and key for testing")

	config := fmt.Sprintf(
		`{"server": {"cert_file": "%s", "key_file": "%s", "client_ca_file": "%s"}}`,
		Cert, Key, Cert)
	tlsConfig, err := signerTLS(configure([]byte(config)), false)
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.NotNil(t, tlsConfig.ClientCAs)
}

// The rest of the functionality of singerTLS depends upon
// utils.ConfigureServerTLS, so this test just asserts that if it fails,
// the error is propogated.
func TestSignerTLSFailure(t *testing.T) {
	config := fmt.Sprintf(
		`{"server": {"cert_file": "%s", "key_file": "%s", "client_ca_file": "%s"}}`,
		Cert, Key, "non-existant")
	tlsConfig, err := signerTLS(configure([]byte(config)), false)
	assert.Error(t, err)
	assert.Nil(t, tlsConfig)
	assert.True(t, strings.Contains(err.Error(), "Unable to set up TLS"))
}
