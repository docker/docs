package timestamp

import (
	"bytes"
	"testing"
	"time"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/server/storage"
)

func TestTimestampExpired(t *testing.T) {
	ts := &data.SignedTimestamp{
		Signatures: nil,
		Signed: data.Timestamp{
			Expires: time.Now().AddDate(-1, 0, 0),
		},
	}
	assert.True(t, timestampExpired(ts), "Timestamp should have expired")
}

func TestTimestampNotExpired(t *testing.T) {
	ts := &data.SignedTimestamp{
		Signatures: nil,
		Signed: data.Timestamp{
			Expires: time.Now().AddDate(1, 0, 0),
		},
	}
	assert.False(t, timestampExpired(ts), "Timestamp should NOT have expired")
}

func TestGetTimestampKey(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()
	k, err := GetOrCreateTimestampKey("gun", store, crypto, data.ED25519Key)
	assert.Nil(t, err, "Expected nil error")
	assert.NotNil(t, k, "Key should not be nil")

	k2, err := GetOrCreateTimestampKey("gun", store, crypto, data.ED25519Key)

	assert.Nil(t, err, "Expected nil error")

	// trying to get the same key again should return the same value
	assert.Equal(t, k, k2, "Did not receive same key when attempting to recreate.")
	assert.NotNil(t, k2, "Key should not be nil")
}

// If there is no previous timestamp or the previous timestamp is corrupt, then
// even if everything else is in place, getting the timestamp fails
func TestGetTimestampNoPreviousTimestamp(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("gun")
	assert.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	assert.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	assert.NoError(t, err)

	for _, timestampJSON := range [][]byte{nil, []byte("invalid JSON")} {
		store := storage.NewMemStorage()

		// so we know it's not a failure in getting root or snapshot
		assert.NoError(t,
			store.UpdateCurrent("gun", storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
		assert.NoError(t,
			store.UpdateCurrent("gun", storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))

		if timestampJSON != nil {
			assert.NoError(t,
				store.UpdateCurrent("gun",
					storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 0, Data: timestampJSON}))
		}

		// create a key to be used by GetOrCreateTimestamp
		key, err := crypto.Create(data.CanonicalTimestampRole, "gun", data.ECDSAKey)
		assert.NoError(t, err)
		assert.NoError(t, store.SetKey("gun", data.CanonicalTimestampRole, key.Algorithm(), key.Public()))

		_, _, err = GetOrCreateTimestamp("gun", store, crypto)
		assert.Error(t, err, "GetTimestamp should have failed")
		if timestampJSON == nil {
			assert.IsType(t, storage.ErrNotFound{}, err)
		} else {
			assert.IsType(t, &json.SyntaxError{}, err)
		}
	}
}

// If there WAS a pre-existing timestamp, and it is not expired, then just return it (it doesn't
// load any other metadata that it doesn't need, like root)
func TestGetTimestampReturnsPreviousTimestampIfUnexpired(t *testing.T) {
	store := storage.NewMemStorage()
	repo, crypto, err := testutils.EmptyRepo("gun")
	assert.NoError(t, err)

	snapJSON, err := json.Marshal(repo.Snapshot)
	assert.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	assert.NoError(t, err)

	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 0, Data: timestampJSON}))

	_, gottenTimestamp, err := GetOrCreateTimestamp("gun", store, crypto)
	assert.NoError(t, err, "GetTimestamp should not have failed")
	assert.True(t, bytes.Equal(timestampJSON, gottenTimestamp))
}

func TestGetTimestampOldTimestampExpired(t *testing.T) {
	store := storage.NewMemStorage()
	repo, crypto, err := testutils.EmptyRepo("gun")
	assert.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	assert.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	assert.NoError(t, err)

	// create an expired timestamp
	_, err = repo.SignTimestamp(time.Now().AddDate(-1, -1, -1))
	assert.True(t, repo.Timestamp.Signed.Expires.Before(time.Now()))
	assert.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	assert.NoError(t, err)

	// set all the metadata
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 1, Data: timestampJSON}))

	_, gottenTimestamp, err := GetOrCreateTimestamp("gun", store, crypto)
	assert.NoError(t, err, "GetTimestamp errored")

	assert.False(t, bytes.Equal(timestampJSON, gottenTimestamp),
		"Timestamp was not regenerated when old one was expired")

	signedMeta := &data.SignedMeta{}
	assert.NoError(t, json.Unmarshal(gottenTimestamp, signedMeta))
	// the new metadata is not expired
	assert.True(t, signedMeta.Signed.Expires.After(time.Now()))
}

