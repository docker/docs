package store

type ErrMetaNotFound struct{}

func (err ErrMetaNotFound) Error() string {
	return "no trust data available"
}

type ErrKeyNotAvailable struct{}

func (err ErrKeyNotAvailable) Error() string {
	return "could not retrieve timestamp public key"
}
