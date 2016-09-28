package jose

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Base64URLEncode encodes the given data using the standard base64 URL
// encoding format but with all trailing '=' characters ommitted in accordance
// with the jose specification.
// http://tools.ietf.org/html/rfc7515#section-2
func Base64URLEncode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

// Base64URLDecode decodes the given string using the standard base64 URL
// decoder but first adds the appropriate number of trailing '=' characters in
// accordance with the jose specification.
// http://tools.ietf.org/html/rfc7515#section-2
func Base64URLDecode(s string) ([]byte, error) {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)

	switch len(s) % 4 {
	case 0:
	case 2:
		s += "=="
	case 3:
		s += "="
	default:
		return nil, errors.New("illegal base64url string")
	}

	return base64.URLEncoding.DecodeString(s)
}

// makeKeyID computes the key ID of the given public key as a simple
// hex-encoded SHA-256 hash of the public key when serialized to ASN.1.
func makeKeyID(key crypto.PublicKey) (string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to encode public key: %s", err)
	}

	hash := crypto.SHA256.New()
	hash.Write(derBytes)

	return hex.EncodeToString(hash.Sum(nil)), nil
}
