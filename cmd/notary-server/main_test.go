package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/signer/client"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/utils"
	_ "github.com/mattn/go-sqlite3"
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

func TestGetAddrAndTLSConfigWithClientTLS(t *testing.T) {
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
	assert.NotNil(t, tlsConf.ClientCAs)
}

// If neither "remote" nor "local" is passed for "trust_service.type", an
// error is returned.
func TestGetInvalidTrustService(t *testing.T) {
	invalids := []string{
		`{"trust_service": {"type": "bruhaha", "key_algorithm": "rsa"}}`,
		`{}`,
	}
	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	for _, config := range invalids {
		_, _, err := getTrustService(configure(config),
			client.NewNotarySigner, fakeRegister)
		assert.Error(t, err)
		assert.Contains(t, err.Error(),
			"must specify either a \"local\" or \"remote\" type for trust_service")
	}
	// no health function ever registered
	assert.Equal(t, 0, registerCalled)
}

// If a local trust service is specified, a local trust service will be used
// with an ED22519 algorithm no matter what algorithm was specified.  No health
// function is configured.
func TestGetLocalTrustService(t *testing.T) {
	localConfig := `{"trust_service": {"type": "local", "key_algorithm": "meh"}}`

	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	trust, algo, err := getTrustService(configure(localConfig),
		client.NewNotarySigner, fakeRegister)
	assert.NoError(t, err)
	assert.IsType(t, &signed.Ed25519{}, trust)
	assert.Equal(t, data.ED25519Key, algo)

	// no health function ever registered
	assert.Equal(t, 0, registerCalled)
}

// Invalid key algorithms result in an error if a remote trust service was
// specified.
func TestGetTrustServiceInvalidKeyAlgorithm(t *testing.T) {
	configTemplate := `
	{
		"trust_service": {
			"type": "remote",
			"hostname": "blah",
			"port": "1234",
			"key_algorithm": "%s"
		}
	}`
	badKeyAlgos := []string{
		fmt.Sprintf(configTemplate, ""),
		fmt.Sprintf(configTemplate, data.ECDSAx509Key),
		fmt.Sprintf(configTemplate, "random"),
	}
	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	for _, config := range badKeyAlgos {
		_, _, err := getTrustService(configure(config),
			client.NewNotarySigner, fakeRegister)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid key algorithm")
	}
	// no health function ever registered
	assert.Equal(t, 0, registerCalled)
}

// template to be used for testing TLS parsing with the trust service
var trustTLSConfigTemplate = `
	{
		"trust_service": {
			"type": "remote",
			"hostname": "notary-signer",
			"port": "1234",
			"key_algorithm": "ecdsa",
			%s
		}
	}`

// Client cert and Key either both have to be empty or both have to be
// provided.
func TestGetTrustServiceTLSMissingCertOrKey(t *testing.T) {
	configs := []string{
		fmt.Sprintf(`"tls_client_cert": "%s"`, Cert),
		fmt.Sprintf(`"tls_client_key": "%s"`, Key),
	}
	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	for _, clientTLSConfig := range configs {
		jsonConfig := fmt.Sprintf(trustTLSConfigTemplate, clientTLSConfig)
		config := configure(jsonConfig)
		_, _, err := getTrustService(config, client.NewNotarySigner,
			fakeRegister)
		assert.Error(t, err)
		assert.True(t,
			strings.Contains(err.Error(), "either pass both client key and cert, or neither"))
	}
	// no health function ever registered
	assert.Equal(t, 0, registerCalled)
}

// If no TLS configuration is provided for the host server, no TLS config will
// be set for the trust service.
func TestGetTrustServiceNoTLSConfig(t *testing.T) {
	config := `{
		"trust_service": {
			"type": "remote",
			"hostname": "notary-signer",
			"port": "1234",
			"key_algorithm": "ecdsa"
		}
	}`
	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	var tlsConfig *tls.Config
	var fakeNewSigner = func(_, _ string, c *tls.Config) *client.NotarySigner {
		tlsConfig = c
		return &client.NotarySigner{}
	}

	trust, algo, err := getTrustService(configure(config),
		fakeNewSigner, fakeRegister)
	assert.NoError(t, err)
	assert.IsType(t, &client.NotarySigner{}, trust)
	assert.Equal(t, "ecdsa", algo)
	assert.Nil(t, tlsConfig.RootCAs)
	assert.Nil(t, tlsConfig.Certificates)
	// health function registered
	assert.Equal(t, 1, registerCalled)
}

