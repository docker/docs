package tuf_test

// package tuf_test to avoid an import cycle since we are using testutils.EmptyRepo

import (
	"testing"

	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/require"
)

var _cachedMeta map[string][]byte

// we just want sample metadata for a role - so we can build cached metadata
// and use it once.
func getSampleMeta(t *testing.T) (map[string][]byte, string) {
	gun := "docker.com/notary"
	delgNames := []string{"targets/a", "targets/a/b", "targets/a/b/force_parent_metadata"}
	if _cachedMeta == nil {
		meta, _, err := testutils.NewRepoMetadata(gun, delgNames...)
		require.NoError(t, err)

		_cachedMeta = meta
	}
	return _cachedMeta, gun
}

// We load only if the rolename is a valid rolename - even if the metadata we provided is valid
func TestBuilderLoadsValidRolesOnly(t *testing.T) {
	meta, gun := getSampleMeta(t)
	builder := tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})
	err := builder.Load("NotRoot", meta[data.CanonicalRootRole], 0, false)
	require.Error(t, err)
	require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
	require.Contains(t, err.Error(), "is an invalid role")
}

func TestBuilderOnlyAcceptsRootFirstWhenLoading(t *testing.T) {
	meta, gun := getSampleMeta(t)
	builder := tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})

	for roleName, content := range meta {
		if roleName != data.CanonicalRootRole {
			err := builder.Load(roleName, content, 0, true)
			require.Error(t, err)
			require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
			require.Contains(t, err.Error(), "root must be loaded first")
			require.False(t, builder.IsLoaded(roleName))
			require.Equal(t, 0, builder.GetLoadedVersion(roleName))
		}
	}

	// we can load the root
	require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
	require.True(t, builder.IsLoaded(data.CanonicalRootRole))
}

func TestBuilderOnlyAcceptsDelegationsAfterParent(t *testing.T) {
	meta, gun := getSampleMeta(t)
	builder := tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})

	// load the root
	require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))

	// delegations can't be loaded without target
	for _, delgName := range []string{"targets/a", "targets/a/b"} {
		err := builder.Load(delgName, meta[delgName], 0, false)
		require.Error(t, err)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "targets must be loaded first")
		require.False(t, builder.IsLoaded(delgName))
		require.Equal(t, 0, builder.GetLoadedVersion(delgName))
	}

	// load the targets
	require.NoError(t, builder.Load(data.CanonicalTargetsRole, meta[data.CanonicalTargetsRole], 0, false))

	// targets/a/b can't be loaded because targets/a isn't loaded
	err := builder.Load("targets/a/b", meta["targets/a/b"], 0, false)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	// targets/a can be loaded now though because targets is loaded
	require.NoError(t, builder.Load("targets/a", meta["targets/a"], 0, false))

	// and now targets/a/b can be loaded because targets/a is loaded
	require.NoError(t, builder.Load("targets/a/b", meta["targets/a/b"], 0, false))
}

func TestBuilderAcceptRoleOnce(t *testing.T) {
	meta, gun := getSampleMeta(t)
	builder := tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})

	for _, roleName := range append(data.BaseRoles, "targets/a", "targets/a/b") {
		// first time loading is ok
		require.NoError(t, builder.Load(roleName, meta[roleName], 0, false))
		require.True(t, builder.IsLoaded(roleName))
		require.Equal(t, 1, builder.GetLoadedVersion(roleName))

		// second time loading is not
		err := builder.Load(roleName, meta[roleName], 0, false)
		require.Error(t, err)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "has already been loaded")

		// still loaded
		require.True(t, builder.IsLoaded(roleName))
	}
}

func TestBuilderStopsAcceptingOrProducingDataOnceDone(t *testing.T) {
	meta, gun := getSampleMeta(t)
	builder := tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})

	for _, roleName := range data.BaseRoles {
		require.NoError(t, builder.Load(roleName, meta[roleName], 0, false))
		require.True(t, builder.IsLoaded(roleName))
	}

	_, err := builder.Finish()
	require.NoError(t, err)

	err = builder.Load("targets/a", meta["targets/a"], 0, false)
	require.Error(t, err)
	require.Equal(t, tuf.ErrBuildDone, err)

	// a new bootstrapped builder can also not have any more input output
	bootstrapped := builder.BootstrapNewBuilder()

	err = bootstrapped.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false)
	require.Error(t, err)
	require.Equal(t, tuf.ErrBuildDone, err)

	for _, b := range []tuf.RepoBuilder{builder, bootstrapped} {
		_, err = b.Finish()
		require.Error(t, err)
		require.Equal(t, tuf.ErrBuildDone, err)

		_, _, err = b.GenerateSnapshot(nil)
		require.Error(t, err)
		require.Equal(t, tuf.ErrBuildDone, err)

		_, _, err = b.GenerateTimestamp(nil)
		require.Error(t, err)
		require.Equal(t, tuf.ErrBuildDone, err)

		for roleName := range meta {
			// a finished builder thinks nothing is loaded
			require.False(t, b.IsLoaded(roleName))
			// checksums are all empty, versions are all zero
			require.Equal(t, 0, b.GetLoadedVersion(roleName))
			require.Equal(t, tuf.ConsistentInfo{RoleName: roleName}, b.GetConsistentInfo(roleName))
		}

	}
}

