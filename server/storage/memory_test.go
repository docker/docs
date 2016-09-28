package storage

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

func TestUpdateCurrent(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})

	k := entryKey("gun", "role")
	gun, ok := s.tufMeta[k]
	v := gun[0]
	require.True(t, ok, "Did not find gun in store")
	require.Equal(t, 1, v.version, "Version mismatch. Expected 1, found %d", v.version)
	require.Equal(t, []byte("test"), v.data, "Data was incorrect")
}

func TestUpdateMany(t *testing.T) {
	s := NewMemStorage()
	require.NoError(t, s.UpdateMany("gun", []MetaUpdate{
		{"role1", 1, []byte("test1")},
		{"role2", 1, []byte("test2")},
	}))

	_, d, err := s.GetCurrent("gun", "role1")
	require.Nil(t, err, "Expected error to be nil")
	require.Equal(t, []byte("test1"), d, "Data was incorrect")

	_, d, err = s.GetCurrent("gun", "role2")
	require.Nil(t, err, "Expected error to be nil")
	require.Equal(t, []byte("test2"), d, "Data was incorrect")

	// updating even one with an equal version fails
	require.IsType(t, &ErrOldVersion{}, s.UpdateMany("gun", []MetaUpdate{
		{"role1", 1, []byte("test1")},
		{"role2", 2, []byte("test2")},
	}))
}

func TestGetCurrent(t *testing.T) {
	s := NewMemStorage()

	_, _, err := s.GetCurrent("gun", "role")
	require.IsType(t, ErrNotFound{}, err, "Expected error to be ErrNotFound")

	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})
	_, d, err := s.GetCurrent("gun", "role")
	require.Nil(t, err, "Expected error to be nil")
	require.Equal(t, []byte("test"), d, "Data was incorrect")
}

func TestDelete(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})
	s.Delete("gun")

	k := entryKey("gun", "role")
	_, ok := s.tufMeta[k]
	require.False(t, ok, "Found gun in store, should have been deleted")
}

func TestGetTimestampKey(t *testing.T) {
	s := NewMemStorage()

	s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))

	c, k, err := s.GetKey("gun", data.CanonicalTimestampRole)
	require.Nil(t, err, "Expected error to be nil")
	require.Equal(t, data.RSAKey, c, "Expected algorithm rsa, received %s", c)
	require.Equal(t, []byte("test"), k, "Key data was wrong")
}

func TestSetKey(t *testing.T) {
	s := NewMemStorage()
	err := s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))
	require.NoError(t, err)

	k := s.keys["gun"][data.CanonicalTimestampRole]
	require.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	require.Equal(t, []byte("test"), k.public, "Public key did not match expected")

}

func TestSetKeyMultipleRoles(t *testing.T) {
	s := NewMemStorage()
	err := s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))
	require.NoError(t, err)

	err = s.SetKey("gun", data.CanonicalSnapshotRole, data.RSAKey, []byte("test"))
	require.NoError(t, err)

	k := s.keys["gun"][data.CanonicalTimestampRole]
	require.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	require.Equal(t, []byte("test"), k.public, "Public key did not match expected")

	k = s.keys["gun"][data.CanonicalSnapshotRole]
	require.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	require.Equal(t, []byte("test"), k.public, "Public key did not match expected")
}

func TestSetKeySameRoleGun(t *testing.T) {
	s := NewMemStorage()
	err := s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))
	require.NoError(t, err)

	// set diff algo and bytes so we can confirm data didn't get replaced
	err = s.SetKey("gun", data.CanonicalTimestampRole, data.ECDSAKey, []byte("test2"))
	require.IsType(t, &ErrKeyExists{}, err, "Expected err to be ErrKeyExists")

	k := s.keys["gun"][data.CanonicalTimestampRole]
	require.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	require.Equal(t, []byte("test"), k.public, "Public key did not match expected")

}

func TestGetChecksumNotFound(t *testing.T) {
	s := NewMemStorage()
	_, _, err := s.GetChecksum("gun", "root", "12345")
	require.Error(t, err)
	require.IsType(t, ErrNotFound{}, err)
}
