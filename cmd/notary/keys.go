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
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/notary/trustmanager"

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
	Short: "Removes trust from a specific certificate authority or certificate.",
	Long:  "remove trust from a specific certificate authority.",
	Run:   keysRemove,
}

var cmdKeysTrust = &cobra.Command{
	Use:   "trust [ certificate ]",
	Short: "Trusts a new certificate.",
	Long:  "adds a the certificate to the trusted certificate authority list.",
	Run:   keysTrust,
}

var cmdKeysGenerate = &cobra.Command{
	Use:   "generate [ GUN ]",
	Short: "Generates a new key for a specific GUN.",
	Long:  "generates a new key for a specific Global Unique Name.",
	Run:   keysGenerate,
}

// keysRemove deletes Certificates based on hash and Private Keys
// based on GUNs.
func keysRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a SHA256 SubjectKeyID of the certificate")
	}

	gunOrID := args[0]

	cert, err := caStore.GetCertificateBykID(gunOrID)
	if err == nil {
		fmt.Printf("Removing: ")
		printCert(cert)

		err = caStore.RemoveCert(cert)
		if err != nil {
			fatalf("failed to remove certificate from KeyStore")
		}
		return
	}

	// Ask for confirmation before adding certificate into repository
	fmt.Printf("Are you sure you want to remove all keys under this Global Unique Name: %s? (yes/no)\n", gunOrID)
	confirmed := askConfirm()
	if !confirmed {
		fatalf("aborting action.")
	}

	err = privKeyStore.RemoveGUN(gunOrID)
	if err != nil {
		fatalf("failed to remove all Private keys under Global Unique Name: %s", gunOrID)
	}
	fmt.Printf("Removing all Private keys form: %s \n", gunOrID)
}

//TODO (diogo): Ask the use if she wants to trust the GUN in the cert
func keysTrust(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("please provide a URL or filename to a certificate")
	}

	certLocationStr := args[0]
	var cert *x509.Certificate

	// Verify if argument is a valid URL
	url, err := url.Parse(certLocationStr)
	if err == nil && url.Scheme != "" {
		cert, err = trustmanager.GetCertFromURL(certLocationStr)
		if err != nil {
			fatalf("error retrieving certificate from url (%s): %v", certLocationStr, err)
		}
	} else if _, err := os.Stat(certLocationStr); err == nil {
		// Try to load the certificate from the file
		cert, err = trustmanager.LoadCertFromFile(certLocationStr)
		if err != nil {
			fatalf("error adding certificate from file: %v", err)
		}
	} else {
		fatalf("please provide a file location or URL for CA certificate.")
	}

	// Ask for confirmation before adding certificate into repository
	fmt.Printf("Are you sure you want to add trust for: %s? (yes/no)\n", cert.Subject.CommonName)
	confirmed := askConfirm()
	if !confirmed {
		fatalf("aborting action.")
	}

	err = caStore.AddCert(cert)
	if err != nil {
		fatalf("error adding certificate from file: %v", err)
	}
	fmt.Printf("Adding: ")
	printCert(cert)

}

func keysList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}

	fmt.Println("# Trusted Certificates:")
	trustedCAs := caStore.GetCertificates()
	for _, c := range trustedCAs {
		printCert(c)
	}

	fmt.Println("")
	fmt.Println("# Signing keys: ")
	for _, k := range privKeyStore.List() {
		k = strings.TrimSuffix(k, filepath.Ext(k))
		k = strings.TrimPrefix(k, viper.GetString("privDir"))

		fingerprint := filepath.Base(k)
		gun := filepath.Dir(k)[1:]
		fmt.Printf("%s %s\n", gun, fingerprint)
	}
}

func keysGenerate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a GUN")
	}

	//TODO (diogo): Validate GUNs. Don't allow '/' or '\' for now.
	gun := args[0]
	if gun[0:1] == "/" || gun[0:1] == "\\" {
		fatalf("invalid Global Unique Name: %s", gun)
	}

	_, cert, err := generateKeyAndCert(gun)
	if err != nil {
		fatalf("could not generate key: %v", err)
	}

	caStore.AddCert(cert)
	fingerprint := trustmanager.FingerprintCert(cert)
	fmt.Println("Generated new keypair with ID: ", string(fingerprint))
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
	fmt.Printf("%s %s (expires in: %v days)\n", cert.Subject.CommonName, string(subjectKeyID), math.Floor(timeDifference.Hours()/24))
}

func askConfirm() bool {
	var res string
	_, err := fmt.Scanln(&res)
	if err != nil {
		return false
	}
	if strings.EqualFold(res, "y") || strings.EqualFold(res, "yes") {
		return true
	}
	return false
}
