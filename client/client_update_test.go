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
	"github.com/gorilla/mux"
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

var metadataDelegations = []string{"targets/a", "targets/a/b", "targets/a/b/c"}

func newServerSwizzler(t *testing.T) (map[string][]byte, *testutils.MetadataSwizzler) {
	serverMeta, cs, err := testutils.NewRepoMetadata("docker.com/notary", metadataDelegations...)
	require.NoError(t, err)

	serverSwizzler := testutils.NewMetadataSwizzler("docker.com/notary", serverMeta, cs)
	require.NoError(t, err)

	return serverMeta, serverSwizzler
}

// bumps the versions of everything in the metadata cache, to force local cache
// to update
func bumpVersions(t *testing.T, s *testutils.MetadataSwizzler) {
	// bump versions of everything on the server, to force everything to update
	for _, r := range s.Roles {
		require.NoError(t, s.OffsetMetadataVersion(r, 1))
	}
	require.NoError(t, s.UpdateSnapshotHashes())
	require.NoError(t, s.UpdateTimestampHash())
}

// create a server that just serves static metadata files from a metaStore
func readOnlyServer(t *testing.T, cache store.MetadataStore, notFoundStatus int) *httptest.Server {
	m := mux.NewRouter()
	m.HandleFunc("/v2/docker.com/notary/_trust/tuf/{role:.*}.json",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			metaBytes, err := cache.GetMeta(vars["role"], maxSize)
			if _, ok := err.(store.ErrMetaNotFound); ok {
				w.WriteHeader(notFoundStatus)
			} else {
				require.NoError(t, err)
				w.Write(metaBytes)
			}
		})
	return httptest.NewServer(m)
}

type unwritableStore struct {
	store.MetadataStore
	roleToNotWrite string
}

func (u *unwritableStore) SetMeta(role string, serverMeta []byte) error {
	if role == u.roleToNotWrite {
		return fmt.Errorf("Non-writable")
	}
	return u.MetadataStore.SetMeta(role, serverMeta)
}

