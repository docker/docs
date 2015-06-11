package signed

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/agl/ed25519"
	"github.com/endophage/gotuf/data"
)

// Verifiers serves as a map of all verifiers available on the system and
// can be injected into a verificationService. For testing and configuration
// purposes, it will not be used by default.
var Verifiers = map[string]Verifier{
	"ed25519": Ed25519Verifier{},
	"rsa":     RSAVerifier{},
	"pycrypto-pkcs#1 pss": RSAPSSVerifier{},
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

type RSAVerifier struct{}

func (v RSAVerifier) Verify(key data.Key, sig []byte, msg []byte) error {
	digest := sha256.Sum256(msg)
	keyReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(key.Public()))
	keyBytes, _ := ioutil.ReadAll(keyReader)
	pub, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		logrus.Infof("Failed to parse public key: %s\n", err)
		return ErrInvalid
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		logrus.Infof("Value returned from ParsePKIXPublicKey was not an RSA public key")
		return ErrInvalid
	}

	if err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest[:], sig); err != nil {
		logrus.Infof("Failed verification: %s", err)
		return ErrInvalid
	}
	return nil
}

// RSAPSSVerifier checks RSASSA-PSS signatures
type RSAPSSVerifier struct{}

// Verify does the actual check.
// N.B. We have not been able to make this work in a way that is compatible
// with PyCrypto.
func (v RSAPSSVerifier) Verify(key data.Key, sig []byte, msg []byte) error {
	digest := sha256.Sum256(msg)

	k, _ := pem.Decode([]byte(key.Public()))
	pub, err := x509.ParsePKIXPublicKey(k.Bytes)
	if err != nil {
		logrus.Infof("Failed to parse public key: %s\n", err)
		return ErrInvalid
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		logrus.Infof("Value returned from ParsePKIXPublicKey was not an RSA public key")
		return ErrInvalid
	}

	opts := rsa.PSSOptions{SaltLength: sha256.Size, Hash: crypto.SHA256}
	if err = rsa.VerifyPSS(rsaPub, crypto.SHA256, digest[:], sig, &opts); err != nil {
		logrus.Infof("Failed verification: %s", err)
		return ErrInvalid
	}
	return nil
}
