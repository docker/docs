package keydbstore

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var retriever = func(string, string, bool, int) (string, bool, error) {
	return "passphrase_1", false, nil
}

var anotherRetriever = func(keyName, alias string, createNew bool, attempts int) (string, bool, error) {
	switch alias {
	case "alias_1":
		return "passphrase_1", false, nil
	case "alias_2":
		return "passphrase_2", false, nil
	}
	return "", false, errors.New("password alias no found")
}

// Create a temporary file, open a database connection to it, and create the
// necessary table.  Return the file name to use and clean up.
func initializeDB(t *testing.T) (tmpfilename string) {
	tmpFile, err := ioutil.TempFile("/tmp", "notary-test-sqlite-db-")
	require.NoError(t, err)
	tmpFile.Close()

	// We are using SQLite for the tests
	gormDB, err := gorm.Open("sqlite3", tmpFile.Name())
	require.NoError(t, err)

	// Ensure that the private_key table exists
	gormDB.CreateTable(&GormPrivateKey{})

	return tmpFile.Name()
}

// gets a key from the DB store, and asserts that the key is the expected key
func testGetSuccess(t *testing.T, dbStore *KeyDBStore, expectedKey data.PrivateKey) {
	retrKey, _, err := dbStore.GetKey(expectedKey.ID())
	require.NoError(t, err)
	require.Equal(t, retrKey, expectedKey)
}

// closes the DB connection first so we can test that the successful get was
// from the cache
func testGetSuccessFromCache(t *testing.T, dbStore *KeyDBStore,
	expectedKey data.PrivateKey) {

	err := dbStore.db.Close()
	require.NoError(t, err)

	testGetSuccess(t, dbStore, expectedKey)
}

// Creating a new KeyDBStore propagates any db opening error
func TestNewKeyDBStorePropagatesDBError(t *testing.T) {
	dbStore, err := NewKeyDBStore(retriever, "ignoredalias", "nodb", "somestring")
	require.Error(t, err)
	require.Nil(t, dbStore)
}

// Creating a key, on succcess, populates the cache.
func TestCreateSuccessPopulatesCache(t *testing.T) {
	testKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err)

	tmpFilename := initializeDB(t)
	defer os.Remove(tmpFilename)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "ignoredalias", "sqlite3", tmpFilename)
	require.NoError(t, err)

	// Test writing new key in database
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun/ignored"}, testKey)
	require.NoError(t, err)

	testGetSuccessFromCache(t, dbStore, testKey)
}

// Getting a key, on succcess, populates the cache.
func TestGetSuccessPopulatesCache(t *testing.T) {
	testKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err)

	tmpFilename := initializeDB(t)
	defer os.Remove(tmpFilename)

	// Create a new KeyDB store and add a key
	dbStore, err := NewKeyDBStore(retriever, "ignoredalias", "sqlite3", tmpFilename)
	require.NoError(t, err)
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun/ignored"}, testKey)
	require.NoError(t, err)

	// delete the cache
	dbStore.cachedKeys = make(map[string]data.PrivateKey)

	testGetSuccess(t, dbStore, testKey)
	testGetSuccessFromCache(t, dbStore, testKey)
}

func TestDoubleCreate(t *testing.T) {
	testKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err)

	anotherTestKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err)

	tmpFilename := initializeDB(t)
	defer os.Remove(tmpFilename)

	// Create a new KeyDB store and add a key
	dbStore, err := NewKeyDBStore(retriever, "ignoredalias", "sqlite3", tmpFilename)
	require.NoError(t, err)

	// Test writing new key in database/cache
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun/ignored"}, testKey)
	require.NoError(t, err)

	// Test writing the same key in the database. Should fail.
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun/ignored"}, testKey)
	require.Error(t, err, "failed to add private key to database:")

	// Test writing new key succeeds
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun/ignored"}, anotherTestKey)
	require.NoError(t, err)
}

func TestCreateDelete(t *testing.T) {
	testKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err)

	tmpFilename := initializeDB(t)
	defer os.Remove(tmpFilename)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "ignoredalias", "sqlite3", tmpFilename)
	require.NoError(t, err)

	// Test writing new key in database/cache
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: "", Gun: ""}, testKey)
	require.NoError(t, err)

	// Test deleting the key from the db
	err = dbStore.RemoveKey(testKey.ID())
	require.NoError(t, err)

	// This should fail, since it is neither in the cache nor the DB
	_, _, err = dbStore.GetKey(testKey.ID())
	require.Error(t, err, "signing key not found:")
}

func TestKeyRotation(t *testing.T) {
	testKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err)

	tmpFilename := initializeDB(t)
	defer os.Remove(tmpFilename)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(anotherRetriever, "alias_1", "sqlite3", tmpFilename)
	require.NoError(t, err)

	// Test writing new key in database/cache
	err = dbStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun/ignored"}, testKey)
	require.NoError(t, err)

	// Try rotating the key to alias-2
	err = dbStore.RotateKeyPassphrase(testKey.ID(), "alias_2")
	require.NoError(t, err)

	// Try rotating the key to alias-3
	err = dbStore.RotateKeyPassphrase(testKey.ID(), "alias_3")
	require.Error(t, err, "there should be no password for alias_3")
}

func TestDBHealthCheck(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "ignoredalias",
		"sqlite3", filepath.Join(tempBaseDir, "test_db"))
	require.NoError(t, err)

	// No key table, health check fails
	err = dbStore.HealthCheck()
	require.Error(t, err, "Cannot access table:")

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Heath check success because the table exists
	err = dbStore.HealthCheck()
	require.NoError(t, err)

	// Close the connection
	err = dbStore.db.Close()
	require.NoError(t, err)

	// Heath check fail because the connection is closed
	err = dbStore.HealthCheck()
	require.Error(t, err, "Cannot access table:")
}
