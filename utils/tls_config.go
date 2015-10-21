package utils

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/docker/notary/trustmanager"
)

// Client TLS cipher suites (dropping CBC ciphers for client preferred suite set)
var clientCipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

// Server TLS cipher suites
var serverCipherSuites = append(clientCipherSuites, []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
}...)

// ConfigureServerTLS specifies a set of ciphersuites, the server cert and key,
// and optionally client authentication.  Note that a tls configuration is
// constructed that either requires and verifies client authentication or
// doesn't deal with client certs at all. Nothing in the middle.
func ConfigureServerTLS(serverCert, serverKey string, clientAuth bool, caCertDir string) (*tls.Config, error) {
	keypair, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites:             serverCipherSuites,
		Certificates:             []tls.Certificate{keypair},
		Rand:                     rand.Reader,
	}

	if clientAuth {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if caCertDir != "" {
		// Check to see if the given directory exists
		fi, err := os.Stat(caCertDir)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("No such directory: %s", caCertDir)
		}

		certStore, err := trustmanager.NewX509FileStore(caCertDir)
		if err != nil {
			return nil, err
		}
		if certStore.Empty() {
			return nil, fmt.Errorf("No certificates in %s", caCertDir)
		}
		tlsConfig.ClientCAs = certStore.GetCertificatePool()
	}

	return tlsConfig, nil
}

// ConfigureClientTLS generates a tls configuration for clients using the
// provided parameters.
func ConfigureClientTLS(rootCA, serverName string, insecureSkipVerify bool, clientCert, clientKey string) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
		CipherSuites:       clientCipherSuites,
		ServerName:         serverName,
	}

	if rootCA != "" {
		rootCert, err := trustmanager.LoadCertFromFile(rootCA)
		if err != nil {
			return nil, fmt.Errorf(
				"Could not load root ca file. %s", err.Error())
		}
		rootPool := x509.NewCertPool()
		rootPool.AddCert(rootCert)
		tlsConfig.RootCAs = rootPool
	}

	if clientCert != "" && clientKey != "" {
		keypair, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{keypair}
	}

	return tlsConfig, nil
}
