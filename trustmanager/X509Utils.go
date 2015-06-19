package trustmanager

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

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
		return nil, errors.New("only HTTPS URLs allowed.")
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

// saveCertificate is an utility function that saves a certificate as a PEM
// encoded block to a file.
func saveCertificate(cert *x509.Certificate, filename string) error {
	block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	pemdata := string(pem.EncodeToMemory(&block))

	err := CreateDirectory(filename)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, []byte(pemdata), 0600)
	if err != nil {
		return err
	}
	return nil
}

func FingerprintCert(cert *x509.Certificate) ID {
	block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	pemdata := string(pem.EncodeToMemory(&block))

	// Create new TUF Key so we can compute the TUF-compliant ID
	tufKey := data.NewTUFKey("RSA", pemdata, "")

	return ID(tufKey.ID())
}

// loadCertsFromDir receives a store and a directory and calls AddCertFromFile
// for each certificate found
func loadCertsFromDir(s *X509FileStore, directory string) {
	filepath.Walk(directory, func(fp string, fi os.FileInfo, err error) error {
		// If there are errors, ignore this particular file
		if err != nil {
			return nil
		}
		// Ignore if it is a directory
		if !!fi.IsDir() {
			return nil
		}
		// Only allow matches that end with our certificate extension (e.g. *.crt)
		matched, _ := filepath.Match("*"+certExtension, fi.Name())

		if matched {
			s.AddCertFromFile(fp)
		}
		return nil
	})
}

// LoadCertFromFile tries to adds a X509 certificate to the store given a filename
func LoadCertFromFile(filename string) (*x509.Certificate, error) {
	// TODO(diogo): handle multiple certificates in one file. Demultiplex into
	// multiple files or load only first
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

func CreateDirectory(dir string) error {
	cleanDir := filepath.Dir(dir)
	if err := os.MkdirAll(cleanDir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}
	return nil
}
