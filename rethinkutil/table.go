package rethinkutil

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type Table struct {
	Name       string
	PrimaryKey string
	// Keys are the index names. If len(value) is 0, it is a simple index
	// on the field matching the key. Otherwise, it is a compound index
	// on the list of fields in the corrensponding slice value.
	SecondaryIndexes map[string][]string
}

func (t Table) Term(db string) rethink.Term {
	return rethink.DB(db).Table(t.Name)
}

func (t Table) Wait(db string, session *rethink.Session) error {
	resp, err := t.Term(db).Wait().Run(session)

	if resp != nil {
		resp.Close()
	}

	return err
}

func (t Table) Configure(db string, session *rethink.Session, numReplicas uint, emergencyRepair bool) error {
	createOpts := rethink.TableCreateOpts{
		PrimaryKey: t.PrimaryKey,
		Durability: "hard",
	}

	log.Debugf("Creating table %q", t.Name)
	if _, err := rethink.DB(db).TableCreate(t.Name, createOpts).RunWrite(session); err != nil {
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
		_, err := t.Term(db).Reconfigure(reconfigureOpts).RunWrite(session)
		if err != nil && !strings.Contains(err.Error(), "table doesn't need to be repaired") {
			return fmt.Errorf("unable to repair table: %s", err)
		}
	}

	reconfigureOpts := rethink.ReconfigureOpts{
		Shards:   1,
		Replicas: numReplicas,
	}

	if _, err := t.Term(db).Reconfigure(reconfigureOpts).RunWrite(session); err != nil {
		return fmt.Errorf("unable to reconfigure table replication: %s", err)
	}

	if err := t.Wait(db, session); err != nil {
		return fmt.Errorf("unable to wait for table to be ready after reconfiguring replication: %s", err)
	}

	log.Debug("Configuring Secondary Indexes")
	for indexName, fieldNames := range t.SecondaryIndexes {
		log.Debugf("Configuring Index %q", indexName)
		if len(fieldNames) == 0 {
			// The field name is the index name.
			fieldNames = []string{indexName}
		}

		if _, err := t.Term(db).IndexCreateFunc(indexName, func(row rethink.Term) interface{} {
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

	if err := t.Wait(db, session); err != nil {
		return fmt.Errorf("unable to wait for table to be ready after creating secondary indexes: %s", err)
	}

	return nil
}

func (t Table) DeleteAll(db string, session *rethink.Session) error {
	_, err := t.Term(db).Delete().RunWrite(session)
	return err
}
