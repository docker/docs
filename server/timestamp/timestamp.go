package timestamp

import (
	"time"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/signed"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/server/snapshot"
	"github.com/docker/notary/server/storage"
)

// GetOrCreateTimestampKey returns the timestamp key for the gun. It uses the store to
// lookup an existing timestamp key and the crypto to generate a new one if none is
// found. It attempts to handle the race condition that may occur if 2 servers try to
// create the key at the same time by simply querying the store a second time if it
// receives a conflict when writing.
func GetOrCreateTimestampKey(gun string, store storage.MetaStore, crypto signed.CryptoService, createAlgorithm string) (data.PublicKey, error) {
	keyAlgorithm, public, err := store.GetKey(gun, data.CanonicalTimestampRole)
	if err == nil {
		return data.NewPublicKey(keyAlgorithm, public), nil
	}

	if _, ok := err.(*storage.ErrNoKey); ok {
		key, err := crypto.Create("timestamp", gun, createAlgorithm)
		if err != nil {
			return nil, err
		}
		logrus.Debug("Creating new timestamp key for ", gun, ". With algo: ", key.Algorithm())
		err = store.SetKey(gun, data.CanonicalTimestampRole, key.Algorithm(), key.Public())
		if err == nil {
			return key, nil
		}

		if _, ok := err.(*storage.ErrKeyExists); ok {
			keyAlgorithm, public, err = store.GetKey(gun, data.CanonicalTimestampRole)
			if err != nil {
				return nil, err
			}
			return data.NewPublicKey(keyAlgorithm, public), nil
		}
		return nil, err
	}
	return nil, err
}

// GetOrCreateTimestamp returns the current timestamp for the gun. This may mean
// a new timestamp is generated either because none exists, or because the current
// one has expired. Once generated, the timestamp is saved in the store.
func GetOrCreateTimestamp(gun string, store storage.MetaStore, cryptoService signed.CryptoService) (
	*time.Time, []byte, error) {

	_, snapshot, err := snapshot.GetOrCreateSnapshot(gun, store, cryptoService)
	if err != nil {
		return nil, nil, err
	}
	lastModified, d, err := store.GetCurrent(gun, data.CanonicalTimestampRole)
	if err != nil {
		if _, ok := err.(storage.ErrNotFound); !ok {
			logrus.Error("error retrieving timestamp: ", err.Error())
			return nil, nil, err
		}
		logrus.Debug("No timestamp found, will proceed to create first timestamp")
	}
	var ts *data.SignedTimestamp
	if d != nil {
		ts = &data.SignedTimestamp{}
		err := json.Unmarshal(d, ts)
		if err != nil {
			logrus.Error("Failed to unmarshal existing timestamp")
			return nil, nil, err
		}
		if !timestampExpired(ts) && !snapshotExpired(ts, snapshot) {
			return lastModified, d, nil
		}
	}
	sgnd, version, err := CreateTimestamp(gun, ts, snapshot, store, cryptoService)
	if err != nil {
		logrus.Error("Failed to create a new timestamp")
		return nil, nil, err
	}
	out, err := json.Marshal(sgnd)
	if err != nil {
		logrus.Error("Failed to marshal new timestamp")
		return nil, nil, err
	}
	err = store.UpdateCurrent(gun, storage.MetaUpdate{Role: "timestamp", Version: version, Data: out})
	if err != nil {
		return nil, nil, err
	}
	c := time.Now()
	return &c, out, nil
}

// timestampExpired compares the current time to the expiry time of the timestamp
func timestampExpired(ts *data.SignedTimestamp) bool {
	return signed.IsExpired(ts.Signed.Expires)
}

// snapshotExpired verifies the checksum(s) for the given snapshot using metadata from the timestamp
func snapshotExpired(ts *data.SignedTimestamp, snapshot []byte) bool {
	// If this check failed, it means the current snapshot was not exactly what we expect
	// via the timestamp. So we can consider it to be "expired."
	return data.CheckHashes(snapshot, ts.Signed.Meta[data.CanonicalSnapshotRole].Hashes) != nil
}

// CreateTimestamp creates a new timestamp. If a prev timestamp is provided, it
// is assumed this is the immediately previous one, and the new one will have a
// version number one higher than prev. The store is used to lookup the current
// snapshot, this function does not save the newly generated timestamp.
func CreateTimestamp(gun string, prev *data.SignedTimestamp, snapshot []byte, store storage.MetaStore, cryptoService signed.CryptoService) (*data.Signed, int, error) {
	kdb := keys.NewDB()
	repo := tuf.NewRepo(kdb, cryptoService)

	// load the current root to ensure we use the correct timestamp key.
	root, err := store.GetCurrent(gun, data.CanonicalRootRole)
	r := &data.SignedRoot{}
	err = json.Unmarshal(root, r)
	if err != nil {
		// couldn't parse root
		return nil, 0, err
	}
	repo.SetRoot(r)

	// load snapshot so we can include it in timestamp
	sn := &data.SignedSnapshot{}
	err = json.Unmarshal(snapshot, sn)
	if err != nil {
		// couldn't parse snapshot
		return nil, 0, err
	}
	repo.SetSnapshot(sn)

	if prev == nil {
		// no previous timestamp: generate first timestamp
		repo.InitTimestamp()
	} else {
		// set repo timestamp to previous timestamp to use as base for
		// generating new one
		repo.SetTimestamp(prev)
	}

	out, err := repo.SignTimestamp(
		data.DefaultExpires(data.CanonicalTimestampRole),
	)
	if err != nil {
		return nil, 0, err
	}
	return out, repo.Timestamp.Signed.Version, nil
}
