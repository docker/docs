---
description: Change log / release notes per release
keywords: pinata, alpha, tutorial
redirect_from:
- /winkit/release-notes/
title: Docker for Windows Release notes
---

Here are the main improvements and issues per release, starting with the current
release. The documentation is always updated for each release.

For system requirements, please see
[What to know before you install](install.md#what-to-know-before-you-install).

Release notes for _stable_ and _beta_ releases are listed below. You can learn
about both kinds of releases, and download stable and beta product installers at
[Download Docker for Windows](install.md#download-docker-for-windows).

## Stable Release Notes

### Docker for Windows 1.13.1, 2017-02-09 (stable)

**Upgrades**

- [Docker 1.13.1](https://github.com/docker/docker/releases/tag/v1.13.1)
- [Docker Compose 1.11.1](https://github.com/docker/compose/releases/tag/1.11.1)
- Linux kernel 4.9.8

**Bug fixes and minor changes**

- Add link to experimental features
- New 1.13 cancellable operations should now be properly handled by the Docker for desktop
- Various typos fixes
- Fix in Hyper-V VM setup (should fix `ObjectNotFound` errors)

### Docker for Windows 1.13.0, 2017-01-19 (stable)

**Upgrades**

- [Docker 1.13.0](https://github.com/docker/docker/releases/tag/v1.13.0)
- [Docker Compose 1.10](https://github.com/docker/compose/releases/tag/1.10.0)
- [Docker Machine 0.9.0](https://github.com/docker/machine/releases/tag/v0.9.0)
- [Notary 0.4.3](https://github.com/docker/notary/releases/tag/v0.4.3)
- Linux kernel 4.9.4

**New**

- Windows containers
- Improved UI for Daemon.json editing
- VHDX file containing images and non-host mounted volumes can be moved
  (using "advanced" tab in the UI)
- Support for arm, aarch64, ppc64le architectures using qemu
- TRIM support for disk (shrinks virtual disk)
- VM's time synchronization is forced after the host wakes from sleep mode
- Docker Experimental mode can be toggled

**Bug fixes and minor changes**

- Improved Proxy UI
- Improvements to Logging and Diagnostics
- About Box is now Copy/Paste enabled
- Improvements in drive sharing code
- Optimized boot process
- Trend Micro Office Scan made the Api proxy think no drive was shared
- Show a link to the virtualization documentation
- Always remove the disk on factory reset
- VPNKit: Improved diagnostics
- VPNKit: Forwarded UDP datagrams should have correct source port numbers
- VPNKit: If one request fails, allow other concurrent requests to succeed.
  For example this allows IPv4 servers to work even if IPv6 is broken.
- VPNKit: Fix bug which could cause the connection tracking to
  underestimate the number of active connections
- VPNKit: add a local cache of DNS responses

### Docker for Windows 1.12.5, 2016-12-20 (stable)

**Upgrades**

- Docker 1.12.5
- Docker Compose 1.9.0

### Skipped Docker for Windows 1.12.4 (stable)

We did not distribute a 1.12.4 stable release

### Docker for Windows 1.12.3, 2016-11-09 (stable)

**New**

- Restore the VM's configuration when it was changed by the user

- Detect firewall configuration that might block the file sharing

- Send more GUI usage statistics to help us improve the product

- The path to HyperV disks is not hardcoded anymore, making the Toolbox import work with non-standard path

- Verify that ALL HyperV features are enabled

- Added Moby console to the logs

- Save the current engine with the other settings

- Notary version 0.4.2 installed

- Reworked the File Sharing dialog and underlying mechanism
  - Pre-fill username
  - Faster and more reliable feedback when the user/password is not valid
  - Better support for domain users
  - Error message in Logs when File Sharing failed for other reasons

**Upgrades**

- Docker 1.12.3
- Linux Kernel 4.4.27
- Docker Machine 0.8.2
- Docker Compose 1.8.1
- aufs 20160912

**Bug fixes and minor changes**

**General**

- Added the settings to the diagnostics

- Make sure we don't use an older Nlog library from the GAC

- Fix a password escaping regression

- Support writing large values to the database, specially for trusted CAs

- Preserve the Powershell stacktraces

- Write OS and Application versions at the top of each log file

- Don't recreate the VM if only the DNS server is set

- The uninstaller now kills the service if it failed to stop it properly

- Improve debug information

**Networking**

- VpnKit is now restarted if it dies

- VpnKit: impose a connection limit to avoid exhausting file descriptors

- VpnKit: handle UDP datagrams larger than 2035 bytes

- VpnKit: reduce the number of file descriptors consumed by DNS

**File sharing**


- Faster mount/unmount of shared drives

- Added a timeout to mounting/unmounting a shared drive

**Hyper-V**

- Make sure invalid "DockerNat" switches are not used

**Moby**

- Increase default ulimit for memlock (fixes [https://github.com/docker/for-mac/issues/801](https://github.com/docker/for-mac/issues/801))

### Docker for Windows 1.12.1, 2016-09-16 (stable)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* To support trusted registry transparently, all trusted CAs (root or intermediate) on the Windows host are automatically copied to Moby

* `Reset Credentials` will also unshare the shared drives

* Logs are now rotated every day

* Support multiple DNS servers

* Added `mfsymlinks` SMB option to support symlinks on bind mounted folder

* Added `nobrl` SMB option to support `sqlite` on bind mounted folders

* Detect outdated versions of Kitematic

**Upgrades**

* Docker 1.12.1
* Docker machine 0.8.1
* Linux kernel 4.4.20
* aufs 20160905

**Bug fixes and minor changes**

**General**

* Uploading a diagnostic now shows a proper status message in the Settings

* Docker will stop asking to import from Toolbox after an upgrade

* Docker can now import from Toolbox just after HyperV is activated

* Added more debug information to the diagnostics

* Sending anonymous statistics shouldn't hang anymore when Mixpanel is not available

* Support newlines in release notes

* Improve error message when Docker daemon is not responding

* The configuration database is now stored in-memory

* Preserve the stacktrace of PowerShell errors

* Display service stacktrace in error windows

**Networking**

* Improve name servers discovery
* VpnKit supports search domains
* VpnKit is now compiled with OCaml 4.03 rather than 4.02.3

**File sharing**

* Set `cifs` version to 3.02

* VnpKit: reduce the number of sockets used by UDP NAT, reduce the probability

* `slirp`: reduce the number of sockets used by UDP NAT, reduce the probability that NAT rules will time out earlier than expected

* Fixed password handling for host file system sharing

**Hyper-V**

* Automatically disable lingering net adapters that prevent Docker from starting or using the network

* Automatically delete duplicated MobyLinuxVMs on a `reset to factory defaults`

* Improved the HyperV detection and activation mechanism

**Moby**

* Fixed Moby Diagnostics and Update Kernel

* Use default `sysfs` settings, transparent huge pages disabled

* `Cgroup` mount to support `systemd` in containers

**Known issues**

* Docker will automatically disable lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Remove stale network adapters](troubleshoot.md#4-remove-stale-network-adapters) under [Networking issues](troubleshoot.md#networking-issues) in Troubleshooting.

### Docker for Windows 1.12.0, 2016-07-28 (stable)

* First stable release

**Components**

* Docker 1.12.0
* Docker Machine 0.8.0
* Docker Compose 1.8.0

## Beta Release Notes

### Docker Community Edition 17.03.0 Release Notes (2017-02-22 17.03.0-ce-rc1)

**New**

- Introduce Docker Community Edition
- Integration with Docker Cloud: control remote Swarms from the local CLI and view your repositories. This feature will be rolled out to all users progressively.

**Upgrades**

- Docker 17.03.0-ce-rc1
- Linux Kernel 4.9.11

**Upgrades**

- VPNKit: Fixed unmarshalling of DNS packets containing pointers to pointers to labels
- Match Hyper-V Integration Services by ID, not name
- Don't consume 100% CPU when the service is stopped
- Log the diagnostic ID when uploading
- Improved Firewall handling: stop listing the rules since it can take a lot of time
- Don't rollback to the previous engine when the desired engine fails to start

### Beta 41 Release Notes (2017-02-07 1.13.1-rc2-beta41)

**Upgrades**

- Docker 1.13.1-rc2
- [Docker Compose 1.11.0-rc1](https://github.com/docker/compose/releases/tag/1.11.0-rc1)
- Linux kernel 4.9.8

**Bug fixes and minor improvements**

- VPNKit: set the Recursion Available bit on DNS responses from the cache
- Don't use port 4222 inside the Linux VM

### Beta 40 Release Notes (2017-01-31 1.13.1-rc1-beta40)

**Upgrades**

- [Docker 1.13.1-rc1](https://github.com/docker/docker/releases/tag/v1.13.1-rc1)
- Linux kernel 4.9.6

**Bug fixes and minor improvements**

- Fix startup error of `ObjectNotFound` in Set-VMFirmware
- Add detailed logs when firewall is configured
- Add a link to the Experimental Features documentation
- Fixed the Copyright in About Dialog
- VPNKit: Avoid diagnostics to capture too much data
- VPNKit: fix a source of occasional packet loss (truncation) on the virtual ethernet link
- Fix negotiation of TimeSync protocol version (via kernel update)

### Beta 39 Release Notes (2017-01-26 1.13.0-beta39)

**Upgrades**

- Linux kernel 4.9.5

**New**

- DNS forwarder ignores responses from malfunctioning servers
- DNS forwarder send all queries in parallel, process results in order
- DNS forwarder includes servers with zones in general searches
- Significantly increased single-stream TCP throughput

**Bug fixes and minor improvements**

- Fix some timeout issues in port forwarding
- Fix for swap not being mounted ([docker/for-win#403](https://github.com/docker/for-win/issues/403))
- Fix aufs xattr delete issue ([docker/docker#30245](https://github.com/docker/docker/issues/30245))

### Beta 38 Release Notes (2017-01-20 1.13.0-beta38)

**Upgrades**

- [Docker 1.13.0](https://github.com/docker/docker/releases/tag/v1.13.0)
- [Docker Compose 1.10](https://github.com/docker/compose/releases/tag/1.10.0)
- [Docker Machine 0.9.0](https://github.com/docker/machine/releases/tag/v0.9.0)
- [Notary 0.4.3](https://github.com/docker/notary/releases/tag/v0.4.3)
- Linux kernel 4.9.4
- qcow-tool 0.7.2

**New**

- The storage location of the Linux volume can now be moved
- Reclaim disk size on reboot
- You can now edit filesharing paths
- Memory can be allocated with 256 MiB steps
- Proxy can now be completely disabled
- Support for arm, aarch64, ppc64le architectures using qemu
- Dedicated preference pane for advanced configuration of the docker daemon (edit daemon.json)
- Docker Experimental mode can be toggled
- Better support for Split DNS VPN configurations
- Use more DNS servers, respect order

**Bug fixes and minor changes**

- You cannot edit settings while docker is restarting
- Support Copy/Paste in About box
- Auto update polling every 24h
- Kernel boots with vsyscall=emulate arg and CONFIG_LEGACY_VSYSCALL is set to NONE in Moby
- Fixed vsock deadlock under heavy write load
- If you opt-out of analytics, you're prompted for approval before a bug report is sent
- Fixed bug where search domain could be read as `DomainName`
- Dedicated preference pane for HTTP proxy settings.
- Dedicated preference pane for CPU & Memory computing resources.
- Privacy settings moved to the general preference pane.
- Fixed an issue where the preference pane disappeared when the welcome whale menu was closed.
- HyperKit: code cleanup and minor fixes
- Improvements to Logging and Diagnostics
- osxfs: switch to libev/kqueue to improve latency
- VPNKit: improvements to DNS handling
- VPNKit: Improved diagnostics
- VPNKit: Forwarded UDP datagrams should have correct source port numbers
- VPNKit: add a local cache of DNS responses
- VPNKit: If one request fails, allow other concurrent requests to succeed.
  For example this allows IPv4 servers to work even if IPv6 is broken.
- VPNKit: Fixed bug which could cause the connection tracking to
  underestimate the number of active connections

### Beta 37 Release Notes (2017-01-16 1.13.0-rc7-beta37)

**Upgrades**

- Docker 1.13.0-rc7
- Notary 0.4.3
- Linux kernel 4.9.3

### Beta 36 Release Notes (2017-01-12 1.13.0-rc6-beta36)

>**Important Note:** Plugins installed using the experimental "managed plugins" feature in Docker 1.12 must be removed/uninstalled before upgrading.

**Upgrades**

- Docker 1.13.0-rc6
- Docker Compose 1.10-rc2
- Linux Kernel 4.9.2

**Bug fixes and minor improvements**

- Uninstall should be more reliable

### Beta 35 Release Notes (2017-01-06 1.13.0-rc5-beta35)

>**Important Note:** Plugins installed using the experimental "managed plugins" feature in Docker 1.12 must be removed/uninstalled before upgrading.

**Upgrades**

- Docker 1.13.0-rc5
- Docker Compose 1.10-rc1

### Beta 34 Release Notes (2016-12-20 1.13.0-rc4-beta34)

**New**

- Basic UI for Daemon.json editing

**Upgrades**

- Docker 1.13.0-rc4
- Linux Kernel 4.8.15

**Bug fixes and minor changes**

- Improved Proxy UI
- Better diagnostics of Windows containers
- Default Experimental/Debug flags are now set on beta for Windows Containers
- Windows Containers Reset to default script improvements
- About Box is now Copy/Paste enabled

### Beta 33 Release Notes (2016-12-15 1.13.0-rc3-beta33)

>**Important Note:** Plugins installed using the experimental "managed plugins" feature in Docker 1.12 must be removed/uninstalled before upgrading.

**New**

- VHDX file containing images and non-host mounted volumes can be moved (using "advanced" tab in the UI)

**Upgrades**

- Linux Kernel 4.8.14

**Bug fixes and minor changes**

- Bugsnag reports should work again
- Fixed a memory leak related to logs and Windows Containers

### Beta 32.1 Release Notes (2016-12-09 1.13.0-rc3-beta32.1)

>**Important Note:** Plugins installed using the experimental "managed plugins" feature in Docker 1.12 must be removed/uninstalled before upgrading.

**Hotfixes**

- Fix Hyper-V Windows Feature detection

**New**

- Windows containers settings panel
- Windows containers: Restart from the settings panel
- Windows containers: Factory default
- Windows containers: Modify Daemon.json
- Windows containers: Proxy settings can be modified
- Support for arm, aarch64, ppc64le architectures using qemu

**Upgrades**

- Docker 1.13.0-rc3
- Docker Machine 0.9.0-rc2
- Linux Kernel 4.8.12

**Bug fixes and minor changes**

- Time drifts between Windows and Linux containers should disapear
- VPNKit: Improved diagnostics
- Improvements in drive sharing code
- Removed the legacy "Disable oplocks" trick for enabling Windows Containers on older insider previews

### Beta 32 Release Notes (2016-12-07 1.13.0-rc3-beta32)

>**Important Note**:
>
>  Plugins installed using the experimental "managed plugins" feature in Docker 1.12 must be removed/uninstalled before upgrading.

**New**

- Windows containers settings panel and options are working. In previous releases, settings were not implemented for [Windows containers
mode](index.md#switch-between-windows-and-linux-containers-beta-feature). (See
[About the Docker Windows containers specific
dialogs](index.md#about-the-docker-windows-containers-specific-dialogs).)
- Windows containers: Restart from the settings panel works
- Windows containers: Factory default
- Windows containers: `Daemon.json` can be modified
- Windows containers: Proxy settings can be modified
- Support for arm, aarch64, ppc64le architectures using qemu

**Upgrades**

- Docker 1.13.0-rc3
- Docker Machine 0.9.0-rc2
- Linux Kernel 4.8.12

**Bug fixes and minor changes**

- Time drifts between Windows and Linux containers should disapear
- VPNKit: Improved diagnostics
- Improvements in drive sharing code
- Removed the legacy "Disable oplocks" trick for enabling Windows Containers on older insider previews

### Beta 31 Release Notes (2016-12-01 1.13.0-rc2-beta31)

**New**

- HTTP/HTTPS proxy settings are used by the Windows Container daemon to pull images
- TRIM support for disk (shrinks virtual disk)
- VM's time synchronization is forced after the host wakes from sleep mode

**Upgrades**

- Docker 1.13.0-rc2
- Dockerd 1.13.0-rc2 (Windows Containers)
- Docker Compose 1.9.0
- Docker Machine 0.9.0-rc1
- Linux kernel 4.8.10

**Bug fixes and minor changes**

- VPNKit: don't permute resource records in responses
- VPNKit: reduced the amount of log spam
- Optimized boot process
- Diagnostics are improved and faster
- Log the error when the GUI fails to initialize
- Trend Micro Office Scan made the API proxy think no drive was shared, fixed
- Show a link to the virtualizaton documentation
- Flush logs to file more often
- Fixed the URL to the SMB/firewall documentation
- Properly remove duplicate firewall rules

### Beta 30 Release Notes (2016-11-10 1.12.3-beta30)

**Upgrades**

- Docker Compose 1.9.0-rc4
- Linux kernel 4.4.30

**Bug fixes and minor changes**

- Optimized disk on stop
- Always remove the disk on factory reset
- Improvements to Logging and Diagnostics

### Beta 29.3 Release Notes (2016-11-02 1.12.3-beta29.3)

**Upgrades**

- Docker Compose 1.9.0-rc2

### Beta 29.2 Release Notes (2016-10-27 1.12.2-beta29.2)

**Hotfixes**

- Upgrade to Docker 1.12.3

### Beta 29.1 Release Notes (2016-10-26 1.12.1-beta29.1)

**Hotfixes**

- Fixed missing `/dev/pty/ptmx`

### Beta 29 Release Notes (2016-10-25 1.12.3-rc1-beta29)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.
>
Windows Container support relies on the Windows 10 container feature, which is
**experimental** at this point.  Windows 10 Pro (1607, build number 14393)
requires update `KB3192366` (soon to be released via Windows Update) to fully
work. Some insider builds may not work.

**New**

- Restore the VM's configuration when it was changed by the user
- Overlay2 is now the default storage driver. After a factory reset overlay2 will automatically be used
- Detect firewall configuration that might block the file sharing
- Send more GUI usage statistics to help us improve the product

**Upgrades**

- Docker 1.12.3-rc1
- Linux Kernel 4.4.27

**Bug fixes and minor changes**

- Faster mount/unmount of shared drives
- Added a timeout to mounting/unmounting a shared drive
- Added the settings to the diagnostics
- Increase default ulimit for memlock (fixes https://github.com/docker/for-mac/issues/801)
- Make sure we don't use an older Nlog library from the GAC


### Beta 28 Release Notes (2016-10-13 1.12.2-rc3-beta28)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.
>
Windows Container support relies on the Windows 10 container feature, which is
**experimental** at this point.  Windows 10 Pro (1607, build number 14393)
requires update `KB3192366` (soon to be released via Windows Update) to fully
work. Some insider builds may not work.

**New**

- Path to HyperV disks in no longer hardcoded, making the Toolbox import work with non-standard path
- Verify that ALL HyperV features are enabled
- Make it clear why user cannot switch to Windows Containers with a tooltip in the systray
- Added Moby console to the logs
- Save the current engine with the other settings
- Notary version 0.4.2 installed


**Upgrades**

- Docker 1.12.2
- Kernel 4.4.24

**Bug fixes and minor changes**

- Fixed a password escaping regression
- Support writing large values to the database, especially for trusted CAs
- VpnKit is now restarted if it dies
- Make sure invalid "DockerNat" switches are not used
- Preserve the Powershell stacktraces
- Write OS and Application versions at the top of each log file

### Beta 27 Release Notes (2016-09-28 1.12.2-rc1-beta27)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* Reworked the File Sharing dialog and underlying mechanism
* Pre-fill username
* Faster and more reliable feedback when the user/password is not valid
* Better support for domain users
* Error message in Logs when File Sharing failed for other reasons

**Upgrades**

* Docker 1.12.2-rc1
* Docker Machine 0.8.2
* Docker Compose 1.8.1
* kernel 4.4.21
* aufs 20160912

**Bug fixes and minor changes**

* Improve the switching between Linux and Windows containers: better errors, more reliable, deal with more edge cases
* Kill lingering dockerd that users might have still around because they played with Windows Containers before
* Don't recreate the VM if only the DNS server is set
* The uninstaller now kills the service if it failed to stop it properly
* Restart VpnKit and DataKit when the processes die
* VpnKit: impose a connection limit to avoid exhausting file descriptors
* VpnKit: handle UDP datagrams larger than 2035 bytes
* VpnKit: reduce the number of file descriptors consumed by DNS
* Improve debug information

### Beta 26 Release Notes (2016-09-14 1.12.1-beta26)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* Basic support for Windows containers. On Windows 10 build >= 14372, a switch in the `systray` icon will change which daemon (Linux or Windows) the Docker CLI talks to

* To support trusted registry transparently, all trusted CAs (root or intermediate) on the Windows host are automatically copied to Moby

* `Reset Credentials` will also unshare the shared drives

* Logs are now rotated every day

**Upgrades**

* Linux kernel 4.4.20
* aufs 20160905

**Bug fixes and minor changes**

* We no longer send the same DNS settings twice to the daemon

* Fixed the lingering net adapters removal on Windows 10 Anniversary Update

* Uploading a diagnostic now shows a proper status message in the Settings

### Beta 25 Release (2016-09-07 1.12.1-beta25)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* Support multiple DNS servers

**Bug fixes and minor changes**

* Improved name servers discovery
* VpnKit supports search domains
* Set CIFS (common internet file system) version to 3.02

**Known issues**

* Only UTF-8 passwords are supported for host filesystem sharing

* Docker will automatically disable lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Remove stale network adapters](troubleshoot.md#4-remove-stale-network-adapters) under [Networking issues](troubleshoot.md#networking-issues) in Troubleshooting.

### Beta 24 Release (2016-08-23 1.12.1-beta24)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**Upgrades**

* Docker 1.12.1
* Docker Machine 0.8.1
* Linux kernel 4.4.19
* aufs 20160822

**Bug fixes and minor changes**

* `slirp`: reduce the number of sockets used by UDP NAT, reduce the probability that NAT rules will time out earlier than expected

**Known issues**

* Only UTF-8 passwords are supported for host filesystem sharing.

* Docker will automatically disable lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Remove stale network adapters](troubleshoot.md#4-remove-stale-network-adapters) under [Networking issues](troubleshoot.md#networking-issues) in Troubleshooting.

### Beta 23 Release (2016-08-16 1.12.1-rc1-beta23)

>**Important Note**:
>
> The auto-update function in Beta 21 will not be able to install this update. To install the latest beta manually if you are still on Beta 21, please download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* Added `mfsymlinks` smb option to support symlinks on bind mounted folder
* Added `nobrl` smb option to support sqlite on bind mounted folders
* Detect outdated versions of Kitematic

**Upgrades**

* Docker 1.12.1-rc1
* Linux kernel 4.4.17
* aufs 20160808

**Bug fixes and minor changes**

* Fixed password handling for host file system sharing
* Automatically disable lingering net adapters that prevent Docker from starting or using the network
* Automatically delete duplicated MobyLinuxVMs on a `reset to factory defaults`
* Docker will stop asking to import from toolbox after an upgrade
* Docker can now import from toolbox just after hyperV is activated
* Fixed Moby Diagnostics and Update Kernel
* Added more debug information to the diagnostics
* Sending anonymous statistics shouldn't hang anymore when Mixpanel is not available
* Improved the HyperV detection and activation mechanism
* VpnKit is now compiled with OCaml 4.03 rather than 4.02.3
* Support newlines in release notes
* Improved error message when docker daemon is not responding
* The configuration database is now stored in-memory
* Preserve the stacktrace of PowerShell errors
* Display service stacktrace in error windows
* Moby: use default sysfs settings, transparent huge pages disabled
* Moby: cgroup mount to support systemd in containers

**Known issues**

* Only UTF-8 passwords are supported for host filesystem sharing
* Docker will automatically disable lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Troubleshooting](troubleshoot.md#networking-issues).

### Beta 22 Release (2016-08-11 1.12.0-beta22)

Unreleased. See Beta 23 for changes.

**Known issues**

* Docker will automatically disable lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Troubleshooting](troubleshoot.md#networking-issues).

### Beta 21 Release (2016-07-28 1.12.0-beta21)

**New**

* Docker for Windows is now available from 2 channels: **stable** and **beta**. New features and bug fixes will go out first in auto-updates to users in the beta channel. Updates to the stable channel are much less frequent and happen in sync with major and minor releases of the Docker engine. Only features that are well-tested and ready for production are added to the stable channel releases. For downloads of both and more information, see the [Getting Started](index.md#download-docker-for-windows).

* Removed the docker host name. Containers with exported ports are reachable via localhost.

* The UI shows whether the user is on beta or stable channel

**Upgrades**

* Docker 1.12.0 with experimental features
* Docker Machine 0.8.0
* Docker Compose 1.8.0

**Bug fixes and minor changes**

* Fixed networking issue when transmitting data to a container via an exposed port.
* Include the sources for qemu-img
* Fixed the migration from toolbox when the user has a space in its login
* Disable the migration from toolbox when hyperV is not yet activated
* More windows can be closed with ESC
* Added the channel to crash reports
* Fixed a path rewriting bug that happens on Windows insider build 14367
* Simplified the MobyLinux.ps1 script

**Known issues**

* Older Kitematic versions are not compatible with Docker for Windows. You need to manually delete the `C:\Program Files\Docker\Kitematic` folder before you click **Open Kitematic...** systray link.

### Beta 20 Release (2016-07-19 1.12.0-rc4-beta20)

**New**

* The UI option to disable port forwarding to `localhost` has been removed

**Bug fixes and minor changes**

* Fixed `docker.sock` permission issues
* Don't check for update when the settings panel opens
* Removed obsolete DNS workaround
* Use the secondary DNS server in more circumstances
* Limit the number of concurrent port forwards to avoid running out of resources
* Store the database as a "bare" git repo to avoid corruption problems

### Beta 19 Release (2016-07-14 1.12.0-rc4-beta19)

**New**

* Added an option to opt-out from sending usage statistics (will be available on the future stable channel)
* New error dialog box to upload crash reports

**Upgrades**

* Docker 1.12.0 RC4
* Docker Compose 1.8.0 RC2
* Docker Machine 0.8.0 RC2
* Linux kernel 4.4.15

**Bug fixes and minor changes**

* `com.docker.slirp`: included the DNS TCP fallback fix, required when UDP responses are truncated
* `docker build/events/logs/stats...` won't leak when interrupted with Ctrl-C
* Disable all buttons on Update Window when a version is downloading

### Beta 18.1 Release (2016-07-07 1.12.0-rc3-beta18.1)

>**Note**: Docker 1.12.0 RC3 release introduces a backward incompatible change from RC2. You can fix this by [recreating or updating your containers](troubleshoot.md#recreate-or-update-your-containers-after-beta-18-upgrade) as described in Troubleshooting.

**Hotfix**

* Fixed issue resulting in error "Hijack is incompatible with use of CloseNotifier", reverts previous fix for `Ctrl-C` during build.

**New**

* Forwarding the ports to localhost is now the default
* Added `http`/`https` proxy configuration to the settings
* The toolbox default machine can be imported on first launch
* Added UI when a crash report is collected and uploaded
* The check for update runs every 6 hours

**Upgrades**

* Docker 1.12.0 RC3

**Bug fixes and minor changes**

* The docker API proxy was failing to deal with 1.12 features (health check for, for example)
* When killing the VM process, ignore when the process is already stopped
* When stopping the VM, always stop the docker proxy
* Prevent the update windows from downloading the `.msi` into `C:\Program Files\Docker`
* All settings should be disabled when Docker is starting. (This regression was introduced in Beta 17)
* VPNKit: Improved scalability as number of network connections increases
* Improve the connection to the database
* Ignore when the shutdown service is not available

### Beta 18 Release (2016-07-06 1.12.0-rc3-beta18)

**New**

* Forwarding the ports to localhost is now the default
* Added `http`/`https` proxy configuration to the settings
* The toolbox default machine can be imported on first launch
* Added UI when a crash report is collected and uploaded
* The check for update runs every 6 hours

**Upgrades**

* Docker 1.12.0 RC3

**Bug fixes and minor changes**

* Interrupting a `docker build` with Ctrl-C will actually stop the build
* The docker API proxy was failing to deal with 1.12 features (health check for, for example)
* When killing the VM process, ignore when the process is already stopped
* When stopping the VM, always stop the docker proxy
* Prevent the update windows from downloading the `.msi` into `C:\Program Files\Docker`
* All settings should be disabled when Docker is starting. (This regression was introduced in Beta 17)
* VPNKit: Improved scalability as number of network connections increases
* Improve the connection to the database
* Ignore when the shutdown service is not available

### Beta 17 Release (2016-06-29 1.12.0-rc2-beta17)

**Upgrades**

* Linux kernel 4.4.14, aufs 20160627

**Bug fixes and minor changes**

* Support users with spaces in their login
* Fix some cases where `dotnet restore` could hang
* Fixed `docker inspect` on an image
* Removed the console from hyper-v manager
* Improved diagnostic for VPN connection and addedlogs for the service port openers
* Improve Moby's boot sequence to adapt to longer boot time when swarm services are running
* Forcefully turn off a VM that won't shut down
* Clicking on a link from the changelog opens a browser
* Fix links to the documentation
* Fix the url to download Kitematic
* Renewed the signing certificates
* Fixed errors with the firewall and the network switch
* Fixed parsing errors in the Powershell script

### Beta 16 Release (2016-06-17 1.12.0-rc2-beta16)

**Upgrades**

* Docker 1.12.0 RC2
* docker-compose 1.8.0 RC1
* docker-machine 0.8.0 RC1
* Alpine 3.4

**Bug fixes and minor changes**

* Fixes to the VPN mode
* Fixed the localhost port forwarding performance issue
* Auto-detect mounted/unmounted drive in the list of shares
	- Changed the name of the application from "DockerforWindows" to "Docker for Windows"
	- Avoid multiple update windows being displayed at the same time

### Beta 15 Release (2016-06-10 1.11.2-beta15)

**New**

* New experimental networking mode, exposing container ports on `localhost`
* New Settings menu to configure sysctl.conf
* New Settings menu to configure http proxies
* The VPN mode setting is removed (VPN mode is now the only supported mode)
* The vSwitch NAT configuration has been removed

**Upgrades**

* Docker 1.11.2
* Linux 4.4.12, aufs 20160530

**Bug fixes and minor changes**

* Moved `Import from toolbox` option to the General Settings
* Increased the timeout to write to the configuration database
* Fixed an issue where sending anonymous stats to Mixpanel made the application stop
* Faster boot time
* All named pipes are now prefixed with the word `docker`
* Full version number is now displayed in the update window
* Default daemon config does not have debug enabled anymore
* More responsive Settings Panel, with new whales also :-)
* Improved logs and debug information

### Beta 14 Release(2016-06-02 1.11.1-beta14)

**New**

* Enabled configuration of the docker daemon (edit `config.json`)
* The VPN mode is enabled by default
* Removed DHCP for VM network configuration
* User configurable NAT prefix and DNS server
* New feedback window to upload diagnostics dialog
* New status indicator in **Settings** window
* VM logs are uploaded with a crash report
* Animated welcome whale

**Bug fixes and minor changes**

* Support non-ASCII characters in passwords
* Fixed unshare a drive operation
* Fixed deserialized of exceptions sent from the service
* If the backend service is not running, the GUI now starts it
* The app no longer complains if the backend service is not running and the user just wants to shut down.


**Known issues**

* Due to limitation in the Windows NAT implementation, co-existence with other NAT prefixes needs to be carefully managed. See [NAT Configuration](troubleshoot.md#nat-configuration) in [Troubleshooting](troubleshoot.md) for more details.

### Beta 13 Release (2016-05-25 1.11.1-beta13)

**New**

This Beta release includes some significant changes:

* Docker communication is over Hyper-V sockets instead of the network
* Experimental VPN mode, also known as `vpnkit`
* Initial support for `datakit` for configuration
* Redesigned Settings panel
* Docker can now be restarted

**Bug fixes and minor changes**

* Support Net adapters with a different name than "vEthernet (DockerNAT)"
* Sharing now has a better support for domain users
* Fixed Toolbox migration (was broken in Beta12)
* Enabling HyperV (was broken in Beta12)
* Fixed error message when invalid labels are passed to `docker run`
* Mixpanel no longer uses roaming App Data
* UI improvements
* Support was added for VMs with other IP addresses out of the `10.0.75.0/24` range
* Improved FAQ

**Known issues**

* Due to limitation in the Windows NAT implementation, co-existence with other NAT prefixes needs to be carefully managed. See [NAT Configuration](troubleshoot.md#nat-configuration) in [Troubleshooting](troubleshoot.md) for more details.

### Beta 12 Release (2016-17-10 1.11.1-beta12)

**New**

* The application is now separated in two parts. A back-end service and a front-end GUI.The front-end GUI no longer asks for elevated access.

**Bug fixes and minor changes**

* Excluded the network drives from the shares list
* Removed the notification when closing the application
* Minor GUI improvements

**Known issues**

* Due to limitation in the Windows NAT implementation, co-existence with other NAT prefixes needs to be carefully managed. See [NAT Configuration](troubleshoot.md#nat-configuration) in [Troubleshooting](troubleshoot.md) for more details.


### Beta 11b Release (2016-05-11 1.11.1-beta11b)

**Hotfixes**

* Fixed an issue with named pipe permissions that prevented Docker from starting

### Beta 11 Release (2016-05-10 1.11.1-beta11)

**New**

* The GUI now runs in non-elevated mode and connects to an elevated Windows service
* Allocate VM memory by 256 MB increments, instead of 1 GB
* Show a meaningful error when the user has an empty password
* Improved [Troubleshooting](troubleshoot.md) page

**Upgrades**

* docker-compose 1.7.1  (see <a href="https://github.com/docker/compose/releases/tag/1.7.1" target="_blank"> changelog</a>)
* Kernel 4.4.9

**Bug fixes and minor changes**

* Report the VM's IP in `docker port`
* Handle passwords with spaces
* Show a clear error message when trying to install on Home editions
* Slower whale animation in the System Tray
* Proxy is restarting itself when it crashes
* DHCP process handles exceptions gracefully
* Moby (Backend) fixes:
  - Fixed `vsock` half closed issue
  - Added NFS support
  - Hostname is now Moby, not Docker
  - Fixes to disk formatting scripts
  - Kernel upgrade to 4.4.9

**Known issues**

* Due to limitation in the Windows NAT implementation, co-existence with other NAT prefixes needs to be carefully managed. See [Troubleshooting](troubleshoot.md) for more details.

* Logs for the windows service are not aggregated with logs from the GUI. This will be fixed in future versions.


## Beta 10 Release (2016-05-03 1.11.0-beta10)

**New**

* Improved Settings panel, allow to configure the VM’s memory and CPUs
* Co-exist with multiple internal Hyper-V switches and improved DHCP handling
* Token validation is now done over HTTPS. This should fix issues with some firewalls and antivirus software.

**Upgrades**

* Docker 1.11.1

**Bug fixes and minor changes**

* Fixed Desktop shortcut name and updated icons
* Preparation to run the backend as service
* Improved logging and Mixpanel events
* Improved code quality
* Improved the build
* New icons

**Known issues**

*  Due to limitation in the Windows NAT implementation, co-existence with other NAT prefixes needs to be carefully managed. See [Troubleshooting](troubleshoot.md) for more details.


### Beta 9 Release (2016-04-26 1.11.0-beta9)

**New**

* Provide one-click dialog to enable Hyper-V
* Report clear underlying Hyper-V errors

**Bug fixes and minor changes**

* Better handling of some networking issues
* Fixed help menu and start menu getting started URLs
* Restored “Docker is Initializing” notification on first run
* Better error messages during authentication
* Improved logging on error conditions
* Improved build and tests

**Known issues**

* If multiple internal Hyper-V switches exist the Moby VM
may not start correctly. We have identified the issue and
are working on a solution.

### Beta 8 Release (2016-04-20 1.11.0-beta8)

**New**

* Auto-update is installed silently, and relaunches the application when it completes
* Uninstaller can be found in Windows menu
* Kitematic can be downloaded from the Dashboard menu

**Bug fixes and minor changes**

* Better UI in the ShareDrive window
* The firewall alert dialog will not come up as often as it was
* Configured MobyLinux VM with a fixed memory of 2GB
* User password is no longer stored on the host-side KVP
* Uninstall shortcut is available in registry

### Beta 7 Release (2016-04-12 1.11.0-beta7)

**New**

  - Multiple drives can be shared
  - New update window
  - Welcome whale

**Upgrades**

* docker 1.11.0-rc5
* docker-machine 0.7.0-rc3
* docker-compose 1.7.0-rc2

**Bug fixes and minor changes**

* Improved networking configuration and error detection: fixed DHCP renewal and rebind issues
* Allow DNS/DHCP processes to restart on bind error
* Less destructive migration from Docker Toolbox
* Improved documentation
* Better error handling: Moby will restart itself if start takes too long.
* Kill proxy and exit docker before a new version is installed
* The application cannot start twice now
* The proxy will stop automatically when the GUI is not running
* Removed existing proxy firewall rules before starting Moby
* The application now collects more and better information on crashes and other issues
* Improved all dialogs and windows
* Added the version to installer's first screen
* Better reset to defaults
* New regression test framework
* The installation MSI is now timestamped
* The Hyper-V install mentions Docker Toolbox only if it is present
* Improved Bugsnag reports: fixed a dependency bug, and added a unique ID to each new report
* Improved the build
* Improved code quality

**Known issues**

- Settings are now serialized in JSON. This install will lose the current settings.

- Docker needs to open ports on the firewall. Sometimes, the user will see a firewall alert dialog. The user should allow the ports to be opened.

- The application was upgraded to 64 bits. The installation path changed to `C:\Program Files\Docker\Docker`. Users might have to close any Powershell/Cmd windows that were already open before the update to get the new `PATH`. In some cases, users may need to log off and on again.

**Bug Fixes**

  - Fixed DHCP renewal and rebind
  - Only mention toolbox on Hyper-V install if it's present
  - The application does not start twice now
  - DNS/DHCP processes are allowed to restart on bind error now
  - Removed the window that opens quickly during bugsnag reports
  - Fixed OS reported by Bugsnag
  - Improved the build
  - Improved code quality

### Beta 6 Release (2016-04-05 1.11.0.1288)

**Enhancements**

- Docs are updated for Beta 6!
- Support roaming: DNS queries are forwarded to the host
- Improved startup times by running a DHCP server on the host
- New settings dialog design
- Support windows paths with -v
- Updated docker CLI and deamon to 1.11.0-rc3
- Updated docker-machine to 0.7.0-rc2
- Updated docker-compose to 1.7.0-rc1
- Now install docker-credential-wincred
- Allow non-root users in containers to create files on volume mounts
- Automatically install HyperV
- The application is now 64bits
- Improved wording in all dialog boxes and error messages
- Removed exit confirmation
- Show clickable URL in the Install HyperV message box
- Dashboard link to Kitematic (as on Mac)
- Moby Kernel updated to 4.4.6
- The registry key was changed to HKLM\SOFTWARE\Docker Inc.\Docker\1.0

**Known issues**

- Migration from Docker Toolbox can fail sometimes. If this happens, the workaround is to restart the application.

- Docker needs to open ports on the firewall, which can activate a firewall alert dialog. Users should allow the ports to be opened.

- The application was upgraded to 64 bits. The installation path changed to `C:\Program Files\Docker\Docker`. If users have Powershell/Cmd windows already open before the update, they might have to close them to catch the new PATH. In some cases, users will need to log off and on again.

**Bug Fixes**

- Kill VMs that cannot be shutdown properly

- Improved the diagnostic information sent with bugsnag reports

- Settings window shows when the drive is shared or not
`C:` drive can be bind mounted with `//c` or `/c`. Used to be `//c/`

- Don't try to submit empty tokens

- Fixed the version shown in the About box

- Fixed a race condition on the logs

- Fixed a race condition on the settings

- Fixed broken links in the documentation

- Replaced `sha1` with actual version in the assemblies

- Don't start the unused agent process

### Beta 5 Release (2016-03-29 1.10.6)

**Enhancements**

* Remove debug console
* Open browser with hyper-v installation instructions
* Added Cloudfront for downloads from Europe
* Capture qemu logs during toolbox upgrades
* Rename alpha distribution channel to beta

**Bug Fixes**

* Fix diagnose section in bugsnag report
* Fix msi version
* Don't truncate Toolbox link

>**Note**: Docker for Windows skipped from Beta 1 to Beta 5 at this point to synch up the version numbering with Docker for Mac, which went into beta cycles a little earlier.

### Beta 1 Release (2016-03-24 1.10.6)

**Enhancements**

- Display the third party licenses
- Display the license agreement
- The application will refuse to start if Hyper-v is not enabled
- Rename `console` to `debug console`
- Remove `machine` from notification
- Open the feedback forum
- Use same MixPanel project for Windows and macOS
- Align MixPanel events with macOS
- Added a script to diagnose problems
- Submit diagnostic with bugsnag reports
- MixPanel heartbeat every hour

**Bug Fixes**

- Accept all versions of Enterprise 10, Pro 10 and Education 10 during installation (Eval, N, ...)
- Fix Linux kernel crashes with certain applications or somesuch
- Fix notifications that are not shown
- Animate the systray whale on reset
- Shorten the enrollment process timeout
- Properly unmount shares when the user un-selects the setting
- Don't install on unsupported builds

## Alpha Release Notes

### Alpha 4 Release (A2016-03-10 1.10.4.0)

- Faster Startup & Shutdown
- Use host DNS parameters
- Enrollment System
- Recreating manually removed vm
- More MixPanel Events
- Various Bug Fixes

### Alpha 3 Release (2016-03-03 1.10.2.14)

**File sharing**

  - Create network share automatically
  - Improve Credentials management
  - Support paths with c and C drive

**Crashes and Analytics**

  - Report crashes with Bugsnag
  - Send analytics through MixPanel

**GUI**

  - Improve layout of About and Settings dialog
  - Improve Updater
  - Link to *Help*
  - Link to *Send Feeback*

**General**

  - Bug fixes

### Alpha 2 Release (2016-02-26 1.10.2.12)

**Installer**

  - Enhancements
  - Auto-update
  - License agreement

**General**

  - Bug fixes

### Alpha 1 Release (2016-02-22 1.10.1.42-1)

**Hypervisor**

  - significant performance improvements

**Security**

  - retrieving Credentials from user

**Filesystem**

  - hot-mounting host filesystem with credential

**General**

  - state management
  - stability, logging
  - bugfixes, eye candies

### Alpha 0 Release (2016-02-09 1.10.0.0-0)

**Hypervision**

  - hyper-v backed virtual machines
  - boots moby in a few seconds
  - installs CLI in `PATH`
  - proxies docker commands to moby

**Filesystem**

  - mounts host filesystem to support `--volume`
  - samba client with a hardcoded password
  - allows live reload

**Networking**

  - live debugging Node.js application
