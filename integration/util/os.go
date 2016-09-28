package util

import (
	"strings"
	"testing"

	"github.com/docker/dhe-deploy/integration/sshclient"
)

type HostOS struct {
	Name        string
	IsSystemd   bool
	SyslogRegex string
	SyslogPath  string
}

var (
	// XXX: we need to handle different versions of ubuntu differently in the future
	// because ubuntu switched to systemd in 15.04
	Ubuntu = HostOS{
		Name:        "ubuntu",
		IsSystemd:   false,
		SyslogRegex: "ERRO",
		SyslogPath:  "/var/log/upstart/docker.log",
	}
	RHEL = HostOS{
		Name:        "rhel_based",
		IsSystemd:   false,
		SyslogRegex: "docker:.+level=error",
		SyslogPath:  "/var/log/messages",
	}
	Arch = HostOS{
		Name:      "arch",
		IsSystemd: true,
	}
	// Sorry, I have a custom kernel
	Viktor = HostOS{
		Name:      "arch",
		IsSystemd: true,
	}
	Boot2Docker = HostOS{
		Name:        "boot2docker",
		IsSystemd:   false,
		SyslogRegex: "level=error",
		SyslogPath:  "/var/lib/boot2docker/docker.log",
	}
)

// DetectHostOS determines the host OS based on the output of uname -a. Systems like redhat and centos will be "rhel_based". Support for new operating systems should go in here first.
func DetectHostOS(t *testing.T, ssh sshclient.SSHClient) *HostOS {
	hostMapping := map[string]*HostOS{
		"Ubuntu":      &Ubuntu,
		"el7":         &RHEL,
		"ARCH":        &Arch,
		"docker-aufs": &Viktor,
		"boot2docker": &Boot2Docker,
	}

	releaseName := Execute(t, ssh, "uname -a", false)

	for supportedReleaseName, hostOS := range hostMapping {
		if strings.Contains(releaseName, supportedReleaseName) {
			return hostOS
		}
	}

	t.Fatal("Compatible host os not found")

	return nil
}
