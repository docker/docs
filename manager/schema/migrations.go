package schema

import (
	"database/sql"

	"gopkg.in/dancannon/gorethink.v2"
)

var (
	FirstRethinkSchemaVersion  uint64 = 4
	SecondRethinkSchemaVersion uint64 = 5
)

// 0 -> 2.0
var migrateFrom0to4 = Migration{
	FromVersion: 0,
	ToVersion:   FirstRethinkSchemaVersion,
	Migrate: func(db *sql.DB, session *gorethink.Session) error {
		if err := SetupDB(
			session,
			repositoriesTable,
			repositoryTeamAccessTable,
			namespaceTeamAccessTable,
			clientTokenTable); err != nil {
			return err
		}
		// TODO add migration
		return nil
	},
	Rollback: func(db *sql.DB, session *gorethink.Session) error {
		// TODO add rollback?
		return nil
	},
}

// DTR 1.4 -> 2.0
var migrateFrom3to4 = Migration{
	FromVersion: 3,
	ToVersion:   FirstRethinkSchemaVersion, // 4
	Migrate: func(db *sql.DB, session *gorethink.Session) error {
		if err := migrateFrom0to4.Migrate(db, session); err != nil {
			return err
		}

		// TODO add migration
		return nil
	},
	Rollback: func(db *sql.DB, session *gorethink.Session) error {
		// Remove the added tables.
		return nil
	},
}

// 2.0 -> 2.1 - initial eventstreams tables
var migrateFrom4to5 = Migration{
	FromVersion: FirstRethinkSchemaVersion,
	ToVersion:   SecondRethinkSchemaVersion, // 5
	Migrate: func(db *sql.DB, session *gorethink.Session) error {
		if err := SetupDB(session, eventsTable); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(db *sql.DB, session *gorethink.Session) error {
		return nil
	},
}

// 2.1 - manifests and tags
var migrateFrom5to6 = Migration{
	FromVersion: 5,
	ToVersion:   6,
	Migrate: func(db *sql.DB, session *gorethink.Session) error {
		if err := SetupDB(
			session,
			tagsTable,
			manifestsTable); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(db *sql.DB, session *gorethink.Session) error {
		return nil
	},
}

func init() {
	registerMigration(migrateFrom0to4)
	registerMigration(migrateFrom3to4)
	registerMigration(migrateFrom4to5)
	registerMigration(migrateFrom5to6)
}
