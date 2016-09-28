package schema

import (
	"archive/tar"
	"database/sql"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/manager/schema/interfaces"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/notary/passphrase"
	notaryserver "github.com/docker/notary/server/storage"
	notarysigner "github.com/docker/notary/signer/keydbstore"
	notaryrethink "github.com/docker/notary/storage/rethinkdb"
	"gopkg.in/dancannon/gorethink.v2"
)

const schemaVersionKey = "schema.version"
const WaitTime = 600

type Manager interface {
	Initialize() error
	Migrate() error
	GetPropertyManager() interfaces.PropertyManager
	SetReplication(replFactor int) error
	DumpAllTables(tw *tar.Writer, backupVersion string) error
	RestoreTableDocument(tableName string, document []byte) error
	DecommissionServer(name string) error
	WaitForReady() error
}

type Migration struct {
	FromVersion uint64
	ToVersion   uint64
	Migrate     func(db *sql.DB, session *gorethink.Session) error
	Rollback    func(db *sql.DB, session *gorethink.Session) error
}

var highestKnownVersion uint64
var migrations map[uint64]Migration

// NOTE: These tablenames should be unique, since we need to look them up for
// backup and restore currently.  The uniqueness will be enforced during the
// manager initialization.  Ideally, the table backups should be namespaced by
// DB, but the backup format is probably going to change soon in the future
// anyway, and we also always need to be able to read old backups.

// DTRTables are the DB tables used by DTR
var DTRTables = map[string]table{
	repositoriesTable.name:         repositoriesTable,
	propertiesTable.name:           propertiesTable,
	repositoryTeamAccessTable.name: repositoryTeamAccessTable,
	namespaceTeamAccessTable.name:  namespaceTeamAccessTable,
	clientTokenTable.name:          clientTokenTable,
	tagsTable.name:                 tagsTable,
	manifestsTable.name:            manifestsTable,
	eventsTable.name:               eventsTable,
}

// NotaryTables are the DB tables used by Notary, which are across 2 DBs and
// in a separate DB from DTR
var NotaryTables = map[string]notaryrethink.Table{
	notaryserver.TUFFilesRethinkTable.Name:    notaryserver.TUFFilesRethinkTable,
	notaryserver.PubKeysRethinkTable.Name:     notaryserver.PubKeysRethinkTable,
	notarysigner.PrivateKeysRethinkTable.Name: notarysigner.PrivateKeysRethinkTable,
}

func registerMigration(m Migration) {
	if migrations == nil {
		migrations = make(map[uint64]Migration)
	}
	if _, ok := migrations[m.FromVersion]; ok {
		panic(fmt.Sprintf("Duplicate migration from version %d!", m.FromVersion))
	}
	migrations[m.FromVersion] = m
	if m.ToVersion > highestKnownVersion {
		highestKnownVersion = m.ToVersion
	}
}

type manager struct {
	// why is there bothe a session and a conn func???
	rethinkConnFunc func() (*gorethink.Session, error)
	db              *sql.DB
	rethinkdb       *gorethink.Session
	// TODO: switch to using rethinkutil.Table
	allTables map[string]table
}

// TODO: implement rethinkutil.Manager
// XXX: don't ask for a connection function - use a session :/
func NewManager(rethinkConnFunc func() (*gorethink.Session, error)) Manager {
	m := &manager{rethinkConnFunc: rethinkConnFunc}

	m.allTables = make(map[string]table)
	for tbleName, tbl := range DTRTables {
		if _, ok := m.allTables[tbleName]; ok {
			panic(fmt.Errorf("Table %s in DB %s does not have a unique table name.",
				tbl.name, tbl.db))
		}
		m.allTables[tbleName] = tbl
	}
	for tbleName, tbl := range NotaryTables {
		dbName := deploy.NotaryServerDBName
		if tbleName == notarysigner.PrivateKeysRethinkTable.Name {
			dbName = deploy.NotarySignerDBName
		}
		if _, ok := m.allTables[tbleName]; ok {
			panic(fmt.Errorf("Table %s in DB %s does not have a unique table name.",
				tbl.Name, dbName))
		}
		m.allTables[tbleName] = table{
			db:               dbName,
			name:             tbl.Name,
			primaryKey:       fmt.Sprintf("%s", tbl.PrimaryKey),
			secondaryIndexes: tbl.SecondaryIndexes,
		}
	}

	return m
}

