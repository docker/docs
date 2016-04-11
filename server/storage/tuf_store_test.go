package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/testutils"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// SetupTUFSQLite creates a sqlite database for testing, wrapped by a TUFMetaStorage
func SetupTUFSQLite(t *testing.T, dbDir string) (*gorm.DB, *TUFMetaStorage) {
	dbStore, err := NewSQLStorage("sqlite3", dbDir+"test_db")
	require.NoError(t, err)

	consistentDBStore := NewTUFMetaStorage(dbStore)

	embeddedDB := dbStore.DB
	// Create the DB tables
	err = CreateTUFTable(embeddedDB)
	require.NoError(t, err)

	err = CreateKeyTable(embeddedDB)
	require.NoError(t, err)

	// verify that the tables are empty
	var count int
	for _, model := range [2]interface{}{&TUFFile{}, &Key{}} {
		query := embeddedDB.Model(model).Count(&count)
		require.NoError(t, query.Error)
		require.Equal(t, 0, count)
	}
	return &embeddedDB, consistentDBStore
}

// TestTUFSQLGetCurrent asserts that GetCurrent walks from the current timestamp metadata
// to the snapshot specified in the checksum, to potentially other role metadata by checksum
func TestTUFSQLGetCurrent(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, tufDBStore := SetupTUFSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)
	defer gormDB.Close()

	ConsistentEmptyGetCurrentTest(t, tufDBStore)

	// put an arbitrary piece of data in the database
	tuf := SampleTUF(1)
	query := gormDB.Create(&tuf)
	require.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	ConsistentMissingTSAndSnapGetCurrentTest(t, tufDBStore)

	// Note that get by checksum succeeds, since it does not try to walk timestamp/snapshot
	GetChecksumFoundTest(t, tufDBStore, tuf)

	// Now setup a valid tuf repo and use it to ensure we walk correctly
	validTUFRepo, _, err := testutils.EmptyRepo("testGUN")
	require.NoError(t, err)

	// Add the timestamp, snapshot, targets, and root to the database
	tufData, err := json.Marshal(validTUFRepo.Timestamp)
	require.NoError(t, err)
	tsTUF := SampleCustomTUF(data.CanonicalTimestampRole, "testGUN", tufData, validTUFRepo.Timestamp.Signed.Version)
	query = gormDB.Create(&tsTUF)
	require.NoError(t, query.Error, "Creating a row for timestamp in DB failed.")

	tufData, err = json.Marshal(validTUFRepo.Snapshot)
	require.NoError(t, err)
	snapTUF := SampleCustomTUF(data.CanonicalSnapshotRole, "testGUN", tufData, validTUFRepo.Snapshot.Signed.Version)
	query = gormDB.Create(&snapTUF)
	require.NoError(t, query.Error, "Creating a row for snapshot in DB failed.")

	tufData, err = json.Marshal(validTUFRepo.Targets[data.CanonicalTargetsRole])
	require.NoError(t, err)
	targetsTUF := SampleCustomTUF(data.CanonicalTargetsRole, "testGUN", tufData, validTUFRepo.Targets[data.CanonicalTargetsRole].Signed.Version)
	query = gormDB.Create(&targetsTUF)
	require.NoError(t, query.Error, "Creating a row for targets in DB failed.")

	tufData, err = json.Marshal(validTUFRepo.Root)
	require.NoError(t, err)
	rootTUF := SampleCustomTUF(data.CanonicalRootRole, "testGUN", tufData, validTUFRepo.Root.Signed.Version)
	query = gormDB.Create(&rootTUF)
	require.NoError(t, query.Error, "Creating a row for root in DB failed.")

	// GetCurrent on all of these roles should succeed
	ConsistentGetCurrentFoundTest(t, tufDBStore, tsTUF)
	ConsistentGetCurrentFoundTest(t, tufDBStore, snapTUF)
	ConsistentGetCurrentFoundTest(t, tufDBStore, targetsTUF)
	ConsistentGetCurrentFoundTest(t, tufDBStore, rootTUF)

	// Delete snapshot
	query = gormDB.Delete(&snapTUF)
	require.NoError(t, query.Error, "Deleting a row for snapshot in DB failed.")

	// GetCurrent snapshot lookup should still succeed because of caching
	ConsistentGetCurrentFoundTest(t, tufDBStore, snapTUF)

	// targets and root lookup on GetCurrent should also still succeed because of caching
	ConsistentGetCurrentFoundTest(t, tufDBStore, targetsTUF)
	ConsistentGetCurrentFoundTest(t, tufDBStore, rootTUF)
}