// Update can succeed even if we cannot write any metadata to the repo (assuming
// no data in the repo)
func TestUpdateSucceedsEvenIfCannotWriteNewRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	serverMeta, _, err := testutils.NewRepoMetadata("docker.com/notary", metadataDelegations...)
	require.NoError(t, err)

	ts := readOnlyServer(t, store.NewMemoryStore(serverMeta, nil), http.StatusNotFound)
	defer ts.Close()

	for role := range serverMeta {
		repo := newBlankRepo(t, ts.URL)
		repo.fileStore = &unwritableStore{MetadataStore: repo.fileStore, roleToNotWrite: role}
		_, err := repo.Update(false)

		if role == data.CanonicalRootRole {
			require.Error(t, err) // because checkRoot loads root from cache to check hashes
			continue
		} else {
			require.NoError(t, err)
		}

		for r, expected := range serverMeta {
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
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	serverMeta, serverSwizzler := newServerSwizzler(t)
	ts := readOnlyServer(t, serverSwizzler.MetadataCache, http.StatusNotFound)
	defer ts.Close()

	// download existing metadata
	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	_, err := repo.Update(false)
	require.NoError(t, err)

	origFileStore := repo.fileStore

	for role := range serverMeta {
		for _, forWrite := range []bool{true, false} {
			// bump versions of everything on the server, to force everything to update
			bumpVersions(t, serverSwizzler)

			// update fileStore
			repo.fileStore = &unwritableStore{MetadataStore: origFileStore, roleToNotWrite: role}
			_, err := repo.Update(forWrite)

			if role == data.CanonicalRootRole {
				require.Error(t, err) // because checkRoot loads root from cache to check hashes
				continue
			}
			require.NoError(t, err)

			for r, expected := range serverMeta {
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

func waysToMessUpLocalMetadata(repoSwizzler *testutils.MetadataSwizzler) map[string]messUpMetadata {
	return map[string]messUpMetadata{
		// for instance if the metadata got truncated or otherwise block corrupted
		"invalid JSON": repoSwizzler.SetInvalidJSON,
		// if the metadata was accidentally deleted
		"missing metadata": repoSwizzler.RemoveMetadata,
		// if the signature was invalid - maybe the user tried to modify something manually
		// that they forgot (add a key, or something)
		"signed with right key but wrong hash": repoSwizzler.InvalidateMetadataSignatures,
		// if the user copied the wrong root.json over it by accident or something
		"signed with wrong key": repoSwizzler.SignMetadataWithInvalidKey,
		// self explanatory
		"expired": repoSwizzler.ExpireMetadata,

		// Not trying any of the other repoSwizzler methods, because those involve modifying
		// and re-serializing, and that means a user has the root and other keys and was trying to
		// actively sabotage and break their own local repo (particularly the root.json)
	}
}

// If a repo has corrupt metadata (in that the hash doesn't match the snapshot) or
// missing metadata, an update will replace all of it
func TestUpdateReplacesCorruptOrMissingMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	serverMeta, cs, err := testutils.NewRepoMetadata("docker.com/notary", metadataDelegations...)
	require.NoError(t, err)

	ts := readOnlyServer(t, store.NewMemoryStore(serverMeta, nil), http.StatusNotFound)
	defer ts.Close()

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	_, err = repo.Update(false) // ensure we have all metadata to start with
	require.NoError(t, err)

	// we want to swizzle the local cache, not the server, so create a new one
	repoSwizzler := testutils.NewMetadataSwizzler("docker.com/notary", serverMeta, cs)
	repoSwizzler.MetadataCache = repo.fileStore

	for _, role := range repoSwizzler.Roles {
		for text, messItUp := range waysToMessUpLocalMetadata(repoSwizzler) {
			for _, forWrite := range []bool{true, false} {
				require.NoError(t, messItUp(role), "could not fuzz %s (%s)", role, text)
				_, err := repo.Update(forWrite)
				require.NoError(t, err)
				for r, expected := range serverMeta {
					actual, err := repo.fileStore.GetMeta(r, maxSize)
					require.NoError(t, err, "problem getting repo metadata for %s", role)
					require.True(t, bytes.Equal(expected, actual),
						"%s for %s: expected to recover after update", text, role)
				}
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
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	serverMeta, serverSwizzler := newServerSwizzler(t)
	origMeta := testutils.CopyRepoMetadata(serverMeta)

	ts := readOnlyServer(t, serverSwizzler.MetadataCache, http.StatusNotFound)
	defer ts.Close()

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	_, err := repo.Update(false) // ensure we have all metadata to start with
	require.NoError(t, err)

	// rotate the server's root.json root key so that they no longer match trust anchors
	require.NoError(t, serverSwizzler.ChangeRootKey())
	// bump versions, update snapshot and timestamp too so it will not fail on a hash
	bumpVersions(t, serverSwizzler)

	// we want to swizzle the local cache, not the server, so create a new one
	repoSwizzler := &testutils.MetadataSwizzler{
		MetadataCache: repo.fileStore,
		CryptoService: serverSwizzler.CryptoService,
		Roles:         serverSwizzler.Roles,
	}

	for text, messItUp := range waysToMessUpLocalMetadata(repoSwizzler) {
		for _, forWrite := range []bool{true, false} {
			require.NoError(t, messItUp(data.CanonicalRootRole), "could not fuzz root (%s)", text)
			messedUpMeta, err := repo.fileStore.GetMeta(data.CanonicalRootRole, maxSize)

			if _, ok := err.(store.ErrMetaNotFound); ok { // one of the ways to mess up is to delete metadata

				_, err = repo.Update(forWrite)
				require.Error(t, err) // the new server has a different root key, update fails

			} else {

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

			}

			// revert our original root metadata
			require.NoError(t,
				repo.fileStore.SetMeta(data.CanonicalRootRole, origMeta[data.CanonicalRootRole]))
		}
	}
}

type updateOpts struct {
	notFoundCode     int    // what code to return when the cache doesn't have the metadata
	serverHasNewData bool   // whether the server should have the same or new version than the local cache
	localCache       bool   // whether the repo should have a local cache before updating
	forWrite         bool   // whether the update is for writing or not (force check remote root.json)
	role             string // the role to mess up on the server
}

// If there's no local cache, we go immediately to check the remote server for
// root, and if it doesn't exist, we return ErrRepositoryNotExist. This happens
// with or without a force check (update for write).
func TestUpdateRemoteRootNotExistNoLocalCache(t *testing.T) {
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusNotFound,
		serverHasNewData: false,
		localCache:       false,
		forWrite:         false,
		role:             data.CanonicalRootRole,
	}, ErrRepositoryNotExist{})
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusNotFound,
		serverHasNewData: false,
		localCache:       false,
		forWrite:         true,
		role:             data.CanonicalRootRole,
	}, ErrRepositoryNotExist{})
}

// If there is a local cache, we use the local root as the trust anchor and we
// then an update. If the server has no root.json, we return an ErrRepositoryNotExist.
// If we force check (update for write), then it hits the server first, and
// still returns an ErrRepositoryNotExist.  This is the
// case where the server has the same data as the client, in which case we might
// be able to just used the cached data and not have to download.
func TestUpdateRemoteRootNotExistCanUseLocalCache(t *testing.T) {
	// if for-write is false, then we don't need to check the root.json on bootstrap,
	// and hence we can just use the cached version on update
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusNotFound,
		serverHasNewData: false,
		localCache:       true,
		forWrite:         false,
		role:             data.CanonicalRootRole,
	}, nil)
	// fails because bootstrap requires a check to remote root.json and fails if
	// the check fails
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusNotFound,
		serverHasNewData: false,
		localCache:       true,
		forWrite:         true,
		role:             data.CanonicalRootRole,
	}, ErrRepositoryNotExist{})
}

