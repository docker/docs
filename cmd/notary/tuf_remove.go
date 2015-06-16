package main

import "github.com/spf13/cobra"

var cmdTufRemove = &cobra.Command{
	Use:   "remove [ QDN ] <target>",
	Short: "Removes a target from the TUF repo.",
	Long:  "removes a target from the local TUF repo identified by a Qualified Docker Name.",
	Run:   tufRemove,
}

func tufRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}
