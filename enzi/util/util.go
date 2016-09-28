package util

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	rethink "gopkg.in/dancannon/gorethink.v2"
)

// CombineErrors combiens multiple errors into a single multiError.
func CombineErrors(errs ...error) error {
	var errSlice []error
	for _, err := range errs {
		if err == nil {
			// do nothing
		} else if m, ok := err.(multiError); ok {
			errSlice = append(errSlice, m.Errors()...)
		} else {
			errSlice = append(errSlice, err)
		}
	}
	if len(errSlice) == 0 {
		return nil
	}
	return &multiErr{errors: errSlice}
}

type multiError interface {
	Error() string
	Errors() []error
}

type multiErr struct {
	errors []error
}

func (e *multiErr) Error() string {
	switch len(e.errors) {
	case 0:
		return ""
	case 1:
		return e.errors[0].Error()
	default:
		errMsg := "errors:"
		for _, err := range e.errors {
			errMsg += "\n\t" + err.Error()
		}
		return errMsg
	}
}

func (e *multiErr) Errors() []error {
	return e.errors
}

// GetTLSConfig creates a TLS config using the certificates and key in the
// default locations:
//    key:  /tls/key.pem
//    cert: /tls/cert.pem
//    ca:   /tls/ca.pem
// When used as config for a server, the clientAuthentication parameter
// indicates whether to request, require, and/or verify a client's certificate.
func GetTLSConfig(clientAuth tls.ClientAuthType) (*tls.Config, error) {
	tlsCert, err := tls.LoadX509KeyPair("/tls/cert.pem", "/tls/key.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to load key pair: %s", err)
	}

	rootCertsPEM, err := ioutil.ReadFile("/tls/ca.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to read root CA certificates: %s", err)
	}

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(rootCertsPEM) {
		return nil, fmt.Errorf("unable to parse root CA certificates")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      rootCAs,
		ClientAuth:   clientAuth,
		ClientCAs:    rootCAs,
	}, nil
}

// GetDBSession attempts to create a session to the RethinkDB cluster with the
// given addresses. If no addrs are specified, localhost:28015 is used as a
// default. At least one RethinkDB node must be available.
func GetDBSession(addrs []string, tlsConfig *tls.Config) (*rethink.Session, error) {
	return rethink.Connect(rethink.ConnectOpts{
		Addresses: addrs,
		// Only connect to the addrs given. Usually just the local
		// RethinkDB instance. We would enable this but we currently
		// do not handle reconnect/retries when a query fails due to
		// a severed connection with another DB in the cluster. This
		// is usually not a big deal, as a retry of the API call would
		// eventually work but we have a bunch of integration tests
		// that do not have retry logic. A long-term approach might be
		// to put the retry logic in enzi/schema.Manager itself but
		// that's a lot of extra work.
		DiscoverHosts: false,
		MaxIdle:       5,
		MaxOpen:       10,
		TLSConfig:     tlsConfig,
	})
}