func (m *manager) Initialize() error {
	var schemaVersion uint64
	if schemaVersionP, err := m.initializeRethink(); err != nil {
		return err
	} else if schemaVersionP != nil {
		schemaVersion = *schemaVersionP
	}

	if schemaVersion > highestKnownVersion {
		log.WithFields(log.Fields{
			"version":             schemaVersion,
			"highestKnownVersion": highestKnownVersion,
		}).Error("Schema version is too high")
		return fmt.Errorf("Schema version is too high")
	}

	if schemaVersion == 0 {
		pm := NewRethinkPropertyManager(m.rethinkdb)
		// start the schema version at 0
		if err := pm.Set(schemaVersionKey, "0"); err != nil {
			return err
		}
	}

	if err := m.BootstrapNotaryTables(); err != nil {
		return err
	}

	return nil
}

func (m *manager) GetPropertyManager() interfaces.PropertyManager {
	return NewRethinkPropertyManager(m.rethinkdb)
}

func (m *manager) initializeRethink() (*uint64, error) {
	if m.rethinkdb == nil {
		session, err := m.rethinkConnFunc()
		if err != nil {
			return nil, err
		}
		m.rethinkdb = session
	}
	if v, err := m.schemaVersion(); err != nil {
		if err := makeDB(m.rethinkdb, deploy.DTRDBName); err != nil {
			return nil, err
		} else {
			err := propertiesTable.create(m.rethinkdb, 1)
			return nil, err
		}
	} else {
		return &v, nil
	}
}

func (m *manager) Migrate() error {
	originalVersion, err := m.schemaVersion()
	if err != nil {
		return err
	}

	if originalVersion > highestKnownVersion {
		log.WithFields(log.Fields{
			"version":             originalVersion,
			"highestKnownVersion": highestKnownVersion,
		}).Error("Schema version is too high")
		return fmt.Errorf("Schema version is too high")
	}

	finalVersion, err := m.runMigrations(originalVersion)
	if err != nil {
		log.WithFields(log.Fields{
			"originalVersion": originalVersion,
			"finalVersion":    finalVersion,
			"error":           err,
		}).Error("Failed to migrate database schema")
		return err
	}

	if finalVersion == originalVersion {
		log.WithField("version", finalVersion).Info("Already at latest database schema version")
	} else {
		pm := NewRethinkPropertyManager(m.rethinkdb)
		if err := pm.Set(schemaVersionKey, strconv.FormatUint(finalVersion, 10)); err != nil {
			log.WithFields(log.Fields{
				"originalVersion": originalVersion,
				"finalVersion":    finalVersion,
				"error":           err,
			}).Error("Failed to update database version")
			return err
		}
		log.Infof("Migrated database from version %d to %d", originalVersion, finalVersion)
	}
	return nil
}

func (m *manager) schemaVersion() (uint64, error) {
	pm := NewRethinkPropertyManager(m.rethinkdb)
	if version, err := pm.Get(schemaVersionKey); err == interfaces.ErrPropertyNotSet {
		log.Debugf("No %s property found in the properties table", schemaVersionKey)
		return 0, fmt.Errorf("unknown schema version")
	} else if err != nil {
		return 0, err
	} else {
		return strconv.ParseUint(version, 10, 64)
	}
}

func (m *manager) runMigrations(originalVersion uint64) (uint64, error) {
	version := originalVersion
	for {
		migration, ok := migrations[version]
		if !ok {
			return version, nil
		}
		log.WithFields(log.Fields{
			"fromVersion": version,
			"toVersion":   migration.ToVersion,
		}).Info("Migrating database schema")
		if err := migration.Migrate(m.db, m.rethinkdb); err != nil {
			log.WithFields(log.Fields{
				"originalVersion": originalVersion,
				"currentVersion":  version,
				"toVersion":       migration.ToVersion,
				"error":           err,
			}).Error("Failed to migrate database schema")
			return version, err
		}
		version = migration.ToVersion
	}
}

