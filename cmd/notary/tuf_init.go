package main

import "github.com/spf13/cobra"

var cmdTufInit = &cobra.Command{
	Use:   "init [ QDN ]",
	Short: "initializes the local TUF repository.",
	Long:  "creates locally the initial set of TUF metadata for the Qualified Docker Name.",
	Run:   tufInit,
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}
