package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/docker/vetinari/trustmanager"
)

const configFileName string = "config"

// Default paths should end with a '/' so directory creation works correctly
const configPath string = ".docker/trust/"
const trustDir string = configPath + "repository_certificates/"
const privDir string = configPath + "private/"
const tufDir string = configPath + "tuf/"

var caStore trustmanager.X509Store
var rawOutput bool

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetOutput(os.Stderr)
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

	// Setup the configuration details
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(path.Join(homeDir, path.Dir(configPath)))
	viper.SetConfigType("json")

	// Find and read the config file
	err = viper.ReadInConfig()
	if err != nil {
		// Ignore if the configuration file doesn't exist, we can use the defaults
		if !os.IsNotExist(err) {
			fatalf("fatal error config file: %v", err)
		}
	}

	// Set up the defaults for our config
	viper.SetDefault("trustDir", path.Join(homeDir, path.Dir(trustDir)))
	viper.SetDefault("privDir", path.Join(homeDir, path.Dir(privDir)))
	viper.SetDefault("tufDir", path.Join(homeDir, path.Dir(tufDir)))

	// Get the final value for the CA directory
	finalTrustDir := viper.GetString("trustDir")
	finalPrivDir := viper.GetString("privDir")

	// Ensure the existence of the CAs directory
	err = trustmanager.CreateDirectory(finalTrustDir)
	if err != nil {
		fatalf("could not create directory: %v", err)
	}
	err = trustmanager.CreateDirectory(finalPrivDir)
	if err != nil {
		fatalf("could not create directory: %v", err)
	}

	// Load all CAs that aren't expired and don't use SHA1
	// We could easily add "return cert.IsCA && cert.BasicConstraintsValid" in order
	// to have only valid CA certificates being loaded
	caStore = trustmanager.NewX509FilteredFileStore(finalTrustDir, func(cert *x509.Certificate) bool {
		return time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
}

func main() {
	var NotaryCmd = &cobra.Command{
		Use:   "notary",
		Short: "notary creates trust for docker",
		Long:  "notary is the main trust-related command for Docker.",
	}

	NotaryCmd.AddCommand(cmdKeys)
	NotaryCmd.AddCommand(cmdTufInit)
	NotaryCmd.AddCommand(cmdTufList)
	cmdTufList.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary list to output a non-pretty printed version of the targets list. Useful if you need to parse the list.")
	NotaryCmd.AddCommand(cmdTufAdd)
	NotaryCmd.AddCommand(cmdTufRemove)
	NotaryCmd.AddCommand(cmdTufPublish)
	cmdTufPublish.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
	NotaryCmd.AddCommand(cmdTufLookup)
	cmdTufLookup.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary lookup to output a non-pretty printed version of the targets list. Useful if you need to parse the list.")
	cmdTufLookup.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")

	NotaryCmd.Execute()
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf("* fatal: "+format+"\n", args...)
	os.Exit(1)
}
