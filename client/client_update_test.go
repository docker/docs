package client

import (
	"os"
	"path"
	"testing"
	"time"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/store"
	json "github.com/jfrazelle/go/canonical/json"
	"github.com/stretchr/testify/require"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/store"
	"github.com/stretchr/testify/require"
)

// If there's no local cache, we go immediately to check the remote server for
// root, and if it doesn't exist, we return ErrRepositoryNotExist. This happens
// with or without a force check (update for write).
func TestUpdateNotExistNoLocalCache(t *testing.T) {
	testUpdateNotExistNoLocalCache(t, false)
	testUpdateNotExistNoLocalCache(t, true)
}

func testUpdateNotExistNoLocalCache(t *testing.T, forWrite bool) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, "docker.com/notary", ts.URL,
		http.DefaultTransport, nil)
	require.NoError(t, err)

	// there is no metadata at all - this is a fresh repo, and the server isn't
	// aware of the root.
	_, err = repo.Update(forWrite)
	require.IsType(t, ErrRepositoryNotExist{}, err)
}

// If there is a local cache, we use the local root as the trust anchor and we
// then an update. If the server has no root.json, we return an ErrRepositoryNotExist.
// If we force check (update for write), then it hits the server first, and
// still returns an ErrRepositoryNotExist.
func TestUpdateNotExistWithLocalCache(t *testing.T) {
	testUpdateNotExistWithLocalCache(t, false)
	testUpdateNotExistWithLocalCache(t, true)
}

func testUpdateNotExistWithLocalCache(t *testing.T, forWrite bool) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// the repo has metadata, but the server is unaware of any metadata
	// whatsoever.
	_, err := repo.Update(forWrite)
	require.IsType(t, ErrRepositoryNotExist{}, err)
}

// If there is a local cache, we use the local root as the trust anchor and we
// then an update. If the server has a root.json, but is missing other data,
// then we propagate the ErrMetaNotFound.  Same if we force check
// (update for write); the root exists, but other metadata doesn't.
func TestUpdateWithLocalCacheRemoteMissingMetadata(t *testing.T) {
	testUpdateWithLocalCacheRemoteMissingMetadata(t, false)
	testUpdateWithLocalCacheRemoteMissingMetadata(t, true)
}

func testUpdateWithLocalCacheRemoteMissingMetadata(t *testing.T, forWrite bool) {
	ts, mux, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	rootJSON, err := repo.fileStore.GetMeta(data.CanonicalRootRole, maxSize)
	require.NoError(t, err)

	// the server should know about the root.json, and nothing else
	mux.HandleFunc(
		fmt.Sprintf("/v2/docker.com/notary/_trust/tuf/root.json"),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, string(rootJSON))
		})

	// the first thing the client tries to get is the timestamp - so that
	// will be the failed metadata update.
	_, err = repo.Update(forWrite)
	require.IsType(t, store.ErrMetaNotFound{}, err)
	metaNotFound, ok := err.(store.ErrMetaNotFound)
	require.True(t, ok)
	require.Equal(t, data.CanonicalTimestampRole, metaNotFound.Resource)
}

type messUpMetadata func(t *testing.T, cs signed.CryptoService, ms store.MetadataStore, role string)

// corrupts metadata into something that is no longer valid JSON
func invalidJSONMetadata(t *testing.T, _ signed.CryptoService, ms store.MetadataStore, role string) {
	require.NoError(t, ms.SetMeta(role, []byte("nope")))
}

// corrupts the metadata into something that is valid JSON, but not unmarshalable at all
func allUnmarshallableMetadata(t *testing.T, _ signed.CryptoService, ms store.MetadataStore, role string) {
	metaBytes, err := json.MarshalCanonical(data.Signed{})
	require.NoError(t, err)
	require.NoError(t, ms.SetMeta(role, metaBytes))
}

// messes up the metadata in such a way that the hash is no longer valid
func invalidateMetadataHash(t *testing.T, _ signed.CryptoService, ms store.MetadataStore, role string) {
	b, err := ms.GetMeta(role, maxSize)
	require.NoError(t, err)

	var unmarshalled map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &unmarshalled))

	signed, ok := unmarshalled["signed"].(map[string]interface{})
	require.True(t, ok)
	signed["boogeyman"] = "exists"

	metaBytes, err := json.MarshalCanonical(unmarshalled)
	require.NoError(t, err)

	require.NoError(t, ms.SetMeta(role, metaBytes))
}