// In practice this might happen if the snapshot is expired, for instance, and
// is re-signed.
func TestGetTimestampIfNewSnapshot(t *testing.T) {
	store := storage.NewMemStorage()
	repo, crypto, err := testutils.EmptyRepo("gun")

	rootJSON, err := json.Marshal(repo.Root)
	assert.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	assert.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	assert.NoError(t, err)

	// set all the metadata
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 0, Data: timestampJSON}))

	c1, ts1, err := GetOrCreateTimestamp("gun", store, crypto)
	assert.Nil(t, err, "GetTimestamp errored")

	// update the snapshot to a new version
	repo.Snapshot.Signed.Version++
	snapJSON, err = json.Marshal(repo.Snapshot)
	assert.NoError(t, err)
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "snapshot", Version: 1, Data: snapJSON})

	c2, ts2, err := GetOrCreateTimestamp("gun", store, crypto)
	assert.NoError(t, err, "GetTimestamp errored")
	assert.NotEqual(t, ts1, ts2, "Timestamp was not regenerated when snapshot changed")
	assert.True(t, c1.Before(*c2), "Timestamp modification time incorrect")
}

// If the root or snapshot is missing or corrupt, no timestamp can be generated
func TestCannotMakeNewTimestampIfNoRootOrSnapshot(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("gun")
	assert.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	assert.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	assert.NoError(t, err)

	// create an expired timestamp
	_, err = repo.SignTimestamp(time.Now().AddDate(-1, -1, -1))
	assert.True(t, repo.Timestamp.Signed.Expires.Before(time.Now()))
	assert.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	assert.NoError(t, err)

	invalids := []map[string][]byte{
		{data.CanonicalRootRole: rootJSON, data.CanonicalSnapshotRole: []byte("invalid JSON")},
		{data.CanonicalRootRole: []byte("invalid JSON"), data.CanonicalSnapshotRole: snapJSON},
		{data.CanonicalRootRole: rootJSON},
		{data.CanonicalSnapshotRole: snapJSON},
	}

	for _, dataToSet := range invalids {
		store := storage.NewMemStorage()
		for roleName, jsonBytes := range dataToSet {
			assert.NoError(t, store.UpdateCurrent("gun",
				storage.MetaUpdate{Role: roleName, Version: 0, Data: jsonBytes}))
		}
		assert.NoError(t, store.UpdateCurrent("gun",
			storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 1, Data: timestampJSON}))

		_, _, err := GetOrCreateTimestamp("gun", store, crypto)
		assert.Error(t, err, "GetTimestamp errored")

		if len(dataToSet) == 1 { // missing metadata
			assert.IsType(t, storage.ErrNotFound{}, err)
		} else {
			assert.IsType(t, &json.SyntaxError{}, err)
		}
	}
}

func TestCreateTimestampNoKeyInCrypto(t *testing.T) {
	store := storage.NewMemStorage()
	repo, _, err := testutils.EmptyRepo("gun")
	assert.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	assert.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	assert.NoError(t, err)

	// create an expired timestamp
	_, err = repo.SignTimestamp(time.Now().AddDate(-1, -1, -1))
	assert.True(t, repo.Timestamp.Signed.Expires.Before(time.Now()))
	assert.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	assert.NoError(t, err)

	// set all the metadata so we know the failure to sign is just because of the key
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	assert.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 1, Data: timestampJSON}))

	// pass it a new cryptoservice without the key
	_, _, err = GetOrCreateTimestamp("gun", store, signed.NewEd25519())
	assert.Error(t, err)
	assert.IsType(t, signed.ErrNoKeys{}, err)
}
