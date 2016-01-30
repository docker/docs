package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configDir        = ".notary/"
	defaultServerURL = "https://notary-server:4443"
)

var (
	debug             bool
	verbose           bool
	roles             []string
	trustDir          string
	configFile        string
	remoteTrustServer string
	configPath        string
	configFileName    = "config"
	configFileExt     = "json"
	retriever         passphrase.Retriever
	getRetriever      = getPassphraseRetriever
	mainViper         = viper.New()
)

func init() {
	retriever = getPassphraseRetriever()
}

func parseConfig() *viper.Viper {
	setVerbosityLevel()

	// Get home directory for current user
	homeDir, err := homedir.Dir()
	if err != nil {
		fatalf("Cannot get current user home directory: %v", err)
	}
	if homeDir == "" {
		fatalf("Cannot get current user home directory")
	}

	// By default our trust directory (where keys are stored) is in ~/.notary/
	mainViper.SetDefault("trust_dir", filepath.Join(homeDir, filepath.Dir(configDir)))

	// If there was a commandline configFile set, we parse that.
	// If there wasn't we attempt to find it on the default location ~/.notary/config
	if configFile != "" {
		configFileExt = strings.TrimPrefix(filepath.Ext(configFile), ".")
		configFileName = strings.TrimSuffix(filepath.Base(configFile), filepath.Ext(configFile))
		configPath = filepath.Dir(configFile)
	} else {
		configPath = filepath.Join(homeDir, filepath.Dir(configDir))
	}

	// Setup the configuration details into viper
	mainViper.SetConfigName(configFileName)
	mainViper.SetConfigType(configFileExt)
	mainViper.AddConfigPath(configPath)

	// Find and read the config file
	err = mainViper.ReadInConfig()
	if err != nil {
		logrus.Debugf("Configuration file not found, using defaults")
		// If we were passed in a configFile via -c, bail if it doesn't exist,
		// otherwise ignore it: we can use the defaults
		if configFile != "" || !os.IsNotExist(err) {
			fatalf("error opening config file %v", err)
		}
	}

	// At this point we either have the default value or the one set by the config.
	// Either way, the command-line flag has precedence and overwrites the value
	if trustDir != "" {
		mainViper.Set("trust_dir", trustDir)
	}

	// Expands all the possible ~/ that have been given, either through -d or config
	// If there is no error, use it, if not, attempt to use whatever the user gave us
	expandedTrustDir, err := homedir.Expand(mainViper.GetString("trust_dir"))
	if err == nil {
		mainViper.Set("trust_dir", expandedTrustDir)
	}
	logrus.Debugf("Using the following trust directory: %s", mainViper.GetString("trust_dir"))

	return mainViper
}

func setupCommand(notaryCmd *cobra.Command) {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of notary",
		Long:  `print the version number of notary`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("notary\n Version:    %s\n Git commit: %s\n", version.NotaryVersion, version.GitCommit)
		},
	}

	notaryCmd.AddCommand(versionCmd)

	notaryCmd.PersistentFlags().StringVarP(&trustDir, "trustDir", "d", "", "Directory where the trust data is persisted to")
	notaryCmd.PersistentFlags().StringVarP(&configFile, "configFile", "c", "", "Path to the configuration file to use")
	notaryCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	notaryCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Debug output")
	notaryCmd.PersistentFlags().StringVarP(&remoteTrustServer, "server", "s", "", "Remote trust server location")

	cmdKeyGenerator := &keyCommander{
		configGetter: parseConfig,
		retriever:    retriever,
	}

	cmdDelegationGenerator := &delegationCommander{
		configGetter: parseConfig,
		retriever:    retriever,
	}

	notaryCmd.AddCommand(cmdKeyGenerator.GetCommand())
	notaryCmd.AddCommand(cmdDelegationGenerator.GetCommand())
	notaryCmd.AddCommand(cmdCert)
	notaryCmd.AddCommand(cmdTufInit)
	cmdTufList.Flags().StringSliceVarP(&roles, "roles", "r", nil, "Delegation roles to list targets for (will shadow targets role)")
	notaryCmd.AddCommand(cmdTufList)
	cmdTufAdd.Flags().StringSliceVarP(&roles, "roles", "r", nil, "Delegation roles to add this target to")
	notaryCmd.AddCommand(cmdTufAdd)
	cmdTufRemove.Flags().StringSliceVarP(&roles, "roles", "r", nil, "Delegation roles to remove this target from")
	notaryCmd.AddCommand(cmdTufRemove)
	notaryCmd.AddCommand(cmdTufStatus)
	notaryCmd.AddCommand(cmdTufPublish)
	notaryCmd.AddCommand(cmdTufLookup)
	notaryCmd.AddCommand(cmdTufVerify)
}

func main() {
	var notaryCmd = &cobra.Command{
		Use:   "notary",
		Short: "Notary allows the creation of trusted collections.",
		Long:  "Notary allows the creation and management of collections of signed targets, allowing the signing and validation of arbitrary content.",
	}
	notaryCmd.SetOutput(os.Stdout)
	setupCommand(notaryCmd)
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
		"targets":  os.Getenv("NOTARY_TARGETS_PASSPHRASE"),
		"snapshot": os.Getenv("NOTARY_SNAPSHOT_PASSPHRASE"),
	}

	return func(keyName string, alias string, createNew bool, numAttempts int) (string, bool, error) {
		if v := env[alias]; v != "" {
			return v, numAttempts > 1, nil
		}
		return baseRetriever(keyName, alias, createNew, numAttempts)
	}
}

// Set the logging level to fatal on default, or the most specific level the user specified (debug or error)
func setVerbosityLevel() {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else if verbose {
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logrus.FatalLevel)
	}
	logrus.SetOutput(os.Stderr)
}
