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

	"github.com/docker/notary/trustmanager"
)

const configFileName string = "config"
const configPath string = ".docker/trust/"
const trustDir string = "trusted_certificates/"
const privDir string = "private/"
const rootKeysDir string = "root_keys/"

var rawOutput bool
var caStore trustmanager.X509Store
var certificateStore trustmanager.X509Store
var privKeyStore trustmanager.FileStore

func init() {
	logrus.SetLevel(logrus.DebugLevel)
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
	viper.SetDefault("baseTrustDir", path.Join(homeDir, path.Dir(configPath)))

	// Get the final value for the CA directory
	finalTrustDir := path.Join(viper.GetString("baseTrustDir"), trustDir)
	finalPrivDir := path.Join(viper.GetString("baseTrustDir"), privDir)

	// Load all CAs that aren't expired and don't use SHA1
	caStore, err = trustmanager.NewX509FilteredFileStore(finalTrustDir, func(cert *x509.Certificate) bool {
		return cert.IsCA && cert.BasicConstraintsValid && cert.SubjectKeyId != nil &&
			time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
	if err != nil {
		fatalf("could not create CA X509FileStore: %v", err)
	}

	// Load all individual (nonCA) certificates that aren't expired and don't use SHA1
	certificateStore, err = trustmanager.NewX509FilteredFileStore(finalTrustDir, func(cert *x509.Certificate) bool {
		return !cert.IsCA &&
			time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
	if err != nil {
		fatalf("could not create Certificate X509FileStore: %v", err)
	}

	//TODO(mccauley): Appears unused? Remove it? Or is it here for early failure?
	privKeyStore, err = trustmanager.NewKeyFileStore(finalPrivDir,
		func(string, string, bool, int) (string, bool, error) { return "", false, nil })
	if err != nil {
		fatalf("could not create KeyFileStore: %v", err)
	}
}

func main() {
	var NotaryCmd = &cobra.Command{
		Use:   "notary",
		Short: "notary allows the creation of trusted collections.",
		Long:  "notary allows the creation and management of collections of signed targets, allowing the signing and validation of arbitrary content.",
	}

	NotaryCmd.AddCommand(cmdKeys)
	NotaryCmd.AddCommand(cmdTufInit)
	NotaryCmd.AddCommand(cmdTufList)
	cmdTufList.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary list to output a nonpretty printed version of the targets list. Useful if you need to parse the list.")
	NotaryCmd.AddCommand(cmdTufAdd)
	NotaryCmd.AddCommand(cmdTufRemove)
	NotaryCmd.AddCommand(cmdTufPublish)
	cmdTufPublish.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
	NotaryCmd.AddCommand(cmdTufLookup)
	cmdTufLookup.Flags().BoolVarP(&rawOutput, "raw", "", false, "Instructs notary lookup to output a nonpretty printed version of the targets list. Useful if you need to parse the list.")
	cmdTufLookup.Flags().StringVarP(&remoteTrustServer, "remote", "r", "", "Remote trust server location")
	NotaryCmd.AddCommand(cmdVerify)

	NotaryCmd.Execute()
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf("* fatal: "+format+"\n", args...)
	os.Exit(1)
}