// deletes the metadata
func deleteMetadata(t *testing.T, _ signed.CryptoService, ms store.MetadataStore, role string) {
	require.NoError(t, ms.DeleteMeta(role))
}

func serializeMetadata(t *testing.T, s *data.Signed, cs signed.CryptoService, role string) []byte {
	// delete the existing signatures
	s.Signatures = []data.Signature{}

	pubKeys := cs.ListKeys(role)
	require.Len(t, pubKeys, 1, "no keys for %s", role)
	pubKey := cs.GetKey(pubKeys[0])
	require.NotNil(t, pubKey, "unable to get %s key %s", role, pubKeys[0])

	require.NoError(t, signed.Sign(cs, s, pubKey))

	metaBytes, err := json.MarshalCanonical(s)
	require.NoError(t, err)

	return metaBytes
}

// signs the metadata with the wrong key
func invalidateMetadataSig(t *testing.T, _ signed.CryptoService, ms store.MetadataStore, role string) {
	b, err := ms.GetMeta(role, maxSize)
	require.NoError(t, err)

	signedThing := data.Signed{}
	require.NoError(t, json.Unmarshal(b, &signedThing), "error unmarshalling data for %s", role)

	// create an invalid key, but not in the existing CryptoService
	cs := signed.NewEd25519()
	_, err = cs.Create("root", data.ED25519Key)
	require.NoError(t, err)

	metaBytes := serializeMetadata(t, &signedThing, cs, "root")
	require.NoError(t, ms.SetMeta(role, metaBytes))
}

func signedMetaFromStore(t *testing.T, ms store.MetadataStore, role string) data.SignedMeta {
	b, err := ms.GetMeta(role, maxSize)
	require.NoError(t, err)

	signedMeta := data.SignedMeta{}
	require.NoError(t, json.Unmarshal(b, &signedMeta), "error unmarshalling data for %s", role)

	return signedMeta
}

func signedMetaToSigned(t *testing.T, signedMeta data.SignedMeta) data.Signed {
	s, err := json.MarshalCanonical(signedMeta.Signed)
	require.NoError(t, err)
	signed := json.RawMessage{}
	require.NoError(t, signed.UnmarshalJSON(s))

	return data.Signed{Signed: signed}
}

// corrupt the metadata in such a way that it is JSON parsable, and correctly signed, but will not
// unmarshal correctly because it has the wrong type
func corruptSignedMetadata(t *testing.T, cs signed.CryptoService, ms store.MetadataStore, role string) {
	if role != data.CanonicalTimestampRole || len(cs.ListKeys(role)) > 0 {
		signedMeta := signedMetaFromStore(t, ms, role)
		signedMeta.Signed.Type = "nonexistent"
		signedThing := signedMetaToSigned(t, signedMeta)
		metaBytes := serializeMetadata(t, &signedThing, cs, role)
		require.NoError(t, ms.SetMeta(role, metaBytes))
	}
}

// decrements the metadata version, which would make it invalid - don't do anything if we don't
// have the timestamp key
func decrementMetadataVersion(t *testing.T, cs signed.CryptoService, ms store.MetadataStore, role string) {
	if role != data.CanonicalTimestampRole || len(cs.ListKeys(role)) > 0 {
		signedMeta := signedMetaFromStore(t, ms, role)
		signedMeta.Signed.Version--
		signedThing := signedMetaToSigned(t, signedMeta)
		metaBytes := serializeMetadata(t, &signedThing, cs, role)
		require.NoError(t, ms.SetMeta(role, metaBytes))
	}
}

// expire the metadata, which would make it invalid - don't do anything if we don't have the
// timestamp key
func expireMetadata(t *testing.T, cs signed.CryptoService, ms store.MetadataStore, role string) {
	if role != data.CanonicalTimestampRole || len(cs.ListKeys(role)) > 0 {
		signedMeta := signedMetaFromStore(t, ms, role)
		signedMeta.Signed.Expires = time.Now().AddDate(-1, -1, -1)
		signedThing := signedMetaToSigned(t, signedMeta)
		metaBytes := serializeMetadata(t, &signedThing, cs, role)
		require.NoError(t, ms.SetMeta(role, metaBytes))
	}
}

