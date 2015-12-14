package snapshot

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"

	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

// GetOrCreateSnapshotKey either creates a new snapshot key, or returns
// the existing one. Only the PublicKey is returned. The private part
// is held by the CryptoService.
func GetOrCreateSnapshotKey(gun string, store storage.KeyStore, crypto signed.CryptoService, createAlgorithm string) (data.PublicKey, error) {
	keyAlgorithm, public, err := store.GetKey(gun, data.CanonicalSnapshotRole)
	if err == nil {
		return data.NewPublicKey(keyAlgorithm, public), nil
	}

	if _, ok := err.(*storage.ErrNoKey); ok {
		key, err := crypto.Create("snapshot", createAlgorithm)
		if err != nil {
			return nil, err
		}
		logrus.Debug("Creating new snapshot key for ", gun, ". With algo: ", key.Algorithm())
		err = store.SetKey(gun, data.CanonicalSnapshotRole, key.Algorithm(), key.Public())
		if err == nil {
			return key, nil
		}

		if _, ok := err.(*storage.ErrKeyExists); ok {
			keyAlgorithm, public, err = store.GetKey(gun, data.CanonicalSnapshotRole)
			if err != nil {
				return nil, err
			}
			return data.NewPublicKey(keyAlgorithm, public), nil
		}
		return nil, err
	}
	return nil, err
}

// GetOrCreateSnapshot either returns the exisiting latest snapshot, or uses
// whatever the most recent snapshot is to create the next one, only updating
// the expiry time and version.
func GetOrCreateSnapshot(gun string, store storage.MetaStore, cryptoService signed.CryptoService) ([]byte, error) {

	d, err := store.GetCurrent(gun, "snapshot")
	if err != nil {
		return nil, err
	}

	sn := &data.SignedSnapshot{}
	if d != nil {
		err := json.Unmarshal(d, sn)
		if err != nil {
			logrus.Error("Failed to unmarshal existing snapshot")
			return nil, err
		}

		if !snapshotExpired(sn) {
			return d, nil
		}
	}

	sgnd, version, err := createSnapshot(gun, sn, store, cryptoService)
	if err != nil {
		logrus.Error("Failed to create a new snapshot")
		return nil, err
	}
	out, err := json.Marshal(sgnd)
	if err != nil {
		logrus.Error("Failed to marshal new snapshot")
		return nil, err
	}
	err = store.UpdateCurrent(gun, storage.MetaUpdate{Role: "snapshot", Version: version, Data: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// snapshotExpired simply checks if the snapshot is past its expiry time
func snapshotExpired(sn *data.SignedSnapshot) bool {
	return signed.IsExpired(sn.Signed.Expires)
}

// createSnapshot uses an existing snapshot to create a new one.
// Important things to be aware of:
//   - It requires that a snapshot already exists. We create snapshots
//     on upload so there should always be an existing snapshot if this
//     gets called.
//   - It doesn't update what roles are present in the snapshot, as those
//     were validated during upload.
func createSnapshot(gun string, sn *data.SignedSnapshot, store storage.MetaStore, cryptoService signed.CryptoService) (*data.Signed, int, error) {
	algorithm, public, err := store.GetKey(gun, data.CanonicalSnapshotRole)
	if err != nil {
		// owner of gun must have generated a snapshot key otherwise
		// we won't proceed with generating everything.
		return nil, 0, err
	}
	key := data.NewPublicKey(algorithm, public)

	// update version and expiry
	sn.Signed.Version = sn.Signed.Version + 1
	sn.Signed.Expires = data.DefaultExpires(data.CanonicalSnapshotRole)

	out, err := sn.ToSigned()
	if err != nil {
		return nil, 0, err
	}
	err = signed.Sign(cryptoService, out, key)
	if err != nil {
		return nil, 0, err
	}
	return out, sn.Signed.Version, nil
}
