package osxTasks

import (
	"github.com/Sirupsen/logrus"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"os"
	"syscall"
)

// AcquireTaskManagerLock acquires a lock on the Task Manager
func AcquireTaskManagerLock() {
	lockPath := appleutil.GetContainerPath() + "/task.lock"

	var lockFh *os.File
	// We first try to create the file with O_EXCL which will fail if it already
	// exists, in which case we open it.
	lockFh, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if !os.IsExist(err) {
			// Permission denied? Directory didn't exist?
			logrus.Fatalf("Creating lock file %s: %#v", lockPath, err)
		}
		// The file already exists so simply open it
		lockFh, err = os.Open(lockPath)
		if err != nil {
			logrus.Fatalf("Opening lock file %s: %#v", lockPath, err)
		}
	}
	if err = syscall.Flock(int(lockFh.Fd()), syscall.LOCK_EX); err != nil {
		logrus.Fatalf("Failed to acquire lock %s: %#v", lockPath, err)
	}
	logrus.Println("Acquired task manager lock")
}
