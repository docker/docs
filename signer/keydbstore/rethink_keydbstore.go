package keydbstore

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dancannon/gorethink"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/storage/rethinkdb"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	jose "github.com/dvsekhvalnov/jose2go"
)

// KeyRethinkDBStore persists and manages private keys on a RethinkDB database
type KeyRethinkDBStore struct {
	sync.Mutex
	session          *gorethink.Session
	defaultPassAlias string
	retriever        passphrase.Retriever
	cachedKeys       map[string]data.PrivateKey
}

// RethinkPrivateKey represents a PrivateKey in the rethink database
type RethinkPrivateKey struct {
	rethinkdb.Timing
	KeyID           string `gorethink:"key_id"`
	EncryptionAlg   string `gorethink:"encryption_alg"`
	KeywrapAlg      string `gorethink:"keywrap_alg"`
	Algorithm       string `gorethink:"algorithm"`
	PassphraseAlias string `gorethink:"passphrase_alias"`
	Public          string `gorethink:"public"`
	Private         string `gorethink:"private"`
}

// TableName sets a specific table name for our RethinkPrivateKey
func (g RethinkPrivateKey) TableName() string {
	return "private_keys"
}

// DatabaseName sets a specific table name for our RethinkPrivateKey
func (g RethinkPrivateKey) DatabaseName() string {
	return "notarysigner"
}

// NewKeyRethinkDBStore returns a new KeyRethinkDBStore backed by a RethinkDB database
func NewKeyRethinkDBStore(passphraseRetriever passphrase.Retriever, defaultPassAlias string, rethinkSession *gorethink.Session) (*KeyRethinkDBStore, error) {
	cachedKeys := make(map[string]data.PrivateKey)

	return &KeyRethinkDBStore{
		session:          rethinkSession,
		defaultPassAlias: defaultPassAlias,
		retriever:        passphraseRetriever,
		cachedKeys:       cachedKeys}, nil
}

// Name returns a user friendly name for the storage location
func (s *KeyRethinkDBStore) Name() string {
	return "RethinkDB"
}

// AddKey stores the contents of a private key. Both role and gun are ignored,
// we always use Key IDs as name, and don't support aliases
func (s *KeyRethinkDBStore) AddKey(keyInfo trustmanager.KeyInfo, privKey data.PrivateKey) error {

	passphrase, _, err := s.retriever(privKey.ID(), s.defaultPassAlias, false, 1)
	if err != nil {
		return err
	}

	encryptedKey, err := jose.Encrypt(string(privKey.Private()), KeywrapAlg, EncryptionAlg, passphrase)
	if err != nil {
		return err
	}

	now := time.Now()
	rethinkPrivKey := RethinkPrivateKey{
		Timing: rethinkdb.Timing{
			CreatedAt: now,
			UpdatedAt: now,
		},
		KeyID:           privKey.ID(),
		EncryptionAlg:   EncryptionAlg,
		KeywrapAlg:      KeywrapAlg,
		PassphraseAlias: s.defaultPassAlias,
		Algorithm:       privKey.Algorithm(),
		Public:          string(privKey.Public()),
		Private:         encryptedKey}

	// Add encrypted private key to the database
	_, err = gorethink.DB(rethinkPrivKey.DatabaseName()).Table(rethinkPrivKey.TableName()).Insert(rethinkPrivKey).RunWrite(s.session)
	if err != nil {
		return fmt.Errorf("failed to add private key to database: %s", privKey.ID())
	}

	// Add the private key to our cache
	s.Lock()
	defer s.Unlock()
	s.cachedKeys[privKey.ID()] = privKey

	return nil
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyRethinkDBStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	cachedKeyEntry, ok := s.cachedKeys[name]
	if ok {
		return cachedKeyEntry, "", nil
	}

	// Retrieve the RethinkDB private key from the database
	dbPrivateKey := RethinkPrivateKey{}
	res, err := gorethink.DB(dbPrivateKey.DatabaseName()).Table(dbPrivateKey.TableName()).Get(RethinkPrivateKey{KeyID: name}).Run(s.session)
	if err != nil {
		return nil, "", trustmanager.ErrKeyNotFound{}
	}
	defer res.Close()

	err = res.One(&dbPrivateKey)
	if err != nil {
		return nil, "", trustmanager.ErrKeyNotFound{}
	}

	// Get the passphrase to use for this key
	passphrase, _, err := s.retriever(dbPrivateKey.KeyID, dbPrivateKey.PassphraseAlias, false, 1)
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
	s.cachedKeys[privKey.ID()] = privKey

	return privKey, "", nil
}

