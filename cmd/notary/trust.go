package main

import (
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var cmdtrust = &cobra.Command{
	Use:   "trust [path/url of certificate to add]",
	Short: "Trusts a new certificate for a specific QDN.",
	Long:  "Adds a the certificate to the trusted certificate authority list for the specified QDN",
	Run:   trust,
}

func trust(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a URL or file.")
	}

	// Verify if argument is a valid URL
	url, err := url.Parse(args[0])
	if err == nil && url.Scheme != "" {
		err = caStore.AddCertFromURL(args[0])
		if err != nil {
			fatalf("error adding certificate to CA Store: %v", err)
		}
		// Verify is argument is a valid file
	} else if _, err := os.Stat(args[0]); err == nil {
		if err := caStore.AddCertFromFile(args[0]); err != nil {
			fatalf("error adding certificate from file: %v", err)
		}
	} else {
		fatalf("please provide a file location or URL for CA certificate.")
	}
}