// Test the cases in which GenerateSnapshot fails
func TestGenerateSnapshotInvalidOperations(t *testing.T) {
	gun := "docker.com/notary"
	repo, cs, err := testutils.EmptyRepo(gun)
	require.NoError(t, err)

	// make snapshot have 2 keys and a threshold of 2
	snapKeys := make([]data.PublicKey, 2)
	for i := 0; i < 2; i++ {
		snapKeys[i], err = cs.Create(data.CanonicalSnapshotRole, gun, data.ECDSAKey)
		require.NoError(t, err)
	}

	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalSnapshotRole, snapKeys...))
	repo.Root.Signed.Roles[data.CanonicalSnapshotRole].Threshold = 2

	meta, err := testutils.SignAndSerialize(repo)
	require.NoError(t, err)

	for _, prevSnapshot := range []*data.SignedSnapshot{nil, repo.Snapshot} {
		// copy keys, since we expect one of these generation attempts to succeed and we do
		// some key deletion tests later
		newCS := testutils.CopyKeys(t, cs, data.CanonicalSnapshotRole)

		// --- we can't generate a snapshot if the root isn't loaded
		builder := tuf.NewRepoBuilder(gun, newCS, trustpinning.TrustPinConfig{})
		_, _, err := builder.GenerateSnapshot(prevSnapshot)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "root must be loaded first")
		require.False(t, builder.IsLoaded(data.CanonicalSnapshotRole))

		// --- we can't generate a snapshot if the targets isn't loaded and we have no previous snapshot,
		// --- but if we have a previous snapshot with a valid targets, we're good even if no snapshot
		// --- is loaded
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		_, _, err = builder.GenerateSnapshot(prevSnapshot)
		if prevSnapshot == nil {
			require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
			require.Contains(t, err.Error(), "targets must be loaded first")
			require.False(t, builder.IsLoaded(data.CanonicalSnapshotRole))
		} else {
			require.NoError(t, err)
		}

		// --- we can't generate a snapshot if we've loaded the timestamp already
		builder = tuf.NewRepoBuilder(gun, newCS, trustpinning.TrustPinConfig{})
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		if prevSnapshot == nil {
			require.NoError(t, builder.Load(data.CanonicalTargetsRole, meta[data.CanonicalTargetsRole], 0, false))
		}
		require.NoError(t, builder.Load(data.CanonicalTimestampRole, meta[data.CanonicalTimestampRole], 0, false))

		_, _, err = builder.GenerateSnapshot(prevSnapshot)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "cannot generate snapshot if timestamp has already been loaded")
		require.False(t, builder.IsLoaded(data.CanonicalSnapshotRole))

		// --- we cannot generate a snapshot if we've already loaded a snapshot
		builder = tuf.NewRepoBuilder(gun, newCS, trustpinning.TrustPinConfig{})
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		if prevSnapshot == nil {
			require.NoError(t, builder.Load(data.CanonicalTargetsRole, meta[data.CanonicalTargetsRole], 0, false))
		}
		require.NoError(t, builder.Load(data.CanonicalSnapshotRole, meta[data.CanonicalSnapshotRole], 0, false))

		_, _, err = builder.GenerateSnapshot(prevSnapshot)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "snapshot has already been loaded")

		// --- we cannot generate a snapshot if we can't satisfy the role threshold
		for i := 0; i < len(snapKeys); i++ {
			require.NoError(t, newCS.RemoveKey(snapKeys[i].ID()))
			builder = tuf.NewRepoBuilder(gun, newCS, trustpinning.TrustPinConfig{})
			require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
			if prevSnapshot == nil {
				require.NoError(t, builder.Load(data.CanonicalTargetsRole, meta[data.CanonicalTargetsRole], 0, false))
			}

			_, _, err = builder.GenerateSnapshot(prevSnapshot)
			require.IsType(t, signed.ErrInsufficientSignatures{}, err)
			require.False(t, builder.IsLoaded(data.CanonicalSnapshotRole))
		}

		// --- we cannot generate a snapshot if we don't have a cryptoservice
		builder = tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		if prevSnapshot == nil {
			require.NoError(t, builder.Load(data.CanonicalTargetsRole, meta[data.CanonicalTargetsRole], 0, false))
		}

		_, _, err = builder.GenerateSnapshot(prevSnapshot)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "cannot generate snapshot without a cryptoservice")
		require.False(t, builder.IsLoaded(data.CanonicalSnapshotRole))
	}

	// --- we can't generate a snapshot if we're given an invalid previous snapshot (for instance, an empty one),
	// --- even if we have a targets loaded
	builder := tuf.NewRepoBuilder(gun, cs, trustpinning.TrustPinConfig{})
	require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
	require.NoError(t, builder.Load(data.CanonicalTargetsRole, meta[data.CanonicalTargetsRole], 0, false))

	_, _, err = builder.GenerateSnapshot(&data.SignedSnapshot{})
	require.IsType(t, data.ErrInvalidMetadata{}, err)
	require.False(t, builder.IsLoaded(data.CanonicalSnapshotRole))
}

