package main

import "github.com/spf13/cobra"

var cmdTufAdd = &cobra.Command{
	Use:   "add [ QDN ] <target> <file path>",
	Short: "pushes local updates.",
	Long:  "pushes all local updates within a specific TUF repo to remote trust server.",
	Run:   tufAdd,
}

func tufAdd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}
