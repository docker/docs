// Usage:
//  Generate CA
//    ./generate_cert --cert ca.pem --key ca-key.pem
//  Generate signed certificate
//    ./generate_cert --host 127.0.0.1 --cert cert.pem --key key.pem --ca ca.pem --ca-key ca-key.pem

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	host     = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	certFile = flag.String("cert", "", "Output file for certificate")
	keyFile  = flag.String("key", "", "Output file for key")
	ca       = flag.String("ca", "", "Certificate authority file to sign with")
	caKey    = flag.String("ca-key", "", "Certificate authority key file to sign with")
)

const (
	RSABITS  = 2048
	VALIDFOR = 1080 * 24 * time.Hour
	ORG      = "Boot2Docker"
)

func main() {
	flag.Parse()

	if *certFile == "" {
		log.Fatalf("Missing required parameter: --cert")
	}

	if *keyFile == "" {
		log.Fatalf("Missing required parameter: --key")
	}

	if *ca == "" {
		if *caKey != "" {
			log.Fatalf("Must provide both --ca and --ca-key")
		}
		if err := GenerateCA(*certFile, *keyFile); err != nil {
			log.Fatalf("Failured to generate CA: %s", err)
		}
	} else {
		if err := GenerateCert(strings.Split(*host, ","), *certFile, *keyFile, *ca, *caKey); err != nil {
			log.Fatalf("Failured to generate cert: %s", err)
		}
	}
}

// newCertificate creates a new template
func newCertificate() *x509.Certificate {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 1080)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{ORG},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
}

// GenerateCA generates a new certificate authority
// and stores the resulting certificate and key file
// in the arguments.
func GenerateCA(certFile, keyFile string) error {
	template := newCertificate()
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	priv, err := rsa.GenerateKey(rand.Reader, RSABITS)
	if err != nil {
		return err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return nil
}

// GenerateCert generates a new certificate signed using the provided
// certificate authority files and stores the result in the certificate
// file and key provided.  The provided host names are set to the
// appropriate certificate fields.
func GenerateCert(hosts []string, certFile, keyFile, caFile, caKeyFile string) error {
	template := newCertificate()
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	tlsCert, err := tls.LoadX509KeyPair(caFile, caKeyFile)
	if err != nil {
		return err
	}

	priv, err := rsa.GenerateKey(rand.Reader, RSABITS)
	if err != nil {
		return err
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, x509Cert, &priv.PublicKey, tlsCert.PrivateKey)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return nil
}
