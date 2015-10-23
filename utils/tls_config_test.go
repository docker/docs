package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/stretchr/testify/assert"
)

const (
	ServerCert = "../fixtures/notary-server.crt"
	ServerKey  = "../fixtures/notary-server.key"
	RootCA     = "../fixtures/root-ca.crt"
)

// generates a multiple-certificate file with both RSA and ECDSA certs and
// some garbage, returns filename.
func generateMultiCert(t *testing.T) string {
	tempFile, err := ioutil.TempFile("/tmp", "cert-test")
	defer tempFile.Close()
	assert.NoError(t, err)

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	ecKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	assert.NoError(t, err)
	template, err := trustmanager.NewCertificate("gun")
	assert.NoError(t, err)

	for _, key := range []crypto.Signer{rsaKey, ecKey} {
		derBytes, err := x509.CreateCertificate(
			rand.Reader, template, template, key.Public(), key)
		assert.NoError(t, err)

		cert, err := x509.ParseCertificate(derBytes)
		assert.NoError(t, err)

		pemBytes := trustmanager.CertToPEM(cert)
		nBytes, err := tempFile.Write(pemBytes)
		assert.NoError(t, err)
		assert.Equal(t, nBytes, len(pemBytes))

		assert.NoError(t, err)
	}

	_, err = tempFile.WriteString(`\n
    -----BEGIN CERTIFICATE-----
    This is some garbage that isnt a cert
    -----END CERTIFICATE-----
    `)

	return tempFile.Name()
}

// If the cert files and directory are provided but are invalid, an error is
// returned.
func TestConfigServerTLSFailsIfUnableToLoadCerts(t *testing.T) {
	for i := 0; i < 3; i++ {
		files := []string{ServerCert, ServerKey, RootCA}
		files[i] = "not-real-file"

		result, err := ConfigureServerTLS(&ServerTLSOpts{
			ServerCertFile:    files[0],
			ServerKeyFile:     files[1],
			RequireClientAuth: true,
			ClientCAFile:      files[2],
		})
		assert.Nil(t, result)
		assert.Error(t, err)
	}
}

// If server cert and key are provided, and client auth is disabled, then
// a valid tls.Config is returned with ClientAuth set to NoClientCert
func TestConfigServerTLSServerCertsOnly(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(&ServerTLSOpts{
		ServerCertFile: ServerCert,
		ServerKeyFile:  ServerKey,
	})
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
	assert.Nil(t, tlsConfig.ClientCAs)
}

