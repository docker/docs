package apple

import (
	"fmt"
	"log"
	"testing"
)

func TestBundleQuery(t *testing.T) {
	t.Skip("the binary must be in a bundle for these to work")
	// NB: the binary must be in a bundle for these to work
	kernel := FindKernel()
	ramdisk := FindRamdisk()
	template := FindTemplate()
	fmt.Println("Kernel = ", kernel)
	fmt.Println("Ramdisk = ", ramdisk)
	fmt.Println("Template = ", template)
}

func TestDNSServers(t *testing.T) {
	s := GetDNSServers()
	if len(s) < 1 {
		t.Fatalf("No DNS Servers Found")
	}
	for i, server := range s {
		fmt.Printf("DNS Server %d: %s\n", i, server)
	}
}

func TestProxyServers(t *testing.T) {
	p := GetProxyServers()
	for k, v := range p {
		fmt.Printf("Proxy %s: %s\n", k, v)
	}
}

func TestListenForConfigChanges(t *testing.T) {
	t.Skip("Find a way to send a stop signal to C")
	ListenForConfigChanges(func(change ConfigChange) {
		switch change {
		case DNSChanged:
			fmt.Println("Got dns changed event")
			s := GetDNSServers()
			if len(s) < 1 {
				t.Fatalf("No DNS Servers Found")
			}
			for i, server := range s {
				fmt.Printf("DNS Server %d: %s\n", i, server)
			}
		case ProxiesChanged:
			fmt.Println("Got proxy changed event")
			p := GetProxyServers()
			for k, v := range p {
				fmt.Printf("Proxy %s: %s\n", k, v)
			}
		}
	})
	log.Println("This happends afterwards")
}
