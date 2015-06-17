package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/vetinari/trustmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var subjectKeyID string

var cmdKeys = &cobra.Command{
	Use:   "keys",
	Short: "Operates on keys.",
	Long:  "operations on signature keys and trusted certificate authorities.",
	Run:   keysList,
}

func init() {
	cmdKeys.AddCommand(cmdKeysTrust)
	cmdKeys.AddCommand(cmdKeysRemove)
	cmdKeys.AddCommand(cmdKeysGenerate)
}

var cmdKeysRemove = &cobra.Command{
	Use:   "remove [ Subject Key ID ]",
	Short: "removes trust from a specific certificate authority or certificate.",
	Long:  "remove trust from a specific certificate authority.",
	Run:   keysRemove,
}

var cmdKeysTrust = &cobra.Command{
	Use:   "trust [ QDN ] [ certificate ]",
	Short: "Trusts a new certificate for a specific QDN.",
	Long:  "Adds a the certificate to the trusted certificate authority list for the specified Qualified Docker Name.",
	Run:   keysTrust,
}

var cmdKeysGenerate = &cobra.Command{
	Use:   "generate [ QDN ]",
	Short: "Generates a new key for a specific QDN.",
	Long:  "generates a new key for a specific QDN. Qualified Docker Name.",
	Run:   keysGenerate,
}

func keysRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a SHA256 SubjectKeyID of the certificate")
	}

	failed := true
	cert, err := caStore.GetCertificateBySKID(args[0])
	if err == nil {
		fmt.Printf("Removing: ")
		printCert(cert)

		err = caStore.RemoveCert(cert)
		if err != nil {
			fatalf("failed to remove certificate for Root KeyStore")
		}
		failed = false
	}

	cert, err = privStore.GetCertificateBySKID(args[0])
	if err == nil {
		fmt.Printf("Removing: ")
		printCert(cert)

		//TODO (diogo): remove associated private key
		err = privStore.RemoveCert(cert)
		if err != nil {
			fatalf("failed to remove certificate for Private KeyStore")
		}
		failed = false
	}
	if failed {
		fatalf("certificate not found in any store")
	}
}

func keysTrust(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("not enough arguments provided")
	}

	qualifiedDN := args[0]
	certLocationStr := args[1]
	// Verify if argument is a valid URL
	url, err := url.Parse(certLocationStr)
	if err == nil && url.Scheme != "" {

		cert, err := trustmanager.GetCertFromURL(certLocationStr)
		if err != nil {
			fatalf("error retreiving certificate from url (%s): %v", certLocationStr, err)
		}
		err = cert.VerifyHostname(qualifiedDN)
		if err != nil {
			fatalf("certificate does not match the Qualified Docker Name: %v", err)
		}
		err = caStore.AddCert(cert)
		if err != nil {
			fatalf("error adding certificate from file: %v", err)
		}
		fmt.Printf("Adding: ")
		printCert(cert)
	} else if _, err := os.Stat(certLocationStr); err == nil {
		if err := caStore.AddCertFromFile(certLocationStr); err != nil {
			fatalf("error adding certificate from file: %v", err)
		}
	} else {
		fatalf("please provide a file location or URL for CA certificate.")
	}
}

func keysList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}

	fmt.Println("# Trusted Root keys: ")
	trustedCAs := caStore.GetCertificates()
	for _, c := range trustedCAs {
		printCert(c)
	}

	fmt.Println("")
	fmt.Println("# Signing keys: ")
	privateCerts := privStore.GetCertificates()
	for _, c := range privateCerts {
		printCert(c)
	}

}

func keysGenerate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}

	// (diogo): Validate QDNs
	qualifiedDN := args[0]

	key, err := generateKey(qualifiedDN)
	if err != nil {
		fatalf("could not generate key: %v", err)
	}

	template := newCertificate(qualifiedDN, qualifiedDN)
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.(crypto.Signer).Public(), key)
	if err != nil {
		fatalf("failed to generate certificate: %s", err)
	}

	certName := filepath.Join(viper.GetString("privDir"), qualifiedDN+".crt")
	certOut, err := os.Create(certName)
	if err != nil {
		fatalf("failed to save certificate: %s", err)
	}

	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		fatalf("failed to save certificate: %s", err)
	}
}

func newCertificate(qualifiedDN, organization string) *x509.Certificate {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365 * 2)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fatalf("failed to generate serial number: %s", err)
	}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   qualifiedDN,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}
}

func generateKey(qualifiedDN string) (crypto.PrivateKey, error) {
	curve := elliptic.P384()
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not generate private key: %v", err)
	}

	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("could not marshal private key: %v", err)
	}

	keyName := filepath.Join(viper.GetString("privDir"), qualifiedDN+".key")
	keyOut, err := os.OpenFile(keyName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not write privatekey: %v", err)
	}
	defer keyOut.Close()

	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	if err != nil {
		return nil, fmt.Errorf("failed to encode key: %v", err)
	}

	return key, nil
}

func printCert(cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	subjectKeyID := trustmanager.FingerprintCert(cert)
	fmt.Printf("Certificate: %s ; Expires in: %v days; SKID: %s\n", printPkix(cert.Subject), math.Floor(timeDifference.Hours()/24), string(subjectKeyID))
}

func printPkix(pkixName pkix.Name) string {
	return fmt.Sprintf("%s - %s", pkixName.CommonName, pkixName.Organization)
}
