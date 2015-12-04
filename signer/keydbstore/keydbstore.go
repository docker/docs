package keydbstore

import (
	"errors"
	"fmt"
	"sync"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/jinzhu/gorm"
)

// Constants
const (
	EncryptionAlg = jose.A256GCM
	KeywrapAlg    = jose.PBES2_HS256_A128KW
)

// KeyDBStore persists and manages private keys on a SQL database
type KeyDBStore struct {
	sync.Mutex
	db               gorm.DB
	defaultPassAlias string
	retriever        passphrase.Retriever
	cachedKeys       map[string]data.PrivateKey
}

// GormPrivateKey represents a PrivateKey in the database
type GormPrivateKey struct {
	gorm.Model
	KeyID           string `sql:"not null;unique;index:key_id_idx"`
	EncryptionAlg   string `sql:"not null"`
	KeywrapAlg      string `sql:"not null"`
	Algorithm       string `sql:"not null"`
	PassphraseAlias string `sql:"not null"`
	Public          string `sql:"not null"`
	Private         string `sql:"not null"`
}

// TableName sets a specific table name for our GormPrivateKey
func (g GormPrivateKey) TableName() string {
	return "private_keys"
}

// NewKeyDBStore returns a new KeyDBStore backed by a SQL database
func NewKeyDBStore(passphraseRetriever passphrase.Retriever, defaultPassAlias string,
	dbDialect string, dbArgs ...interface{}) (*KeyDBStore, error) {
	cachedKeys := make(map[string]data.PrivateKey)

	db, err := gorm.Open(dbDialect, dbArgs...)
	if err != nil {
		return nil, err
	}

	return &KeyDBStore{
		db:               db,
		defaultPassAlias: defaultPassAlias,
		retriever:        passphraseRetriever,
		cachedKeys:       cachedKeys}, nil
}

// Name returns a user friendly name for the storage location
func (s *KeyDBStore) Name() string {
	return "database"
}

// AddKey stores the contents of a private key. Both name and alias are ignored,
// we always use Key IDs as name, and don't support aliases
func (s *KeyDBStore) AddKey(name, alias string, privKey data.PrivateKey) error {

	passphrase, _, err := s.retriever(privKey.ID(), s.defaultPassAlias, false, 1)
	if err != nil {
		return err
	}

	encryptedKey, err := jose.Encrypt(string(privKey.Private()), KeywrapAlg, EncryptionAlg, passphrase)
	if err != nil {
		return err
	}

	gormPrivKey := GormPrivateKey{
		KeyID:           privKey.ID(),
		EncryptionAlg:   EncryptionAlg,
		KeywrapAlg:      KeywrapAlg,
		PassphraseAlias: s.defaultPassAlias,
		Algorithm:       privKey.Algorithm(),
		Public:          string(privKey.Public()),
		Private:         encryptedKey}

	// Add encrypted private key to the database
	s.db.Create(&gormPrivKey)
	// Value will be false if Create succeeds
	failure := s.db.NewRecord(gormPrivKey)
	if failure {
		return fmt.Errorf("failed to add private key to database: %s", privKey.ID())
	}

	// Add the private key to our cache
	s.Lock()
	defer s.Unlock()
	s.cachedKeys[privKey.ID()] = privKey

	return nil
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyDBStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	cachedKeyEntry, ok := s.cachedKeys[name]
	if ok {
		return cachedKeyEntry, "", nil
	}

	// Retrieve the GORM private key from the database
	dbPrivateKey := GormPrivateKey{}
	if s.db.Where(&GormPrivateKey{KeyID: name}).First(&dbPrivateKey).RecordNotFound() {
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

// ListKeys always returns nil. This method is here to satisfy the KeyStore interface
func (s *KeyDBStore) ListKeys() map[string]string {
	return nil
}

// RemoveKey removes the key from the keyfilestore
func (s *KeyDBStore) RemoveKey(name string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.cachedKeys, name)

	// Retrieve the GORM private key from the database
	dbPrivateKey := GormPrivateKey{}
	if s.db.Where(&GormPrivateKey{KeyID: name}).First(&dbPrivateKey).RecordNotFound() {
		return trustmanager.ErrKeyNotFound{}
	}

	// Delete the key from the database
	s.db.Delete(&dbPrivateKey)

	return nil
}

// RotateKeyPassphrase rotates the key-encryption-key
func (s *KeyDBStore) RotateKeyPassphrase(name, newPassphraseAlias string) error {
	// Retrieve the GORM private key from the database
	dbPrivateKey := GormPrivateKey{}
	if s.db.Where(&GormPrivateKey{KeyID: name}).First(&dbPrivateKey).RecordNotFound() {
		return trustmanager.ErrKeyNotFound{}
	}

	// Get the current passphrase to use for this key
	passphrase, _, err := s.retriever(dbPrivateKey.KeyID, dbPrivateKey.PassphraseAlias, false, 1)
	if err != nil {
		return err
	}

	// Decrypt private bytes from the gorm key
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
	s.db.Save(dbPrivateKey)

	return nil
}

// ExportKey is currently unimplemented and will always return an error
func (s *KeyDBStore) ExportKey(name string) ([]byte, error) {
	return nil, errors.New("Exporting from a KeyDBStore is not supported.")
}

// ImportKey is currently unimplemented and will always return an error
func (s *KeyDBStore) ImportKey(pemBytes []byte, alias string) error {
	return errors.New("Importing into a KeyDBStore is not supported")
}

// HealthCheck verifies that DB exists and is query-able
func (s *KeyDBStore) HealthCheck() error {
	dbPrivateKey := GormPrivateKey{}
	tableOk := s.db.HasTable(&dbPrivateKey)
	switch {
	case s.db.Error != nil:
		return s.db.Error
	case !tableOk:
		return fmt.Errorf(
			"Cannot access table: %s", dbPrivateKey.TableName())
	}
	return nil
}
