package schema

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

const dbName = "enzi"

func makeDB(session *rethink.Session, name string) error {
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

type table struct {
	db         string
	name       string
	primaryKey string
	// Keys are the index names. If len(value) is 0, it is a simple index
	// on the field matching the key. Otherwise, it is a compound index
	// on the list of fields in the corrensponding slice value.
	secondaryIndexes map[string][]string
}

func (t table) Term() rethink.Term {
	return rethink.DB(t.db).Table(t.name)
}

func (t table) wait(session *rethink.Session) error {
	resp, err := t.Term().Wait().Run(session)

	if resp != nil {
		resp.Close()
	}

	return err
}

func (t table) configure(session *rethink.Session, numReplicas uint, emergencyRepair bool) error {
	createOpts := rethink.TableCreateOpts{
		PrimaryKey: t.primaryKey,
		Durability: "hard",
	}

	log.Debugf("Creating table %q", t.name)
	if _, err := rethink.DB(t.db).TableCreate(t.name, createOpts).RunWrite(session); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("unable to run table creation: %s", err)
		}

		log.Debug("Table Already Exists")
	}

	if emergencyRepair {
		reconfigureOpts := rethink.ReconfigureOpts{
			// This operation is considered "unsafe" because if you
			// have multiple shards and aren't 100% certain that
			// all shards are available then if you run this repair
			// then all of the data in those unavailable shards are
			// lost.
			// In our case, there's only 1 shard for all tables
			// with a replica on each available server so this
			// operation will leave at least one copy of all of the
			// data.
			EmergencyRepair: "unsafe_rollback",
		}
		_, err := t.Term().Reconfigure(reconfigureOpts).RunWrite(session)
		if err != nil && !strings.Contains(err.Error(), "table doesn't need to be repaired") {
			return fmt.Errorf("unable to repair table: %s", err)
		}
	}

	reconfigureOpts := rethink.ReconfigureOpts{
		Shards:   1,
		Replicas: numReplicas,
	}

	if _, err := t.Term().Reconfigure(reconfigureOpts).RunWrite(session); err != nil {
		return fmt.Errorf("unable to reconfigure table replication: %s", err)
	}

	if err := t.wait(session); err != nil {
		return fmt.Errorf("unable to wait for table to be ready after reconfiguring replication: %s", err)
	}

	log.Debug("Configuring Secondary Indexes")
	for indexName, fieldNames := range t.secondaryIndexes {
		log.Debugf("Configuring Index %q", indexName)
		if len(fieldNames) == 0 {
			// The field name is the index name.
			fieldNames = []string{indexName}
		}

		if _, err := t.Term().IndexCreateFunc(indexName, func(row rethink.Term) interface{} {
			fields := make([]interface{}, len(fieldNames))

			for i, fieldName := range fieldNames {
				term := row
				for _, subfield := range strings.Split(fieldName, ".") {
					term = term.Field(subfield)
				}

				fields[i] = term
			}

			if len(fields) == 1 {
				return fields[0]
			}

			return fields
		}).RunWrite(session); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("unable to create secondary index %q: %s", indexName, err)
			}

			log.Debug("Index Already Exists")
		}
	}

	if err := t.wait(session); err != nil {
		return fmt.Errorf("unable to wait for table to be ready after creating secondary indexes: %s", err)
	}

	return nil
}

func (t table) deleteAll(session *rethink.Session) error {
	_, err := t.Term().Delete().RunWrite(session)
	return err
}

func isDuplicatePrimaryKeyErr(resp rethink.WriteResponse) bool {
	return strings.HasPrefix(resp.FirstError, "Duplicate primary key")
}

var allTables = []table{
	accountsTable,
	orgMembershipTable,
	teamsTable,
	teamMembershipTable,
	propertiesTable,
	sessionsTable,
	servicesTable,
	serviceSessionsTable,
	serviceAuthCodesTable,
	serviceKeysTable,
	signingKeysTable,
	workersTable,
	jobsTable,
}

