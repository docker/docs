package main

import (
	"crypto/x509"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"time"

	"github.com/docker/notary/certs"
	"github.com/docker/notary/trustmanager"
	"github.com/olekukonko/tablewriter"

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
func certRemove(cmd *cobra.Command, args []string) {
	// If the user hasn't provided -g with a gun, or a cert ID, show usage
	// If the user provided -g and a cert ID, also show usage
	if (len(args) < 1 && certRemoveGUN == "") || (len(args) > 0 && certRemoveGUN != "") {
		cmd.Usage()
		fatalf("Must specify the cert ID or the GUN of the certificates to remove")
	}
	parseConfig()

	trustDir := mainViper.GetString("trust_dir")
	certManager, err := certs.NewManager(trustDir)
	if err != nil {
		fatalf("Failed to create a new truststore manager with directory: %s", trustDir)
	}

	var certsToRemove []*x509.Certificate

	// If there is no GUN, we expect a cert ID
	if certRemoveGUN == "" {
		certID := args[0]
		// This is an invalid ID
		if len(certID) != idSize {
			fatalf("Invalid certificate ID provided: %s", certID)
		}
		// Attempt to find this certificates
		cert, err := certManager.TrustedCertificateStore().GetCertificateByCertID(certID)
		if err != nil {
			fatalf("Unable to retrieve certificate with cert ID: %s", certID)
		}
		certsToRemove = append(certsToRemove, cert)
	} else {
		// We got the -g flag, it's a GUN
		toRemove, err := certManager.TrustedCertificateStore().GetCertificatesByCN(
			certRemoveGUN)
		if err != nil {
			fatalf("%v", err)
		}
		certsToRemove = append(certsToRemove, toRemove...)
	}

	// List all the keys about to be removed
	cmd.Printf("The following certificates will be removed:\n\n")
	for _, cert := range certsToRemove {
		// This error can't occur because we're getting certs off of an
		// x509 store that indexes by ID.
		certID, _ := trustmanager.FingerprintCert(cert)
		cmd.Printf("%s - %s\n", cert.Subject.CommonName, certID)
	}
	cmd.Println("\nAre you sure you want to remove these certificates? (yes/no)")

	// Ask for confirmation before removing certificates, unless -y is provided
	if !certRemoveYes {
		confirmed := askConfirm()
		if !confirmed {
			fatalf("Aborting action.")
		}
	}

	// Remove all the certs
	for _, cert := range certsToRemove {
		err = certManager.TrustedCertificateStore().RemoveCert(cert)
		if err != nil {
			fatalf("Failed to remove root certificate for %s", cert.Subject.CommonName)
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
	certManager, err := certs.NewManager(trustDir)
	if err != nil {
		fatalf("Failed to create a new truststore manager with directory: %s", trustDir)
	}

	trustedCerts := certManager.TrustedCertificateStore().GetCertificates()

	cmd.Println("")
	prettyPrintCerts(trustedCerts, cmd.Out())
	cmd.Println("")
}

func printCert(cmd *cobra.Command, cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	certID, err := trustmanager.FingerprintCert(cert)
	if err != nil {
		fatalf("Could not fingerprint certificate: %v", err)
	}

	cmd.Printf("%s %s (expires in: %v days)\n", cert.Subject.CommonName, certID, math.Floor(timeDifference.Hours()/24))
}

// cert by repo name then expiry time.  Don't bother sorting by fingerprint.
type certSorter []*x509.Certificate

func (t certSorter) Len() int      { return len(t) }
func (t certSorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t certSorter) Less(i, j int) bool {
	if t[i].Subject.CommonName < t[j].Subject.CommonName {
		return true
	} else if t[i].Subject.CommonName > t[j].Subject.CommonName {
		return false
	}

	return t[i].NotAfter.Before(t[j].NotAfter)
}

// Given a list of Ceritifcates in order of listing preference, pretty-prints
// the cert common name, fingerprint, and expiry
func prettyPrintCerts(certs []*x509.Certificate, writer io.Writer) {
	if len(certs) == 0 {
		writer.Write([]byte("\nNo trusted root certificates present.\n\n"))
		return
	}

	sort.Stable(certSorter(certs))

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{
		"GUN", "Fingerprint of Trusted Root Certificate", "Expires In"})
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetAutoWrapText(false)

	for _, c := range certs {
		days := math.Floor(c.NotAfter.Sub(time.Now()).Hours() / 24)
		expiryString := "< 1 day"
		if days == 1 {
			expiryString = "1 day"
		} else if days > 1 {
			expiryString = fmt.Sprintf("%d days", int(days))
		}

		certID, err := trustmanager.FingerprintCert(c)
		if err != nil {
			fatalf("Could not fingerprint certificate: %v", err)
		}

		table.Append([]string{c.Subject.CommonName, certID, expiryString})
	}
	table.Render()
}
