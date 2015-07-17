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

var (
	// ErrValidationFail is returned when there is no trusted certificate in any of the
	// root keys available in the roots.json
	ErrValidationFail = errors.New("could not validate the path to a trusted root")
	// ErrRootRotationFail is returned when we fail to do a full root key rotation
	// by either failing to add the new root certificate, or delete the old ones
	ErrRootRotationFail = errors.New("could not rotate trust to a new trusted root")
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

Currently this method operates on a Trust On First Use (TOFU) model: if we
have never seen a certificate for a particular CN, we trust it. If later we see
a different certificate for that certificate, we return an ErrValidationFailed error.

Note that since we only allow trust data to be downloaded over an HTTPS channel
we are using the current web-of-trust to validate the first download of the certificate
adding an extra layer of security over the normal (SSH style) trust model.
We shall call this: TOFUS.

ValidateRoot also supports root key rotation, trusting a new certificate that has
been included in the roots.json, and removing trust in the old one.
*/
func (km *KeyStoreManager) ValidateRoot(root *data.Signed, dnsName string) error {
	logrus.Debugf("entered ValidateRoot with dns: %s", dnsName)
	rootSigned, err := data.RootFromSigned(root)
	if err != nil {
		return err
	}

	// validKeys will store all the keys that were considered valid either by
	// direct certificate match, or CA chain path
	validKeys := make(map[string]*data.PublicKey)

	// allCerts will keep a list of all leafCerts that were found, and is used
	// to aid on root certificate rotation
	allCerts := make(map[string]*x509.Certificate)

	// Before we loop through all root keys available, make sure any exist
	rootRoles, ok := rootSigned.Signed.Roles["root"]
	if !ok {
		return errors.New("no root roles found in tuf metadata")
	}

	logrus.Debugf("found the following root keys in roots.json: %v", rootRoles.KeyIDs)
	// Iterate over every keyID for the root role inside of roots.json
	for _, keyID := range rootRoles.KeyIDs {
		// Decode all the x509 certificates that were bundled with this
		// Specific root key
		decodedCerts, err := trustmanager.LoadCertBundleFromPEM([]byte(rootSigned.Signed.Keys[keyID].Public()))
		if err != nil {
			logrus.Debugf("error while parsing root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// Get all non-CA certificates in the decoded certificates
		leafCerts := trustmanager.GetLeafCerts(decodedCerts)

		// If we got no leaf certificates or we got more than one, fail
		if len(leafCerts) != 1 {
			logrus.Debugf("wasn't able to find a leaf certificate in the chain of keyID: %s", keyID)
			continue
		}

		// Get the ID of the leaf certificate
		leafCert := leafCerts[0]
		leafID, err := trustmanager.FingerprintCert(leafCert)
		if err != nil {
			logrus.Debugf("error while fingerprinting root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// Validate that this leaf certificate has a CN that matches the exact gun
		if leafCert.Subject.CommonName != dnsName {
			logrus.Debugf("error leaf certificate CN: %s doesn't match the given dns name: %s", leafCert.Subject.CommonName, dnsName)
			continue
		}

		// Add all the valid leafs to the certificates map so we can refer to them later
		allCerts[leafID] = leafCert

		// Retrieve all the trusted certificates that match this dns Name
		certsForCN, err := km.certificateStore.GetCertificatesByCN(dnsName)
		if err != nil {
			// If the error that we get back is different than ErrNoCertificatesFound
			// we couldn't check if there are any certificates with this CN already
			// trusted. Let's take the conservative approach and not trust this key
			if _, ok := err.(*trustmanager.ErrNoCertificatesFound); !ok {
				logrus.Debugf("error retrieving certificates for: %s, %v", dnsName, err)
				continue
			}
		}

		// If there are no certificates with this CN, lets TOFUS!
		// Note that this logic should only exist in docker 1.8
		if len(certsForCN) == 0 {
			km.certificateStore.AddCert(leafCert)
			certsForCN = append(certsForCN, leafCert)
			logrus.Debugf("using TOFUS on %s with keyID: %s", dnsName, leafID)
		}

		// Iterate over all known certificates for this CN and see if any are trusted
		for _, cert := range certsForCN {
			// Check to see if there is an exact match of this certificate.
			certID, err := trustmanager.FingerprintCert(cert)
			if err == nil && certID == leafID {
				validKeys[keyID] = rootSigned.Signed.Keys[keyID]
				logrus.Debugf("found an exact match for %s with keyID: %s", dnsName, keyID)
			}
		}

		// Check to see if this leafCertificate has a chain to one of the Root
		// CAs of our CA Store.
		err = trustmanager.Verify(km.caStore, dnsName, decodedCerts)
		if err == nil {
			validKeys[keyID] = rootSigned.Signed.Keys[keyID]
			logrus.Debugf("found a CA path for %s with keyID: %s", dnsName, keyID)
		}
	}

	if len(validKeys) < 1 {
		logrus.Debugf("wasn't able to trust any of the root keys")
		return ErrValidationFail
	}

	// TODO(david): change hardcoded minversion on TUF.
	newRootKey, err := signed.VerifyRoot(root, 0, validKeys, 1)
	if err != nil {
		return err
	}

	// VerifyRoot returns a non-nil value if there is a root key rotation happening.
	// If this happens, we should replace the old root of trust with the new one
	if newRootKey != nil {
		logrus.Debugf("got a new root key to rotate to: %s", newRootKey.ID())

		// Retrieve the certificate associated with the new root key and trust it
		newRootKeyCert, ok := allCerts[newRootKey.ID()]
		// Paranoid check for the certificate still being in the map
		if !ok {
			logrus.Debugf("error while retrieving new root certificate with keyID: %s, %v", newRootKey.ID(), err)
			return ErrRootRotationFail
		}

		// Add the new root certificate to our certificate store
		err := km.certificateStore.AddCert(newRootKeyCert)
		if err != nil {
			// Ignore the error if the certificate already exists
			if _, ok := err.(*trustmanager.ErrCertExists); !ok {
				logrus.Debugf("error while adding new root certificate with keyID: %s, %v", newRootKey.ID(), err)
				return ErrRootRotationFail
			}
			logrus.Debugf("root certificate already exists in keystore: %s", newRootKey.ID())
		}

		// Remove the new root certificate from the certificate mapping so we
		// can remove trust from all of the remaining ones
		delete(allCerts, newRootKey.ID())

		// Iterate over all old valid certificates and remove them, essentially
		// finishing the rotation of the currently trusted root certificate
		for _, cert := range allCerts {
			err := km.certificateStore.RemoveCert(cert)
			if err != nil {
				logrus.Debugf("error while removing old root certificate: %v", err)
				return ErrRootRotationFail
			}
			logrus.Debugf("removed trust from old root certificate")
		}
	}

	logrus.Debugf("Root validation succeeded")
	return nil
}
