package storage

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	"github.com/endophage/gotuf/data"
	"github.com/go-sql-driver/mysql"
)

// MySQLStorage implements a versioned store using a relational database.
// The database table must look like:
// CREATE TABLE `tuf_files` (
//   `id` INT AUTO_INCREMENT,
//   `gun` VARCHAR(255) NOT NULL
//   `role` VARCHAR(255) NOT NULL
//   `version` INT
//   `data` LONGBLOB
//   PRIMARY KEY (`id`)
//   UNIQUE INDEX (`gun`, `role`, `version`)
// ) DEFAULT CHARSET=utf8;
//
// CREATE TABLE `timestamp_keys` (
//   `gun` VARCHAR(255),
//   `cipher` VARCHAR(30),
//   `public` BLOB NOT NULL,
// ) DEFAULT CHARSET=utf8;
type MySQLStorage struct {
	sql.DB
}

// NewMySQLStorage is a convenience method to create a MySQLStorage
func NewMySQLStorage(db *sql.DB) *MySQLStorage {
	return &MySQLStorage{
		DB: *db,
	}
}

// UpdateCurrent updates multiple TUF records in a single transaction.
// Always insert a new row. The unique constraint will ensure there is only ever
func (db *MySQLStorage) UpdateCurrent(gun string, update MetaUpdate) error {
	insertStmt := "INSERT INTO `tuf_files` (`gun`, `role`, `version`, `data`) VALUES (?,?,?,?) WHERE (SELECT count(*) FROM `tuf_files` WHERE `gun`=? AND `role`=? AND `version`>=?) = 0"

	// attempt to insert. Due to race conditions with the check this could fail.
	// That's OK, we're doing first write wins. The client will be messaged it
	// needs to rebase.
	_, err := db.Exec(insertStmt, gun, update.Role, update.Version, update.Data, gun, update.Role, update.Version)
	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1022 { // duplicate key error
				return &ErrOldVersion{}
			}
		}
		// need to check error type for duplicate key exception
		// and return ErrOldVersion if duplicate
		return err
	}
	return nil
}

// UpdateMany atomically updates many TUF records in a single transaction
func (db *MySQLStorage) UpdateMany(gun string, updates []MetaUpdate) error {
	insertStmt := "INSERT INTO `tuf_files` (`gun`, `role`, `version`, `data`) VALUES (?,?,?,?) WHERE (SELECT count(*) FROM `tuf_files` WHERE `gun`=? AND `role`=? AND `version`>=?) = 0;"

	tx, err := db.Begin()
	for _, u := range updates {
		// attempt to insert. Due to race conditions with the check this could fail.
		// That's OK, we're doing first write wins. The client will be messaged it
		// needs to rebase.
		_, err = tx.Exec(insertStmt, gun, u.Role, u.Version, u.Data, gun, u.Role, u.Version)
		if err != nil {
			// need to check error type for duplicate key exception
			// and return ErrOldVersion if duplicate
			rbErr := tx.Rollback()
			if rbErr != nil {
				logrus.Panic("Failed on Tx rollback with error: ", err.Error())
			}
			if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1022 { // duplicate key error
				return &ErrOldVersion{}
			}
			return err
		}
	}
	return tx.Commit()
}

// GetCurrent gets a specific TUF record
func (db *MySQLStorage) GetCurrent(gun, tufRole string) (data []byte, err error) {
	stmt := "SELECT `data` FROM `tuf_files` WHERE `gun`=? AND `role`=? ORDER BY `version` DESC LIMIT 1;"
	rows, err := db.Query(stmt, gun, tufRole) // this should be a QueryRow()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// unique constraint on (gun, role) will ensure only one row is returned (or none if no match is found)
	if !rows.Next() {
		return nil, &ErrNotFound{}
	}

	err = rows.Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Delete deletes all the records for a specific GUN
func (db *MySQLStorage) Delete(gun string) error {
	stmt := "DELETE FROM `tuf_files` WHERE `gun`=?;"
	_, err := db.Exec(stmt, gun)
	return err
}

// GetTimestampKey returns the timestamps Public Key data
func (db *MySQLStorage) GetTimestampKey(gun string) (algorithm data.KeyAlgorithm, public []byte, err error) {
	stmt := "SELECT `cipher`, `public` FROM `timestamp_keys` WHERE `gun`=?;"
	row := db.QueryRow(stmt, gun)

	var cipher string
	err = row.Scan(&cipher, &public)
	if err == sql.ErrNoRows {
		return "", nil, &ErrNoKey{gun: gun}
	} else if err != nil {
		return "", nil, err
	}

	return data.KeyAlgorithm(cipher), public, err
}

// SetTimestampKey attempts to write a TimeStamp key and returns an error if it already exists
func (db *MySQLStorage) SetTimestampKey(gun string, algorithm data.KeyAlgorithm, public []byte) error {
	stmt := "INSERT INTO `timestamp_keys` (`gun`, `cipher`, `public`) VALUES (?,?,?);"
	_, err := db.Exec(stmt, gun, string(algorithm), public)
	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1022 {
			return &ErrTimestampKeyExists{gun: gun}
		}
		return err
	}
	return nil
}
