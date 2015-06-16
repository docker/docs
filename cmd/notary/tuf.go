package main

import "github.com/spf13/cobra"

var cmdTuf = &cobra.Command{
	Use:   "tuf",
	Short: "Manages trust of data for notary.",
	Long:  "manages signed repository metadata.",
	Run:   nil,
}

var remoteTrustServer string
