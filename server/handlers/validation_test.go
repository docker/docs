package handlers

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/docker/notary/tuf/validation"
	"github.com/stretchr/testify/require"
)

// this is a fake storage that serves errors
type getFailStore struct {
	errsToReturn map[string]error
	storage.MetaStore
}

// GetCurrent returns the current metadata, or an error depending on whether
// getFailStore is configured to return an error for this role
func (f getFailStore) GetCurrent(gun, tufRole string) (*time.Time, []byte, error) {
	err := f.errsToReturn[tufRole]
	if err == nil {
		return f.MetaStore.GetCurrent(gun, tufRole)
	}
	return nil, nil, err
}

// GetChecksum returns the metadata with this checksum, or an error depending on
// whether getFailStore is configured to return an error for this role
func (f getFailStore) GetChecksum(gun, tufRole, checksum string) (*time.Time, []byte, error) {
	err := f.errsToReturn[tufRole]
	if err == nil {
		return f.MetaStore.GetChecksum(gun, tufRole, checksum)
	}
	return nil, nil, err
}

func copyKeys(t *testing.T, from signed.CryptoService, roles ...string) signed.CryptoService {
	memKeyStore := trustmanager.NewKeyMemoryStore(passphrase.ConstantRetriever("pass"))
	for _, role := range roles {
		for _, keyID := range from.ListKeys(role) {
			key, _, err := from.GetPrivateKey(keyID)
			require.NoError(t, err)
			memKeyStore.AddKey(trustmanager.KeyInfo{Role: role}, key)
		}
	}
	return cryptoservice.NewCryptoService(memKeyStore)
}

// Returns a mapping of role name to `MetaUpdate` objects
func getUpdates(r, tg, sn, ts *data.Signed) (
	root, targets, snapshot, timestamp storage.MetaUpdate, err error) {

	rs, tgs, sns, tss, err := testutils.Serialize(r, tg, sn, ts)
	if err != nil {
		return
	}

	root = storage.MetaUpdate{
		Role:    data.CanonicalRootRole,
		Version: 1,
		Data:    rs,
	}
	targets = storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgs,
	}
	snapshot = storage.MetaUpdate{
		Role:    data.CanonicalSnapshotRole,
		Version: 1,
		Data:    sns,
	}
	timestamp = storage.MetaUpdate{
		Role:    data.CanonicalTimestampRole,
		Version: 1,
		Data:    tss,
	}
	return
}

func TestValidateEmptyNew(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)

	// we generated our own timestamp, and did not take the other timestamp,
	// but all other metadata should come from updates
	founds := make(map[string]bool)
	for _, update := range updates {
		switch update.Role {
		case data.CanonicalRootRole:
			require.True(t, bytes.Equal(update.Data, root.Data))
			founds[data.CanonicalRootRole] = true
		case data.CanonicalSnapshotRole:
			require.True(t, bytes.Equal(update.Data, snapshot.Data))
			founds[data.CanonicalSnapshotRole] = true
		case data.CanonicalTargetsRole:
			require.True(t, bytes.Equal(update.Data, targets.Data))
			founds[data.CanonicalTargetsRole] = true
		case data.CanonicalTimestampRole:
			require.False(t, bytes.Equal(update.Data, timestamp.Data))
			founds[data.CanonicalTimestampRole] = true
		}
	}
	require.Len(t, founds, 4)
}

func TestValidatePrevTimestamp(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot}

	store := storage.NewMemStorage()
	store.UpdateCurrent("testGUN", timestamp)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)

	// we generated our own timestamp, and did not take the other timestamp,
	// but all other metadata should come from updates
	var foundTimestamp bool
	for _, update := range updates {
		if update.Role == data.CanonicalTimestampRole {
			foundTimestamp = true
			oldTimestamp, newTimestamp := &data.SignedTimestamp{}, &data.SignedTimestamp{}
			require.NoError(t, json.Unmarshal(timestamp.Data, oldTimestamp))
			require.NoError(t, json.Unmarshal(update.Data, newTimestamp))
			require.Equal(t, oldTimestamp.Signed.Version+1, newTimestamp.Signed.Version)
		}
	}
	require.True(t, foundTimestamp)
}

func TestValidatePreviousTimestampCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot}

	// have corrupt timestamp data in the storage already
	store := storage.NewMemStorage()
	timestamp.Data = timestamp.Data[1:]
	store.UpdateCurrent("testGUN", timestamp)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, &json.SyntaxError{}, err)
}

