package trustmanager

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
)

// saveCertificate is an utility function that saves a certificate as a PEM
// encoded block to a file.
func saveCertificate(cert *x509.Certificate, filename string) error {
	block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	pemdata := string(pem.EncodeToMemory(&block))

	err := ioutil.WriteFile(filename, []byte(pemdata), 0600)
	if err != nil {
		return err
	}
	return nil
}

func fingerprintCert(cert *x509.Certificate) ID {
	fingerprintBytes := sha256.Sum256(cert.Raw)
	return ID(hex.EncodeToString(fingerprintBytes[:]))
}

// loadCertsFromDir receives a store and a directory and calls loadCertFromFile
// for each certificate found
func loadCertsFromDir(s *X509FileStore, directory string) {
	certFiles, _ := filepath.Glob(path.Join(directory, fmt.Sprintf("*%s", certExtension)))
	for _, f := range certFiles {
		cert, err := loadCertFromFile(f)
		// Ignores files that do not contain valid certificates
		if err == nil {
			s.addNamedCert(cert, f)
		}
	}
}

// loadCertFromFile tries to adds a X509 certificate to the store given a filename
func loadCertFromFile(filename string) (*x509.Certificate, error) {
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
// if nothing is found
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
