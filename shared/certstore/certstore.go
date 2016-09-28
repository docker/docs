package certstore

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CertStore struct {
	Path string
}

func New(path string) *CertStore {
	return &CertStore{path}
}

func (cs CertStore) Mkdirp() error {
	return os.MkdirAll(cs.Path, 0700)
}

func (cs CertStore) LoadCACertPool() (*x509.CertPool, error) {
	content, err := ioutil.ReadFile(cs.CertPath())
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(content)
	return caCertPool, nil
}

func (cs CertStore) LoadKeyPair() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(cs.CertPath(), cs.KeyPath())
}

func (cs CertStore) LoadCert() (*x509.Certificate, error) {
	content, err := ioutil.ReadFile(cs.CertPath())
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (cs CertStore) LoadKey() (*ecdsa.PrivateKey, error) {
	content, err := ioutil.ReadFile(cs.KeyPath())
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	pk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (cs CertStore) Exists() (bool, error) {
	exists1, err := fileExists(cs.CertPath())
	if err != nil {
		return false, fmt.Errorf("Cannot determine if cert exists: %s", err)
	}
	exists2, err := fileExists(cs.KeyPath())
	if err != nil {
		return false, fmt.Errorf("Cannot determine if key exists: %s", err)
	}
	if exists1 && exists2 {
		return true, nil
	}
	return false, nil
}

func (cs CertStore) SaveCertBytes(derBytes []byte) error {
	certOut, err := os.Create(cs.CertPath())
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %s", err)
	}

	// write out public key
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %s", err)
	}
	certOut.Close()
	return nil
}

func (cs CertStore) SaveKey(priv *ecdsa.PrivateKey) error {
	// write out private key
	keyOut, err := os.OpenFile(cs.KeyPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing:", err)
	}
	bytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("failed to marshal private key:", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: bytes})
	keyOut.Close()
	return nil
}

func (cs CertStore) CertPath() string {
	return filepath.Join(cs.Path, "cert.pem")
}
func (cs CertStore) KeyPath() string {
	return filepath.Join(cs.Path, "key.pem")
}