func TestValidateGetCurrentTimestampBroken(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot}

	store := getFailStore{
		MetaStore:    storage.NewMemStorage(),
		errsToReturn: map[string]error{data.CanonicalTimestampRole: data.ErrNoSuchRole{}},
	}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, data.ErrNoSuchRole{}, err)
}

func TestValidateNoNewRoot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store.UpdateCurrent("testGUN", root)
	updates := []storage.MetaUpdate{targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)
}

func TestValidateNoNewTargets(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store.UpdateCurrent("testGUN", targets)
	updates := []storage.MetaUpdate{root, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)
}

func TestValidateOnlySnapshot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store.UpdateCurrent("testGUN", root)
	store.UpdateCurrent("testGUN", targets)

	updates := []storage.MetaUpdate{snapshot}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)
}

func TestValidateOldRoot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store.UpdateCurrent("testGUN", root)
	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)
}

func TestValidateOldRootCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	badRoot := storage.MetaUpdate{
		Version: root.Version,
		Role:    root.Role,
		Data:    root.Data[1:],
	}
	store.UpdateCurrent("testGUN", badRoot)
	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, &json.SyntaxError{}, err)
}

// We cannot validate a new root if the old root is corrupt, because there might
// have been a root key rotation.
func TestValidateOldRootCorruptRootRole(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// so a valid root, but missing the root role
	signedRoot, err := data.RootFromSigned(r)
	require.NoError(t, err)
	delete(signedRoot.Signed.Roles, data.CanonicalRootRole)
	badRootJSON, err := json.Marshal(signedRoot)
	require.NoError(t, err)
	badRoot := storage.MetaUpdate{
		Version: root.Version,
		Role:    root.Role,
		Data:    badRootJSON,
	}
	store.UpdateCurrent("testGUN", badRoot)
	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

// We cannot validate a new root if we cannot get the old root from the DB (
// and cannot detect whether there was an old root or not), because there might
// have been an old root and we can't determine if the new root represents a
// root key rotation.
func TestValidateRootGetCurrentRootBroken(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := getFailStore{
		MetaStore:    storage.NewMemStorage(),
		errsToReturn: map[string]error{data.CanonicalRootRole: data.ErrNoSuchRole{}},
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, data.ErrNoSuchRole{}, err)
}

// A valid root rotation only cares about the immediately previous old root keys,
// whether or not there are old root roles
func TestValidateRootRotationWithOldSigs(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	serverCrypto := copyKeys(t, crypto, data.CanonicalTimestampRole)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// set the original root in the store
	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}
	require.NoError(t, store.UpdateMany("testGUN", updates))

	// rotate the root key, sign with both keys, and update - update should succeed
	newRootKey, err := crypto.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	newRootID := newRootKey.ID()

	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalRootRole, newRootKey))
	r, _, sn, _, err = testutils.Sign(repo)
	require.NoError(t, err)
	root, _, snapshot, _, err = getUpdates(r, tg, sn, ts)
	require.NoError(t, err)
	root.Version = repo.Root.Signed.Version
	snapshot.Version = repo.Snapshot.Signed.Version

	updates, err = validateUpdate(serverCrypto, "testGUN", []storage.MetaUpdate{root, snapshot}, store)
	require.NoError(t, err)
	require.NoError(t, store.UpdateMany("testGUN", updates))

	// the next root does NOT need to be signed by both keys, because we only care
	// about signing with both keys if the root keys have changed (signRoot again to bump the version)

	r, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	// delete all signatures except the one with the new key
	for _, sig := range repo.Root.Signatures {
		if sig.KeyID == newRootID {
			r.Signatures = []data.Signature{sig}
			repo.Root.Signatures = r.Signatures
			break
		}
	}
	sn, err = repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	require.NoError(t, err)

	root, _, snapshot, _, err = getUpdates(r, tg, sn, ts)
	require.NoError(t, err)
	root.Version = repo.Root.Signed.Version
	snapshot.Version = repo.Snapshot.Signed.Version
	updates, err = validateUpdate(serverCrypto, "testGUN", []storage.MetaUpdate{root, snapshot}, store)
	require.NoError(t, err)
	require.NoError(t, store.UpdateMany("testGUN", updates))

	// another root rotation requires only the previous and new keys, and not the
	// original root key even though that original role is still in the metadata

	newRootKey2, err := crypto.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	newRootID2 := newRootKey2.ID()

	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalRootRole, newRootKey2))
	r, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	// delete all signatures except the ones with the first and second new keys
	sigs := make([]data.Signature, 0, 2)
	for _, sig := range repo.Root.Signatures {
		if sig.KeyID == newRootID || sig.KeyID == newRootID2 {
			sigs = append(sigs, sig)
		}
	}
	require.Len(t, sigs, 2)
	repo.Root.Signatures = sigs
	r.Signatures = sigs

	sn, err = repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	require.NoError(t, err)

	root, _, snapshot, _, err = getUpdates(r, tg, sn, ts)
	require.NoError(t, err)
	root.Version = repo.Root.Signed.Version
	snapshot.Version = repo.Snapshot.Signed.Version
	_, err = validateUpdate(serverCrypto, "testGUN", []storage.MetaUpdate{root, snapshot}, store)
	require.NoError(t, err)
}

