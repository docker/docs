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

	_, err := s.GetCurrent("gun", "role")
	assert.IsType(t, &ErrNotFound{}, err, "Expected error to be ErrNotFound")

	s.UpdateCurrent("gun", MetaUpdate{"role", 1, []byte("test")})
	d, err := s.GetCurrent("gun", "role")
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

	//_, _, err := s.GetTimestampKey("gun")
	//assert.IsType(t, &ErrNoKey{}, err, "Expected err to be ErrNoKey")

	s.SetTimestampKey("gun", data.RSAKey, []byte("test"))

	c, k, err := s.GetTimestampKey("gun")
	assert.Nil(t, err, "Expected error to be nil")
	assert.Equal(t, data.RSAKey, c, "Expected algorithm rsa, received %s", c)
	assert.Equal(t, []byte("test"), k, "Key data was wrong")
}

func TestSetTimestampKey(t *testing.T) {
	s := NewMemStorage()
	s.SetTimestampKey("gun", data.RSAKey, []byte("test"))

	err := s.SetTimestampKey("gun", data.RSAKey, []byte("test2"))
	assert.IsType(t, &ErrTimestampKeyExists{}, err, "Expected err to be ErrTimestampKeyExists")

	k := s.tsKeys["gun"]
	assert.Equal(t, data.RSAKey, k.algorithm, "Expected algorithm to be rsa, received %s", k.algorithm)
	assert.Equal(t, []byte("test"), k.public, "Public key did not match expected")

}
