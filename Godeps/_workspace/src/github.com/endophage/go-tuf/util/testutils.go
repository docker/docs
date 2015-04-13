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
	conn, err := sqlite3.Open("/Users/david/gopath/src/github.com/endophage/go-tuf/db/files.db")
	if err != nil {
		panic("can't connect to db")
	}
	return conn
}

func FlushDB(db *sqlite3.Conn) {
	db.Exec("DELETE FROM `filemeta`")
	db.Exec("DELETE FROM `filehashes`")
	db.Exec("DELETE FROM `keys`")

	os.RemoveAll("/tmp/tuf")
}
