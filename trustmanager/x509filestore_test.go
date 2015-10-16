package trustmanager

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	store, err := NewX509FileStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create a new X509FileStore: %v", store)
	}
}

// NewX509FileStore loads any existing certs from the directory, and does
// not overwrite any of the.
func TestNewX509FileStoreLoadsExistingCerts(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	assert.NoError(t, err, "couldn't open temp directory")
	defer os.RemoveAll(tempDir)

	certBytes, err := ioutil.ReadFile("../fixtures/root-ca.crt")
	assert.NoError(t, err, "couldn't read fixtures/root-ca.crt")
	out, err := os.Create(filepath.Join(tempDir, "root-ca.crt"))
	assert.NoError(t, err, "couldn't create a file in the temp dir")

	// to distinguish it from the canonical format
	distinguishingBytes := []byte{'\n', '\n', '\n', '\n', '\n', '\n'}
	nBytes, err := out.Write(distinguishingBytes)
	assert.NoError(t, err, "could not write newlines to the temporary file")
	assert.Equal(t, len(distinguishingBytes), nBytes,
		"didn't write all bytes to temporary file")

	nBytes, err = out.Write(certBytes)
	assert.NoError(t, err, "could not write cert to the temporary file")
	assert.Equal(t, len(certBytes), nBytes,
		"didn't write all bytes to temporary file")

	err = out.Close()
	assert.NoError(t, err, "could not close temporary file")

	store, err := NewX509FileStore(tempDir)
	assert.NoError(t, err, "failed to create a new X509FileStore")

	expectedCert, err := LoadCertFromFile("../fixtures/root-ca.crt")
	assert.NoError(t, err, "could not load root-ca.crt")
	assert.Equal(t, store.GetCertificates(), []*x509.Certificate{expectedCert},
		"did not load certificate already in the directory")

	outBytes, err := ioutil.ReadFile(filepath.Join(tempDir, "root-ca.crt"))
	assert.NoError(t, err, "couldn't read temporary file")
	assert.Equal(t, distinguishingBytes, outBytes[:6], "original file overwritten")
	assert.Equal(t, certBytes, outBytes[6:], "original file overwritten")
}

func TestAddCertX509FileStore(t *testing.T) {
	// Read certificate from file
	b, err := ioutil.ReadFile("../fixtures/root-ca.crt")
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
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	// Create a Store and add the certificate to it
	store, _ := NewX509FileStore(tempDir)
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

func TestAddCertFromFileX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	assert.NoError(t, err, "failed to create temporary directory")

	store, err := NewX509FileStore(tempDir)
	assert.NoError(t, err, "failed to load x509 filestore")

	err = store.AddCertFromFile("../fixtures/root-ca.crt")
	assert.NoError(t, err, "failed to add certificate from file")
	assert.Len(t, store.GetCertificates(), 1)

	// Now load the x509 filestore with the same path and expect the same result
	newStore, err := NewX509FileStore(tempDir)
	assert.NoError(t, err, "failed to load x509 filestore")
	assert.Len(t, newStore.GetCertificates(), 1)

	// Test that adding the same certificate returns an error
	err = newStore.AddCert(newStore.GetCertificates()[0])
	if assert.Error(t, err, "expected error when adding certificate twice") {
		assert.Equal(t, err, &ErrCertExists{})
	}
}

// TestNewX509FileStoreEmpty verifies the behavior of the Empty function
func TestNewX509FileStoreEmpty(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	store, err := NewX509FileStore(tempDir)
	assert.NoError(t, err, "failed to create a new X509FileStore: %v", store)
	assert.True(t, store.Empty())

	err = store.AddCertFromFile("../fixtures/root-ca.crt")
	assert.NoError(t, err, "failed to add certificate from file")
	assert.False(t, store.Empty())
}

func TestAddCertFromPEMX509FileStore(t *testing.T) {
	b, err := ioutil.ReadFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	store, _ := NewX509FileStore(tempDir)
	err = store.AddCertFromPEM(b)
	if err != nil {
		t.Fatalf("failed to load certificate from PEM: %v", err)
	}
	numCerts := len(store.GetCertificates())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}
}

func TestRemoveCertX509FileStore(t *testing.T) {
	b, err := ioutil.ReadFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	var block *pem.Block
	block, _ = pem.Decode(b)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}

	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	store, _ := NewX509FileStore(tempDir)
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

