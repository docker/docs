package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/docker/libtrust"
)

func generateTrustCA() (libtrust.PrivateKey, *x509.Certificate) {
	key, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		panic(err)
	}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: "CA Root",
		},
		NotBefore:             time.Now().Add(-time.Second),
		NotAfter:              time.Now().Add(24 * 7 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader, cert, cert,
		key.CryptoPublicKey(), key.CryptoPrivateKey(),
	)
	if err != nil {
		panic(err)
	}

	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		panic(err)
	}

	return key, cert
}

func generateIntermediate(key libtrust.PublicKey, parentKey libtrust.PrivateKey, parent *x509.Certificate) *x509.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: "Intermediate",
		},
		NotBefore:             time.Now().Add(-time.Second),
		NotAfter:              time.Now().Add(24 * 7 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader, cert, parent,
		key.CryptoPublicKey(), parentKey.CryptoPrivateKey(),
	)
	if err != nil {
		panic(err)
	}

	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		panic(err)
	}

	return cert
}

func generateTrustCert(key libtrust.PublicKey, parentKey libtrust.PrivateKey, parent *x509.Certificate) *x509.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: "Trust Cert",
		},
		NotBefore:             time.Now().Add(-time.Second),
		NotAfter:              time.Now().Add(24 * 7 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader, cert, parent,
		key.CryptoPublicKey(), parentKey.CryptoPrivateKey(),
	)
	if err != nil {
		panic(err)
	}

	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		panic(err)
	}

	return cert
}

func generateTrustChain(key libtrust.PrivateKey, ca *x509.Certificate) (libtrust.PrivateKey, []*x509.Certificate) {
	parent := ca
	parentKey := key
	chain := make([]*x509.Certificate, 6)
	for i := 5; i > 0; i-- {
		intermediatekey, err := libtrust.GenerateECP256PrivateKey()
		if err != nil {
			panic(err)
		}
		chain[i] = generateIntermediate(intermediatekey, parentKey, parent)
		parent = chain[i]
		parentKey = intermediatekey
	}
	trustKey, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		panic(err)
	}
	chain[0] = generateTrustCert(trustKey, parentKey, parent)

	return trustKey, chain
}

func main() {
	caKey, caCert := generateTrustCA()
	key, chain := generateTrustChain(caKey, caCert)

	caf, err := os.OpenFile("ca.pem", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Error opening ca.pem: %s", err)
	}
	defer caf.Close()
	pem.Encode(caf, &pem.Block{Type: "CERTIFICATE", Bytes: caCert.Raw})

	chainf, err := os.OpenFile("chain.pem", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Error opening ca.pem: %s", err)
	}
	defer chainf.Close()
	for _, c := range chain[1:] {
		pem.Encode(chainf, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
	}

	certf, err := os.OpenFile("cert.pem", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Error opening ca.pem: %s", err)
	}
	defer certf.Close()
	pem.Encode(certf, &pem.Block{Type: "CERTIFICATE", Bytes: chain[0].Raw})

	err = libtrust.SaveKey("key.pem", key)
	if err != nil {
		log.Fatalf("Error saving key: %s", err)
	}

}
