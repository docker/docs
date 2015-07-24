package trustmanager

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var retriever = func(string, string, bool, int) (string, bool, error) {
	return "abcgdhfjdhfjhfgdhejnfhdfgshdjfbv", false, nil
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
	dbStore, err := NewKeyDBStore(retriever, "sqlite3", db)
	assert.NoError(t, err)

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Test writing new key in database/cache
	err = dbStore.AddKey("", "", testKey)
	fmt.Println(err)
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

func TestCreateDelete(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	testKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// We are using SQLite for the tests
	db, err := sql.Open("sqlite3", tempBaseDir+"test_db")
	assert.NoError(t, err)

	// Create a new KeyDB store
	dbStore, err := NewKeyDBStore(retriever, "sqlite3", db)
	assert.NoError(t, err)

	// Ensure that the private_key table exists
	dbStore.db.CreateTable(&GormPrivateKey{})

	// Test writing new key in database/cache
	err = dbStore.AddKey("", "", testKey)
	fmt.Println(err)
	assert.NoError(t, err)

	// Test deleting the key from the db
	err = dbStore.RemoveKey(testKey.ID())
	assert.NoError(t, err)

	// This should fail
	_, _, err = dbStore.GetKey(testKey.ID())
	assert.Error(t, err, "signing key not found:")
}
