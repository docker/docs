package schema

import (
	"fmt"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/manager/schema/interfaces"

	rethink "gopkg.in/dancannon/gorethink.v2"
)

type property struct {
	Key   string `gorethink:"key"`
	Value string `gorethink:"value"`
}

var propertiesTable = table{
	db:         deploy.DTRDBName,
	name:       "properties",
	primaryKey: "key", // Guarantees uniqueness. Quick lookups by key.
}

// RethinkPropertyManager exports CRUDy methods for properties in the database.
type RethinkPropertyManager struct {
	rethinkdb *rethink.Session
}

// NewPropertyManager creates a new property manager using the given database.
func NewRethinkPropertyManager(session *rethink.Session) interfaces.PropertyManager {
	return &RethinkPropertyManager{rethinkdb: session}
}

// Set assigns the given value to the property key in the database.
func (p *RethinkPropertyManager) Set(key, value string) error {
	prop := property{
		Key:   key,
		Value: value,
	}

	if _, err := propertiesTable.Term().Get(key).Replace(prop).RunWrite(p.rethinkdb); err != nil {
		return fmt.Errorf("unable to set property in database: %s", err)
	}

	return nil
}

// Get retrieves the value assigned to the given proeprty key in the database.
func (p *RethinkPropertyManager) Get(key string) (string, error) {
	cursor, err := propertiesTable.Term().Get(key).Run(p.rethinkdb)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return "", interfaces.ErrPropertyNotSet
		} else {
			return "", fmt.Errorf("unable to query DB: %s", err)
		}
	}

	defer cursor.Close()

	if cursor.IsNil() {
		return "", interfaces.ErrPropertyNotSet
	}

	var prop property
	if err := cursor.One(&prop); err != nil {
		return "", fmt.Errorf("unable to scan property value: %s", err)
	}

	return prop.Value, nil
}
