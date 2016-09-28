package pki

type (
	PKIClient interface {
		Address() string
		SignCSR(csr *CertificateSigningRequest) (*CertificateResponse, error)
		GetRootCertificate() (string, error)
	}
)
