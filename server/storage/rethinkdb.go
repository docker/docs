package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/docker/notary/storage/rethinkdb"
	"github.com/docker/notary/tuf/data"
	"gopkg.in/dancannon/gorethink.v2"
)

// RDBTUFFile is a tuf file record
type RDBTUFFile struct {
	rethinkdb.Timing
	GunRoleVersion []interface{} `gorethink:"gun_role_version"`
	Gun            string        `gorethink:"gun"`
	Role           string        `gorethink:"role"`
	Version        int           `gorethink:"version"`
	Sha256         string        `gorethink:"sha256"`
	Data           []byte        `gorethink:"data"`
	TSchecksum     string        `gorethink:"timestamp_checksum"`
}

// TableName returns the table name for the record type
func (r RDBTUFFile) TableName() string {
	return "tuf_files"
}

// RDBKey is the public key record
type RDBKey struct {
	rethinkdb.Timing
	Gun    string `gorethink:"gun"`
	Role   string `gorethink:"role"`
	Cipher string `gorethink:"cipher"`
	Public []byte `gorethink:"public"`
}

// TableName returns the table name for the record type
func (r RDBKey) TableName() string {
	return "tuf_keys"
}

// RethinkDB implements a MetaStore against the Rethink Database
type RethinkDB struct {
	dbName string
	sess   *gorethink.Session
}

// NewRethinkDBStorage initializes a RethinkDB object
func NewRethinkDBStorage(dbName string, sess *gorethink.Session) RethinkDB {
	return RethinkDB{
		dbName: dbName,
		sess:   sess,
	}
}

// GetKey returns the cipher and public key for the given GUN and role.
// If the GUN+role don't exist, returns an error.
func (rdb RethinkDB) GetKey(gun, role string) (cipher string, public []byte, err error) {
	var key RDBKey
	res, err := gorethink.DB(rdb.dbName).Table(key.TableName()).GetAllByIndex(
		rdbGunRoleIdx, []string{gun, role},
	).Run(rdb.sess)
	if err != nil {
		return "", nil, err
	}
	defer res.Close()
	err = res.One(&key)
	if err == gorethink.ErrEmptyResult {
		return "", nil, &ErrNoKey{gun: gun}
	}
	return key.Cipher, key.Public, err
}

