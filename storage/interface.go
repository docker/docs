package storage

// Bootstrapper is a thing that can set itself up
type Bootstrapper interface {
	// Bootstrap instructs a configured Bootstrapper to perform
	// its setup operations.
	Bootstrap() error
}
