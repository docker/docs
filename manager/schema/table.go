package schema

import (
	"fmt"
	"strings"

	"github.com/docker/dhe-deploy"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// this file is mostly taken from
// https://github.com/docker/enzi/blob/master/schema/table.go

func makeDB(session *rethink.Session, name string) error {
	_, err := rethink.DBCreate(name).RunWrite(session)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	return nil
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

func (t table) create(session *rethink.Session, numReplicas uint) error {
	opts := rethink.TableCreateOpts{
		PrimaryKey: t.primaryKey,
		Durability: "hard",
		DataCenter: numReplicas,
	}

	if _, err := rethink.DB(t.db).TableCreate(t.name, opts).RunWrite(session); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("unable to run table creation: %s", err)
		}
	}

	for indexName, fieldNames := range t.secondaryIndexes {
		var err error
		if len(fieldNames) == 0 {
			_, err = t.Term().IndexCreate(indexName).RunWrite(session)
		} else {
			_, err = t.Term().IndexCreateFunc(indexName, func(row rethink.Term) interface{} {
				fields := make([]interface{}, len(fieldNames))

				for i, fieldName := range fieldNames {
					term := row
					for _, subfield := range strings.Split(fieldName, ".") {
						term = term.Field(subfield)
					}

					fields[i] = term
				}

				return fields
			}).RunWrite(session)
		}

		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("unable to create secondary index %q: %s", indexName, err)
			}
		}
	}

	if err := t.wait(session); err != nil {
		return fmt.Errorf("unable to wait for table to be ready: %s", err)
	}

	return nil
}

func (t table) deleteAll(session *rethink.Session) error {
	_, err := t.Term().Delete().RunWrite(session)
	return err
}

func (t table) getRowByIndexVal(session *rethink.Session, indexName string, val interface{}, row interface{}, noSuchErr error) error {
	cursor, err := t.Term().GetAllByIndex(indexName, val).Run(session)
	if err != nil {
		return fmt.Errorf("unable to query db: %s", err)
	}

	if err := cursor.One(row); err != nil {
		if err == rethink.ErrEmptyResult {
			return noSuchErr
		}

		return fmt.Errorf("unable to get query result: %s", err)
	}

	return nil
}

func isDuplicatePrimaryKeyErr(resp rethink.WriteResponse) bool {
	return strings.HasPrefix(resp.FirstError, "Duplicate primary key")
}

// SetupDB hadles creating the database and creating all tables and indexes.
func SetupDB(session *rethink.Session, allTables ...table) error {
	if err := makeDB(session, deploy.DTRDBName); err != nil {
		return fmt.Errorf("unable to create database: %s", err)
	}

	cursor, err := rethink.DB("rethinkdb").Table("server_config").Count().Run(session)
	if err != nil {
		return fmt.Errorf("unable to query db server config: %s", err)
	}

	var replicaCount uint
	if err := cursor.One(&replicaCount); err != nil {
		return fmt.Errorf("unable to scan db server config count: %s", err)
	}

	for _, table := range allTables {
		if err = table.create(session, replicaCount); err != nil {
			return fmt.Errorf("unable to create table %q: %s", table.name, err)
		}
	}

	return nil
}
