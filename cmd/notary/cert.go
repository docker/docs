package main

import (
	"crypto/x509"
	"os"
	"path/filepath"

	"github.com/docker/notary"
	notaryclient "github.com/docker/notary/client"
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
	Long:  `Operations on certificates.`,
}

var cmdCertList = &cobra.Command{
	Use:   "list",
	Short: "Lists certificates.",
	Long:  "Lists root certificates known to notary.",
	Run:   certList,
}

var certRemoveGUN string
var certRemoveYes bool

var cmdCertRemove = &cobra.Command{
	Use:   "remove [ certID ]",
	Short: "Removes the certificate with the given cert ID.",
	Long:  "Remove the certificate with the given cert ID from the local host.",
	Run:   certRemove,
}

// certRemove deletes a certificate given a cert ID or a gun
// If given a gun, certRemove will also remove local TUF data
func certRemove(cmd *cobra.Command, args []string) {
	// If the user hasn't provided -g with a gun, or a cert ID, show usage
	// If the user provided -g and a cert ID, also show usage
	if (len(args) < 1 && certRemoveGUN == "") || (len(args) > 0 && certRemoveGUN != "") {
		cmd.Usage()
		fatalf("Must specify the cert ID or the GUN of the certificates to remove")
	}
	parseConfig()

	trustDir := mainViper.GetString("trust_dir")
	certPath := filepath.Join(trustDir, notary.TrustedCertsDir)
	certStore, err := trustmanager.NewX509FilteredFileStore(
		certPath,
		trustmanager.FilterCertsExpiredSha1,
	)
	if err != nil {
		fatalf("Failed to create a new truststore with directory: %s", trustDir)
	}

	var certsToRemove []*x509.Certificate
	var certFoundByID *x509.Certificate
	var removeTrustData bool

	// If there is no GUN, we expect a cert ID
	if certRemoveGUN == "" {
		certID := args[0]
		// Attempt to find this certificate
		certFoundByID, err = certStore.GetCertificateByCertID(certID)
		if err != nil {
			// This is an invalid ID, the user might have forgotten a character
			if len(certID) != notary.Sha256HexSize {
				fatalf("Unable to retrieve certificate with invalid certificate ID provided: %s", certID)
			}
			fatalf("Unable to retrieve certificate with cert ID: %s", certID)
		}
		// the GUN is the CN from the certificate
		certRemoveGUN = certFoundByID.Subject.CommonName
		certsToRemove = []*x509.Certificate{certFoundByID}
	}

	toRemove, err := certStore.GetCertificatesByCN(certRemoveGUN)
	// We could not find any certificates matching the user's query, so propagate the error
	if err != nil {
		fatalf("%v", err)
	}

	// If we specified a GUN or if the ID we specified is the only certificate with its CN, remove all GUN certs and trust data too
	if certFoundByID == nil || len(toRemove) == 1 {
		removeTrustData = true
		certsToRemove = toRemove
	}

	// List all the certificates about to be removed
	cmd.Printf("The following certificates will be removed:\n\n")
	for _, cert := range certsToRemove {
		// This error can't occur because we're getting certs off of an
		// x509 store that indexes by ID.
		certID, _ := trustmanager.FingerprintCert(cert)
		cmd.Printf("%s - %s\n", cert.Subject.CommonName, certID)
	}
	// If we were given a GUN or the last ID for a GUN, inform the user that we'll also delete all TUF data
	if removeTrustData {
		cmd.Printf("\nAll local trust data will be removed for %s\n", certRemoveGUN)
	}
	cmd.Println("\nAre you sure you want to remove these certificates? (yes/no)")

	// Ask for confirmation before removing certificates, unless -y is provided
	if !certRemoveYes {
		confirmed := askConfirm()
		if !confirmed {
			fatalf("Aborting action.")
		}
	}

	if removeTrustData {
		// Remove all TUF data, so call RemoveTrustData on a NotaryRepository with the GUN
		// no online operations are performed so the transport argument is nil
		nRepo, err := notaryclient.NewNotaryRepository(trustDir, certRemoveGUN, getRemoteTrustServer(mainViper), nil, retriever)
		if err != nil {
			fatalf("Could not establish trust data for GUN %s", certRemoveGUN)
		}
		// DeleteTrustData will pick up all of the same certificates by GUN (CN) and remove them
		err = nRepo.DeleteTrustData()
		if err != nil {
			fatalf("Failed to delete trust data for %s", certRemoveGUN)
		}
	} else {
		for _, cert := range certsToRemove {
			err = certStore.RemoveCert(cert)
			if err != nil {
				fatalf("Failed to remove cert %s", cert)
			}
		}
	}
}

func certList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}
	parseConfig()

	trustDir := mainViper.GetString("trust_dir")
	certPath := filepath.Join(trustDir, notary.TrustedCertsDir)
	// Load all individual (non-CA) certificates that aren't expired and don't use SHA1
	certStore, err := trustmanager.NewX509FilteredFileStore(
		certPath,
		trustmanager.FilterCertsExpiredSha1,
	)
	if err != nil {
		fatalf("Failed to create a new truststore with directory: %s", trustDir)
	}

	trustedCerts := certStore.GetCertificates()

	cmd.Println("")
	prettyPrintCerts(trustedCerts, cmd.Out())
	cmd.Println("")
}
