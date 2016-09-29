package util

import (
	"log"
	"os/user"
	"path/filepath"
)

// GetContainerPath returns the sandbox writable path
func GetContainerPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, "Library", "Containers", "com.docker.docker", "Data")
}
