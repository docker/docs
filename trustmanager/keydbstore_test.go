package trustmanager

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var retriever = func(string, string, bool, int) (string, bool, error) {
	return "passphrase-1", false, nil
}

var anotherRetriever = func(keyName, alias string, createNew bool, attempts int) (string, bool, error) {
	switch alias {
	case "alias-1":
		return "passphrase-1", false, nil
	case "alias-2":
		return "passphrase-2", false, nil
	}
	return "", false, errors.New("password alias no found")
}

func TestCreateRead(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	testKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// We are using SQLite for the tests
	db, err := sql.Open("sqlite3", tempBaseDir+"test_db")
	assert.NoError(t, err)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "", "sqlite3", db)
	assert.NoError(t, err)

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Test writing new key in database/cache
	err = dbStore.AddKey("", "", testKey)
	assert.NoError(t, err)

	// Test retrieval of key from DB
	delete(dbStore.cachedKeys, testKey.ID())

	retrKey, _, err := dbStore.GetKey(testKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, retrKey, testKey)

	// Tests retrieval of key from Cache
	// Close database connection
	err = dbStore.db.Close()
	assert.NoError(t, err)

	retrKey, _, err = dbStore.GetKey(testKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, retrKey, testKey)
}

func TestDoubleCreate(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	testKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	anotherTestKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// We are using SQLite for the tests
	db, err := sql.Open("sqlite3", tempBaseDir+"test_db")
	assert.NoError(t, err)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "", "sqlite3", db)
	assert.NoError(t, err)

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Test writing new key in database/cache
	err = dbStore.AddKey("", "", testKey)
	assert.NoError(t, err)

	// Test writing the same key in the database. Should fail.
	err = dbStore.AddKey("", "", testKey)
	assert.Error(t, err, "failed to add private key to database:")

	// Test writing new key succeeds
	err = dbStore.AddKey("", "", anotherTestKey)
	assert.NoError(t, err)
}

func TestCreateDelete(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	testKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// We are using SQLite for the tests
	db, err := sql.Open("sqlite3", tempBaseDir+"test_db")
	assert.NoError(t, err)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "", "sqlite3", db)
	assert.NoError(t, err)

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Test writing new key in database/cache
	err = dbStore.AddKey("", "", testKey)
	assert.NoError(t, err)

	// Test deleting the key from the db
	err = dbStore.RemoveKey(testKey.ID())
	assert.NoError(t, err)

	// This should fail
	_, _, err = dbStore.GetKey(testKey.ID())
	assert.Error(t, err, "signing key not found:")
}

func TestKeyRotation(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	testKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// We are using SQLite for the tests
	db, err := sql.Open("sqlite3", tempBaseDir+"test_db")
	assert.NoError(t, err)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(anotherRetriever, "alias-1", "sqlite3", db)
	assert.NoError(t, err)

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Test writing new key in database/cache
	err = dbStore.AddKey("", "", testKey)
	assert.NoError(t, err)

	// Try rotating the key to alias-2
	err = dbStore.RotateKeyPassphrase(testKey.ID(), "alias-2")
	assert.NoError(t, err)

	// Try rotating the key to alias-3
	err = dbStore.RotateKeyPassphrase(testKey.ID(), "alias-3")
	assert.Error(t, err, "password alias no found")
}
