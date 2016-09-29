package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/pinata/v1/apple"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"github.com/docker/pinata/v1/docker_proxy"
	"github.com/docker/pinata/v1/osxTasks"
	"github.com/docker/pinata/v1/pinataSockets"
	"github.com/docker/pinata/v1/reportError"
	"github.com/twinj/uuid"
	"golang.org/x/net/context"
)

const (
	driver = "com.docker.driver.amd64-linux"
)

var (
	addr               string
	dockerPath         string // path to the docker socket
	detach             bool   // not used (delete once UI argument is removed)
	debug              bool
	ncpu               int
	memory             int
	hypervisor         string
	hostname           string
	sysctl             string
	insecure           string // insecure-registry setting
	daemonConf         string
	driverDir          string // path to the driver directory
	dbName             string // path in database for config
	dbSocketPath       string // path to the com.docker.db socket
	vsockSocketPath    string // path to the vsock directory (equals container folder by default)
	filesystem         string // tag selecting the shared FS implementation
	network            string // either 'native' or 'slirp'
	vmnetFailure       bool   // true if we will simulate a vmnet failure
	slirpPath          string // path to the slirp ethernet service
	vmnetPath          string // path to the vmnet ethernet service
	slirpMode          bool
	portPath           string // path to the port forwarding service
	dockerIP           string // if slirp: IP to give moby
	hostIP             string // if slirp: IP to give host
	nativePorts        bool   // true if exposing ports on native networking
	dockerSocket       bool   // true if we allow exposing the docker socket to containers
	transfusedPort     int    // vsock port for Moby transfused
	osxfsVolumePath    string // path to the com.docker.osxfs volume control service
	proxyVerbose       bool   // true if the docker-proxy will log all data
	nativeBootProtocol string // 'direct' or 'uefi'
	nativeUefiBootDisk string // path to UEFI boot disk image (XXX should come from Bundle)
	diskFlushSync      bool   // true if a flush request should hit the physical disk
	slirpMaxConns      int    // maximum number of connections to accept
	systemHTTPProxy    string // address:port of the http_proxy if set
	systemHTTPSProxy   string // address:port of the https_proxy if set
	systemNoProxy      string // comma separated list of URLs that should be excluded from the proxy
)

const (
	vsockHostCID  = 2
	vsockGuestCID = 3
	dockerPort    = 2376 // TCP and VSock
)

func init() {

	flag.StringVar(&addr, "addr", "unix:"+pinataSockets.GetDockerSocketPath(), "bind addr for docker socket, prefix with unix: for unix socket")
	flag.StringVar(&dockerPath, "docker", pinataSockets.GetDockerSocketPath(), "path of the docker socket")
	flag.BoolVar(&detach, "detach", false, "detach from terminal")
	flag.BoolVar(&debug, "debug", false, "write verbose debug logging to stdout")
	flag.IntVar(&ncpu, "ncpu", 0, "number of CPUs (0 means autodetect)")
	flag.IntVar(&memory, "memory", 2, "memory in GiB")
	flag.StringVar(&hypervisor, "hypervisor", "native", "hypervisor to use")
	flag.StringVar(&hostname, "hostname", "docker", "hostname")
	flag.StringVar(&sysctl, "sysctl", "", "sysctl configuration")
	flag.StringVar(&insecure, "insecure-registry", "", "insecure registry configuration")
	flag.StringVar(&daemonConf, "daemon-conf", `{}`, "json config for docker daemon")
	flag.StringVar(&driverDir, "driver", filepath.Join(appleutil.GetContainerPath(), driver), "path to the driver specific files (e.g. disk image)")
	flag.StringVar(&dbName, "database", driver, "database configuration path")
	flag.StringVar(&dbSocketPath, "db", pinataSockets.GetDBSocketPath(), "path to the database socket")
	flag.StringVar(&filesystem, "fs", "osxfs", "shared filesystem to use (osxfs)")
	flag.StringVar(&vsockSocketPath, "vsock", pinataSockets.GetVsockDirPath(), "path to the vsock directory")
	flag.StringVar(&network, "network", "hybrid", "network technology to use: native, slirp or hybrid")
	flag.StringVar(&slirpPath, "slirp", pinataSockets.GetSlirpSocketPath(), "path to slirp ethernet socket")
	flag.StringVar(&vmnetPath, "vmnet", "/var/tmp/com.docker.vmnetd.socket", "path to vmnet ethernet socket")
	flag.StringVar(&portPath, "port", pinataSockets.GetPortSocketPath(), "path to port forwarding socket")
	flag.StringVar(&dockerIP, "docker-ip", "192.168.65.2", "IP address to assign docker in slirp mode")
	flag.StringVar(&hostIP, "host-ip", "192.168.65.1", "IP address to assign the host in slirp mode")
	flag.BoolVar(&nativePorts, "native-ports", true, "true if we want to expose ports on the Mac when using native networking")
	flag.BoolVar(&dockerSocket, "docker-socket", false, "true if we want to allow the docker socket to be exposed to containers")
	flag.StringVar(&osxfsVolumePath, "osxfs-volume", pinataSockets.GetOsxfsVolumeSocketPath(), "path to the 9P volume control service")
	flag.IntVar(&transfusedPort, "transfuse", 1525, "vsock port for Moby transfused")
	flag.BoolVar(&proxyVerbose, "proxy-verbose", false, "enable verbose docker-proxy debugging")
	flag.StringVar(&nativeBootProtocol, "native-boot-protocol", "direct", "boot protocol for native hypervisor (\"direct\" or \"uefi\")")
	flag.StringVar(&nativeUefiBootDisk, "native-uefi-boot-disk", os.Getenv("HOME")+"/UefiBoot.qcow2", "boot protocol for native hypervisor (\"direct\" or \"uefi\")")
}

