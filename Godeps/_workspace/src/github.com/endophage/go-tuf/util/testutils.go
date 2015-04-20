package util

import (
	"os"

	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/endophage/go-tuf/data"
)

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

func GetSqliteDB() *sqlite3.Conn {
	conn, err := sqlite3.Open("")
	if err != nil {
		panic("can't connect to db")
	}
	conn.Exec("CREATE TABLE keys (id int auto_increment, namespace varchar(255) not null, role varchar(255) not null, key text not null, primary key (id));")
	conn.Exec("CREATE TABLE filehashes(namespace varchar(255) not null, path varchar(255) not null, alg varchar(10) not null, hash varchar(128) not null, primary key (namespace, path, alg));")
	conn.Exec("CREATE TABLE filemeta(namespace varchar(255) not null, path varchar(255) not null, size int not null, custom text default null, primary key (namespace, path));")
	conn.Commit()
	return conn
}

func FlushDB(db *sqlite3.Conn) {
	db.Exec("DELETE FROM `filemeta`")
	db.Exec("DELETE FROM `filehashes`")
	db.Exec("DELETE FROM `keys`")

	os.RemoveAll("/tmp/tuf")
}