// GetKeyInfo always returns empty and an error. This method is here to satisfy the KeyStore interface
func (s *KeyRethinkDBStore) GetKeyInfo(name string) (trustmanager.KeyInfo, error) {
	return trustmanager.KeyInfo{}, fmt.Errorf("GetKeyInfo currently not supported for KeyRethinkDBStore, as it does not track roles or GUNs")
}

// ListKeys always returns nil. This method is here to satisfy the KeyStore interface
func (s *KeyRethinkDBStore) ListKeys() map[string]trustmanager.KeyInfo {
	return nil
}

// RemoveKey removes the key from the table
func (s *KeyRethinkDBStore) RemoveKey(keyID string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.cachedKeys, keyID)

	// Delete the key from the database
	dbPrivateKey := RethinkPrivateKey{KeyID: keyID}
	_, err := gorethink.DB(dbPrivateKey.DatabaseName()).Table(dbPrivateKey.TableName()).Get(dbPrivateKey).Delete().RunWrite(s.session)
	if err != nil {
		return fmt.Errorf("unable to delete private key from database: %s", err.Error())
	}

	return nil
}

// RotateKeyPassphrase rotates the key-encryption-key
func (s *KeyRethinkDBStore) RotateKeyPassphrase(name, newPassphraseAlias string) error {
	// Retrieve the RethinkDB private key from the database
	dbPrivateKey := RethinkPrivateKey{KeyID: name}
	res, err := gorethink.DB(dbPrivateKey.DatabaseName()).Table(dbPrivateKey.TableName()).Get(dbPrivateKey).Run(s.session)
	if err != nil {
		return trustmanager.ErrKeyNotFound{}
	}
	defer res.Close()

	err = res.One(&dbPrivateKey)
	if err != nil {
		return trustmanager.ErrKeyNotFound{}
	}

	// Get the current passphrase to use for this key
	passphrase, _, err := s.retriever(dbPrivateKey.KeyID, dbPrivateKey.PassphraseAlias, false, 1)
	if err != nil {
		return err
	}

	// Decrypt private bytes from the rethinkDB key
	decryptedPrivKey, _, err := jose.Decode(dbPrivateKey.Private, passphrase)
	if err != nil {
		return err
	}

	// Get the new passphrase to use for this key
	newPassphrase, _, err := s.retriever(dbPrivateKey.KeyID, newPassphraseAlias, false, 1)
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
	if _, err := gorethink.DB(dbPrivateKey.DatabaseName()).Table(dbPrivateKey.TableName()).Get(RethinkPrivateKey{KeyID: name}).Update(dbPrivateKey).RunWrite(s.session); err != nil {
		return err
	}

	return nil
}

// ExportKey is currently unimplemented and will always return an error
func (s *KeyRethinkDBStore) ExportKey(keyID string) ([]byte, error) {
	return nil, errors.New("Exporting from a KeyRethinkDBStore is not supported.")
}

// HealthCheck verifies that DB exists and is query-able
func (s *KeyRethinkDBStore) HealthCheck() error {
	var tableOk bool
	dbPrivateKey := RethinkPrivateKey{}
	res, err := gorethink.DB(dbPrivateKey.DatabaseName()).TableList().Contains(dbPrivateKey.TableName()).Run(s.session)
	if err != nil {
		return err
	}
	defer res.Close()
	err = res.One(tableOk)
	if err != nil || !tableOk {
		return fmt.Errorf(
			"Cannot access table: %s", dbPrivateKey.TableName())
	}
	return nil
}
