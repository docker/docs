// version implementes a versioned store for TUF metadata
package version

import (
	"database/sql"
	"fmt"
)

// VersionDB implements a versioned store using a relational database.
// The database table must look like:
// CREATE TABLE `tuf_files` (
//   `id` INT AUTO_INCREMENT,
//   `qdn` VARCHAR(255) NOT NULL
//   `role` VARCHAR(255) NOT NULL
//   `version` INT
//   `data` LONGBLOB
//   PRIMARY KEY (`id`)
//   UNIQUE INDEX (`qdn`, `role`, `version`)
// ) DEFAULT CHARSET=utf8;
type VersionDB struct {
	sql.DB
}

func NewVersionDB(db *sql.DB) *VersionDB {
	return &VersionDB{
		DB: *db,
	}
}

// Update multiple TUF records in a single transaction.
// Always insert a new row. The unique constraint will ensure there is only ever
func (vdb *VersionDB) UpdateCurrent(qdn, role string, version int, data []byte) error {
	checkStmt := "SELECT 1 FROM `tuf_files` WHERE `qdn`=? AND `role`=? AND `version`=?;"
	insertStmt := "INSERT INTO `tuf_files` (`qdn`, `role`, `version`, `data`) VALUES (?,?,?,?) ;"

	// ensure immediately previous version exists
	row := vdb.QueryRow(checkStmt, qdn, role, version-1)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Attempting to increment version by more than 1 for QDN: %s, role: %s, version: %d", qdn, role, version)
	}

	// attempt to insert. Due to race conditions with the check this could fail.
	// That's OK, we're doing first write wins. The client will be messaged it
	// needs to rebase.
	_, err = vdb.Exec(insertStmt, qdn, role, version, data)
	if err != nil {
		return err
	}
	return nil
}

// Get a specific TUF record
func (vdb *VersionDB) GetCurrent(qdn, tufRole string) (data []byte, err error) {
	stmt := "SELECT `data` FROM `tuf_files` WHERE `qdn`=? AND `role`=? ORDER BY `version` DESC LIMIT 1;"
	rows, err := vdb.Query(stmt, qdn, tufRole) // this should be a QueryRow()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// unique constraint on (qdn, role) will ensure only one row is returned (or none if no match is found)
	if !rows.Next() {
		return nil, nil
	}

	err = rows.Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
