package main

// The Driver interface encapsulates the necessary behaviour to be a pitfall driver
type Driver interface {
	RevertVMToSnapshot() (string, error)
	CloneVM() (string, error)
	GetIP() (string, error)
}
