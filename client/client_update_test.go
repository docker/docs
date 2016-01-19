package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/store"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/require"
)

func newBlankRepo(t *testing.T, url string) *NotaryRepository {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo, err := NewNotaryRepository(tempBaseDir, "docker.com/notary", url,
		http.DefaultTransport, passphrase.ConstantRetriever("pass"))
	require.NoError(t, err)
	return repo
}

// If there's no local cache, we go immediately to check the remote server for
// root, and if it doesn't exist, we return ErrRepositoryNotExist. This happens
// with or without a force check (update for write).
func TestUpdateNotExistNoLocalCache(t *testing.T) {
	testUpdateNotExistNoLocalCache(t, false)
	testUpdateNotExistNoLocalCache(t, true)
}

func testUpdateNotExistNoLocalCache(t *testing.T, forWrite bool) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	// there is no metadata at all - this is a fresh repo, and the server isn't
	// aware of the root.
	_, err := repo.Update(forWrite)
	require.Error(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
	require.IsType(t, store.ErrMetaNotFound{}, err)
	metaNotFound, ok := err.(store.ErrMetaNotFound)
	require.True(t, ok)
	require.Equal(t, data.CanonicalTimestampRole, metaNotFound.Resource)
}

// create a server that just serves static metadata files from a metaStore
func readOnlyServer(t *testing.T, cache store.MetadataStore, roles []string) *httptest.Server {
	mux := http.NewServeMux()

	// serve all the metadata files from the metadata cache
	for _, roleName := range roles {
		localRoleName := roleName
		mux.HandleFunc(
			fmt.Sprintf("/v2/docker.com/notary/_trust/tuf/%s.json", localRoleName),
			func(w http.ResponseWriter, r *http.Request) {
				metaBytes, err := cache.GetMeta(localRoleName, maxSize)
				require.NoError(t, err)
				w.Write(metaBytes)
			})

	}

	return httptest.NewServer(mux)
}

type unwritableStore struct {
	store.MetadataStore
	roleToNotWrite string
}

func (u *unwritableStore) SetMeta(role string, meta []byte) error {
	if role == u.roleToNotWrite {
		return fmt.Errorf("Non-writable")
	}
	return u.MetadataStore.SetMeta(role, meta)
}

// Update can succeed even if we cannot write any metadata to the repo (assuming
// no data in the repo)
func TestUpdateSucceedsEvenIfCannotWriteNewRepo(t *testing.T) {
	s, err := testutils.NewMetadataSwizzler("docker.com/notary")
	require.NoError(t, err)

	ts := readOnlyServer(t, s.MetadataCache, s.Roles)
	defer ts.Close()

	for _, role := range s.Roles {
		repo := newBlankRepo(t, ts.URL)
		repo.fileStore = &unwritableStore{MetadataStore: repo.fileStore, roleToNotWrite: role}
		_, err := repo.Update(false)

		if role == data.CanonicalRootRole {
			require.Error(t, err) // because checkRoot loads root from cache to check hashes
			continue
		} else {
			require.NoError(t, err)
		}

		for _, r := range s.Roles {
			expected, err := s.MetadataCache.GetMeta(r, maxSize)
			require.NoError(t, err, "problem getting expected metadata for %s", r)
			actual, err := repo.fileStore.GetMeta(r, maxSize)
			if r == role {
				require.Error(t, err)
				require.IsType(t, store.ErrMetaNotFound{}, err,
					"expected no data because unable to write for %s", role)
			} else {
				require.NoError(t, err, "problem getting repo metadata for %s", r)
				require.True(t, bytes.Equal(expected, actual),
					"%s: expected to update since only %s was unwritable", r, role)
			}
		}

		os.RemoveAll(repo.baseDir)
	}
}

// Update can succeed even if we cannot write any metadata to the repo (assuming
// existing data in the repo)
func TestUpdateSucceedsEvenIfCannotWriteExistingRepo(t *testing.T) {
	s, err := testutils.NewMetadataSwizzler("docker.com/notary")
	require.NoError(t, err)

	ts := readOnlyServer(t, s.MetadataCache, s.Roles)
	defer ts.Close()

	// download existing metadata
	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)
	_, err = repo.Update(false)
	require.NoError(t, err)

	origFileStore := repo.fileStore
	for _, role := range s.Roles {
		for _, forWrite := range []bool{true, false} {
			// bump versions of everything on the server, to force everything to update
			for _, r := range s.Roles {
				require.NoError(t, s.OffsetMetadataVersion(r, 1))
			}
			require.NoError(t, s.UpdateSnapshotHashes())
			require.NoError(t, s.UpdateTimestampHash())

			// update fileStore
			repo.fileStore = &unwritableStore{MetadataStore: origFileStore, roleToNotWrite: role}
			_, err := repo.Update(forWrite)

			if role == data.CanonicalRootRole {
				require.Error(t, err) // because checkRoot loads root from cache to check hashes
				continue
			}
			require.NoError(t, err)

			for _, r := range s.Roles {
				expected, err := s.MetadataCache.GetMeta(r, maxSize)
				require.NoError(t, err, "problem getting expected metadata for %s", r)
				actual, err := repo.fileStore.GetMeta(r, maxSize)
				require.NoError(t, err, "problem getting repo metadata for %s", r)
				if role == r {
					require.False(t, bytes.Equal(expected, actual),
						"%s: expected to not update because %s was unwritable", r, role)
				} else {
					require.True(t, bytes.Equal(expected, actual),
						"%s: expected to update since only %s was unwritable", r, role)
				}
			}
		}
	}
}

