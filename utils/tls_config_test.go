package utils

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ServerCert = "../fixtures/notary-server.crt"
	ServerKey  = "../fixtures/notary-server.key"
	RootCA     = "../fixtures/root-ca.crt"
)

// copies the provided certificate into a temporary directory
func makeTempCertDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("/tmp", "cert-test")
	assert.NoError(t, err, "couldn't open temp directory")

	in, err := os.Open(RootCA)
	assert.NoError(t, err, "cannot open %s", RootCA)
	defer in.Close()
	copiedCert := filepath.Join(tempDir, filepath.Base(RootCA))
	out, err := os.Create(copiedCert)
	assert.NoError(t, err, "cannot open %s", copiedCert)
	defer out.Close()
	_, err = io.Copy(out, in)
	assert.NoError(t, err, "cannot copy %s to %s", RootCA, copiedCert)

	return tempDir
}

// If the cert files and directory are provided but are invalid, an error is
// returned.
func TestConfigServerTLSFailsIfUnableToLoadCerts(t *testing.T) {
	tempDir := makeTempCertDir(t)

	for i := 0; i < 3; i++ {
		files := []string{ServerCert, ServerKey, tempDir}
		files[i] = "not-real-file"

		result, err := ConfigureServerTLS(files[0], files[1], true, files[2])
		assert.Nil(t, result)
		assert.Error(t, err)
	}
}

// If server cert and key are provided, and client auth is disabled, then
// a valid tls.Config is returned with ClientAuth set to NoClientCert
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

// If a valid client cert directory is provided, but it contains no client
// certs, an error is returned.
func TestConfigServerTLSWithEmptyCACertDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("/tmp", "cert-test")
	assert.NoError(t, err, "couldn't open temp directory")

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, false, tempDir)
	assert.Nil(t, tlsConfig)
	assert.Error(t, err)
}

// If server cert and key are provided, and client cert directory is provided,
// a valid tls.Config is returned with the clientCAs set to the certs in that
// directory.
func TestConfigServerTLSWithCACerts(t *testing.T) {
	tempDir := makeTempCertDir(t)
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, false, tempDir)
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
	assert.Len(t, tlsConfig.ClientCAs.Subjects(), 1)
}

// If server cert and key are provided, and client auth is disabled, then
// a valid tls.Config is returned with ClientAuth set to
// RequireAndVerifyClientCert
func TestConfigServerTLSClientAuthEnabled(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(ServerCert, ServerKey, true, "")
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
	assert.Nil(t, tlsConfig.ClientCAs)
}

// The skipVerify boolean gets set on the tls.Config's InsecureSkipBoolean
func TestConfigClientTLSNoVerify(t *testing.T) {
	for _, skip := range []bool{true, false} {
		tlsConfig, err := ConfigureClientTLS("", "", skip, "", "")
		assert.NoError(t, err)
		assert.Nil(t, tlsConfig.Certificates)
		assert.Equal(t, skip, tlsConfig.InsecureSkipVerify)
		assert.Equal(t, "", tlsConfig.ServerName)
		assert.Nil(t, tlsConfig.RootCAs)
	}
}

// The skipVerify boolean gets set on the tls.Config's InsecureSkipBoolean
func TestConfigClientServerName(t *testing.T) {
	for _, name := range []string{"", "myname"} {
		tlsConfig, err := ConfigureClientTLS("", name, false, "", "")
		assert.NoError(t, err)
		assert.Nil(t, tlsConfig.Certificates)
		assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
		assert.Equal(t, name, tlsConfig.ServerName)
		assert.Nil(t, tlsConfig.RootCAs)
	}
}

// The RootCA is set if it is provided and valid
func TestConfigClientTLSValidRootCA(t *testing.T) {
	tlsConfig, err := ConfigureClientTLS(RootCA, "", false, "", "")
	assert.NoError(t, err)
	assert.Nil(t, tlsConfig.Certificates)
	assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, "", tlsConfig.ServerName)
	assert.Len(t, tlsConfig.RootCAs.Subjects(), 1)
}

// An error is returned if a root CA is provided but not valid
func TestConfigClientTLSInValidRootCA(t *testing.T) {
	tlsConfig, err := ConfigureClientTLS("not-a-file.crt", "", false, "", "")
	assert.Error(t, err)
	assert.Nil(t, tlsConfig)
}

// An error is returned if either the client cert or the key are provided
// but invalid.
func TestConfigClientTLSClientCertOrKeyInvalid(t *testing.T) {
	for i := 0; i < 2; i++ {
		files := []string{ServerCert, ServerKey}
		files[i] = "not-a-file.crt"
		tlsConfig, err := ConfigureClientTLS("", "", false, files[0], files[1])
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
	}
}

// The certificate is set if the client cert and client key are provided and
// valid.
func TestConfigClientTLSValidClientCertAndKey(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureClientTLS("", "", false, ServerCert, ServerKey)
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, "", tlsConfig.ServerName)
	assert.Nil(t, tlsConfig.RootCAs)
}
