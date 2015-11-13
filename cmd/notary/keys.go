package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"

	"github.com/docker/notary/tuf/data"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	cmdKey.AddCommand(cmdKeyList)
	cmdKey.AddCommand(cmdKeyGenerateRootKey)

	cmdKeyExport.Flags().StringVarP(&keysExportGUN, "gun", "g", "", "Globally Unique Name to export keys for")
	cmdKey.AddCommand(cmdKeyExport)
	cmdKey.AddCommand(cmdKeyExportRoot)
	cmdKeyExportRoot.Flags().BoolVarP(&keysExportRootChangePassphrase, "change-passphrase", "p", false, "Set a new passphrase for the key being exported")
	cmdKey.AddCommand(cmdKeyImport)
	cmdKey.AddCommand(cmdKeyImportRoot)
	cmdKey.AddCommand(cmdRotateKey)
}

var cmdKey = &cobra.Command{
	Use:   "key",
	Short: "Operates on keys.",
	Long:  `Operations on private keys.`,
}

var cmdKeyList = &cobra.Command{
	Use:   "list",
	Short: "Lists keys.",
	Long:  "Lists all keys known to notary.",
	Run:   keysList,
}

var cmdRotateKey = &cobra.Command{
	Use:   "rotate [ GUN ]",
	Short: "Rotate all the signing (non-root) keys for the given Globally Unique Name.",
	Long:  "Removes all old signing (non-root) keys for the given Globally Unique Name, and generates new ones.  This only makes local changes - please use then `notary publish` to push the key rotation changes to the remote server.",
	Run:   keysRotate,
}

var cmdKeyGenerateRootKey = &cobra.Command{
	Use:   "generate [ algorithm ]",
	Short: "Generates a new root key with a given algorithm.",
	Long:  "Generates a new root key with a given algorithm. If a hardware smartcard is available, the key will be stored both on hardware and on disk.  Please make sure to back up the key that is written to disk, and to then take the on-disk key offline.",
	Run:   keysGenerateRootKey,
}

var keysExportGUN string

var cmdKeyExport = &cobra.Command{
	Use:   "export [ filename ]",
	Short: "Exports keys to a ZIP file.",
	Long:  "Exports a collection of keys. The keys are reencrypted with a new passphrase. The output is a ZIP file.  If the --gun option is passed, only signing keys and no root keys will be exported.  Does not work on keys that are only in hardware (smartcards).",
	Run:   keysExport,
}

var keysExportRootChangePassphrase bool

var cmdKeyExportRoot = &cobra.Command{
	Use:   "export-root [ keyID ] [ filename ]",
	Short: "Exports given root key to a file.",
	Long:  "Exports a root key, without reencrypting. The output is a PEM file. Does not work on keys that are only in hardware (smartcards).",
	Run:   keysExportRoot,
}

var cmdKeyImport = &cobra.Command{
	Use:   "import [ filename ]",
	Short: "Imports keys from a ZIP file.",
	Long:  "Imports one or more keys from a ZIP file. If a hardware smartcard is available, the root key will be imported into the smartcard but not to disk.",
	Run:   keysImport,
}

var cmdKeyImportRoot = &cobra.Command{
	Use:   "import-root [ filename ]",
	Short: "Imports root key.",
	Long:  "Imports a root key from a PEM file. If a hardware smartcard is available, the root key will be imported into the smartcard but not to disk.",
	Run:   keysImportRoot,
}

func truncateWithEllipsis(str string, maxWidth int, leftTruncate bool) string {
	if len(str) <= maxWidth {
		return str
	}
	if leftTruncate {
		return fmt.Sprintf("...%s", str[len(str)-(maxWidth-3):])
	}
	return fmt.Sprintf("%s...", str[:maxWidth-3])
}

const (
	maxGUNWidth = 25
	maxLocWidth = 40
)

type keyInfo struct {
	gun      string // assumption that this is "" if role is root
	role     string
	keyID    string
	location string
}

// We want to sort by gun, then by role, then by keyID, then by location
// In the case of a root role, then there is no GUN, and a root role comes
// first.
type keyInfoSorter []keyInfo

func (k keyInfoSorter) Len() int      { return len(k) }
func (k keyInfoSorter) Swap(i, j int) { k[i], k[j] = k[j], k[i] }
func (k keyInfoSorter) Less(i, j int) bool {
	// special-case role
	if k[i].role != k[j].role {
		if k[i].role == data.CanonicalRootRole {
			return true
		}
		if k[j].role == data.CanonicalRootRole {
			return false
		}
		// otherwise, neither of them are root, they're just different, so
		// go with the traditional sort order.
	}

	// sort order is GUN, role, keyID, location.
	orderedI := []string{k[i].gun, k[i].role, k[i].keyID, k[i].location}
	orderedJ := []string{k[j].gun, k[j].role, k[j].keyID, k[j].location}

	for x := 0; x < 4; x++ {
		switch {
		case orderedI[x] < orderedJ[x]:
			return true
		case orderedI[x] > orderedJ[x]:
			return false
		}
		// continue on and evalulate the next item
	}
	// this shouldn't happen - that means two values are exactly equal
	return false
}

