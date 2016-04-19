package snapshot

import (
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf"
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
		key, err := crypto.Create(data.CanonicalSnapshotRole, gun, createAlgorithm)
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

// GetOrCreateSnapshot either returns the existing latest snapshot, or uses
// whatever the most recent snapshot is to generate the next one, only updating
// the expiry time and version.  Note that this function does not write generated
// snapshots to the underlying data store, and will either return the latest snapshot time
// or nil as the time modified
func GetOrCreateSnapshot(gun, checksum string, store storage.MetaStore, cryptoService signed.CryptoService) (
	*time.Time, []byte, error) {

	lastModified, currentJSON, err := store.GetChecksum(gun, data.CanonicalSnapshotRole, checksum)
	if err != nil {
		return nil, nil, err
	}

	prev := new(data.SignedSnapshot)
	if err := json.Unmarshal(currentJSON, prev); err != nil {
		logrus.Error("Failed to unmarshal existing snapshot for GUN ", gun)
		return nil, nil, err
	}

	if !snapshotExpired(prev) {
		return lastModified, currentJSON, nil
	}

	repo := tuf.NewRepo(cryptoService)

	// load the current root to ensure we use the correct snapshot key.
	_, rootJSON, err := store.GetCurrent(gun, data.CanonicalRootRole)

	if err != nil {
		return nil, nil, err
	}
	root := &data.SignedRoot{}
	if err := json.Unmarshal(rootJSON, root); err != nil {
		logrus.Error("Failed to unmarshal existing root for GUN ", gun)
		return nil, nil, err
	}
	repo.SetRoot(root)

	snapshotUpdate, err := NewSnapshotUpdate(prev, repo)
	if err != nil {
		logrus.Error("Failed to create a new snapshot")
		return nil, nil, err
	}
	return nil, snapshotUpdate.Data, nil
}

// snapshotExpired simply checks if the snapshot is past its expiry time
func snapshotExpired(sn *data.SignedSnapshot) bool {
	return signed.IsExpired(sn.Signed.Expires)
}

// NewSnapshotUpdate produces a new snapshot and returns it as a metadata update, given the
// previous snapshot and the TUF repo.
func NewSnapshotUpdate(prev *data.SignedSnapshot, repo *tuf.Repo) (*storage.MetaUpdate, error) {
	if prev != nil {
		repo.SetSnapshot(prev) // SetSnapshot never errors
	} else {
		// this will only occur if no snapshot has ever been created for the repository
		if err := repo.InitSnapshot(); err != nil {
			return nil, err
		}
	}
	sgnd, err := repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	if err != nil {
		return nil, err
	}
	sgndJSON, err := json.Marshal(sgnd)
	if err != nil {
		return nil, err
	}
	return &storage.MetaUpdate{
		Role:    data.CanonicalSnapshotRole,
		Version: repo.Snapshot.Signed.Version,
		Data:    sgndJSON,
	}, nil
}
