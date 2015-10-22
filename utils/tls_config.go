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

// ServerTLSOpts generates a tls configuration for servers using the
// provided parameters.
type ServerTLSOpts struct {
	ServerCertFile    string
	ServerKeyFile     string
	RequireClientAuth bool
	ClientCADirectory string
}

// ConfigureServerTLS specifies a set of ciphersuites, the server cert and key,
// and optionally client authentication.  Note that a tls configuration is
// constructed that either requires and verifies client authentication or
// doesn't deal with client certs at all. Nothing in the middle.
func ConfigureServerTLS(opts *ServerTLSOpts) (*tls.Config, error) {
	keypair, err := tls.LoadX509KeyPair(
		opts.ServerCertFile, opts.ServerKeyFile)
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

	if opts.RequireClientAuth {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if opts.ClientCADirectory != "" {
		// Check to see if the given directory exists
		fi, err := os.Stat(opts.ClientCADirectory)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("No such directory: %s", opts.ClientCADirectory)
		}

		certStore, err := trustmanager.NewX509FileStore(opts.ClientCADirectory)
		if err != nil {
			return nil, err
		}
		if certStore.Empty() {
			return nil, fmt.Errorf("No certificates in %s", opts.ClientCADirectory)
		}
		tlsConfig.ClientCAs = certStore.GetCertificatePool()
	}

	return tlsConfig, nil
}

// ClientTLSOpts is a struct that contains options to pass to
// ConfigureClientTLS
type ClientTLSOpts struct {
	RootCAFile         string
	ServerName         string
	InsecureSkipVerify bool
	ClientCertFile     string
	ClientKeyFile      string
}

// ConfigureClientTLS generates a tls configuration for clients using the
// provided parameters.
func ConfigureClientTLS(opts *ClientTLSOpts) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: opts.InsecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
		CipherSuites:       clientCipherSuites,
		ServerName:         opts.ServerName,
	}

	if opts.RootCAFile != "" {
		rootCert, err := trustmanager.LoadCertFromFile(opts.RootCAFile)
		if err != nil {
			return nil, fmt.Errorf(
				"Could not load root ca file. %s", err.Error())
		}
		rootPool := x509.NewCertPool()
		rootPool.AddCert(rootCert)
		tlsConfig.RootCAs = rootPool
	}

	if opts.ClientCertFile != "" || opts.ClientKeyFile != "" {
		keypair, err := tls.LoadX509KeyPair(
			opts.ClientCertFile, opts.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{keypair}
	}

	return tlsConfig, nil
}
