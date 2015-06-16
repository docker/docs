package trustmanager

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"testing"

	"github.com/docker/vetinari/trustmanager"
)

func TestAddCert(t *testing.T) {
	// Read certificate from file
	b, err := ioutil.ReadFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	// Decode PEM block
	var block *pem.Block
	block, _ = pem.Decode(b)

	// Load X509 Certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}
	// Create a Store and add the certificate to it
	store := trustmanager.NewX509MemStore()
	err = store.AddCert(cert)
	if err != nil {
		t.Fatalf("failed to load certificate: %v", err)
	}
	// Retrieve all the certificates
	certs := store.GetCertificates()
	// Check to see if certificate is present and total number of certs is correct
	numCerts := len(certs)
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}
	if certs[0] != cert {
		t.Fatalf("expected certificates to be the same")
	}
}

func TestAddCertFromFile(t *testing.T) {
	store := trustmanager.NewX509MemStore()
	err := store.AddCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}
	numCerts := len(store.GetCertificates())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}
}

func TestAddCertFromPEM(t *testing.T) {
	b, err := ioutil.ReadFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	store := trustmanager.NewX509MemStore()
	err = store.AddCertFromPEM(b)
	if err != nil {
		t.Fatalf("failed to load certificate from PEM: %v", err)
	}
	numCerts := len(store.GetCertificates())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}
}

func TestRemoveCert(t *testing.T) {
	b, err := ioutil.ReadFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	var block *pem.Block
	block, _ = pem.Decode(b)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}

	store := trustmanager.NewX509MemStore()
	err = store.AddCert(cert)
	if err != nil {
		t.Fatalf("failed to load certificate: %v", err)
	}

	// Number of certificates should be 1 since we added the cert
	numCerts := len(store.GetCertificates())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}

	// Remove the cert from the store
	err = store.RemoveCert(cert)
	if err != nil {
		t.Fatalf("failed to remove certificate: %v", err)
	}
	// Number of certificates should be 0 since we added and removed the cert
	numCerts = len(store.GetCertificates())
	if numCerts != 0 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}
}

func TestInexistentGetCertificateBySKID(t *testing.T) {
	store := trustmanager.NewX509MemStore()
	err := store.AddCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	_, err = store.GetCertificateBySKID("4d06afd30b8bed131d2a84c97d00b37f422021598bfae34285ce98e77b708b5a")
	if err == nil {
		t.Fatalf("no error returned for inexistent certificate")
	}
}

func TestGetCertificateBySKID(t *testing.T) {
	b, err := ioutil.ReadFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	var block *pem.Block
	block, _ = pem.Decode(b)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}

	store := trustmanager.NewX509MemStore()
	err = store.AddCert(cert)
	if err != nil {
		t.Fatalf("failed to load certificate from PEM: %v", err)
	}

	// Calculate SHA256 fingerprint for cert
	fingerprintBytes := sha256.Sum256(cert.Raw)
	certFingerprint := hex.EncodeToString(fingerprintBytes[:])

	// Tries to retreive cert by Subject Key IDs
	_, err = store.GetCertificateBySKID(certFingerprint)
	if err != nil {
		t.Fatalf("expected certificate in store: %s", certFingerprint)
	}
}

func TestGetVerifyOpsErrorsWithoutCerts(t *testing.T) {
	// Create empty Store
	store := trustmanager.NewX509MemStore()

	// Try to get VerifyOptions without certs added
	_, err := store.GetVerifyOptions("docker.com")
	if err == nil {
		t.Fatalf("expecting an error when getting empty VerifyOptions")
	}
}

func TestVerifyLeafCertFromIntermediate(t *testing.T) {
	// Create a store and add a root
	store := trustmanager.NewX509MemStore()
	err := store.AddCertFromFile("../fixtures/notary/ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("secure.docker.com")

	// Get leaf certificate
	b, err := ioutil.ReadFile("../fixtures/notary/secure.docker.com.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	var block *pem.Block
	block, _ = pem.Decode(b)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}

	// Try to find a valid chain for cert
	_, err = cert.Verify(opts)
	if err != nil {
		t.Fatalf("couldn't find a valid chain for this certificate: %v", err)
	}
}

func TestVerifyIntermediateFromRoot(t *testing.T) {
	// Create a store and add a root
	store := trustmanager.NewX509MemStore()
	err := store.AddCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("Docker CA")

	// Get leaf certificate
	b, err := ioutil.ReadFile("../fixtures/notary/ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	var block *pem.Block
	block, _ = pem.Decode(b)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}

	// Try to find a valid chain for cert
	_, err = cert.Verify(opts)
	if err != nil {
		t.Fatalf("couldn't find a valid chain for this certificate: %v", err)
	}
}

func TestNewX509FilteredMemStore(t *testing.T) {
	store := trustmanager.NewX509FilteredMemStore(func(cert *x509.Certificate) bool {
		return cert.IsCA
	})

	// AddCert should succeed because this is a CA being added
	err := store.AddCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}
	numCerts := len(store.GetCertificates())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}

	// AddCert should fail because this is a leaf cert being added
	err = store.AddCertFromFile("../fixtures/notary/secure.docker.com.crt")
	if err == nil {
		t.Fatalf("was expecting non-CA certificate to be rejected")
	}
}

func TestGetCertificatePool(t *testing.T) {
	// Create a store and add a root
	store := trustmanager.NewX509MemStore()
	err := store.AddCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	pool := store.GetCertificatePool()
	numCerts := len(pool.Subjects())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in pool: %d", numCerts)
	}
}
