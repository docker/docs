package main

import "github.com/spf13/cobra"

var cmdTufList = &cobra.Command{
	Use:   "list [ QDN ]",
	Short: "Lists all targets in a TUF repository.",
	Long:  "lists all the targets in the TUF repository identified by the Qualified Docker Name.",
	Run:   tufList,
}

func tufList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}
