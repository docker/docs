package main

import (
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/docker/notary"
	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdCertTemplate = usageTemplate{
	Use:   "cert",
	Short: "Operates on certificates.",
	Long:  `Operations on certificates.`,
}

var cmdCertListTemplate = usageTemplate{
	Use:   "list",
	Short: "Lists certificates.",
	Long:  "Lists root certificates known to notary.",
}

var cmdCertRemoveTemplate = usageTemplate{
	Use:   "remove [ certID ]",
	Short: "Removes the certificate with the given cert ID.",
	Long:  "Remove the certificate with the given cert ID from the local host.",
}

type certCommander struct {
	// these need to be set
	configGetter func() (*viper.Viper, error)
	retriever    passphrase.Retriever

	// these are for command line parsing - no need to set
	certRemoveGUN string
	certRemoveYes bool
}

func (c *certCommander) GetCommand() *cobra.Command {
	cmd := cmdCertTemplate.ToCommand(nil)
	cmd.AddCommand(cmdCertListTemplate.ToCommand(c.certList))

	cmdCertRemove := cmdCertRemoveTemplate.ToCommand(c.certRemove)
	cmdCertRemove.Flags().StringVarP(
		&c.certRemoveGUN, "gun", "g", "", "Globally unique name to delete certificates for")
	cmdCertRemove.Flags().BoolVarP(
		&c.certRemoveYes, "yes", "y", false, "Answer yes to the removal question (no confirmation)")

	cmd.AddCommand(cmdCertRemove)

	return cmd
}

// certRemove deletes a certificate given a cert ID or a gun
// If given a gun, certRemove will also remove local TUF data
func (c *certCommander) certRemove(cmd *cobra.Command, args []string) error {
	// If the user hasn't provided -g with a gun, or a cert ID, show usage
	// If the user provided -g and a cert ID, also show usage
	if (len(args) < 1 && c.certRemoveGUN == "") || (len(args) > 0 && c.certRemoveGUN != "") {
		cmd.Usage()
		return fmt.Errorf("Must specify the cert ID or the GUN of the certificates to remove")
	}
	config, err := c.configGetter()
	if err != nil {
		return err
	}

	trustDir := config.GetString("trust_dir")
	certPath := filepath.Join(trustDir, notary.TrustedCertsDir)
	certStore, err := trustmanager.NewX509FilteredFileStore(
		certPath,
		trustmanager.FilterCertsExpiredSha1,
	)
	if err != nil {
		return fmt.Errorf("Failed to create a new truststore with directory: %s", trustDir)
	}

	var certsToRemove []*x509.Certificate
	var certFoundByID *x509.Certificate
	var removeTrustData bool

	// If there is no GUN, we expect a cert ID
	if c.certRemoveGUN == "" {
		certID := args[0]
		// Attempt to find this certificate
		certFoundByID, err = certStore.GetCertificateByCertID(certID)
		if err != nil {
			// This is an invalid ID, the user might have forgotten a character
			if len(certID) != notary.Sha256HexSize {
				return fmt.Errorf("Unable to retrieve certificate with invalid certificate ID provided: %s", certID)
			}
			return fmt.Errorf("Unable to retrieve certificate with cert ID: %s", certID)
		}
		// the GUN is the CN from the certificate
		c.certRemoveGUN = certFoundByID.Subject.CommonName
		certsToRemove = []*x509.Certificate{certFoundByID}
	}

	toRemove, err := certStore.GetCertificatesByCN(c.certRemoveGUN)
	// We could not find any certificates matching the user's query, so propagate the error
	if err != nil {
		return fmt.Errorf("%v", err)
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
		cmd.Printf("\nAll local trust data will be removed for %s\n", c.certRemoveGUN)
	}
	cmd.Println("\nAre you sure you want to remove these certificates? (yes/no)")

	// Ask for confirmation before removing certificates, unless -y is provided
	if !c.certRemoveYes {
		confirmed := askConfirm()
		if !confirmed {
			return fmt.Errorf("Aborting action.")
		}
	}

	if removeTrustData {
		// Remove all TUF data, so call RemoveTrustData on a NotaryRepository with the GUN
		// no online operations are performed so the transport argument is nil
		nRepo, err := notaryclient.NewNotaryRepository(
			trustDir, c.certRemoveGUN, getRemoteTrustServer(config), nil, c.retriever)
		if err != nil {
			return fmt.Errorf("Could not establish trust data for GUN %s", c.certRemoveGUN)
		}
		// DeleteTrustData will pick up all of the same certificates by GUN (CN) and remove them
		err = nRepo.DeleteTrustData()
		if err != nil {
			return fmt.Errorf("Failed to delete trust data for %s", c.certRemoveGUN)
		}
	} else {
		for _, cert := range certsToRemove {
			err = certStore.RemoveCert(cert)
			if err != nil {
				return fmt.Errorf("Failed to remove cert %s", cert)
			}
		}
	}
	return nil
}

func (c *certCommander) certList(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		cmd.Usage()
		return fmt.Errorf("")
	}
	config, err := c.configGetter()
	if err != nil {
		return err
	}

	trustDir := config.GetString("trust_dir")
	certPath := filepath.Join(trustDir, notary.TrustedCertsDir)
	// Load all individual (non-CA) certificates that aren't expired and don't use SHA1
	certStore, err := trustmanager.NewX509FilteredFileStore(
		certPath,
		trustmanager.FilterCertsExpiredSha1,
	)
	if err != nil {
		return fmt.Errorf("Failed to create a new truststore with directory: %s", trustDir)
	}

	trustedCerts := certStore.GetCertificates()

	cmd.Println("")
	prettyPrintCerts(trustedCerts, cmd.Out())
	cmd.Println("")
	return nil
}