// If a valid client cert file is provided, but it contains no client
// certs, an error is returned.
func TestConfigServerTLSWithEmptyCACertFile(t *testing.T) {
	tempFile, err := ioutil.TempFile("/tmp", "cert-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempFile.Name())
	tempFile.Close()

	tlsConfig, err := ConfigureServerTLS(&ServerTLSOpts{
		ServerCertFile: ServerCert,
		ServerKeyFile:  ServerKey,
		ClientCAFile:   tempFile.Name(),
	})
	assert.Nil(t, tlsConfig)
	assert.Error(t, err)
}

// If server cert and key are provided, and client cert file is provided with
// one cert, a valid tls.Config is returned with the clientCAs set to that
// cert.
func TestConfigServerTLSWithOneCACert(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(&ServerTLSOpts{
		ServerCertFile: ServerCert,
		ServerKeyFile:  ServerKey,
		ClientCAFile:   RootCA,
	})
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
	assert.Len(t, tlsConfig.ClientCAs.Subjects(), 1)
}

// If server cert and key are provided, and client cert file is provided with
// multiple certs (and garbage), a valid tls.Config is returned with the
// clientCAs set to the valid cert and the garbage is ignored (but only
// because the garbage is at the end - actually CertPool.AppendCertsFromPEM
// aborts as soon as it finds an invalid cert)
func TestConfigServerTLSWithMultipleCACertsAndGarbage(t *testing.T) {
	tempFilename := generateMultiCert(t)
	defer os.RemoveAll(tempFilename)

	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(&ServerTLSOpts{
		ServerCertFile: ServerCert,
		ServerKeyFile:  ServerKey,
		ClientCAFile:   tempFilename,
	})
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
	assert.Len(t, tlsConfig.ClientCAs.Subjects(), 2)
}

// If server cert and key are provided, and client auth is disabled, then
// a valid tls.Config is returned with ClientAuth set to
// RequireAndVerifyClientCert
func TestConfigServerTLSClientAuthEnabled(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureServerTLS(&ServerTLSOpts{
		ServerCertFile:    ServerCert,
		ServerKeyFile:     ServerKey,
		RequireClientAuth: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
	assert.Nil(t, tlsConfig.ClientCAs)
}

// The skipVerify boolean gets set on the tls.Config's InsecureSkipBoolean
func TestConfigClientTLSNoVerify(t *testing.T) {
	for _, skip := range []bool{true, false} {
		tlsConfig, err := ConfigureClientTLS(
			&ClientTLSOpts{InsecureSkipVerify: skip})
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
		tlsConfig, err := ConfigureClientTLS(&ClientTLSOpts{ServerName: name})
		assert.NoError(t, err)
		assert.Nil(t, tlsConfig.Certificates)
		assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
		assert.Equal(t, name, tlsConfig.ServerName)
		assert.Nil(t, tlsConfig.RootCAs)
	}
}

// The RootCA is set if the file provided has a single CA cert.
func TestConfigClientTLSRootCAFileWithOneCert(t *testing.T) {
	tlsConfig, err := ConfigureClientTLS(&ClientTLSOpts{RootCAFile: RootCA})
	assert.NoError(t, err)
	assert.Nil(t, tlsConfig.Certificates)
	assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, "", tlsConfig.ServerName)
	assert.Len(t, tlsConfig.RootCAs.Subjects(), 1)
}

// If the root CA file provided has multiple CA certs and garbage, only the
// valid certs are read (but only because the garbage is at the end - actually
// CertPool.AppendCertsFromPEM aborts as soon as it finds an invalid cert)
func TestConfigClientTLSRootCAFileMultipleCertsAndGarbage(t *testing.T) {
	tempFilename := generateMultiCert(t)
	defer os.RemoveAll(tempFilename)

	tlsConfig, err := ConfigureClientTLS(
		&ClientTLSOpts{RootCAFile: tempFilename})
	assert.NoError(t, err)
	assert.Nil(t, tlsConfig.Certificates)
	assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, "", tlsConfig.ServerName)
	assert.Len(t, tlsConfig.RootCAs.Subjects(), 2)
}

// An error is returned if a root CA is provided but the file doesn't exist.
func TestConfigClientTLSNonexistentRootCAFile(t *testing.T) {
	tlsConfig, err := ConfigureClientTLS(
		&ClientTLSOpts{RootCAFile: "not-a-file"})
	assert.Error(t, err)
	assert.Nil(t, tlsConfig)
}

// An error is returned if either the client cert or the key are provided
// but invalid or blank.
func TestConfigClientTLSClientCertOrKeyInvalid(t *testing.T) {
	for i := 0; i < 2; i++ {
		for _, invalid := range []string{"not-a-file", ""} {
			files := []string{ServerCert, ServerKey}
			files[i] = invalid
			tlsConfig, err := ConfigureClientTLS(&ClientTLSOpts{
				ClientCertFile: files[0], ClientKeyFile: files[1]})
			assert.Error(t, err)
			assert.Nil(t, tlsConfig)
		}
	}
}

// The certificate is set if the client cert and client key are provided and
// valid.
func TestConfigClientTLSValidClientCertAndKey(t *testing.T) {
	keypair, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	assert.NoError(t, err)

	tlsConfig, err := ConfigureClientTLS(&ClientTLSOpts{
		ClientCertFile: ServerCert, ClientKeyFile: ServerKey})
	assert.NoError(t, err)
	assert.Equal(t, []tls.Certificate{keypair}, tlsConfig.Certificates)
	assert.Equal(t, false, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, "", tlsConfig.ServerName)
	assert.Nil(t, tlsConfig.RootCAs)
}