func TestRemoveAllX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}

	// Add three certificates to store
	store, _ := NewX509FileStore(tempDir)
	certFiles := [3]string{"../fixtures/root-ca.crt",
		"../fixtures/intermediate-ca.crt",
		"../fixtures/secure.example.com.crt"}
	for _, file := range certFiles {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("couldn't load fixture: %v", err)
		}
		var block *pem.Block
		block, _ = pem.Decode(b)

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			t.Fatalf("couldn't parse certificate: %v", err)
		}
		err = store.AddCert(cert)
		if err != nil {
			t.Fatalf("failed to load certificate: %v", err)
		}
	}

	// Number of certificates should be 3 since we added the cert
	numCerts := len(store.GetCertificates())
	if numCerts != 3 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}

	// Remove the cert from the store
	err = store.RemoveAll()
	if err != nil {
		t.Fatalf("failed to remove all certificates: %v", err)
	}
	// Number of certificates should be 0 since we added and removed the cert
	numCerts = len(store.GetCertificates())
	if numCerts != 0 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}
}
func TestInexistentGetCertificateByKeyIDX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	store, _ := NewX509FileStore(tempDir)
	err = store.AddCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	_, err = store.GetCertificateByCertID("4d06afd30b8bed131d2a84c97d00b37f422021598bfae34285ce98e77b708b5a")
	if err == nil {
		t.Fatalf("no error returned for inexistent certificate")
	}
}

func TestGetCertificateByKeyIDX509FileStore(t *testing.T) {
	b, err := ioutil.ReadFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}
	var block *pem.Block
	block, _ = pem.Decode(b)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("couldn't parse certificate: %v", err)
	}

	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	store, _ := NewX509FileStore(tempDir)
	err = store.AddCert(cert)
	if err != nil {
		t.Fatalf("failed to load certificate from PEM: %v", err)
	}

	keyID, err := FingerprintCert(cert)
	if err != nil {
		t.Fatalf("failed to fingerprint the certificate: %v", err)
	}

	// Tries to retrieve cert by Subject Key IDs
	_, err = store.GetCertificateByCertID(keyID)
	if err != nil {
		t.Fatalf("expected certificate in store: %s", keyID)
	}
}

func TestGetVerifyOpsErrorsWithoutCertsX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	// Create empty Store
	store, _ := NewX509FileStore(tempDir)

	// Try to get VerifyOptions without certs added
	_, err = store.GetVerifyOptions("example.com")
	if err == nil {
		t.Fatalf("expecting an error when getting empty VerifyOptions")
	}
}

func TestVerifyLeafCertFromIntermediateX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	// Create a store and add a root
	store, _ := NewX509FileStore(tempDir)
	err = store.AddCertFromFile("../fixtures/intermediate-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("secure.example.com")

	// Get leaf certificate
	b, err := ioutil.ReadFile("../fixtures/secure.example.com.crt")
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

func TestVerifyIntermediateFromRootX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	// Create a store and add a root
	store, _ := NewX509FileStore(tempDir)
	err = store.AddCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("Notary Testing CA")

	// Get leaf certificate
	b, err := ioutil.ReadFile("../fixtures/intermediate-ca.crt")
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

func TestNewX509FilteredFileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewX509FilteredFileStore(tempDir, func(cert *x509.Certificate) bool {
		return cert.IsCA
	})
	if err != nil {
		t.Fatalf("failed to create new X509FilteredFileStore: %v", err)
	}
	// AddCert should succeed because this is a CA being added
	err = store.AddCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}
	numCerts := len(store.GetCertificates())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in store: %d", numCerts)
	}

	// AddCert should fail because this is a leaf cert being added
	err = store.AddCertFromFile("../fixtures/secure.example.com.crt")
	if err == nil {
		t.Fatalf("was expecting non-CA certificate to be rejected")
	}
}

func TestGetCertificatePoolX509FileStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	// Create a store and add a root
	store, _ := NewX509FileStore(tempDir)
	err = store.AddCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	pool := store.GetCertificatePool()
	numCerts := len(pool.Subjects())
	if numCerts != 1 {
		t.Fatalf("unexpected number of certificates in pool: %d", numCerts)
	}
}
