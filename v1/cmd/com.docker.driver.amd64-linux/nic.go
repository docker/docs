package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	//Vmnetd sends ethernet to vmnetd running as root
	Vmnetd = iota
	// Vmnet uses the vmnet.framework directly, needs entitlement
	Vmnet = iota
	// Slirp uses the slirp for networking
	Slirp = iota
)

// NIC contains the properties of a network interface
type NIC struct {
	idx  int
	net  int
	uuid string
}

// GetNICs returns a list of NICs for the provided network mode
func GetNICs(mode string, nativePorts bool) []NIC {
	var results []NIC

	switch strings.Trim(mode, " \r\n") {
	case "slirp":
		uuid := getUUID(driverDir + "/nic1.uuid")
		results = append(results, NIC{idx: 0, net: Slirp, uuid: uuid})
	case "hybrid":
		if nativePorts {
			// If in hybrid mode and using localhost for port forwarding, don't
			// use vmnet.framework at all.
			uuid := getUUID(driverDir + "/nic1.uuid")
			results = append(results, NIC{idx: 0, net: Slirp, uuid: uuid})
		} else {
			uuid := getUUID(driverDir + "/nic1.uuid")
			results = append(results, NIC{idx: 0, net: Vmnetd, uuid: uuid})
			uuid = getUUID(driverDir + "/nic2.uuid")
			results = append(results, NIC{idx: 1, net: Slirp, uuid: uuid})
		}
	default:
		uuid := getUUID(driverDir + "/nic1.uuid")
		results = append(results, NIC{idx: 0, net: Vmnetd, uuid: uuid})
	}
	return results
}

// Hyperkit command-line argument describing this NIC
func (n NIC) Hyperkit(slot int) string {
	networkPath := vmnetPath
	if n.net == Slirp {
		networkPath = slirpPath
	}
	macPath := driverDir + fmt.Sprintf("/mac.%d", n.idx)
	virtioDriver := "virtio-vpnkit"
	if n.net == Vmnet {
		virtioDriver = "virtio-net"
	}
	simulateFailure := ""
	if n.net == Vmnetd && vmnetFailure {
		simulateFailure = ",simulate-failure"
	}
	return fmt.Sprintf("%d:0,%s%s,uuid=%s,path=%s,macfile=%s",
		slot, virtioDriver, simulateFailure, n.uuid, networkPath, macPath)
}

// HasFailed returns true if the NIC setup has failed
func (n NIC) HasFailed() bool {
	errorPath := driverDir + fmt.Sprintf("/error.%d", n.idx)
	bytes, err := ioutil.ReadFile(errorPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalf("Error reading NIC error file %s: %#v", errorPath, err)
	}
	log.Printf("NIC %d reports error: %s", n.idx, string(bytes))
	return true
}