// Given a list of KeyStores in order of listing preference, pretty-prints the
// root keys and then the signing keys.
func prettyPrintKeys(keyStores []trustmanager.KeyStore, writer io.Writer) {
	var info []keyInfo

	for _, store := range keyStores {
		for keyPath, role := range store.ListKeys() {
			gun := ""
			if role != data.CanonicalRootRole {
				gun = filepath.Dir(keyPath)
			}
			info = append(info, keyInfo{
				role:     role,
				location: store.Name(),
				gun:      gun,
				keyID:    filepath.Base(keyPath),
			})
		}
	}
	if len(info) == 0 {
		writer.Write([]byte("No signing keys found.\n"))
		return
	}

	sort.Stable(keyInfoSorter(info))

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"ROLE", "GUN", "KEY ID", "LOCATION"})
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetAutoWrapText(false)

	for _, oneKeyInfo := range info {
		table.Append([]string{
			oneKeyInfo.role,
			truncateWithEllipsis(oneKeyInfo.gun, maxGUNWidth, true),
			oneKeyInfo.keyID,
			truncateWithEllipsis(oneKeyInfo.location, maxLocWidth, true),
		})
	}
	table.Render()
}

func keysList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(1)
	}

	parseConfig()

	stores := getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, true)
	cmd.Println("")
	prettyPrintKeys(stores, cmd.Out())
	cmd.Println("")
}

func keysGenerateRootKey(cmd *cobra.Command, args []string) {
	// We require one or no arguments (since we have a default value), but if the
	// user passes in more than one argument, we error out.
	if len(args) > 1 {
		cmd.Usage()
		fatalf("Please provide only one Algorithm as an argument to generate (rsa, ecdsa)")
	}

	parseConfig()

	// If no param is given to generate, generates an ecdsa key by default
	algorithm := data.ECDSAKey

	// If we were provided an argument lets attempt to use it as an algorithm
	if len(args) > 0 {
		algorithm = args[0]
	}

	allowedCiphers := map[string]bool{
		data.ECDSAKey: true,
		data.RSAKey:   true,
	}

	if !allowedCiphers[strings.ToLower(algorithm)] {
		fatalf("Algorithm not allowed, possible values are: RSA, ECDSA")
	}

	parseConfig()

	cs := cryptoservice.NewCryptoService(
		"",
		getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, true)...,
	)

	pubKey, err := cs.Create(data.CanonicalRootRole, algorithm)
	if err != nil {
		fatalf("Failed to create a new root key: %v", err)
	}

	cmd.Printf("Generated new %s root key with keyID: %s\n", algorithm, pubKey.ID())
}

// keysExport exports a collection of keys to a ZIP file
func keysExport(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify output filename for export")
	}

	parseConfig()
	exportFilename := args[0]

	cs := cryptoservice.NewCryptoService(
		"",
		getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, false)...,
	)

	exportFile, err := os.Create(exportFilename)
	if err != nil {
		fatalf("Error creating output file: %v", err)
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
		fatalf("Error exporting keys: %v", err)
	}
}

// keysExportRoot exports a root key by ID to a PEM file
func keysExportRoot(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("Must specify key ID and output filename for export")
	}

	parseConfig()

	keyID := args[0]
	exportFilename := args[1]

	if len(keyID) != idSize {
		fatalf("Please specify a valid root key ID")
	}

	parseConfig()

	cs := cryptoservice.NewCryptoService(
		"",
		getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, false)...,
	)

	exportFile, err := os.Create(exportFilename)
	if err != nil {
		fatalf("Error creating output file: %v", err)
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
		fatalf("Error exporting root key: %v", err)
	}
}

// keysImport imports keys from a ZIP file
func keysImport(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify input filename for import")
	}

	importFilename := args[0]

	parseConfig()

	cs := cryptoservice.NewCryptoService(
		"",
		getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, true)...,
	)

	zipReader, err := zip.OpenReader(importFilename)
	if err != nil {
		fatalf("Opening file for import: %v", err)
	}
	defer zipReader.Close()

	err = cs.ImportKeysZip(zipReader.Reader)

	if err != nil {
		fatalf("Error importing keys: %v", err)
	}
}

// keysImportRoot imports a root key from a PEM file
func keysImportRoot(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		fatalf("Must specify input filename for import")
	}

	parseConfig()

	cs := cryptoservice.NewCryptoService(
		"",
		getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, true)...,
	)

	importFilename := args[0]

	importFile, err := os.Open(importFilename)
	if err != nil {
		fatalf("Opening file for import: %v", err)
	}
	defer importFile.Close()

	err = cs.ImportRootKey(importFile)

	if err != nil {
		fatalf("Error importing root key: %v", err)
	}
}

func printKey(cmd *cobra.Command, keyPath, alias, loc string) {
	keyID := filepath.Base(keyPath)
	gun := filepath.Dir(keyPath)
	cmd.Printf("%s - %s - %s - %s\n", gun, alias, keyID, loc)
}

func keysRotate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN and target")
	}
	parseConfig()

	gun := args[0]
	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, remoteTrustServer, nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}
	if err := nRepo.RotateKeys(); err != nil {
		fatalf(err.Error())
	}
}

func getKeyStores(cmd *cobra.Command, directory string,
	ret passphrase.Retriever, withHardware bool) []trustmanager.KeyStore {

	fileKeyStore, err := trustmanager.NewKeyFileStore(directory, ret)
	if err != nil {
		fatalf("Failed to create private key store in directory: %s", directory)
	}

	ks := []trustmanager.KeyStore{fileKeyStore}

	if withHardware {
		yubiStore, err := getYubiKeyStore(fileKeyStore, ret)
		if err == nil && yubiStore != nil {
			// Note that the order is important, since we want to prioritize
			// the yubikey store
			ks = []trustmanager.KeyStore{yubiStore, fileKeyStore}
		}
	}

	return ks
}
