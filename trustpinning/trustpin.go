package trustpinning

import (
	"crypto/x509"
	"fmt"
	"github.com/docker/notary"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/utils"
	"strings"
)

// TrustPinChecker handles logic for trust pinning for a gun, given a TrustPinConfig
type TrustPinChecker struct {
	mode         TrustPinMode
	gun          string
	config       notary.TrustPinConfig
	pinnedCAPool *x509.CertPool
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
	if _, ok := trustPinConfig.Certs[gun]; ok {
		trustPinChecker.mode = certs
		return trustPinChecker, nil
	}

	if caFilepath, err := trustPinChecker.getCAFilepathByPrefix(gun); err == nil {
		trustPinChecker.mode = ca
		// Try to add the CA certs from its bundle file to our certificate store,
		// and use it to validate certs in the root.json later
		caCerts, err := trustmanager.LoadCertBundleFromFile(caFilepath)
		if err != nil {
			return TrustPinChecker{}, fmt.Errorf("could not load root cert from CA path")
		}
		// Now only consider certificates that are direct children from this CA cert, overwriting allValidCerts
		caRootPool := x509.NewCertPool()
		for _, caCert := range caCerts {
			if err = trustmanager.ValidateCertificate(caCert); err != nil {
				continue
			}
			caRootPool.AddCert(caCert)
		}
		// If we didn't have any valid CA certs, error out
		if len(caRootPool.Subjects()) == 0 {
			return TrustPinChecker{}, fmt.Errorf("invalid CA certs provided")
		}
		trustPinChecker.pinnedCAPool = caRootPool
		return trustPinChecker, nil
	}

	if trustPinConfig.TOFU {
		trustPinChecker.mode = tofus
		return trustPinChecker, nil
	}
	return TrustPinChecker{}, fmt.Errorf("invalid trust pinning specified")
}

func (t TrustPinChecker) checkCert(leafCert *x509.Certificate, intCerts []*x509.Certificate) bool {
	switch t.mode {
	case certs:
		certID, err := trustmanager.FingerprintCert(leafCert)
		if err != nil {
			return false
		}
		return utils.StrSliceContains(t.config.Certs[t.gun], certID)
	case ca:
		// Use intermediate certificates included in the root TUF metadata for our validation
		caIntPool := x509.NewCertPool()
		for _, intCert := range intCerts {
			caIntPool.AddCert(intCert)
		}
		// Attempt to find a valid certificate chain from the leaf cert to CA root
		// Use this certificate if such a valid chain exists (possibly using intermediates)
		if _, err := leafCert.Verify(x509.VerifyOptions{Roots: t.pinnedCAPool, Intermediates: caIntPool}); err == nil {
			return true
		}
	case tofus:
		return true
	}
	return false
}

// Will return the CA filepath corresponding to the most specific (longest) entry in the map that is still a prefix
// of the provided gun.  Returns an error if no entry matches this GUN as a prefix.
func (t TrustPinChecker) getCAFilepathByPrefix(gun string) (string, error) {
	specificGUN := ""
	specificCAFilepath := ""
	foundCA := false
	for gunPrefix, caFilepath := range t.config.CA {
		if strings.HasPrefix(gun, gunPrefix) && len(gunPrefix) >= len(specificGUN) {
			specificGUN = gunPrefix
			specificCAFilepath = caFilepath
			foundCA = true
		}
	}
	if !foundCA {
		return "", fmt.Errorf("could not find pinned CA for GUN: %s\n", gun)
	}
	return specificCAFilepath, nil
}
