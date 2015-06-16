package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math"
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
	Run:   nil,
}

func init() {
	cmdKeys.AddCommand(cmdKeysTrust)
	cmdKeys.AddCommand(cmdKeysList)
	cmdKeys.AddCommand(cmdKeysRemove)
}

var cmdKeysList = &cobra.Command{
	Use:   "list",
	Short: "List the currently trusted certificate authorities.",
	Long:  "lists the currently trusted certificate authorities.",
	Run:   keysList,
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

func keysRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a SHA256 SubjectKeyID of the certificate")
	}

	cert, err := caStore.GetCertificateBySKID(args[0])
	if err != nil {
		fatalf("certificate not found")
	}

	fmt.Printf("Removing: ")
	printCert(cert)

	err = caStore.RemoveCert(cert)
	if err != nil {
		fatalf("failed to remove certificate for Key Store")
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
	// Load all the certificates
	trustedCAs := caStore.GetCertificates()

	for _, c := range trustedCAs {
		printCert(c)
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
