package keydbstore

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/storage/rethinkdb"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/utils"
	jose "github.com/dvsekhvalnov/jose2go"
	"gopkg.in/dancannon/gorethink.v2"
)

// RethinkDBKeyStore persists and manages private keys on a RethinkDB database
type RethinkDBKeyStore struct {
	lock             *sync.Mutex
	sess             *gorethink.Session
	dbName           string
	defaultPassAlias string
	retriever        passphrase.Retriever
	cachedKeys       map[string]data.PrivateKey
}

// RDBPrivateKey represents a PrivateKey in the rethink database
type RDBPrivateKey struct {
	rethinkdb.Timing
	KeyID           string `gorethink:"key_id"`
	EncryptionAlg   string `gorethink:"encryption_alg"`
	KeywrapAlg      string `gorethink:"keywrap_alg"`
	Algorithm       string `gorethink:"algorithm"`
	PassphraseAlias string `gorethink:"passphrase_alias"`
	Public          string `gorethink:"public"`
	Private         string `gorethink:"private"`
}

var privateKeys = rethinkdb.Table{
	Name:       RDBPrivateKey{}.TableName(),
	PrimaryKey: RDBPrivateKey{}.KeyID,
}

// TableName sets a specific table name for our RDBPrivateKey
func (g RDBPrivateKey) TableName() string {
	return "private_keys"
}

// NewRethinkDBKeyStore returns a new RethinkDBKeyStore backed by a RethinkDB database
func NewRethinkDBKeyStore(dbName string, passphraseRetriever passphrase.Retriever, defaultPassAlias string, rethinkSession *gorethink.Session) *RethinkDBKeyStore {
	cachedKeys := make(map[string]data.PrivateKey)

	return &RethinkDBKeyStore{
		lock:             &sync.Mutex{},
		sess:             rethinkSession,
		defaultPassAlias: defaultPassAlias,
		dbName:           dbName,
		retriever:        passphraseRetriever,
		cachedKeys:       cachedKeys,
	}
}

// Name returns a user friendly name for the storage location
func (rdb *RethinkDBKeyStore) Name() string {
	return "RethinkDB"
}

// AddKey stores the contents of a private key. Both role and gun are ignored,
// we always use Key IDs as name, and don't support aliases
func (rdb *RethinkDBKeyStore) AddKey(keyInfo trustmanager.KeyInfo, privKey data.PrivateKey) error {

	passphrase, _, err := rdb.retriever(privKey.ID(), rdb.defaultPassAlias, false, 1)
	if err != nil {
		return err
	}

	encryptedKey, err := jose.Encrypt(string(privKey.Private()), KeywrapAlg, EncryptionAlg, passphrase)
	if err != nil {
		return err
	}

	now := time.Now()
	rethinkPrivKey := RDBPrivateKey{
		Timing: rethinkdb.Timing{
			CreatedAt: now,
			UpdatedAt: now,
		},
		KeyID:           privKey.ID(),
		EncryptionAlg:   EncryptionAlg,
		KeywrapAlg:      KeywrapAlg,
		PassphraseAlias: rdb.defaultPassAlias,
		Algorithm:       privKey.Algorithm(),
		Public:          string(privKey.Public()),
		Private:         encryptedKey}

	// Add encrypted private key to the database
	_, err = gorethink.DB(rdb.dbName).Table(rethinkPrivKey.TableName()).Insert(rethinkPrivKey).RunWrite(rdb.sess)
	if err != nil {
		return fmt.Errorf("failed to add private key to database: %s", privKey.ID())
	}

	// Add the private key to our cache
	rdb.lock.Lock()
	defer rdb.lock.Unlock()
	rdb.cachedKeys[privKey.ID()] = privKey

	return nil
}

// GetKey returns the PrivateKey given a KeyID
func (rdb *RethinkDBKeyStore) GetKey(name string) (data.PrivateKey, string, error) {
	rdb.lock.Lock()
	defer rdb.lock.Unlock()
	cachedKeyEntry, ok := rdb.cachedKeys[name]
	if ok {
		return cachedKeyEntry, "", nil
	}

	// Retrieve the RethinkDB private key from the database
	dbPrivateKey := RDBPrivateKey{}
	res, err := gorethink.DB(rdb.dbName).Table(dbPrivateKey.TableName()).Filter(gorethink.Row.Field("key_id").Eq(name)).Run(rdb.sess)
	if err != nil {
		return nil, "", trustmanager.ErrKeyNotFound{}
	}
	defer res.Close()

	err = res.One(&dbPrivateKey)
	if err != nil {
		return nil, "", trustmanager.ErrKeyNotFound{}
	}

	// Get the passphrase to use for this key
	passphrase, _, err := rdb.retriever(dbPrivateKey.KeyID, dbPrivateKey.PassphraseAlias, false, 1)
	if err != nil {
		return nil, "", err
	}

	// Decrypt private bytes from the gorm key
	decryptedPrivKey, _, err := jose.Decode(dbPrivateKey.Private, passphrase)
	if err != nil {
		return nil, "", err
	}

	pubKey := data.NewPublicKey(dbPrivateKey.Algorithm, []byte(dbPrivateKey.Public))
	// Create a new PrivateKey with unencrypted bytes
	privKey, err := data.NewPrivateKey(pubKey, []byte(decryptedPrivKey))
	if err != nil {
		return nil, "", err
	}

	// Add the key to cache
	rdb.cachedKeys[privKey.ID()] = privKey

	return privKey, "", nil
}

