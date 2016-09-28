package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/pem"
	"fmt"
)

func GetFingerprint(pemCert []byte) (string, error) {
	der, _ := pem.Decode([]byte(pemCert))
	if der == nil {
		return "", fmt.Errorf("Malformed cert: %s", pemCert)
	}
	hash := sha1.Sum(der.Bytes)
	hexified := make([][]byte, len(hash))
	for i, data := range hash {
		hexified[i] = []byte(fmt.Sprintf("%02X", data))
	}
	return fmt.Sprintf("SHA1 Fingerprint=%s", string(bytes.Join(hexified, []byte(":")))), nil
}
