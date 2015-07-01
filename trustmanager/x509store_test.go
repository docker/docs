package trustmanager

import (
	"crypto/x509"
	"fmt"
	"testing"
)

func TestVerifyLeafSuccessfully(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get intermediate certificate
	intermediateCA, err := LoadCertFromFile("../fixtures/notary/ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get leaf certificate
	leafCert, err := LoadCertFromFile("../fixtures/notary/secure.docker.com.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Create a store and add the CA root
	store := NewX509MemStore()
	err = store.AddCert(rootCA)
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get our certList with Leaf Cert and Intermediate
	certList := []*x509.Certificate{leafCert, intermediateCA}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("secure.docker.com")
	fmt.Println(opts)

	// Try to find a valid chain for cert
	err = Verify(store, "secure.docker.com", certList)
	if err != nil {
		t.Fatalf("expected to find a valid chain for this certificate: %v", err)
	}
}

func TestVerifyLeafSuccessfullyWithMultipleIntermediates(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get intermediate certificate
	intermediateCA, err := LoadCertFromFile("../fixtures/notary/ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get leaf certificate
	leafCert, err := LoadCertFromFile("../fixtures/notary/secure.docker.com.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Create a store and add the CA root
	store := NewX509MemStore()
	err = store.AddCert(rootCA)
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get our certList with Leaf Cert and Intermediate
	certList := []*x509.Certificate{leafCert, intermediateCA, intermediateCA, rootCA}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("secure.docker.com")
	fmt.Println(opts)

	// Try to find a valid chain for cert
	err = Verify(store, "secure.docker.com", certList)
	if err != nil {
		t.Fatalf("expected to find a valid chain for this certificate: %v", err)
	}
}

func TestVerifyLeafWithNoIntermediate(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get leaf certificate
	leafCert, err := LoadCertFromFile("../fixtures/notary/secure.docker.com.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Create a store and add the CA root
	store := NewX509MemStore()
	err = store.AddCert(rootCA)
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get our certList with Leaf Cert and Intermediate
	certList := []*x509.Certificate{leafCert, leafCert}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("secure.docker.com")
	fmt.Println(opts)

	// Try to find a valid chain for cert
	err = Verify(store, "secure.docker.com", certList)
	if err == nil {
		t.Fatalf("expected error due to more than one leaf certificate")
	}
}

func TestVerifyLeafWithNoLeaf(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/notary/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get intermediate certificate
	intermediateCA, err := LoadCertFromFile("../fixtures/notary/ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Create a store and add the CA root
	store := NewX509MemStore()
	err = store.AddCert(rootCA)
	if err != nil {
		t.Fatalf("failed to load certificate from file: %v", err)
	}

	// Get our certList with Leaf Cert and Intermediate
	certList := []*x509.Certificate{intermediateCA, intermediateCA}

	// Get the VerifyOptions from our Store
	opts, err := store.GetVerifyOptions("secure.docker.com")
	fmt.Println(opts)

	// Try to find a valid chain for cert
	err = Verify(store, "secure.docker.com", certList)
	if err == nil {
		t.Fatalf("expected error due to no leafs provided")
	}
}