func ParseConfig(config map[string]interface{}) ([]string, string, error) {
	// This hurt. I died a little on the inside.
	shards, ok := config["shards"].([]interface{})
	if !ok {
		return nil, "", fmt.Errorf("Failed shards array cast")
	}
	shard, ok := shards[0].(map[string]interface{})
	if !ok {
		return nil, "", fmt.Errorf("Failed shard cast")
	}
	primaryReplica, ok := shard["primary_replica"].(string)
	if !ok {
		return nil, "", fmt.Errorf("Failed primary_replica cast")
	}
	replicas, ok := shard["replicas"].([]interface{})
	if !ok {
		return nil, "", fmt.Errorf("Failed replicas cast")
	}
	replicasStr := []string{}
	for _, replica := range replicas {
		replica, ok := replica.(string)
		if !ok {
			return nil, "", fmt.Errorf("Failed replica cast")
		}
		replicasStr = append(replicasStr, replica)
	}
	return replicasStr, primaryReplica, nil
}

func (m *manager) WaitForReady() error {
	// do we need to call this everywhere? :/
	if _, err := m.initializeRethink(); err != nil {
		return fmt.Errorf("Couldn't initialize rethink: %s", err)
	}
	waitOpts := gorethink.WaitOpts{
		WaitFor: "all_replicas_ready",
		Timeout: WaitTime,
	}

	dbNames, err := m.dbList()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve database list from rethink: %s", err)
	}

	for _, dbName := range dbNames {
		// database rethinkdb is special, system tables in it are always available
		// and don't to be waited on
		if dbName == "rethinkdb" {
			continue
		}

		tableNames, err := m.tableList(dbName)
		if err != nil {
			return fmt.Errorf("Couldn't retrieve table list on db %s from rethink: %s", dbName, err)
		}

		// Wait for all existing tables rather than what we know about. During
		// reconfigure step in upgrade, we only care that rethinkdb is ready and up
		for _, tableName := range tableNames {
			_, err := gorethink.DB(dbName).Table(tableName).Wait(waitOpts).Run(m.rethinkdb)
			if err != nil {
				return fmt.Errorf("Failed to wait for replica to settle on table %s in db %s: %s", tableName, dbName, err)
			}
		}
	}

	return nil
}

// TODO: consider moving to rethinkutil
func (m *manager) DecommissionServer(name string) error {
	if _, err := m.initializeRethink(); err != nil {
		return fmt.Errorf("Couldn't initialize rethink: %s", err)
	}
	_, err := gorethink.DB("rethinkdb").Table("server_config").Filter(map[string]interface{}{
		"name": name,
	}).Update(map[string]interface{}{
		"tags": []string{},
	}).RunWrite(m.rethinkdb)
	if err != nil {
		return err
	}
	return m.WaitForReady()
}

func (m *manager) SetReplication(replFactor int) error {
	if _, err := m.initializeRethink(); err != nil {
		return fmt.Errorf("Couldn't initialize rethink: %s", err)
	}
	opts := gorethink.ReconfigureOpts{
		Shards:   1,
		Replicas: replFactor,
	}
	for _, tbl := range m.allTables {
		log.Infof("setting replication factor for table %s to %d", tbl.name, replFactor)
		_, err := tbl.Term().Reconfigure(opts).Run(m.rethinkdb)
		if err != nil {
			return fmt.Errorf("There was an error changing the replication factor for %s: %s", tbl.name, err)
		}
	}
	return m.WaitForReady()
}

