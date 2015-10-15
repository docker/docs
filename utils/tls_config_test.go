package utils

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ServerCert = "../fixtures/notary-server.crt"
	ServerKey  = "../fixtures/notary-server.key"
	CACert     = "../fixtures/root-ca.crt"
	FakeFile   = "../fixtures/not-a-file"
)

// TestTLSConfigFailsIfUnableToLoadCerts fails if unable to load either of the
// server files or the client cert info
func TestConfigServerTLSFailsIfUnableToLoadCerts(t *testing.T) {
	for i := 0; i < 3; i++ {
		files := []string{ServerCert, ServerKey, CACert}
		files[i] = FakeFile

		result, err := ConfigureServerTLS(files[0], files[1], true, files[2])
		assert.Nil(t, result)
		assert.Error(t, err)
	}
}

// TestConfigServerTLSServerCertsOnly returns a valid tls config with the
// provided server certificate, and since clientAuth was false, no client auth
// or CAs configured.
func TestConfigServerTLSServerCertsOnly(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, false, "")
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
	assert.Nil(t, tlsConfig.ClientCAs)
}

// TestConfigServerTLSNoCACertsIfNoClientAuth returns a valid tls config with
// the provided server certificate, and since clientAuth was false, no client
// auth or CAs configured even though a client CA cert was provided.
func TestConfigServerTLSNoCACertsIfNoClientAuth(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, false, CACert)
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
	assert.Nil(t, tlsConfig.ClientCAs)
}

// TestTLSConfigClientAuthEnabledNoCACerts returns a valid tls config with the
// provided server certificate client auth enabled, but no CAs configured.
func TestTLSConfigClientAuthEnabledNoCACerts(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, true, "")
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
	assert.Nil(t, tlsConfig.ClientCAs)
}

// TestTLSConfigClientAuthEnabledWithCACert returns a valid tls config with the
// provided server certificate, client auth enabled, and a client CA.
func TestTLSConfigClientAuthEnabledWithCACert(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, true, CACert)
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
	assert.Equal(t, 1, len(tlsConfig.ClientCAs.Subjects()))
}
