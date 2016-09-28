package rethinkutil

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type Manager interface {
	DB() string
	Session() *rethink.Session
	Tables() []Table
	SetupDB(expectedReplicaCount uint) error
	ScaleDB(uint, bool) error
}

type GenericManager struct {
	db      string
	session *rethink.Session
	tables  []Table
}

var _ Manager = &GenericManager{}

func NewGenericManager(db string, session *rethink.Session, tables []Table) Manager {
	return &GenericManager{
		db:      db,
		session: session,
		tables:  tables,
	}
}

func (m *GenericManager) DB() string {
	return m.db
}

func (m *GenericManager) Session() *rethink.Session {
	return m.session
}

func (m *GenericManager) Tables() []Table {
	return m.tables
}

func makeDB(name string, session *rethink.Session) error {
	_, err := rethink.DBCreate(name).RunWrite(session)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Debug("Database Already Exists")
			return nil
		}

		return err
	}

	resp, err := rethink.DB(name).Wait().Run(session)
	if resp != nil {
		resp.Close()
	}

	return err
}

func getNumberOfServers(session *rethink.Session) (uint, error) {
	cursor, err := rethink.DB("rethinkdb").Table("server_config").Count().Run(session)
	if err != nil {
		return 0, fmt.Errorf("unable to query db server config: %s", err)
	}
	defer cursor.Close()

	var serverCount uint
	if err := cursor.One(&serverCount); err != nil {
		return 0, fmt.Errorf("unable to scan db server config count: %s", err)
	}

	return serverCount, nil
}

// SetupDB hadles creating the database and creating all tables and indexes.
func (m *GenericManager) SetupDB(expectedReplicaCount uint) error {
	log.Info("Determining Number of Nodes in DB Cluster")
	replicaCount, err := getNumberOfServers(m.session)
	if err != nil {
		return fmt.Errorf("unable to get number of servers: %s", err)
	}

	log.Infof("Nodes in DB Cluster: %d", replicaCount)

	if expectedReplicaCount > 0 && replicaCount != expectedReplicaCount {
		return fmt.Errorf("unexpected number of nodes in DB cluster: got %d, expected %d", replicaCount, expectedReplicaCount)
	}

	return m.ScaleDB(expectedReplicaCount, false)
}

func (m *GenericManager) ScaleDB(replicaCount uint, emergencyRepair bool) error {
	log.Infof("(Re)configuring Database - replicaCount=%d", replicaCount)
	if err := makeDB(m.db, m.session); err != nil {
		return fmt.Errorf("unable to create database: %s", err)
	}

	type tableConfigResult struct {
		tableName string
		err       error
	}

	results := make(chan tableConfigResult)
	numDone := 0

	allTables := m.Tables()
	log.Infof("(%02d/%02d) Configuring Tables...", numDone, len(allTables))
	for _, tbl := range allTables {
		go func(tbl Table) {
			results <- tableConfigResult{
				tableName: tbl.Name,
				err:       tbl.Configure(m.db, m.session, replicaCount, emergencyRepair),
			}
		}(tbl)
	}

	for numDone < len(allTables) {
		result := <-results

		if result.err != nil {
			return fmt.Errorf("unable to configure table %q: %s", result.tableName, result.err)
		}

		numDone++

		log.Infof("(%02d/%02d) Configured Table %q", numDone, len(allTables), result.tableName)
	}

	return nil
}
