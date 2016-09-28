package manager

import (
	"testing"
)

func TestParseDriverStatus(t *testing.T) {
	driverStatus := [][2]string{
		[2]string{"Strategy", "spread"},
		[2]string{"Filters", "affinity, health, constraint"},
		[2]string{"Nodes", "1"},
		[2]string{" localhost", "127.0.0.1:2375"},
		[2]string{"  └ Containers", "10"},
		[2]string{"  └ Reserved CPUs", "1 / 4"},
		[2]string{"  └ Reserved Memory", "2 / 8.083GiB"},
		[2]string{"  └ Labels", "executiondriver=native-0.2, kernelversion=3.16.0-4-amd64, operatingsystem=Debian GNU/Linux 8 (jessie), storagedriver=btrfs"},
		[2]string{" remote", "1.2.3.4:2375"},
		[2]string{"  └ Containers", "3"},
		[2]string{"  └ Reserved CPUs", "0 / 4"},
		[2]string{"  └ Reserved Memory", "0 / 8.083GiB"},
		[2]string{"  └ Labels", "executiondriver=native-0.2, kernelversion=3.16.0-4-amd64, operatingsystem=Debian GNU/Linux 8 (jessie), storagedriver=aufs"},
	}

	nodes, err := parseClusterNodes(driverStatus)
	if err != nil {
		t.Fatal(err)
	}

	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes; received %d", len(nodes))
	}

	n1 := nodes[0]
	n2 := nodes[1]

	if n1.Name != "localhost" {
		t.Fatalf("expected name 'localhost'; received %q", n1.Name)
	}

	if n1.Addr != "127.0.0.1:2375" {
		t.Fatalf("expected addr '127.0.0.1:2375'; received %q", n1.Addr)
	}

	if n1.Containers != "10" {
		t.Fatalf("expected 10 containers; received %q", n1.Containers)
	}

	if n1.ReservedCPUs != "1 / 4" {
		t.Fatalf("expected 1 / 4 cpus; received %q", n1.ReservedCPUs)
	}

	if n2.Name != "remote" {
		t.Fatalf("expected name 'remote'; received %q", n2.Name)
	}

	if n2.Addr != "1.2.3.4:2375" {
		t.Fatalf("expected addr '1.2.3.4:2375'; received %q", n2.Addr)
	}

	if n2.Containers != "3" {
		t.Fatalf("expected 3 containers; received %q", n1.Containers)
	}

	if n2.ReservedCPUs != "0 / 4" {
		t.Fatalf("expected 0 / 4 cpus; received %q", n2.ReservedCPUs)
	}

}
