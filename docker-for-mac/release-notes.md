---
aliases:
- /mackit/release-notes/
description: Change log / release notes per release
keywords:
- pinata, alpha, tutorial
menu:
  main:
    identifier: docker-mac-relnotes
    parent: pinata_mac_menu
    weight: 10
title: Release Notes
---

# Docker for Mac Release Notes

Here are the main improvements and issues per release, starting with the current release. The documentation is always updated for each release.

For system requirements, please see the Getting Started topic on [What to know before you install](index.md#what-to-know-before-you-install).

In the list below, release notes for stable releases start with "_Docker for Mac.._" and are further identified with **(stable)** after the topic title. You can download and run _stable_ and _beta_ product installers at [Download Docker for Mac](index.md#download-docker-for-mac).

## Beta 22 Release Notes (2016-08-11 1.12.0-beta22)

**Upgrades**

*  Linux kernel to 4.4.16

**Bug fixes and minor changes**

* Increase Moby fs.file-max to 524288
* Use Mac System Configuration database to detect DNS
* HyperKit updated with dtrace support and lock fixes
* Fix Moby Diagnostics and Update Kernel
* UI Fixes
* osxfs: fix socket chowns

**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app

* There are a number of issues with the performance of directories bind-mounted with `osxfs`. In particular, writes of small blocks and traversals of large directories are currently slow. Additionally, containers that perform large numbers of directory operations, such as repeated scans of large directory trees, may suffer from poor performance. More information is available in [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

* Under some unhandled error conditions, inotify event delivery can fail and become permanently disabled. The workaround is to restart Docker.app

## Beta 21.1 Release Notes (2016-08-03 1.12.0-beta21.1)

**Hotfixes**

* osxfs: fixed an issue causing access to children of renamed directories to fail (symptoms: npm failures, apt-get failures) (docker/for-mac)

* osxfs: fixed an issue causing some ATTRIB and CREATE inotify events to fail delivery and other inotify events to stop

* osxfs: fixed an issue causing all inotify events to stop when an ancestor directory of a mounted directory was mounted

* osxfs: fixed an issue causing volumes mounted under other mounts to spontaneously unmount (docker/docker#24503)

### Docker for Mac 1.12.0 (2016-07-28 1.12.0-beta21)

**New**

* Docker for Mac is now available from 2 channels: **stable** and **beta**. New features and bug fixes will go out first in auto-updates to users in the beta channel. Updates to the stable channel are much less frequent and happen in sync with major and minor releases of the Docker engine. Only features that are well-tested and ready for production are added to the stable channel releases. For downloads of both and more information, see the [Getting Started](index.md#download-docker-for-mac).

**Upgrades**

* Docker 1.12.0 with experimental features
* Docker Machine 0.8.0
* Docker Compose 1.8.0

**Bug fixes and minor changes**

* Check for updates, auto-update and diagnose can be run by non-admin users
* osxfs: fixed an issue causing occasional incorrect short reads
* osxfs: fixed an issue causing occasional EIO errors
* osxfs: fixed an issue causing inotify creation events to fail
* osxfs: increased the fs.inotify.max_user_watches limit in Moby to 524288
* The UI shows documentation link for sharing volumes
* Clearer error message when running with outdated Virtualbox version
* Added link to sources for qemu-img

**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app

* There are a number of issues with the performance of directories bind-mounted with `osxfs`.  In particular, writes of small blocks, and traversals of large directories are currently slow.  Additionally, containers that perform large numbers of directory operations, such as repeated scans of large directory trees, may suffer from poor  performance. Applications that behave in this way include:

 - `rake`
 - `ember build`
 - Symfony
 - Magento

 As a work-around for this behavior, you can put vendor or third-party library directories in Docker volumes, perform temporary file system operations outside of `osxfs` mounts, and use third-party tools like Unison or `rsync` to synchronize between container directories and bind-mounted directories. We are actively working on `osxfs` performance using a number of different techniques and we look forward to sharing improvements with you soon. More information is available in [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md).

* Under some unhandled error conditions, inotify event delivery can fail and become permanently disabled. The workaround is to restart Docker.app

## Beta 20 Release Notes (2016-07-19 1.12.0-rc4-beta20)

**Bug fixes and minor changes**

* Fixed `docker.sock` permission issues
* Don't check for update when the settings panel opens
* Removed obsolete DNS workaround
* Use the secondary DNS server in more circumstances
* Limit the number of concurrent port forwards to avoid running out of resources
* Store the database as a "bare" git repo to avoid corruption problems

**Known issues**

*  `Docker.app` sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker for Mac (`Docker.app`).

## Beta 19 Release Notes (2016-07-14 1.12.0-rc4-beta19)

**New**

* Added privacy tab in settings
* Allow the definition of HTTP proxy overrides in the UI

**Upgrades**

* Docker 1.12.0 RC4
* Docker Compose 1.8.0 RC2
* Docker Machine 0.8.0 RC2
* Linux kernel 4.4.15

**Bug fixes and minor changes**

* Filesystem sharing permissions can only be configured in the UI (no more `/Mac` in moby)
* `com.docker.osx.xhyve.hyperkit`: increased max number of fds to 10240
* Improved Moby syslog facilities
* Improved file-sharing tab
* `com.docker.slirp`: included the DNS TCP fallback fix, required when UDP responses are truncated
* `docker build/events/logs/stats... ` won't leak when iterrupted with Ctrl-C

**Known issues**

* See [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

## Beta 18.1 Release Notes (2016-07-07 1.12.0-rc3-beta18.1)

>**Note**: Docker 1.12.0 RC3 release introduces a backward incompatible change from RC2. You can fix this by [recreating or updating your containers](troubleshoot.md#recreate-or-update-your-containers-after-beta-18-upgrade) as described in Troubleshooting.

**Hotfix**

* Fixed issue resulting in error "Hijack is incompatible with use of CloseNotifier", reverts previous fix for `Ctrl-C` during build.

**New**

* New host/container file sharing UI
* `/Mac` bind mount prefix is deprecated and will be removed soon

**Upgrades**

* Docker 1.12.0 RC3

**Bug fixes and minor changes**

* VPNKit: Improved scalability as number of network connections increases
* The docker API proxy was failing to deal with some 1.12 features (e.g. health check)

**Known issues**

* See [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

## Beta 18 Release Notes (2016-07-06 1.12.0-rc3-beta18)

**New**

* New host/container file sharing UI
* `/Mac` bind mount prefix is deprecated and will be removed soon

**Upgrades**

* Docker 1.12.0 RC3

**Bug fixes and minor changes**

* VPNKit: Improved scalability as number of network connections increases
* Interrupting a `docker build` with Ctrl-C will actually stop the build
* The docker API proxy was failing to deal with some 1.12 features (e.g. health check)

**Known issues**

* See [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

## Beta 17 Release Notes (2016-06-29 1.12.0-rc2-beta17)

**Upgrades**

* Linux kernel 4.4.14, aufs 20160627

**Bug fixes and minor changes**

* Documentation moved to https://docs.docker.com/docker-for-mac/
* Allow non-admin users to launch the app for the first time (using admin creds)
* Prompt non-admin users for admin password when needed in Preferences
* Fixed download links, documentation links
* Fixed "failure: No error" message in diagnostic panel
* Improved diagnostics for networking and logs for the service port openers

**Known issues**

* See [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

## Beta 16 Release Notes (2016-06-17 1.12.0-rc2-beta16)

**Upgrades**

* Docker 1.12.0 RC2
* docker-compose 1.8.0 RC1
* docker-machine 0.8.0 RC1
* notary 0.3
* Alpine 3.4

**Bug fixes and minor changes**

* VPNKit: Fixed a regressed error message when a port is in use
* Fixed UI crashing with `NSInternalInconsistencyException` / fixed leak
* HyperKit API: Improved error reporting
* osxfs: fix sporadic EBADF due to fd access/release races (#3683)


**Known issues**

* See [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

## Beta 15 Release Notes (2016-06-10 1.11.2-beta15)

**New**

* Registry mirror and insecure registries can now be configured from Preferences
* VM can now be restarted from Preferences
* `sysctl.conf` can be edited from Preferences

**Upgrades**

* Docker 1.11.2
* Linux 4.4.12, `aufs` 20160530

**Bug fixes and minor changes**

* Timekeeping in Moby VM improved
* Number of concurrent TCP/UDP connections increased in VPNKit
* Hyperkit: `vsock` stability improvements
* Fixed crash when user is admin

**Known issues**

* See [Known Issues](troubleshoot.md#known-issues) in [Troubleshooting](troubleshoot.md)

## Beta 14 Release Notes (2016-06-02 1.11.1-beta14)

**New**

* New settings menu item, **Diagnose & Feedback**, is available to run diagnostics and upload logs to Docker.

**Known issues**

* `Docker.app` sometimes uses 200% CPU after OS X wakes up from sleep mode with OSX 10.10. The issue is being investigated. The workaround is to restart `Docker.app`.

**Bug fixes and minor changes**

* `osxfs`: now support `statfs`
* **Preferences**: updated toolbar icons
* Fall back to secondary DNS server if primary fails.
* Added a link to the documentation from menu.

## Beta 13.1 Release Notes (2016-05-28 1.11.1-beta13.1)

**Hotfixes**

* `osxfs`:
  - Fixed sporadic EBADF errors and End_of_file crashes due to a race corrupting node table invariants
  - Fixed a crash after accessing a sibling of a file moved to another directory caused by a node table invariant violation
* Fixed issue where Proxy settings were applied on network change, causing docker daemon to restart too often
* Fixed issue where log file sizes doubled on docker daemon restart

## Beta 13 Release Notes (2016-05-25 1.11.1-beta13)

**New**

* `osxfs`: Enabled 10ms dcache for 3x speedup on a `go list ./...` test against docker/machine. Workloads heavy in file system path resolution (common among dynamic languages and build systems) will have those resolutions performed in amortized constant time rather than time linear in the depth of the path so speedups of 2-10x will be common.

* Support multiple users on the same machine, non-admin users can use the app as long as `vmnetd` has been installed. Currently, only one user can be logged in at the same time.

* Basic support for using system HTTP/HTTPS proxy in docker daemon

**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app.

**Bug fixes and minor changes**

* `osxfs`:
  - setting `atime` and `mtime` of nodes is now supported
  - Fixed major regression in Beta 12 with ENOENT, ENOTEMPY, and other spurious errors after a directory rename. This manifested as `npm install` failure and other directory traversal issues.
  - Fixed temporary file ENOENT errors
  - Fixed in-place editing file truncation error (e.g. `perl -i`)w
* improved time synchronisation after sleep

## Beta 12 Release (2016-05-17 1.11.1-beta12)

**Upgrades**

* FUSE 7.23 for [osxfs](osxfs.md)

**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app.

**Bug fixes and minor changes**

* UI improvements
* Fixed a problem in [osxfs](osxfs.md) where`mkdir` returned EBUSY but directory was created.

## Beta 11 Release (2016-05-10 1.11.1-beta11)

**New**

The `osxfs` file system now persists ownership changes in an extended attribute. (See the topic on [ownership](osxfs.md#ownership) in [Sharing the OS X file system with Docker containers](osxfs.md).)

**Upgrades**

* docker-compose 1.7.1 (see <a href="https://github.com/docker/compose/releases/tag/1.7.1" target="_blank"> changelog</a>)
* Linux kernel 4.4.9

**Bug fixes and minor changes**

* Desktop notifications after successful update
* No "update available" popup during install process
* Fixed repeated bind of privileged ports
* `osxfs`: Fixed the block count reported by stat
* Moby (Backend) fixes:
  - Fixed `vsock` half closed issue
  - Added NFS support
  - Hostname is now Moby, not Docker
  - Fixes to disk formatting scripts
  - Linux kernel upgrade to 4.4.9

## Beta 10 Release (2016-05-03 1.11.0-beta10)

**New**

* Token validation is now done over an actual SSL tunnel (HTTPS). (This should fix issues with antivirus applictions.)

**Upgrades**

* Docker 1.11.1

**Bug fixes and minor changes**

* UCP now starts again
* Include debugging symbols in HyperKit
* vsock stability improvements
* Addressed glitches in **Preferences** panel
* Fixed issues impacting the “whale menu”
* Fixed uninstall process
* HyperKit vcpu state machine improvements, may improve suspend/resume


## Beta 9 Release (2016-04-26 1.11.0-beta9)

**New**

* New Preferences window - memory and vCPUs now adjustable
* `localhost` is now used for port forwarding by default.`docker.local` will no longer work as of Beta 9.

**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app.

**Bug fixes and minor changes**

* Fix loopback device naming
* Improved docker socket download and osxfs sequential write by 20%
* `com.docker.osxfs`
  - improved sequential read throughput by up to 20%
  - improved `readdir` performance by up to 6x
  - log all fatal exceptions
* More reliable DNS forwarding over UDP and TCP
* UDP ports can be proxied over vsock
* Fixed EADDRINUSE (manifesting as errno 526) when ports are re-used
* Send ICMP when asked to not fragment and we can’t guarantee it
* Fixed parsing of UDP datagrams with IP socket options
* Drop abnormally large ethernet frames
* Improved HyperKit logging
* Record VM start and stop events

## Beta 8 Release (2016-04-20 1.11.0-beta8)

**New**

* Networking mode switched to VPN compatible by default, and as part of this change the overall experience has been improved:
 - `docker.local` now works in VPN compatibility mode
 - exposing ports on the Mac is available in both networking modes
 - port forwarding of privileged ports now works in both networking modes
 - traffic to external DNS servers is no longer dropped in VPN mode


* `osxfs` now uses `AF_VSOCK` for transport giving ~1.8x speedup for large sequential read/write workloads but increasing latency by ~1.3x. `osxfs` performance engineering work continues.


**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart `Docker.app`

**Bug fixes and minor changes**

* Apple System Log now used for most logs instead of direct filesystem logging
* `docker_proxy` fixes
* Merged HyperKit upstream patches
* Improved error reporting in `nat` network mode
* `osxfs` `transfused` client now logs over `AF_VSOCK`
* Fixed a `com.docker.osx.HyperKit.linux` supervisor deadlock if processes exit during a controlled shutdown
* Fixed VPN mode malformed DNS query bug preventing some resolutions


## Beta 7 Release (2016-04-12 1.11.0-beta7)

**New**

* Docs are updated per the Beta 7 release
* Use AF_VSOCK for docker socket transport

**Upgrades**

* docker 1.11.0-rc5
* docker-machine 0.7.0-rc3
* docker-compose 1.7.0rc2


**Known issues**

* Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app

* If VPN mode is enabled and then disabled and then re-enabled again, `docker ps` will block for 90s

**Bug fixes and minor changes**

* Logging improvements
* Improve process management

## Beta 6 Release (2016-04-05 1.11.0-beta6)

**New**

* Docs are updated per the Beta 6 release
* Added uninstall option in user interface

**Upgrades**

* docker 1.11.0-rc5
* docker-machine 0.7.0-rc3
* docker-compose 1.7.0rc2

**Known issues**

* `Docker.app` sometimes uses 200% CPU after OS X wakes up from sleep mode.
The issue is being investigated. The workaround is to restart
`Docker.app`.

* If VPN mode is enabled, then disabled and re-enabled again,
`docker ps` will block for 90 seconds.

**Bug fixes and minor changes**

* Fixed osxfs multiple same directory bind mounts stopping inotify
* Fixed osxfs `setattr` on mode 0 files (`sed` failures)
* Fixed osxfs blocking all operations during `readdir`
* Fixed osxfs mishandled errors which crashed the file system and VM
* Removed outdated `lofs`/`9p` support
* Added more debugging info to logs uploaded by `pinata diagnose`
* Improved diagnostics from within the virtual machine
* VirtualBox version check now also works without VBoxManage in path
* VPN mode now uses same IP range as NAT mode
* Tokens are now verified on port 443
* Removed outdated uninstall scripts
* Increased default ulimits
* Port forwarding with `-p` and `-P` should work in VPN mode
* Fixed a memory leak in `com.docker.db`
* Fixed a race condition on startup between Docker and networking which can
lead to `Docker.app` not starting on reboot

## Beta 5 Release (2016-03-29 1.10.3-beta5)

**New**

- Docs are updated per the Beta 5 release!

**Known issues**

- There is a race on startup between docker and networking which can lead to Docker.app not starting on reboot. The workaround is to restart the application manually.

- Docker.app sometimes uses 200% CPU after OS X wakes up from sleep mode. The issue is being investigated. The workaround is to restart Docker.app.

- In VPN mode, the `-p` option needs to be explicitly of the form `-p <host port>:<container port>`. `-p <port>` and `-P` will not work yet.

**Bug fixes and minor changes**

- Updated DMG background image
- Show correct VM memory in Preferences
- Feedback opens forum, not email
- Fixed RAM amount error message
- Fixed wording of CPU error dialog
- Removed status from Preferences
- Check for incompatible versions of Virtualbox

## Beta 4 Release (2016-03-22 1.10.3-beta4)

**New Features and Upgrades**

  <style type="text/css">
  .tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
  .tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
  .tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
  .tg .tg-vn4c{background-color:#D2E4FC}
  </style>
  <table class="tg">
    <tr>
      <th class="tg-031e">Component</th>
      <th class="tg-031e">Description</th>
    </tr>
    <tr>
      <td class="tg-vn4c">File System/Sharing</td>
      <td class="tg-vn4c">Support `inotify` events so that file system events on the
      Mac will trigger file system activations inside Linux containers</td>
    </tr>
    <tr>
      <td class="tg-031e">Docker Machine</td>
      <td class="tg-031e">Install Docker Machine as a part of Docker for Mac install in `/usr/local`</td>
    </tr>
    <tr>
      <td class="tg-vn4c">Getting Started and About</td>
      <td class="tg-vn4c">- Added  animated popover window to help first-time users get started<br>- Added a Beta icon to  About box</td>
    </tr>
  </table>

**Known issues**

  <style type="text/css">
  .tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
  .tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
  .tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
  .tg .tg-vn4c{background-color:#D2E4FC}
  </style>
  <table class="tg">
    <tr>
      <th class="tg-031e">Component</th>
      <th class="tg-031e">Description</th>
    </tr>
    <tr>
      <td class="tg-vn4c">Starting Docker</td>
      <td class="tg-vn4c">There is a race on startup between Docker and networking that can lead to `Docker.app` not starting on reboot.<br><br>The workaround is to restart the application manually.
      </td>
    </tr>
    <tr>
      <td class="tg-031e">OS X version support</td>
      <td class="tg-031e">`Docker.app` sometimes uses 200% CPU after OS X wakes up from sleep mode.
      The issue is being investigated. <br><br>The workaround is to restart
      `Docker.app`. <br>
      </td>
    </tr>
    <tr>
      <td class="tg-vn4c">VPN/Hostnet</td>
      <td class="tg-vn4c">In VPN mode, the `-p` option needs to be explicitly of the form
      `-p <host port>:<container port>`. `-p <port>` and `-P` will not
      work yet.
      </td>
    </tr>
  </table>

**Bug fixes and minor changes**

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
.tg .tg-vn4c{background-color:#D2E4FC}
</style>
<table class="tg">
  <tr>
    <th class="tg-031e">Component</th>
    <th class="tg-031e">Description</th>
  </tr>
  <tr>
    <td class="tg-vn4c">Hostnet/VPN mode</td>
    <td class="tg-vn4c">Fixed Moby DNS resolver failures by proxying the "Recursion Available" flag.<br>
    </td>
  </tr>
  <tr>
    <td class="tg-031e">IP addresses</td>
    <td class="tg-031e">`docker ps` shows IP address rather than `docker.local`.
    </td>
  </tr>
  <tr>
    <td class="tg-vn4c">OS X version support</td>
    <td class="tg-vn4c">- Re-enabled support for OS X Yosemite version 10.10 <br>- Ensured binaries are built for 10.10 rather than 10.11.
    </td>
  </tr>
  <tr>
    <td class="tg-031e">Application startup</td>
    <td class="tg-031e">- Fixed "Notification Center"-related crash on startup<br>- Fixed watchdog crash on startup</td>
  </tr>
</table>

## Beta 3 Release (2016-03-15 1.10.3-beta3)

**New Features and Upgrades**

  <style type="text/css">
  .tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
  .tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
  .tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
  .tg .tg-vn4c{background-color:#D2E4FC}
  </style>
  <table class="tg">
    <tr>
      <th class="tg-031e">Component</th>
      <th class="tg-031e">Description</th>
    </tr>
    <tr>
      <td class="tg-vn4c">File System</td>
      <td class="tg-vn4c">Improved file sharing write speed in OSXFS</td>
    </tr>
    <tr>
      <td class="tg-031e">User space networking</td>
      <td class="tg-031e">Renamed `bridged` mode to `nat` mode</td>
    </tr>
    <tr>
      <td class="tg-vn4c">Debugging</td>
      <td class="tg-vn4c">Docker runs in debug mode by default for new installs
      </td>
    </tr>
    <tr>
      <td class="tg-031e">Docker Engine</td>
      <td class="tg-031e">Upgraded to 1.10.3</td>
    </tr>
  </table>

**Bug fixes and minor changes**

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
.tg .tg-vn4c{background-color:#D2E4FC}
</style>
<table class="tg">
  <tr>
    <th class="tg-031e">Component</th>
    <th class="tg-031e">Description</th>
  </tr>
  <tr>
    <td class="tg-vn4c">GUI</td>
    <td class="tg-vn4c">Auto update automatically checks for new versions again<br>
    </td>
  </tr>
  <tr>
    <td class="tg-031e">File System</td>
    <td class="tg-031e">- Fixed OSXFS chmod on sockets <br>
    - FixED OSXFS EINVAL from `open` using O_NOFOLLOW</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Hypervisor</td>
    <td class="tg-vn4c">Hypervisor stability fixes, resynced with upstream repository <br></td>
  </tr>
  <tr>
    <td class="tg-031e">Hostnet/VPN mode</td>
    <td class="tg-031e">- Fixed get/set VPN mode in Preferences (GUI) <br>
    - Added more verbose logging on errors in `nat` mode<br>
    - Show correct forwarding details in `docker ps/inspect/port` in `nat` mode<br>
    </td>
  </tr>
  <tr>
    <td class="tg-vn4c">Tokens</td>
    <td class="tg-vn4c">New lines ignored in token entry field</td>
  </tr>
  <tr>
    <td class="tg-031e">Feedback</td>
    <td class="tg-031e">Feedback mail has app version in subject field</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Licensing</td>
    <td class="tg-vn4c">Clarified open source licenses</td>
  </tr>
  <tr>
    <td class="tg-031e">Crash reporting and error handling</td>
    <td class="tg-031e">- Fixed HockeyApp crash reporting <br>
    - Fatal GUI errors now correctly terminate the app again<br>
    - Fix proxy panics on EOF when decoding JSON<br>
    - Fix long delay/crash when switching from `hostnet` to `nat` mode
    </td>
  </tr>
  <tr>
    <td class="tg-vn4c">Logging</td>
    <td class="tg-vn4c">- Moby logs included in diagnose upload<br>
    - App version included in logs on startup
    </td>
  </tr>
</table>


## Beta 2 Release (2016-03-08 1.10.2-beta2)

**New Features and Upgrades**

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
.tg .tg-vn4c{background-color:#D2E4FC}
</style>
<table class="tg">
  <tr>
    <th class="tg-031e">Component</th>
    <th class="tg-031e">Description</th>
  </tr>
  <tr>
    <td class="tg-vn4c">GUI</td>
    <td class="tg-vn4c">Add VPN mode/`hostnet` to Preferences<br></td>
  </tr>
  <tr>
    <td class="tg-031e"></td>
    <td class="tg-031e">Add disable Time Machine backups of VM disk image to Preferences</td>
  </tr>
  <tr>
    <td class="tg-vn4c">CLI</td>
    <td class="tg-vn4c">Added `pinata` configuration tool for experimental Preferences</td>
  </tr>
  <tr>
    <td class="tg-031e">File System</td>
    <td class="tg-031e">Add guest-to-guest FIFO and socket file support</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Notary</td>
    <td class="tg-vn4c">Upgraded to version 0.2</td>
  </tr>
</table>

**Bug fixes**

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
.tg .tg-vn4c{background-color:#D2E4FC}
</style>
<table class="tg">
  <tr>
    <th class="tg-031e">Component</th>
    <th class="tg-031e">Description</th>
  </tr>
  <tr>
    <td class="tg-vn4c">File System</td>
    <td class="tg-vn4c">Fixed data corruption bug during cp (use of sendfile/splice)</td>
  </tr>
  <tr>
    <td class="tg-031e">GUI</td>
    <td class="tg-031e">Fixed About box to contain correct version string</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Hostnet/VPN mode</td>
    <td class="tg-vn4c">- Stability fixes and tests<br>- Fixed DNS issues when changing networks</td>
  </tr>
  <tr>
    <td class="tg-031e">Moby</td>
    <td class="tg-031e">Cleaned up Docker startup code</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Linking and dependencies</td>
    <td class="tg-vn4c">Fixed various problems</td>
  </tr>
  <tr>
    <td class="tg-031e">Logging</td>
    <td class="tg-031e">Various improvements</td>
  </tr>
</table>

<p>&nbsp;</p>

## Beta 1 Release (2016-03-01 1.10.2-b1)

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#999;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#444;background-color:#F7FDFA;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#999;color:#fff;background-color:#26ADE4;}
.tg .tg-vn4c{background-color:#D2E4FC}
</style>
<table class="tg">
  <tr>
    <th class="tg-031e">Component</th>
    <th class="tg-031e">Description </th>
  </tr>
  <tr>
    <td class="tg-vn4c">GUI</td>
    <td class="tg-vn4c">Added dialog to explain why we need admin rights</td>
  </tr>
  <tr>
    <td class="tg-031e"></td>
    <td class="tg-031e">Removed shutdown/quit window</td>
  </tr>
  <tr>
    <td class="tg-031e"></td>
    <td class="tg-031e">Improved machine migration</td>
  </tr>
  <tr>
    <td class="tg-031e"></td>
    <td class="tg-031e">Added "Help" option in menu to open documentation web pages</td>
  </tr>
  <tr>
    <td class="tg-031e"></td>
    <td class="tg-031e">Added license agreement</td>
  </tr>
  <tr>
    <td class="tg-031e"></td>
    <td class="tg-031e">Added MixPanel support</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Crash Reports</td>
    <td class="tg-vn4c">Add HockeyApp crash reporting</td>
  </tr>
  <tr>
    <td class="tg-yw4l">Task Manager</td>
    <td class="tg-yw4l">Improve signal handling</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Logging</td>
    <td class="tg-vn4c">Use ISO timestamps with microsecond precision</td>
  </tr>
  <tr>
    <td class="tg-yw4l"></td>
    <td class="tg-yw4l">Clean up logging format</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Packaging</td>
    <td class="tg-vn4c">Create /usr/local if it doesn't exist</td>
  </tr>
  <tr>
    <td class="tg-yw4l"></td>
    <td class="tg-yw4l">docker-uninstall improvements</td>
  </tr>
  <tr>
    <td class="tg-yw4l"></td>
    <td class="tg-yw4l">Remove docker-select as it's no longer used</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Hypervisor</td>
    <td class="tg-vn4c">Add PID file</td>
  </tr>
  <tr>
    <td class="tg-yw4l"></td>
    <td class="tg-yw4l">Networking reliability improvements</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Hostnet</td>
    <td class="tg-vn4c">Fixed port forwarding issue</td>
  </tr>
  <tr>
    <td class="tg-yw4l"></td>
    <td class="tg-yw4l">Stability fixes</td>
  </tr>
  <tr>
    <td class="tg-yw4l"></td>
    <td class="tg-yw4l">Fixed setting hostname</td>
  </tr>
  <tr>
    <td class="tg-vn4c">Symlinks</td>
    <td class="tg-vn4c">Fixed permissions on `usr/local` symbolic links</td>
  </tr>
</table>


<p style="margin-bottom:300px">&nbsp;</p>
