package storage

// MetaStore holds the methods that are used for a Metadata Store
type MetaStore interface {
	UpdateCurrent(gun, role string, version int, data []byte) error
	GetCurrent(gun, tufRole string) (data []byte, err error)
	Delete(gun string) error
	GetTimestampKey(gun string) (cipher string, public []byte, err error)
	SetTimestampKey(gun, cipher string, public []byte) error
}