// A valid root rotation requires the immediately previous root ROLE be satisfied,
// not just that there is a single root signature.  So if there were 2 keys, either
// of which can sign the root rotation, then either one of those keys can be used
// to sign the root rotation - not necessarily the one that signed the previous root.
func TestValidateRootRotationMultipleKeysThreshold1(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	serverCrypto := copyKeys(t, crypto, data.CanonicalTimestampRole)
	store := storage.NewMemStorage()

	// add a new root key to the root so that either can sign have to sign
	additionalRootKey, err := crypto.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	additionalRootID := additionalRootKey.ID()
	repo.Root.Signed.Keys[additionalRootID] = additionalRootKey
	repo.Root.Signed.Roles[data.CanonicalRootRole].KeyIDs = append(
		repo.Root.Signed.Roles[data.CanonicalRootRole].KeyIDs, additionalRootID)

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	require.Len(t, r.Signatures, 2)
	// make sure the old root was just signed with the first key
	for _, sig := range r.Signatures {
		if sig.KeyID != additionalRootID {
			r.Signatures = []data.Signature{sig}
			break
		}
	}
	require.Len(t, r.Signatures, 1)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// set the original root in the store
	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}
	require.NoError(t, store.UpdateMany("testGUN", updates))

	// replace the keys with just 1 key
	rotatedRootKey, err := crypto.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	rotatedRootID := rotatedRootKey.ID()
	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalRootRole, rotatedRootKey))

	r, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	require.Len(t, r.Signatures, 3)
	// delete all signatures except the additional key (which didn't sign the
	// previous root) and the new key
	sigs := make([]data.Signature, 0, 2)
	for _, sig := range repo.Root.Signatures {
		if sig.KeyID == additionalRootID || sig.KeyID == rotatedRootID {
			sigs = append(sigs, sig)
		}
	}
	require.Len(t, sigs, 2)
	repo.Root.Signatures = sigs
	r.Signatures = sigs

	sn, err = repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	require.NoError(t, err)

	root, _, snapshot, _, err = getUpdates(r, tg, sn, ts)
	require.NoError(t, err)
	root.Version = repo.Root.Signed.Version
	snapshot.Version = repo.Snapshot.Signed.Version
	_, err = validateUpdate(serverCrypto, "testGUN", []storage.MetaUpdate{root, snapshot}, store)
	require.NoError(t, err)
}

// A root rotation must be signed with old and new root keys, otherwise the
// new root fails to validate
func TestRootRotationNotSignedWithOldKeys(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store.UpdateCurrent("testGUN", root)

	rootKey, err := crypto.Create("root", "testGUN", data.ED25519Key)
	require.NoError(t, err)
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil)
	require.NoError(t, err)

	repo.Root.Signed.Roles["root"] = &rootRole.RootRole
	repo.Root.Signed.Keys[rootKey.ID()] = rootKey

	r, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	err = signed.Sign(crypto, r, []data.PublicKey{rootKey}, 1, nil)
	require.NoError(t, err)

	rt, err := data.RootFromSigned(r)
	require.NoError(t, err)
	repo.SetRoot(rt)

	sn, err = repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err = getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, crypto, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.Contains(t, err.Error(), "new root was not signed with at least 1 old keys")
}

// An update is not valid without the root metadata.
func TestValidateNoRoot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	_, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrValidation{}, err)
}

func TestValidateSnapshotMissingNoSnapshotKey(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadHierarchy{}, err)
}

func TestValidateSnapshotGenerateNoPrev(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	require.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		require.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)
}

