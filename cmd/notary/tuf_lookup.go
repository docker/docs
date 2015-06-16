package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdTufLookup = &cobra.Command{
	Use:   "lookup [ QDN ] <target name>",
	Short: "Looks up a specific TUF target in a repository.",
	Long:  "looks up a TUF target in a repository given a Qualified Docker Name.",
	Run:   tufLookup,
}

func init() {
	cmdTufLookup.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
}

func tufLookup(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}

	fmt.Println("Remote trust server configured: " + remoteTrustServer)

}