// The rest of the functionality of getTrustService depends upon
// utils.ConfigureClientTLS, so this test just asserts that if successful,
// the correct tls.Config is returned based on all the configuration parameters
func TestGetTrustServiceTLSSuccess(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(Cert, Key)
	assert.NoError(t, err, "Unable to load cert and key for testing")

	tlspart := fmt.Sprintf(`"tls_client_cert": "%s", "tls_client_key": "%s"`,
		Cert, Key)

	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	var tlsConfig *tls.Config
	var fakeNewSigner = func(_, _ string, c *tls.Config) *client.NotarySigner {
		tlsConfig = c
		return &client.NotarySigner{}
	}

	trust, algo, err := getTrustService(
		configure(fmt.Sprintf(trustTLSConfigTemplate, tlspart)),
		fakeNewSigner, fakeRegister)
	assert.NoError(t, err)
	assert.IsType(t, &client.NotarySigner{}, trust)
	assert.Equal(t, "ecdsa", algo)
	assert.Len(t, tlsConfig.Certificates, 1)
	assert.True(t, reflect.DeepEqual(keypair, tlsConfig.Certificates[0]))
	// health function registered
	assert.Equal(t, 1, registerCalled)
}

// The rest of the functionality of getTrustService depends upon
// utils.ConfigureServerTLS, so this test just asserts that if it fails,
// the error is propagated.
func TestGetTrustServiceTLSFailure(t *testing.T) {
	tlspart := fmt.Sprintf(`"tls_client_cert": "none", "tls_client_key": "%s"`,
		Key)

	var registerCalled = 0
	var fakeRegister = func(_ string, _ func() error, _ time.Duration) {
		registerCalled++
	}

	_, _, err := getTrustService(
		configure(fmt.Sprintf(trustTLSConfigTemplate, tlspart)),
		client.NewNotarySigner, fakeRegister)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(),
		"Unable to configure TLS to the trust service"))

	// no health function ever registered
	assert.Equal(t, 0, registerCalled)
}

// Just to ensure that errors are propagated
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

	config := fmt.Sprintf(`{"storage": {"backend": "%s", "db_url": "%s"}}`,
		utils.SqliteBackend, tmpFile.Name())

	store, err := getStore(configure(config), []string{utils.SqliteBackend})
	assert.NoError(t, err)
	_, ok := store.(*storage.SQLStorage)
	assert.True(t, ok)
}

func TestGetMemoryStore(t *testing.T) {
	config := fmt.Sprintf(`{"storage": {"backend": "%s"}}`, utils.MemoryBackend)
	store, err := getStore(configure(config),
		[]string{utils.MySQLBackend, utils.MemoryBackend})
	assert.NoError(t, err)
	_, ok := store.(*storage.MemStorage)
	assert.True(t, ok)
}

func TestGetCacheConfig(t *testing.T) {
	valid := `{"caching": {"max_age": {"current_metadata": 0, "metadata_by_checksum": 31536000}}}`
	invalids := []string{
		`{"caching": {"max_age": {"current_metadata": 0, "metadata_by_checksum": 31539000}}}`,
		`{"caching": {"max_age": {"current_metadata": -1, "metadata_by_checksum": 300}}}`,
		`{"caching": {"max_age": {"current_metadata": "hello", "metadata_by_checksum": 300}}}`,
	}

	current, consistent, err := getCacheConfig(configure(valid))
	assert.NoError(t, err)
	assert.IsType(t, utils.NoCacheControl{}, current)
	assert.IsType(t, utils.PublicCacheControl{}, consistent)

	for _, invalid := range invalids {
		_, _, err := getCacheConfig(configure(invalid))
		assert.Error(t, err)
	}
}
