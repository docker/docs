package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdTufPush = &cobra.Command{
	Use:   "push [ QDN ]",
	Short: "initializes the local TUF repository.",
	Long:  "creates locally the initial set of TUF metadata for the Qualified Docker Name.",
	Run:   tufPush,
}

func init() {
	cmdTufPush.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
}

func tufPush(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}

	fmt.Println("Remote trust server configured: " + remoteTrustServer)
}