type messUpMetadata func(role string) error

// If a repo has corrupt metadata (in that the hash doesn't match the snapshot) or
// missing metadata, an update will replace all of it
func TestUpdateReplacesCorruptOrMissingMetadata(t *testing.T) {
	s, err := testutils.NewMetadataSwizzler("docker.com/notary")
	require.NoError(t, err)

	ts := readOnlyServer(t, s.MetadataCache, s.Roles)
	defer ts.Close()

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	_, err = repo.Update(false) // ensure we have all metadata to start with
	require.NoError(t, err)

	// we want to swizzle the local cache, not the server, so create a new one
	swizzler := testutils.MetadataSwizzler{
		MetadataCache: repo.fileStore,
		CryptoService: repo.CryptoService,
		Roles:         s.Roles,
	}

	waysToMessUp := map[string]messUpMetadata{
		"invalid JSON":     swizzler.SetInvalidJSON,
		"missing metadata": swizzler.DeleteMetadata,
	}
	for _, role := range s.Roles {
		for text, messItUp := range waysToMessUp {
			for _, forWrite := range []bool{true, false} {
				require.NoError(t, messItUp(role), "could not fuzz %s (%s)", role, text)
				_, err := repo.Update(forWrite)
				require.NoError(t, err)
				for _, role := range s.Roles {
					expected, err := s.MetadataCache.GetMeta(role, maxSize)
					require.NoError(t, err, "problem getting expected metadata for %s", role)
					actual, err := repo.fileStore.GetMeta(role, maxSize)
					require.NoError(t, err, "problem getting repo metadata for %s", role)
					require.True(t, bytes.Equal(expected, actual),
						"%s for %s: expected to recover after update", text, role)
				}
			}
		}
	}
}

// If a repo has an invalid root (signed by wrong key, expired, invalid version, etc.),
// the repo will just get the new root from the server, whether or not the update
// is for writing (forced update)
func TestUpdateWhenLocalRootRecoverablyCorrupt(t *testing.T) {
	s, err := testutils.NewMetadataSwizzler("docker.com/notary")
	require.NoError(t, err)

	ts := readOnlyServer(t, s.MetadataCache, s.Roles)
	defer ts.Close()

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	_, err = repo.Update(false) // ensure we have all metadata to start with
	require.NoError(t, err)

	// we want to swizzle the local cache, not the server, so create a new one
	swizzler := testutils.MetadataSwizzler{
		MetadataCache: repo.fileStore,
		CryptoService: s.CryptoService,
		Roles:         s.Roles,
	}

	waysToMessUp := map[string]messUpMetadata{
		// TODO: If invalid threshold fails because it's an invalid role, then this
		// should also fail because it's an invalid role (if the metadata type is wrong)
		"wrong metadata type":                  swizzler.SetInvalidMetadataType,
		"signed with right key but wrong hash": swizzler.InvalidateMetadataSignatures,
		"signed with wrong key":                swizzler.SignMetadataWithInvalidKey,
		"expired":                              swizzler.ExpireMetadata,

		"negative version": func(r string) error { return swizzler.OffsetMetadataVersion(r, -100) },
		// // TODO: This fails at bootstrapClient, when we do SetRoot - then the only way to recover
		// // from this type of local data corruption/change is to probably wipe out the local cache.
		// // Is that ok?  It seems like this is a client corruption where we can just redownload
		// "invalid threshold": func(r string) error { return swizzler.SetThreshold(r, 0) },
		// // TODO: this makes more sense, but the on-disk root itself does not have sufficient sigs.
		// // Is it ok to use it as the trust anchor?
		// "insufficient signatures": func(r string) error { return swizzler.SetThreshold(r, 5) },
	}

	role := data.CanonicalRootRole
	for text, messItUp := range waysToMessUp {
		for _, forWrite := range []bool{true, false} {
			repo.baseURL = ts.URL
			require.NoError(t, messItUp(role), "could not fuzz %s (%s)", role, text)
			_, err := repo.Update(forWrite)
			require.NoError(t, err, "unable to update after locally fuzzing: %s", text)

			for _, role := range s.Roles {
				expected, err := s.MetadataCache.GetMeta(role, maxSize)
				require.NoError(t, err, "problem getting expected metadata for %s", role)
				actual, err := repo.fileStore.GetMeta(role, maxSize)
				require.NoError(t, err, "problem getting repo metadata for %s", role)
				require.True(t, bytes.Equal(expected, actual),
					"%s for %s: expected to recover after update", text, role)
			}
		}
	}
}

