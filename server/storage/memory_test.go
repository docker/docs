package storage

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCurrent(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})

	k := entryKey("gun", "role")
	gun, ok := s.tufMeta[k]
	v := gun[0]
	assert.True(t, ok, "Did not find gun in store")
	assert.Equal(t, 1, v.version, "Version mismatch. Expected 1, found %d", v.version)
	assert.Equal(t, []byte("test"), v.data, "Data was incorrect")
}

func TestGetCurrent(t *testing.T) {
	s := NewMemStorage()

	_, _, err := s.GetCurrent("gun", "role")
	assert.IsType(t, ErrNotFound{}, err, "Expected error to be ErrNotFound")

	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})
	_, d, err := s.GetCurrent("gun", "role")
	assert.Nil(t, err, "Expected error to be nil")
	assert.Equal(t, []byte("test"), d, "Data was incorrect")
}

func TestDelete(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})
	s.Delete("gun")

	k := entryKey("gun", "role")
	_, ok := s.tufMeta[k]
	assert.False(t, ok, "Found gun in store, should have been deleted")
}

func TestGetTimestampKey(t *testing.T) {
	s := NewMemStorage()

	s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))

	c, k, err := s.GetKey("gun", data.CanonicalTimestampRole)
	assert.Nil(t, err, "Expected error to be nil")
	assert.Equal(t, data.RSAKey, c, "Expected algorithm rsa, received %s", c)
	assert.Equal(t, []byte("test"), k, "Key data was wrong")
}

func TestSetKey(t *testing.T) {
	s := NewMemStorage()
	err := s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))
	assert.NoError(t, err)

	k := s.keys["gun"][data.CanonicalTimestampRole]
	assert.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	assert.Equal(t, []byte("test"), k.public, "Public key did not match expected")

}

func TestSetKeyMultipleRoles(t *testing.T) {
	s := NewMemStorage()
	err := s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))
	assert.NoError(t, err)

	err = s.SetKey("gun", data.CanonicalSnapshotRole, data.RSAKey, []byte("test"))
	assert.NoError(t, err)

	k := s.keys["gun"][data.CanonicalTimestampRole]
	assert.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	assert.Equal(t, []byte("test"), k.public, "Public key did not match expected")

	k = s.keys["gun"][data.CanonicalSnapshotRole]
	assert.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	assert.Equal(t, []byte("test"), k.public, "Public key did not match expected")
}

func TestSetKeySameRoleGun(t *testing.T) {
	s := NewMemStorage()
	err := s.SetKey("gun", data.CanonicalTimestampRole, data.RSAKey, []byte("test"))
	assert.NoError(t, err)

	// set diff algo and bytes so we can confirm data didn't get replaced
	err = s.SetKey("gun", data.CanonicalTimestampRole, data.ECDSAKey, []byte("test2"))
	assert.IsType(t, &ErrKeyExists{}, err, "Expected err to be ErrKeyExists")

	k := s.keys["gun"][data.CanonicalTimestampRole]
	assert.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	assert.Equal(t, []byte("test"), k.public, "Public key did not match expected")

}

func TestGetChecksumNotFound(t *testing.T) {
	s := NewMemStorage()
	_, _, err := s.GetChecksum("gun", "root", "12345")
	assert.Error(t, err)
	assert.IsType(t, ErrNotFound{}, err)
}
