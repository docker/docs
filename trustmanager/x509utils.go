package trustmanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
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
	cert, err := LoadCertFromPEM(certBytes)
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

// KeyToPEM returns a PEM encoded key from a Private Key
func KeyToPEM(privKey *data.PrivateKey) ([]byte, error) {
	if privKey.Cipher() != "RSA" {
		return nil, errors.New("only RSA keys are currently supported")
	}

	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKey.Private()}), nil
}

// EncryptPrivateKey returns an encrypted PEM key given a Privatekey
// and a passphrase
func EncryptPrivateKey(key *data.PrivateKey, passphrase string) ([]byte, error) {
	// TODO(diogo): Currently only supports RSA Private keys
	if key.Cipher() != "RSA" {
		return nil, errors.New("only RSA keys are currently supported")
	}

	password := []byte(passphrase)
	cipherType := x509.PEMCipherAES256
	blockType := "RSA PRIVATE KEY"

	encryptedPEMBlock, err := x509.EncryptPEMBlock(rand.Reader,
		blockType,
		key.Private(),
		password,
		cipherType)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(encryptedPEMBlock), nil
}

// LoadCertFromPEM returns the first certificate found in a bunch of bytes or error
// if nothing is found. Taken from https://golang.org/src/crypto/x509/cert_pool.go#L85.
func LoadCertFromPEM(pemBytes []byte) (*x509.Certificate, error) {
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
	certFiles := s.fileStore.ListFiles(true)
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

// ParsePEMPrivateKey returns a data.PrivateKey from a PEM encoded private key. It
// only supports RSA (PKCS#1) and attempts to decrypt using the passphrase, if encrypted.
func ParsePEMPrivateKey(pemBytes []byte, passphrase string) (*data.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no valid private key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		var privKeyBytes []byte
		var err error

		if x509.IsEncryptedPEMBlock(block) {
			privKeyBytes, err = x509.DecryptPEMBlock(block, []byte(passphrase))
			if err != nil {
				return nil, errors.New("could not decrypt private key")
			}
		} else {
			privKeyBytes = block.Bytes
		}

		rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse DER encoded key: %v", err)
		}

		tufRSAPrivateKey, err := RSAToPrivateKey(rsaPrivKey)
		if err != nil {
			return nil, fmt.Errorf("could not convert rsa.PrivateKey to data.PrivateKey: %v", err)
		}

		return tufRSAPrivateKey, nil
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}

// GenerateRSAKey generates an RSA Private key and returns a TUF PrivateKey
func GenerateRSAKey(random io.Reader, bits int) (*data.PrivateKey, error) {
	rsaPrivKey, err := rsa.GenerateKey(random, bits)
	if err != nil {
		return nil, fmt.Errorf("could not generate private key: %v", err)
	}

	return RSAToPrivateKey(rsaPrivKey)
}

// RSAToPrivateKey converts an rsa.Private key to a TUF data.PrivateKey type
func RSAToPrivateKey(rsaPrivKey *rsa.PrivateKey) (*data.PrivateKey, error) {
	// Get a DER-encoded representation of the PublicKey
	rsaPubBytes, err := x509.MarshalPKIXPublicKey(&rsaPrivKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %v", err)
	}

	// Get a DER-encoded representation of the PrivateKey
	rsaPrivBytes := x509.MarshalPKCS1PrivateKey(rsaPrivKey)

	return data.NewPrivateKey("RSA", rsaPubBytes, rsaPrivBytes), nil
}

// GenerateECDSAKey generates an ECDSA Private key and returns a TUF PrivateKey
func GenerateECDSAKey(random io.Reader) (*data.PrivateKey, error) {
	// TODO(diogo): For now hardcode P256. There were timming attacks on the other
	// curves, but I can't seem to find the issue.
	ecdsaPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), random)
	if err != nil {
		return nil, err
	}

	return ECDSAToPrivateKey(ecdsaPrivKey)
}

// ECDSAToPrivateKey converts an rsa.Private key to a TUF data.PrivateKey type
func ECDSAToPrivateKey(ecdsaPrivKey *ecdsa.PrivateKey) (*data.PrivateKey, error) {
	// Get a DER-encoded representation of the PublicKey
	ecdsaPubBytes, err := x509.MarshalPKIXPublicKey(&ecdsaPrivKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %v", err)
	}

	// Get a DER-encoded representation of the PrivateKey
	ecdsaPrivKeyBytes, err := x509.MarshalECPrivateKey(ecdsaPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}

	return data.NewPrivateKey("ECDSA", ecdsaPubBytes, ecdsaPrivKeyBytes), nil
}

// NewCertificate returns an X509 Certificate following a template, given a GUN.
func NewCertificate(gun string) (*x509.Certificate, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365 * 2)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new certificate: %v", err)
	}

	// TODO(diogo): Currently hard coding organization to be the gun. Revisit.
	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{gun},
			CommonName:   gun,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}, nil
}
