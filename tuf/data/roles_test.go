package data

import (
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalRole(t *testing.T) {

	testRoles := map[string]string{
		CanonicalRootRole:      "testRoot",
		CanonicalTargetsRole:   "testTargets",
		CanonicalSnapshotRole:  "testSnapshot",
		CanonicalTimestampRole: "testTimestamp",
		"garbageRole":          "testGarbageRole",
	}

	SetValidRoles(testRoles)

	// make sure roles were set correctly
	assert.Equal(t, "testRoot", ValidRoles[CanonicalRootRole])
	assert.Equal(t, "testTargets", ValidRoles[CanonicalTargetsRole])
	assert.Equal(t, "testSnapshot", ValidRoles[CanonicalSnapshotRole])
	assert.Equal(t, "testTimestamp", ValidRoles[CanonicalTimestampRole])
	// check SetValidRoles doesn't allow non-valid roles in
	assert.Equal(t, "", ValidRoles["garbageRole"])

	// check when looking up CanonicalRole from configured role
	assert.Equal(t, CanonicalRootRole, CanonicalRole("testRoot"))
	assert.Equal(t, CanonicalTargetsRole, CanonicalRole("testTargets"))
	assert.Equal(t, CanonicalSnapshotRole, CanonicalRole("testSnapshot"))
	assert.Equal(t, CanonicalTimestampRole, CanonicalRole("testTimestamp"))
	assert.Equal(t, "", CanonicalRole("testGarbageRole"))

	// check when looking up CanonicalRole with canonical role
	assert.Equal(t, CanonicalRootRole, CanonicalRole(CanonicalRootRole))
	assert.Equal(t, CanonicalTargetsRole, CanonicalRole(CanonicalTargetsRole))
	assert.Equal(t, CanonicalSnapshotRole, CanonicalRole(CanonicalSnapshotRole))
	assert.Equal(t, CanonicalTimestampRole, CanonicalRole(CanonicalTimestampRole))
	assert.Equal(t, "", CanonicalRole("garbageRole"))

	assert.Equal(t, "", CanonicalRole("not found"))

	// reset ValidRoles so other tests aren't messed up
	ValidRoles = map[string]string{
		CanonicalRootRole:      CanonicalRootRole,
		CanonicalTargetsRole:   CanonicalTargetsRole,
		CanonicalSnapshotRole:  CanonicalSnapshotRole,
		CanonicalTimestampRole: CanonicalTimestampRole,
	}
}

func TestMergeStrSlicesExclusive(t *testing.T) {
	orig := []string{"a"}
	new := []string{"b"}

	res := mergeStrSlices(orig, new)
	assert.Len(t, res, 2)
	assert.Equal(t, "a", res[0])
	assert.Equal(t, "b", res[1])
}

func TestMergeStrSlicesOverlap(t *testing.T) {
	orig := []string{"a"}
	new := []string{"a", "b"}

	res := mergeStrSlices(orig, new)
	assert.Len(t, res, 2)
	assert.Equal(t, "a", res[0])
	assert.Equal(t, "b", res[1])
}

func TestMergeStrSlicesEqual(t *testing.T) {
	orig := []string{"a"}
	new := []string{"a"}

	res := mergeStrSlices(orig, new)
	assert.Len(t, res, 1)
	assert.Equal(t, "a", res[0])
}

func TestSubtractStrSlicesExclusive(t *testing.T) {
	orig := []string{"a"}
	new := []string{"b"}

	res := subtractStrSlices(orig, new)
	assert.Len(t, res, 1)
	assert.Equal(t, "a", res[0])
}

func TestSubtractStrSlicesOverlap(t *testing.T) {
	orig := []string{"a", "b"}
	new := []string{"a"}

	res := subtractStrSlices(orig, new)
	assert.Len(t, res, 1)
	assert.Equal(t, "b", res[0])
}

func TestSubtractStrSlicesEqual(t *testing.T) {
	orig := []string{"a"}
	new := []string{"a"}

	res := subtractStrSlices(orig, new)
	assert.Len(t, res, 0)
}

func TestNewRolePathsAndHashPrefixRejection(t *testing.T) {
	_, err := NewRole("targets/level1", 1, []string{"abc"}, nil, nil)
	assert.Error(t, err)
	assert.IsType(t, ErrInvalidRole{}, err)

	_, err = NewRole("targets/level1", 1, []string{"abc"}, []string{""}, []string{""})
	assert.Error(t, err)
	assert.IsType(t, ErrInvalidRole{}, err)
}

func TestAddRemoveKeys(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, []string{""}, nil)
	assert.NoError(t, err)
	role.AddKeys([]string{"abc"})
	assert.Equal(t, []string{"abc"}, role.KeyIDs)
	role.AddKeys([]string{"def"})
	assert.Equal(t, []string{"abc", "def"}, role.KeyIDs)
	role.RemoveKeys([]string{"abc"})
	assert.Equal(t, []string{"def"}, role.KeyIDs)
}

