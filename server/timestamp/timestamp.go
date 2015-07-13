package timestamp

import (
	"encoding/json"
	"time"

	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
	cjson "github.com/tent/canonical-json-go"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/server/storage"
)

// GetOrCreateTimestampKey returns the timestamp key for the gun. It uses the store to
// lookup an existing timestamp key and the crypto to generate a new one if none is
// found. It attempts to handle the race condition that may occur if 2 servers try to
// create the key at the same time by simply querying the store a second time if it
// receives a conflict when writing.
func GetOrCreateTimestampKey(gun string, store storage.MetaStore, crypto signed.CryptoService, fallBackAlgorithm data.KeyAlgorithm) (*data.TUFKey, error) {
	keyAlgorithm, public, err := store.GetTimestampKey(gun)
	if err == nil {
		return data.NewTUFKey(keyAlgorithm, public, nil), nil
	}

	if _, ok := err.(*storage.ErrNoKey); ok {
		key, err := crypto.Create("timestamp", fallBackAlgorithm)
		if err != nil {
			return nil, err
		}
		err = store.SetTimestampKey(gun, key.Algorithm(), key.Public())
		if err == nil {
			return &key.TUFKey, nil
		}

		if _, ok := err.(*storage.ErrTimestampKeyExists); ok {
			keyAlgorithm, public, err = store.GetTimestampKey(gun)
			if err != nil {
				return nil, err
			}
			return data.NewTUFKey(keyAlgorithm, public, nil), nil
		}
		return nil, err
	}
	return nil, err
}

// GetOrCreateTimestamp returns the current timestamp for the gun. This may mean
// a new timestamp is generated either because none exists, or because the current
// one has expired. Once generated, the timestamp is saved in the store.
func GetOrCreateTimestamp(gun string, store storage.MetaStore, signer *signed.Signer) ([]byte, error) {
	d, err := store.GetCurrent(gun, "timestamp")
	if err != nil {
		if _, ok := err.(*storage.ErrNotFound); !ok {
			// If we received an ErrNotFound, we're going to
			// generate the first timestamp, any other error
			// should be returned here
			return nil, err
		}
	}
	ts := &data.SignedTimestamp{}
	if d != nil {
		err := json.Unmarshal(d, ts)
		if err != nil {
			logrus.Error("Failed to unmarshal existing timestamp")
			return nil, err
		}
		if !timestampExpired(ts) {
			return d, nil
		}
	}
	sgnd, version, err := createTimestamp(gun, ts, store, signer)
	if err != nil {
		logrus.Error("Failed to create a new timestamp")
		return nil, err
	}
	out, err := json.Marshal(sgnd)
	if err != nil {
		logrus.Error("Failed to marshal new timestamp")
		return nil, err
	}
	err = store.UpdateCurrent(gun, "timestamp", version, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// timestampExpired compares the current time to the expiry time of the timestamp
func timestampExpired(ts *data.SignedTimestamp) bool {
	return time.Now().After(ts.Signed.Expires)
}

// createTimestamp creates a new timestamp. If a prev timestamp is provided, it
// is assumed this is the immediately previous one, and the new one will have a
// version number one higher than prev. The store is used to lookup the current
// snapshot, this function does not save the newly generated timestamp.
func createTimestamp(gun string, prev *data.SignedTimestamp, store storage.MetaStore, signer *signed.Signer) (*data.Signed, int, error) {
	algorithm, public, err := store.GetTimestampKey(gun)
	if err != nil {
		// owner of gun must have generated a timestamp key otherwise
		// we won't proceed with generating everything.
		return nil, 0, err
	}
	key := data.NewPublicKey(algorithm, public)
	snapshot, err := store.GetCurrent(gun, "snapshot")
	if err != nil {
		return nil, 0, err
	}
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
	sgndTs, err := cjson.Marshal(ts.Signed)
	if err != nil {
		return nil, 0, err
	}
	out := &data.Signed{
		Signatures: ts.Signatures,
		Signed:     sgndTs,
	}
	err = signer.Sign(out, key)
	if err != nil {
		return nil, 0, err
	}
	return out, ts.Signed.Version, nil
}
