Database Schema
===============

## Version: 8

Summary of changes:

- Removed `insecure-registry`
- Removed `expose-docker-socket`
- Renamed `vmnet-simulate-failure` to `debug/simulate-vmnet-failure`
- Renamed `proxy-verbose` to `debug/proxy-verbose`
- Renamed `native` to `hyperkit` and changed the default for mac to reflect this.
- Added `net/*` keys for static interface configuration in moby
- Added `mount` key for Windows CIFS mounts and OSX mounts.
- Removed `/etc/hostname`
- 20160824 amended `etc/docker/daemon.json` from `{"storage-driver":"aufs","debug":true}`
to `{"debug":true}` without bumping schema as this is a gradual change will not affect users.
- 20160903 amended `etc/docker/daemon.json` to `{}` without bumping schema, Moby now sets the
default if not set.

> Note: Mounts are written as `path`:`mountpoint` on OSX and `drive`:`path`:`mountpoint` on Windows

| Key                     | Description                                        | Default     |
|-------------------------|----------------------------------------------------|-------------|
| **schema-version**      | Current DB Schema version                          |             |
| **filesystem**          | filesystem used by D4M/D4W                         | mac: `osxfs`|
| **hypervisor**          | hypervisor used by D4M/D4M                         | mac: `hyperkit` |
| **memory**              | Amount of memory in GB of the guest VM             | 2           |
| **ncpu**                | Number of CPUs                                     | 0           |
| **network**             | Network mode of D4M/D4W                            | `hybrid`    |
| **mounts**              | Host paths that are mounted in the VM              |             | 
| **etc**                 | contents of etc files for moby                     |             |
| etc/sysctl.conf         | sysctl parameters                                  |             |
| etc/docker/daemon.json  | docker daemon json configuration                   | `{}`        |
| etc/resolv.conf         | dns nameservers                                    |             |
| **hyperkit**            | Hyperkit configuration                             |             |
| hyperkit/boot-protocol  | Boot protocol of native hypervisor                 | direct      |
| hyperkit/port-forwarding | Expose container ports on mac instead of VM       | true        |
| hyperkit/uefi-boot-disk | Uefi boot disk                                     |             |
| hyperkit/on-sleep       | Behaviour on-sleep                                 | `freeze`    |
| **net**                 | Configuration of moby network interfaces           |             |
| net/config              | eth0 configuration. static or not set -> dhcp      |             |
| net/address             | eth0 ip address                                    |             |
| net/netmask             | eth0 subnet mask                                   |             |
| net/gateway             | default gateway                                    |             |
| **proxy**               | Proxy settings                                     |             |
| proxy/detect            | Enable proxy detection or manual                   | detect      |
| proxy/auto              | URI of Pac File                                    |             |
| proxy/exclude           | Excluded paths from proxy lookup                   |             |
| proxy/http              | HTTP Proxy                                         |             |
| proxy/https             | HTTPS Proxy                                        |             |
| **slirp**               | Slirp configuration                                |             |
| slirp/docker            | IP Address on Docker side                          | 192.168.65.2 |
| slirp/host              | IP Address on host side                            | 192.168.65.1 |
| **state**               | Persisted state information                        |             |
| state/last-shutdown-time | Time D4M was last shut down                       |             |
| state/last-startup-time | Time D4M was last started                          |             |
| **debug**               | Vmnet networking settings                          |             |
| simulate-vmnet-failure  | Simulate vmnet failure                             | false       |
| proxy-verbose           | Verbose logging for Docker Proxy                   |             |

> Note: Bold denotes a top-level key

## Version: 7

| Key                     | Description                                        | Default     |
|-------------------------|----------------------------------------------------|-------------|
| **etc**                 | contents of etc files for moby                     |             |
| etc/hostname            | hostname                                           | docker      |
| etc/sysctl.conf         | sysctl parameters                                  |             |
| etc/docker/daemon.json  | docker daemon json configuration                   | `{"storage-driver":"aufs","debug":true}` |
| expose-docker-socket    | expose docker socket to containers *unused?*       | false       |
| filesystem              | filesystem used by D4M/D4W                         | mac: `osxfs`|
| hypervisor              | hypervisor used by D4M/D4M                         | native      |
| insecure-registry       | **Deprecated**                                     |             |
| memory                  | Amount of memory in GB of the guest VM             | 2           |
| **native**              | Native hypervisor configuration                    |             |
| native/boot-protocol    | Boot protocol of native hypervisor                 | direct      |
| native/port-forwarding  | Expose container ports on mac instead of VM        | true        |
| native/uefi-boot-disk   | Uefi boot disk                                     |             |
| ncpu                    | Number of CPUs                                     | 0           |
| network                 | Network mode of D4M/D4W                            | `hybrid`    |
| on-sleep                | Xhyve behaviour on-sleep                           | `freeze`    |
| **proxy**               | Proxy settings                                     |             |
| proxy/exclude           | Excluded paths from proxy lookup                   |             |
| proxy/http              | HTTP Proxy                                         |             |
| proxy/https             | HTTPS Proxy                                        |             |
| proxy-system/exclude    | Excluded paths from proxy lookup                   |             |
| proxy-system/http       | HTTP Proxy                                         |             |
| proxy-system/https      | HTTPS Proxy                                        |             |
| proxy-override/exclude  | Excluded paths from proxy lookup                   |             |
| proxy-override/http     | HTTP Proxy                                         |             |
| proxy-override/https    | HTTPS Proxy                                        |             |
| proxy-verbose           | Verbose logging for Docker Proxy                   |             |
| schema-version          | Current DB Schema version                          |             |
| **slirp**               | Slirp configuration                                |             |
| slirp/docker            | IP Address on Docker side                          | 192.168.65.2 |
| slirp/host              | IP Address on host side                            | 192.168.65.1 |
| **state**               | Persisted state information                        |             |
| state/last-shutdown-time | Time D4M was last shut down                       |             |
| state/last-startup-time | Time D4M was last started                          |             |
| vmnet-simulate-failure  | Simulate vmnet failure                             | false       |
