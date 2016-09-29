package pinataSockets

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"path/filepath"
)

// socket names
const (
	// bind mount notification from docker proxy
	osxFsVolumeSocketName = "s30" // formerly com.docker.osxfs.volume.socket
	// osxfs query and control interface for frontend
	osxfsControlSocketName = "s31"
	// hostnet ethernet frames
	slirpSocketName = "s50" // formerly com.docker.slirp.socket
	// port forwarding control
	portSocketName = "s51" // formerly com.docker.port.socket
	// docker API endpoint created by com.docker.driver.amd64-linux
	dockerSocketName = "s60" // formerly /var/tmp/docker.sock
)

// GetOsxfsVolumeSocketPath returns path to osx fs volume unix socket
func GetOsxfsVolumeSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), osxFsVolumeSocketName)
}

// GetOsxfsControlSocketPath returns the path to the osx control socket
func GetOsxfsControlSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), osxfsControlSocketName)
}

// GetSlirpSocketPath returns path to slirp unix socket
func GetSlirpSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), slirpSocketName)
}

// GetPortSocketPath returns path to port forwarding control unix socket
func GetPortSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), portSocketName)
}

// GetDockerSocketPath returns path to docker API endpoint socket
func GetDockerSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), dockerSocketName)
}

// GetVsockDirPath returns the path to the vsock dir
func GetVsockDirPath() string {
	return appleutil.GetContainerPath()
}

// GetOsxfsSocketPath returns the path to the Osxfs socket
func GetOsxfsSocketPath() string {
	return GetVsockSocketPath(2, 1524)
}

// GetVsockSocketPath returns the path to the vsock socket
func GetVsockSocketPath(callerID uint32, port uint32) string {
	callerIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(callerIDBytes, callerID)
	portBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(portBytes, port)
	socketName := "*" + hex.EncodeToString(callerIDBytes) + "." + hex.EncodeToString(portBytes)
	return filepath.Join(appleutil.GetContainerPath(), socketName)
}

// GetVsockAliasSocketPath returns the path to the vsock alias socket
func GetVsockAliasSocketPath(alias string, port uint32) (string, error) {
	// don't allow aliases with name longer than 6
	if len([]byte(alias)) > 6 {
		return "", errors.New("vsock aliases can't accept alias names longer than 6")
	}
	portBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(portBytes, port)
	socketName := "@" + alias + "." + hex.EncodeToString(portBytes)
	return filepath.Join(appleutil.GetContainerPath(), socketName), nil
}

// GetVsockConnectSocketPath returns the path to the vscok connect socket
func GetVsockConnectSocketPath() string {
	return filepath.Join(appleutil.GetContainerPath(), "@connect")
}
