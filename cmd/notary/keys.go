package main

import "github.com/spf13/cobra"

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Operates on keys",
	Long:  "operations on signature keys and trusted certificate authorities",
}
