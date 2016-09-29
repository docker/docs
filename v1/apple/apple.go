// Package apple allows access to OSX APIs in Go.
package apple

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -framework CoreFoundation -framework Foundation -framework ServiceManagement -framework SystemConfiguration -framework CoreServices -framework IOKit

#include "util.h"
#include "bundle.h"
#include "config.h"
#include "paths.h"
#include "power.h"
#include "asl_logger.h"
#include <stdlib.h>
*/
import "C"

// NOTE(aduermael): I'm including vmnetd's utils.h to log using the ASL API
// We should define a better location for C utils that we may be using in
// different processes.

// FindKernel will return a full path to the default kernel
// shipped within the application bundle.
func FindKernel() string {
	p := C.find_kernel()
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s
}

// FindRamdisk will return a full path to the default ramdisk
// shipped within the application bundle.
func FindRamdisk() string {
	p := C.find_ramdisk()
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s
}

// FindTemplate will return a full path to the default template
// shipped within the application bundle.
func FindTemplate() string {
	p := C.find_template()
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s
}

// FindUefi will return a full path to the default uefi
// shipped within the application bundle.
func FindUefi() string {
	p := C.find_uefi()
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s
}

// FindDocuments will return a full path to the directory
// where we can put the boot2docker writable disk image
func FindDocuments() string {
	p := C.find_documents()
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s
}

// FindBundle returns the full path to the bundle directory.
// FIXME: it's not ideal to construct paths manually rather than use
// the OSX path search functions.
func FindBundle() string {
	return path.Dir(path.Dir(path.Dir(path.Dir(FindKernel()))))
}

// PowerEvent is an integer representing an OSX power event
type PowerEvent int

const (
	// GoingToSleep is an event issued by OSX when the machine is going to sleep
	GoingToSleep PowerEvent = iota
	// GoingToWake is an event issued by OSX when the machine is going to wake from sleep
	GoingToWake
)

// ListenForPowerEvents will listen forever for power events
// from OSX
func ListenForPowerEvents(callback func(PowerEvent)) {
	eventR, eventW, err := os.Pipe()
	if err != nil {
		log.Fatalln("Failed to create a pipe")
	}
	defer eventR.Close()
	defer eventW.Close()
	ackR, ackW, err := os.Pipe()
	if err != nil {
		log.Fatalln("Failed to create a pipe")
	}
	defer ackR.Close()
	defer eventW.Close()
	go func() {
		b := make([]byte, 1)
		for {
			n, err := eventR.Read(b)
			if err != nil {
				log.Printf("Error reading power events: %#v. Disabling power notifications.\n", err)
				return
			}
			if n == 0 {
				log.Printf("EOF reading power events. Disabling power notifications.\n")
				return
			}
			switch b[0] {
			case 'S':
				callback(GoingToSleep)
				n, err = ackW.Write(b)
				if err != nil {
					log.Printf("Error writing power event acknowledgement: %#v. Disabling power notifications.\n", err)
					return
				}
				if n == 0 {
					log.Printf("EOF writing power event acknowledgements. Disabling power notifications.\n")
				}
			case 'W':
				callback(GoingToWake)
			default:
				log.Printf("Unknown power event: %c", b[0])
			}
		}
	}()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	C.listen_for_power_events(C.int(eventW.Fd()), C.int(ackR.Fd()))
}

// ConfigChange represents configuration change events
type ConfigChange int

const (
	// DNSChanged events occur when OSX's DNS servers have changed
	DNSChanged = iota
	// ProxiesChanged events occur when OSX's proxy settings have changed
	ProxiesChanged
)

// ListenForConfigChanges will listen forever for config events
// from OSX
func ListenForConfigChanges(callback func(ConfigChange)) {
	eventR, eventW, err := os.Pipe()
	if err != nil {
		log.Fatalln("Failed to create a pipe")
	}
	defer eventW.Close()
	defer eventR.Close()
	go func() {
		b := make([]byte, 1)
		for {
			n, err := eventR.Read(b)
			if err != nil {
				log.Printf("Error reading config changes: %#v\n", err)
				return
			}
			if n == 0 {
				log.Printf("EOF reading config changes\n")
				return
			}
			switch b[0] {
			case 'D':
				callback(DNSChanged)
			case 'P':
				callback(ProxiesChanged)
			default:
				log.Printf("Unknown config change event: %c", b[0])
			}
		}
	}()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	C.listen_for_config_changes(C.int(eventW.Fd()))
}

// GetDNSServers returns the current configured DNS servers
func GetDNSServers() []string {
	maxDNSServers := 10
	a := C.new_char_array(C.int(maxDNSServers))
	defer C.free_char_array(a, C.int(maxDNSServers))
	C.get_dns_servers(a)

	var servers []string
	for i := 0; i < maxDNSServers; i++ {
		p := C.get_array_string(a, C.int(i))
		if p != nil {
			s := C.GoString(p)
			if s != "" {
				servers = append(servers, s)
			}
		}
	}
	return servers
}

// GetDNSSearchDomains returns the current configured DNS search domains
func GetDNSSearchDomains() []string {
	maxDNSSearchDomains := 10
	a := C.new_char_array(C.int(maxDNSSearchDomains))
	defer C.free_char_array(a, C.int(maxDNSSearchDomains))
	C.get_dns_search_domains(a)

	var servers []string
	for i := 0; i < maxDNSSearchDomains; i++ {
		p := C.get_array_string(a, C.int(i))
		if p != nil {
			s := C.GoString(p)
			if s != "" {
				servers = append(servers, s)
			}
		}
	}
	return servers
}

