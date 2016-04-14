package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/docker/notary/tuf/data"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// SampleTUF returns a sample TUFFile with the given Version (ID will have
// to be set independently)
func SampleTUF(version int) TUFFile {
	return SampleCustomTUF(data.CanonicalRootRole, "testGUN", []byte("1"), version)
}

func SampleCustomTUF(role, gun string, data []byte, version int) TUFFile {
	checksum := sha256.Sum256(data)
	hexChecksum := hex.EncodeToString(checksum[:])
	return TUFFile{
		Gun:     gun,
		Role:    role,
		Version: version,
		Sha256:  hexChecksum,
		Data:    data,
	}
}

func SampleUpdate(version int) MetaUpdate {
	return MetaUpdate{
		Role:    "root",
		Version: version,
		Data:    []byte("1"),
	}
}

// SetUpSQLite creates a sqlite database for testing
func SetUpSQLite(t *testing.T, dbDir string) (*gorm.DB, *SQLStorage) {
	dbStore, err := NewSQLStorage("sqlite3", dbDir+"test_db")
	require.NoError(t, err)

	// Create the DB tables
	err = CreateTUFTable(dbStore.DB)
	require.NoError(t, err)

	err = CreateKeyTable(dbStore.DB)
	require.NoError(t, err)

	// verify that the tables are empty
	var count int
	for _, model := range [2]interface{}{&TUFFile{}, &Key{}} {
		query := dbStore.DB.Model(model).Count(&count)
		require.NoError(t, query.Error)
		require.Equal(t, 0, count)
	}
	return &dbStore.DB, dbStore
}

// TestSQLUpdateCurrent asserts that UpdateCurrent will add a new TUF file
// if no previous version existed.
func TestSQLUpdateCurrentNew(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	// Adding a new TUF file should succeed
	err = dbStore.UpdateCurrent("testGUN", SampleUpdate(0))
	require.NoError(t, err, "Creating a row in an empty DB failed.")

	// There should just be one row
	var rows []TUFFile
	query := gormDB.Select("ID, Gun, Role, Version, Sha256, Data").Find(&rows)
	require.NoError(t, query.Error)

	expected := SampleTUF(0)
	expected.ID = 1
	require.Equal(t, []TUFFile{expected}, rows)
}

// TestSQLUpdateCurrentNewVersion asserts that UpdateCurrent will add a
// new (higher) version of an existing TUF file
func TestSQLUpdateCurrentNewVersion(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	// insert row
	oldVersion := SampleTUF(0)
	query := gormDB.Create(&oldVersion)
	require.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	// UpdateCurrent with a newer version should succeed
	update := SampleUpdate(2)
	err = dbStore.UpdateCurrent("testGUN", update)
	require.NoError(t, err, "Creating a row in an empty DB failed.")

	// There should just be one row
	var rows []TUFFile
	query = gormDB.Select("ID, Gun, Role, Version, Sha256, Data").Find(&rows)
	require.NoError(t, query.Error)

	oldVersion.Model = gorm.Model{ID: 1}
	expected := SampleTUF(2)
	expected.Model = gorm.Model{ID: 2}
	require.Equal(t, []TUFFile{oldVersion, expected}, rows)
}

// TestSQLUpdateCurrentOldVersionError asserts that an error is raised if
// trying to update to an older version of a TUF file.
func TestSQLUpdateCurrentOldVersionError(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	// insert row
	newVersion := SampleTUF(3)
	query := gormDB.Create(&newVersion)
	require.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	// UpdateCurrent should fail due to the version being lower than the
	// previous row
	err = dbStore.UpdateCurrent("testGUN", SampleUpdate(0))
	require.Error(t, err, "Error should not be nil")
	require.IsType(t, &ErrOldVersion{}, err,
		"Expected ErrOldVersion error type, got: %v", err)

	// There should just be one row
	var rows []TUFFile
	query = gormDB.Select("ID, Gun, Role, Version, Sha256, Data").Find(&rows)
	require.NoError(t, query.Error)

	newVersion.Model = gorm.Model{ID: 1}
	require.Equal(t, []TUFFile{newVersion}, rows)

	dbStore.DB.Close()
}

