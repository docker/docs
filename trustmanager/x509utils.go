package trustmanager

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/endophage/gotuf/data"
)

// GetCertFromURL tries to get a X509 certificate given a HTTPS URL
func GetCertFromURL(urlStr string) (*x509.Certificate, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Check if we are adding via HTTPS
	if url.Scheme != "https" {
		return nil, errors.New("only HTTPS URLs allowed")
	}

	// Download the certificate and write to directory
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	// Copy the content to certBytes
	defer resp.Body.Close()
	certBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Try to extract the first valid PEM certificate from the bytes
	cert, err := loadCertFromPEM(certBytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// ToPEM is an utility function returns a PEM encoded x509 Certificate
func ToPEM(cert *x509.Certificate) []byte {
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})

	return pemCert
}

// KeyToPEM is an utility function returns a PEM encoded Key
func KeyToPEM(key crypto.PrivateKey) ([]byte, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("only RSA keys are currently supported")
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(rsaKey)
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}), nil
}

// KeyToEncryptedPEM is an utility function returns a PEM encoded Key
func KeyToEncryptedPEM(key crypto.PrivateKey, passphrase string) ([]byte, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("only RSA keys are currently supported")
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(rsaKey)

	//TODO(diogo): if we do keystretching, where do we keep the salt + params?
	password := []byte(passphrase)
	cipherType := x509.PEMCipherAES256
	blockType := "RSA PRIVATE KEY"

	encryptedPEMBlock, err := x509.EncryptPEMBlock(rand.Reader,
		blockType,
		keyBytes,
		password,
		cipherType)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(encryptedPEMBlock), nil
}

// loadCertFromPEM returns the first certificate found in a bunch of bytes or error
// if nothing is found. Taken from https://golang.org/src/crypto/x509/cert_pool.go#L85.
func loadCertFromPEM(pemBytes []byte) (*x509.Certificate, error) {
	for len(pemBytes) > 0 {
		var block *pem.Block
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			return nil, errors.New("no certificates found in PEM data")
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}

		return cert, nil
	}

	return nil, errors.New("no certificates found in PEM data")
}

// FingerprintCert returns a TUF compliant fingerprint for a X509 Certificate
func FingerprintCert(cert *x509.Certificate) string {
	return string(fingerprintCert(cert))
}

func fingerprintCert(cert *x509.Certificate) CertID {
	block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	pemdata := pem.EncodeToMemory(&block)

	// Create new TUF Key so we can compute the TUF-compliant CertID
	tufKey := data.NewTUFKey("RSA", pemdata, nil)

	return CertID(tufKey.ID())
}

// loadCertsFromDir receives a store AddCertFromFile for each certificate found
func loadCertsFromDir(s *X509FileStore) {
	certFiles := s.fileStore.ListAll()
	for _, c := range certFiles {
		s.AddCertFromFile(c)
	}
}

// LoadCertFromFile tries to adds a X509 certificate to the store given a filename
func LoadCertFromFile(filename string) (*x509.Certificate, error) {
	// TODO(diogo): handle multiple certificates in one file.
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var block *pem.Block
	block, b = pem.Decode(b)
	for ; block != nil; block, b = pem.Decode(b) {
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				return cert, nil
			}
		}
	}

	return nil, errors.New("could not load certificate from file")
}

// LoadKeyFromFile returns a PrivateKey given a filename
func LoadKeyFromFile(filename string) (crypto.PrivateKey, error) {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	key, err := ParseRawPrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// ParseRawPrivateKey returns a private key from a PEM encoded private key. It
// only supports RSA (PKCS#1).
func ParseRawPrivateKey(pemBytes []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no valid key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}

// ParseRawEncryptedPrivateKey returns a private key from a PEM encrypted private key. It
// only supports RSA (PKCS#1).
func ParseRawEncryptedPrivateKey(pemBytes []byte, passphrase string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no valid private key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		if !x509.IsEncryptedPEMBlock(block) {
			return nil, errors.New("private key is not encrypted")
		}

		decryptedPEMBlock, err := x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, errors.New("could not decrypt private key")
		}

		return x509.ParsePKCS1PrivateKey(decryptedPEMBlock)
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}
