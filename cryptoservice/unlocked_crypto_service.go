package cryptoservice

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
)

// GenerateCertificate generates an X509 Certificate from a template, given a GUN
func GenerateCertificate(rootKey data.PrivateKey, gun string) (*x509.Certificate, error) {

	algorithm := rootKey.Algorithm()
	var (
		publicKey  crypto.PublicKey
		privateKey crypto.PrivateKey
		err        error
	)
	switch algorithm {
	case data.RSAKey:
		var rsaPrivateKey *rsa.PrivateKey
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(rootKey.Private())
		privateKey = rsaPrivateKey
		publicKey = rsaPrivateKey.Public()
	case data.ECDSAKey:
		var ecdsaPrivateKey *ecdsa.PrivateKey
		ecdsaPrivateKey, err = x509.ParseECPrivateKey(rootKey.Private())
		privateKey = ecdsaPrivateKey
		publicKey = ecdsaPrivateKey.Public()
	default:
		return nil, fmt.Errorf("only RSA or ECDSA keys are currently supported. Found: %s", algorithm)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse root key: %s (%v)", gun, err)
	}

	template, err := trustmanager.NewCertificate(gun)
	if err != nil {
		return nil, fmt.Errorf("failed to create the certificate template for: %s (%v)", gun, err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
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
