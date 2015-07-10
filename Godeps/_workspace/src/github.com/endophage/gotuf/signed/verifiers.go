package signed

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/agl/ed25519"
	"github.com/endophage/gotuf/data"
)

// Verifiers serves as a map of all verifiers available on the system and
// can be injected into a verificationService. For testing and configuration
// purposes, it will not be used by default.
var Verifiers = map[string]Verifier{
	"ed25519":             Ed25519Verifier{},
	"rsassa-pss":          RSAPSSVerifier{},
	"rsassa-pss-x509":     RSAPSSX509Verifier{},
	"pycrypto-pkcs#1 pss": RSAPyCryptoVerifier{},
}

// RegisterVerifier provides a convenience function for init() functions
// to register additional verifiers or replace existing ones.
func RegisterVerifier(name string, v Verifier) {
	curr, ok := Verifiers[name]
	if ok {
		typOld := reflect.TypeOf(curr)
		typNew := reflect.TypeOf(v)
		logrus.Debugf(
			"Replacing already loaded verifier %s:%s with %s:%s",
			typOld.PkgPath(), typOld.Name(),
			typNew.PkgPath(), typNew.Name(),
		)
	} else {
		logrus.Debug("Adding verifier for: ", name)
	}
	Verifiers[name] = v
}

type Ed25519Verifier struct{}

func (v Ed25519Verifier) Verify(key data.Key, sig []byte, msg []byte) error {
	var sigBytes [ed25519.SignatureSize]byte
	if len(sig) != len(sigBytes) {
		logrus.Infof("Signature length is incorrect, must be %d, was %d.", ed25519.SignatureSize, len(sig))
		return ErrInvalid
	}
	copy(sigBytes[:], sig)

	var keyBytes [ed25519.PublicKeySize]byte
	copy(keyBytes[:], key.Public())

	if !ed25519.Verify(&keyBytes, msg, &sigBytes) {
		logrus.Infof("Failed ed25519 verification")
		return ErrInvalid
	}
	return nil
}

func verifyPSS(key interface{}, digest, sig []byte) error {
	rsaPub, ok := key.(*rsa.PublicKey)
	if !ok {
		logrus.Infof("Value was not an RSA public key")
		return ErrInvalid
	}

	opts := rsa.PSSOptions{SaltLength: sha256.Size, Hash: crypto.SHA256}
	if err := rsa.VerifyPSS(rsaPub, crypto.SHA256, digest[:], sig, &opts); err != nil {
		logrus.Infof("Failed verification: %s", err)
		return ErrInvalid
	}
	return nil
}

// RSAPSSVerifier checks RSASSA-PSS signatures
type RSAPSSVerifier struct{}

// Verify does the actual check.
func (v RSAPSSVerifier) Verify(key data.Key, sig []byte, msg []byte) error {
	digest := sha256.Sum256(msg)

	pub, err := x509.ParsePKIXPublicKey(key.Public())
	if err != nil {
		logrus.Infof("Failed to parse public key: %s\n", err)
		return ErrInvalid
	}

	return verifyPSS(pub, digest[:], sig)
}

// RSAPSSVerifier checks RSASSA-PSS signatures
type RSAPyCryptoVerifier struct{}

// Verify does the actual check.
// N.B. We have not been able to make this work in a way that is compatible
// with PyCrypto.
func (v RSAPyCryptoVerifier) Verify(key data.Key, sig []byte, msg []byte) error {
	digest := sha256.Sum256(msg)

	k, _ := pem.Decode([]byte(key.Public()))
	if k == nil {
		logrus.Infof("Failed to decode PEM-encoded x509 certificate")
		return ErrInvalid
	}

	pub, err := x509.ParsePKIXPublicKey(k.Bytes)
	if err != nil {
		logrus.Infof("Failed to parse public key: %s\n", err)
		return ErrInvalid
	}

	return verifyPSS(pub, digest[:], sig)
}

// RSAPSSPEMVerifier checks RSASSA-PSS signatures, extracting the public key
// from an X509 certificate.
type RSAPSSX509Verifier struct{}

// Verify does the actual check.
func (v RSAPSSX509Verifier) Verify(key data.Key, sig []byte, msg []byte) error {
	digest := sha256.Sum256(msg)

	k, _ := pem.Decode([]byte(key.Public()))
	if k == nil {
		logrus.Infof("Failed to decode PEM-encoded x509 certificate")
		return ErrInvalid
	}
	cert, err := x509.ParseCertificate(k.Bytes)
	if err != nil {
		logrus.Infof("Failed to parse x509 certificate: %s\n", err)
		return ErrInvalid
	}

	return verifyPSS(cert.PublicKey, digest[:], sig)
}