// If there is a local cache, we use the local root as the trust anchor and we
// then an update. If the server has no root.json, we return an ErrRepositoryNotExist.
// If we force check (update for write), then it hits the server first, and
// still returns an ErrRepositoryNotExist. This is the case where the server
// has new updated data, so we cannot default to cached data.
func TestUpdateRemoteRootNotExistCannotUseLocalCache(t *testing.T) {
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusNotFound,
		serverHasNewData: true,
		localCache:       true,
		forWrite:         false,
		role:             data.CanonicalRootRole,
	}, ErrRepositoryNotExist{})
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusNotFound,
		serverHasNewData: true,
		localCache:       true,
		forWrite:         true,
		role:             data.CanonicalRootRole,
	}, ErrRepositoryNotExist{})
}

// If there's no local cache, we go immediately to check the remote server for
// root, and if it 50X's, we return ErrServerUnavailable. This happens
// with or without a force check (update for write).
func TestUpdateRemoteRoot50XNoLocalCache(t *testing.T) {
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusServiceUnavailable,
		serverHasNewData: false,
		localCache:       false,
		forWrite:         false,
		role:             data.CanonicalRootRole,
	}, store.ErrServerUnavailable{})
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusServiceUnavailable,
		serverHasNewData: false,
		localCache:       false,
		forWrite:         true,
		role:             data.CanonicalRootRole,
	}, store.ErrServerUnavailable{})
}

// If there is a local cache, we use the local root as the trust anchor and we
// then an update. If the server 50X's on root.json, we return an ErrServerUnavailable.
// This happens with or without a force check (update for write).  This is the
// case where the server has the same data as the client, in which case we might
// be able to just used the cached data and not have to download.
func TestUpdateRemoteRoot50XCanUseLocalCache(t *testing.T) {
	// if for-write is false, then we don't need to check the root.json on bootstrap,
	// and hence we can just use the cached version on update
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusServiceUnavailable,
		serverHasNewData: false,
		localCache:       true,
		forWrite:         false,
		role:             data.CanonicalRootRole,
	}, nil)
	// fails because bootstrap requires a check to remote root.json and fails if
	// the check fails
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusServiceUnavailable,
		serverHasNewData: false,
		localCache:       true,
		forWrite:         true,
		role:             data.CanonicalRootRole,
	}, store.ErrServerUnavailable{})
}

// If there is a local cache, we use the local root as the trust anchor and we
// then an update. If the server 50X's on root.json, we return an ErrServerUnavailable.
// This happens with or without a force check (update for write)
func TestUpdateRemoteRoot50XCannotUseLocalCache(t *testing.T) {
	// if for-write is false, then we don't need to check the root.json on bootstrap,
	// and hence we can just use the cached version on update
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusServiceUnavailable,
		serverHasNewData: true,
		localCache:       true,
		forWrite:         false,
		role:             data.CanonicalRootRole,
	}, store.ErrServerUnavailable{})
	// fails because of bootstrap
	testUpdateRemoteError(t, updateOpts{
		notFoundCode:     http.StatusServiceUnavailable,
		serverHasNewData: true,
		localCache:       true,
		forWrite:         true,
		role:             data.CanonicalRootRole,
	}, store.ErrServerUnavailable{})
}

// If there is no local cache, we just update. If the server has a root.json,
// but is missing other data, then we propagate the ErrMetaNotFound.  Skipping
// force check, because that only matters for root.
func TestUpdateNonRootRemoteMissingMetadataNoLocalCache(t *testing.T) {
	for _, role := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		if role == data.CanonicalRootRole {
			continue
		}
		testUpdateRemoteError(t, updateOpts{
			notFoundCode:     http.StatusNotFound,
			serverHasNewData: false,
			localCache:       false,
			forWrite:         false,
			role:             role,
		}, store.ErrMetaNotFound{})
	}
}

