package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"os"
	"time"

	"github.com/docker/vetinari/trustmanager"

	"github.com/spf13/cobra"
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
	Use:   "trust [ certificate ]",
	Short: "Trusts a new certificate for a specific GUN.",
	Long:  "Adds a the certificate to the trusted certificate authority list for the specified Global Unique Name.",
	Run:   keysTrust,
}

var cmdKeysGenerate = &cobra.Command{
	Use:   "generate [ GUN ]",
	Short: "Generates a new key for a specific GUN.",
	Long:  "generates a new key for a specific GUN. Global Unique Name.",
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

	gun := args[0]
	certLocationStr := args[1]
	// Verify if argument is a valid URL
	url, err := url.Parse(certLocationStr)
	if err == nil && url.Scheme != "" {

		cert, err := trustmanager.GetCertFromURL(certLocationStr)
		if err != nil {
			fatalf("error retreiving certificate from url (%s): %v", certLocationStr, err)
		}
		err = cert.VerifyHostname(gun)
		if err != nil {
			fatalf("certificate does not match the Global Unique Name: %v", err)
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
		fatalf("must specify a GUN")
	}

	// (diogo): Validate GUNs
	gun := args[0]

	_, cert, err := generateKeyAndCert(gun)
	if err != nil {
		fatalf("could not generate key: %v", err)
	}

	caStore.AddCertFromPEM(cert)
}

func newCertificate(gun, organization string) *x509.Certificate {
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
			CommonName:   gun,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}
}

func printCert(cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	subjectKeyID := trustmanager.FingerprintCert(cert)
	fmt.Printf("Certificate: %s ; Expires in: %v days; SKID: %s\n", printPkix(cert.Subject), math.Floor(timeDifference.Hours()/24), string(subjectKeyID))
}

func printPkix(pkixName pkix.Name) string {
	return fmt.Sprintf("%s - %s", pkixName.CommonName, pkixName.Organization)
}
