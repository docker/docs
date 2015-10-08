package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/docker/notary/pkg/passphrase"
	"github.com/docker/notary/version"
)

const configFileName string = "config"
const defaultTrustDir string = ".notary/"
const defaultServerURL = "https://notary-server:4443"
const idSize = 64

var rawOutput bool
var trustDir string
var remoteTrustServer string
var verbose bool
var retriever passphrase.Retriever

func init() {
	retriever = getPassphraseRetriever()
}

func parseConfig() {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	}

	if trustDir == "" {
		// Get home directory for current user
		homeDir, err := homedir.Dir()
		if err != nil {
			fatalf("cannot get current user home directory: %v", err)
		}
		if homeDir == "" {
			fatalf("cannot get current user home directory")
		}
		trustDir = filepath.Join(homeDir, filepath.Dir(defaultTrustDir))

		logrus.Debugf("no trust directory provided, using default: %s", trustDir)
	} else {
		logrus.Debugf("trust directory provided: %s", trustDir)
	}

	// Setup the configuration details
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(trustDir)
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
	serverURL := os.Getenv("NOTARY_SERVER_URL")
	if serverURL == "" {
		serverURL = defaultServerURL
	}

	var notaryCmd = &cobra.Command{
		Use:   "notary",
		Short: "notary allows the creation of trusted collections.",
		Long:  "notary allows the creation and management of collections of signed targets, allowing the signing and validation of arbitrary content.",
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of notary",
		Long:  `print the version number of notary`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("notary\n Version:    %s\n Git commit: %s\n", version.NotaryVersion, version.GitCommit)
		},
	}

	notaryCmd.AddCommand(versionCmd)

	notaryCmd.PersistentFlags().StringVarP(&trustDir, "trustdir", "d", "", "directory where the trust data is persisted to")
	notaryCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	notaryCmd.AddCommand(cmdKey)
	notaryCmd.AddCommand(cmdCert)
	notaryCmd.AddCommand(cmdTufInit)
	cmdTufInit.Flags().StringVarP(&remoteTrustServer, "server", "s", serverURL, "Remote trust server location")
	notaryCmd.AddCommand(cmdTufList)
	cmdTufList.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary list to output a nonpretty printed version of the targets list. Useful if you need to parse the list.")
	cmdTufList.Flags().StringVarP(&remoteTrustServer, "server", "s", serverURL, "Remote trust server location")
	notaryCmd.AddCommand(cmdTufAdd)
	notaryCmd.AddCommand(cmdTufRemove)
	notaryCmd.AddCommand(cmdTufPublish)
	cmdTufPublish.Flags().StringVarP(&remoteTrustServer, "server", "s", serverURL, "Remote trust server location")
	notaryCmd.AddCommand(cmdTufLookup)
	cmdTufLookup.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary lookup to output a nonpretty printed version of the targets list. Useful if you need to parse the list.")
	cmdTufLookup.Flags().StringVarP(&remoteTrustServer, "server", "s", serverURL, "Remote trust server location")
	notaryCmd.AddCommand(cmdVerify)
	cmdVerify.Flags().StringVarP(&remoteTrustServer, "server", "s", serverURL, "Remote trust server location")

	notaryCmd.Execute()
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf("* fatal: "+format+"\n", args...)
	os.Exit(1)
}

func askConfirm() bool {
	var res string
	_, err := fmt.Scanln(&res)
	if err != nil {
		return false
	}
	if strings.EqualFold(res, "y") || strings.EqualFold(res, "yes") {
		return true
	}
	return false
}

func getPassphraseRetriever() passphrase.Retriever {
	baseRetriever := passphrase.PromptRetriever()
	env := map[string]string{
		"root":     os.Getenv("NOTARY_ROOT_PASSPHRASE"),
		"targets":  os.Getenv("NOTARY_TARGET_PASSPHRASE"),
		"snapshot": os.Getenv("NOTARY_SNAPSHOT_PASSPHRASE"),
	}

	return func(keyName string, alias string, createNew bool, numAttempts int) (string, bool, error) {
		if v := env[alias]; v != "" {
			return v, numAttempts > 1, nil
		}
		return baseRetriever(keyName, alias, createNew, numAttempts)
	}
}
