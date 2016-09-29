package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func waitForDockerUp() {
	dockerIsUpM.Lock()
	defer dockerIsUpM.Unlock()

	if dockerIsUp {
		return
	}
	path := getDockerLocalPath()

	// Suggested by http://stackoverflow.com/questions/26223839/go-net-http-unix-domain-socket-connection
	unixDial := func(proto, addr string) (conn net.Conn, err error) {
		return net.Dial("unix", path)
	}
	tr := &http.Transport{
		Dial: unixDial,
	}

	// Wait until we get a 200 response from Docker, the result
	// should be a wodge of JSON but don't worry about checking
	// that.
	for {
		client := &http.Client{Transport: tr}
		resp, err := client.Get("http://./info")
		if err != nil {
			logrus.Printf("Docker is not responding: %s: waiting 0.5s", err)
		} else if resp.StatusCode == 200 {
			logrus.Println("Docker is responding")
			dockerIsUp = true
			return
		} else {
			logrus.Printf("Docker is not responding: %s: waiting 0.5s", resp.Status)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func getDockerLocalPath() string {
	return fmt.Sprintf("%s/*%08x.%08x", vsockSocketPath, vsockGuestCID, dockerPort)
}

func getDockerDaemonIP() string {
	if dockerDaemonIP != "" {
		return dockerDaemonIP
	}
	switch network {
	case "slirp":
		dockerDaemonIP = dockerIP // We control the DHCP server
	case "hybrid":
		mac := readMACFile(driverDir + "/mac.0")
		if mac.String() == "de:ad:be:ef:de:ad" {
			// This MAC address means the vmnet interface failed. There is therefore
			// no reachable IP for the VM ports
			logrus.Printf("NIC1 reports MAC %s, which means the vmnet connection failed.", mac)
			logrus.Printf("There is no reachable IP address for this VM.")
			dockerDaemonIP = "None"
		} else {
			dockerDaemonIP = readIPFromLeaseDB(mac)
		}
	default:
		mac := readMACFile(driverDir + "/mac.0")
		dockerDaemonIP = readIPFromLeaseDB(mac)
	}
	return dockerDaemonIP
}

func getDockerDaemonIPFromVM() string {
	return strings.Trim(readFileWait(driverDir+"/ip"), " \r\n\t")
}

func readIPFromLeaseDB(mac net.HardwareAddr) string {
	// parse the /var/db/dhcpd_leases file
	// look for lines like:
	// ip_address=192.168.64.120
	// hw_address=1,2e:89:ac:46:cf:32
	filename := "/var/db/dhcpd_leases"
	// Note the entry will be created lazily
	for {
		leases := readFileWait(filename)
		ip, err := readDHCPLease(strings.NewReader(leases), mac)
		if err == nil {
			return ip
		}
		logrus.Printf("No entry for MAC %s in %s", mac, filename)
		time.Sleep(100 * time.Millisecond)
	}
}

func readDHCPLease(leasesFile io.Reader, mac net.HardwareAddr) (string, error) {
	scanner := bufio.NewScanner(leasesFile)
	scanner.Split(bufio.ScanLines)
	ip := ""
	for scanner.Scan() {
		line := scanner.Text()
		bits := strings.SplitN(line, "=", 2)
		if len(bits) == 2 {
			key := strings.Trim(bits[0], " \r\n\t")
			val := strings.Trim(bits[1], " \r\n\t")
			if strings.Compare(key, "ip_address") == 0 {
				ip = val
			} else if strings.Compare(key, "hw_address") == 0 {
				// of the form n,mac
				bits = strings.SplitN(val, ",", 2)
				if len(bits) == 2 {
					// need to compare without leading zeroes on the bytes
					string := strings.Trim(bits[1], " \r\n\t")
					leaseMAC, err := parseMAC(string)
					if err != nil {
						logrus.Printf("Failed to parse MAC address: '%s': %#v\n", string, err)
						continue
					}
					if bytes.Compare(mac, leaseMAC) == 0 {
						return ip, nil
					}
				}
			}
		}
	}
	return "", errors.New("Failed to discover IP address corresponding to MAC " + mac.String())
}

// net.ParseMAC fails to parse MAC addresses which have bytes with leading
// zeroes removed. Unfortunately /var/db/dhcpd_leases is full of these
func parseMAC(mac string) (net.HardwareAddr, error) {
	bits := []string{}
	for _, bit := range strings.Split(mac, ":") {
		nibbles := []byte(bit)
		if len(nibbles) == 1 {
			bits = append(bits, "0"+bit)
		} else {
			bits = append(bits, bit)
		}
	}
	return net.ParseMAC(strings.Join(bits, ":"))
}

func readMACFile(filename string) net.HardwareAddr {
	text := strings.Trim(readFileWait(filename), " \r\n\t")
	mac, err := parseMAC(text)
	if err != nil {
		logrus.Fatal("Failed to parse VM MAC", text, err)
	}
	return mac
}

func readFileWait(filename string) string {
	retries := 100
	for {
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				time.Sleep(100 * time.Millisecond)
				retries = retries - 1
				if retries == 0 {
					retries = 100
					logrus.Printf("File %s does not exist after 10s", filename)
				}
				continue
			}
			logrus.Fatalln("Failed to read", filename)
		}
		return string(bytes)
	}
}

func removeStaleDockerIPMAC() {
	for _, file := range []string{"ip", "mac", "mac.0", "mac.1", "error.0", "error.1"} {
		filename := driverDir + "/" + file
		err := os.Remove(filename)
		if err != nil && !os.IsNotExist(err) {
			logrus.Fatalln("Failed to remove", filename, err)
		}
	}
}

// Cached IP of the docker daemon
var dockerDaemonIP string

// True when docker has started
var dockerIsUp bool
var dockerIsUpM sync.Mutex

// CloseWriter allows for a write to be closed
type CloseWriter interface {
	CloseWrite() error
}

// CloseReader allows for a Read to be closed
type CloseReader interface {
	CloseRead() error
}