// FatalError reports a fatal error
func FatalError(err error, format string, a ...interface{}) {
	reportError.FatalError(err, format, a)
}

func getUUID(filename string) string {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			u := uuid.NewV4().String()
			err = ioutil.WriteFile(filename, []byte(u), 0644)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			logrus.Fatal(err)
		}
	}
	u, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Fatal(err)
	}
	return (string(u))
}

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		logrus.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		logrus.Fatal(err)
	}
}

func detectCPUs() int {
	// autodetect a better value
	n, err := syscall.SysctlUint32("hw.ncpu")
	if err != nil {
		n = 2
	}
	// ignore hyperthreading, do not use all resources
	if n > 1 {
		n = n / 2
	}
	return int(n)
}

func main() {
	logrus.AddHook(apple.NewLogrusASLHook())
	reportError.Initialize()
	flag.Parse()

	ctx := context.Background()

	var err error

	vmlinuz := apple.FindKernel()
	initrd := apple.FindRamdisk()
	template := apple.FindTemplate()
	uefi := apple.FindUefi()
	hyperkit := apple.FindBundle() + "/Contents/MacOS/com.docker.hyperkit"

	err = os.MkdirAll(driverDir, 0755)
	if err != nil {
		FatalError(err, "Creating directory %s", driverDir)
	}
	// Delete any stale console-ring
	err = os.Remove(driverDir + "/console-ring")
	if err != nil && !(os.IsNotExist(err)) {
		FatalError(err, "Deleting previous console-ring %s", driverDir+"/console-ring")
	}
	// Create the log subdirectory for Moby
	err = os.MkdirAll(driverDir+"/log", 0755)
	if err != nil {
		FatalError(err, "Creating directory %s", driverDir+"/log")
	}

	dbSocketDir := filepath.Dir(dbSocketPath)
	err = os.MkdirAll(dbSocketDir, 0755)
	if err != nil {
		FatalError(err, "Creating directory %s", dbSocketDir)
	}

	image := driverDir + "/Docker.qcow2"
	_, err = os.Stat(image)
	if err != nil {
		if os.IsNotExist(err) {
			copyFile(template, image)
		} else {
			FatalError(err, "Looking for %s", image)
		}
	}

	// There is a 108 byte limit on the length of Unix domain socket paths
	// so we can't put them in the Library/Application Support/ directory
	// (symptom: EINVAL from net.Listen)
	ptyPath := driverDir + "/tty"
	consoleRingPath := driverDir + "/console-ring"
	lockPath := driverDir + "/lock"

	var lockFh *os.File
	// We first try to create the file with O_EXCL which will fail if it already
	// exists, in which case we open it.
	lockFh, err = os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if !os.IsExist(err) {
			// Permission denied? Directory didn't exist?
			FatalError(err, "Creating lock file", lockPath)
		}
		// The file already exists so simply open it
		lockFh, err = os.Open(lockPath)
		if err != nil {
			FatalError(err, "Opening lock file", lockPath)
		}
	}
	if err = syscall.Flock(int(lockFh.Fd()), syscall.LOCK_EX); err != nil {
		FatalError(err, "Failed to acquire lock", lockPath, "perhaps another hypervisor is running?")
	}
	logrus.Println("Acquired hypervisor lock")
	// Closing the file would release the lock, so keep holding it until
	// the process quits.

	if ncpu == 0 {
		ncpu = detectCPUs()
	}

	proxies := apple.GetProxyServers()

	dir := driver

	stateDir := "state"
	startKey := "last-start-time"
	shutdownKey := "last-shutdown-time"
	networkKey := "network"
	onSleep := "freeze"

	var client *datakit.Client
	attemptsRemaining := 100
	for attemptsRemaining > 0 {
		client, err = datakit.Dial(ctx, "unix", dbSocketPath)
		if err == nil {
			break
		}
		logrus.Printf("Failed to connect to db: %#v\n", err)
		time.Sleep(100 * time.Millisecond)
		attemptsRemaining = attemptsRemaining - 1
	}
	if client == nil {
		logrus.Fatalln("Failed to connect to the database after 10s")
	}
	go func() {
		waitForDockerUp()
		t, err := datakit.NewTransaction(ctx, client, "master", dir+".startup")
		if err != nil {
			FatalError(err, "NewTransaction")
		}
		if err = t.Write(ctx, []string{dir, stateDir, startKey}, fmt.Sprintf("%d", time.Now().Unix())); err != nil {
			FatalError(err, "transaction.Write")
		}
		if err = t.Commit(ctx, fmt.Sprintf("Docker started %d", time.Now().Unix())); err != nil {
			FatalError(err, "transaction.Commit")
		}
	}()
	// For the sleep workaround, change the default based on OSX version. We should
	// remove this completely for beta 27.
	osxVersion, err := apple.GetOSXVersion()
	if err != nil {
		logrus.Fatalln("Failed to discover OSX version", err)
	}
	defaultSleep := "no not freeze"
	if osxVersion.Minor == 10 || osxVersion.Minor == 11 {
		defaultSleep = "freeze" // Remove for beta 27
	}
	logrus.Printf("OSX version = %s, default value of on-sleep = %s", osxVersion.String(), defaultSleep)

	// Read configuration
	config, err := datakit.NewRecord(ctx, client, "master", []string{dir})
	if err != nil {
		logrus.Fatalln("NewRecord failed", err)
	}
	memoryF := config.IntField("memory", memory)
	ncpuF := config.IntField("ncpu", ncpu)
	hypervisorF := config.StringField("hypervisor", hypervisor)
	hostnameF := config.StringField("etc/hostname", hostname)
	sysctlF := config.StringField("etc/sysctl.conf", sysctl)
	// this field is deprecated
	insecureF := config.StringField("insecure-registry", insecure)
	daemonF := config.StringField("etc/docker/daemon.json", daemonConf)
	proxyVerboseF := config.BoolField("proxy-verbose", proxyVerbose)
	onSleepF := config.StringField("on-sleep", defaultSleep)
	filesystemF := config.StringField("filesystem", filesystem)
	networkF := config.StringField("network", network)
	vmnetFailureF := config.BoolField("vmnet-simulate-failure", vmnetFailure)
	dockerIPF := config.StringField("slirp/docker", dockerIP)
	hostIPF := config.StringField("slirp/host", hostIP)
	nativePortsF := config.BoolField("native/port-forwarding", nativePorts)
	dockerSocketF := config.BoolField("expose-docker-socket", dockerSocket)
	// The settings in use by the VM
	httpProxyF := config.StringField("proxy/http", proxies[apple.HTTPProxy])
	httpsProxyF := config.StringField("proxy/https", proxies[apple.HTTPSProxy])
	noProxyF := config.StringField("proxy/exclude", proxies[apple.NoProxy])
	// The system settings
	systemHTTPProxyF := config.StringField("proxy-system/http", proxies[apple.HTTPProxy])
	systemHTTPSProxyF := config.StringField("proxy-system/https", proxies[apple.HTTPSProxy])
	systemNoProxyF := config.StringField("proxy-system/exclude", proxies[apple.NoProxy])
	// Optional user-supplied overrides
	overrideHTTPProxyF := config.StringRefField("proxy-override/http", nil)
	overrideHTTPSProxyF := config.StringRefField("proxy-override/https", nil)
	overrideNoProxyF := config.StringRefField("proxy-override/exclude", nil)

	nativeBootProtocolF := config.StringField("native/boot-protocol", nativeBootProtocol)
	nativeUefiBootDiskF := config.StringField("native/uefi-boot-disk", nativeUefiBootDisk)

	diskFlushSyncF := config.BoolField("disk/full-sync-on-flush", true)
	_ = config.IntField("slirp/max-connections", 900)

	// sync the current values
	config.Wait(ctx)

	config.Upgrade(ctx, 9)

	var memoryV, ncpuV datakit.Version

	memory, memoryV = memoryF.Get()
	ncpu, ncpuV = ncpuF.Get()

	var hypervisorV, hostnameV, sysctlV, insecureV, daemonV, proxyVerboseV, onSleepV datakit.Version
	var networkV, dockerIPV, hostIPV, nativePortsV datakit.Version
	var httpProxyV, httpsProxyV, noProxyV datakit.Version
	var systemHTTPProxyV, systemHTTPSProxyV, systemNoProxyV datakit.Version
	var overrideHTTPProxyV, overrideHTTPSProxyV, overrideNoProxyV datakit.Version

	var nativeBootProtocolV, nativeUefiBootDiskV datakit.Version
	var diskFlushSyncV datakit.Version

	hypervisor, hypervisorV = hypervisorF.Get()
	hypervisor = strings.Trim(hypervisor, " \n\t")
	logrus.Println("hypervisor:", hypervisor)

	hostname, hostnameV = hostnameF.Get()
	hostname = strings.Trim(hostname, " \n\t")

	sysctl, sysctlV = sysctlF.Get()
	// this field is deprecated
	insecure, insecureV = insecureF.Get()
	daemonConf, daemonV = daemonF.Get()

	proxyVerbose, proxyVerboseV = proxyVerboseF.Get()
	proxy.SetVerboseLogging(proxyVerbose)

	nativeBootProtocol, nativeBootProtocolV = nativeBootProtocolF.Get()
	nativeUefiBootDisk, nativeUefiBootDiskV = nativeUefiBootDiskF.Get()

	diskFlushSync, diskFlushSyncV = diskFlushSyncF.Get()

	_, httpProxyV = httpProxyF.Get()
	_, httpsProxyV = httpsProxyF.Get()
	_, noProxyV = noProxyF.Get()

	systemHTTPProxy, systemHTTPProxyV = systemHTTPProxyF.Get()
	systemHTTPSProxy, systemHTTPSProxyV = systemHTTPSProxyF.Get()
	systemNoProxy, systemNoProxyV = systemNoProxyF.Get()

	_, overrideHTTPProxyV = overrideHTTPProxyF.Get()
	_, overrideHTTPSProxyV = overrideHTTPSProxyF.Get()
	_, overrideNoProxyV = overrideNoProxyF.Get()

	var filesystemV datakit.Version

	filesystem, filesystemV = filesystemF.Get()
	filesystem = strings.Trim(filesystem, " \n\t")
	logrus.Println("filesystem:", filesystem)
	if filesystem == "lofs" {
		// This will force an almost immediate restart of the process.
		err := filesystemF.Set("lofs removed, use osxfs", "osxfs")
		if err != nil {
			logrus.Println("Couldn't set filesystem key to 'osxfs'")
		}
		filesystem = "osxfs"
	}

	onSleep, onSleepV = onSleepF.Get()

	network, networkV = networkF.Get()
	network = strings.Trim(network, " \n\t")
	logrus.Println("network:", network)
	if network == "native" {
		// This will force an almost immediate restart of the process.
		err := networkF.Set("native removed, use hybrid", "hybrid")
		if err != nil {
			logrus.Println("Couldn't set network key to 'hybrid'")
		}
		network = "hybrid"
	}

	slirpMode = strings.Compare(network, "slirp") == 0
	vmnetFailure, _ = vmnetFailureF.Get()

	dockerIP, dockerIPV = dockerIPF.Get()
	dockerIP = strings.Trim(dockerIP, " \n\t")

	hostIP, hostIPV = hostIPF.Get()
	hostIP = strings.Trim(hostIP, " \n\t")

	nativePorts, nativePortsV = nativePortsF.Get()
	dockerSocket, _ = dockerSocketF.Get()

	// Write to c to cleanly shutdown hyperkit and exit the process
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	// Advertise our pid so other people can see it and signal us
	pidFilePath := driverDir + "/pid"
	err = ioutil.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		logrus.Fatalln("Failed to write", pidFilePath, err)
	}

	go func() {
		for {
			rebootNeeded := false
			err := config.Wait(ctx)
			if err != nil {
				logrus.Printf("Error from config.Wait: %s", err)
				rebootNeeded = true
			}
			if memoryF.HasChanged(memoryV) {
				logrus.Println("Memory settings have changed")
				rebootNeeded = true
			}
			if ncpuF.HasChanged(ncpuV) {
				logrus.Println("CPU settings have changed")
				rebootNeeded = true
			}
			if hypervisorF.HasChanged(hypervisorV) {
				logrus.Println("Hypervisor has changed")
				rebootNeeded = true
			}
			if nativeBootProtocolF.HasChanged(nativeBootProtocolV) {
				logrus.Println("Native boot protocol has changed")
				rebootNeeded = true
			}
			if nativeUefiBootDiskF.HasChanged(nativeUefiBootDiskV) {
				logrus.Println("UEFI boot disk has changed")
				rebootNeeded = true
			}
			if diskFlushSyncF.HasChanged(diskFlushSyncV) {
				logrus.Println("Disk flush sync behaviour has changed")
				rebootNeeded = true
			}
			if hostnameF.HasChanged(hostnameV) {
				logrus.Println("Hostname has changed")
				rebootNeeded = true
			}
			if sysctlF.HasChanged(sysctlV) {
				logrus.Println("Sysctl conf has changed")
				rebootNeeded = true
			}

			// this field is deprecated
			if insecureF.HasChanged(insecureV) {
				logrus.Println("Insecure registry has changed")
				rebootNeeded = true
			}

			if daemonF.HasChanged(daemonV) {
				logrus.Println("Daemon config has changed")
				rebootNeeded = true
			}

			if proxyVerboseF.HasChanged(proxyVerboseV) {
				proxyVerbose, proxyVerboseV = proxyVerboseF.Get()
				logrus.Println("Verbose flag has changed to", proxyVerbose)
				proxy.SetVerboseLogging(proxyVerbose)
			}

			if filesystemF.HasChanged(filesystemV) {
				logrus.Println("Filesystem has changed")
				rebootNeeded = true
			}

			if networkF.HasChanged(networkV) || dockerIPF.HasChanged(dockerIPV) || hostIPF.HasChanged(hostIPV) || nativePortsF.HasChanged(nativePortsV) {
				logrus.Println("Network settings have changed")
				rebootNeeded = true
			}

			if onSleepF.HasChanged(onSleepV) {
				logrus.Println("on-sleep has changed")
				onSleep, onSleepV = onSleepF.Get()
			}

			// If either the system proxy settings have changed or if the user-supplied
			// overrides have changed, recompute the new settings to use.
			if systemHTTPProxyF.HasChanged(systemHTTPProxyV) || systemHTTPSProxyF.HasChanged(systemHTTPSProxyV) || systemNoProxyF.HasChanged(systemNoProxyV) || overrideHTTPProxyF.HasChanged(overrideHTTPProxyV) || overrideHTTPSProxyF.HasChanged(overrideHTTPSProxyV) || overrideNoProxyF.HasChanged(overrideNoProxyV) {
				http, _ := systemHTTPProxyF.Get()
				ohttp, _ := overrideHTTPProxyF.Get()
				if ohttp != nil {
					http = *ohttp
				}
				https, _ := systemHTTPSProxyF.Get()
				ohttps, _ := overrideHTTPSProxyF.Get()
				if ohttps != nil {
					https = *ohttps
				}
				no, _ := systemNoProxyF.Get()
				ono, _ := overrideNoProxyF.Get()
				if ono != nil {
					no = *ono
				}
				logrus.Printf("Recomputed effective proxy settings: http=%s https=%s noproxy=%s", http, https, no)
				fields := []*datakit.StringField{httpProxyF, httpsProxyF, noProxyF}
				values := []string{http, https, no}
				config.SetMultiple("Recomputed effective proxy settings", fields, values)
			}
			if httpProxyF.HasChanged(httpProxyV) || httpsProxyF.HasChanged(httpsProxyV) || noProxyF.HasChanged(noProxyV) {
				logrus.Println("Effective proxy settings have changed")
				rebootNeeded = true
			}

			if rebootNeeded {
				logrus.Println("I should reboot")
				// Rely on the process supervisor to restart us
				c <- syscall.SIGTERM
			}
		}
	}()

	removeStaleDockerIPMAC()

	writeCertificatesToDb(ctx, client, dir)

	// Acquire the listening socket (i.e. the user-facing docker socket)
	dockerAPIListeners := []net.Listener{}
	dockerAPIListener := apple.ListenOn(addr)
	dockerAPIListeners = append(dockerAPIListeners, dockerAPIListener)

	adders := []Adder{}     // when containers start
	removers := []Remover{} // when containers die

	// connect to the OSXFS Volume sharing service
	OSXFSVolumes := proxy.NewOSXFSConn(osxfsVolumePath)
	var portRewriter proxy.PortRewriter

	if !slirpMode && !nativePorts {
		portRewriter = proxy.NewIPRewriter(getDockerDaemonIP)
	}

	proxyRewriter := proxy.NewHTTPProxyRewriter()

	for _, dockerAPIListener := range dockerAPIListeners {
		// Connections from listeners get sent to the vsock socket
		go proxy.Serve(dockerAPIListener, getDockerLocalPath(), OSXFSVolumes, portRewriter, proxyRewriter)
	}

	if strings.Compare(hypervisor, "qemu") == 0 {
		qemu(memory, ncpu, image)
		logrus.Println("qemu shutdown with no error")
		os.Exit(0)
	}

	var vsockPorts string
	vsockPorts = fmt.Sprintf("%d", dockerPort)
	if strings.Compare(filesystem, "osxfs") == 0 {
		vsockPorts += fmt.Sprintf(";%d", transfusedPort)
		removers = append(removers, OSXFSVolumes)
	}

	go startWatchingEvents(adders, removers)

	go handleVSyslog()

	logrus.Printf("Hypervisor: %s; BootProtocol: %s; UefiBootDisk: %s", hypervisor, nativeBootProtocol, nativeUefiBootDisk)
	var acpiEnabled bool
	var kernelArg, bootDiskArg string
	switch nativeBootProtocol {
	case "direct":
		acpiEnabled = true
		kernelArg = fmt.Sprintf("kexec,%s,%s,earlyprintk=serial console=ttyS0 com.docker.driver=\"%s\", com.docker.database=\"%s\" ntp=gateway mobyplatform=mac", vmlinuz, initrd, driver, dbName)
	case "uefi":
		acpiEnabled = false
		kernelArg = fmt.Sprintf("bootrom,%s,,", uefi)
		bootDiskArg = nativeUefiBootDisk
	default:
		FatalError(err, "Unknown native boot protocol %s", nativeBootProtocol)
	}

	args := make([]string, 0, 32)
	if acpiEnabled {
		args = append(args, "-A")
	}
	args = append(args, "-m", fmt.Sprintf("%dG", memory))
	args = append(args, "-c", fmt.Sprintf("%d", ncpu))
	args = append(args, "-u")
	args = append(args, "-s", "0:0,hostbridge")
	args = append(args, "-s", "31,lpc")
	nextSlot := 2
	for _, nic := range GetNICs(network, nativePorts) {
		args = append(args, "-s", nic.Hyperkit(nextSlot))
		nextSlot++
	}
	sync := "0"
	if diskFlushSync {
		sync = "1"
	}
	if bootDiskArg != "" {
		args = append(args, "-s", fmt.Sprintf("%d,virtio-blk,file://%s?sync=%s&buffered=1,format=qcow", nextSlot, bootDiskArg, sync))
		nextSlot++
	}
	args = append(args, "-s", fmt.Sprintf("%d,virtio-blk,file://%s?sync=%s&buffered=1,format=qcow", nextSlot, image, sync))
	nextSlot++
	args = append(args, "-s", fmt.Sprintf("%d,virtio-9p,path=%s,tag=db", nextSlot, dbSocketPath))
	nextSlot++
	args = append(args, "-s", fmt.Sprintf("%d,virtio-rnd", nextSlot))
	nextSlot++
	args = append(args, "-s", fmt.Sprintf("%d,virtio-9p,path=%s,tag=port", nextSlot, portPath))
	nextSlot++
	args = append(args, "-s", fmt.Sprintf("%d,virtio-sock,guest_cid=%d,path=%s,guest_forwards=%s",
		nextSlot, vsockGuestCID, vsockSocketPath, vsockPorts))
	nextSlot++
	args = append(args, "-l", fmt.Sprintf("com1,autopty=%s,log=%s", ptyPath, consoleRingPath))
	args = append(args, "-f", kernelArg)

	hypervisorPidfile := driverDir + "/hypervisor.pid"
	args = append(args, "-F", hypervisorPidfile)
	err = os.Remove(hypervisorPidfile) // hyperkit requires it to not already exist

	cmd := exec.Command(hyperkit, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Fatalf("creating a stdout pipe: %#v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Fatalf("creating a stderr pipe: %#v", err)
	}
	if err := cmd.Start(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Printf("Launched[%d]: %s %s", cmd.Process.Pid, hyperkit, strings.Join(args, " "))
	// Ensure that the main process supervisor will kill any leaked processes
	task := osxTasks.NewTask(hyperkit, "com.docker.hyperkit", args, uint(1), []osxTasks.ListeningSocket{}, "")
	osxTasks.WritePidFile(task, cmd.Process.Pid)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logrus.Println(scanner.Text())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logrus.Println(scanner.Text())
		}
	}()
	// If our subprocess quits, so do we: our process supervisor can
	// restart us.
	go func() {
		err := cmd.Wait()
		if err != nil && err.(*exec.ExitError) == nil {
			logrus.Fatalln("Failed to Wait() for hypervisor subprocess", err)
		}
		t, err := datakit.NewTransaction(ctx, client, "master", dir+".shutdown")
		if err != nil {
			FatalError(err, "NewTransaction")
		}
		if err = t.Write(ctx, []string{dir, stateDir, shutdownKey}, fmt.Sprintf("%d", time.Now().Unix())); err != nil {
			FatalError(err, "transaction.Write")
		}
		if err = t.Commit(ctx, fmt.Sprintf("Stopped at %d", time.Now().Unix())); err != nil {
			FatalError(err, "transaction.Commit")
		}

		// Check if any of the NICs failed to setup
		anyFailed := false
		for _, nic := range GetNICs(network, nativePorts) {
			anyFailed = anyFailed || nic.HasFailed()
		}
		// We switch to hybrid mode which has 2 NICs and should work more reliably.
		// Note we still rely on vmnet for docker.local so this aspect will fail.
		if anyFailed && network != "hybrid" {
			logrus.Printf("A NIC failed, so switching to hybrid networking mode")
			t, err := datakit.NewTransaction(ctx, client, "master", dir+".hybid")
			if err != nil {
				logrus.Fatalf("Failed to open a new transaction: %#v", err)
			}
			if err = t.Write(ctx, []string{dir, networkKey}, "hybrid"); err != nil {
				logrus.Fatalf("Failed to write %s = %s: %#v", networkKey, "hybrid", err)
			}
			if err = t.Commit(ctx, "Switching to hybrid mode"); err != nil {
				logrus.Fatalf("Failed to commit transaction: %#v", err)
			}
		}
		os.Exit(0)
	}()

	go func() {
		<-c
		// propagate the SIGTERM to our child
		logrus.Println("sending SIGTERM to com.docker.hyperkit pid", cmd.Process.Pid)
		syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}()

	// Synchronise once at program start, and then again when we receive
	// DNSChanged events
	refreshDNSServers(ctx, client)

	go apple.ListenForConfigChanges(func(change apple.ConfigChange) {
		switch change {
		case apple.DNSChanged:
			refreshDNSServers(ctx, client)
		case apple.ProxiesChanged:
			proxies := apple.GetProxyServers()
			commitMsg := fmt.Sprintf("Settings Changed %s", time.Now().Format(time.RFC822Z))
			haveUpdatedProxies := false

			var f []*datakit.StringField
			var v []string

			p1, err := parseHTTPProxy(proxies[apple.HTTPProxy])
			if err != nil {
				logrus.Warn(err)
			}

			if systemHTTPProxy != p1 {
				f = append(f, systemHTTPProxyF)
				v = append(v, p1)
				haveUpdatedProxies = true
			}

			p2, err := parseHTTPProxy(proxies[apple.HTTPSProxy])
			if err != nil {
				logrus.Warn(err)
			}

			if systemHTTPSProxy != p2 {
				f = append(f, systemHTTPSProxyF)
				v = append(v, p2)
				haveUpdatedProxies = true
			}

			if haveUpdatedProxies {
				// If both HTTP Proxy and HTTPS Proxy are unset, clear the exclusions
				if p1 == "" && p2 == "" {
					f = append(f, systemNoProxyF)
					v = append(v, "")
				}
				e := parseNoProxy(proxies[apple.NoProxy])
				if systemNoProxy != e {
					f = append(f, systemNoProxyF)
					v = append(v, e)
				}
			}

			if len(f) > 0 {
				err = config.SetMultiple(commitMsg, f, v)
				if err != nil {
					logrus.Warn(err)
				}
			}
		}
	})

	apple.ListenForPowerEvents(func(event apple.PowerEvent) {
		switch event {
		case apple.GoingToSleep:
			logrus.Println("System wants to go to sleep")
			policy := strings.Trim(onSleep, " \n\t")

			if policy == "freeze" {
				if cmd.Process != nil {
					logrus.Println("Asking com.docker.hyperkit to freeze vcpus")
					syscall.Kill(cmd.Process.Pid, syscall.SIGUSR1)
					time.Sleep(time.Millisecond * 100) // TODO We should ideally get a signal back from hyperkit when ready
					logrus.Println("vcpu freeze complete: allowing sleep")
					return
				}
				logrus.Println("There is no running com.docker.hyperkit process")
				return
			}
			logrus.Printf("com.docker.hyperkit is running and not shutting down because on-sleep = %s\n", policy)
			return
		case apple.GoingToWake:
			logrus.Println("System wants to wake up")
			policy := strings.Trim(onSleep, " \n\t")
			if policy == "freeze" {
				logrus.Println("Asking com.docker.hyperkit to thaw vcpus")
				syscall.Kill(cmd.Process.Pid, syscall.SIGUSR2)
			} else {
				logrus.Println("com.docker.hyperkit is already running so nothing to do")
			}
			return
		}
	})
}

func writeCertificatesToDb(ctx context.Context, client *datakit.Client, dir string) {
	certs, err := apple.GetCertificates()
	if err != nil {
		logrus.Printf("Failed to read certificates: %s", err)
		return
	}
	t, err := datakit.NewTransaction(ctx, client, "master", dir+".certificates")
	if err != nil {
		logrus.Printf("Failed to create a transaction: %s", err)
		return
	}
	if err = t.Write(ctx, []string{dir, "etc", "ssl", "certs", "ca-certificates.crt"}, certs); err != nil {
		logrus.Printf("Failed to write certificates: %s", err)
		return
	}
	if err = t.Commit(ctx, "Setting certificates"); err != nil {
		logrus.Printf("Failed to commit certificates: %s", err)
		return
	}
}
