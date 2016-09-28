package enzi

import (
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	libkv "github.com/docker/libkv/store"
	"github.com/docker/orca/enzi/jose"
)

// SigningKeyCacheMaxAge the TTL for JWT signing keys in the KV store.
const SigningKeyCacheMaxAge = 12 * time.Hour

// UCP's Service Signing Keys are stored in the following directory
// structure:
//     orca/v1/auth/
//         signing_keys/
//             {key_id} -> JWK
//             {key_id} -> JWK
//             ...
const kvStoreSigningKeysDir = "orca/v1/auth/signing_keys"

func kvStoreSigningKeyPath(keyID string) string {
	return path.Join(kvStoreSigningKeysDir, keyID)
}

func (a *Authenticator) ListSigningKeys() (*jose.JWKSet, error) {
	kvPairs, err := a.kvStore.List(kvStoreSigningKeysDir)
	if err != nil {
		return nil, err
	}

	signingKeys := make([]*jose.PublicKey, 0, len(kvPairs))
	for _, kvPair := range kvPairs {
		signingKey := new(jose.PublicKey)
		if err := signingKey.UnmarshalJSON(kvPair.Value); err != nil {
			log.Errorf("unable to decode signing key JWK: %s", err)
			continue
		}

		signingKeys = append(signingKeys, signingKey)
	}

	return &jose.JWKSet{
		Keys: signingKeys,
	}, nil
}

func (a *Authenticator) SaveSigningKey() error {
	jwk, err := a.signingKey.PublicKey.MarshalJSON()
	if err != nil {
		return fmt.Errorf("unable to encode signing key to JWK: %s", err)
	}

	// This will extend the TTL of the key if it already exists.
	if err := a.kvStore.Put(kvStoreSigningKeyPath(a.signingKey.ID), jwk, &libkv.WriteOptions{TTL: SigningKeyCacheMaxAge}); err != nil {
		return err
	}

	return nil
}
