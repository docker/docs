package proxy

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
)

// Approver approves shared directories
type Approver interface {
	// Approve(containerID, paths) is called when trying to share directories
	// with the guest. Once this returns (without error), the request is
	// forwarded to Docker.
	Approve(string, []string) error
}

func approveMounts(app Approver, containerID string, hc *container.HostConfig) error {
	if runtime.GOOS == "windows" {
		return nil
	}

	macMounts := make([]string, len(hc.Binds))
	macMountsUsed := 0
	for _, spec := range hc.Binds {
		arr := strings.SplitN(spec, ":", 2)
		if len(arr) > 0 {
			source := arr[0]
			macMounts[macMountsUsed] = source
			macMountsUsed++
		}
	}
	if macMountsUsed == 0 {
		return nil
	}

	return app.Approve(containerID, macMounts[:macMountsUsed])
}

var ctx context.Context

// OSXFSConn is a connection to OSXFS
type OSXFSConn struct {
	forwarder *datakit.Client
	files     map[string]*datakit.File
	m         *sync.Mutex
}

// AddFile adds a file
func (o *OSXFSConn) AddFile(containerID string, file *datakit.File) {
	o.m.Lock()
	defer o.m.Unlock()
	o.files[containerID] = file
}

// DelFile deletes a file
func (o *OSXFSConn) DelFile(containerID string) {
	o.m.Lock()
	defer o.m.Unlock()
	file, ok := o.files[containerID]
	if ok {
		file.Close(ctx)
		delete(o.files, containerID)
	}
}

// Approve approves the supplied paths
func (o *OSXFSConn) Approve(containerID string, hostPaths []string) error {
	// This file:
	filename := containerID
	// will have these contents:
	spec := fmt.Sprintf("%s:%s", containerID, strings.Join(hostPaths, ":"))

	log.Printf("approve %s\n", spec)
	err := o.forwarder.Mkdir(ctx, filename)
	if err != nil {
		log.Printf("Approve failed to create %s: %#v\n", filename, err)
		return err
	}
	ctl, err := o.forwarder.Open(ctx, p9p.OREAD, filename, "ctl")
	if err != nil {
		log.Printf("Approve failed to open %s/ctl: %#v\n", filename, err)
		return err
	}
	// NB the entry will be removed when the fid is finally clunked

	// Read any error from a prevoius session
	bytes := make([]byte, 256)
	n, err := ctl.Read(ctx, bytes, 0)
	if err != nil {
		log.Printf("Approve %s: failed to read response from ctl: %#v\n", spec, err)
		ctl.Close(ctx)
		return err
	}
	_, _ = ctl.Read(ctx, bytes, int64(n))

	response := string(bytes)
	if !strings.HasPrefix(response, "ERROR no request received") {
		log.Printf("Approve %s: read error from previous operation: %s\n", spec, response[0:n])
	}

	request := []byte(spec)
	_, err = ctl.Write(ctx, request, 0)
	if err != nil {
		log.Printf("Approve %s: failed to write to ctl: %#v\n", spec, err)
		ctl.Close(ctx)
		return err
	}

	n, err = ctl.Read(ctx, bytes, 0)
	if err != nil {
		log.Printf("Approve %s: failed to read response from ctl: %#v\n", spec, err)
		ctl.Close(ctx)
		return err
	}

	_, _ = ctl.Read(ctx, bytes, int64(n))
	response = string(bytes)
	if strings.HasPrefix(response, "OK ") {
		log.Printf("Approve %s: succeeded\n", spec)
		response = strings.Trim(response[3:n], " \t\r\n")
		o.AddFile(containerID, ctl)
		return nil
	}

	log.Printf("Approve %s: failed: %s\n", spec, response[0:n])
	if strings.HasPrefix(response, "ERROR ") {
		response = response[6:n]
	}
	ctl.Close(ctx)

	return errors.New(response)
}

// Remove removes the file corresponding to a containerID
func (o *OSXFSConn) Remove(containerID string) error {
	filename := containerID
	log.Printf("remove %s\n", filename)

	o.DelFile(containerID)

	log.Printf("remove %s: succeeded\n", filename)
	return nil
}

// NewOSXFSConn returns a new OSXFSConn
func NewOSXFSConn(osxfsVolumePath string) *OSXFSConn {
	ctx = context.Background()

	for {
		ninep, err := datakit.Dial(ctx, "unix", osxfsVolumePath)
		if err != nil {
			log.Println("Failed to connect to volume forwarding service", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		files := make(map[string]*datakit.File, 0)
		var m sync.Mutex
		return &OSXFSConn{forwarder: ninep, files: files, m: &m}
	}
}
