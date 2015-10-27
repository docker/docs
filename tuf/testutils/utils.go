package testutils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/docker/notary/tuf/data"
	// need to initialize sqlite for tests
	_ "github.com/mattn/go-sqlite3"
)

var counter = 1

// SampleMeta returns a static, fake (and invalid) FileMeta object
func SampleMeta() data.FileMeta {
	meta := data.FileMeta{
		Length: 1,
		Hashes: data.Hashes{
			"sha256": []byte{0x01, 0x02},
			"sha512": []byte{0x03, 0x04},
		},
	}
	return meta
}

// GetSqliteDB creates and initializes a sqlite db
func GetSqliteDB() *sql.DB {
	os.Mkdir("/tmp/sqlite", 0755)
	conn, err := sql.Open("sqlite3", fmt.Sprintf("/tmp/sqlite/file%d.db", counter))
	if err != nil {
		panic("can't connect to db")
	}
	counter++
	tx, err := conn.Begin()
	if err != nil {
		panic("can't begin db transaction")
	}
	tx.Exec("CREATE TABLE keys (id int auto_increment, namespace varchar(255) not null, role varchar(255) not null, key text not null, primary key (id));")
	tx.Exec("CREATE TABLE filehashes(namespace varchar(255) not null, path varchar(255) not null, alg varchar(10) not null, hash varchar(128) not null, primary key (namespace, path, alg));")
	tx.Exec("CREATE TABLE filemeta(namespace varchar(255) not null, path varchar(255) not null, size int not null, custom text default null, primary key (namespace, path));")
	tx.Commit()
	return conn
}

// FlushDB deletes a sqliteDB
func FlushDB(db *sql.DB) {
	tx, _ := db.Begin()
	tx.Exec("DELETE FROM `filemeta`")
	tx.Exec("DELETE FROM `filehashes`")
	tx.Exec("DELETE FROM `keys`")
	tx.Commit()
	os.RemoveAll("/tmp/tuf")
}
