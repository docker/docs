package data

import (
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeStrSlicesExclusive(t *testing.T) {
	orig := []string{"a"}
	new := []string{"b"}

	res := mergeStrSlices(orig, new)
	require.Len(t, res, 2)
	require.Equal(t, "a", res[0])
	require.Equal(t, "b", res[1])
}

func TestMergeStrSlicesOverlap(t *testing.T) {
	orig := []string{"a"}
	new := []string{"a", "b"}

	res := mergeStrSlices(orig, new)
	require.Len(t, res, 2)
	require.Equal(t, "a", res[0])
	require.Equal(t, "b", res[1])
}

func TestMergeStrSlicesEqual(t *testing.T) {
	orig := []string{"a"}
	new := []string{"a"}

	res := mergeStrSlices(orig, new)
	require.Len(t, res, 1)
	require.Equal(t, "a", res[0])
}

func TestSubtractStrSlicesExclusive(t *testing.T) {
	orig := []string{"a"}
	new := []string{"b"}

	res := subtractStrSlices(orig, new)
	require.Len(t, res, 1)
	require.Equal(t, "a", res[0])
}

func TestSubtractStrSlicesOverlap(t *testing.T) {
	orig := []string{"a", "b"}
	new := []string{"a"}

	res := subtractStrSlices(orig, new)
	require.Len(t, res, 1)
	require.Equal(t, "b", res[0])
}

func TestSubtractStrSlicesEqual(t *testing.T) {
	orig := []string{"a"}
	new := []string{"a"}

	res := subtractStrSlices(orig, new)
	require.Len(t, res, 0)
}

func TestAddRemoveKeys(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, []string{""})
	require.NoError(t, err)
	role.AddKeys([]string{"abc"})
	require.Equal(t, []string{"abc"}, role.KeyIDs)
	role.AddKeys([]string{"def"})
	require.Equal(t, []string{"abc", "def"}, role.KeyIDs)
	role.RemoveKeys([]string{"abc"})
	require.Equal(t, []string{"def"}, role.KeyIDs)
}

func TestAddRemovePaths(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, []string{"123"})
	require.NoError(t, err)
	err = role.AddPaths([]string{"123"})
	require.NoError(t, err)
	require.Equal(t, []string{"123"}, role.Paths)
	err = role.AddPaths([]string{"456"})
	require.NoError(t, err)
	require.Equal(t, []string{"123", "456"}, role.Paths)
	role.RemovePaths([]string{"123"})
	require.Equal(t, []string{"456"}, role.Paths)
}

func TestAddPathNil(t *testing.T) {
	role, err := NewRole("targets", 1, []string{"abc"}, nil)
	require.NoError(t, err)
	err = role.AddPaths(nil)
	require.NoError(t, err)
}

func TestErrNoSuchRole(t *testing.T) {
	var err error = ErrNoSuchRole{Role: "test"}
	require.True(t, strings.HasSuffix(err.Error(), "test"))
}

func TestErrInvalidRole(t *testing.T) {
	var err error = ErrInvalidRole{Role: "test"}
	require.False(t, strings.Contains(err.Error(), "Reason"))
}

func TestIsDelegation(t *testing.T) {
	require.True(t, IsDelegation(path.Join(CanonicalTargetsRole, "level1")))
	require.True(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "level1", "level2", "level3")))
	require.True(t, IsDelegation(path.Join(CanonicalTargetsRole, "under_score")))
	require.True(t, IsDelegation(path.Join(CanonicalTargetsRole, "hyphen-hyphen")))
	require.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, strings.Repeat("x", 255-len(CanonicalTargetsRole)))))

	require.False(t, IsDelegation(""))
	require.False(t, IsDelegation(CanonicalRootRole))
	require.False(t, IsDelegation(path.Join(CanonicalRootRole, "level1")))

	require.False(t, IsDelegation(CanonicalTargetsRole))
	require.False(t, IsDelegation(CanonicalTargetsRole+"/"))
	require.False(t, IsDelegation(path.Join(CanonicalTargetsRole, "level1")+"/"))
	require.False(t, IsDelegation(path.Join(CanonicalTargetsRole, "UpperCase")))

	require.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "directory")+"/../../traversal"))

	require.False(t, IsDelegation(CanonicalTargetsRole+"///test/middle/slashes"))

	require.False(t, IsDelegation(CanonicalTargetsRole+"/./././"))

	require.False(t, IsDelegation(
		path.Join("  ", CanonicalTargetsRole, "level1")))

	require.False(t, IsDelegation(
		path.Join("  "+CanonicalTargetsRole, "level1")))

	require.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "level1"+"  ")))

	require.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, "white   space"+"level2")))

	require.False(t, IsDelegation(
		path.Join(CanonicalTargetsRole, strings.Repeat("x", 256-len(CanonicalTargetsRole)))))
}

func TestValidRoleFunction(t *testing.T) {
	require.True(t, ValidRole(CanonicalRootRole))
	require.True(t, ValidRole(CanonicalTimestampRole))
	require.True(t, ValidRole(CanonicalSnapshotRole))
	require.True(t, ValidRole(CanonicalTargetsRole))
	require.True(t, ValidRole(path.Join(CanonicalTargetsRole, "level1")))
	require.True(t, ValidRole(
		path.Join(CanonicalTargetsRole, "level1", "level2", "level3")))

	require.False(t, ValidRole(""))
	require.False(t, ValidRole(CanonicalRootRole+"/"))
	require.False(t, ValidRole(CanonicalTimestampRole+"/"))
	require.False(t, ValidRole(CanonicalSnapshotRole+"/"))
	require.False(t, ValidRole(CanonicalTargetsRole+"/"))

	require.False(t, ValidRole(path.Join(CanonicalRootRole, "level1")))

	require.False(t, ValidRole(path.Join("role")))
}

func TestBaseRoleEquals(t *testing.T) {
	fakeKeyHello := NewRSAPublicKey([]byte("hello"))
	fakeKeyThere := NewRSAPublicKey([]byte("there"))

	keys := map[string]PublicKey{"hello": fakeKeyHello, "there": fakeKeyThere}
	baseRole := BaseRole{Name: "name", Threshold: 1, Keys: keys}

	require.True(t, BaseRole{}.Equals(BaseRole{}))
	require.True(t, baseRole.Equals(BaseRole{Name: "name", Threshold: 1, Keys: keys}))
	require.False(t, baseRole.Equals(BaseRole{}))
	require.False(t, baseRole.Equals(BaseRole{Name: "notName", Threshold: 1, Keys: keys}))
	require.False(t, baseRole.Equals(BaseRole{Name: "name", Threshold: 2, Keys: keys}))
	require.False(t, baseRole.Equals(BaseRole{Name: "name", Threshold: 1,
		Keys: map[string]PublicKey{"hello": fakeKeyThere, "there": fakeKeyHello}}))
	require.False(t, baseRole.Equals(BaseRole{Name: "name", Threshold: 1,
		Keys: map[string]PublicKey{"hello": fakeKeyHello, "there": fakeKeyHello}}))
	require.False(t, baseRole.Equals(BaseRole{Name: "name", Threshold: 1,
		Keys: map[string]PublicKey{"hello": fakeKeyHello}}))
	require.False(t, baseRole.Equals(BaseRole{Name: "name", Threshold: 1,
		Keys: map[string]PublicKey{"hello": fakeKeyHello, "there": fakeKeyThere, "again": fakeKeyHello}}))
}
