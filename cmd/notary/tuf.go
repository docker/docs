package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdTuf = &cobra.Command{
	Use:   "tuf",
	Short: "Manages trust of data for notary.",
	Long:  "manages signed repository metadata.",
	Run:   nil,
}

var remoteTrustServer string

func init() {
	cmdTuf.AddCommand(cmdTufInit)
	cmdTuf.AddCommand(cmdTufAdd)
	cmdTuf.AddCommand(cmdTufRemove)
	cmdTuf.AddCommand(cmdTufPush)
	cmdTufPush.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
	cmdTuf.AddCommand(cmdTufLookup)
	cmdTufLookup.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
	cmdTuf.AddCommand(cmdTufList)
}

var cmdTufAdd = &cobra.Command{
	Use:   "add [ QDN ] <target> <file path>",
	Short: "pushes local updates.",
	Long:  "pushes all local updates within a specific TUF repo to remote trust server.",
	Run:   tufAdd,
}

var cmdTufRemove = &cobra.Command{
	Use:   "remove [ QDN ] <target>",
	Short: "Removes a target from the TUF repo.",
	Long:  "removes a target from the local TUF repo identified by a Qualified Docker Name.",
	Run:   tufRemove,
}

var cmdTufInit = &cobra.Command{
	Use:   "init [ QDN ]",
	Short: "initializes the local TUF repository.",
	Long:  "creates locally the initial set of TUF metadata for the Qualified Docker Name.",
	Run:   tufInit,
}

var cmdTufList = &cobra.Command{
	Use:   "list [ QDN ]",
	Short: "Lists all targets in a TUF repository.",
	Long:  "lists all the targets in the TUF repository identified by the Qualified Docker Name.",
	Run:   tufList,
}

var cmdTufLookup = &cobra.Command{
	Use:   "lookup [ QDN ] <target name>",
	Short: "Looks up a specific TUF target in a repository.",
	Long:  "looks up a TUF target in a repository given a Qualified Docker Name.",
	Run:   tufLookup,
}

var cmdTufPush = &cobra.Command{
	Use:   "push [ QDN ]",
	Short: "initializes the local TUF repository.",
	Long:  "creates locally the initial set of TUF metadata for the Qualified Docker Name.",
	Run:   tufPush,
}

func tufAdd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}

func tufList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}

func tufLookup(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}

	fmt.Println("Remote trust server configured: " + remoteTrustServer)

}

func tufPush(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}

	fmt.Println("Remote trust server configured: " + remoteTrustServer)
}

func tufRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a QDN")
	}
}