// If a repo has an invalid root (signed by wrong key, expired, invalid version,
// invalid number of signatures, etc.), the repo will just get the new root from
// the server, whether or not the update is for writing (forced update), but
// it will refuse to update if the root key has changed and the new root is
// not signed by the old and new key
func TestUpdateFailsIfServerRootKeyChangedWithoutMultiSign(t *testing.T) {
	serverSwizzler, err := testutils.NewMetadataSwizzler("docker.com/notary")
	require.NoError(t, err)

	ts := readOnlyServer(t, serverSwizzler.MetadataCache, serverSwizzler.Roles)
	defer ts.Close()

	origMeta := make(map[string][]byte)
	for _, role := range serverSwizzler.Roles {
		origMetadata, err := serverSwizzler.MetadataCache.GetMeta(role, maxSize)
		require.NoError(t, err)
		origMeta[role] = origMetadata
	}

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	_, err = repo.Update(false) // ensure we have all metadata to start with
	require.NoError(t, err)
	ts.Close()

	// rotate the server's root.json root key so that they no longer match trust anchors
	require.NoError(t, serverSwizzler.ChangeRootKey())
	// bump versions, update snapshot and timestamp too so it will not fail on a hash
	require.NoError(t, serverSwizzler.OffsetMetadataVersion(data.CanonicalRootRole, 1))
	require.NoError(t, serverSwizzler.OffsetMetadataVersion(data.CanonicalSnapshotRole, 1))
	require.NoError(t, serverSwizzler.OffsetMetadataVersion(data.CanonicalTimestampRole, 1))
	require.NoError(t, serverSwizzler.UpdateSnapshotHashes(data.CanonicalRootRole))
	require.NoError(t, serverSwizzler.UpdateTimestampHash())

	// we want to swizzle the local cache, not the server, so create a new one
	swizzler := testutils.MetadataSwizzler{
		MetadataCache: repo.fileStore,
		CryptoService: serverSwizzler.CryptoService,
		Roles:         serverSwizzler.Roles,
	}

	waysToMessUp := map[string]messUpMetadata{
		"wrong metadata type":                  swizzler.SetInvalidMetadataType,
		"signed with right key but wrong hash": swizzler.InvalidateMetadataSignatures,
		"signed with wrong key":                swizzler.SignMetadataWithInvalidKey,
		"expired":                              swizzler.ExpireMetadata,

		"negative version":        func(r string) error { return swizzler.OffsetMetadataVersion(r, -100) },
		"invalid threshold":       func(r string) error { return swizzler.SetThreshold(r, 0) },
		"insufficient signatures": func(r string) error { return swizzler.SetThreshold(r, 5) },
	}

	role := data.CanonicalRootRole
	for text, messItUp := range waysToMessUp {
		for _, forWrite := range []bool{true, false} {
			require.NoError(t, messItUp(role), "could not fuzz %s (%s)", role, text)
			messedUpMeta, err := repo.fileStore.GetMeta(data.CanonicalRootRole, maxSize)
			require.NoError(t, err)

			_, err = repo.Update(forWrite)
			require.Error(t, err) // the new server has a different root, update fails

			// we can't test that all the metadata is the same, because we probably would
			// have downloaded a new timestamp and maybe snapshot.  But the root should be the
			// same because it has failed to update.
			for role, expected := range origMeta {
				if role != data.CanonicalTimestampRole && role != data.CanonicalSnapshotRole {
					actual, err := repo.fileStore.GetMeta(role, maxSize)
					require.NoError(t, err, "problem getting repo metadata for %s", role)

					if role == data.CanonicalRootRole {
						expected = messedUpMeta
					}
					require.True(t, bytes.Equal(expected, actual),
						"%s for %s: expected to not have updated", text, role)
				}
			}
			// revert our original root metadata
			require.NoError(t,
				repo.fileStore.SetMeta(data.CanonicalRootRole, origMeta[data.CanonicalRootRole]))
		}
	}
}
