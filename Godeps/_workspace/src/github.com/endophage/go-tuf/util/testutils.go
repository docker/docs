package util

import (
	"database/sql"
	"fmt"
	"os"

	_ "code.google.com/p/gosqlite/sqlite3"
	"github.com/endophage/go-tuf/data"
)

var counter int = 1

func SampleMeta() data.FileMeta {
	meta := data.FileMeta{
		Length: 1,
		Hashes: data.Hashes{
			"sha256": data.HexBytes{0x01, 0x02},
			"sha512": data.HexBytes{0x03, 0x04},
		},
	}
	return meta
}

func GetSqliteDB() *sql.DB {
	conn, err := sql.Open("sqlite3", fmt.Sprintf("/tmp/sqlite/file%d.db", counter))
	if err != nil {
		panic("can't connect to db")
	}
	counter++
	tx, _ := conn.Begin()
	tx.Exec("CREATE TABLE keys (id int auto_increment, namespace varchar(255) not null, role varchar(255) not null, key text not null, primary key (id));")
	tx.Exec("CREATE TABLE filehashes(namespace varchar(255) not null, path varchar(255) not null, alg varchar(10) not null, hash varchar(128) not null, primary key (namespace, path, alg));")
	tx.Exec("CREATE TABLE filemeta(namespace varchar(255) not null, path varchar(255) not null, size int not null, custom text default null, primary key (namespace, path));")
	tx.Commit()
	return conn
}

func FlushDB(db *sql.DB) {
	tx, _ := db.Begin()
	tx.Exec("DELETE FROM `filemeta`")
	tx.Exec("DELETE FROM `filehashes`")
	tx.Exec("DELETE FROM `keys`")
	tx.Commit()
	os.RemoveAll("/tmp/tuf")
}
