package storage

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySQLUpdateCurrent(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)
	update := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	sqlmock.ExpectExec("INSERT INTO `tuf_files` \\(`gun`, `role`, `version`, `data`\\) VALUES \\(\\?,\\?,\\?,\\?\\) WHERE \\(SELECT count\\(\\*\\) FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? AND `version`>=\\?\\) = 0").WithArgs(
		"testGUN",
		update.Role,
		update.Version,
		update.Data,
		"testGUN",
		update.Role,
		update.Version,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err = s.UpdateCurrent(
		"testGUN",
		update,
	)
	assert.Nil(t, err, "UpdateCurrent errored unexpectedly: %v", err)

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLUpdateCurrentError(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)
	update := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	sqlmock.ExpectExec("INSERT INTO `tuf_files` \\(`gun`, `role`, `version`, `data`\\) VALUES \\(\\?,\\?,\\?,\\?\\) WHERE \\(SELECT count\\(\\*\\) FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? AND `version`>=\\?\\) = 0").WithArgs(
		"testGUN",
		update.Role,
		update.Version,
		update.Data,
		"testGUN",
		update.Role,
		update.Version,
	).WillReturnError(
		&mysql.MySQLError{
			Number:  1022,
			Message: "Duplicate key error",
		},
	)

	err = s.UpdateCurrent(
		"testGUN",
		update,
	)
	assert.NotNil(t, err, "Error should not be nil")
	assert.IsType(t, &ErrOldVersion{}, err, "Expected ErrOldVersion error type")

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLUpdateMany(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)
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
	// start transation
	sqlmock.ExpectBegin()

	// insert first update
	sqlmock.ExpectExec("INSERT INTO `tuf_files` \\(`gun`, `role`, `version`, `data`\\) VALUES \\(\\?,\\?,\\?,\\?\\) WHERE \\(SELECT count\\(\\*\\) FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? AND `version`>=\\?\\) = 0").WithArgs(
		"testGUN",
		update1.Role,
		update1.Version,
		update1.Data,
		"testGUN",
		update1.Role,
		update1.Version,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	// insert second update
	sqlmock.ExpectExec("INSERT INTO `tuf_files` \\(`gun`, `role`, `version`, `data`\\) VALUES \\(\\?,\\?,\\?,\\?\\) WHERE \\(SELECT count\\(\\*\\) FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? AND `version`>=\\?\\) = 0").WithArgs(
		"testGUN",
		update2.Role,
		update2.Version,
		update2.Data,
		"testGUN",
		update2.Role,
		update2.Version,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// expect commit
	sqlmock.ExpectCommit()

	err = s.UpdateMany(
		"testGUN",
		[]MetaUpdate{update1, update2},
	)
	assert.Nil(t, err, "UpdateMany errored unexpectedly: %v", err)

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLUpdateManyRollback(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)
	update1 := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	execError := mysql.MySQLError{}
	// start transation
	sqlmock.ExpectBegin()

	// insert first update
	sqlmock.ExpectExec("INSERT INTO `tuf_files` \\(`gun`, `role`, `version`, `data`\\) VALUES \\(\\?,\\?,\\?,\\?\\) WHERE \\(SELECT count\\(\\*\\) FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? AND `version`>=\\?\\) = 0").WithArgs(
		"testGUN",
		update1.Role,
		update1.Version,
		update1.Data,
		"testGUN",
		update1.Role,
		update1.Version,
	).WillReturnError(&execError)

	// expect commit
	sqlmock.ExpectRollback()

	err = s.UpdateMany(
		"testGUN",
		[]MetaUpdate{update1},
	)
	assert.IsType(t, &execError, err, "UpdateMany returned wrong error type")

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLUpdateManyDuplicate(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)
	update1 := MetaUpdate{
		Role:    "root",
		Version: 0,
		Data:    []byte("1"),
	}
	execError := mysql.MySQLError{Number: 1022}
	// start transation
	sqlmock.ExpectBegin()

	// insert first update
	sqlmock.ExpectExec("INSERT INTO `tuf_files` \\(`gun`, `role`, `version`, `data`\\) VALUES \\(\\?,\\?,\\?,\\?\\) WHERE \\(SELECT count\\(\\*\\) FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? AND `version`>=\\?\\) = 0").WithArgs(
		"testGUN",
		update1.Role,
		update1.Version,
		update1.Data,
		"testGUN",
		update1.Role,
		update1.Version,
	).WillReturnError(&execError)

	// expect commit
	sqlmock.ExpectRollback()

	err = s.UpdateMany(
		"testGUN",
		[]MetaUpdate{update1},
	)
	assert.IsType(t, &ErrOldVersion{}, err, "UpdateMany returned wrong error type")

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLGetCurrent(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)

	sqlmock.ExpectQuery(
		"SELECT `data` FROM `tuf_files` WHERE `gun`=\\? AND `role`=\\? ORDER BY `version` DESC LIMIT 1;",
	).WithArgs("testGUN", "root").WillReturnRows(
		sqlmock.RowsFromCSVString(
			[]string{"data"},
			"1",
		),
	)

	byt, err := s.GetCurrent("testGUN", "root")
	assert.Nil(t, err, "Expected nil error from GetCurrent")
	assert.Equal(t, []byte("1"), byt, "Returned data was no correct")

	// TODO(endophage): these two lines are breaking because there
	//                  seems to be some problem with go-sqlmock
	//err = db.Close()
	//assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLDelete(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)

	sqlmock.ExpectExec(
		"DELETE FROM `tuf_files` WHERE `gun`=\\?;",
	).WithArgs("testGUN").WillReturnResult(sqlmock.NewResult(0, 1))

	err = s.Delete("testGUN")
	assert.Nil(t, err, "Expected nil error from Delete")

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLGetTimestampKeyNoKey(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)

	sqlmock.ExpectQuery(
		"SELECT `cipher`, `public` FROM `timestamp_keys` WHERE `gun`=\\?;",
	).WithArgs("testGUN").WillReturnError(sql.ErrNoRows)

	_, _, err = s.GetTimestampKey("testGUN")
	assert.IsType(t, &ErrNoKey{}, err, "Expected ErrNoKey from GetTimestampKey")

	//err = db.Close()
	//assert.Nil(t, err, "Expectation not met: %v", err)
}

func TestMySQLSetTimestampKeyExists(t *testing.T) {
	db, err := sqlmock.New()
	assert.Nil(t, err, "Could not initialize mock DB")
	s := NewMySQLStorage(db)

	sqlmock.ExpectExec(
		"INSERT INTO `timestamp_keys` \\(`gun`, `cipher`, `public`\\) VALUES \\(\\?,\\?,\\?\\);",
	).WithArgs(
		"testGUN",
		"testCipher",
		[]byte("1"),
	).WillReturnError(
		&mysql.MySQLError{Number: 1022},
	)

	err = s.SetTimestampKey("testGUN", "testCipher", []byte("1"))
	assert.IsType(t, &ErrTimestampKeyExists{}, err, "Expected ErrTimestampKeyExists from SetTimestampKey")

	err = db.Close()
	assert.Nil(t, err, "Expectation not met: %v", err)
}
