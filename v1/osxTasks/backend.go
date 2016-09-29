package osxTasks

import (
	appleutil "github.com/docker/pinata/v1/apple/util"
	"github.com/docker/pinata/v1/pinataSockets"
	"syscall"

	"github.com/Sirupsen/logrus"
)

const (
	//APIServerCommandName is the name of the API server command
	APIServerCommandName = "com.docker.osx.hyperkit.linux"
)

// ListTasks returns a list of Tasks
func ListTasks(bundle string) []Task {
	MacOS := bundle + "/Contents/MacOS/"

	dbListener := ListenUnix(pinataSockets.GetDBSocketPath())
	db := NewTask(MacOS+"com.docker.db", "com.docker.db", []string{
		"--url", "fd:" + dbListener.String(),
		"--git", appleutil.GetContainerPath() + "/database",
	}, uint(4), []ListeningSocket{dbListener}, "")

	osxfsPath := pinataSockets.GetOsxfsSocketPath()
	osxfsListener := ListenUnix(osxfsPath)
	osxfsControlListener := ListenUnix(pinataSockets.GetOsxfsControlSocketPath())
	osxfsVolumeListener := ListenUnix(pinataSockets.GetOsxfsVolumeSocketPath())
	osxfs := NewTask(MacOS+"com.docker.osxfs", "com.docker.osxfs", []string{
		"--address", "fd:" + osxfsListener.String(),
		"--connect", pinataSockets.GetVsockConnectSocketPath(),
		"--control", "fd:" + osxfsControlListener.String(),
		"--volume-control", "fd:" + osxfsVolumeListener.String(),
		"--database", pinataSockets.GetDBSocketPath(),
	}, uint(2), []ListeningSocket{osxfsListener, osxfsControlListener, osxfsVolumeListener},
		"--debug")

	slirpListener := ListenUnix(pinataSockets.GetSlirpSocketPath())
	// launchd registers this one:
	vmnetPath := "/var/tmp/com.docker.vmnetd.socket"

	portListener := ListenUnix(pinataSockets.GetPortSocketPath())
	slirp := NewTask(MacOS+"com.docker.slirp", "com.docker.slirp", []string{
		"--db", pinataSockets.GetDBSocketPath(),
		"--ethernet", "fd:" + slirpListener.String(),
		"--port", "fd:" + portListener.String(),
		"--vsock-path", pinataSockets.GetVsockConnectSocketPath(),
	}, uint(2), []ListeningSocket{slirpListener, portListener},
		"--debug")

	api := NewTask(MacOS+APIServerCommandName, APIServerCommandName, []string{}, uint(3), []ListeningSocket{},
		"-debug")

	dockerListener := ListenUnix(pinataSockets.GetDockerSocketPath())

	xhyve := NewTask(MacOS+"com.docker.driver.amd64-linux", "com.docker.driver.amd64-linux", []string{
		"-db", pinataSockets.GetDBSocketPath(),
		"-osxfs-volume", pinataSockets.GetOsxfsVolumeSocketPath(),
		"-slirp", pinataSockets.GetSlirpSocketPath(),
		"-vmnet", vmnetPath,
		"-port", pinataSockets.GetPortSocketPath(),
		"-vsock", pinataSockets.GetVsockDirPath(),
		"-docker", pinataSockets.GetDockerSocketPath(),
		"-addr", "fd:" + dockerListener.String(), "-debug",
	}, uint(1), []ListeningSocket{dockerListener}, "")

	return []Task{db, osxfs, slirp, api, xhyve}
}

// IncreaseFdLimit is a best-effort increase the maximum number of files per process for this
// process and its children.
func IncreaseFdLimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logrus.Println("Error getting rlimits: ", err)
		return
	}
	rLimit.Max = 10240
	rLimit.Cur = 10240
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logrus.Println("Error setting rlimits: ", err)
		return
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logrus.Println("Error getting rlimits: ", err)
		return
	}
	logrus.Printf("Maximum number of file descriptors is %d", rLimit.Cur)
}
