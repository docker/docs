package storage

import (
	"github.com/docker/notary/storage/rethinkdb"
)

var (
	tufFiles = rethinkdb.Table{
		Name:       RDBTUFFile{}.TableName(),
		PrimaryKey: []string{"gun", "role", "version"},
		SecondaryIndexes: map[string][]string{
			"sha256": nil,
		},
	}

	keys = rethinkdb.Table{
		Name:             RDBKey{}.TableName(),
		PrimaryKey:       "id",
		SecondaryIndexes: nil,
	}
)