func TestAddRemovePaths(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, []string{"123"}, nil)
	assert.NoError(t, err)
	err = role.AddPaths([]string{"123"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"123"}, role.Paths)
	err = role.AddPaths([]string{"456"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"123", "456"}, role.Paths)
	role.RemovePaths([]string{"123"})
	assert.Equal(t, []string{"456"}, role.Paths)
}

func TestAddRemovePathHashPrefixes(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, nil, []string{"123"})
	assert.NoError(t, err)
	err = role.AddPathHashPrefixes([]string{"123"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"123"}, role.PathHashPrefixes)
	err = role.AddPathHashPrefixes([]string{"456"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"123", "456"}, role.PathHashPrefixes)
	role.RemovePathHashPrefixes([]string{"123"})
	assert.Equal(t, []string{"456"}, role.PathHashPrefixes)
}

func TestAddPathConflict(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, nil, []string{"123"})
	assert.NoError(t, err)
	err = role.AddPaths([]string{"123"})
	assert.Error(t, err)
}

func TestAddPathHashPrefixesConflict(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, []string{"123"}, nil)
	assert.NoError(t, err)
	err = role.AddPathHashPrefixes([]string{"123"})
	assert.Error(t, err)
}

func TestAddPathNil(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, nil, []string{"123"})
	assert.NoError(t, err)
	err = role.AddPaths(nil)
	assert.NoError(t, err)
}

func TestAddPathHashPrefixesNil(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, []string{"123"}, nil)
	assert.NoError(t, err)
	err = role.AddPathHashPrefixes(nil)
	assert.NoError(t, err)
}

func TestErrNoSuchRole(t *testing.T) {
	var err error = ErrNoSuchRole{Role: "test"}
	assert.True(t, strings.HasSuffix(err.Error(), "test"))
}

func TestErrInvalidRole(t *testing.T) {
	var err error = ErrInvalidRole{Role: "test"}
	assert.False(t, strings.Contains(err.Error(), "Reason"))
}

func TestIsDelegation(t *testing.T) {
	assert.True(t, IsDelegation(path.Join(CanonicalTargetsRole, "level1")))
	assert.True(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "level1", "level2", "level3")))
	assert.True(t, IsDelegation(path.Join(CanonicalTargetsRole, "under_score")))
	assert.True(t, IsDelegation(path.Join(CanonicalTargetsRole, "hyphen-hyphen")))
	assert.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, strings.Repeat("x", 255-len(CanonicalTargetsRole)))))

	assert.False(t, IsDelegation(""))
	assert.False(t, IsDelegation(CanonicalRootRole))
	assert.False(t, IsDelegation(path.Join(CanonicalRootRole, "level1")))

	assert.False(t, IsDelegation(CanonicalTargetsRole))
	assert.False(t, IsDelegation(CanonicalTargetsRole+"/"))
	assert.False(t, IsDelegation(path.Join(CanonicalTargetsRole, "level1")+"/"))
	assert.False(t, IsDelegation(path.Join(CanonicalTargetsRole, "UpperCase")))

	assert.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "directory")+"/../../traversal"))

	assert.False(t, IsDelegation(CanonicalTargetsRole+"///test/middle/slashes"))

	assert.False(t, IsDelegation(CanonicalTargetsRole+"/./././"))

	assert.False(t, IsDelegation(
		path.Join("  ", CanonicalTargetsRole, "level1")))

	assert.False(t, IsDelegation(
		path.Join("  "+CanonicalTargetsRole, "level1")))

	assert.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "level1"+"  ")))

	assert.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "white   space"+"level2")))

	assert.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, strings.Repeat("x", 256-len(CanonicalTargetsRole)))))
}

func TestValidRoleFunction(t *testing.T) {
	assert.True(t, ValidRole(CanonicalRootRole))
	assert.True(t, ValidRole(CanonicalTimestampRole))
	assert.True(t, ValidRole(CanonicalSnapshotRole))
	assert.True(t, ValidRole(CanonicalTargetsRole))
	assert.True(t, ValidRole(path.Join(CanonicalTargetsRole, "level1")))
	assert.True(t, ValidRole(
		path.Join(CanonicalTargetsRole, "level1", "level2", "level3")))

	assert.False(t, ValidRole(""))
	assert.False(t, ValidRole(CanonicalRootRole+"/"))
	assert.False(t, ValidRole(CanonicalTimestampRole+"/"))
	assert.False(t, ValidRole(CanonicalSnapshotRole+"/"))
	assert.False(t, ValidRole(CanonicalTargetsRole+"/"))

	assert.False(t, ValidRole(path.Join(CanonicalRootRole, "level1")))

	assert.False(t, ValidRole(path.Join("role")))
}
