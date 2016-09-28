package config

import (
	"encoding/pem"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libtrust"
)

// Generate a new ID, storing it in the config
func GetNewID() error {
	trustKey, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return fmt.Errorf("error generating key: %s", err)
	}
	pemKey, err := trustKey.PEMBlock()
	if err != nil {
		return fmt.Errorf("error generating PEM key: %s", err)
	}

	OrcaInstanceID = trustKey.PublicKey().KeyID()
	log.Debugf(`New UCP Instance ID will be "%s"`, OrcaInstanceID)
	OrcaInstanceKey = string(pem.EncodeToMemory(pemKey))
	return nil
}