// TestSQLUpdateMany asserts that inserting multiple updates succeeds if the
// updates do not conflict with each.
func TestSQLUpdateMany(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.UpdateMany("testGUN", []MetaUpdate{
		SampleUpdate(0),
		{
			Role:    "targets",
			Version: 1,
			Data:    []byte("2"),
		},
		SampleUpdate(2),
	})
	require.NoError(t, err, "UpdateMany errored unexpectedly: %v", err)

	gorm1 := SampleTUF(0)
	gorm1.ID = 1
	data := []byte("2")
	checksum := sha256.Sum256(data)
	hexChecksum := hex.EncodeToString(checksum[:])
	gorm2 := TUFFile{
		Model: gorm.Model{ID: 2}, Gun: "testGUN", Role: "targets",
		Version: 1, Sha256: hexChecksum, Data: data}
	gorm3 := SampleTUF(2)
	gorm3.ID = 3
	expected := []TUFFile{gorm1, gorm2, gorm3}

	var rows []TUFFile
	query := gormDB.Select("ID, Gun, Role, Version, Sha256, Data").Find(&rows)
	require.NoError(t, query.Error)
	require.Equal(t, expected, rows)

	dbStore.DB.Close()
}

// TestSQLUpdateManyVersionOrder asserts that inserting updates with
// non-monotonic versions still succeeds.
func TestSQLUpdateManyVersionOrder(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.UpdateMany(
		"testGUN", []MetaUpdate{SampleUpdate(2), SampleUpdate(0)})
	require.NoError(t, err)

	// the whole transaction should have rolled back, so there should be
	// no entries.
	gorm1 := SampleTUF(2)
	gorm1.ID = 1
	gorm2 := SampleTUF(0)
	gorm2.ID = 2

	var rows []TUFFile
	query := gormDB.Select("ID, Gun, Role, Version, Sha256, Data").Find(&rows)
	require.NoError(t, query.Error)
	require.Equal(t, []TUFFile{gorm1, gorm2}, rows)

	dbStore.DB.Close()
}

// TestSQLUpdateManyDuplicateRollback asserts that inserting duplicate
// updates fails.
func TestSQLUpdateManyDuplicateRollback(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	update := SampleUpdate(0)
	err = dbStore.UpdateMany("testGUN", []MetaUpdate{update, update})
	require.Error(
		t, err, "There should be an error updating the same data twice.")
	require.IsType(t, &ErrOldVersion{}, err,
		"UpdateMany returned wrong error type")

	// the whole transaction should have rolled back, so there should be
	// no entries.
	var count int
	query := gormDB.Model(&TUFFile{}).Count(&count)
	require.NoError(t, query.Error)
	require.Equal(t, 0, count)

	dbStore.DB.Close()
}

func TestSQLGetCurrent(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	_, byt, err := dbStore.GetCurrent("testGUN", "root")
	require.Nil(t, byt)
	require.Error(t, err, "There should be an error Getting an empty table")
	require.IsType(t, ErrNotFound{}, err, "Should get a not found error")

	tuf := SampleTUF(0)
	query := gormDB.Create(&tuf)
	require.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	cDate, byt, err := dbStore.GetCurrent("testGUN", "root")
	require.NoError(t, err, "There should not be any errors getting.")
	require.Equal(t, []byte("1"), byt, "Returned data was incorrect")
	// the update date was sometime wthin the last minute
	fmt.Println(cDate)
	require.True(t, cDate.After(time.Now().Add(-1*time.Minute)))
	require.True(t, cDate.Before(time.Now().Add(5*time.Second)))

	dbStore.DB.Close()
}

func TestSQLDelete(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	tuf := SampleTUF(0)
	query := gormDB.Create(&tuf)
	require.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	err = dbStore.Delete("testGUN")
	require.NoError(t, err, "There should not be any errors deleting.")

	// verify deletion
	var count int
	query = gormDB.Model(&TUFFile{}).Count(&count)
	require.NoError(t, query.Error)
	require.Equal(t, 0, count)

	dbStore.DB.Close()
}

