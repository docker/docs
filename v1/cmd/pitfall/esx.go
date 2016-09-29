package main

import (
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/twinj/uuid"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

var requiredConfig = []string{esxUserFlag, esxPassFlag, esxHostFlag}

// ESXDriver is a pitfall driver for VMware ESX
type ESXDriver struct {
	Driver
	Config map[string]string
	ctx    context.Context
	c      *govmomi.Client
	folder *object.Folder
	dsr    types.ManagedObjectReference
	rpr    types.ManagedObjectReference
}

// NewESXDriver creates a new ESXDriver
func NewESXDriver(config map[string]string) (*ESXDriver, error) {
	for _, k := range requiredConfig {
		if config[k] == "" {
			return nil, fmt.Errorf("%s not set", k)
		}
	}
	return &ESXDriver{Config: config}, nil
}

// GetIP returns the IP address of a virtual machine
func (d *ESXDriver) GetIP() (string, error) {
	vm, err := d.getVM()
	if err != nil {
		return "", err
	}

	ip, err := vm.WaitForIP(d.ctx)
	return ip, err
}

// RevertVMToSnapshot reverts a VM to the required snapshot
func (d *ESXDriver) RevertVMToSnapshot() (string, error) {
	vm, err := d.getVM()
	if err != nil {
		return "", err
	}

	state, err := vm.PowerState(d.ctx)
	if err != nil {
		return "", err
	}

	var t *object.Task
	if state != "poweredOff" {
		log.Info("Stopping VM...")
		// Stop the VM, ignore error if it's already powered off
		t, err = vm.PowerOff(d.ctx)
		if err != nil {
			return "", err
		}

		// wait for ze result
		_, err = t.WaitForResult(d.ctx, nil)
		if err != nil {
			return "", err
		}
	}

	log.Info("Reverting to Snapshot...")
	// ToDo - Implement RevertCurrentSnapshot instead
	// Use RevertToSnapshot to get specific minor revisions of an OS
	switch d.Config[osFlag] {
	case OSX:
		t, err = vm.RevertToSnapshot(d.ctx, osxSnapshot, true)
		if err != nil {
			return "", err
		}
	case Win:
		t, err = vm.RevertToSnapshot(d.ctx, winSnapshot, true)
		if err != nil {
			return "", err
		}
	}

	if _, err = t.WaitForResult(d.ctx, nil); err != nil {
		return "", err
	}

	log.Info("Starting VM...")
	t, err = vm.PowerOn(d.ctx)
	if err != nil {
		return "", err
	}

	if _, err = t.WaitForResult(d.ctx, nil); err != nil {
		return "", err
	}

	log.Info("Waiting for IP Address...")
	ip, err := vm.WaitForIP(d.ctx)
	if err != nil {
		return "", err
	}

	log.Infof("IP Address: %s", ip)

	return ip, nil
}

// CloneVM - this will only work if we buy a vCenter server - grrr....
func (d *ESXDriver) CloneVM() (string, error) {
	vm, err := d.getVM()
	if err != nil {
		return "", err
	}

	// Check template VM devices
	devices, err := vm.Device(d.ctx)
	if err != nil {
		panic(err)
	}

	// Grab the key of the disk
	var key int32
	for _, d := range devices {
		if devices.Type(d) == "disk" {
			key = d.GetVirtualDevice().Key
		}
	}

	// Generate relocationspec for cloning
	relocSpec := types.VirtualMachineRelocateSpec{
		Datastore:    &d.dsr,
		Pool:         &d.rpr,
		DiskMoveType: "moveAllDiskBackingsAndDisallowSharing",
		Disk: []types.VirtualMachineRelocateSpecDiskLocator{
			{
				Datastore: d.dsr,
				DiskBackingInfo: &types.VirtualDiskFlatVer2BackingInfo{
					DiskMode:        "persistent",
					ThinProvisioned: types.NewBool(true),
					EagerlyScrub:    types.NewBool(false),
				},
				DiskId: key,
			},
		},
	}

	uuid := uuid.NewV4().String()[0:12]
	name := d.Config[osFlag] + uuid

	var vmMo mo.VirtualMachine
	err = vm.Properties(d.ctx, vm.Reference(), []string{"config.guestId", "snapshot"}, &vmMo)

	// make config spec
	configSpec := types.VirtualMachineConfigSpec{
		GuestId:           vmMo.Config.GuestId,
		Name:              name,
		NumCPUs:           2,
		NumCoresPerSocket: 1,
		MemoryMB:          4096,
	}

	// make clone spec
	cloneSpec := types.VirtualMachineCloneSpec{
		Location: relocSpec,
		Template: false,
		Config:   &configSpec,
		PowerOn:  true,
		Snapshot: vmMo.Snapshot.CurrentSnapshot,
	}

	// do some cloning
	task, err := vm.Clone(d.ctx, d.folder, name, cloneSpec)
	if err != nil {
		panic(err)
	}

	if _, err = task.WaitForResult(d.ctx, nil); err != nil {
		panic(err)
	}

	log.Info("Waiting for IP Address...")
	ip, err := vm.WaitForIP(d.ctx)
	if err != nil {
		return "", err
	}

	log.Infof("IP Address: %s", ip)

	return ip, nil

}

func (d *ESXDriver) getVM() (*object.VirtualMachine, error) {
	d.ctx, _ = context.WithCancel(context.Background())

	esxURL := fmt.Sprintf("https://%s:%s@%s:443/sdk",
		d.Config[esxUserFlag],
		d.Config[esxPassFlag],
		d.Config[esxHostFlag])

	u, err := url.Parse(esxURL)
	if err != nil {
		return nil, err
	}
	d.c, err = govmomi.NewClient(d.ctx, u, true)
	if err != nil {
		return nil, err
	}

	f := find.NewFinder(d.c.Client, true)

	// Find one and only dc
	dc, err := f.DefaultDatacenter(d.ctx)
	if err != nil {
		return nil, err
	}

	// Set DC
	f.SetDatacenter(dc)

	// Set Folder
	dcFolders, err := dc.Folders(d.ctx)
	if err != nil {
		panic(err)
	}

	d.folder = dcFolders.VmFolder
	if d.folder == nil {
		return nil, fmt.Errorf("Can't get VM Folder")
	}

	// Find datastore
	ds, err := f.DatastoreOrDefault(d.ctx, d.Config[esxDatastoreFlag])
	if err != nil {
		return nil, err
	}
	d.dsr = ds.Reference()

	// Find HostSystem
	hs, err := f.HostSystemOrDefault(d.ctx, "")
	if err != nil {
		return nil, err
	}

	// FindResourcePool
	rp, err := hs.ResourcePool(d.ctx)
	if err != nil {
		return nil, err
	}
	d.rpr = rp.Reference()

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
		return nil, fmt.Errorf("Unsupported OS")
	}
	if vmName == "" {
		return nil, fmt.Errorf("VM not found for %s %s", testOS, osVersion)
	}
	vm, err := f.VirtualMachine(d.ctx, vmName)
	if err != nil {
		return nil, err
	}
	return vm, nil
}