// If there is a local cache, we update anyway and see if anything's different
// (assuming remote has a root.json).  If anything's missing, and nothing has
// changed, we don't need to try to download and can just use the local cache.
// Skipping force check, because that only matters for root.
func TestUpdateNonRootRemoteMissingMetadataCanUseLocalCache(t *testing.T) {
	// TODO: fix timestamp
	for _, role := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		if role == data.CanonicalRootRole {
			continue
		}
		testUpdateRemoteError(t, updateOpts{
			notFoundCode:     http.StatusNotFound,
			serverHasNewData: false,
			localCache:       true,
			forWrite:         false,
			role:             role,
		}, nil)
	}
}

// If there is a local cache, we update anyway and see if anything's different
// (assuming remote has a root.json).  If the server has new data, we cannot
// use the local cache so if the server is missing any metadata we cannot update.
// Skipping force check, because that only matters for root.
func TestUpdateNonRootRemoteMissingMetadataCannotUseLocalCache(t *testing.T) {
	for _, role := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		if role == data.CanonicalRootRole {
			continue
		}
		testUpdateRemoteError(t, updateOpts{
			notFoundCode:     http.StatusNotFound,
			serverHasNewData: true,
			localCache:       true,
			forWrite:         false,
			role:             role,
		}, store.ErrMetaNotFound{})
	}
}

// If there is no local cache, we just update. If the server 50X's when getting
// metadata, we propagate ErrServerUnavailable.
func TestUpdateNonRootRemote50XNoLocalCache(t *testing.T) {
	// TODO: fix json syntax error
	for _, role := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		if role == data.CanonicalRootRole {
			continue
		}
		testUpdateRemoteError(t, updateOpts{
			notFoundCode:     http.StatusServiceUnavailable,
			serverHasNewData: false,
			localCache:       false,
			forWrite:         false,
			role:             role,
		}, store.ErrServerUnavailable{})
	}
}

// If there is a local cache, we update anyway and see if anything's different
// (assuming remote has a root.json).  If anything 50X's, and nothing has
// changed, we don't need to try to download and can just use the local cache.
// This happens whether or not we force a remote check (because that's on the
// root)
func TestUpdateNonRootRemote50XCanUseLocalCache(t *testing.T) {
	for _, role := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		if role == data.CanonicalRootRole {
			continue
		}
		testUpdateRemoteError(t, updateOpts{
			notFoundCode:     http.StatusServiceUnavailable,
			serverHasNewData: false,
			localCache:       true,
			forWrite:         false,
			role:             role,
		}, nil)
	}
}

// If there is a local cache, we update anyway and see if anything's different
// (assuming remote has a root.json).  If the server has new data, we cannot
// use the local cache so if the server 50X's on any metadata we cannot update.
// This happens whether or not we force a remote check (because that's on the
// root)
func TestUpdateNonRootRemote50XCannotUseLocalCache(t *testing.T) {
	for _, role := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		if role == data.CanonicalRootRole {
			continue
		}

		var errExpected interface{} = store.ErrServerUnavailable{}
		if role == data.CanonicalTimestampRole {
			// if we can't download the timestamp, we use the cached timestamp.
			// it says that we have all the local data already, so we download
			// nothing.  So the update won't error, it will just fail to update
			// to the latest version.

			// TODO: should we warn that we may not have the latest version?
			errExpected = nil
		}

		testUpdateRemoteError(t, updateOpts{
			notFoundCode:     http.StatusServiceUnavailable,
			serverHasNewData: true,
			localCache:       true,
			forWrite:         false,
			role:             role,
		}, errExpected)
	}
}

func testUpdateRemoteError(t *testing.T, opts updateOpts, errExpected interface{}) {
	_, serverSwizzler := newServerSwizzler(t)
	ts := readOnlyServer(t, serverSwizzler.MetadataCache, opts.notFoundCode)
	defer ts.Close()

	repo := newBlankRepo(t, ts.URL)
	defer os.RemoveAll(repo.baseDir)

	if opts.localCache {
		_, err := repo.Update(false) // acquire local cache
		require.NoError(t, err)
	}

	if opts.serverHasNewData {
		bumpVersions(t, serverSwizzler)
	}

	require.NoError(t, serverSwizzler.MetadataCache.RemoveMeta(opts.role))

	_, err := repo.Update(opts.forWrite)
	if errExpected == nil {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		require.IsType(t, errExpected, err)
		if metaNotFound, ok := err.(store.ErrMetaNotFound); ok {
			require.True(t, ok)
			require.Equal(t, opts.role, metaNotFound.Resource)
		}
	}
}
