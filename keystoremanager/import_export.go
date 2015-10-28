package keystoremanager

import (
	"archive/zip"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/notary/pkg/passphrase"
	"github.com/docker/notary/trustmanager"
)

const zipMadeByUNIX = 3 << 8

var (
	// ErrNoValidPrivateKey is returned if a key being imported doesn't
	// look like a private key
	ErrNoValidPrivateKey = errors.New("no valid private key found")

	// ErrRootKeyNotEncrypted is returned if a root key being imported is
	// unencrypted
	ErrRootKeyNotEncrypted = errors.New("only encrypted root keys may be imported")

	// ErrNoKeysFoundForGUN is returned if no keys are found for the
	// specified GUN during export
	ErrNoKeysFoundForGUN = errors.New("no keys found for specified GUN")
)

// ExportRootKey exports the specified root key to an io.Writer in PEM format.
// The key's existing encryption is preserved.
func (km *KeyStoreManager) ExportRootKey(dest io.Writer, keyID string) error {
	pemBytes, err := km.KeyStore.ExportKey(keyID)
	if err != nil {
		return err
	}
	nBytes, err := dest.Write(pemBytes)
	if err != nil {
		return err
	}
	if nBytes != len(pemBytes) {
		return errors.New("Unable to finish writing exported key.")
	}
	return nil
}

// ExportRootKeyReencrypt exports the specified root key to an io.Writer in
// PEM format. The key is reencrypted with a new passphrase.
func (km *KeyStoreManager) ExportRootKeyReencrypt(dest io.Writer, keyID string, newPassphraseRetriever passphrase.Retriever) error {
	privateKey, alias, err := km.KeyStore.GetKey(keyID)
	if err != nil {
		return err
	}

	// Create temporary keystore to use as a staging area
	tempBaseDir, err := ioutil.TempDir("", "notary-key-export-")
	defer os.RemoveAll(tempBaseDir)

	tempKeysPath := filepath.Join(tempBaseDir, privDir)
	tempKeyStore, err := trustmanager.NewKeyFileStore(tempKeysPath, newPassphraseRetriever)
	if err != nil {
		return err
	}

	err = tempKeyStore.AddKey(keyID, alias, privateKey)
	if err != nil {
		return err
	}

	pemBytes, err := tempKeyStore.ExportKey(keyID)
	if err != nil {
		return err
	}
	nBytes, err := dest.Write(pemBytes)
	if err != nil {
		return err
	}
	if nBytes != len(pemBytes) {
		return errors.New("Unable to finish writing exported key.")
	}
	return nil
}

// checkRootKeyIsEncrypted makes sure the root key is encrypted. We have
// internal assumptions that depend on this.
func checkRootKeyIsEncrypted(pemBytes []byte) error {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return ErrNoValidPrivateKey
	}

	if !x509.IsEncryptedPEMBlock(block) {
		return ErrRootKeyNotEncrypted
	}

	return nil
}

// ImportRootKey imports a root in PEM format key from an io.Reader
// The key's existing encryption is preserved. The keyID parameter is
// necessary because otherwise we'd need the passphrase to decrypt the key
// in order to compute the ID.
func (km *KeyStoreManager) ImportRootKey(source io.Reader, keyID string) error {
	pemBytes, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}

	if err = checkRootKeyIsEncrypted(pemBytes); err != nil {
		return err
	}

	if err = km.KeyStore.ImportKey(pemBytes, "root"); err != nil {
		return err
	}

	return err
}

func moveKeys(oldKeyStore, newKeyStore *trustmanager.KeyFileStore) error {
	for f := range oldKeyStore.ListKeys() {
		privateKey, alias, err := oldKeyStore.GetKey(f)
		if err != nil {
			return err
		}

		err = newKeyStore.AddKey(f, alias, privateKey)

		if err != nil {
			return err
		}
	}

	return nil
}

func addKeysToArchive(zipWriter *zip.Writer, newKeyStore *trustmanager.KeyFileStore, subDir string) error {
	for _, relKeyPath := range newKeyStore.ListFiles() {
		fullKeyPath := filepath.Join(newKeyStore.BaseDir(), relKeyPath)

		fi, err := os.Lstat(fullKeyPath)
		if err != nil {
			return err
		}

		infoHeader, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		infoHeader.Name = filepath.Join(subDir, relKeyPath)

		zipFileEntryWriter, err := zipWriter.CreateHeader(infoHeader)
		if err != nil {
			return err
		}

		fileContents, err := ioutil.ReadFile(fullKeyPath)
		if err != nil {
			return err
		}

		if _, err = zipFileEntryWriter.Write(fileContents); err != nil {
			return err
		}
	}

	return nil
}

