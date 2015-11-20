package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/notary/server/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	Cert = "../../fixtures/notary-server.crt"
	Key  = "../../fixtures/notary-server.key"
	Root = "../../fixtures/root-ca.crt"
)

// initializes a viper object with test configuration
func configure(jsonConfig string) *viper.Viper {
	config := viper.New()
	config.SetConfigType("json")
	config.ReadConfig(bytes.NewBuffer([]byte(jsonConfig)))
	return config
}

func TestGetAddrAndTLSConfigInvalidTLS(t *testing.T) {
	invalids := []string{
		`{"server": {
				"http_addr": ":1234",
				"tls_key_file": "nope"
		}}`,
	}
	for _, configJSON := range invalids {
		_, _, err := getAddrAndTLSConfig(configure(configJSON))
		assert.Error(t, err)
	}
}

func TestGetAddrAndTLSConfigNoHTTPAddr(t *testing.T) {
	_, _, err := getAddrAndTLSConfig(configure(fmt.Sprintf(`{
		"server": {
			"tls_cert_file": "%s",
			"tls_key_file": "%s"
		}
	}`, Cert, Key)))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "http listen address required for server")
}

func TestGetAddrAndTLSConfigSuccessWithTLS(t *testing.T) {
	httpAddr, tlsConf, err := getAddrAndTLSConfig(configure(fmt.Sprintf(`{
		"server": {
			"http_addr": ":2345",
			"tls_cert_file": "%s",
			"tls_key_file": "%s"
		}
	}`, Cert, Key)))
	assert.NoError(t, err)
	assert.Equal(t, ":2345", httpAddr)
	assert.NotNil(t, tlsConf)
}

func TestGetAddrAndTLSConfigSuccessWithoutTLS(t *testing.T) {
	httpAddr, tlsConf, err := getAddrAndTLSConfig(configure(
		`{"server": {"http_addr": ":2345"}}`))
	assert.NoError(t, err)
	assert.Equal(t, ":2345", httpAddr)
	assert.Nil(t, tlsConf)
}

// We don't support client CAs yet on notary server
func TestGetAddrAndTLSConfigSkipClientTLS(t *testing.T) {
	httpAddr, tlsConf, err := getAddrAndTLSConfig(configure(fmt.Sprintf(`{
		"server": {
			"http_addr": ":2345",
			"tls_cert_file": "%s",
			"tls_key_file": "%s",
			"client_ca_file": "%s"
		}
	}`, Cert, Key, Root)))
	assert.NoError(t, err)
	assert.Equal(t, ":2345", httpAddr)
	assert.Nil(t, tlsConf.ClientCAs)
}

// Client cert and Key either both have to be empty or both have to be
// provided.
func TestGrpcTLSMissingCertOrKey(t *testing.T) {
	configs := []string{
		fmt.Sprintf(`"tls_client_cert": "%s"`, Cert),
		fmt.Sprintf(`"tls_client_key": "%s"`, Key),
	}
	for _, trustConfig := range configs {
		jsonConfig := fmt.Sprintf(
			`{"trust_service": {"hostname": "notary-signer", %s}}`,
			trustConfig)
		config := configure(jsonConfig)
		tlsConfig, err := grpcTLS(config)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.True(t,
			strings.Contains(err.Error(), "Partial TLS configuration found."))
	}
}

// If no TLS configuration is provided for the host server, a tls config with
// the provided serverName is still returned.
func TestGrpcTLSNoConfig(t *testing.T) {
	tlsConfig, err := grpcTLS(
		configure(`{"trust_service": {"hostname": "notary-signer"}}`))
	assert.NoError(t, err)
	assert.Equal(t, "notary-signer", tlsConfig.ServerName)
	assert.Nil(t, tlsConfig.RootCAs)
	assert.Nil(t, tlsConfig.Certificates)
}

// The rest of the functionality of grpcTLS depends upon
// utils.ConfigureClientTLS, so this test just asserts that if successful,
// the correct tls.Config is returned based on all the configuration parameters
func TestGrpcTLSSuccess(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(Cert, Key)
	assert.NoError(t, err, "Unable to load cert and key for testing")

	config := fmt.Sprintf(
		`{"trust_service": {
            "hostname": "notary-server",
            "tls_client_cert": "%s",
            "tls_client_key": "%s"}}`,
		Cert, Key)
	tlsConfig, err := grpcTLS(configure(config))
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
}

// The rest of the functionality of grpcTLS depends upon
// utils.ConfigureServerTLS, so this test just asserts that if it fails,
// the error is propogated.
func TestGrpcTLSFailure(t *testing.T) {
	config := fmt.Sprintf(
		`{"trust_service": {
            "hostname": "notary-server",
            "tls_client_cert": "no-exist",
            "tls_client_key": "%s"}}`,
		Key)
	tlsConfig, err := grpcTLS(configure(config))
	assert.Error(t, err)
	assert.Nil(t, tlsConfig)
	assert.True(t, strings.Contains(err.Error(),
		"Unable to configure TLS to the trust service"))
}

// Just to ensure that errors are propogated
func TestGetStoreInvalid(t *testing.T) {
	config := `{"storage": {"backend": "asdf", "db_url": "/tmp/1234"}}`

	_, err := getStore(configure(config), []string{"mysql"})
	assert.Error(t, err)
}

func TestGetStoreDBStore(t *testing.T) {
	tmpFile, err := ioutil.TempFile("/tmp", "sqlite3")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	config := fmt.Sprintf(`{"storage": {"backend": "sqlite3", "db_url": "%s"}}`,
		tmpFile.Name())

	store, err := getStore(configure(config), []string{"sqlite3"})
	assert.NoError(t, err)
	_, ok := store.(*storage.SQLStorage)
	assert.True(t, ok)
}

func TestGetMemoryStore(t *testing.T) {
	config := fmt.Sprintf(`{"storage": {}}`)
	store, err := getStore(configure(config), []string{"mysql"})
	assert.NoError(t, err)
	_, ok := store.(*storage.MemStorage)
	assert.True(t, ok)
}
