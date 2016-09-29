// +build darwin

package main

import (
	"net"

	"github.com/docker/pinata/v1/pinataSockets"
)

// Dial returns a connection for a given unix path
func Dial(path string) (net.Conn, error) {
	return net.Dial("unix", path)
}

// GetDefaultDBPath returns the default database path
func GetDefaultDBPath() string {
	return pinataSockets.GetDBSocketPath()
}
