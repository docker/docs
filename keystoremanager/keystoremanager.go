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

// ErrValidationFail is returned when there is no valid trusted certificates
// being served inside of the roots.json
type ErrValidationFail struct {
	Reason string
}

// ErrValidationFail is returned when there is no valid trusted certificates
// being served inside of the roots.json
func (err ErrValidationFail) Error() string {
	return fmt.Sprintf("could not validate the path to a trusted root: %s", err.Reason)
}

// ErrRootRotationFail is returned when we fail to do a full root key rotation
// by either failing to add the new root certificate, or delete the old ones
type ErrRootRotationFail struct {
	Reason string
}

// ErrRootRotationFail is returned when we fail to do a full root key rotation
// by either failing to add the new root certificate, or delete the old ones
func (err ErrRootRotationFail) Error() string {
	return fmt.Sprintf("could not rotate trust to a new trusted root: %s", err.Reason)
}

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
func (km *KeyStoreManager) ValidateRoot(root *data.Signed, gun string) error {
	logrus.Debugf("entered ValidateRoot with dns: %s", gun)
	signedRoot, err := data.RootFromSigned(root)
	if err != nil {
		return err
	}

	// Retrieve all the trusted certificates that match this gun
	certsForCN, err := km.trustedCertificateStore.GetCertificatesByCN(gun)
	if err != nil {
		// If the error that we get back is different than ErrNoCertificatesFound
		// we couldn't check if there are any certificates with this CN already
		// trusted. Let's take the conservative approach and return a failed validation
		if _, ok := err.(*trustmanager.ErrNoCertificatesFound); !ok {
			logrus.Debugf("error retrieving trusted certificates for: %s, %v", gun, err)
			return &ErrValidationFail{Reason: "unable to retrieve trusted certificates"}
		}
	}

	// Retrieve all the leaf certificates in which the CN matches the GUN
	allValidCerts, err := validRootLeafCerts(signedRoot, gun)
	if err != nil {
		logrus.Debugf("error retrieving valid leaf certificates for: %s, %v", gun, err)
		return &ErrValidationFail{Reason: "unable to retrieve valid leaf certificates"}
	}

	// If there are no certificates with this CN, lets TOFUS!
	// Note that this logic should only exist in docker 1.8
	var keysToCheck map[string]data.PublicKey
	if len(certsForCN) == 0 {
		// Convert all the valid certificates into public keys to do root validation
		keysToCheck = trustmanager.CertsToKeys(allValidCerts)
		logrus.Debugf("found no valid certificates for %s, entering TOFU mode", gun)
	} else {
		// Convert all the currently trusted certificates into public keys for root validation
		keysToCheck = trustmanager.CertsToKeys(certsForCN)
		logrus.Debugf("found valid certificates for %s", gun)
	}

	fmt.Printf("LEN: %d\n", len(certsForCN))

	// Verify the self consistency of the root keys that we're about to trust
	// or validate this new file with the old keys
	err = signed.VerifyRoot(root, 0, keysToCheck)
	if err != nil {
		logrus.Debugf("failed to verify TUF data for: %s, %v", gun, err)
		return &ErrValidationFail{Reason: "failed to validate integrity of roots"}
	}

	logrus.Debugf("validation of the root succeeded for: %s", gun)

	// We validated that this new root file has been signed by at least one
	// of our trusted keys, so let's only trust the currently available keys
	// First, we delete all the old certificates
	// rotationFailed keeps track of wether we have been able to delete all old certificates
	var rotationFailed bool
	for _, cert := range certsForCN {
		// Getting the CertID for debug only, we don't care about the error
		certID, _ := trustmanager.FingerprintCert(cert)
		logrus.Debugf("removing certificate with certID: %s", certID)

		err = km.trustedCertificateStore.RemoveCert(cert)
		if err != nil {
			logrus.Debugf("failed to remove trusted certificate with keyID: %s, %v", certID, err)
			rotationFailed = true
		}
	}

	// We've delete all the old certificates, let's add trust for all the new ones
	for _, cert := range allValidCerts {
		km.trustedCertificateStore.AddCert(cert)
		// Getting the CertID for debug only, don't care about the error
		certID, _ := trustmanager.FingerprintCert(cert)
		logrus.Debugf("adding trust certificate for with certID %s for %s", certID, gun)
	}

	// If we failed to delete any of the old certificates return an error.
	// We do it this way, so we're sure to first add the new trusted certificates
	// before we fail and make the repository unusable
	if rotationFailed {
		return &ErrRootRotationFail{Reason: "failed to rotate root keys"}
	}

	// Check to see if this leafCertificate has a chain to one of the Root
	// CAs of our CA Store.
	// err = trustmanager.Verify(km.trustedCAStore, dnsName, decodedCerts)
	// if err == nil {
	// 	validKeys[keyID] = signedRoot.Signed.Keys[keyID]
	// 	logrus.Debugf("found a CA path for %s with keyID: %s", dnsName, keyID)
	// }

	logrus.Debugf("Root validation succeeded")
	return nil
}

