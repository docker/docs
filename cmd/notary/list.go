package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math"
	"time"

	"github.com/docker/vetinari/trustmanager"
	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the currently trusted certificate authorities.",
	Long:  "lists the currently trusted certificate authorities.",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	// Load all the certificates
	trustedCAs := caStore.GetCertificates()

	fmt.Println("CAs Loaded:")
	for _, c := range trustedCAs {
		print_cert(c)
	}

}

func print_cert(cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	subjectKeyID := trustmanager.FingerprintCert(cert)
	fmt.Printf("Certificate: %s ; Expires in: %v days; SKID: %s\n", printPkix(cert.Subject), math.Floor(timeDifference.Hours()/24), string(subjectKeyID))
}

func printPkix(pkixName pkix.Name) string {
	return fmt.Sprintf("%s - %s", pkixName.CommonName, pkixName.Organization)
}
