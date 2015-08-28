package main

import (
	"github.com/Sirupsen/logrus"
	notaryclient "github.com/docker/notary/client"
	"github.com/spf13/cobra"
)

func init() {
	cmdMeta.AddCommand(cmdRoleDisplay)
	cmdMeta.AddCommand(cmdReplaceKey)
	cmdMeta.AddCommand(cmdAddKey)
	cmdMeta.AddCommand(cmdRemoveKey)
}

var cmdMeta = &cobra.Command{
	Use:   "meta",
	Short: "Operates on repository metadata.",
	Long:  "Operations to manage key usage and delegations within a repository.",
}

var cmdRoleDisplay = &cobra.Command{
	Use:   "display [ GUN ] <role>",
	Short: "Shows metadata about a role",
	Long:  "Display all metadata about a role including the associated keys, the role name, and the owner name if applicable.",
	Run:   metaRoleDisplay,
}

var cmdReplaceKey = &cobra.Command{
	Use:   "replace [ GUN ] <role>",
	Short: "Replace all keys for role.",
	Long:  "Replaces all keys for the given role.",
	Run:   metaReplaceKey,
}

var cmdAddKey = &cobra.Command{
	Use:   "add [ GUN ] <role>",
	Short: "Add key to role.",
	Long:  "Adds a key to the given role.",
	Run:   metaAddKey,
}

var cmdRemoveKey = &cobra.Command{
	Use:   "remove [ GUN ] <role>",
	Short: "Remove a key role.",
	Long:  "Removes a key from the given role.",
	Run:   metaRemoveKey,
}

func metaRoleDisplay(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and role")
	}

	gun := args[0]
	parseConfig()

	logrus.Debug("Displaying info")
	_, err := notaryclient.NewNotaryRepository(trustDir, gun, remoteTrustServer, getTransport(), retriever)
	if err != nil {
		fatalf(err.Error())
	}
}

// metaReplaceKey rotates the key for the given GUN and role.
func metaReplaceKey(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and role")
	}

	gun := args[0]
	role := args[1]
	parseConfig()

	r, err := notaryclient.NewNotaryRepository(trustDir, gun, remoteTrustServer, getTransport(), retriever)
	if err != nil {
		fatalf(err.Error())
	}
	// args[2:] should be a list of paths to public key files.
	// they will be imported into the TUF repo for the given
	// role, replaceing any existing keys.
	if err := r.ReplaceKeys(role, args[2:]...); err != nil {
		fatalf(err.Error())
	}
}

func metaAddKey(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and role")
	}

	gun := args[0]
	parseConfig()

	_, err := notaryclient.NewNotaryRepository(trustDir, gun, remoteTrustServer, getTransport(), retriever)
	if err != nil {
		fatalf(err.Error())
	}

}

func metaRemoveKey(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and role")
	}

	gun := args[0]
	parseConfig()

	_, err := notaryclient.NewNotaryRepository(trustDir, gun, remoteTrustServer, getTransport(), retriever)
	if err != nil {
		fatalf(err.Error())
	}
}
