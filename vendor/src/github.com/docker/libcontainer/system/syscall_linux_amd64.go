// +build linux,amd64
package system

import (
	"syscall"
)

// Setuid sets the uid of the calling thread to the specified uid.
func Setuid(uid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETUID, uintptr(uid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

// Setgid sets the gid of the calling thread to the specified gid.
func Setgid(gid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETGID, uintptr(gid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
