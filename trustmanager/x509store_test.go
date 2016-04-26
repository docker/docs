package trustmanager

import (
	"crypto/x509"
	"testing"
)

func TestVerifyLeafSuccessfully(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get intermediate certificate
	intermediateCA, err := LoadCertFromFile("../fixtures/intermediate-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get leaf certificate
	leafCert, err := LoadCertFromFile("../fixtures/secure.example.com.crt")
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

	// Try to find a valid chain for cert
	err = Verify(store, "secure.example.com", certList)
	if err != nil {
		t.Fatalf("expected to find a valid chain for this certificate: %v", err)
	}
}

func TestVerifyLeafSuccessfullyWithMultipleIntermediates(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get intermediate certificate
	intermediateCA, err := LoadCertFromFile("../fixtures/intermediate-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get leaf certificate
	leafCert, err := LoadCertFromFile("../fixtures/secure.example.com.crt")
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

	// Try to find a valid chain for cert
	err = Verify(store, "secure.example.com", certList)
	if err != nil {
		t.Fatalf("expected to find a valid chain for this certificate: %v", err)
	}
}

func TestVerifyLeafWithNoIntermediate(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get leaf certificate
	leafCert, err := LoadCertFromFile("../fixtures/secure.example.com.crt")
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

	// Try to find a valid chain for cert
	err = Verify(store, "secure.example.com", certList)
	if err == nil {
		t.Fatalf("expected error due to more than one leaf certificate")
	}
}

func TestVerifyLeafWithNoLeaf(t *testing.T) {
	// Get root certificate
	rootCA, err := LoadCertFromFile("../fixtures/root-ca.crt")
	if err != nil {
		t.Fatalf("couldn't load fixture: %v", err)
	}

	// Get intermediate certificate
	intermediateCA, err := LoadCertFromFile("../fixtures/intermediate-ca.crt")
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

	// Try to find a valid chain for cert
	err = Verify(store, "secure.example.com", certList)
	if err == nil {
		t.Fatalf("expected error due to no leaves provided")
	}
}