// increments a threshold for a metadata role - invalidates the metadata for which the threshold
// is increased, since there is only 1 signature for each
func incrementThreshold(t *testing.T, cs signed.CryptoService, ms store.MetadataStore, role string) {
	roleSpecifyingThreshold := data.CanonicalRootRole
	if data.IsDelegation(role) {
		roleSpecifyingThreshold = path.Dir(role)
	}

	b, err := ms.GetMeta(roleSpecifyingThreshold, maxSize)
	require.NoError(t, err)

	signedThing := &data.Signed{}
	require.NoError(t, json.Unmarshal(b, signedThing), "error unmarshalling data for %s",
		roleSpecifyingThreshold)

	if roleSpecifyingThreshold == data.CanonicalRootRole {
		signedRoot, err := data.RootFromSigned(signedThing)
		require.NoError(t, err)
		signedRoot.Signed.Roles[role].Threshold++
		signedThing, err = signedRoot.ToSigned()
		require.NoError(t, err)
	} else {
		signedTargets, err := data.TargetsFromSigned(signedThing)
		require.NoError(t, err)
		for _, roleObject := range signedTargets.Signed.Delegations.Roles {
			if roleObject.Name == role {
				roleObject.Threshold++
				break
			}
		}
		signedThing, err = signedTargets.ToSigned()
		require.NoError(t, err)
	}

	metaBytes := serializeMetadata(t, signedThing, cs, roleSpecifyingThreshold)
	require.NoError(t, ms.SetMeta(roleSpecifyingThreshold, metaBytes))
}

// If a repo has corrupt metadata, an update will replace all the metadata
func TestUpdateReplacesCorruptOrMissingMetadata(t *testing.T) {
	// create repo with 2 level delegations
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	delegatedRoles := []string{"targets/a", "targets/a/b"}
	for _, delgName := range delegatedRoles {
		delgKey, err := repo.CryptoService.Create(delgName, data.ECDSAKey)
		require.NoError(t, err, "error creating delegation key")

		require.NoError(t,
			repo.AddDelegation(delgName, 1, []data.PublicKey{delgKey}, []string{""}),
			"error creating delegation")
	}
	// add a target so the second level delegation is created
	addTarget(t, repo, "first", "../fixtures/root-ca.crt", "targets/a/b")
	require.NoError(t, repo.Publish())
	_, err := repo.Update() // ensure we have all metadata to start with
	require.NoError(t, err)

	// corrupt any number of roles - an update should fix all of them
	roles := []string{
		data.CanonicalTimestampRole,
		data.CanonicalSnapshotRole,
		"targets/a",
		"targets/a/b",
		data.CanonicalTargetsRole,
		data.CanonicalRootRole,
	}

	// store original metadata
	origMeta := make(map[string][]byte)
	for _, role := range roles {
		b, err := repo.fileStore.GetMeta(role, maxSize)
		require.NoError(t, err)
		require.NotNil(t, b)
		origMeta[role] = b
	}

	// mess up metadata in different ways, update the repo, and assert that the metadata is fixed.
	waysToMessUp := map[string]messUpMetadata{
		"corrupted/invalid JSON":       invalidJSONMetadata,
		"metadata has invalid hash":    invalidateMetadataHash,
		"missing metadata":             deleteMetadata,
		"metadata signed by wrong key": invalidateMetadataSig,
		"expired metadata":             expireMetadata,
		"insufficient signatures":      incrementThreshold,
		// decremented version just tests that updates do not need to increment
		// by 1, only increment at all
		"version much lower": decrementMetadataVersion,
	}
	for i := range roles {
		for text, messItUp := range waysToMessUp {
			for _, role := range roles[:i+1] {
				messItUp(t, repo.CryptoService, repo.fileStore, role)
			}
			_, err := repo.Update()
			require.NoError(t, err)
			for role, origBytes := range origMeta {
				b, err := repo.fileStore.GetMeta(role, maxSize)
				require.NoError(t, err, "problem getting metadata for %s", role)
				require.Equal(t, origBytes, b, "%s for %s expected to recover after update", text, role)
			}
		}
	}
}
