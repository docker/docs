package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/pinata/v1/apple"
	"github.com/mgutz/str"
	"os/exec"
	"path/filepath"
)

// qemu starts a qemu running Moby, with the given .qcow2 file as a root
// disk and the given amount of memory.
func qemu(memory int, ncpu int, qcow2Path string) {
	vmlinuz := apple.FindKernel()
	ramdisk := apple.FindRamdisk()

	graphics := "-nographic -vga none"
	m := fmt.Sprintf("-m size=%d", memory*1024)
	smp := fmt.Sprintf("-smp cpus=%d", ncpu)
	network := "" // FIXME: need vmnet
	drive := ""   // FIXME: need FUSE/9P
	cdrom := ""
	kernel := fmt.Sprintf("-kernel %s", vmlinuz)
	initrd := fmt.Sprintf("-initrd %s", ramdisk)
	root := fmt.Sprintf("-drive file=\"%s\",if=virtio,media=disk,format=qcow2,index=0", qcow2Path)
	params := "earlyprintk=serial console=ttyAMA0"

	// FIXME: we should have a library function to get the bundle path
	bundle := filepath.Dir(filepath.Dir(filepath.Dir(vmlinuz)))
	qemu := filepath.Join(bundle, "MacOS", "qemu-system-x86_64")
	args := str.ToArgv(graphics + " " + m + " " + smp + " " + network + " " + drive + " " + cdrom + " " + kernel + " " + initrd + " " + root + " " + fmt.Sprintf("-append \"%s\"", params))
	logrus.Printf("exec: %s %#v\n", qemu, args)
	cmd := exec.Command(qemu, args...)
	err := cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}
