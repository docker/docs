package main

import (
	"archive/zip"
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/pkg/passphrase"
	"github.com/docker/notary/trustmanager"

	"github.com/spf13/cobra"
)

func init() {
	cmdKeys.AddCommand(cmdKeysRemoveRootKey)
	cmdKeys.AddCommand(cmdKeysGenerateRootKey)

	cmdKeysExport.Flags().StringVarP(&keysExportGUN, "gun", "g", "", "Globally unique name to export keys for. A new password will be set for all the keys. Output format is a zip archive.")
	cmdKeys.AddCommand(cmdKeysExport)
	cmdKeys.AddCommand(cmdKeysExportRoot)
	cmdKeys.AddCommand(cmdKeysImport)
	cmdKeys.AddCommand(cmdKeysImportRoot)
}

var cmdKeys = &cobra.Command{
	Use:   "keys",
	Short: "Operates on keys.",
	Long:  "operations on private keys.",
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

var keysExportGUN string

var cmdKeysExport = &cobra.Command{
	Use:   "export [ filename ]",
	Short: "Exports keys to a ZIP file.",
	Long:  "exports a collection of keys. The keys are reencrypted with a new passphrase. The output is a ZIP file.",
	Run:   keysExport,
}

var cmdKeysExportRoot = &cobra.Command{
	Use:   "export-root [ keyID ] [ filename ]",
	Short: "Exports given root key to a file.",
	Long:  "exports a root key, without reencrypting. The output is a PEM file.",
	Run:   keysExportRoot,
}

var cmdKeysImport = &cobra.Command{
	Use:   "import [ filename ]",
	Short: "Imports keys from a ZIP file.",
	Long:  "imports one or more keys from a ZIP file.",
	Run:   keysImport,
}

var cmdKeysImportRoot = &cobra.Command{
	Use:   "import-root [ keyID ] [ filename ]",
	Short: "Imports root key.",
	Long:  "imports a root key from a PEM file.",
	Run:   keysImportRoot,
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

// keysExport exports a collection of keys to a ZIP file
func keysExport(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify output filename for export")
	}

	exportFilename := args[0]

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	exportFile, err := os.Create(exportFilename)
	if err != nil {
		fatalf("error creating output file: %v", err)
	}

	// Must use a different passphrase retriever to avoid caching the
	// unlocking passphrase and reusing that.
	exportRetriever := passphrase.PromptRetriever()
	if keysExportGUN != "" {
		err = keyStoreManager.ExportKeysByGUN(exportFile, keysExportGUN, exportRetriever)
	} else {
		err = keyStoreManager.ExportAllKeys(exportFile, exportRetriever)
	}

	exportFile.Close()

	if err != nil {
		fatalf("error exporting keys: %v", err)
		os.Remove(exportFilename)
	}
}

// keysExportRoot exports a root key by ID to a PEM file
func keysExportRoot(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify key ID and output filename for export")
	}

	keyID := args[0]
	exportFilename := args[1]

	if len(keyID) != 64 {
		fatalf("please specify a valid root key ID")
	}

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	exportFile, err := os.Create(exportFilename)
	if err != nil {
		fatalf("error creating output file: %v", err)
	}
	err = keyStoreManager.ExportRootKey(exportFile, keyID)
	exportFile.Close()
	if err != nil {
		fatalf("error exporting root key: %v", err)
		os.Remove(exportFilename)
	}
}

// keysImport imports keys from a ZIP file
func keysImport(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify input filename for import")
	}

	importFilename := args[0]

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	zipReader, err := zip.OpenReader(importFilename)
	if err != nil {
		fatalf("opening file for import: %v", err)
	}
	defer zipReader.Close()

	err = keyStoreManager.ImportKeysZip(zipReader.Reader)

	if err != nil {
		fatalf("error importing keys: %v", err)
	}
}

// keysImportRoot imports a root key from a PEM file
func keysImportRoot(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify key ID and input filename for import")
	}

	keyID := args[0]
	importFilename := args[1]

	if len(keyID) != 64 {
		fatalf("please specify a valid root key ID")
	}

	parseConfig()

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(trustDir, retriever)
	if err != nil {
		fatalf("failed to create a new truststore manager with directory: %s", trustDir)
	}

	importFile, err := os.Open(importFilename)
	if err != nil {
		fatalf("opening file for import: %v", err)
	}
	defer importFile.Close()

	err = keyStoreManager.ImportRootKey(importFile, keyID)

	if err != nil {
		fatalf("error importing root key: %v", err)
	}
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
