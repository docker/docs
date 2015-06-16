package trustmanager

import "crypto/x509"

const certExtension string = ".crt"

// X509Store is the interface for all X509Stores
type X509Store interface {
	AddCert(cert *x509.Certificate) error
	AddCertFromPEM(pemCerts []byte) error
	AddCertFromFile(filename string) error
	RemoveCert(cert *x509.Certificate) error
	GetCertificateBySKID(hexSKID string) (*x509.Certificate, error)
	GetCertificates() []*x509.Certificate
	GetCertificatePool() *x509.CertPool
	GetVerifyOptions(dnsName string) (x509.VerifyOptions, error)
}

type ID string

// Validator is a convenience type to create validating function that filters
// certificates that get added to the store
type Validator interface {
	Validate(cert *x509.Certificate) bool
}

// ValidatorFunc is a convenience type to create functions that implement
// the Validator interface
type ValidatorFunc func(cert *x509.Certificate) bool

// Validate implements the Validator interface to allow for any func() bool method
// to be passed as a Validator
func (vf ValidatorFunc) Validate(cert *x509.Certificate) bool {
	return vf(cert)
}
