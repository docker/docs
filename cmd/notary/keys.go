package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/docker/notary"
	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/spf13/cobra"
)

func init() {
	cmdKey.AddCommand(cmdKeyList)
	cmdKey.AddCommand(cmdKeyGenerateRootKey)

	cmdKeyExport.Flags().StringVarP(&keysExportGUN, "gun", "g", "", "Globally unique name to export keys for")
	cmdKey.AddCommand(cmdKeyExport)
	cmdKey.AddCommand(cmdKeyExportRoot)
	cmdKeyExportRoot.Flags().BoolVarP(&keysExportRootChangePassphrase, "change-passphrase", "p", false, "set a new passphrase for the key being exported")
	cmdKey.AddCommand(cmdKeyImport)
	cmdKey.AddCommand(cmdKeyImportRoot)
	cmdKey.AddCommand(cmdRotateKey)
}

var cmdKey = &cobra.Command{
	Use:   "key",
	Short: "Operates on keys.",
	Long:  `operations on private keys.`,
}

var cmdKeyList = &cobra.Command{
	Use:   "list",
	Short: "Lists keys.",
	Long:  "lists keys known to notary.",
	Run:   keysList,
}

var cmdRotateKey = &cobra.Command{
	Use:   "rotate [ GUN ]",
	Short: "Rotate all keys for role.",
	Long:  "Removes all old keys for the given role and generates 1 new key.",
	Run:   keysRotate,
}

var cmdKeyGenerateRootKey = &cobra.Command{
	Use:   "generate [ algorithm ]",
	Short: "Generates a new root key with a given algorithm.",
	Long:  "generates a new root key with a given algorithm.",
	Run:   keysGenerateRootKey,
}

var keysExportGUN string

var cmdKeyExport = &cobra.Command{
	Use:   "export [ filename ]",
	Short: "Exports keys to a ZIP file.",
	Long:  "exports a collection of keys. The keys are reencrypted with a new passphrase. The output is a ZIP file.",
	Run:   keysExport,
}

var keysExportRootChangePassphrase bool

var cmdKeyExportRoot = &cobra.Command{
	Use:   "export-root [ keyID ] [ filename ]",
	Short: "Exports given root key to a file.",
	Long:  "exports a root key, without reencrypting. The output is a PEM file.",
	Run:   keysExportRoot,
}

var cmdKeyImport = &cobra.Command{
	Use:   "import [ filename ]",
	Short: "Imports keys from a ZIP file.",
	Long:  "imports one or more keys from a ZIP file.",
	Run:   keysImport,
}

var cmdKeyImportRoot = &cobra.Command{
	Use:   "import-root [ filename ]",
	Short: "Imports root key.",
	Long:  "imports a root key from a PEM file.",
	Run:   keysImportRoot,
}

func keysList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}

	parseConfig()

	keysPath := filepath.Join(trustDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		fatalf("failed to create private key store in directory: %s", keysPath)
	}
	yubiStore, _ := api.NewYubiKeyStore(fileKeyStore, retriever)
	var cs signed.CryptoService
	if yubiStore == nil {
		cs = cryptoservice.NewCryptoService("", fileKeyStore)
	} else {
		cs = cryptoservice.NewCryptoService("", yubiStore, fileKeyStore)
	}

	// Get a map of all the keys/roles
	keysMap := cs.ListAllKeys()

	cmd.Println("")
	cmd.Println("# Root keys: ")
	for k, v := range keysMap {
		if v == "root" {
			cmd.Println(k)
		}
	}

	cmd.Println("")
	cmd.Println("# Signing keys: ")

	// Get a list of all the keys
	var sortedKeys []string
	for k := range keysMap {
		sortedKeys = append(sortedKeys, k)
	}
	// Sort the list of all the keys
	sort.Strings(sortedKeys)

	// Print a sorted list of the key/role
	for _, k := range sortedKeys {
		if keysMap[k] != "root" {
			printKey(cmd, k, keysMap[k])
		}
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

	keysPath := filepath.Join(trustDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		fatalf("failed to create private key store in directory: %s", keysPath)
	}
	yubiStore, err := api.NewYubiKeyStore(fileKeyStore, retriever)
	var cs signed.CryptoService
	if err != nil {
		cmd.Printf("No Yubikey detected, importing to local filesystem.")
		cs = cryptoservice.NewCryptoService("", fileKeyStore)
	} else {
		cs = cryptoservice.NewCryptoService("", yubiStore, fileKeyStore)
	}

	pubKey, err := cs.Create(data.CanonicalRootRole, algorithm)
	if err != nil {
		fatalf("failed to create a new root key: %v", err)
	}

	cmd.Printf("Generated new %s root key with keyID: %s\n", algorithm, pubKey.ID())
}

