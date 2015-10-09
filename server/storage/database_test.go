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
	Gun    string `sql:"type:varchar(255);unique" gorm:"primary key"`
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

	// verify that the tables are empty
	var count int
	for _, model := range [2]interface{}{&GormTUFFile{}, &GormTimestampKey{}} {
		query = gormDB.Model(model).Count(&count)
		assert.NoError(t, query.Error)
		assert.Equal(t, 0, count)
	}

	return &gormDB, NewMySQLStorage(db)
}

func TestMySQLUpdateCurrent(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

    // UpdateCurrent should succeed
	update := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	err := dbStore.UpdateCurrent("testGUN", update)
	assert.NoError(t, err, "Creating a row in an empty DB failed.")

    // There should just be one row
    var rows []GormTUFFile
    query := gormDB.Model(&GormTUFFile{}).Find(&rows)
    assert.NoError(t, query.Error)
    assert.Equal(
        t,
        []GormTUFFile{
            GormTUFFile{ID: 1, Gun: "testGUN", Role: "root", Version: 0,
                Data: []byte("1")},
        },
        rows)

	dbStore.DB.Close()
}

func TestMySQLUpdateCurrentError(t *testing.T) {
    gormDB, dbStore := SetUpSQLite(t)

    // insert row
    query := gormDB.Create(&GormTUFFile{
        Gun:     "testGUN",
        Role:    "root",
        Version: 0,
        Data:    []byte("1"),
    })
    assert.NoError(t, query.Error, "Creating a row in an empty DB failed.")

    // UpdateCurrent should fail due to clash with prevoius row
    update := MetaUpdate{
        Role:    "root",
        Version: 0,
        Data:    []byte("1"),
    }
    err := dbStore.UpdateCurrent("testGUN", update)
    assert.Error(t, err, "Error should not be nil")
    assert.IsType(t, &ErrOldVersion{}, err,
        "Expected ErrOldVersion error type, got: %v", err)

    // There should just be one row
    var rows []GormTUFFile
    query = gormDB.Model(&GormTUFFile{}).Find(&rows)
    assert.NoError(t, query.Error)
    assert.Equal(
        t,
        []GormTUFFile{
            GormTUFFile{ID: 1, Gun: "testGUN", Role: "root", Version: 0,
                Data: []byte("1")},
        },
        rows)

    dbStore.DB.Close()
}

func TestMySQLUpdateMany(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

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

	expected := []GormTUFFile{
		GormTUFFile{ID: 1, Gun: "testGUN", Role: "root", Version: 0,
			Data: []byte("1")},
		GormTUFFile{ID: 2, Gun: "testGUN", Role: "targets", Version: 1,
			Data: []byte("2")},
	}

	err := dbStore.UpdateMany("testGUN", []MetaUpdate{update1, update2})
	assert.NoError(t, err, "UpdateMany errored unexpectedly: %v", err)

	var rows []GormTUFFile
	query := gormDB.Model(&GormTUFFile{}).Find(&rows)
	assert.NoError(t, query.Error)
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, expected, rows)

	dbStore.DB.Close()
}

func TestMySQLUpdateManyDuplicateRollback(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

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
	var count int
	query := gormDB.Model(&GormTUFFile{}).Count(&count)
	assert.NoError(t, query.Error)
	assert.Equal(t, 0, count)

	dbStore.DB.Close()
}

func TestMySQLGetCurrent(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

	byt, err := dbStore.GetCurrent("testGUN", "root")
	assert.Nil(t, byt)
	assert.Error(t, err, "There should be an error Getting an empty table")
	assert.IsType(t, &ErrNotFound{}, err, "Should get a not found error")

	query := gormDB.Create(&GormTUFFile{
		Gun:     "testGUN",
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	})
	assert.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	byt, err = dbStore.GetCurrent("testGUN", "root")
	assert.NoError(t, err, "There should not be any errors getting.")
	assert.Equal(t, []byte("1"), byt, "Returned data was incorrect")

	dbStore.DB.Close()
}

func TestMySQLDelete(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

	// Not testing deleting from an empty table, because that's not an error
	// in SQLite3

	// use UpdateCurrent to create one and test GetCurrent
	query := gormDB.Create(&GormTUFFile{
		Gun:     "testGUN",
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	})
	assert.NoError(t, query.Error, "Creating a row in an empty DB failed.")

	err := dbStore.Delete("testGUN")
	assert.NoError(t, err, "There should not be any errors deleting.")

	// verify deletion
	var count int
	query = gormDB.Model(&GormTUFFile{}).Count(&count)
	assert.NoError(t, query.Error)
	assert.Equal(t, 0, count)

	dbStore.DB.Close()
}

func TestMySQLGetTimestampKeyNoKey(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

	cipher, public, err := dbStore.GetTimestampKey("testGUN")
	assert.Equal(t, data.KeyAlgorithm(""), cipher)
	assert.Nil(t, public)
	assert.IsType(t, &ErrNoKey{}, err,
		"Expected ErrNoKey from GetTimestampKey")

	query := gormDB.Create(&GormTimestampKey{
		Gun:    "testGUN",
		Cipher: "testCipher",
		Public: []byte("1"),
	})
	assert.NoError(
		t, query.Error, "Inserting timestamp into empty DB should succeed")

	cipher, public, err = dbStore.GetTimestampKey("testGUN")
	assert.Equal(t, data.KeyAlgorithm("testCipher"), cipher,
		"Returned cipher was incorrect")
	assert.Equal(t, []byte("1"), public, "Returned pubkey was incorrect")
}

func TestMySQLSetTimestampKeyExists(t *testing.T) {
	gormDB, dbStore := SetUpSQLite(t)

	err := dbStore.SetTimestampKey("testGUN", "testCipher", []byte("1"))
	assert.NoError(t, err, "Inserting timestamp into empty DB should succeed")

	err = dbStore.SetTimestampKey("testGUN", "testCipher", []byte("1"))
	assert.Error(t, err)
	// sqlite3 error and mysql error aren't compatible

	// assert.IsType(t, &ErrTimestampKeyExists{}, err,
	//               "Expected ErrTimestampKeyExists from SetTimestampKey")

	var rows []GormTimestampKey
	query := gormDB.Model(&GormTimestampKey{}).Find(&rows)
	assert.NoError(t, query.Error)
	assert.Equal(t, 1, len(rows))
	assert.Equal(
		t,
		GormTimestampKey{Gun: "testGUN", Cipher: "testCipher",
			Public: []byte("1")},
		rows[0])

	dbStore.DB.Close()
}
