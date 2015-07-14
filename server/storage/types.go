package storage

// MetaUpdate packages up the fields required to update a TUF record
type MetaUpdate struct {
	Role    string
	Version int
	Data    []byte
}
