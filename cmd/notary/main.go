package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/docker/notary/pkg/passphrase"
)

const configFileName string = "config"
const defaultTrustDir string = ".notary/"
const defaultServerURL = "https://notary-server:4443"

var rawOutput bool
var TrustDir string
var RemoteTrustServer string
var Verbose bool
var retriever passphrase.Retriever

func init() {
	retriever = passphrase.PromptRetriever()
}

func parseConfig() {
	if Verbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	}

	if TrustDir == "" {
		// Retrieve current user to get home directory
		usr, err := user.Current()
		if err != nil {
			fatalf("cannot get current user: %v", err)
		}

		// Get home directory for current user
		homeDir := usr.HomeDir
		if homeDir == "" {
			fatalf("cannot get current user home directory")
		}
		TrustDir = filepath.Join(homeDir, filepath.Dir(defaultTrustDir))

		logrus.Debugf("no trust directory provided, using default: %s", TrustDir)
	} else {
		logrus.Debugf("trust directory provided: %s", TrustDir)
	}

	// Setup the configuration details
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(TrustDir)
	viper.SetConfigType("json")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Debugf("configuration file not found, using defaults")
		// Ignore if the configuration file doesn't exist, we can use the defaults
		if !os.IsNotExist(err) {
			fatalf("fatal error config file: %v", err)
		}
	}
}

func main() {
	var NotaryCmd = &cobra.Command{
		Use:   "notary",
		Short: "notary allows the creation of trusted collections.",
		Long:  "notary allows the creation and management of collections of signed targets, allowing the signing and validation of arbitrary content.",
	}

	NotaryCmd.PersistentFlags().StringVarP(&TrustDir, "trustdir", "d", "", "directory where the trust data is persisted to")
	NotaryCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	NotaryCmd.AddCommand(cmdKeys)
	NotaryCmd.AddCommand(cmdTufInit)
	cmdTufInit.Flags().StringVarP(&RemoteTrustServer, "server", "s", defaultServerURL, "Remote trust server location")
	NotaryCmd.AddCommand(cmdTufList)
	cmdTufList.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary list to output a nonpretty printed version of the targets list. Useful if you need to parse the list.")
	cmdTufList.Flags().StringVarP(&RemoteTrustServer, "server", "s", defaultServerURL, "Remote trust server location")
	NotaryCmd.AddCommand(cmdTufAdd)
	NotaryCmd.AddCommand(cmdTufRemove)
	NotaryCmd.AddCommand(cmdTufPublish)
	cmdTufPublish.Flags().StringVarP(&RemoteTrustServer, "server", "s", defaultServerURL, "Remote trust server location")
	NotaryCmd.AddCommand(cmdTufLookup)
	cmdTufLookup.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary lookup to output a nonpretty printed version of the targets list. Useful if you need to parse the list.")
	cmdTufLookup.Flags().StringVarP(&RemoteTrustServer, "server", "s", defaultServerURL, "Remote trust server location")
	NotaryCmd.AddCommand(cmdVerify)

	NotaryCmd.Execute()
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf("* fatal: "+format+"\n", args...)
	os.Exit(1)
}
