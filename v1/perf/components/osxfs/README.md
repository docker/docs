# osxfs transport latency harness

1. Run pinata
2. Wait until load settles and the VM is fully booted.
3. Execute `killall -USR2 com.docker.osxfs`
4. Run `make` in this directory
5. Open `vsock.png` and `ping.png`

# What it measures

`vsock.png` contains the RTT distribution for the osxfs-transfused link
with 50th and 90th percentiles marked. This measures the latency
overhead of virtualization, hypervisor, AF_VSOCK, and the last-mile UNIX
domain socket on OS X.

`ping.png` contains 4 distributions:

1. "event RTT" is the RTT distribution for an IN_CREATE event (vsock
message) triggered from OS X to reach the first placation stage in osxfs
after traversing the VFS and FUSE as a `symlink` (or `lookup` prelude to
`symlink`) request. This is a vsock roundtrip and a FUSE request.

2. "error RTT" is the RTT distribution for an ENOTRECOVERABLE error to
return to the Moby userspace in response to the `symlink` syscall and
for the failure of the syscall to be logged back to osxfs. This is a
vsock roundtrip and a FUSE reply.

3. "event+error RTT" is the RTT distribution for the temporal pairwise
summation of 1 and 2. This is 2 vsock roundtrips and a FUSE
request/reply pair.

4. "FUSE symlink error RTT" is 3 with 2 median (50th percentile) vsock
roundtrips subtracted as computed from the median element of `vsock.png`.
