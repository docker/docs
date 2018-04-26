---
description: Change log / release notes per edge release
keywords: Docker for Windows, edge, release notes
title: Docker for Windows Edge Release notes
---

Here are the main improvements and issues per edge release, starting with the
current release. The documentation is always updated for each release.

For system requirements, see
[What to know before you install](install.md#what-to-know-before-you-install).

Release notes for _edge_ releases are listed below, [_stable_ release
notes](release-notes) are also available. (Following the CE release model,
'beta' releases are called 'edge' releases.)  You can learn about both kinds of
releases, and download stable and edge product installers at [Download Docker
for Windows](install.md#download-docker-for-windows).

## Edge Releases of 2018

### Docker Community Edition 18.05.0-ce-rc1-win63 2018-04-26

[Download](https://download.docker.com/win/edge/17439/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.05.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v18.05.0-ce-rc1)
  - [Notary 0.6.1](https://github.com/docker/notary/releases/tag/v0.6.1)

* Bug fixes and minor changes
  - Fix startup issue due to incompatibility with other programs (like Razer Synapse 3). Fixes [docker/for-win#1723](https://github.com/docker/for-win/issues/1723)
  - Fix Kubernetes hostPath translation for Persistent Volume Claim. Previously failing PVCs must be deleted and recreated. Fixes [docker/for-win#1758](https://github.com/docker/for-win/issues/1758)
  - Fix Kubernetes status when resetting to factory defaults.
  

### Docker Community Edition 18.04.0-ce-win62 2018-04-12

[Download](https://download.docker.com/win/edge/17151/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.04.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.04.0-ce)
  - [Docker compose 1.21.0](https://github.com/docker/compose/releases/tag/1.21.0)

### Docker Community Edition 18.04.0-ce-rc2-win61 2018-04-09

[Download](https://download.docker.com/win/edge/17070/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.04.0-ce-rc2](https://github.com/docker/docker-ce/releases/tag/v18.04.0-ce-rc2)
  - [Kubernetes 1.9.6](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.9.md#v196). If Kubernetes is enabled, the upgrade will be performed automatically when starting Docker for Windows.

* New 
  - Enable ceph & rbd modules in LinuxKit VM.

* Bug fixes and minor changes
  - Fix ApyProxy not starting properly when Docker for Windows is started with the `HOME` environment variable already defined (typically started from the command line). Fixes [docker/for-win#1880](https://github.com/docker/for-win/issues/1880)

### Docker Community Edition 18.03.0-ce-win58 2018-03-26

[Download](https://download.docker.com/win/edge/16761/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce)
  - [Docker compose 1.20.1](https://github.com/docker/compose/releases/tag/1.20.1)

* Bug fixes and minor changes
  - Adding Docker for Windows icon on desktop is optional in the installer. Fixes [docker/for-win#246](https://github.com/docker/for-win/issues/246)

### Docker Community Edition 18.03.0-ce-rc4-win57 2018-03-15

[Download](https://download.docker.com/win/edge/16511/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.0-ce-rc4](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce-rc4)
  - AUFS 20180312

* Bug fixes and minor changes
  - Fix support for AUFS. Fixes [docker/for-win#1831](https://github.com/docker/for-win/issues/1831)

### Docker Community Edition 18.03.0-ce-rc3-win56 2018-03-13

[Download](https://download.docker.com/win/edge/16433/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.0-ce-rc3](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce-rc3)
  - [Docker Machine 0.14.0](https://github.com/docker/machine/releases/tag/v0.14.0)
  - [Docker compose 1.20.0-rc2](https://github.com/docker/compose/releases/tag/1.20.0-rc2)
  - [Notary 0.6.0](https://github.com/docker/notary/releases/tag/v0.6.0)
  - Linux Kernel 4.9.87

* Bug fixes and minor changes
  - Fix port Windows Containers port forwarding on windows 10 build 16299 post KB4074588. Fixes [docker/for-win#1707](https://github.com/docker/for-win/issues/1707), [docker/for-win#1737](https://github.com/docker/for-win/issues/1737)
  - Fix for the HTTP/S transparent proxy when using "localhost" names (e.g. "host.docker.internal", "docker.for.win.host.internal", "docker.for.win.localhost").
  - If Kubernetes is enabled, switch CLI orchestrator option back to "swarm" when switching to Windows Containers.
  - Fix daemon not starting properly when setting TLS-related options.

### Docker Community Edition 18.03.0-ce-rc1-win54 2018-02-27

[Download](https://download.docker.com/win/edge/16164/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce-rc1)

* New
  - VM Swap size can be changed in settings.
  - Support NFS Volume sharing. Also works in Kubernetes. See [docker/for-win#1700](https://github.com/docker/for-win/issues/1700)
  - Allow users to activate Windows container during installation (avoid vm disk creation and vm boot when working only on win containers). See [docker/for-win#217](https://github.com/docker/for-win/issues/217).

* Bug fixes and minor changes
  - DNS name `host.docker.internal` shoud be used for host resolution from containers. Older aliases (still valid) are deprecated in favor of this one. (See https://tools.ietf.org/html/draft-west-let-localhost-be-localhost-06).
  - Fix Linuxkit start on Windows Insider. Fixes [docker/for-win#1458](https://github.com/docker/for-win/issues/1458), [docker/for-win#1514](https://github.com/docker/for-win/issues/1514), [docker/for-win#1640](https://github.com/docker/for-win/issues/1640)
  - Fix risk of privilege escalation. (https://www.tenable.com/sc-report-templates/microsoft-windows-unquoted-service-path-vulnerability)
  - All users present in the docker-users group are now able to use docker. Fixes [docker/for-win#1732](https://github.com/docker/for-win/issues/1732)
  - Kubernetes Load balanced services are no longer marked as `Pending`.
  - Fix hostPath mounts in Kubernetes.
  - Update Compose on Kubernetes to v0.3.0 rc4. Existing Kubernetes stacks will be removed during migration and need to be re-deployed on the cluster.


### Docker Community Edition 18.02.0-ce-win52 2018-02-08

[Download](https://download.docker.com/win/edge/15732/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.02.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.02.0-ce)
  - [Docker compose 1.19.0](https://github.com/docker/compose/releases/tag/1.19.0)

### Docker Community Edition 18.02.0-ce-rc2-win51 2018-02-02

* Upgrades
  - [Docker 18.02.0-ce-rc2](https://github.com/docker/docker-ce/releases/tag/v18.02.0-ce-rc2)
  - [Docker compose 1.19.0-rc2](https://github.com/docker/compose/releases/tag/1.19.0-rc2)
  - [Kubernetes 1.9.2](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.9.md#v192). If you have Kubernetes enabled, the upgrade will be performed automatically when starting Docker for Windows.

* New
  - VM disk size can be changed in settings. Fixes [docker/for-win#105](https://github.com/docker/for-win/issues/105)
  - New menu item to restart Docker.

* Bug fixes and minor changes
  - Migration of Docker Toolbox images is not proposed anymore in Docker For Windows installer (still possible to migrate Toolbox images manually).

### Docker Community Edition 18.02.0-ce-rc1-win50 2018-01-26

* Upgrades
  - [Docker 18.02.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v18.02.0-ce-rc1)

* New
  - Experimental Kubernetes Support. You can now run a single-node Kubernetes cluster from the "Kubernetes" Pane in Docker for Windows settings and use kubectl commands as well as docker commands. See https://docs.docker.com/docker-for-windows/kubernetes/
  - LCOW containers can now be run next to Windows containers (on Windows RS3 build 16299 and later). Use `--platform=linux` in Windows container mode to run Linux Containers On Windows. Note that LCOW is still experimental, it requires daemon `experimental` option.

* Bug fixes and minor changes
  - Better cleanup for Windows containers and images on reset/uninstall. Fixes [docker/for-win#1580](https://github.com/docker/for-win/issues/1580), [docker/for-win#1544](https://github.com/docker/for-win/issues/1544), [docker/for-win#191](https://github.com/docker/for-win/issues/191)
  - Do not recreate Desktop icon on upgrade (effective on next upgrade). Fixes [docker/for-win#246](https://github.com/docker/for-win/issues/246), [docker/for-win#925](https://github.com/docker/for-win/issues/925), [docker/for-win#1551](https://github.com/docker/for-win/issues/1551)
  - Fix proxy for docker.for.win.localhost & docker.for.win.host.internal. Fixes [docker/for-win#1130](https://github.com/docker/for-win/issues/1130)

### Docker Community Edition 18.01.0-ce-win48 2018-01-19

[Download](https://download.docker.com/win/edge/15285/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.01.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.01.0-ce)
  - Linux Kernel 4.9.75

* Bug fixes and minor changes
  - Fix linuxKit port-forwarder sometimes not being able to start. Fixes [docker/for-win#1506](https://github.com/docker/for-win/issues/1506)
  - Fix certificate management when connecting to a private registry. Fixes [docker/for-win#1512](https://github.com/docker/for-win/issues/1512)
  - Fix Mount compatibility when mounting drives with `-v //c/...`, now mounted in /host_mnt/c in the LinuxKit VM. Fixes [docker/for-win#1509](https://github.com/docker/for-win/issues/1509), [docker/for-win#1516](https://github.com/docker/for-win/issues/1516), [docker/for-win#1497](https://github.com/docker/for-win/issues/1497)

### Docker Community Edition 17.12.0-ce-win45 2018-01-05

[Download](https://download.docker.com/win/edge/15017/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 17.12.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce)

## Edge Releases of 2017
### Docker Community Edition 17.12.0-ce-rc4-win44 2017-12-21

* Upgrades
  - [Docker 17.12.0-ce-rc4](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce-rc4)
  - [Docker compose 1.18.0](https://github.com/docker/compose/releases/tag/1.18.0)

* Bug fixes and minor changes
  - Fix DNS "search domain" and "domain name" settings. See [docker/for-win#1437](https://github.com/docker/for-win/issues/1437).
  - Fix Vpnkit issue when username has spaces. See [docker/for-win#1429](https://github.com/docker/for-win/issues/1429).
  - Diagnostic improvements to get VM logs before VM shutdown.

### Docker Community Edition 17.12.0-ce-rc3-win43 2017-12-15

* Upgrades
  - [Docker 17.12.0-ce-rc3](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce-rc3)

* Bug fixes and minor changes
  - Fix installer check for not supported Windows `CoreCountrySpecific` Edition.

### Docker Community Edition 17.12.0-ce-rc2-win41 2017-12-13

* Upgrades
  - [Docker 17.12.0-ce-rc2](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce-rc2)
  - [Docker compose 1.18.0-rc2](https://github.com/docker/compose/releases/tag/1.18.0-rc2)


* Bug fixes and minor changes
  - Fix a class of startup failures where the database fails to start, see [docker/for-win#498](https://github.com/docker/for-win/issues/498)
  - Display various component versions in About box
  - Better removal of LCOW images & containers when uninstalling Docker
  - Links in Update changelog open the default browser instead of IE (fixes [docker/for-win#1311](https://github.com/docker/for-win/issues/1311))

### Docker Community Edition 17.11.0-ce-win40 2017-11-22

[Download](https://download.docker.com/win/edge/14328/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 17.11.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.11.0-ce)

### Docker Community Edition 17.11.0-ce-rc4-win39 2017-11-17

* Upgrades
  - [Docker 17.11.0-ce-rc4](https://github.com/docker/docker-ce/releases/tag/v17.11.0-ce-rc4)
  - [Docker compose 1.17.1](https://github.com/docker/compose/releases/tag/1.17.1)
  - Linux Kernel 4.9.60

* Bug fixes and minor changes
  - Increased timeout for VM boot startup to 2 minutes.

### Docker Community Edition 17.11.0-ce-rc3-win38 2017-11-09

* Upgrades
  - [Docker 17.11.0-ce-rc3](https://github.com/docker/docker-ce/releases/tag/v17.11.0-ce-rc3)

* New
  - Use Microsoft localhost port forwarder for Windows and LCOW Containers when it is available (insider build RS4)

* Bug fixes and minor changes
  - Fix docker build exits successfully but fails to build image [moby/#35413](https://github.com/moby/moby/issues/35413)
  - Fix bug during Windows fast-startup process. Fixes [for-win/#953](https://github.com/docker/for-win/issues/953)
  - Fix uninstaller issue (in some specific cases dockerd process was not killed properly)
  - Do not propose toolbox migration popup after clicking "Try LCOW" on first startup
  - Fix `docker.for.win.localhost` not working in proxy settings. Fixes [for-win/#1130](https://github.com/docker/for-win/issues/1130)

### Docker Community Edition 17.11.0-ce-rc2-win37 2017-11-02

* Upgrades
  - [Docker 17.11.0-ce-rc2](https://github.com/docker/docker-ce/releases/tag/v17.11.0-ce-rc2)
  - [Docker compose 1.17.0](https://github.com/docker/compose/releases/tag/1.17.0)
  - Linuxkit blueprint updated to [linuxkit/linuxkit#2633](https://github.com/linuxkit/linuxkit/pull/2633), fixes CVE-2017-15650

* New
  - Add localhost port forwarder for Windows and LCOW Containers (thanks @simonferquel)

* Bug fixes and minor changes
  - Fix centos:5 & centos:6 images not starting properly with LinuxKit VM (fixes [docker/for-win#1245](https://github.com/docker/for-win/issues/1245)).

### Docker Community Edition 17.10.0-ce-win36 2017-10-24

* Upgrades
  - [Docker 17.10.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.10.0-ce)
  - [Docker Machine 0.13.0](https://github.com/docker/machine/releases/tag/v0.13.0)
  - [Docker compose 1.17.0-rc1](https://github.com/docker/compose/releases/tag/1.17.0-rc1)

* New
  - VM entirely built with Linuxkit
  - Experimental support for Microsoft Linux Containers On Windows, on Windows 10 RS3.


### Docker Community Edition 17.09.0-ce-win34 2017-10-06

* Bug fixes
  - Fix Docker For Windows unable to start in some cases : removed use of libgmp sometimes causing the vpnkit process to die.


### Docker Community Edition 17.09.0-ce-win31 2017-09-29

* Upgrades
  - [Docker 17.09.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce)

* Bug fixes and minor changes
  - VPNKit: security fix to reduce the risk of DNS cache poisoning attack (reported by Hannes Mehnert https://hannes.nqsb.io/)


### Docker Community Edition 17.09.0-ce-rc3-win30 2017-09-22

* Upgrades
  - [Docker 17.09.0-ce-rc3](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce-rc3)

### Docker Community Edition 17.09.0-ce-rc2-win29 2017-09-19

* Upgrades
  - [Docker 17.09.0-ce-rc2](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce-rc2)
  - Linux Kernel 4.9.49
  - AUFS 20170911

* Bug fixes and minor changes
  - Kernel: Enable TASK_XACCT and TASK_IO_ACCOUNTING
  - Rotate logs in the VM more often (docker/for-win#244)
  - Vpnkit : do not block startup when ICMP permission is denied. (Fixes docker/for-win#1036, docker/for-win#1035, docker/for-win#1040)
  - Fix minor bug on update checks

### Docker Community Edition 17.09.0-ce-rc1-win28 2017-09-07

* Upgrades
  - [Docker 17.09.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce-rc1)
  - [Docker compose 1.16.1](https://github.com/docker/compose/releases/tag/1.16.1)
  - Linux Kernel 4.9.46

* New
  - Add `Skip This version` button in update window

* Bug fixes and minor changes
  - VPNKit: change protocol to support error messages reported back from the server
  - Reset to default stops all engines and removes settings including all daemon.json files
  - Better backend service checks (related to https://github.com/docker/for-win/issues/953)
  - Fix auto updates checkbox, no need to restart the application
  - Fix check for updates menu when auto updates was disable

### Docker Community Edition 17.07.0-win26 Release Notes (2017-09-01 17.07.0-win26)

* Upgrades
  - [Docker 17.07.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.07.0-ce)
  - [Docker compose 1.16.0](https://github.com/docker/compose/releases/tag/1.16.0)
  - [Docker Credential Helpers 0.6.0](https://github.com/docker/docker-credential-helpers/releases/tag/v0.6.0)

### Docker Community Edition 17.07.0-rc4-win25 Release Notes (2017-08-24 17.07.0-win25)

**Upgrades**

- [Docker 17.07.0-ce-rc4](https://github.com/docker/docker-ce/releases/tag/v17.07.0-ce-rc4)
- [Docker compose 1.16.0-rc1](https://github.com/docker/compose/releases/tag/1.16.0-rc1)

**Bug fixes and minor changes**

- VPNKit: Fixed a bug which causes a socket to leak if the corresponding
TCP connection is idle for more than 5 minutes (related to
[docker/for-mac#1374](https://github.com/docker/for-mac/issues/1374))

> **Note**: The link above goes to Docker for Mac issues because a
Mac user reported this problem, which applied to both Mac and Windows
and was fixed on both.

### Docker Community Edition 17.07.0-rc3-win23 Release Notes (2017-08-21 17.07.0-win23)

**Upgrades**

- [Docker 17.07.0-ce-rc3](https://github.com/docker/docker-ce/releases/tag/v17.07.0-ce-rc3)

**New**

- Store Linux daemon configuration in `~\.docker\daemon.json` instead of settings file
- Store Windows daemon configuration in `C:\ProgramData\Docker\config\daemon.json` instead of settings file
- VPNKit: Added support for ping!
- VPNKit: Added `slirp/port-max-idle-timeout` to allow the timeout to be adjusted or even disabled
- VPNKit: Bridge mode is default everywhere now

**Bug fixes and minor changes**

- VPNKit: Improved the logging around the Unix domain socket connections
- VPNKit: Automatically trim whitespace from `int` or `bool` database keys

### Docker Community Edition 17.07.0-ce-rc2-win22 Release Notes (2017-08-11 17.06.0-win22)

**Upgrades**

- [Docker 17.07.0-ce-rc2](https://github.com/docker/docker-ce/releases/tag/v17.07.0-ce-rc2)
- Linux Kernel 4.9.41

### Docker Community Edition 17.07.0-ce-rc1-win21 Release Notes (2017-07-31 17.07.0-win21)

**Upgrades**

- [Docker 17.07.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v17.07.0-ce-rc1)
- [Docker compose 1.15.0](https://github.com/docker/compose/releases/tag/1.15.0)
- [Docker Machine 0.12.2](https://github.com/docker/machine/releases/tag/v0.12.2)
- Linux Kernel 4.9.38

**New**

- Windows Docker daemon is now started as service for better lifecycle management

**Bug fixes and minor changes**

- Keep Docker info in the same place as before in the registry, used by Visual Studio 2017 (Fixes [docker/for-win#939](https://github.com/docker/for-win/issues/939))
- Fix `config.json` not being released properly (Fixes [docker/for-win#867](https://github.com/docker/for-win/issues/867))
- Do not anymore move credentials in credential store at startup

### Docker Community Edition 17.06.1-ce-rc1-win20 Release Notes (2017-07-18 17.06.1-win20)

**Upgrades**

- [Docker 17.06.1-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v17.06.1-ce-rc1)
- Linux Kernel 4.9.36
- AUFS 20170703

**Bug fixes and minor changes**

- Fix locked container id file (Fixes [docker/for-win#818](https://github.com/docker/for-win/issues/818))
- Avoid expanding variables in PATH env variable (Fixes [docker/for-win#859](https://github.com/docker/for-win/issues/859))

### Docker Community Edition 17.06.0-win17 Release Notes (2017-06-28 17.06.0-win17)

**Upgrades**

- [Docker 17.06.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce)

### Docker Community Edition 17.06.0-win16 Release Notes (2017-06-21 17.06.0-rc5-ce-win16)

**Upgrades**

- [Docker 17.06.0-ce-rc5](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce-rc5)
 - [Docker compose 1.14.0](https://github.com/docker/compose/releases/tag/1.14.0)

### Docker Community Edition 17.06.0-win15 Release Notes (2017-06-16 17.06.0-rc4-ce-win15)

**Upgrades**

- [Docker 17.06.0-rc4-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce-rc4)
- [Docker Credential Helpers 0.5.2](https://github.com/docker/docker-credential-helpers/releases/tag/v0.5.2)
- Linux Kernel 4.9.31


### Docker Community Edition 17.06.0-win14 Release Notes (2017-06-08 17.06.0-rc2-ce-win14)

**Upgrades**

  - [Docker 17.06.0-rc2-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce-rc2)
  - [Docker Machine 0.12.0](https://github.com/docker/machine/releases/tag/v0.12.0)
  - [Docker compose 1.14.0-rc2](https://github.com/docker/compose/releases/tag/1.14.0-rc2)

**New**

- Added reset to defaults button in error report window
- Unified login between Docker CLI and Docker Hub, Docker Cloud.

**Bug fixes and minor changes**

- Fixed group access check for users logged in with Active Directory (fixes [docker/for-win#785](https://github.com/docker/for-win/issues/785))
- Check environment variables and add some warnings in logs if they can cause docker to fail

### Docker Community Edition 17.06.0-win13 Release Notes (2017-06-01 17.06.0-rc1-ce-win13)

**Upgrades**

- [Docker 17.06.0-rc1-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce-rc1)
- [Docker Credential Helpers 0.5.1](https://github.com/docker/docker-credential-helpers/releases/tag/v0.5.1)
- Linux Kernel 4.9.30

**New**

- Sharing a drive can be done on demand, the first time a mount is requested
- Add an experimental DNS name for the host: `docker.for.win.localhost`
- Support for client (i.e. "login") certificates for authenticating
registry access (fixes [docker/for-win#569](https://github.com/docker/for-win/issues/569))

**Bug fixes and minor changes**

- Many processes that were running in admin mode are now running within the user identity
- Cloud federation command lines now opens in the user home directory
- Named pipes are now created with more constrained security descriptors to improve security
- Security fix : Users must be part of a specific group "docker-users" to run Docker for windows

### Docker Community Edition 17.0.5-win11 Release Notes (2017-05-12 17.05.0-ce-win11)

**Upgrades**

- Security fix for CVE-2017-7308

### Docker Community Edition 17.0.5-win9 Release Notes (2017-05-09 17.05.0-ce-win9)

**Upgrades**

- [Docker 17.05.0-ce](https://github.com/docker/docker/releases/tag/v17.05.0-ce)
- [Docker Compose 1.13.0](https://github.com/docker/compose/releases/tag/1.13.0)
- [Docker Machine 1.11.0](https://github.com/docker/machine/releases/tag/v0.11.0)

**Security**

- Disable TCP exposition of the Daemon (tcp://localhost:2375), now an opt-in feature.

**Bug fixes and minor changes**

- Reset to default / uninstall also reset docker cli settings and logout user from Docker Cloud and registries
- Detect a bitlocker policy preventing windows containers to work
- Fixed an issue on filesharing when explicitly disabled on vmswitch interface
- Fixed VM not starting when machine had very long name
- Fixed a bug where Windows daemon.json file was not written (fixes [docker/for-win#670](https://github.com/docker/for-win/issues/670))

### Docker Community Edition 17.0.5-win8 Release Notes (2017-04-13 17.05.0-ce-rc1-win8)

**Upgrades**

- [Docker 17.05.0-ce-rc1](https://github.com/docker/docker/releases/tag/v17.05.0-ce-rc1)

### Docker Community Edition 17.0.4-win7 Release Notes (2017-04-06 17.04.0-ce-win7)

**New**

- New installer experience
- Experimental Windows Server 2016 support

**Upgrades**

- [Docker 17.04.0-ce](https://github.com/docker/docker/releases/tag/v17.04.0-ce)
- [Docker Compose 1.12.0](https://github.com/docker/compose/releases/tag/1.12.0)
- Linux Kernel 4.9.19

**Bug fixes and minor changes**

- Added patches to the kernel to fix VMBus crash

### Docker Community Edition 17.04.0-ce-win6 Release Notes (2017-04-03 17.04.0-ce-rc2-win6)

**Upgrades**

- [Docker 17.04.0-ce-rc2](https://github.com/docker/docker/releases/tag/v17.04.0-ce-rc2)
- [Docker Compose 1.12.0-rc2](https://github.com/docker/compose/releases/tag/1.12.0-rc2)
- Linux Kernel 4.9.18

**Bug fixes and minor changes**

- Named pipe client connection should not trigger dead locks on `docker run` with data in stdin anymore
- Buffered data should be treated correctly when docker client requests are upgraded to raw streams

### Docker Community Edition 17.03.1 Release Notes (2017-03-28 17.03.1-ce-rc1-win3)

**Upgrades**

- [Docker 17.03.1-ce-rc1](https://github.com/docker/docker/releases/tag/v17.03.1-ce-rc1)
- [Docker Credential Helpers 0.5.0](https://github.com/docker/docker-credential-helpers/releases/tag/v0.5.0)
- Linux Kernel 4.9.14

**Bug fixes and minor changes**

- VPNKit: capture up to 64KiB of NTP traffic for diagnostics, better handling of DNS

### Docker Community Edition 17.03.0 Release Notes (2017-03-06 17.03.0-ce-win1)

**New**

- Renamed to Docker Community Edition
- Integration with Docker Cloud: control remote Swarms from the local CLI and view your repositories. This feature is going to be rolled out to all users progressively

**Upgrades**

- [Docker 17.03.0-ce](https://github.com/docker/docker/releases/tag/v17.03.0-ce)
- [Docker Compose 1.11.2](https://github.com/docker/compose/releases/tag/1.11.2)
- [Docker Machine 0.10.0](https://github.com/docker/machine/releases/tag/v0.10.0)
- Linux Kernel 4.9.12

**Bug fixes and minor changes**

- VPNKit: fix unmarshalling of DNS packets containing pointers to pointers to labels
- Match Hyper-V Integration Services by ID, not name
- Don't consume 100% CPU when the service is stopped
- Log the diagnostic id when uploading
- Improved Firewall handling: stop listing the rules since it can take a lot of time
- Don't rollback to the previous engine when the desired engine fails to start

### Docker Community Edition 17.03.0 Release Notes (2017-02-22 17.03.0-ce-rc1-win1)

**New**

- Introduce Docker Community Edition
- Integration with Docker Cloud: control remote Swarms from the local CLI and view your repositories. This feature is being rolled out to all users progressively.

**Upgrades**

- Docker 17.03.0-ce-rc1
- Linux Kernel 4.9.11

**Bug fixes and minor changes**

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

## Edge Releases of 2016
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

- Time drifts between Windows and Linux containers should disappear
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

- Time drifts between Windows and Linux containers should disappear
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
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.
>
Windows Container support relies on the Windows 10 container feature, which is
**experimental** at this point.  Windows 10 Pro (1607, build number 14393)
requires update `KB3192366` (soon to be released via Windows Update) to fully
work. Some insider builds may not work.

**New**

- Restore the VM's configuration when it was changed by the user
- Overlay2 is now the default storage driver. After a factory reset overlay2 is automatically used
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
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

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
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

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
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* Basic support for Windows containers. On Windows 10 build >= 14372, a switch in the `systray` icon changes which daemon (Linux or Windows) the Docker CLI talks to

* To support trusted registry transparently, all trusted CAs (root or intermediate) on the Windows host are automatically copied to Moby

* `Reset Credentials` also unshares the shared drives

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
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

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

* Docker automatically disables lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Remove stale network adapters](troubleshoot.md#4-remove-stale-network-adapters) under [Networking issues](troubleshoot.md#networking-issues) in Troubleshooting.

### Beta 24 Release (2016-08-23 1.12.1-beta24)

>**Important Note**:
>
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**Upgrades**

* Docker 1.12.1
* Docker Machine 0.8.1
* Linux kernel 4.4.19
* aufs 20160822

**Bug fixes and minor changes**

* `slirp`: reduce the number of sockets used by UDP NAT, reduce the probability that NAT rules time out earlier than expected

**Known issues**

* Only UTF-8 passwords are supported for host filesystem sharing.

* Docker automatically disables lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Remove stale network adapters](troubleshoot.md#4-remove-stale-network-adapters) under [Networking issues](troubleshoot.md#networking-issues) in Troubleshooting.

### Beta 23 Release (2016-08-16 1.12.1-rc1-beta23)

>**Important Note**:
>
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

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
* Docker stops asking to import from toolbox after an upgrade
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
* Docker automatically disables lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Troubleshooting](troubleshoot.md#networking-issues).

### Beta 22 Release (2016-08-11 1.12.0-beta22)

Unreleased. See Beta 23 for changes.

**Known issues**

* Docker automatically disables lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Troubleshooting](troubleshoot.md#networking-issues).

### Beta 21 Release (2016-07-28 1.12.0-beta21)

**New**

* Docker for Windows is now available from 2 channels: **stable** and **beta**. New features and bug fixes go out first in auto-updates to users in the beta channel. Updates to the stable channel are much less frequent and happen in sync with major and minor releases of the Docker engine. Only features that are well-tested and ready for production are added to the stable channel releases. For downloads of both and more information, see the [Getting Started](index.md#download-docker-for-windows).

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

* Added an option to opt-out from sending usage statistics (only on the Edge channel for now)
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

* Interrupting a `docker build` with Ctrl-C actually stops the build
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

* Logs for the windows service are not aggregated with logs from the GUI. This is expected to be fixed in future versions.


### Beta 10 Release (2016-05-03 1.11.0-beta10)

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
* The firewall alert dialog does not come up as often as previously
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
* Better error handling: Moby restarts itself if start takes too long.
* Kill proxy and exit docker before a new version is installed
* The application cannot start twice now
* The proxy stops automatically when the GUI is not running
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

- Settings are now serialized in JSON. This install loses the current settings.

- Docker needs to open ports on the firewall. Sometimes, the user sees a firewall alert dialog. The user should allow the ports to be opened.

- The application was upgraded to 64 bits. The installation path changed to `C:\Program Files\Docker\Docker`. Users might need to close any Powershell/Cmd windows that were already open before the update to get the new `PATH`. In some cases, users may need to log off and on again.

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

- The application was upgraded to 64 bits. The installation path changed to `C:\Program Files\Docker\Docker`. If users have Powershell/Cmd windows already open before the update, they might need to close them to catch the new PATH. In some cases, users need to log off and on again.

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
- The application refuses to start if Hyper-v is not enabled
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

### Alpha 4 Release (2016-03-10 1.10.4.0)

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
