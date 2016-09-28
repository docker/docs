package enzi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	libkv "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/enzi/api/client/openid"
	"github.com/docker/orca/enzi/jose"
)

// The API endpoint on the provider that issues openid tokens.
const TokenAPIEndpointPath = "/enzi/v0/id/token"

type Authenticator struct {
	kvStore      libkv.Store
	httpClient   *http.Client
	openidClient *openid.Client

	defaultOrgID    string
	providerAddr    string
	userDefaultRole auth.Role
	providerCAs     *x509.CertPool
	signingKey      *jose.PrivateKey
}

var _ auth.Authenticator = (*Authenticator)(nil)

func NewAuthenticator(config auth.EnziConfig, hostAddr string, kvStore libkv.Store, signingKey *jose.PrivateKey, providerRootCAs string, signingKeyCertChain ...string) (auth.Authenticator, error) {
	providerCAs := x509.NewCertPool()
	if !providerCAs.AppendCertsFromPEM([]byte(providerRootCAs)) {
		return nil, fmt.Errorf("unable to parse any certificates from provider CA bundle")
	}

	tlsConfig := &tls.Config{
		RootCAs: providerCAs,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			// The default is 2 which is too small. We may need to
			// adjust this value as we get results from load/stress
			// tests.
			MaxIdleConnsPerHost: 5,
		},
	}

	providerAddr, err := getProviderAddr(hostAddr, config.ProviderAddrs)
	if err != nil {
		return nil, fmt.Errorf("unable to get auth provider address: %s", err)
	}

	openidClient := openid.NewClient(
		httpClient,
		signingKey,
		config.ServiceID,
		"redirectURI", // NOT REQUIRED: We do not use the authorizaiton_code grant type.
		providerAddr,
		"authorizationPath", // NOT REQUIRED: We do not use the authorizaiton_code grant type.
		TokenAPIEndpointPath,
		signingKeyCertChain...,
	)

	authenticator := &Authenticator{
		kvStore:         kvStore,
		httpClient:      httpClient,
		openidClient:    openidClient,
		defaultOrgID:    config.DefaultOrgID,
		providerAddr:    providerAddr,
		userDefaultRole: config.UserDefaultRole,
		providerCAs:     providerCAs,
		signingKey:      signingKey,
	}

	// Save the signing key in the KV Store.
	if err := authenticator.SaveSigningKey(); err != nil {
		return nil, fmt.Errorf("unable to save signing key in KV Store: %s", err)
	}

	return authenticator, nil
}

func getProviderAddr(hostAddr string, providerAddrs []string) (string, error) {
	if len(providerAddrs) == 0 {
		return "", fmt.Errorf("no auth provider addrs configured")
	}

	for _, providerAddr := range providerAddrs {
		host, _, err := net.SplitHostPort(providerAddr)
		if err != nil {
			if addrErr, ok := err.(*net.AddrError); ok && addrErr.Err == "missing port in address" {
				host = providerAddr // No port specified.
			} else {
				// Ignore the error and go to the next provider
				// address.
				log.Warnf("unable to parse provider address host/port: %s", err)
				continue
			}
		}

		if host == hostAddr {
			// We've found the provider API server on this host.
			return providerAddr, nil
		}
	}

	// Didn't find our host addr in the list. Fallback to using the first
	// entry in the list.
	log.Infof("unable to find local host addr in list of auth provider addrs (defaulting to first entry): %s", providerAddrs)
	return providerAddrs[0], nil
}

func (a *Authenticator) ProviderAddr() string {
	return a.providerAddr
}

func (a *Authenticator) ProviderCAs() *x509.CertPool {
	return a.providerCAs
}

func (a *Authenticator) Name() string {
	return "eNZi"
}
