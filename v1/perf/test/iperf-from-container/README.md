Run iperf from the container to the host
========================================

This test measures average throughput in Gbit/sec to the host from the
container.

Note this test requires the host to have a real kernel network interface
on the interface to the VM, which we don't have if using userspace networking
("slirp")
