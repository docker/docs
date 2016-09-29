package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FusionDriver is a pitfall driver for VMware Fusion
type FusionDriver struct {
	Driver
	Config map[string]string
	vm     string
}

func vmrun(args ...string) error {
	path := filepath.Join("/Applications/VMware Fusion.app/Contents/Library/", "vmrun")
	cmd := exec.Command(path, args...)
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = fmt.Errorf("vmrun not found. Please make sure VMware Fusion is installed")
		}
	}
	return err
}

// NewFusionDriver returns a new FusionDriver
func NewFusionDriver(config map[string]string) (*FusionDriver, error) {
	return &FusionDriver{Config: config}, nil
}

// GetIP returns a VM's IP
func (d *FusionDriver) GetIP() (string, error) {
	var err error
	var ip string

	err = d.getVM()
	if err != nil {
		return "", err
	}
	ip, err = d.getIP()
	return ip, err
}

// RevertVMToSnapshot reverts a VM to the required snapshot
func (d *FusionDriver) RevertVMToSnapshot() (string, error) {
	err := d.getVM()
	if err != nil {
		return "", err
	}

	log.Info("Stopping VM...")
	vmrun("stop", d.vm, "hard")

	var snapshot string
	switch d.Config[osFlag] {
	case OSX:
		snapshot = osxSnapshot
	case Win:
		snapshot = winSnapshot
	}

	log.Info("Restoring Snapshot...")
	err = vmrun("revertToSnapshot", d.vm, snapshot)
	if err != nil {
		return "", err
	}

	vmrun("start", d.vm)

	log.Info("Waiting for VM to start...")
	retries := 6
	var ip string
	for i := 0; i < retries; i++ {
		time.Sleep(10 * time.Second)
		ip, err = d.getIP()
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	log.Infof("VM Started. IP Address: %s", ip)
	log.Info("Waiting for VM to boot...")
	time.Sleep(60 * time.Second)
	return ip, err
}

// CloneVM is not supported on this driver
func (d *FusionDriver) CloneVM() (string, error) {
	return "", nil
}

func (d *FusionDriver) getVM() error {
	path := filepath.Join(os.Getenv("HOME"), "Documents/Virtual Machines")
	pathLocalized := filepath.Join(os.Getenv("HOME"), "Documents/Virtual Machines.localized")
	// Get the template VM
	testOS := d.Config[osFlag]
	osVersion := d.Config[osVersionFlag]
	var vmName string

	switch testOS {
	case OSX:
		vmName = osxVMs[osVersion]
	case Win:
		vmName = osxVMs[osVersion]
	default:
		return fmt.Errorf("Unsupported OS")
	}

	var vm string
	vm = fmt.Sprintf("%s.vmwarevm/%s.vmx", filepath.Join(path, vmName), vmName)
	file, err := os.Open(vm)
	if err != nil {
		vm = fmt.Sprintf("%s.vmwarevm/%s.vmx", filepath.Join(pathLocalized, vmName), vmName)
		file, err := os.Open(vm)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	defer file.Close()
	d.vm = vm
	return nil
}

// EVEYYTHING BELOW THIS POINT HAS BEEN BORROWED FROM DOCKER MACHINE fusion_darwin.go

func (d *FusionDriver) getIP() (string, error) {
	// determine MAC address for VM
	macaddr, err := d.getMacAddressFromVmx()
	if err != nil {
		return "", err
	}

	// attempt to find the address in the vmnet configuration
	if ip, err := d.getIPfromVmnetConfiguration(macaddr); err == nil {
		return ip, err
	}

	// address not found in vmnet so look for a DHCP lease
	ip, err := d.getIPfromDHCPLease(macaddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}

// borrowed from machine
func (d *FusionDriver) getMacAddressFromVmx() (string, error) {
	var vmxfh *os.File
	var vmxcontent []byte
	var err error

	if vmxfh, err = os.Open(d.vm); err != nil {
		return "", err
	}
	defer vmxfh.Close()

	if vmxcontent, err = ioutil.ReadAll(vmxfh); err != nil {
		return "", err
	}
	// Look for generatedAddress as we're passing a VMX with addressType = "generated".
	var macaddr string
	vmxparse := regexp.MustCompile(`^ethernet0.generatedAddress\s*=\s*"(.*?)"\s*$`)
	for _, line := range strings.Split(string(vmxcontent), "\n") {
		if matches := vmxparse.FindStringSubmatch(line); matches == nil {
			continue
		} else {
			macaddr = strings.ToLower(matches[1])
		}
	}
	if macaddr == "" {
		return "", fmt.Errorf("couldn't find MAC address in VMX file %s", d.vm)
	}
	return macaddr, nil
}

func (d *FusionDriver) getIPfromVmnetConfiguration(macaddr string) (string, error) {

	// DHCP lease table for NAT vmnet interface
	confFiles, _ := filepath.Glob("/Library/Preferences/VMware Fusion/vmnet*/dhcpd.conf")
	for _, conffile := range confFiles {
		if ipaddr, err := d.getIPfromVmnetConfigurationFile(conffile, macaddr); err == nil {
			return ipaddr, err
		}
	}

	return "", fmt.Errorf("IP not found for MAC %s in vmnet configuration files", macaddr)
}

func (d *FusionDriver) getIPfromVmnetConfigurationFile(conffile, macaddr string) (string, error) {
	var conffh *os.File
	var confcontent []byte

	var currentip string
	var lastipmatch string
	var lastmacmatch string

	var err error

	if conffh, err = os.Open(conffile); err != nil {
		return "", err
	}
	defer conffh.Close()

	if confcontent, err = ioutil.ReadAll(conffh); err != nil {
		return "", err
	}

	// find all occurences of 'host .* { .. }' and extract
	// out of the inner block the MAC and IP addresses

	// key = MAC, value = IP
	m := make(map[string]string)

	// Begin of a host block, that contains the IP, MAC
	hostbegin := regexp.MustCompile(`^host (.+?) {`)
	// End of a host block
	hostend := regexp.MustCompile(`^}`)

	// Get the IP address.
	ip := regexp.MustCompile(`^\s*fixed-address (.+?);$`)
	// Get the MAC address associated.
	mac := regexp.MustCompile(`^\s*hardware ethernet (.+?);$`)

	// we use a block depth so that just in case inner blocks exists
	// we are not being fooled by them
	blockdepth := 0
	for _, line := range strings.Split(string(confcontent), "\n") {

		if matches := hostbegin.FindStringSubmatch(line); matches != nil {
			blockdepth = blockdepth + 1
			continue
		}

		// we are only in intressted in endings if we in a block. Otherwise we will count
		// ending of non host blocks as well
		if matches := hostend.FindStringSubmatch(line); blockdepth > 0 && matches != nil {
			blockdepth = blockdepth - 1

			if blockdepth == 0 {
				// add data
				m[lastmacmatch] = lastipmatch

				// reset all temp var holders
				lastipmatch = ""
				lastmacmatch = ""
			}

			continue
		}

		// only if we are within the first level of a block
		// we are looking for addresses to extract
		if blockdepth == 1 {
			if matches := ip.FindStringSubmatch(line); matches != nil {
				lastipmatch = matches[1]
				continue
			}

			if matches := mac.FindStringSubmatch(line); matches != nil {
				lastmacmatch = strings.ToLower(matches[1])
				continue
			}
		}
	}

	// map is filled to now lets check if we have a MAC associated to an IP
	currentip, ok := m[strings.ToLower(macaddr)]

	if !ok {
		return "", fmt.Errorf("IP not found for MAC %s in vmnet configuration", macaddr)
	}

	return currentip, nil

}

func (d *FusionDriver) getIPfromDHCPLease(macaddr string) (string, error) {

	// DHCP lease table for NAT vmnet interface
	leasesFiles, _ := filepath.Glob("/var/db/vmware/*.leases")
	for _, dhcpfile := range leasesFiles {
		if ipaddr, err := d.getIPfromDHCPLeaseFile(dhcpfile, macaddr); err == nil {
			return ipaddr, err
		}
	}

	return "", fmt.Errorf("IP not found for MAC %s in DHCP leases", macaddr)
}

func (d *FusionDriver) getIPfromDHCPLeaseFile(dhcpfile, macaddr string) (string, error) {

	var dhcpfh *os.File
	var dhcpcontent []byte
	var lastipmatch string
	var currentip string
	var lastleaseendtime time.Time
	var currentleadeendtime time.Time
	var err error

	if dhcpfh, err = os.Open(dhcpfile); err != nil {
		return "", err
	}
	defer dhcpfh.Close()

	if dhcpcontent, err = ioutil.ReadAll(dhcpfh); err != nil {
		return "", err
	}

	// Get the IP from the lease table.
	leaseip := regexp.MustCompile(`^lease (.+?) {$`)
	// Get the lease end date time.
	leaseend := regexp.MustCompile(`^\s*ends \d (.+?);$`)
	// Get the MAC address associated.
	leasemac := regexp.MustCompile(`^\s*hardware ethernet (.+?);$`)

	for _, line := range strings.Split(string(dhcpcontent), "\n") {

		if matches := leaseip.FindStringSubmatch(line); matches != nil {
			lastipmatch = matches[1]
			continue
		}

		if matches := leaseend.FindStringSubmatch(line); matches != nil {
			lastleaseendtime, _ = time.Parse("2006/01/02 15:04:05", matches[1])
			continue
		}

		if matches := leasemac.FindStringSubmatch(line); matches != nil && matches[1] == macaddr && currentleadeendtime.Before(lastleaseendtime) {
			currentip = lastipmatch
			currentleadeendtime = lastleaseendtime
		}
	}

	if currentip == "" {
		return "", fmt.Errorf("IP not found for MAC %s in DHCP leases", macaddr)
	}

	return currentip, nil
}
