package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

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