func TestValidateSnapshotGenerateWithPrev(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	require.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		require.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	// set the current snapshot in the store manually so we find it when generating
	// the next version
	store.UpdateCurrent("testGUN", snapshot)

	prev, err := data.SnapshotFromSigned(sn)
	require.NoError(t, err)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)

	for _, u := range updates {
		if u.Role == data.CanonicalSnapshotRole {
			curr := &data.SignedSnapshot{}
			err = json.Unmarshal(u.Data, curr)
			require.Equal(t, prev.Signed.Version+1, curr.Signed.Version)
			require.Equal(t, u.Version, curr.Signed.Version)
		}
	}
}

func TestValidateSnapshotGeneratePrevCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	require.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		require.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	// corrupt the JSON structure of prev snapshot
	snapshot.Data = snapshot.Data[1:]
	// set the current snapshot in the store manually so we find it when generating
	// the next version
	store.UpdateCurrent("testGUN", snapshot)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, &json.SyntaxError{}, err)
}

// Store is broken when getting the current snapshot
func TestValidateSnapshotGenerateStoreGetCurrentSnapshotBroken(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := getFailStore{
		MetaStore:    storage.NewMemStorage(),
		errsToReturn: map[string]error{data.CanonicalSnapshotRole: data.ErrNoSuchRole{}},
	}
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	require.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		require.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, data.ErrNoSuchRole{}, err)
}

func TestValidateSnapshotGenerateNoTargets(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	require.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		require.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, _, _, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
}

func TestValidateSnapshotGenerate(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	require.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		require.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{targets}

	store.UpdateCurrent("testGUN", root)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.NoError(t, err)
}

// If there is no timestamp key in the store, validation fails.  This could
// happen if pushing an existing repository from one server to another that
// does not have the repo.
func TestValidateRootNoTimestampKey(t *testing.T) {
	oldRepo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{root, targets, snapshot}

	// do not copy the targets key to the storage, and try to update the root
	serverCrypto := signed.NewEd25519()
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)

	// there should still be no timestamp keys - one should not have been
	// created
	require.Empty(t, serverCrypto.ListAllKeys())
}

// If the timestamp key in the store does not match the timestamp key in
// the root.json, validation fails.  This could happen if pushing an existing
// repository from one server to another that had already initialized the same
// repo.
func TestValidateRootInvalidTimestampKey(t *testing.T) {
	oldRepo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{root, targets, snapshot}

	serverCrypto := signed.NewEd25519()
	_, err = serverCrypto.Create(data.CanonicalTimestampRole, "testGUN", data.ED25519Key)
	require.NoError(t, err)

	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

// If the timestamp role has a threshold > 1, validation fails.
func TestValidateRootInvalidTimestampThreshold(t *testing.T) {
	oldRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	tsKey2, err := cs.Create("timestamp2", "", data.ED25519Key)
	require.NoError(t, err)
	oldRepo.AddBaseKeys(data.CanonicalTimestampRole, tsKey2)
	tsRole, ok := oldRepo.Root.Signed.Roles[data.CanonicalTimestampRole]
	require.True(t, ok)
	tsRole.Threshold = 2

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	require.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{root, targets, snapshot}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.Contains(t, err.Error(), "timestamp role has invalid threshold")
}

// If any role has a threshold < 1, validation fails
func TestValidateRootInvalidZeroThreshold(t *testing.T) {
	for _, role := range data.BaseRoles {
		oldRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
		require.NoError(t, err)
		tsRole, ok := oldRepo.Root.Signed.Roles[role]
		require.True(t, ok)
		tsRole.Threshold = 0

		r, tg, sn, ts, err := testutils.Sign(oldRepo)
		require.NoError(t, err)
		root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
		require.NoError(t, err)

		store := storage.NewMemStorage()
		updates := []storage.MetaUpdate{root, targets, snapshot}

		serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
		_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid threshold")
	}
}

// ### Role missing negative tests ###
// These tests remove a role from the Root file and
// check for a validation.ErrBadRoot
func TestValidateRootRoleMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "root")

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsRoleMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "targets")

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateSnapshotRoleMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "snapshot")

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

// ### End role missing negative tests ###

// ### Signature missing negative tests ###
func TestValidateRootSigMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "snapshot")

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	r.Signatures = nil

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsSigMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	tg.Signatures = nil

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadTargets{}, err)
}

func TestValidateSnapshotSigMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	sn.Signatures = nil

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End signature missing negative tests ###

