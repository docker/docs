package storage

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

    "github.com/endophage/gotuf/data"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// GormTUFFile represents a TUF file in the database
type GormTUFFile struct {
	ID      int    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Gun     string `sql:"type:varchar(255);not null"`
	Role    string `sql:"type:varchar(255);not null"`
	Version int
	Data    []byte `sql:"type:longblob"`
}

// TableName sets a specific table name for GormTUFFile
func (g GormTUFFile) TableName() string {
	return "tuf_files"
}

// GormTimestampKey represents a single timestamp key in the database
type GormTimestampKey struct {
	Gun    string `sql:"type:varchar(255)" gorm:"primary key"`
	Cipher string `sql:"type:varchar(30)"`
	Public []byte `sql:"type:blob;not null"`
}

// TableName sets a specific table name for our GormTimestampKey
func (g GormTimestampKey) TableName() string {
	return "timestamp_keys"
}

// SetUpSQLite creates a sqlite database for testing
func SetUpSQLite(t *testing.T) (*gorm.DB, *MySQLStorage) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	// We are using SQLite for the tests
	db, err := sql.Open("sqlite3", tempBaseDir+"test_db")
	assert.NoError(t, err)

	// Create the DB tables
	gormDB, _ := gorm.Open("sqlite3", db)
	query := gormDB.CreateTable(&GormTUFFile{})
	assert.NoError(t, query.Error)
	query = gormDB.Model(&GormTUFFile{}).AddUniqueIndex(
		"idx_gun", "gun", "role", "version")
	assert.NoError(t, query.Error)
	query = gormDB.CreateTable(&GormTimestampKey{})
	assert.NoError(t, query.Error)

	return &gormDB, NewMySQLStorage(db)
}

func TestMySQLUpdateCurrent(t *testing.T) {
	_, dbStore := SetUpSQLite(t)

	update := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	err := dbStore.UpdateCurrent("testGUN", update)
	assert.NoError(t, err, "Creating a row in an empty DB failed.")

    err = dbStore.UpdateCurrent("testGUN", update)
	assert.Error(t, err, "Error should not be nil")
    assert.IsType(t, &ErrOldVersion{}, err,
                  "Expected ErrOldVersion error type, got: %v", err)
    dbStore.DB.Close()
}

func TestMySQLUpdateMany(t *testing.T) {
    _, dbStore := SetUpSQLite(t)

	update1 := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	update2 := MetaUpdate{
		Role:    "targets",
		Version: 1,
		Data:    []byte("2"),
	}

    err := dbStore.UpdateMany("testGUN", []MetaUpdate{update1, update2})
    assert.NoError(t, err, "UpdateMany errored unexpectedly: %v", err)
    dbStore.DB.Close()
}


func TestMySQLUpdateManyDuplicateRollback(t *testing.T) {
    _, dbStore := SetUpSQLite(t)

    update := MetaUpdate{
        Role:    "root",
        Version: 0,
        Data:    []byte("1"),
    }

    err := dbStore.UpdateMany("testGUN", []MetaUpdate{update, update})
    assert.Error(t, err, "There should be an error updating twice.")
    // sqlite3 error and mysql error aren't compatible
    // assert.IsType(t, &ErrOldVersion{}, err,
    //               "UpdateMany returned wrong error type")

    // the whole transaction should have rolled back, so there should be
    // no entries.
    byt, err := dbStore.GetCurrent("testGUN", "root")
    assert.Nil(t, byt)
    assert.Error(t, err, "There should be an error Getting any entries")
    assert.IsType(t, &ErrNotFound{}, err, "Should get a not found error")

    dbStore.DB.Close()
}


func TestMySQLGetCurrent(t *testing.T) {
    _, dbStore := SetUpSQLite(t)

    byt, err := dbStore.GetCurrent("testGUN", "root")
    assert.Nil(t, byt)
    assert.Error(t, err, "There should be an error Getting an empty table")
    assert.IsType(t, &ErrNotFound{}, err, "Should get a not found error")

    // use UpdateCurrent to create one and test GetCurrent
    update := MetaUpdate{
        Role:    "root",
        Version: 0,
        Data:    []byte("1"),
    }
    err = dbStore.UpdateCurrent("testGUN", update)
    assert.NoError(t, err, "Creating a row in an empty DB failed.")

    byt, err = dbStore.GetCurrent("testGUN", "root")
    assert.NoError(t, err, "There should not be any errors getting.")
	assert.Equal(t, []byte("1"), byt, "Returned data was incorrect")

    dbStore.DB.Close()
}

func TestMySQLDelete(t *testing.T) {
    _, dbStore := SetUpSQLite(t)

    // Not testing deleting from an empty table, because that's not an error
    // in SQLite3

    // use UpdateCurrent to create one and test GetCurrent
    update := MetaUpdate{
        Role:    "root",
        Version: 0,
        Data:    []byte("1"),
    }
    err := dbStore.UpdateCurrent("testGUN", update)
    assert.NoError(t, err, "Creating a row in an empty DB failed.")

    err = dbStore.Delete("testGUN")
    assert.NoError(t, err, "There should not be any errors deleting.")

    dbStore.DB.Close()
}

func TestMySQLGetTimestampKeyNoKey(t *testing.T) {
    _, dbStore := SetUpSQLite(t)

    cipher, public, err := dbStore.GetTimestampKey("testGUN")
    assert.Equal(t, data.KeyAlgorithm(""), cipher)
    assert.Nil(t, public)
    assert.IsType(t, &ErrNoKey{}, err,
                  "Expected ErrNoKey from GetTimestampKey")

    err = dbStore.SetTimestampKey("testGUN", "testCipher", []byte("1"))
    assert.NoError(t, err, "Inserting timestamp into empty DB should succeed")

    cipher, public, err = dbStore.GetTimestampKey("testGUN")
    assert.Equal(t, data.KeyAlgorithm("testCipher"), cipher,
                 "Returned cipher was incorrect")
    assert.Equal(t, []byte("1"), public, "Returned pubkey was incorrect")
}

func TestMySQLSetTimestampKeyExists(t *testing.T) {
    _, dbStore := SetUpSQLite(t)

    err := dbStore.SetTimestampKey("testGUN", "testCipher", []byte("1"))
    assert.NoError(t, err, "Inserting timestamp into empty DB should succeed")

    err = dbStore.SetTimestampKey("testGUN", "testCipher", []byte("1"))
    // sqlite3 error and mysql error aren't compatible

    // assert.IsType(t, &ErrTimestampKeyExists{}, err,
    //               "Expected ErrTimestampKeyExists from SetTimestampKey")

    dbStore.DB.Close()
}
