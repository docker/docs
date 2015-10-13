package storage

import (
	"github.com/Sirupsen/logrus"
	"github.com/endophage/gotuf/data"
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

	return db.Create(&TUFFile{
		Gun:     gun,
		Role:    update.Role,
		Version: update.Version,
		Data:    update.Data,
	}).Error
}

// UpdateMany atomically updates many TUF records in a single transaction
func (db *SQLStorage) UpdateMany(gun string, updates []MetaUpdate) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	rollback := func(err error) error {
		query := tx.Rollback()
		if query.Error != nil {
			logrus.Panic("Failed on Tx rollback with error: ", err.Error())
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
		query = tx.Where(map[string]interface{}{
			"gun":     gun,
			"role":    update.Role,
			"version": update.Version,
		}).Attrs("data", update.Data).FirstOrCreate(&row)

		if query.Error != nil {
			return rollback(query.Error)
		}
		// it's previously been added, which means it's a duplicate entry
		// in the same transaction
		if added[row.ID] {
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

	if q.RecordNotFound() {
		return nil, &ErrNotFound{}
	} else if q.Error != nil {
		return nil, q.Error
	}
	return row.Data, nil
}

// Delete deletes all the records for a specific GUN
func (db *SQLStorage) Delete(gun string) error {
	return db.Where(&TUFFile{Gun: gun}).Delete(TUFFile{}).Error
}

// GetTimestampKey returns the timestamps Public Key data
func (db *SQLStorage) GetTimestampKey(gun string) (algorithm data.KeyAlgorithm, public []byte, err error) {
	logrus.Debug("retrieving timestamp key for ", gun)

	var row TimestampKey
	query := db.Select("cipher, public").Where(&TimestampKey{Gun: gun}).Find(&row)

	if query.RecordNotFound() {
		return "", nil, &ErrNoKey{gun: gun}
	} else if query.Error != nil {
		return "", nil, query.Error
	}

	return data.KeyAlgorithm(row.Cipher), row.Public, nil
}

// SetTimestampKey attempts to write a TimeStamp key and returns an error if it already exists
func (db *SQLStorage) SetTimestampKey(gun string, algorithm data.KeyAlgorithm, public []byte) error {

	entry := TimestampKey{
		Gun:    gun,
		Cipher: string(algorithm),
		Public: public,
	}

	if !db.Where(&entry).First(&TimestampKey{}).RecordNotFound() {
		return &ErrTimestampKeyExists{gun: gun}
	}

	return db.FirstOrCreate(&TimestampKey{}, &entry).Error
}
