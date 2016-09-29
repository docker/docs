package pinataSockets

import (
	appleutil "github.com/docker/pinata/v1/apple/util"
	"path/filepath"
)

const (
	// 9P interface to database
	dbSocketName = "s40" // formerly com.docker.db.socket
)

// GetDBSocketPath returns path to db unix socket
func GetDBSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), dbSocketName)
}
