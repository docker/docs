package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	"github.com/docker/dhe-deploy/hubconfig"

	"github.com/docker/libtrust"
)

type CertificateGenerator interface {
	Generate() error
}

type FSCertificateGenerator struct {
	SettingsStore hubconfig.SettingsStore
}

func createRootCert(key libtrust.PrivateKey, duration time.Duration) (*x509.Certificate, error) {
	keyID := key.KeyID()
	now := time.Now()

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{CommonName: keyID},
		NotBefore:             now,
		NotAfter:              now.Add(duration),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		IsCA:         true,
		MaxPathLen:   1,
		SubjectKeyId: []byte(keyID),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.CryptoPublicKey(), key.CryptoPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("unable to create root certificate: %s", err.Error())
	}

	rootCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse root certificate: %s", err.Error())
	}

	return rootCert, nil
}

func (g *FSCertificateGenerator) Generate() error {
	// TODO: the enzi signing keys should be per node in the long run
	enziPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("unable to generate 2048-bit RSA key for enzi: %s", err)
	}

	cryptoPrivateKey, err := libtrust.FromCryptoPrivateKey(enziPrivateKey)
	if err != nil {
		return fmt.Errorf("unable to parse crypto private key: %s", err)
	}

	err = g.SettingsStore.SetEnziSigningKey(cryptoPrivateKey)
	if err != nil {
		return err
	}

	// TODO: we probably shouldn't be generating new certs on every start

	rootKey, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return fmt.Errorf("unable to generate root key: %s", err)
	}

	rootCert, err := createRootCert(rootKey, time.Hour*24*365*3) // Valid for 3 years.
	if err != nil {
		return fmt.Errorf("unable to generate root cert: %s", err)
	}

	err = g.SettingsStore.SetGarantRootCert(rootCert)
	if err != nil {
		return err
	}

	return g.SettingsStore.SetGarantSigningKey(rootKey)
}
