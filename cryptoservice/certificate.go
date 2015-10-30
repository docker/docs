package cryptoservice

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
)

// GenerateCertificate generates an X509 Certificate from a template, given a GUN
func GenerateCertificate(rootKey data.PrivateKey, gun string) (*x509.Certificate, error) {

	switch rootKey.(type) {
	case *data.RSAPrivateKey, *data.ECDSAPrivateKey:
		// go doesn't fall through
	default:
		return nil, fmt.Errorf("only bare RSA or ECDSA keys (not x509 variants) are currently supported. Found: %s", rootKey.Algorithm())
	}

	template, err := trustmanager.NewCertificate(gun)
	if err != nil {
		return nil, fmt.Errorf("failed to create the certificate template for: %s (%v)", gun, err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, rootKey.CryptoSigner().Public(), rootKey.CryptoSigner())
	if err != nil {
		return nil, fmt.Errorf("failed to create the certificate for: %s (%v)", gun, err)
	}

	// Encode the new certificate into PEM
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the certificate for key: %s (%v)", gun, err)
	}

	return cert, nil
}