func (m *manager) BootstrapNotaryTables() error {
	// For now, share the same database as DTR but use separate tables
	// TODO(cyli): don't use hardcoded server/signer passwords
	log.Info("Bootstrapping notary tables")
	notaryServerRethink := notaryserver.NewRethinkDBStorage(deploy.NotaryServerDBName, defaultconfigs.DefaultNotaryServerConfig.Storage.Username, defaultconfigs.DefaultNotaryServerConfig.Storage.Password, m.rethinkdb)
	if err := notaryServerRethink.Bootstrap(); err != nil {
		return err
	}
	// TODO(riyazdf): handle configurable default alias
	notarySignerRethink := notarysigner.NewRethinkDBKeyStore(deploy.NotarySignerDBName, defaultconfigs.DefaultNotarySignerConfig.Storage.Username, defaultconfigs.DefaultNotarySignerConfig.Storage.Password, passphrase.PromptRetriever(), "timestamp_1", m.rethinkdb)
	if err := notarySignerRethink.Bootstrap(); err != nil {
		return err
	}
	return nil
}

func (m *manager) DumpAllTables(tw *tar.Writer, backupVersion string) error {
	// it's safe to assume even in extreme cases an instance has fewer than
	// ~4.3 billion repos
	var row_index uint64
	for _, tbl := range m.allTables {
		row_index = 0
		var documents []interface{}
		query := tbl.Term()
		res, err := query.Run(m.rethinkdb)
		if err != nil {
			log.Errorf("Couldn't run the query: %s", err)
			return err
		}

		err = res.All(&documents)
		if err != nil {
			log.Errorf("Couldn't parse rethinkdb's resposne: %s", err)
			return err
		}

		tableName := fmt.Sprintf("%s/%s/%s/", backupVersion, deploy.RethinkdbDirectory, tbl.name)
		tableFolderHeader := &tar.Header{
			Name:     tableName,
			Typeflag: tar.TypeDir,
			Mode:     0700,
		}
		if err := tw.WriteHeader(tableFolderHeader); err != nil {
			log.Errorf("Couldn't write header to indicate rethinkdb table directory: %s", err)
			return err
		}

		// convert the map into bytes
		for _, doc := range documents {
			documentBytes, err := dtrutil.GetBytes(doc)
			if err != nil {
				log.Errorf("Couldn't convert rethink table's content into bytes: %s", err)
				return err
			}
			// the slash (/) is already appended to the tableName
			name := fmt.Sprintf("%s%d", tableName, row_index)
			row_index += 1
			if err = dtrutil.AddBytesToTar(tw, documentBytes, name); err != nil {
				log.Errorf("Couldn't write rethink table's data to tar archive: %s", err)
				return err
			}
		}
	}
	return nil
}

func (m *manager) RestoreTableDocument(tableName string, document []byte) error {
	// first convert the document back into something
	var (
		doc interface{}
		err error
	)
	if notaryTable, ok := NotaryTables[tableName]; ok {
		doc, err = notaryTable.JSONUnmarshaller(document)
	} else {
		doc, err = dtrutil.GetMap(document)
	}

	if err != nil {
		log.Errorf("Can't convert bytes to document: %s", err)
		return err
	}

	// check to see if the tableName is registered in the map of valid tables -
	// if not, for now log a warning and assume that it's a DTR table
	dbName := deploy.DTRDBName
	if tbl, ok := m.allTables[tableName]; ok {
		dbName = tbl.db
	} else {
		log.Warnf("Unrecognized table %s in backup.  Assuming it should go in %s", tableName, dbName)
	}
	log.Debugf("Restoring document from %s %s %s", dbName, tableName)
	query := gorethink.DB(dbName).Table(tableName).Insert(doc)
	_, err = query.Run(m.rethinkdb)
	if err != nil {
		log.Errorf("Can't restore rethinkdb row: %s", err)
		return err
	}
	return nil
}

func (m *manager) dbList() ([]string, error) {
	var dbNames []string
	res, err := gorethink.DBList().Run(m.rethinkdb)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if err := res.All(&dbNames); err != nil {
		return nil, err
	}

	return dbNames, nil
}

func (m *manager) tableList(dbName string) ([]string, error) {
	var tables []string
	res, err := gorethink.DB(dbName).TableList().Run(m.rethinkdb)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if err := res.All(&tables); err != nil {
		return nil, err
	}

	return tables, nil
}
