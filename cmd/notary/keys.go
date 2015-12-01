package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"

	"github.com/docker/notary/tuf/data"
	"github.com/spf13/cobra"
)

func init() {
	cmdKey.AddCommand(cmdKeyList)
	cmdKey.AddCommand(cmdKeyGenerateRootKey)

	cmdKeysBackup.Flags().StringVarP(&keysExportGUN, "gun", "g", "", "Globally Unique Name to export keys for")
	cmdKey.AddCommand(cmdKeysBackup)
	cmdKey.AddCommand(cmdKeyExportRoot)
	cmdKeyExportRoot.Flags().BoolVarP(&keysExportRootChangePassphrase, "change-passphrase", "p", false, "Set a new passphrase for the key being exported")
	cmdKey.AddCommand(cmdKeysRestore)
	cmdKey.AddCommand(cmdKeyImportRoot)
	cmdKey.AddCommand(cmdRotateKey)
	cmdKey.AddCommand(cmdKeyRemove)
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
	Long:  "Generates a new root key with a given algorithm. If hardware key storage (e.g. a Yubikey) is available, the key will be stored both on hardware and on disk (so that it can be backed up).  Please make sure to back up and then remove this on-key disk immediately afterwards.",
	Run:   keysGenerateRootKey,
}

var keysExportGUN string

var cmdKeysBackup = &cobra.Command{
	Use:   "backup [ zipfilename ]",
	Short: "Backs up all your on-disk keys to a ZIP file.",
	Long:  "Backs up all of your accessible of keys. The keys are reencrypted with a new passphrase. The output is a ZIP file.  If the --gun option is passed, only signing keys and no root keys will be backed up.  Does not work on keys that are only in hardware (e.g. Yubikeys).",
	Run:   keysBackup,
}

var keysExportRootChangePassphrase bool

var cmdKeyExportRoot = &cobra.Command{
	Use:   "export [ keyID ] [ pemfilename ]",
	Short: "Export a root key on disk to a PEM file.",
	Long:  "Exports a single root key on disk, without reencrypting. The output is a PEM file. Does not work on keys that are only in hardware (e.g. Yubikeys).",
	Run:   keysExportRoot,
}

var cmdKeysRestore = &cobra.Command{
	Use:   "restore [ zipfilename ]",
	Short: "Restore multiple keys from a ZIP file.",
	Long:  "Restores one or more keys from a ZIP file. If hardware key storage (e.g. a Yubikey) is available, root keys will be imported into the hardware, but not backed up to disk in the same location as the other, non-root keys.",
	Run:   keysRestore,
}

var cmdKeyImportRoot = &cobra.Command{
	Use:   "import [ pemfilename ]",
	Short: "Imports a root key from a PEM file.",
	Long:  "Imports a single root key from a PEM file. If a hardware key storage (e.g. Yubikey) is available, the root key will be imported into the hardware but not backed up on disk again.",
	Run:   keysImportRoot,
}

var cmdKeyRemove = &cobra.Command{
	Use:   "remove [ keyID ]",
	Short: "Removes the key with the given keyID.",
	Long:  "Removes the key with the given keyID.  If the key is stored in more than one location, you will be asked which one to remove.",
	Run:   keyRemove,
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

// keysBackup exports a collection of keys to a ZIP file
func keysBackup(cmd *cobra.Command, args []string) {
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

// keysRestore imports keys from a ZIP file
func keysRestore(cmd *cobra.Command, args []string) {
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

func removeKeyInteractively(keyStores []trustmanager.KeyStore, keyID string,
	in io.Reader, out io.Writer) error {

	var foundKeys [][]string
	var storesByIndex []trustmanager.KeyStore

	for _, store := range keyStores {
		for keypath, role := range store.ListKeys() {
			if filepath.Base(keypath) == keyID {
				foundKeys = append(foundKeys,
					[]string{keypath, role, store.Name()})
				storesByIndex = append(storesByIndex, store)
			}
		}
	}

	if len(foundKeys) == 0 {
		return fmt.Errorf("No key with ID %s found.", keyID)
	}

	readIn := bufio.NewReader(in)

	if len(foundKeys) > 1 {
		for {
			// ask the user for which key to delete
			fmt.Fprintf(out, "Found the following matching keys:\n")
			for i, info := range foundKeys {
				fmt.Fprintf(out, "\t%d. %s: %s (%s)\n", i+1, info[0], info[1], info[2])
			}
			fmt.Fprint(out, "Which would you like to delete?  Please enter a number:  ")
			result, err := readIn.ReadBytes('\n')
			if err != nil {
				return err
			}
			index, err := strconv.Atoi(strings.TrimSpace(string(result)))

			if err != nil || index > len(foundKeys) || index < 1 {
				fmt.Fprintf(out, "\nInvalid choice: %s\n", string(result))
				continue
			}
			foundKeys = [][]string{foundKeys[index-1]}
			storesByIndex = []trustmanager.KeyStore{storesByIndex[index-1]}
			fmt.Fprintln(out, "")
			break
		}
	}
	// Now the length must be 1 - ask for confirmation.
	keyDescription := fmt.Sprintf("%s (role %s) from %s", foundKeys[0][0],
		foundKeys[0][1], foundKeys[0][2])

	fmt.Fprintf(out, "Are you sure you want to remove %s?  [Y/n]  ",
		keyDescription)
	result, err := readIn.ReadBytes('\n')
	if err != nil {
		return err
	}
	yesno := strings.ToLower(strings.TrimSpace(string(result)))

	if !strings.HasPrefix("yes", yesno) && yesno != "" {
		fmt.Fprintln(out, "\nAborting action.")
		return nil
	}

	err = storesByIndex[0].RemoveKey(foundKeys[0][0])
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "\nDeleted %s.\n", keyDescription)
	return nil
}

// keyRemove deletes a private key based on ID
func keyRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify the key ID of the key to remove")
	}

	parseConfig()
	keyID := args[0]

	// This is an invalid ID
	if len(keyID) != idSize {
		fatalf("invalid key ID provided: %s", keyID)
	}
	stores := getKeyStores(cmd, mainViper.GetString("trust_dir"), retriever, true)
	cmd.Println("")
	err := removeKeyInteractively(stores, keyID, os.Stdin, cmd.Out())
	cmd.Println("")
	if err != nil {
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
