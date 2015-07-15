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

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

func moveKeysWithNewPassphrase(oldKeyStore, newKeyStore *trustmanager.KeyFileStore, outputPassphrase string) error {
	// List all files but no symlinks
	for _, f := range oldKeyStore.ListFiles(false) {
		fullKeyPath := strings.TrimSpace(strings.TrimSuffix(f, filepath.Ext(f)))
		relKeyPath := strings.TrimPrefix(fullKeyPath, oldKeyStore.BaseDir())
		relKeyPath = strings.TrimPrefix(relKeyPath, string(filepath.Separator))

		pemBytes, err := oldKeyStore.Get(relKeyPath)
		if err != nil {
			return err
		}

		block, _ := pem.Decode(pemBytes)
		if block == nil {
			return errors.New("no valid private key found")
		}

		if !x509.IsEncryptedPEMBlock(block) {
			// Key is not encrypted. Parse it, and add it
			// to the temporary store as an encrypted key.
			privKey, err := trustmanager.ParsePEMPrivateKey(pemBytes, "")
			if err != nil {
				return err
			}
			err = newKeyStore.AddEncryptedKey(relKeyPath, privKey, outputPassphrase)
		} else {
			// Encrypted key - pass it through without
			// decrypting
			err = newKeyStore.Add(relKeyPath, pemBytes)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func addKeysToArchive(zipWriter *zip.Writer, newKeyStore *trustmanager.KeyFileStore, tempBaseDir string, dedup map[string]struct{}) error {
	// List all files but no symlinks
	for _, fullKeyPath := range newKeyStore.ListFiles(false) {
		if _, present := dedup[fullKeyPath]; present {
			continue
		}
		dedup[fullKeyPath] = struct{}{}

		relKeyPath := strings.TrimPrefix(fullKeyPath, tempBaseDir)
		relKeyPath = strings.TrimPrefix(relKeyPath, string(filepath.Separator))

		fi, err := os.Stat(fullKeyPath)
		if err != nil {
			return err
		}

		infoHeader, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		infoHeader.Name = relKeyPath
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
// outputPassphrase is the new passphrase to use to encrypt the existing keys.
// If blank, the keys will not be encrypted. Note that keys which are already
// encrypted are not re-encrypted. They will be included in the zip with their
// original encryption.
func (km *KeyStoreManager) ExportAllKeys(dest io.Writer, outputPassphrase string) error {
	tempBaseDir, err := ioutil.TempDir("", "notary-key-export-")
	defer os.RemoveAll(tempBaseDir)

	// Create temporary keystores to use as a staging area
	tempNonRootKeysPath := filepath.Join(tempBaseDir, privDir)
	tempNonRootKeyStore, err := trustmanager.NewKeyFileStore(tempNonRootKeysPath)
	if err != nil {
		return err
	}

	tempRootKeysPath := filepath.Join(tempBaseDir, privDir, rootKeysSubdir)
	tempRootKeyStore, err := trustmanager.NewKeyFileStore(tempRootKeysPath)
	if err != nil {
		return err
	}

	if err := moveKeysWithNewPassphrase(km.rootKeyStore, tempRootKeyStore, outputPassphrase); err != nil {
		return err
	}
	if err := moveKeysWithNewPassphrase(km.nonRootKeyStore, tempNonRootKeyStore, outputPassphrase); err != nil {
		return err
	}

	zipWriter := zip.NewWriter(dest)

	// Root and non-root stores overlap, so we need to dedup files
	dedup := make(map[string]struct{})

	if err := addKeysToArchive(zipWriter, tempRootKeyStore, tempBaseDir, dedup); err != nil {
		return err
	}
	if err := addKeysToArchive(zipWriter, tempNonRootKeyStore, tempBaseDir, dedup); err != nil {
		return err
	}

	zipWriter.Close()

	return nil
}

// ImportKeysZip imports keys from a zip file provided as an io.ReaderAt. The
// keys in the root_keys directory are left encrypted, but the other keys are
// decrypted with the specified passphrase.
func (km *KeyStoreManager) ImportKeysZip(zipReader zip.Reader, passphrase string) error {
	// Temporarily store the keys in maps, so we can bail early if there's
	// an error (for example, wrong passphrase), without leaving the key
	// store in an inconsistent state
	newRootKeys := make(map[string][]byte)
	newNonRootKeys := make(map[string]*data.PrivateKey)

	// Iterate through the files in the archive. Don't add the keys
	for _, f := range zipReader.File {
		fNameTrimmed := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
		// Note that using / as a separator is okay here - the zip
		// package guarantees that the separator will be /
		keysPrefix := privDir + "/"

		if !strings.HasPrefix(fNameTrimmed, keysPrefix) {
			// This path inside the zip archive doesn't start with
			// "private". That's unexpected, because all keys
			// should be in that subdirectory. To avoid adding a
			// file to the filestore that we won't be able to use,
			// skip this file in the import.
			logrus.Warnf("skipping import of key with a path that doesn't begin with %s: %s", keysPrefix, f.Name)
			continue
		}
		fNameTrimmed = strings.TrimPrefix(fNameTrimmed, keysPrefix)

		rc, err := f.Open()
		if err != nil {
			return err
		}

		pemBytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil
		}

		// Is this in the root_keys directory?
		// Note that using / as a separator is okay here - the zip
		// package guarantees that the separator will be /
		rootKeysPrefix := rootKeysSubdir + "/"
		if strings.HasPrefix(fNameTrimmed, rootKeysPrefix) {
			// Root keys are preserved without decrypting
			keyName := strings.TrimPrefix(fNameTrimmed, rootKeysPrefix)
			newRootKeys[keyName] = pemBytes
		} else {
			// Non-root keys need to be decrypted
			key, err := trustmanager.ParsePEMPrivateKey(pemBytes, passphrase)
			if err != nil {
				return err
			}
			newNonRootKeys[fNameTrimmed] = key
		}

		rc.Close()
	}

	for keyName, pemBytes := range newRootKeys {
		if err := km.rootKeyStore.Add(keyName, pemBytes); err != nil {
			return err
		}
	}

	for keyName, privKey := range newNonRootKeys {
		if err := km.nonRootKeyStore.AddKey(keyName, privKey); err != nil {
			return err
		}
	}

	return nil
}