// Test the cases in which GenerateTimestamp fails
func TestGenerateTimestampInvalidOperations(t *testing.T) {
	gun := "docker.com/notary"
	repo, cs, err := testutils.EmptyRepo(gun)
	require.NoError(t, err)

	// make timsetamp have 2 keys and a threshold of 2
	tsKeys := make([]data.PublicKey, 2)
	for i := 0; i < 2; i++ {
		tsKeys[i], err = cs.Create(data.CanonicalTimestampRole, gun, data.ECDSAKey)
		require.NoError(t, err)
	}

	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalTimestampRole, tsKeys...))
	repo.Root.Signed.Roles[data.CanonicalTimestampRole].Threshold = 2

	meta, err := testutils.SignAndSerialize(repo)
	require.NoError(t, err)

	for _, prevTimestamp := range []*data.SignedTimestamp{nil, repo.Timestamp} {
		// --- we can't generate a timestamp if the root isn't loaded
		builder := tuf.NewRepoBuilder(gun, cs, trustpinning.TrustPinConfig{})
		_, _, err := builder.GenerateTimestamp(prevTimestamp)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "root must be loaded first")
		require.False(t, builder.IsLoaded(data.CanonicalTimestampRole))

		// --- we can't generate a timestamp if the snapshot isn't loaded, no matter if we have a previous
		// --- timestamp or not
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		_, _, err = builder.GenerateTimestamp(prevTimestamp)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "snapshot must be loaded first")
		require.False(t, builder.IsLoaded(data.CanonicalTimestampRole))

		// --- we can't generate a timestamp if we've loaded the timestamp already
		builder = tuf.NewRepoBuilder(gun, cs, trustpinning.TrustPinConfig{})
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		require.NoError(t, builder.Load(data.CanonicalSnapshotRole, meta[data.CanonicalSnapshotRole], 0, false))
		require.NoError(t, builder.Load(data.CanonicalTimestampRole, meta[data.CanonicalTimestampRole], 0, false))

		_, _, err = builder.GenerateTimestamp(prevTimestamp)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "timestamp has already been loaded")

		// --- we cannot generate a timestamp if we can't satisfy the role threshold
		for i := 0; i < len(tsKeys); i++ {
			require.NoError(t, cs.RemoveKey(tsKeys[i].ID()))
			builder = tuf.NewRepoBuilder(gun, cs, trustpinning.TrustPinConfig{})
			require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
			require.NoError(t, builder.Load(data.CanonicalSnapshotRole, meta[data.CanonicalSnapshotRole], 0, false))

			_, _, err = builder.GenerateTimestamp(prevTimestamp)
			require.IsType(t, signed.ErrInsufficientSignatures{}, err)
			require.False(t, builder.IsLoaded(data.CanonicalTimestampRole))
		}

		// --- we cannot generate a timestamp if we don't have a cryptoservice
		builder = tuf.NewRepoBuilder(gun, nil, trustpinning.TrustPinConfig{})
		require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
		require.NoError(t, builder.Load(data.CanonicalSnapshotRole, meta[data.CanonicalSnapshotRole], 0, false))

		_, _, err = builder.GenerateTimestamp(prevTimestamp)
		require.IsType(t, tuf.ErrInvalidBuilderInput{}, err)
		require.Contains(t, err.Error(), "cannot generate timestamp without a cryptoservice")
		require.False(t, builder.IsLoaded(data.CanonicalTimestampRole))
	}

	// --- we can't generate a timsetamp if we're given an invalid previous timestamp (for instance, an empty one),
	// --- even if we have a snapshot loaded
	builder := tuf.NewRepoBuilder(gun, cs, trustpinning.TrustPinConfig{})
	require.NoError(t, builder.Load(data.CanonicalRootRole, meta[data.CanonicalRootRole], 0, false))
	require.NoError(t, builder.Load(data.CanonicalSnapshotRole, meta[data.CanonicalSnapshotRole], 0, false))

	_, _, err = builder.GenerateTimestamp(&data.SignedTimestamp{})
	require.IsType(t, data.ErrInvalidMetadata{}, err)
	require.False(t, builder.IsLoaded(data.CanonicalTimestampRole))
}
