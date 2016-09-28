// +build !darwin
package client

import (
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/bootstrap/config"
	version "github.com/hashicorp/go-version"
)

func (c *EngineClient) CheckKernelVersion() error {
	log.Debug("Checking for compatible kernel version")
	uname := &syscall.Utsname{}
	if err := syscall.Uname(uname); err != nil {
		return err
	}

	s := make([]byte, len(uname.Release))
	var lens int
	for ; lens < len(uname.Release); lens++ {
		if uname.Release[lens] == 0 {
			break
		}
		s[lens] = uint8(uname.Release[lens])
	}
	v := string(s[0:lens])

	// Drop the portion after the "-" since it doesn't follow semantic versioning
	kernelVer, err := version.NewVersion(strings.Split(v, "-")[0])
	if err != nil {
		return err
	}
	minVer, err := version.NewVersion(config.MinKernelVersion)
	if err != nil {
		return err
	}
	if kernelVer.LessThan(minVer) {
		log.Warnf("Your kernel is too old.  You may experience problems with UCP.  Consider upgrading to %v or newer", minVer)
		time.Sleep(4 * time.Second) // Make sure the user sees the warning
	}
	log.Debugf("Kernel version %s is compatible", v)
	return nil
}
