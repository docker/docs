package main

import (
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/trustmanager"

	"github.com/spf13/cobra"
)

func init() {
	cmdCert.AddCommand(cmdCertList)

	cmdCertRemove.Flags().StringVarP(&certRemoveGUN, "gun", "g", "", "Globally unique name to delete certificates for")
	cmdCertRemove.Flags().BoolVarP(&certRemoveYes, "yes", "y", false, "Answer yes to the removal question (no confirmation)")
	cmdCert.AddCommand(cmdCertRemove)
}

var cmdCert = &cobra.Command{
	Use:   "cert",
	Short: "Operates on certificates.",
	Long:  `operations on certificates.`,
}

var cmdCertList = &cobra.Command{
	Use:   "list",
	Short: "Lists certificates.",
	Long:  "lists root certificates known to notary.",
	Run:   certList,
}

var certRemoveGUN string
var certRemoveYes bool

var cmdCertRemove = &cobra.Command{
	Use:   "remove [ certID ]",
	Short: "Removes the certificate with the given cert ID.",
	Long:  "remove the certificate with the given cert ID from the local host.",
	Run:   certRemove,
}

// certRemove deletes a certificate given a cert ID or a gun
func certRemove(cmd *cobra.Command, args []string) {
	// If the user hasn't provided -g with a gun, or a cert ID, show usage
	// If the user provided -g and a cert ID, also show usage
	if (len(args) < 1 && certRemoveGUN == "") || (len(args) > 0 && certRemoveGUN != "") {
		cmd.Usage()
		fatalf("must specify the cert ID or the GUN of the certificates to remove")
	}

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	var certsToRemove []*x509.Certificate

	// If there is no GUN, we expect a cert ID
	if certRemoveGUN == "" {
		certID := args[0]
		// This is an invalid ID
		if len(certID) != idSize {
			fatalf("invalid certificate ID provided: %s", certID)
		}
		// Attempt to find this certificates
		cert, err := keyStoreManager.TrustedCertificateStore().GetCertificateByCertID(certID)
		if err != nil {
			fatalf("unable to retrieve certificate with cert ID: %s", certID)
		}
		certsToRemove = append(certsToRemove, cert)
	} else {
		// We got the -g flag, it's a GUN
		certs, err := keyStoreManager.TrustedCertificateStore().GetCertificatesByCN(certRemoveGUN)
		if err != nil {
			fatalf("%v", err)
		}
		certsToRemove = append(certsToRemove, certs...)
	}

	// List all the keys about to be removed
	fmt.Printf("The following certificates will be removed:\n\n")
	for _, cert := range certsToRemove {
		// This error can't occur because we're getting certs off of an
		// x509 store that indexes by ID.
		certID, _ := trustmanager.FingerprintCert(cert)
		fmt.Printf("%s - %s\n", cert.Subject.CommonName, certID)
	}
	fmt.Println("\nAre you sure you want to remove these certificates? (yes/no)")

	// Ask for confirmation before removing certificates, unless -y is provided
	if !certRemoveYes {
		confirmed := askConfirm()
		if !confirmed {
			fatalf("aborting action.")
		}
	}

	// Remove all the certs
	for _, cert := range certsToRemove {
		err = keyStoreManager.TrustedCertificateStore().RemoveCert(cert)
		if err != nil {
			fatalf("failed to remove root certificate for %s", cert.Subject.CommonName)
		}
	}
}

func certList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	fmt.Println("")
	fmt.Println("# Trusted Certificates:")
	trustedCerts := keyStoreManager.TrustedCertificateStore().GetCertificates()
	for _, c := range trustedCerts {
		printCert(c)
	}
}

func printCert(cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	certID, err := trustmanager.FingerprintCert(cert)
	if err != nil {
		fatalf("could not fingerprint certificate: %v", err)
	}

	fmt.Printf("%s %s (expires in: %v days)\n", cert.Subject.CommonName, certID, math.Floor(timeDifference.Hours()/24))
}
