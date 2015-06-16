package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/docker/vetinari/trustmanager"
	"github.com/spf13/cobra"
)

var cmdKeysTrust = &cobra.Command{
	Use:   "trust [ QDN ] [ certificate ]",
	Short: "Trusts a new certificate for a specific QDN.",
	Long:  "Adds a the certificate to the trusted certificate authority list for the specified Qualified Docker Name.",
	Run:   keysTrust,
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
		print_cert(cert)
	} else if _, err := os.Stat(certLocationStr); err == nil {
		if err := caStore.AddCertFromFile(certLocationStr); err != nil {
			fatalf("error adding certificate from file: %v", err)
		}
	} else {
		fatalf("please provide a file location or URL for CA certificate.")
	}
}
