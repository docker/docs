package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	rethink "gopkg.in/dancannon/gorethink.v2"
)

// NOTE: We could just as easily use something like etcd for storing properties
// but that would add more required infrastructure.

var (
	// ErrNoSuchProperty indicates that a property with a given key does
	// not exist.
	ErrNoSuchProperty = errors.New("no such property")
)

// Property is a key/value pair for JSON configuration objects.
type Property struct {
	Key       string `gorethink:"key"`
	JSONValue string `gorethink:"jsonValue"`
}

var propertiesTable = table{
	db:         dbName,
	name:       "properties",
	primaryKey: "key", // Guarantees uniqueness. Quick lookups by key.
}

// GetProperty looks up the property with the given key and unmarshals its
// value into the given val interface which must be a pointer to a type which
// can be unmarshalled from JSON.
func (m *manager) GetProperty(key string, val interface{}) error {
	cursor, err := propertiesTable.Term().Get(key).Run(m.session)
	if err != nil {
		return fmt.Errorf("unable to query DB: %s", err)
	}

	defer cursor.Close()

	if cursor.IsNil() {
		return ErrNoSuchProperty
	}

	var prop Property
	if err := cursor.One(&prop); err != nil {
		return fmt.Errorf("unable to scan property value: %s", err)
	}

	if err := json.Unmarshal([]byte(prop.JSONValue), val); err != nil {
		return fmt.Errorf("unable to unmarshal property JSON value: %s", err)
	}

	return nil
}

// SetProperty sets the property with the given key to the given value. The
// given value must be a type which can be marshalled to JSON.
func (m *manager) SetProperty(key string, val interface{}) error {
	jsonValBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("unable to marshal value to JSON: %s", err)
	}

	prop := Property{
		Key:       key,
		JSONValue: string(jsonValBytes),
	}

	if _, err := propertiesTable.Term().Get(key).Replace(prop).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to set property in database: %s", err)
	}

	return nil
}

// PropertyChange is used to deliver old and new values of a cron as part of a
// changes stream.
type PropertyChange struct {
	OldValue *Property `gorethink:"old_val"`
	NewValue *Property `gorethink:"new_val"`
}

// GetPropertyChanges begins listening for any changes to properties. Returns a
// channel on which the caller may receive a stream of PropertyChange objects
// and an io.Closer which performs necessary cleanup to end the stream's
// underlying goroutine. Only changes for the given keys are streamed. If no
// keys are specified, all property changes are streamed. After closing, the
// changeStream should be checked for a possible remaining value.
func (m *manager) GetPropertyChanges(keys ...string) (changeStream <-chan PropertyChange, streamCloser io.Closer, err error) {
	// Start with all properties.
	query := propertiesTable.Term()

	if len(keys) > 0 {
		// Filter to only the requested property keys. We need to make
		// a slice of interface{} explicitly because Golang is weird
		// like that.
		keyVals := make([]interface{}, len(keys))
		for i, key := range keys {
			keyVals[i] = key
		}

		query = query.GetAllByIndex("key", keyVals...)
	}

	cursor, err := query.Changes(
		rethink.ChangesOpts{IncludeInitial: true},
	).Run(m.session)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query db: %s", err)
	}

	changes := make(chan PropertyChange)
	cursor.Listen(changes)

	return changes, cursor, nil
}

// DeleteProperty deletes the property with the given key if it exists.
func (m *manager) DeleteProperty(key string) error {
	if _, err := propertiesTable.Term().Get(key).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete property in database: %s", err)
	}

	return nil
}