// GetKeyInfo always returns empty and an error. This method is here to satisfy the KeyStore interface
func (rdb RethinkDBKeyStore) GetKeyInfo(name string) (trustmanager.KeyInfo, error) {
	return trustmanager.KeyInfo{}, fmt.Errorf("GetKeyInfo currently not supported for RethinkDBKeyStore, as it does not track roles or GUNs")
}

// ListKeys always returns nil. This method is here to satisfy the KeyStore interface
func (rdb RethinkDBKeyStore) ListKeys() map[string]trustmanager.KeyInfo {
	return nil
}

// RemoveKey removes the key from the table
func (rdb RethinkDBKeyStore) RemoveKey(keyID string) error {
	rdb.lock.Lock()
	defer rdb.lock.Unlock()

	delete(rdb.cachedKeys, keyID)

	// Delete the key from the database
	dbPrivateKey := RDBPrivateKey{KeyID: keyID}
	_, err := gorethink.DB(rdb.dbName).Table(dbPrivateKey.TableName()).Filter(gorethink.Row.Field("key_id").Eq(keyID)).Delete().RunWrite(rdb.sess)
	if err != nil {
		return fmt.Errorf("unable to delete private key from database: %s", err.Error())
	}

	return nil
}

// RotateKeyPassphrase rotates the key-encryption-key
func (rdb RethinkDBKeyStore) RotateKeyPassphrase(name, newPassphraseAlias string) error {
	// Retrieve the RethinkDB private key from the database
	dbPrivateKey := RDBPrivateKey{KeyID: name}
	res, err := gorethink.DB(rdb.dbName).Table(dbPrivateKey.TableName()).Get(dbPrivateKey).Run(rdb.sess)
	if err != nil {
		return trustmanager.ErrKeyNotFound{}
	}
	defer res.Close()

	err = res.One(&dbPrivateKey)
	if err != nil {
		return trustmanager.ErrKeyNotFound{}
	}

	// Get the current passphrase to use for this key
	passphrase, _, err := rdb.retriever(dbPrivateKey.KeyID, dbPrivateKey.PassphraseAlias, false, 1)
	if err != nil {
		return err
	}

	// Decrypt private bytes from the rethinkDB key
	decryptedPrivKey, _, err := jose.Decode(dbPrivateKey.Private, passphrase)
	if err != nil {
		return err
	}

	// Get the new passphrase to use for this key
	newPassphrase, _, err := rdb.retriever(dbPrivateKey.KeyID, newPassphraseAlias, false, 1)
	if err != nil {
		return err
	}

	// Re-encrypt the private bytes with the new passphrase
	newEncryptedKey, err := jose.Encrypt(decryptedPrivKey, KeywrapAlg, EncryptionAlg, newPassphrase)
	if err != nil {
		return err
	}

	// Update the database object
	dbPrivateKey.Private = newEncryptedKey
	dbPrivateKey.PassphraseAlias = newPassphraseAlias
	if _, err := gorethink.DB(rdb.dbName).Table(dbPrivateKey.TableName()).Get(RDBPrivateKey{KeyID: name}).Update(dbPrivateKey).RunWrite(rdb.sess); err != nil {
		return err
	}

	return nil
}

// ExportKey is currently unimplemented and will always return an error
func (rdb RethinkDBKeyStore) ExportKey(keyID string) ([]byte, error) {
	return nil, errors.New("Exporting from a RethinkDBKeyStore is not supported.")
}

// Bootstrap sets up the database and tables
func (rdb RethinkDBKeyStore) Bootstrap() error {
	return rethinkdb.SetupDB(rdb.sess, rdb.dbName, []rethinkdb.Table{
		privateKeys,
	})
}

// CheckHealth verifies that DB exists and is query-able
func (rdb RethinkDBKeyStore) CheckHealth() error {
	var tables []string
	dbPrivateKey := RDBPrivateKey{}
	res, err := gorethink.DB(rdb.dbName).TableList().Run(rdb.sess)
	if err != nil {
		return err
	}
	defer res.Close()
	err = res.All(&tables)
	if err != nil || !utils.StrSliceContains(tables, dbPrivateKey.TableName()) {
		return fmt.Errorf(
			"Cannot access table: %s", dbPrivateKey.TableName())
	}
	return nil
}
