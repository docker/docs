package main

import (
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/trustmanager"

	"github.com/spf13/cobra"
)

func init() {
	cmdKeys.AddCommand(cmdKeysRemoveRootKey)
	cmdKeys.AddCommand(cmdKeysGenerateRootKey)
}

var cmdKeys = &cobra.Command{
	Use:   "keys",
	Short: "Operates on root keys.",
	Long:  "operations on private root keys.",
	Run:   keysList,
}

var cmdKeysRemoveRootKey = &cobra.Command{
	Use:   "remove [ keyID ]",
	Short: "Removes the root key with the given keyID.",
	Long:  "remove the root key with the given keyID from the local host.",
	Run:   keysRemoveRootKey,
}

var cmdKeysGenerateRootKey = &cobra.Command{
	Use:   "generate [ algorithm ]",
	Short: "Generates a new root key with a given algorithm.",
	Long:  "generates a new root key with a given algorithm.",
	Run:   keysGenerateRootKey,
}

// keysRemoveRootKey deletes a root private key based on ID
func keysRemoveRootKey(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify the key ID of the root key to remove")
	}

	keyID := args[0]
	if len(keyID) != 64 {
		fatalf("please enter a valid root key ID")
	}
	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	// List all the keys about to be removed
	fmt.Printf("Are you sure you want to remove the following key?\n%s\n (yes/no)\n", keyID)

	// Ask for confirmation before removing keys
	confirmed := askConfirm()
	if !confirmed {
		fatalf("aborting action.")
	}

	// Remove all the keys under the Global Unique Name
	err = keyStoreManager.RootKeyStore().RemoveKey(keyID)
	if err != nil {
		fatalf("failed to remove root key with key ID: %s", keyID)
	}

	fmt.Printf("Root key %s removed\n", keyID)
}

func keysList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	fmt.Println("")
	fmt.Println("# Trusted Certificates:")
	trustedCerts := keyStoreManager.TrustedCertificateStore().GetCertificates()
	for _, c := range trustedCerts {
		printCert(c)
	}

	fmt.Println("")
	fmt.Println("# Root keys: ")
	for _, k := range keyStoreManager.RootKeyStore().ListKeys() {
		fmt.Println(k)
	}

	fmt.Println("")
	fmt.Println("# Signing keys: ")
	for _, k := range keyStoreManager.NonRootKeyStore().ListKeys() {
		printKey(k)
	}
}

func keysGenerateRootKey(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify an Algorithm (RSA, ECDSA)")
	}

	algorithm := args[0]
	allowedCiphers := map[string]bool{
		"rsa":   true,
		"ecdsa": true,
	}

	if !allowedCiphers[strings.ToLower(algorithm)] {
		fatalf("algorithm not allowed, possible values are: RSA, ECDSA")
	}

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	keyID, err := keyStoreManager.GenRootKey(algorithm)
	if err != nil {
		fatalf("failed to create a new root key: %v", err)
	}

	fmt.Printf("Generated new %s key with keyID: %s\n", algorithm, keyID)
}

func printCert(cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	certID, err := trustmanager.FingerprintCert(cert)
	if err != nil {
		fatalf("could not fingerprint certificate: %v", err)
	}

	fmt.Printf("%s %s (expires in: %v days)\n", cert.Subject.CommonName, certID, math.Floor(timeDifference.Hours()/24))
}

func printKey(keyPath string) {
	keyID := filepath.Base(keyPath)
	gun := filepath.Dir(keyPath)
	fmt.Printf("%s %s\n", gun, keyID)
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