// ### Corrupted metadata negative tests ###
func TestValidateRootCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// flip all the bits in the first byte
	root.Data[0] = root.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// flip all the bits in the first byte
	targets.Data[0] = targets.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadTargets{}, err)
}

func TestValidateSnapshotCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// flip all the bits in the first byte
	snapshot.Data[0] = snapshot.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End corrupted metadata negative tests ###

// ### Snapshot size mismatch negative tests ###
func TestValidateRootModifiedSize(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	// add another copy of the signature so the hash is different
	r.Signatures = append(r.Signatures, r.Signatures[0])

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	// flip all the bits in the first byte
	root.Data[0] = root.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsModifiedSize(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	// add another copy of the signature so the hash is different
	tg.Signatures = append(tg.Signatures, tg.Signatures[0])

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End snapshot size mismatch negative tests ###

// ### Snapshot hash mismatch negative tests ###
func TestValidateRootModifiedHash(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	snap, err := data.SnapshotFromSigned(sn)
	require.NoError(t, err)
	snap.Signed.Meta["root"].Hashes["sha256"][0] = snap.Signed.Meta["root"].Hashes["sha256"][0] ^ 0xff

	sn, err = snap.ToSigned()
	require.NoError(t, err)

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadSnapshot{}, err)
}

func TestValidateTargetsModifiedHash(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)

	snap, err := data.SnapshotFromSigned(sn)
	require.NoError(t, err)
	snap.Signed.Meta["targets"].Hashes["sha256"][0] = snap.Signed.Meta["targets"].Hashes["sha256"][0] ^ 0xff

	sn, err = snap.ToSigned()
	require.NoError(t, err)

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	require.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End snapshot hash mismatch negative tests ###

// ### generateSnapshot tests ###
func TestGenerateSnapshotRootNotLoaded(t *testing.T) {
	repo := tuf.NewRepo(nil)
	_, err := generateSnapshot("gun", repo, storage.NewMemStorage())
	require.Error(t, err)
	require.IsType(t, validation.ErrValidation{}, err)
}

func TestGenerateSnapshotNoKey(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	// delete snapshot key in the cryptoservice
	for _, keyID := range cs.ListKeys(data.CanonicalSnapshotRole) {
		require.NoError(t, cs.RemoveKey(keyID))
	}

	_, err = generateSnapshot("gun", repo, store)
	require.Error(t, err)
	require.IsType(t, validation.ErrBadHierarchy{}, err)
}

// ### End generateSnapshot tests ###

// ### Target validation with delegations tests
func TestLoadTargetsFromStore(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	st, err := repo.SignTargets(
		data.CanonicalTargetsRole,
		data.DefaultExpires(data.CanonicalTargetsRole),
	)
	require.NoError(t, err)

	tgs, err := json.Marshal(st)
	require.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgs,
	}
	store.UpdateCurrent("gun", update)

	generated := repo.Targets[data.CanonicalTargetsRole]
	delete(repo.Targets, data.CanonicalTargetsRole)
	_, ok := repo.Targets[data.CanonicalTargetsRole]
	require.False(t, ok)

	err = loadTargetsFromStore("gun", data.CanonicalTargetsRole, repo, store)
	require.NoError(t, err)
	loaded, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.True(t, reflect.DeepEqual(generated.Signatures, loaded.Signatures))
	require.Len(t, loaded.Signed.Targets, 0)
	require.Equal(t, len(generated.Signed.Targets), len(loaded.Signed.Targets))
	require.Len(t, loaded.Signed.Delegations.Roles, 0)
	require.Equal(t, len(generated.Signed.Delegations.Roles), len(loaded.Signed.Delegations.Roles))
	require.Len(t, loaded.Signed.Delegations.Keys, 0)
	require.Equal(t, len(generated.Signed.Delegations.Keys), len(loaded.Signed.Delegations.Keys))
	require.True(t, generated.Signed.Expires.Equal(loaded.Signed.Expires))
	require.Equal(t, generated.Signed.Type, loaded.Signed.Type)
	require.Equal(t, generated.Signed.Version, loaded.Signed.Version)
}