func TestSQLGetKeyNoKey(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	cipher, public, err := dbStore.GetKey("testGUN", data.CanonicalTimestampRole)
	require.Equal(t, "", cipher)
	require.Nil(t, public)
	require.IsType(t, &ErrNoKey{}, err,
		"Expected ErrNoKey from GetKey")

	query := gormDB.Create(&Key{
		Gun:    "testGUN",
		Role:   data.CanonicalTimestampRole,
		Cipher: "testCipher",
		Public: []byte("1"),
	})
	require.NoError(
		t, query.Error, "Inserting timestamp into empty DB should succeed")

	cipher, public, err = dbStore.GetKey("testGUN", data.CanonicalTimestampRole)
	require.Equal(t, "testCipher", cipher,
		"Returned cipher was incorrect")
	require.Equal(t, []byte("1"), public, "Returned pubkey was incorrect")
}

func TestSQLSetKeyExists(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.SetKey("testGUN", data.CanonicalTimestampRole, "testCipher", []byte("1"))
	require.NoError(t, err, "Inserting timestamp into empty DB should succeed")

	err = dbStore.SetKey("testGUN", data.CanonicalTimestampRole, "testCipher", []byte("1"))
	require.Error(t, err)
	require.IsType(t, &ErrKeyExists{}, err,
		"Expected ErrKeyExists from SetKey")

	var rows []Key
	query := gormDB.Select("ID, Gun, Cipher, Public").Find(&rows)
	require.NoError(t, query.Error)

	expected := Key{Gun: "testGUN", Cipher: "testCipher",
		Public: []byte("1")}
	expected.Model = gorm.Model{ID: 1}

	require.Equal(t, []Key{expected}, rows)

	dbStore.DB.Close()
}

func TestSQLSetKeyMultipleRoles(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.SetKey("testGUN", data.CanonicalTimestampRole, "testCipher", []byte("1"))
	require.NoError(t, err, "Inserting timestamp into empty DB should succeed")

	err = dbStore.SetKey("testGUN", data.CanonicalSnapshotRole, "testCipher", []byte("1"))
	require.NoError(t, err, "Inserting snapshot key into DB with timestamp key should succeed")

	var rows []Key
	query := gormDB.Select("ID, Gun, Role, Cipher, Public").Find(&rows)
	require.NoError(t, query.Error)

	expectedTS := Key{Gun: "testGUN", Role: "timestamp", Cipher: "testCipher",
		Public: []byte("1")}
	expectedTS.Model = gorm.Model{ID: 1}

	expectedSN := Key{Gun: "testGUN", Role: "snapshot", Cipher: "testCipher",
		Public: []byte("1")}
	expectedSN.Model = gorm.Model{ID: 2}

	require.Equal(t, []Key{expectedTS, expectedSN}, rows)

	dbStore.DB.Close()
}

func TestSQLSetKeyMultipleGuns(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	gormDB, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.SetKey("testGUN", data.CanonicalTimestampRole, "testCipher", []byte("1"))
	require.NoError(t, err, "Inserting timestamp into empty DB should succeed")

	err = dbStore.SetKey("testAnotherGUN", data.CanonicalTimestampRole, "testCipher", []byte("1"))
	require.NoError(t, err, "Inserting snapshot key into DB with timestamp key should succeed")

	var rows []Key
	query := gormDB.Select("ID, Gun, Role, Cipher, Public").Find(&rows)
	require.NoError(t, query.Error)

	expected1 := Key{Gun: "testGUN", Role: "timestamp", Cipher: "testCipher",
		Public: []byte("1")}
	expected1.Model = gorm.Model{ID: 1}

	expected2 := Key{Gun: "testAnotherGUN", Role: "timestamp", Cipher: "testCipher",
		Public: []byte("1")}
	expected2.Model = gorm.Model{ID: 2}

	require.Equal(t, []Key{expected1, expected2}, rows)

	dbStore.DB.Close()
}

