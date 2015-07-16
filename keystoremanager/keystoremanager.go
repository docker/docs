package keystoremanager

import (
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
)

// KeyStoreManager is an abstraction around the root and non-root key stores,
// and related CA stores
type KeyStoreManager struct {
	rootKeyStore    *trustmanager.KeyFileStore
	nonRootKeyStore *trustmanager.KeyFileStore

	trustedCAStore          trustmanager.X509Store
	trustedCertificateStore trustmanager.X509Store
}

const (
	trustDir          = "trusted_certificates"
	privDir           = "private"
	rootKeysSubdir    = "root_keys"
	nonRootKeysSubdir = "tuf_keys"
	rsaRootKeySize    = 4096 // Used for new root keys
)

// NewKeyStoreManager returns an initialized KeyStoreManager, or an error
// if it fails to create the KeyFileStores or load certificates
func NewKeyStoreManager(baseDir string) (*KeyStoreManager, error) {
	nonRootKeysPath := filepath.Join(baseDir, privDir, nonRootKeysSubdir)
	nonRootKeyStore, err := trustmanager.NewKeyFileStore(nonRootKeysPath)
	if err != nil {
		return nil, err
	}

	// Load the keystore that will hold all of our encrypted Root Private Keys
	rootKeysPath := filepath.Join(baseDir, privDir, rootKeysSubdir)
	rootKeyStore, err := trustmanager.NewKeyFileStore(rootKeysPath)
	if err != nil {
		return nil, err
	}

	trustPath := filepath.Join(baseDir, trustDir)

	// Load all CAs that aren't expired and don't use SHA1
	trustedCAStore, err := trustmanager.NewX509FilteredFileStore(trustPath, func(cert *x509.Certificate) bool {
		return cert.IsCA && cert.BasicConstraintsValid && cert.SubjectKeyId != nil &&
			time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
	if err != nil {
		return nil, err
	}

	// Load all individual (non-CA) certificates that aren't expired and don't use SHA1
	trustedCertificateStore, err := trustmanager.NewX509FilteredFileStore(trustPath, func(cert *x509.Certificate) bool {
		return !cert.IsCA &&
			time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
	if err != nil {
		return nil, err
	}

	return &KeyStoreManager{
		rootKeyStore:            rootKeyStore,
		nonRootKeyStore:         nonRootKeyStore,
		trustedCAStore:          trustedCAStore,
		trustedCertificateStore: trustedCertificateStore,
	}, nil
}

// RootKeyStore returns the root key store being managed by this
// KeyStoreManager
func (km *KeyStoreManager) RootKeyStore() *trustmanager.KeyFileStore {
	return km.rootKeyStore
}

// NonRootKeyStore returns the non-root key store being managed by this
// KeyStoreManager
func (km *KeyStoreManager) NonRootKeyStore() *trustmanager.KeyFileStore {
	return km.nonRootKeyStore
}

// TrustedCertificateStore returns the trusted certificate store being managed
// by this KeyStoreManager
func (km *KeyStoreManager) TrustedCertificateStore() trustmanager.X509Store {
	return km.trustedCertificateStore
}

// TrustedCAStore returns the CA store being managed by this KeyStoreManager
func (km *KeyStoreManager) TrustedCAStore() trustmanager.X509Store {
	return km.trustedCAStore
}

// AddTrustedCert adds a cert to the trusted certificate store (not the CA
// store)
func (km *KeyStoreManager) AddTrustedCert(cert *x509.Certificate) {
	km.trustedCertificateStore.AddCert(cert)
}

// AddTrustedCACert adds a cert to the trusted CA certificate store
func (km *KeyStoreManager) AddTrustedCACert(cert *x509.Certificate) {
	km.trustedCAStore.AddCert(cert)
}

// GenRootKey generates a new root key protected by a given passphrase
// TODO(diogo): show not create keys manually, should use a cryptoservice instead
func (km *KeyStoreManager) GenRootKey(algorithm, passphrase string) (string, error) {
	var err error
	var privKey data.PrivateKey

	// We don't want external API callers to rely on internal TUF data types, so
	// the API here should continue to receive a string algorithm, and ensure
	// that it is downcased
	switch data.KeyAlgorithm(strings.ToLower(algorithm)) {
	case data.RSAKey:
		privKey, err = trustmanager.GenerateRSAKey(rand.Reader, rsaRootKeySize)
	case data.ECDSAKey:
		privKey, err = trustmanager.GenerateECDSAKey(rand.Reader)
	default:
		return "", fmt.Errorf("only RSA or ECDSA keys are currently supported. Found: %s", algorithm)

	}
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// Changing the root
	km.rootKeyStore.AddEncryptedKey(privKey.ID(), privKey, passphrase)

	return privKey.ID(), nil
}

// GetRootCryptoService retreives a root key and a cryptoservice to use with it
func (km *KeyStoreManager) GetRootCryptoService(rootKeyID, passphrase string) (*cryptoservice.UnlockedCryptoService, error) {
	privKey, err := km.rootKeyStore.GetDecryptedKey(rootKeyID, passphrase)
	if err != nil {
		return nil, fmt.Errorf("could not get decrypted root key with keyID: %s, %v", rootKeyID, err)
	}

	cryptoService := cryptoservice.NewCryptoService("", km.rootKeyStore, passphrase)

	return cryptoservice.NewUnlockedCryptoService(privKey, cryptoService), nil
}

/*
ValidateRoot iterates over every root key included in the TUF data and
attempts to validate the certificate by first checking for an exact match on
the certificate store, and subsequently trying to find a valid chain on the
trustedCAStore.

When this is being used with a notary repository, the dnsName parameter should
be the GUN associated with the repository.

Example TUF Content for root role:
"roles" : {
  "root" : {
    "threshold" : 1,
      "keyids" : [
        "e6da5c303d572712a086e669ecd4df7b785adfc844e0c9a7b1f21a7dfc477a38"
      ]
  },
 ...
}

Example TUF Content for root key:
"e6da5c303d572712a086e669ecd4df7b785adfc844e0c9a7b1f21a7dfc477a38" : {
	"keytype" : "RSA",
	"keyval" : {
	  "private" : "",
	  "public" : "Base64-encoded, PEM encoded x509 Certificate"
	}
}
*/
func (km *KeyStoreManager) ValidateRoot(root *data.Signed, dnsName string) error {
	rootSigned, err := data.RootFromSigned(root)
	if err != nil {
		return err
	}

	// iterate over every keyID for the root role inside of roots.json
	validKeys := make(map[string]*data.PublicKey)
	for _, keyID := range rootSigned.Signed.Roles["root"].KeyIDs {
		// decode all the x509 certificates that were bundled with this
		// specific root key
		decodedCerts, err := trustmanager.LoadCertBundleFromPEM([]byte(rootSigned.Signed.Keys[keyID].Public()))
		if err != nil {
			logrus.Debugf("error while parsing root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// get all non-CA certificates in the decoded certificates
		leafCerts := trustmanager.GetLeafCerts(decodedCerts)

		// gf we got no leaf certificates or we got more than one, fail
		if len(leafCerts) != 1 {
			logrus.Debugf("error while parsing root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// get the ID of the leaf certificate
		leafCert := leafCerts[0]
		leafID, err := trustmanager.FingerprintCert(leafCert)
		if err != nil {
			logrus.Debugf("error while fingerprinting root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// retrieve all the trusted certificates that match this dns Name
		certsForCN, err := km.certificateStore.GetCertificatesByCN(dnsName)

		// if there are no certificates with this CN, lets TOFU!
		// note that this logic should only exist in docker 1.8
		if len(certsForCN) == 0 {
			km.certificateStore.AddCert(leafCert)
			certsForCN = append(certsForCN, leafCert)
		}

		// iterate over all known certificates for this CN and see if any are trusted
		for _, cert := range certsForCN {
			// Check to see if there is an exact match of this certificate.
			certID, err := trustmanager.FingerprintCert(cert)
			if err == nil && certID == leafID {
				validKeys[keyID] = rootSigned.Signed.Keys[keyID]
			}
		}

		// Check to see if this leafCertificate has a chain to one of the Root
		// CAs of our CA Store.
		err = trustmanager.Verify(km.caStore, dnsName, decodedCerts)
		if err == nil {
			validKeys[keyID] = rootSigned.Signed.Keys[keyID]
		}
	}

	if len(validKeys) < 1 {
		return errors.New("could not validate the path to a trusted root")
	}

	// TODO(david): change hardcoded minversion on TUF.
	newRootKey, err := signed.VerifyRoot(root, 0, validKeys, 1)
	if err != nil {
		return err
	}

	// VerifyRoot returns a non-nil value if there is a root key rotation happening
	// if this happens, we should replace the old root of trust with the new one
	if newRootKey != nil {
		// retrieve all the certificates associated with the new root key
		keyID := newRootKey.ID()
		decodedCerts, err := trustmanager.LoadCertBundleFromPEM([]byte(rootSigned.Signed.Keys[keyID].Public()))
		if err != nil {
			logrus.Debugf("error while parsing root certificate with keyID: %s, %v", keyID, err)
			return err
		}

		// adds trust on the certificate of the new root key
		leafCerts := trustmanager.GetLeafCerts(decodedCerts)
		err = km.certificateStore.AddCert(leafCerts[0])
		if err != nil {
			return err
		}

		// iterate over all old valid keys and removes the associated certificates
		// were previously valid
		for _, key := range validKeys {
			cert, err := km.certificateStore.GetCertificateByCertID(key.ID())
			if err != nil {
				return err
			}
			// Remove the old certificate
			km.certificateStore.RemoveCert(cert)
		}
	}

	return nil
}
