package storage

import (
	"github.com/docker/notary/storage/rethinkdb"
)

// These consts are the index names we've defined for RethinkDB
const (
	rdbSha256Idx        = "sha256"
	rdbGunRoleIdx       = "gun_role"
	rdbGunRoleSha256Idx = "gun_role_sha256"
)

var (
	tufFiles = rethinkdb.Table{
		Name:       RDBTUFFile{}.TableName(),
		PrimaryKey: []string{"gun", "role", "version"},
		SecondaryIndexes: map[string][]string{
			rdbSha256Idx:        nil,
			"gun":               nil,
			rdbGunRoleIdx:       {"gun", "role"},
			rdbGunRoleSha256Idx: {"gun", "role", "sha256"},
		},
	}

	keys = rethinkdb.Table{
		Name:       RDBKey{}.TableName(),
		PrimaryKey: "id",
		SecondaryIndexes: map[string][]string{
			rdbGunRoleIdx: {"gun", "role"},
		},
	}
)
