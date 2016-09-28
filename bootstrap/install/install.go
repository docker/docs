package install

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/shared/certstore"
)

// The bootstrap works in two phases.  The first phase is used to start another
// copy of the bootstrap container in the correct location.
func install(c *cli.Context) error {
	if !bootstrap.IsPhase2() {
		_, err := Phase1(c)
		return err
	} else {
		return phase2(c)
	}
}

// XXX: CA cert is valid for 100 years for now...
var validFor = time.Hour * 24 * 365 * 100

func MakeSignedClientCert(caCertStore, childCertStore *certstore.CertStore, altNames []string) error {
	if exists, err := caCertStore.Exists(); !exists {
		return fmt.Errorf("CA cert doesn't exist at %s", caCertStore.Path)
	} else if err != nil {
		return fmt.Errorf("Can't determine if CA cert exists: %s", err)
	}
	// read out the CA cert and key
	caCert, err := caCertStore.LoadCert()
	if err != nil {
		return fmt.Errorf("Can't load CA cert: %s", err)
	}
	caKey, err := caCertStore.LoadKey()
	if err != nil {
		return fmt.Errorf("Can't load CA key: %s", err)
	}
	// generate private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %s", err)
	}
	// make the cert template
	template, err := makeCertTemplate()
	if err != nil {
		return fmt.Errorf("Failed to make cert template: %s", err)
	}
	template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	template.DNSNames = altNames

	// sign it
	derBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &priv.PublicKey, caKey)
	// save it all
	err = childCertStore.SaveCertBytes(derBytes)
	if err != nil {
		return fmt.Errorf("Failed to save CA certificate: %s", err)
	}
	err = childCertStore.SaveKey(priv)
	if err != nil {
		return fmt.Errorf("Failed to save key: %s", err)
	}
	log.Debugf("Client cert created for %v at %s, signed by CA at %s", altNames, childCertStore.Path, caCertStore.Path)
	return nil
}

var CertTemplate = x509.Certificate{
	Subject: pkix.Name{
		Organization: []string{"dtr-ha"},
	},

	KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	BasicConstraintsValid: true,
}

func makeCertTemplate() (*x509.Certificate, error) {
	template := CertTemplate
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}
	template.SerialNumber = serialNumber
	template.NotBefore = notBefore
	template.NotAfter = notAfter
	return &template, nil
}

// based on https://golang.org/src/crypto/tls/generate_cert.go
func bootstrapCA(certStore *certstore.CertStore, cn string) error {
	err := certStore.Mkdirp()
	if err != nil {
		return fmt.Errorf("Failed to create dir for CA: %s", err)
	}
	if exists, err := certStore.Exists(); exists {
		log.Debug("CA exists")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to determine if CA cert exists: %s", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %s", err)
	}

	template, err := makeCertTemplate()
	if err != nil {
		return fmt.Errorf("Failed to make cert template: %s", err)
	}
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	template.Subject.CommonName = cn
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("Failed to create certificate: %s", err)
	}
	err = certStore.SaveCertBytes(derBytes)
	if err != nil {
		return fmt.Errorf("Failed to save CA certificate: %s", err)
	}
	err = certStore.SaveKey(priv)
	if err != nil {
		return fmt.Errorf("Failed to save key: %s", err)
	}
	log.Debugf("CA created for %s at %s", cn, certStore.Path)
	return nil
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	err := install(c)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
