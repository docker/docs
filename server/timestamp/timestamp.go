package timestamp

import (
	"bytes"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/tuf/data"
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
		key, err := crypto.Create("timestamp", createAlgorithm)
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
func GetOrCreateTimestamp(gun string, store storage.MetaStore, cryptoService signed.CryptoService) ([]byte, error) {
	snapshot, err := snapshot.GetOrCreateSnapshot(gun, store, cryptoService)
	if err != nil {
		return nil, err
	}
	d, err := store.GetCurrent(gun, "timestamp")
	if err != nil {
		if _, ok := err.(storage.ErrNotFound); !ok {
			logrus.Error("error retrieving timestamp: ", err.Error())
			return nil, err
		}
		logrus.Debug("No timestamp found, will proceed to create first timestamp")
	}
	ts := &data.SignedTimestamp{}
	if d != nil {
		err := json.Unmarshal(d, ts)
		if err != nil {
			logrus.Error("Failed to unmarshal existing timestamp")
			return nil, err
		}
		if !timestampExpired(ts) && !snapshotExpired(ts, snapshot) {
			return d, nil
		}
	}
	sgnd, version, err := CreateTimestamp(gun, ts, snapshot, store, cryptoService)
	if err != nil {
		logrus.Error("Failed to create a new timestamp")
		return nil, err
	}
	out, err := json.Marshal(sgnd)
	if err != nil {
		logrus.Error("Failed to marshal new timestamp")
		return nil, err
	}
	err = store.UpdateCurrent(gun, storage.MetaUpdate{Role: "timestamp", Version: version, Data: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// timestampExpired compares the current time to the expiry time of the timestamp
func timestampExpired(ts *data.SignedTimestamp) bool {
	return signed.IsExpired(ts.Signed.Expires)
}

func snapshotExpired(ts *data.SignedTimestamp, snapshot []byte) bool {
	meta, err := data.NewFileMeta(bytes.NewReader(snapshot), "sha256")
	if err != nil {
		// if we can't generate FileMeta from the current snapshot, we should
		// continue to serve the old timestamp if it isn't time expired
		// because we won't be able to generate a new one.
		return false
	}
	hash := meta.Hashes["sha256"]
	return !bytes.Equal(hash, ts.Signed.Meta["snapshot"].Hashes["sha256"])
}

// CreateTimestamp creates a new timestamp. If a prev timestamp is provided, it
// is assumed this is the immediately previous one, and the new one will have a
// version number one higher than prev. The store is used to lookup the current
// snapshot, this function does not save the newly generated timestamp.
func CreateTimestamp(gun string, prev *data.SignedTimestamp, snapshot []byte, store storage.MetaStore, cryptoService signed.CryptoService) (*data.Signed, int, error) {
	algorithm, public, err := store.GetKey(gun, data.CanonicalTimestampRole)
	if err != nil {
		// owner of gun must have generated a timestamp key otherwise
		// we won't proceed with generating everything.
		return nil, 0, err
	}
	key := data.NewPublicKey(algorithm, public)
	sn := &data.Signed{}
	err = json.Unmarshal(snapshot, sn)
	if err != nil {
		// couldn't parse snapshot
		return nil, 0, err
	}
	ts, err := data.NewTimestamp(sn)
	if err != nil {
		return nil, 0, err
	}
	if prev != nil {
		ts.Signed.Version = prev.Signed.Version + 1
	}
	sgndTs, err := json.MarshalCanonical(ts.Signed)
	if err != nil {
		return nil, 0, err
	}
	out := &data.Signed{
		Signatures: ts.Signatures,
		Signed:     sgndTs,
	}
	err = signed.Sign(cryptoService, out, key)
	if err != nil {
		return nil, 0, err
	}
	return out, ts.Signed.Version, nil
}
