package snapshot

import (
	"bytes"
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
// whatever the most recent snapshot is to generate a new one.
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

		// want to ensure we always execute both of these such that if snapExp == true,
		// we update the meta in preparation for resigning
		snapExp := snapshotExpired(sn)
		contExp := contentExpired(gun, sn, store)

		if !snapExp && !contExp {
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

// contentExpired checks to see if any of the roles already in the snapshot
// have been updated. It will update any roles that have changed as it goes
// so that we don't have to run through all this again a second time.
func contentExpired(gun string, sn *data.SignedSnapshot, store storage.MetaStore) bool {
	expired := false
	updatedMeta := make(data.Files)
	for role, meta := range sn.Signed.Meta {
		curr, err := store.GetCurrent(gun, role)
		if err != nil {
			return false
		}
		roleExp, newMeta := roleExpired(curr, meta)
		if roleExp {
			updatedMeta[role] = newMeta
		} else {
			updatedMeta[role] = meta
		}
		expired = expired || roleExp
	}
	if expired {
		sn.Signed.Meta = updatedMeta
	}
	return expired
}

// roleExpired checks if the content for a specific role differs from
// the snapshot
func roleExpired(roleData []byte, meta data.FileMeta) (bool, data.FileMeta) {
	currMeta, err := data.NewFileMeta(bytes.NewReader(roleData), "sha256")
	if err != nil {
		// if we can't generate FileMeta from the current roleData, we should
		// continue to serve the old role if it isn't time expired
		// because we won't be able to generate a new one.
		return false, data.FileMeta{}
	}
	hash := currMeta.Hashes["sha256"]
	return !bytes.Equal(hash, meta.Hashes["sha256"]), currMeta
}

// createSnapshot uses an existing snapshot to create a new one.
// Important things to be aware of:
//   - It requires that a snapshot already exists. We create snapshots
//     on upload so there should always be an existing snapshot if this
//     gets called.
//   - It doesn't update what roles are present in the snapshot, as those
//     were validated during upload. We also updated the hashes of the
//     already present roles as part of our checks on whether we could
//     serve the previous version of the snapshot
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