const (
	//HTTPProxy is the name of the HTTP Proxy environment variable
	HTTPProxy = "http_proxy"
	//HTTPSProxy is the name of the HTTPS Proxy environment variable
	HTTPSProxy = "https_proxy"
	//NoProxy is the name of the No Proxy environment variable
	NoProxy = "no_proxy"
)

// GetProxyServers returns the current configured proxy servers
func GetProxyServers() map[string]string {
	numEntries := 3
	a := C.new_char_array(C.int(numEntries))
	defer C.free_char_array(a, C.int(numEntries))
	C.get_proxy_servers(a)

	proxies := make(map[string]string)
	for i := 0; i < numEntries; i++ {
		p := C.get_array_string(a, C.int(i))
		if p != nil {
			s := C.GoString(p)
			switch i {
			case 0:
				proxies[HTTPProxy] = s
			case 1:
				proxies[HTTPSProxy] = s
			case 2:
				proxies[NoProxy] = s
			}
		}
	}
	return proxies
}

// GetLogWriter returns the standard log file location
func GetLogWriter() io.Writer {
	now := time.Now()
	date := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	dir := appleutil.GetContainerPath() + "/logs/" + date
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Failed to create log directory %s: %#v", dir, err)
	}
	file := dir + "/" + path.Base(os.Args[0]) + ".log"
	w, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %#v", file, err)
	}
	return w
}

// ListenOn listens on the given address string
func ListenOn(addr string) net.Listener {
	if strings.HasPrefix(addr, "unix:") {
		addr = addr[5:]
		err := os.Remove(addr)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("os.Remove %s: %#v", addr, err)
		}
		l, err := net.Listen("unix", addr)
		if err != nil {
			log.Fatalln("net.Listen: %#v", err)
		}
		return l
	}
	if strings.HasPrefix(addr, "fd:") {
		f := addr[3:]
		fd, err := strconv.Atoi(f)
		if err != nil {
			log.Fatalf("Failed to parse address %s: %#v", f, err)
		}
		file := os.NewFile(uintptr(fd), addr)
		l, err := net.FileListener(file)
		if err != nil {
			log.Fatalln("Unable to convert fd", fd, "to net.Listener")
		}
		return l
	}
	log.Fatalf("Failed to parse address %s", addr)
	return nil
}

// LogrusASLHook defines a hook for Logrus that redirects logs
// to ASL API (to be displayed in Console application)
type LogrusASLHook struct {
}

// NewLogrusASLHook returns a new LogrusASLHook
func NewLogrusASLHook() *LogrusASLHook {
	hook := new(LogrusASLHook)
	C.apple_asl_logger_init(C.CString("Docker"), C.CString(path.Base(os.Args[0])))
	return hook
}

// Levels returns the available ASL log levels
func (t *LogrusASLHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

// Fire sends a log entry to ASL
func (t *LogrusASLHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.PanicLevel:
		C.apple_asl_logger_log(C.ASL_LEVEL_ALERT, C.CString(entry.Message))
	case logrus.FatalLevel:
		C.apple_asl_logger_log(C.ASL_LEVEL_CRIT, C.CString(entry.Message))
	case logrus.ErrorLevel:
		C.apple_asl_logger_log(C.ASL_LEVEL_ERR, C.CString(entry.Message))
	case logrus.WarnLevel:
		C.apple_asl_logger_log(C.ASL_LEVEL_WARNING, C.CString(entry.Message))
	case logrus.InfoLevel:
		C.apple_asl_logger_log(C.ASL_LEVEL_NOTICE, C.CString(entry.Message))
	case logrus.DebugLevel:
		C.apple_asl_logger_log(C.ASL_LEVEL_DEBUG, C.CString(entry.Message))
	default:
		// unknown level, don't do anything
	}
	return nil
}

// OSXVersion is the Major.Micro.Micro of an OSX host
type OSXVersion struct {
	Major int
	Minor int
	Micro int
}

// GetOSXVersion returns the version of the local host
func GetOSXVersion() (*OSXVersion, error) {
	out, err := exec.Command("/usr/bin/sw_vers", "-productVersion").Output()
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(out))
	bits := strings.Split(trimmed, ".")
	if len(bits) < 2 {
		return nil, errors.New("Failed to parse version: " + trimmed)
	}
	Major, err := strconv.Atoi(bits[0])
	if err != nil {
		return nil, err
	}
	Minor, err := strconv.Atoi(bits[1])
	if err != nil {
		return nil, err
	}
	Micro := 0
	if len(bits) > 2 {
		Micro, err = strconv.Atoi(bits[2])
		if err != nil {
			return nil, err
		}
	}
	return &OSXVersion{Major, Minor, Micro}, nil
}

// String converts an OSXVersion to a human-readable string Major.Minor.Micro
func (v *OSXVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Micro)
}

// GetCertificates returns the .pem files from the keychain
func GetCertificates() (string, error) {
	bytes, err := exec.Command("/usr/bin/security", "find-certificate", "-a", "-p").Output()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