// ExportAllKeys exports all keys to an io.Writer in zip format.
// newPassphraseRetriever will be used to obtain passphrases to use to encrypt the existing keys.
func (km *KeyStoreManager) ExportAllKeys(dest io.Writer, newPassphraseRetriever passphrase.Retriever) error {
	tempBaseDir, err := ioutil.TempDir("", "notary-key-export-")
	defer os.RemoveAll(tempBaseDir)

	// Create temporary keystore to use as a staging area
	tempKeysPath := filepath.Join(tempBaseDir, privDir)
	tempKeyStore, err := trustmanager.NewKeyFileStore(tempKeysPath, newPassphraseRetriever)
	if err != nil {
		return err
	}

	if err := moveKeys(km.KeyStore, tempKeyStore); err != nil {
		return err
	}

	zipWriter := zip.NewWriter(dest)

	if err := addKeysToArchive(zipWriter, tempKeyStore, privDir); err != nil {
		return err
	}

	zipWriter.Close()

	return nil
}

// ImportKeysZip imports keys from a zip file provided as an zip.Reader. The
// keys in the root_keys directory are left encrypted, but the other keys are
// decrypted with the specified passphrase.
func (km *KeyStoreManager) ImportKeysZip(zipReader zip.Reader) error {
	// Temporarily store the keys in maps, so we can bail early if there's
	// an error (for example, wrong passphrase), without leaving the key
	// store in an inconsistent state
	newKeys := make(map[string][]byte)

	// Iterate through the files in the archive. Don't add the keys
	for _, f := range zipReader.File {
		fNameTrimmed := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))

		rc, err := f.Open()
		if err != nil {
			return err
		}

		fileBytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil
		}

		// Note that using / as a separator is okay here - the zip
		// package guarantees that the separator will be /
		if strings.HasPrefix(fNameTrimmed, privDir) {
			if fNameTrimmed[len(fNameTrimmed)-5:] == "_root" {
				if err = checkRootKeyIsEncrypted(fileBytes); err != nil {
					rc.Close()
					return err
				}
			}
			keyName := strings.TrimPrefix(fNameTrimmed, privDir)
			newKeys[keyName] = fileBytes
		} else {
			// This path inside the zip archive doesn't look like a
			// root key, non-root key, or alias. To avoid adding a file
			// to the filestore that we won't be able to use, skip
			// this file in the import.
			rc.Close()
			continue
		}

		rc.Close()
	}

	for keyName, pemBytes := range newKeys {
		if err := km.KeyStore.Add(keyName, pemBytes); err != nil {
			return err
		}
	}

	return nil
}

func moveKeysByGUN(oldKeyStore, newKeyStore *trustmanager.KeyFileStore, gun string) error {
	for relKeyPath := range oldKeyStore.ListKeys() {
		// Skip keys that aren't associated with this GUN
		if !strings.HasPrefix(relKeyPath, filepath.FromSlash(gun)) {
			continue
		}

		privKey, alias, err := oldKeyStore.GetKey(relKeyPath)
		if err != nil {
			return err
		}

		err = newKeyStore.AddKey(relKeyPath, alias, privKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExportKeysByGUN exports all keys associated with a specified GUN to an
// io.Writer in zip format. passphraseRetriever is used to select new passphrases to use to
// encrypt the keys.
func (km *KeyStoreManager) ExportKeysByGUN(dest io.Writer, gun string, passphraseRetriever passphrase.Retriever) error {
	tempBaseDir, err := ioutil.TempDir("", "notary-key-export-")
	defer os.RemoveAll(tempBaseDir)

	// Create temporary keystore to use as a staging area
	tempKeysPath := filepath.Join(tempBaseDir, privDir)
	tempKeyStore, err := trustmanager.NewKeyFileStore(tempKeysPath, passphraseRetriever)
	if err != nil {
		return err
	}

	if err := moveKeysByGUN(km.KeyStore, tempKeyStore, gun); err != nil {
		return err
	}

	zipWriter := zip.NewWriter(dest)

	if len(tempKeyStore.ListKeys()) == 0 {
		return ErrNoKeysFoundForGUN
	}

	if err := addKeysToArchive(zipWriter, tempKeyStore, privDir); err != nil {
		return err
	}

	zipWriter.Close()

	return nil
}
