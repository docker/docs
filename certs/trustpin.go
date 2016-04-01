package certs

import (
	"crypto/x509"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/notary"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/utils"
	"strings"
)

// TrustPinChecker handles logic for trust pinning for a gun, given a TrustPinConfig
type TrustPinChecker struct {
	mode   TrustPinMode
	gun    string
	config notary.TrustPinConfig
}

// TrustPinMode is a type to distinguish possible trust pinning modes
type TrustPinMode string

const (
	tofus TrustPinMode = "tofu"
	ca    TrustPinMode = "ca"
	certs TrustPinMode = "certs"
)

// NewTrustPinChecker returns a new TrustPinChecker from a TrustPinConfig for a GUN
func NewTrustPinChecker(trustPinConfig notary.TrustPinConfig, gun string) (TrustPinChecker, error) {
	trustPinChecker := TrustPinChecker{config: trustPinConfig, gun: gun}
	// Determine the mode, and if it's even valid
	if trustPinConfig.TOFU {
		trustPinChecker.mode = tofus
	} else if _, ok := trustPinConfig.Certs[gun]; ok {
		trustPinChecker.mode = certs
	} else if utils.ContainsKeyPrefix(trustPinConfig.CA, gun) {
		trustPinChecker.mode = ca
	} else {
		return TrustPinChecker{}, fmt.Errorf("no trust pinning specified")
	}
	return trustPinChecker, nil
}

func (t TrustPinChecker) checkCert(leafCert *x509.Certificate, intCerts []*x509.Certificate) bool {
	switch t.mode {
	case tofus:
		return true
	case certs:
		certID, err := trustmanager.FingerprintCert(leafCert)
		if err != nil {
			return false
		}
		return t.config.Certs[t.gun] == certID
	case ca:
		for caGunPrefix, caFilepath := range t.config.CA {
			if strings.HasPrefix(t.gun, caGunPrefix) {
				// Try to add the CA cert to our certificate store,
				// and use it to validate certs in the root.json later
				caCert, err := trustmanager.LoadCertFromFile(caFilepath)
				if err != nil {
					return false
				}
				if err = trustmanager.ValidateCertificate(caCert); err != nil {
					return false
				}
				// Now only consider certificates that are direct children from this CA cert, overwriting allValidCerts
				caRootPool := x509.NewCertPool()
				caRootPool.AddCert(caCert)
				if err != nil {
					logrus.Debugf("error retrieving valid leaf certificates for: %s, %v", t.gun, err)
					return false
				}

				// Use intermediate certificates included in the root TUF metadata for our validation
				caIntPool := x509.NewCertPool()
				for _, intCert := range intCerts {
					caIntPool.AddCert(intCert)
				}
				// Attempt to find a valid certificate chain from the leaf cert to CA root
				// Use this certificate if such a valid chain exists (possibly using intermediates)
				if _, err = leafCert.Verify(x509.VerifyOptions{Roots: caRootPool, Intermediates: caIntPool}); err == nil {
					return true
				}
			}
		}
	}
	return false
}