func validRootLeafCerts(root *data.SignedRoot, gun string) ([]*x509.Certificate, error) {
	// Get a list of all of the leaf certificates present in root
	allLeafCerts, _ := parseAllCerts(root)
	var validLeafCerts []*x509.Certificate

	// Go through every leaf certificate and check that the CN matches the gun
	for _, cert := range allLeafCerts {
		// Validate that this leaf certificate has a CN that matches the exact gun
		if cert.Subject.CommonName != gun {
			logrus.Debugf("error leaf certificate CN: %s doesn't match the given dns name: %s", cert.Subject.CommonName)
			continue
		}
		validLeafCerts = append(validLeafCerts, cert)
	}

	if len(validLeafCerts) < 1 {
		return validLeafCerts, errors.New("no valid leaf certificates found in any of the root keys")
	}

	return validLeafCerts, nil
}

func parseAllCerts(signedRoot *data.SignedRoot) (map[string]*x509.Certificate, map[string][]*x509.Certificate) {
	leafCerts := make(map[string]*x509.Certificate)
	intCerts := make(map[string][]*x509.Certificate)

	// Before we loop through all root keys available, make sure any exist
	rootRoles, ok := signedRoot.Signed.Roles["root"]
	if !ok {
		logrus.Debugf("tried to parse certificates from invalid root signed data")
		return nil, nil
	}

	logrus.Debugf("found the following root keys: %v", rootRoles.KeyIDs)
	// Iterate over every keyID for the root role inside of roots.json
	for _, keyID := range rootRoles.KeyIDs {
		// Decode all the x509 certificates that were bundled with this
		// Specific root key
		decodedCerts, err := trustmanager.LoadCertBundleFromPEM([]byte(signedRoot.Signed.Keys[keyID].Public()))
		if err != nil {
			logrus.Debugf("error while parsing root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// Get all non-CA certificates in the decoded certificates
		leafCertList := trustmanager.GetLeafCerts(decodedCerts)

		// If we got no leaf certificates or we got more than one, fail
		if len(leafCertList) != 1 {
			logrus.Debugf("invalid chain due to leaf certificate missing or too many leaf certificates for keyID: %s", keyID)
			continue
		}

		// Get the ID of the leaf certificate
		leafCert := leafCertList[0]
		leafID, err := trustmanager.FingerprintCert(leafCert)
		if err != nil {
			logrus.Debugf("error while fingerprinting root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// Store the leaf cert in the map
		leafCerts[leafID] = leafCert

		// Get all the remainder certificates marked as a CA to be used as intermediates
		intermediateCerts := trustmanager.GetIntermediateCerts(decodedCerts)
		intCerts[leafID] = intermediateCerts
	}

	return leafCerts, intCerts
}
