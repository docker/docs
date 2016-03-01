package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// SQLStorage implements a versioned store using a relational database.
// See server/storage/models.go
type SQLStorage struct {
	gorm.DB
}

// NewSQLStorage is a convenience method to create a SQLStorage
func NewSQLStorage(dialect string, args ...interface{}) (*SQLStorage, error) {
	gormDB, err := gorm.Open(dialect, args...)
	if err != nil {
		return nil, err
	}
	return &SQLStorage{
		DB: gormDB,
	}, nil
}

// translateOldVersionError captures DB errors, and attempts to translate
// duplicate entry - currently only supports MySQL and Sqlite3
func translateOldVersionError(err error) error {
	switch err := err.(type) {
	case *mysql.MySQLError:
		// https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html
		// 1022 = Can't write; duplicate key in table '%s'
		// 1062 = Duplicate entry '%s' for key %d
		if err.Number == 1022 || err.Number == 1062 {
			return &ErrOldVersion{}
		}
	}
	return err
}

// UpdateCurrent updates a single TUF.
func (db *SQLStorage) UpdateCurrent(gun string, update MetaUpdate) error {
	// ensure we're not inserting an immediately old version - can't use the
	// struct, because that only works with non-zero values, and Version
	// can be 0.
	exists := db.Where("gun = ? and role = ? and version >= ?",
		gun, update.Role, update.Version).First(&TUFFile{})

	if !exists.RecordNotFound() {
		return &ErrOldVersion{}
	}
	checksum := sha256.Sum256(update.Data)
	return translateOldVersionError(db.Create(&TUFFile{
		Gun:     gun,
		Role:    update.Role,
		Version: update.Version,
		Sha256:  hex.EncodeToString(checksum[:]),
		Data:    update.Data,
	}).Error)
}

// UpdateMany atomically updates many TUF records in a single transaction
func (db *SQLStorage) UpdateMany(gun string, updates []MetaUpdate) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	rollback := func(err error) error {
		if rxErr := tx.Rollback().Error; rxErr != nil {
			logrus.Error("Failed on Tx rollback with error: ", rxErr.Error())
			return rxErr
		}
		return err
	}

	var (
		query *gorm.DB
		added = make(map[uint]bool)
	)
	for _, update := range updates {
		// This looks like the same logic as UpdateCurrent, but if we just
		// called, version ordering in the updates list must be enforced
		// (you cannot insert the version 2 before version 1).  And we do
		// not care about monotonic ordering in the updates.
		query = db.Where("gun = ? and role = ? and version >= ?",
			gun, update.Role, update.Version).First(&TUFFile{})

		if !query.RecordNotFound() {
			return rollback(&ErrOldVersion{})
		}

		var row TUFFile
		checksum := sha256.Sum256(update.Data)
		hexChecksum := hex.EncodeToString(checksum[:])
		query = tx.Where(map[string]interface{}{
			"gun":     gun,
			"role":    update.Role,
			"version": update.Version,
		}).Attrs("data", update.Data).Attrs("sha256", hexChecksum).FirstOrCreate(&row)

		if query.Error != nil {
			return rollback(translateOldVersionError(query.Error))
		}
		// it's previously been added, which means it's a duplicate entry
		// in the same transaction
		if _, ok := added[row.ID]; ok {
			return rollback(&ErrOldVersion{})
		}
		added[row.ID] = true
	}
	return tx.Commit().Error
}

// GetCurrent gets a specific TUF record
func (db *SQLStorage) GetCurrent(gun, tufRole string) ([]byte, error) {
	var row TUFFile
	q := db.Select("data").Where(&TUFFile{Gun: gun, Role: tufRole}).Order("version desc").Limit(1).First(&row)
	return returnRead(q, row)
}

// GetChecksum gets a specific TUF record by its hex checksum
func (db *SQLStorage) GetChecksum(gun, tufRole, checksum string) ([]byte, error) {
	var row TUFFile
	q := db.Select("data").Where(
		&TUFFile{
			Gun:    gun,
			Role:   tufRole,
			Sha256: checksum,
		},
	).First(&row)
	return returnRead(q, row)
}

func returnRead(q *gorm.DB, row TUFFile) ([]byte, error) {
	if q.RecordNotFound() {
		return nil, ErrNotFound{}
	} else if q.Error != nil {
		return nil, q.Error
	}
	return row.Data, nil
}

// Delete deletes all the records for a specific GUN
func (db *SQLStorage) Delete(gun string) error {
	return db.Where(&TUFFile{Gun: gun}).Delete(TUFFile{}).Error
}

// GetKey returns the Public Key data for a gun+role
func (db *SQLStorage) GetKey(gun, role string) (algorithm string, public []byte, err error) {
	logrus.Debugf("retrieving timestamp key for %s:%s", gun, role)

	var row Key
	query := db.Select("cipher, public").Where(&Key{Gun: gun, Role: role}).Find(&row)

	if query.RecordNotFound() {
		return "", nil, &ErrNoKey{gun: gun}
	} else if query.Error != nil {
		return "", nil, query.Error
	}

	return row.Cipher, row.Public, nil
}

// SetKey attempts to write a key and returns an error if it already exists for the gun and role
func (db *SQLStorage) SetKey(gun, role, algorithm string, public []byte) error {

	entry := Key{
		Gun:  gun,
		Role: role,
	}

	if !db.Where(&entry).First(&Key{}).RecordNotFound() {
		return &ErrKeyExists{gun: gun, role: role}
	}

	entry.Cipher = algorithm
	entry.Public = public

	return translateOldVersionError(
		db.FirstOrCreate(&Key{}, &entry).Error)
}

// CheckHealth asserts that both required tables are present
func (db *SQLStorage) CheckHealth() error {
	interfaces := []interface {
		TableName() string
	}{&TUFFile{}, &Key{}}

	for _, model := range interfaces {
		tableOk := db.HasTable(model)
		if db.Error != nil {
			return db.Error
		}
		if !tableOk {
			return fmt.Errorf(
				"Cannot access table: %s", model.TableName())
		}
	}
	return nil
}
