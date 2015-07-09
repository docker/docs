package trustmanager

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"time"

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

// CertToPEM is an utility function returns a PEM encoded x509 Certificate
func CertToPEM(cert *x509.Certificate) []byte {
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})

	return pemCert
}

// KeyToPEM returns a PEM encoded key from a crypto.PrivateKey
func KeyToPEM(key crypto.PrivateKey) ([]byte, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("only RSA keys are currently supported")
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(rsaKey)
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}), nil
}

// EncryptPrivateKey returns an encrypted PEM encoded key given a Private key
// and a passphrase
func EncryptPrivateKey(key crypto.PrivateKey, passphrase string) ([]byte, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("only RSA keys are currently supported")
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(rsaKey)

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

	key, err := ParsePEMPrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// ParsePEMPrivateKey returns a private key from a PEM encoded private key. It
// only supports RSA (PKCS#1).
func ParsePEMPrivateKey(pemBytes []byte) (crypto.PrivateKey, error) {
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

// TufParsePEMPrivateKey returns a data.PrivateKey from a PEM encoded private key. It
// only supports RSA (PKCS#1).
func TufParsePEMPrivateKey(pemBytes []byte) (*data.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no valid key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaPrivKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse PEM: %v", err)
		}

		tufRSAPrivateKey, err := RSAToPrivateKey(rsaPrivKey)
		if err != nil {
			return nil, fmt.Errorf("could not convert crypto.PrivateKey to PrivateKey: %v", err)
		}
		return tufRSAPrivateKey, nil
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}

func RSAToPrivateKey(rsaPrivKey *rsa.PrivateKey) (*data.PrivateKey, error) {
	// Get a DER-encoded representation of the PublicKey
	rsaPubBytes, err := x509.MarshalPKIXPublicKey(&rsaPrivKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}

	// Get a DER-encoded representation of the PrivateKey
	rsaPrivBytes := x509.MarshalPKCS1PrivateKey(rsaPrivKey)

	return data.NewPrivateKey("RSA", rsaPubBytes, rsaPrivBytes), nil
}

// ParsePEMEncryptedPrivateKey returns a private key from a PEM encrypted private key. It
// only supports RSA (PKCS#1).
func ParsePEMEncryptedPrivateKey(pemBytes []byte, passphrase string) (crypto.PrivateKey, error) {
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

func NewCertificate(gun, organization string) *x509.Certificate {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365 * 2)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	// TODO(diogo): Don't silently ignore this error
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   gun,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}
}
