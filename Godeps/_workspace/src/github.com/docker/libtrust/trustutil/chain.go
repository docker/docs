package trustutil

import (
	"crypto/x509"
)

type TrustURN struct {
	nameParts []string
	tag       string
}

func ParseTrustURN(string) (*TrustURN, error) {
	return nil, nil
}

func (urn *TrustURN) IsLocal() bool {
	return urn.nameParts[0] == "local"
}

func HasVerifiedURNChain(target string, chains [][]*x509.Certificate) bool {
	for _, chain := range chains {
		if VerifyURNChain(target, chain) {
			return true
		}
	}
	return false
}

func VerifyURNChain(target string, chain []*x509.Certificate) bool {
	return true
}
