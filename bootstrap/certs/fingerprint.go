package certs

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	ErrMalformedCert = errors.New("Malformed certificate")
)

func getFingerprint(der []byte, algo string) string {
	var hash []byte
	algo = strings.ToUpper(algo)
	if algo == "SHA256" || algo == "SHA-256" {
		v := sha256.Sum256(der)
		hash = v[:]
	} else {
		algo = "SHA1"
		v := sha1.Sum(der)
		hash = v[:]
	}
	hexified := make([][]byte, len(hash))
	for i, data := range hash {
		hexified[i] = []byte(fmt.Sprintf("%02X", data))
	}
	return fmt.Sprintf("%s Fingerprint=%s", algo, string(bytes.Join(hexified, []byte(":"))))
}

func GetFingerprint(pemCertFilename string) (string, error) {
	pemCert, err := ioutil.ReadFile(pemCertFilename)
	if err != nil {
		return "", err
	}

	der, _ := pem.Decode([]byte(pemCert))
	if der == nil {
		return "", ErrMalformedCert
	}

	return getFingerprint(der.Bytes, "SHA-256"), nil
}
