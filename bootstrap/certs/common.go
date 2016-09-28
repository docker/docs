package certs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/docker/orca/bootstrap/config"
	caconfig "github.com/docker/orca/ca/config"
)

func WriteCert(filename string, content []byte, perms os.FileMode) error {
	log.Debugf("Writing %s", filename)
	if err := ioutil.WriteFile(filename, content, perms); err != nil {
		return err
	}
	return nil
}

func GenerateCSR(cn, ou string, hosts []string) ([]byte, []byte, error) {
	req := &csr.CertificateRequest{
		CN: cn,
		KeyRequest: &csr.BasicKeyRequest{
			A: config.KeyAlgo,
			S: config.KeySize,
		},
		Hosts: hosts,
	}

	if ou != "" {
		req.Names = []csr.Name{{OU: ou}}
	}

	log.Debugf(`CSR Generated for "%s" with hostnames %v`, cn, hosts)
	return csr.ParseRequest(req)

}

func InitLocalNode(caMount, certMount, cn, ou, nodeName string, hostnames []string, uid, gid int) error {
	// Check to see if we can skip regeneration
	if existing, err := verifyExisting(certMount, false); err != nil {
		// We could just clobber them and regenerate, but that might mask real failures
		return fmt.Errorf("Existing certs for %s appear to be corrupt: %s", nodeName, err)
	} else if existing {
		log.Debugf("Reusing existing certs for %s", nodeName)
		return nil
	} // else proceed and generate them

	s, err := local.NewSignerFromFile(
		filepath.Join(caMount, config.CertFilename),
		filepath.Join(caMount, config.KeyFilename), caconfig.Signing())
	if err != nil {
		log.Debug("Failed to load signer")
		return err
	}

	csr, key, err := GenerateCSR(cn, ou, hostnames)
	if err != nil {
		log.Debug("Failed to parse csr")
		return err
	}

	cert, err := s.Sign(signer.SignRequest{
		Request: string(csr),
		Profile: "node",
	})
	if err != nil {
		log.Debug("Sign failure")
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(certMount, config.CertFilename), cert, 0644); err != nil {
		return err
	}
	if err := os.Chown(filepath.Join(certMount, config.CertFilename), uid, gid); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(certMount, config.KeyFilename), key, 0600); err != nil {
		return err
	}
	if err := os.Chown(filepath.Join(certMount, config.KeyFilename), uid, gid); err != nil {
		return err
	}

	caChain, err := ioutil.ReadFile(filepath.Join(caMount, config.CertFilename))
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(certMount, config.CAFilename), caChain, 0644); err != nil {
		return err
	}
	if err := os.Chown(filepath.Join(certMount, config.CAFilename), uid, gid); err != nil {
		return err
	}

	return nil
}

func verifyExisting(mount string, isCA bool) (bool, error) {
	// Check for existing files, skip init if all present
	expected := []string{
		config.CertFilename,
		config.KeyFilename,
	}
	if !isCA {
		expected = append(expected, config.CAFilename)
	}
	found := []string{}
	for _, filename := range expected {
		if _, err := os.Stat(filepath.Join(mount, filename)); err == nil {
			found = append(found, filename)
		}
	}
	if len(found) == len(expected) {
		log.Debugf("Reusing existing certs in %s", mount)
		return true, nil
	} else if len(found) > 0 {
		return false, fmt.Errorf("Missing one or more files in %s, unable to re-use.  Found: %v, expected %v",
			mount, found, expected)
	}
	return false, nil
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