func TestSQLSetKeySameRoleGun(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	_, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.SetKey("testGUN", data.CanonicalTimestampRole, "testCipher", []byte("1"))
	require.NoError(t, err, "Inserting timestamp into empty DB should succeed")

	err = dbStore.SetKey("testGUN", data.CanonicalTimestampRole, "testCipher", []byte("2"))
	require.Error(t, err)
	require.IsType(t, &ErrKeyExists{}, err,
		"Expected ErrKeyExists from SetKey")

	dbStore.DB.Close()
}

// TestDBCheckHealthTableMissing asserts that the health check fails if one or
// both the tables are missing.
func TestDBCheckHealthTableMissing(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	_, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	dbStore.DropTable(&TUFFile{})
	dbStore.DropTable(&Key{})

	// No tables, health check fails
	err = dbStore.CheckHealth()
	require.Error(t, err, "Cannot access table:")

	// only one table existing causes health check to fail
	CreateTUFTable(dbStore.DB)
	err = dbStore.CheckHealth()
	require.Error(t, err, "Cannot access table:")
	dbStore.DropTable(&TUFFile{})

	CreateKeyTable(dbStore.DB)
	err = dbStore.CheckHealth()
	require.Error(t, err, "Cannot access table:")
}

// TestDBCheckHealthDBCOnnection asserts that if the DB is not connectable, the
// health check fails.
func TestDBCheckHealthDBConnectionFail(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	_, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.Close()
	require.NoError(t, err)

	err = dbStore.CheckHealth()
	require.Error(t, err, "Cannot access table:")
}

// TestDBCheckHealthSuceeds asserts that if the DB is connectable and both
// tables exist, the health check succeeds.
func TestDBCheckHealthSucceeds(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	_, dbStore := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	err = dbStore.CheckHealth()
	require.NoError(t, err)
}

func TestDBGetChecksum(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	_, store := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	ts := data.SignedTimestamp{
		Signatures: make([]data.Signature, 0),
		Signed: data.Timestamp{
			SignedCommon: data.SignedCommon{
				Type:    data.TUFTypes[data.CanonicalTimestampRole],
				Version: 1,
				Expires: data.DefaultExpires(data.CanonicalTimestampRole),
			},
		},
	}
	j, err := json.Marshal(&ts)
	require.NoError(t, err)
	update := MetaUpdate{
		Role:    data.CanonicalTimestampRole,
		Version: 1,
		Data:    j,
	}
	checksumBytes := sha256.Sum256(j)
	checksum := hex.EncodeToString(checksumBytes[:])

	store.UpdateCurrent("gun", update)

	// create and add a newer timestamp. We're going to try and get the one
	// created above by checksum
	ts = data.SignedTimestamp{
		Signatures: make([]data.Signature, 0),
		Signed: data.Timestamp{
			SignedCommon: data.SignedCommon{
				Type:    data.TUFTypes[data.CanonicalTimestampRole],
				Version: 2,
				Expires: data.DefaultExpires(data.CanonicalTimestampRole),
			},
		},
	}
	newJ, err := json.Marshal(&ts)
	require.NoError(t, err)
	update = MetaUpdate{
		Role:    data.CanonicalTimestampRole,
		Version: 2,
		Data:    newJ,
	}

	store.UpdateCurrent("gun", update)

	cDate, data, err := store.GetChecksum("gun", data.CanonicalTimestampRole, checksum)
	require.NoError(t, err)
	require.EqualValues(t, j, data)
	// the creation date was sometime wthin the last minute
	require.True(t, cDate.After(time.Now().Add(-1*time.Minute)))
	require.True(t, cDate.Before(time.Now().Add(5*time.Second)))
}

func TestDBGetChecksumNotFound(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	_, store := SetUpSQLite(t, tempBaseDir)
	defer os.RemoveAll(tempBaseDir)

	_, _, err = store.GetChecksum("gun", data.CanonicalTimestampRole, "12345")
	require.Error(t, err)
	require.IsType(t, ErrNotFound{}, err)
}