// SetKey sets the cipher and public key for the given GUN and role if
// it doesn't already exist.  Otherwise an error is returned.
func (rdb RethinkDB) SetKey(gun, role, cipher string, public []byte) error {
	now := time.Now()
	key := RDBKey{
		Timing: rethinkdb.Timing{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Gun:    gun,
		Role:   role,
		Cipher: cipher,
		Public: public,
	}
	_, err := gorethink.DB(rdb.dbName).Table(key.TableName()).Insert(key).RunWrite(rdb.sess)
	return err
}

// UpdateCurrent adds new metadata version for the given GUN if and only
// if it's a new role, or the version is greater than the current version
// for the role. Otherwise an error is returned.
func (rdb RethinkDB) UpdateCurrent(gun string, update MetaUpdate) error {
	now := time.Now()
	checksum := sha256.Sum256(update.Data)
	file := RDBTUFFile{
		Timing: rethinkdb.Timing{
			CreatedAt: now,
			UpdatedAt: now,
		},
		GunRoleVersion: []interface{}{gun, update.Role, update.Version},
		Gun:            gun,
		Role:           update.Role,
		Version:        update.Version,
		Sha256:         hex.EncodeToString(checksum[:]),
		Data:           update.Data,
	}
	_, err := gorethink.DB(rdb.dbName).Table(file.TableName()).Insert(
		file,
		gorethink.InsertOpts{
			Conflict: "error", // default but explicit for clarity of intent
		},
	).RunWrite(rdb.sess)
	if err != nil && gorethink.IsConflictErr(err) {
		return &ErrOldVersion{}
	}
	return err
}

// UpdateCurrentWithTSChecksum adds new metadata version for the given GUN with an associated
// checksum for the timestamp it belongs to, to afford us transaction-like functionality
func (rdb RethinkDB) UpdateCurrentWithTSChecksum(gun, tsChecksum string, update MetaUpdate) error {
	now := time.Now()
	checksum := sha256.Sum256(update.Data)
	file := RDBTUFFile{
		Timing: rethinkdb.Timing{
			CreatedAt: now,
			UpdatedAt: now,
		},
		GunRoleVersion: []interface{}{gun, update.Role, update.Version},
		Gun:            gun,
		Role:           update.Role,
		Version:        update.Version,
		Sha256:         hex.EncodeToString(checksum[:]),
		TSchecksum:     tsChecksum,
		Data:           update.Data,
	}
	_, err := gorethink.DB(rdb.dbName).Table(file.TableName()).Insert(
		file,
		gorethink.InsertOpts{
			Conflict: "error", // default but explicit for clarity of intent
		},
	).RunWrite(rdb.sess)
	if err != nil && gorethink.IsConflictErr(err) {
		return &ErrOldVersion{}
	}
	return err
}

// Used for sorting updates alphabetically by role name, such that timestamp is always last:
// Ordering: root, snapshot, targets, targets/* (delegations), timestamp
type updateSorter []MetaUpdate

func (u updateSorter) Len() int      { return len(u) }
func (u updateSorter) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
func (u updateSorter) Less(i, j int) bool {
	return u[i].Role < u[j].Role
}

// UpdateMany adds multiple new metadata for the given GUN. RethinkDB does
// not support transactions, therefore we will attempt to insert the timestamp
// last as this represents a published version of the repo.  However, we will
// insert all other role data in alphabetical order first, and also include the
// associated timestamp checksum so that we can easily roll back this pseudotransaction
func (rdb RethinkDB) UpdateMany(gun string, updates []MetaUpdate) error {
	// find the timestamp first and save its checksum
	// then apply the updates in alphabetic role order with the timestamp last
	// if there are any failures, we roll back in the same alphabetic order
	var tsChecksum string
	for _, up := range updates {
		if up.Role == data.CanonicalTimestampRole {
			tsChecksumBytes := sha256.Sum256(up.Data)
			tsChecksum = hex.EncodeToString(tsChecksumBytes[:])
			break
		}
	}

	// alphabetize the updates by Role name
	sort.Stable(updateSorter(updates))

	for _, up := range updates {
		if err := rdb.UpdateCurrentWithTSChecksum(gun, tsChecksum, up); err != nil {
			// roll back with best-effort deletion, and then error out
			rdb.deleteByTSChecksum(tsChecksum)
			return err
		}
	}
	return nil
}

// GetCurrent returns the modification date and data part of the metadata for
// the latest version of the given GUN and role.  If there is no data for
// the given GUN and role, an error is returned.
func (rdb RethinkDB) GetCurrent(gun, role string) (created *time.Time, data []byte, err error) {
	file := RDBTUFFile{}
	res, err := gorethink.DB(rdb.dbName).Table(file.TableName(), gorethink.TableOpts{ReadMode: "majority"}).GetAllByIndex(
		rdbGunRoleIdx, []string{gun, role},
	).OrderBy(gorethink.Desc("version")).Run(rdb.sess)
	if err != nil {
		return nil, nil, err
	}
	defer res.Close()
	if res.IsNil() {
		return nil, nil, ErrNotFound{}
	}
	err = res.One(&file)
	if err == gorethink.ErrEmptyResult {
		return nil, nil, ErrNotFound{}
	}
	return &file.CreatedAt, file.Data, err
}

// GetChecksum returns the given TUF role file and creation date for the
// GUN with the provided checksum. If the given (gun, role, checksum) are
// not found, it returns storage.ErrNotFound
func (rdb RethinkDB) GetChecksum(gun, role, checksum string) (created *time.Time, data []byte, err error) {
	var file RDBTUFFile
	res, err := gorethink.DB(rdb.dbName).Table(file.TableName(), gorethink.TableOpts{ReadMode: "majority"}).GetAllByIndex(
		rdbGunRoleSha256Idx, []string{gun, role, checksum},
	).Run(rdb.sess)
	if err != nil {
		return nil, nil, err
	}
	defer res.Close()
	if res.IsNil() {
		return nil, nil, ErrNotFound{}
	}
	err = res.One(&file)
	if err == gorethink.ErrEmptyResult {
		return nil, nil, ErrNotFound{}
	}
	return &file.CreatedAt, file.Data, err
}

// Delete removes all metadata for a given GUN.  It does not return an
// error if no metadata exists for the given GUN.
func (rdb RethinkDB) Delete(gun string) error {
	_, err := gorethink.DB(rdb.dbName).Table(RDBTUFFile{}.TableName()).GetAllByIndex(
		"gun", []string{gun},
	).Delete().RunWrite(rdb.sess)
	if err != nil {
		return fmt.Errorf("unable to delete %s from database: %s", gun, err.Error())
	}
	return nil
}

// deleteByTSChecksum removes all metadata by a timestamp checksum, used for rolling back a "transaction"
// from a call to rethinkdb's UpdateMany
func (rdb RethinkDB) deleteByTSChecksum(tsChecksum string) error {
	_, err := gorethink.DB(rdb.dbName).Table(RDBTUFFile{}.TableName()).GetAllByIndex(
		"timestamp_checksum", []string{tsChecksum},
	).Delete().RunWrite(rdb.sess)
	if err != nil {
		return fmt.Errorf("unable to delete timestamp checksum data: %s from database: %s", tsChecksum, err.Error())
	}
	return nil
}

// Bootstrap sets up the database and tables
func (rdb RethinkDB) Bootstrap() error {
	return rethinkdb.SetupDB(rdb.sess, rdb.dbName, []rethinkdb.Table{
		tufFiles,
		keys,
	})
}

// CheckHealth is currently a noop
func (rdb RethinkDB) CheckHealth() error {
	return nil
}
