package manager

import (
	"crypto/sha1"
	"encoding/pem"
	"fmt"
	"io"
	"path"

	log "github.com/Sirupsen/logrus"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libtrust"
	"github.com/nu7hatch/gouuid"

	"github.com/docker/orca/utils"
)

func LoadOrCreateTrustKey(kv kvstore.Store) (libtrust.PrivateKey, error) {
	// Check KV store first
	kvpair, err := kv.Get(path.Join(datastoreVersion, "_trust", "key"))
	if err == nil {
		return libtrust.UnmarshalPrivateKeyPEM([]byte(kvpair.Value))
	}

	// Wasn't found, so generate a new one
	log.Debug("Generating new libtrust key")
	trustKey, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return nil, fmt.Errorf("error generating key: %s", err)
	}
	block, err := trustKey.PEMBlock()
	if err != nil {
		// Should never happen
		return nil, err
	}
	if err := kv.Put(path.Join(datastoreVersion, "_trust", "key"), pem.EncodeToMemory(block), nil); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return nil, fmt.Errorf("Error saving key to KV store: %s", err)
	}

	return trustKey, nil
}

func GetKeyHash(data string) string {
	h := sha1.New()
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Helper routine to get a single value from the KV store
func (m DefaultManager) getKvValue(key string) (string, error) {
	kv := m.Datastore()
	if kvpair, err := kv.Get(key); kvpair != nil {
		return string(kvpair.Value), nil
	} else {
		return "", utils.MaybeWrapEtcdClusterErr(err)
	}
}

// genRandomUUID generates a random UUID V4 string
func (m DefaultManager) genRandomUUID() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u4.String(), nil
}
