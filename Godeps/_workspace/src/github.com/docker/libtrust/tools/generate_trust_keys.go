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

func writeCertFile(filename string, certs []*x509.Certificate) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, cert := range certs {
		err = pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
		if err != nil {
			return err
		}
	}

	return nil
}

func generateIntermediate(parent *x509.Certificate, key libtrust.PublicKey, parentKey libtrust.PrivateKey) (*x509.Certificate, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: key.KeyID(),
		},
		NotBefore:             time.Now().Add(-time.Second),
		NotAfter:              time.Now().Add(90 * 24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, cert, parent, key.CryptoPublicKey(), parentKey.CryptoPrivateKey())
	if err != nil {
		return nil, err
	}

	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func generateLeaf(parent *x509.Certificate, key libtrust.PublicKey, parentKey libtrust.PrivateKey) (*x509.Certificate, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: key.KeyID(),
		},
		NotBefore:             time.Now().Add(-time.Second),
		NotAfter:              time.Now().Add(90 * 24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, cert, parent, key.CryptoPublicKey(), parentKey.CryptoPrivateKey())
	if err != nil {
		return nil, err
	}

	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func main() {
	caKey, _ := libtrust.GenerateECP256PrivateKey()
	masterKey, _ := libtrust.GenerateECP256PrivateKey()
	trustKey, _ := libtrust.GenerateECP256PrivateKey()

	log.Printf("Generated keys:\n\tCA:     %s\n\tMaster: %s\n\tTrust:  %s", caKey.KeyID(), masterKey.KeyID(), trustKey.KeyID())

	libtrust.SaveKey("ca-key.json", caKey)
	libtrust.SaveKey("master-key.json", masterKey)
	libtrust.SaveKey("trust-key.json", trustKey)

	// TODO better CA function
	ca, err := libtrust.GenerateCACert(caKey, caKey.PublicKey())
	if err != nil {
		log.Fatalf("Error generating CA: %s", err)
	}

	err = writeCertFile("ca.pem", []*x509.Certificate{ca})
	if err != nil {
		log.Fatalf("Error writing CA pem file: %s", err)
	}

	masterCert, err := generateIntermediate(ca, masterKey.PublicKey(), caKey)
	if err != nil {
		log.Fatalf("Error generating master certificate: %s", err)
	}
	// Generate Master Server certificate, signed by CA
	// Output master-key.json

	leafCert, err := generateLeaf(masterCert, trustKey.PublicKey(), masterKey)
	if err != nil {
		log.Fatalf("Error generating leaf certificate: %s", err)
	}
	// Generate key, from key trust Server certificate, signed by master
	// Output cert.pem (both trust server and master certificate), key.json

	err = writeCertFile("cert.pem", []*x509.Certificate{leafCert, masterCert})
	if err != nil {
		log.Fatalf("Error generating cert pem file")
	}

}
