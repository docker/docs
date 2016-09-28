package pki

import (
	"errors"
)

type DummyClient struct {
}

func NewDummyClient() (*DummyClient, error) {
	// Load up the root CA for this server at the beginning
	return &DummyClient{}, nil
}

func (c *DummyClient) Address() string {
	return ""
}

func (c *DummyClient) SignCSR(csr *CertificateSigningRequest) (*CertificateResponse, error) {
	return nil, errors.New("UCP has not been configured with a CA, so signing is not supported")
}

func (c *DummyClient) GetRootCertificate() (string, error) {
	return "", nil
}
