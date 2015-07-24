package trustmanager

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/docker/notary/pkg/passphrase"
	"github.com/endophage/gotuf/data"
	"github.com/jinzhu/gorm"
	gojose "github.com/square/go-jose"
)

// KeyDBStore persists and manages private keys on a SQL database
type KeyDBStore struct {
	sync.Mutex
	db         gorm.DB
	passphrase string
	encrypter  gojose.Encrypter
	cachedKeys map[string]data.PrivateKey
}

// GormPrivateKey represents a PrivateKey in the database
type GormPrivateKey struct {
	gorm.Model
	KeyID      string `gorm:"not null;unique_index"`
	Encryption string `gorm:"type:varchar(50);not null"`
	Algorithm  string `gorm:"not null"`
	Public     []byte `gorm:"not null"`
	Private    string `gorm:"not null"`
}

// TableName sets a specific table name for our GormPrivateKey
func (g GormPrivateKey) TableName() string {
	return "private_keys"
}

// NewKeyDBStore returns a new KeyDBStore backed by a SQL database
func NewKeyDBStore(passphraseRetriever passphrase.Retriever, dbType string, dbSQL *sql.DB) (*KeyDBStore, error) {
	cachedKeys := make(map[string]data.PrivateKey)

	// Retreive the passphrase that will be used to encrypt the keys
	passphrase, _, err := passphraseRetriever("", "", false, 0)
	if err != nil {
		return nil, err
	}

	// Setup our encrypted object
	encrypter, err := gojose.NewEncrypter(gojose.A256GCMKW, gojose.A256GCM, []byte(passphrase))
	if err != nil {
		return nil, err
	}

	// Open a connection to our database
	db, _ := gorm.Open(dbType, dbSQL)

	return &KeyDBStore{db: db,
		passphrase: passphrase,
		encrypter:  encrypter,
		cachedKeys: cachedKeys}, nil
}

// AddKey stores the contents of a private key. Both name and alias are ignored,
// we always use Key IDs as name, and don't support aliases
func (s *KeyDBStore) AddKey(name, alias string, privKey data.PrivateKey) error {
	encryptedKey, err := s.encrypter.Encrypt(privKey.Private())
	if err != nil {
		return err
	}

	// Encrypt the private key material
	encryptedPrivKeyStr := encryptedKey.FullSerialize()

	gormPrivKey := GormPrivateKey{
		KeyID:      privKey.ID(),
		Encryption: string(gojose.PBES2_HS512_A256KW),
		Algorithm:  privKey.Algorithm().String(),
		Public:     privKey.Public(),
		Private:    encryptedPrivKeyStr}

	// Add encrypted private key to the database
	s.db.Create(&gormPrivKey)
	// Value will be false if Create suceeds
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
		return nil, "", ErrKeyNotFound{}
	}

	// Decrypt private bytes from the gorm key
	encryptedPrivKeyJWE, err := gojose.ParseEncrypted(dbPrivateKey.Private)
	if err != nil {
		return nil, "", err
	}
	decryptedPrivKeyBytes, err := encryptedPrivKeyJWE.Decrypt([]byte(s.passphrase))
	if err != nil {
		return nil, "", err
	}

	// Create a new PrivateKey with unencrypted bytes
	privKey := data.NewPrivateKey(data.KeyAlgorithm(dbPrivateKey.Algorithm), dbPrivateKey.Public, decryptedPrivKeyBytes)

	// Add the key to cache
	s.cachedKeys[privKey.ID()] = privKey

	return privKey, "", nil
}

// ListKeys always returns nil. This method is here to satisfy the KeyStore interface
func (s *KeyDBStore) ListKeys() []string {
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
		return ErrKeyNotFound{}
	}

	// Delete the key from the database
	s.db.Delete(&dbPrivateKey)

	return nil
}