// keysExport exports a collection of keys to a ZIP file
func keysExport(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify output filename for export")
	}

	exportFilename := args[0]

	parseConfig()

	keysPath := filepath.Join(trustDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		fatalf("failed to create private key store in directory: %s", keysPath)
	}
	cs := cryptoservice.NewCryptoService("", fileKeyStore)

	exportFile, err := os.Create(exportFilename)
	if err != nil {
		fatalf("error creating output file: %v", err)
	}

	// Must use a different passphrase retriever to avoid caching the
	// unlocking passphrase and reusing that.
	exportRetriever := getRetriever()
	if keysExportGUN != "" {
		err = cs.ExportKeysByGUN(exportFile, keysExportGUN, exportRetriever)
	} else {
		err = cs.ExportAllKeys(exportFile, exportRetriever)
	}

	exportFile.Close()

	if err != nil {
		os.Remove(exportFilename)
		fatalf("error exporting keys: %v", err)
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

	if len(keyID) != idSize {
		fatalf("please specify a valid root key ID")
	}

	parseConfig()

	keysPath := filepath.Join(trustDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		fatalf("failed to create private key store in directory: %s", keysPath)
	}
	cs := cryptoservice.NewCryptoService("", fileKeyStore)

	exportFile, err := os.Create(exportFilename)
	if err != nil {
		fatalf("error creating output file: %v", err)
	}
	if keysExportRootChangePassphrase {
		// Must use a different passphrase retriever to avoid caching the
		// unlocking passphrase and reusing that.
		exportRetriever := getRetriever()
		err = cs.ExportRootKeyReencrypt(exportFile, keyID, exportRetriever)
	} else {
		err = cs.ExportRootKey(exportFile, keyID)
	}
	exportFile.Close()
	if err != nil {
		os.Remove(exportFilename)
		fatalf("error exporting root key: %v", err)
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

	keysPath := filepath.Join(trustDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		fatalf("failed to create private key store in directory: %s", keysPath)
	}
	cs := cryptoservice.NewCryptoService("", fileKeyStore)

	zipReader, err := zip.OpenReader(importFilename)
	if err != nil {
		fatalf("opening file for import: %v", err)
	}
	defer zipReader.Close()

	err = cs.ImportKeysZip(zipReader.Reader)

	if err != nil {
		fatalf("error importing keys: %v", err)
	}
}

// keysImportRoot imports a root key from a PEM file
func keysImportRoot(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		fatalf("must specify input filename for import")
	}

	importFilename := args[0]

	parseConfig()

	keysPath := filepath.Join(trustDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		fatalf("failed to create private key store in directory: %s", keysPath)
	}
	yubiStore, err := api.NewYubiKeyStore(fileKeyStore, retriever)
	var cs signed.CryptoService
	if err != nil {
		cmd.Printf("No Yubikey detected, importing to local filesystem.")
		cs = cryptoservice.NewCryptoService("", fileKeyStore)
	} else {
		cs = cryptoservice.NewCryptoService("", yubiStore, fileKeyStore)
	}

	importFile, err := os.Open(importFilename)
	if err != nil {
		fatalf("opening file for import: %v", err)
	}
	defer importFile.Close()

	err = cs.ImportRootKey(importFile)

	if err != nil {
		fatalf("error importing root key: %v", err)
	}
}

func printKey(cmd *cobra.Command, keyPath, alias string) {
	keyID := filepath.Base(keyPath)
	gun := filepath.Dir(keyPath)
	cmd.Printf("%s - %s - %s\n", gun, alias, keyID)
}

func keysRotate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	parseConfig()

	gun := args[0]
	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, remoteTrustServer, nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}
	if err := nRepo.RotateKeys(); err != nil {
		fatalf(err.Error())
	}
}
