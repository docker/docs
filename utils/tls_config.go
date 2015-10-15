package utils

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"

	"github.com/docker/notary/trustmanager"
)

// ConfigureServerTLS specifies a set of ciphersuites, the server cert and key,
// and optionally client authentication.  Note that a tls configuration is
// constructed that either requires and verifies client authentication or
// doesn't deal with client certs at all. Nothing in the middle.
func ConfigureServerTLS(serverCert string, serverKey string, clientAuth bool, caCert string) (*tls.Config, error) {
	keypair, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{keypair},
		Rand:         rand.Reader,
	}

	if clientAuth {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		if caCert != "" {
			cert, err := trustmanager.LoadCertFromFile(caCert)
			if err != nil {
				return nil, err
			}
			caPool := x509.NewCertPool()
			caPool.AddCert(cert)
			tlsConfig.ClientCAs = caPool
		}
	}

	return tlsConfig, nil
}
