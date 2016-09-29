// A small utility to download the diagnostics tarball from Moby via Hyper-V sockets

package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/rneugeba/virtsock/go/hvsock"
)

var (
	vmIdStr  string
	svcIdStr string
	outFile  string
)

func init() {
	flag.StringVar(&vmIdStr, "vmid", "", "VM Id to download from")
	flag.StringVar(&svcIdStr, "svcid", "445BA2CB-E69B-4912-8B42-D7F494D007EA", "Service Id to use")
	flag.StringVar(&outFile, "o", "moby-diag.tar", "File to write output to")
}

func main() {
	log.SetFlags(log.LstdFlags)
	flag.Parse()

	vmid, err := hvsock.GuidFromString(vmIdStr)
	if err != nil {
		log.Fatalln("Can't parse VM ID GUID: ", vmIdStr)
	}
	svcid, err := hvsock.GuidFromString(svcIdStr)
	if err != nil {
		log.Fatalln("Can't parse Service ID GUID: ", svcIdStr)
	}

	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalln("Can't open file", outFile, err)
	}
	defer f.Close()

	sa := hvsock.HypervAddr{VmId: vmid, ServiceId: svcid}
	s, err := hvsock.Dial(sa)
	if err != nil {
		log.Fatalln("Failed to Dial", sa.VmId.String(), sa.ServiceId.String(), ":", err)
	}
	defer s.Close()

	// Copy what we receive from the socket to the file
	// We sometimes get the benign error of the Socket being closed already. Ignore
	_, err = io.Copy(f, s)
	if err != nil && err != hvsock.ErrSocketClosed {
		log.Printf("Failed to copy: %s\n", err)
	}
}