// SetupDB hadles creating the database and creating all tables and indexes.
func SetupDB(session *rethink.Session, expectedReplicaCount uint, emergencyRepair bool) error {
	log.Info("Determining Number of Nodes in DB Cluster")
	replicaCount, err := getNumberOfServers(session)
	if err != nil {
		return fmt.Errorf("unable to get number of servers: %s", err)
	}

	log.Infof("Nodes in DB Cluster: %d", replicaCount)

	if expectedReplicaCount > 0 && replicaCount != expectedReplicaCount {
		return fmt.Errorf("unexpected number of nodes in DB cluster: got %d, expected %d", replicaCount, expectedReplicaCount)
	}

	return configureDB(session, replicaCount, emergencyRepair)
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

func configureDB(session *rethink.Session, replicaCount uint, emergencyRepair bool) error {
	log.Infof("(Re)configuring Database - replicaCount=%d", replicaCount)
	if err := makeDB(session, dbName); err != nil {
		return fmt.Errorf("unable to create database: %s", err)
	}

	type tableConfigResult struct {
		tableName string
		err       error
	}

	results := make(chan tableConfigResult)
	numDone := 0

	log.Infof("(%02d/%02d) Configuring Tables...", numDone, len(allTables))
	for _, tbl := range allTables {
		go func(tbl table) {
			results <- tableConfigResult{
				tableName: tbl.name,
				err:       tbl.configure(session, replicaCount, emergencyRepair),
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

// DrainServer removes all replicas from the server with the given hostname and
// port. The cluster's server list is searched for a matching server and all
// server tags are removed from that server. Next, all tables are reconfigured
// to a replica count that is one fewer than the number of servers that are
// currently in the cluster.
func DrainServer(session *rethink.Session, host string, port uint) error {
	// First, find the server with this canonical address.
	server, err := findServer(session, host, port)
	if err != nil {
		return fmt.Errorf("unable to find server: %s", err)
	}

	// Now that we have the server, update it's config to remove all tags.
	if _, err := rethink.DB("rethinkdb").Table("server_config").Get(server.ID).Update(
		map[string]interface{}{"tags": []string{}},
	).RunWrite(session); err != nil {
		return fmt.Errorf("unable to update server config: %s", err)
	}

	log.Info("Determining Number of Servers in DB Cluster")
	numberOfServers, err := getNumberOfServers(session)
	if err != nil {
		return fmt.Errorf("unable to get number of servers: %s", err)
	}

	log.Infof("Nodes in DB Cluster: %d", numberOfServers)

	// Finally, attempt to reconfigure all database tables to use one fewer
	// replicas then there are servers. Since the target server has had all
	// tags removed (specifically the 'default' tag), no replicas will be
	// placed on it.

	if numberOfServers < 2 {
		return fmt.Errorf("refusing to reduce replication to less than 1")
	}

	return configureDB(session, numberOfServers-1, false)
}

func findServer(session *rethink.Session, host string, port uint) (*serverStatus, error) {
	cursor, err := rethink.DB("rethinkdb").Table("server_status").Run(session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db server status: %s", err)
	}
	defer cursor.Close()

	var server serverStatus
	for cursor.Next(&server) {
		if server.hasHostPort(host, port) {
			return &server, nil
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("unable to iterate through server status records: %s", err)
	}

	return nil, fmt.Errorf("no such server with host %s, port %d", host, port)
}

type serverStatus struct {
	ID      string        `gorethink:"id"`
	Network networkConfig `gorethink:"network"`
}

func (s *serverStatus) hasHostPort(host string, port uint) bool {
	for _, addr := range s.Network.CanonicalAddresses {
		if addr.Host == host && addr.Port == port {
			return true
		}
	}

	return false
}

type networkConfig struct {
	CanonicalAddresses []canonicalAddr `gorethink:"canonical_addresses"`
}

type canonicalAddr struct {
	Host string `gorethink:"host"`
	Port uint   `gorethink:"port"`
}

// TableWaitTarget is a type used with the WaitForReadyTables function.
type TableWaitTarget string

// These constants are valid TableWaitTarget values.
const (
	ReadyForOutdatedReads TableWaitTarget = "ready_for_outdated_reads"
	ReadyForReads                         = "ready_for_reads"
	ReadyForWrites                        = "ready_for_writes"
	AllReplicasReady                      = "all_replicas_ready"
)

// WaitForReadyTables waits for all tables to be ready. If the given waitFor
// TableWaitTarget value is an empty string, the default of AllReplicasReady is
// used.
func WaitForReadyTables(session *rethink.Session, waitFor TableWaitTarget) error {
	type tableWaitResult struct {
		tableName string
		err       error
	}

	results := make(chan tableWaitResult, len(allTables))
	numReady := 0

	log.Debugf("(%02d/%02d) Waiting for Tables to be Ready...", numReady, len(allTables))
	for _, tbl := range allTables {
		go func(tbl table) {
			resp, err := tbl.Term().Wait(rethink.WaitOpts{
				WaitFor: waitFor,
			}).Run(session)
			if resp != nil {
				resp.Close()
			}

			results <- tableWaitResult{
				tableName: tbl.name,
				err:       err,
			}
		}(tbl)
	}

	for numReady < len(allTables) {
		result := <-results

		if result.err != nil {
			return fmt.Errorf("unable to wait for table %s: %s", result.tableName, result.err)
		}

		numReady++

		log.Debugf("(%02d/%02d) Table %s is Ready", numReady, len(allTables), result.tableName)
	}

	return nil
}
