package keystoremanager

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
)

// KeyStoreManager is an abstraction around the root and non-root key stores,
// and related CA stores
type KeyStoreManager struct {
	rootKeyStore    *trustmanager.KeyFileStore
	nonRootKeyStore *trustmanager.KeyFileStore

	caStore          trustmanager.X509Store
	certificateStore trustmanager.X509Store
}

const (
	trustDir       = "trusted_certificates"
	privDir        = "private"
	rootKeysSubdir = "root_keys"
)

// NewKeyStoreManager returns an initialized KeyStoreManager, or an error
// if it fails to create the KeyFileStores or load certificates
func NewKeyStoreManager(baseDir string) (*KeyStoreManager, error) {
	nonRootKeyStore, err := trustmanager.NewKeyFileStore(filepath.Join(baseDir, privDir))
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
	caStore, err := trustmanager.NewX509FilteredFileStore(trustPath, func(cert *x509.Certificate) bool {
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
	certificateStore, err := trustmanager.NewX509FilteredFileStore(trustPath, func(cert *x509.Certificate) bool {
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
		rootKeyStore:     rootKeyStore,
		nonRootKeyStore:  nonRootKeyStore,
		caStore:          caStore,
		certificateStore: certificateStore,
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

// CertificateStore returns the certificate store being managed by this
// KeyStoreManager
func (km *KeyStoreManager) CertificateStore() trustmanager.X509Store {
	return km.certificateStore
}

// CAStore returns the CA store being managed by this KeyStoreManager
func (km *KeyStoreManager) CAStore() trustmanager.X509Store {
	return km.caStore
}

/*
ValidateRoot iterates over every root key included in the TUF data and
attempts to validate the certificate by first checking for an exact match on
the certificate store, and subsequently trying to find a valid chain on the
caStore.

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
	rootSigned := &data.Root{}
	err := json.Unmarshal(root.Signed, rootSigned)
	if err != nil {
		return err
	}

	certs := make(map[string]*data.PublicKey)
	for _, keyID := range rootSigned.Roles["root"].KeyIDs {
		// TODO(dlaw): currently assuming only one cert contained in
		// public key entry. Need to fix when we want to pass in chains.
		k, _ := pem.Decode([]byte(rootSigned.Keys[keyID].Public()))
		decodedCerts, err := x509.ParseCertificates(k.Bytes)
		if err != nil {
			logrus.Debugf("error while parsing root certificate with keyID: %s, %v", keyID, err)
			continue
		}
		// TODO(diogo): Assuming that first certificate is the leaf-cert. Need to
		// iterate over all decodedCerts and find a non-CA one (should be the last).
		leafCert := decodedCerts[0]

		leafID, err := trustmanager.FingerprintCert(leafCert)
		if err != nil {
			logrus.Debugf("error while fingerprinting root certificate with keyID: %s, %v", keyID, err)
			continue
		}

		// Check to see if there is an exact match of this certificate.
		// Checking the CommonName is not required since ID is calculated over
		// Cert.Raw. It's included to prevent breaking logic with changes of how the
		// ID gets computed.
		_, err = km.certificateStore.GetCertificateByKeyID(leafID)
		if err == nil && leafCert.Subject.CommonName == dnsName {
			certs[keyID] = rootSigned.Keys[keyID]
		}

		// Check to see if this leafCertificate has a chain to one of the Root CAs
		// of our CA Store.
		certList := []*x509.Certificate{leafCert}
		err = trustmanager.Verify(km.caStore, dnsName, certList)
		if err == nil {
			certs[keyID] = rootSigned.Keys[keyID]
		}
	}

	if len(certs) < 1 {
		return errors.New("could not validate the path to a trusted root")
	}

	_, err = signed.VerifyRoot(root, 0, certs, 1)

	return err
}
