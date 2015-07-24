package storage

import "github.com/endophage/gotuf/data"

// MetaStore holds the methods that are used for a Metadata Store
type MetaStore interface {
	UpdateCurrent(gun string, update MetaUpdate) error
	UpdateMany(gun string, updates []MetaUpdate) error
	GetCurrent(gun, tufRole string) (data []byte, err error)
	Delete(gun string) error
	GetTimestampKey(gun string) (algorithm data.KeyAlgorithm, public []byte, err error)
	SetTimestampKey(gun string, algorithm data.KeyAlgorithm, public []byte) error
}

type PrivateKeyStore interface {
	GetPrivateKey(keyID string) (algorithm data.KeyAlgorithm, public []byte, private []byte, err error)
	SetPrivateKey(keyID string, algorithm data.KeyAlgorithm, public, private []byte) error
}