func TestValidateTargetsLoadParent(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	k, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	require.NoError(t, err)

	err = baseRepo.UpdateDelegationKeys("targets/level1", []data.PublicKey{k}, []string{}, 1)
	require.NoError(t, err)
	err = baseRepo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	require.NoError(t, err)

	// no targets file is created for the new delegations, so force one
	baseRepo.InitTargets("targets/level1")

	// we're not going to validate things loaded from storage, so no need
	// to sign the base targets, just Marshal it and set it into storage
	tgtsJSON, err := json.Marshal(baseRepo.Targets["targets"])
	require.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgtsJSON,
	}
	store.UpdateCurrent("gun", update)

	// generate the update object we're doing to use to call loadAndValidateTargets
	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	require.NoError(t, err)
	delJSON, err := json.Marshal(del)
	require.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	roles := map[string]storage.MetaUpdate{"targets/level1": delUpdate}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	updates, err := loadAndValidateTargets("gun", valRepo, roles, store)
	require.NoError(t, err)
	require.Len(t, updates, 1)
	require.Equal(t, "targets/level1", updates[0].Role)
	require.Equal(t, delJSON, updates[0].Data)
}

func TestValidateTargetsParentInUpdate(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	k, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	require.NoError(t, err)

	err = baseRepo.UpdateDelegationKeys("targets/level1", []data.PublicKey{k}, []string{}, 1)
	require.NoError(t, err)
	err = baseRepo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	require.NoError(t, err)

	// no targets file is created for the new delegations, so force one
	baseRepo.InitTargets("targets/level1")

	targets, err := baseRepo.SignTargets("targets", data.DefaultExpires(data.CanonicalTargetsRole))

	tgtsJSON, err := json.Marshal(targets)
	require.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgtsJSON,
	}
	store.UpdateCurrent("gun", update)

	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	require.NoError(t, err)
	delJSON, err := json.Marshal(del)
	require.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	roles := map[string]storage.MetaUpdate{
		"targets/level1": delUpdate,
		"targets":        update,
	}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	// because we sort the roles, the list of returned updates
	// will contain shallower roles first, in this case "targets",
	// and then "targets/level1"
	updates, err := loadAndValidateTargets("gun", valRepo, roles, store)
	require.NoError(t, err)
	require.Len(t, updates, 2)
	require.Equal(t, "targets", updates[0].Role)
	require.Equal(t, tgtsJSON, updates[0].Data)
	require.Equal(t, "targets/level1", updates[1].Role)
	require.Equal(t, delJSON, updates[1].Data)
}

func TestValidateTargetsParentNotFound(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	k, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	require.NoError(t, err)

	err = baseRepo.UpdateDelegationKeys("targets/level1", []data.PublicKey{k}, []string{}, 1)
	require.NoError(t, err)
	err = baseRepo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	require.NoError(t, err)

	// no targets file is created for the new delegations, so force one
	baseRepo.InitTargets("targets/level1")

	// generate the update object we're doing to use to call loadAndValidateTargets
	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	require.NoError(t, err)
	delJSON, err := json.Marshal(del)
	require.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	roles := map[string]storage.MetaUpdate{"targets/level1": delUpdate}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	_, err = loadAndValidateTargets("gun", valRepo, roles, store)
	require.Error(t, err)
	require.IsType(t, storage.ErrNotFound{}, err)
}

func TestValidateTargetsRoleNotInParent(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	store := storage.NewMemStorage()

	level1Key, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	require.NoError(t, err)
	r, err := data.NewRole("targets/level1", 1, []string{level1Key.ID()}, []string{""})

	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Roles = []*data.Role{r}
	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Keys = data.Keys{
		level1Key.ID(): level1Key,
	}

	baseRepo.InitTargets("targets/level1")

	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	require.NoError(t, err)
	delJSON, err := json.Marshal(del)
	require.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	// set back to empty so stored targets doesn't have reference to level1
	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Roles = nil
	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Keys = nil
	targets, err := baseRepo.SignTargets(data.CanonicalTargetsRole, data.DefaultExpires(data.CanonicalTargetsRole))

	tgtsJSON, err := json.Marshal(targets)
	require.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgtsJSON,
	}
	store.UpdateCurrent("gun", update)

	roles := map[string]storage.MetaUpdate{
		"targets/level1":          delUpdate,
		data.CanonicalTargetsRole: update,
	}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	// because we sort the roles, the list of returned updates
	// will contain shallower roles first, in this case "targets",
	// and then "targets/level1"
	updates, err := loadAndValidateTargets("gun", valRepo, roles, store)
	require.NoError(t, err)
	require.Len(t, updates, 1)
	require.Equal(t, data.CanonicalTargetsRole, updates[0].Role)
	require.Equal(t, tgtsJSON, updates[0].Data)
}

// ### End target validation with delegations tests
