---
description: Change log / release notes for Docker Desktop Mac
keywords: Docker Desktop for Mac, release notes
redirect_from:
- /mackit/release-notes/
- /docker-for-mac/edge-release-notes/
title: Docker for Mac release notes
toc_min: 1
toc_max: 2
---

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

> **Important**
>
> Starting with Docker Desktop 3.0.0, Stable and Edge releases are combined into a single release stream for all users. Updates to Docker Desktop will now be available automatically as delta updates from the previous version. This means, when there is a newer version of Docker Desktop, it will be automatically downloaded to your machine. All you need to do is to click **Update and restart** from the Docker menu to install the latest update.
{: .important }

## Docker Desktop 3.2.2
2021-03-15

> [Download](https://desktop.docker.com/mac/stable/amd64/Docker.dmg)

### Bug fixes and minor changes

- Fixed an issue that stopped containers binding to port 53. Fixes [docker/for-mac#5416](https://github.com/docker/for-mac/issues/5416).
- Fixed an issue that 32-bit Intel binaries were emulated on Intel CPUs. Fixes [docker/for-win#10594](https://github.com/docker/for-win/issues/10594).
- Fixed an issue related to high CPU consumption and frozen UI when the network connection is lost. Fixes [for-win/#10563](https://github.com/docker/for-win/issues/10563).
- Fixed an issue opening a terminal in iTerm2 when it has no other windows open. Fixes [docker/roadmap#98](https://github.com/docker/roadmap/issues/98#issuecomment-791927788).

## Docker Desktop 3.2.1
2021-03-05

> [Download](https://desktop.docker.com/mac/stable/amd64/61626/Docker.dmg)

### Upgrades

- [Docker Engine 20.10.5](https://docs.docker.com/engine/release-notes/#20105)

### Bug fixes and minor changes

- Fixed an issue that sometimes caused Docker Desktop to fail to start after updating to version 3.2.0. Fixes [docker/for-mac#5406](https://github.com/docker/for-mac/issues/5406). If you are still experiencing this issue when trying to update from 3.2.0 to 3.2.1, we recommend that you uninstall 3.2.0 and manually install Docker Desktop 3.2.1.

## Docker Desktop 3.2.0
2021-03-01

> [Download](https://desktop.docker.com/mac/stable/amd64/61504/Docker.dmg)

### New

- The Docker Dashboard opens automatically when you start Docker Desktop.
- The Docker Dashboard displays a tip once a week.
- Docker Desktop uses iTerm2 to launch the terminal on the container if it is installed. Otherwise, it launches the default Terminal.App. [docker/roadmap#98](https://github.com/docker/roadmap/issues/98)
- Add experimental support to use the new Apple Virtualization framework (requires macOS Big Sur 11.1 or later)
- BuildKit is now the default builder for all users, not just for new installations. To turn this setting off, go to **Preferences** > **Docker Engine** and add the following block to the Docker daemon configuration file:
```json
"features": {
    "buildkit": false
}
```

### Upgrades

- [Docker Engine 20.10.3](https://docs.docker.com/engine/release-notes/#20103)
- [Docker Compose 1.28.5](https://github.com/docker/compose/releases/tag/1.28.5)
- [Compose CLI v1.0.9](https://github.com/docker/compose-cli/tree/v1.0.9)
- [Docker Hub Tool v0.3.0](https://github.com/docker/hub-tool/releases/tag/v0.3.0)
- [QEMU 5.0.1](https://wiki.qemu.org/ChangeLog/5.0)
- [Amazon ECR Credential Helper v0.5.0](https://github.com/awslabs/amazon-ecr-credential-helper/releases/tag/v0.5.0)
- [Alpine 3.13](https://alpinelinux.org/posts/Alpine-3.13.0-released.html)
- [Kubernetes 1.19.7](https://github.com/kubernetes/kubernetes/releases/tag/v1.19.7)
- [Go 1.16](https://golang.org/doc/go1.16)

### Bug fixes and minor changes

- Fixed an issue on the container detail screen where the buttons would disappear when scrolling the logs. Fixes [docker/for-mac#5290](https://github.com/docker/for-mac/issues/5290)
- Fixed an issue when port forwarding multiple ports with an IPv6 container network. Fixes [docker/for-mac#5247](https://github.com/docker/for-mac/issues/5247)
- Fixed a regression where `docker load` could not use an xz archive anymore. Fixes [docker/for-mac#5271](https://github.com/docker/for-mac/issues/5271)
- Fixed a navigation issue in the **Containers / Apps** view. Fixes [docker/for-win#10160](https://github.com/docker/for-win/issues/10160#issuecomment-764660660)
- Fixed container instance view with long container/image name. Fixes [docker/for-mac#5290](https://github.com/docker/for-mac/issues/5290)
- Fixed an issue when binding ports on specific IPs. Note: It may now take a bit of time before the `docker inspect` command shows the open ports. Fixes [docker/for-mac#4541](https://github.com/docker/for-mac/issues/4541)
- Fixed an issue where an image deleted from the Docker dashboard was still displayed on the **Images** view.

### Known issue

Docker Desktop can sometimes fail to start after updating to version 3.2.0. If you are experiencing this issue, we recommend that you uninstall 3.2.0 and manually install [Docker Desktop 3.2.1](#docker-desktop-321). See [docker/for-mac#5406](https://github.com/docker/for-mac/issues/5406).

## Docker Desktop 3.1.0
2021-01-14

> [Download](https://desktop.docker.com/mac/stable/51484/Docker.dmg)

### New

- Docker daemon now runs within a Debian Buster based container (instead of Alpine).

### Upgrades

- [Compose CLI v1.0.7](https://github.com/docker/compose-cli/tree/v1.0.7)

### Bug fixes and minor changes

  - Fixed UI reliability issues when users create or delete a lot of objects in batches.
  - Fixed an issue with DNS address resolution in Alpine containers. Fixes [docker/for-mac#5020](https://github.com/docker/for-mac/issues/5020).
  - Redesigned the **Support** UI to improve usability.

## Docker Desktop 3.0.4
2021-01-06

> [Download](https://desktop.docker.com/mac/stable/51218/Docker.dmg)

### Upgrades

- [Docker Engine 20.10.2](https://docs.docker.com/engine/release-notes/#20102)

### Bug fixes and minor changes

- Avoid timeouts during `docker-compose up` by making cache invalidation faster. Fixes [docker/for-mac#4957](https://github.com/docker/for-mac/issues/4957).
- Avoid generating a spurious filesystem DELETE event while invalidating caches. Fixes [docker/for-mac#5124](https://github.com/docker/for-mac/issues/5124).

### Known issues

- Some DNS addresses fail to resolve within containers based on Alpine Linux 3.13. See [docker/for-mac#5020](https://github.com/docker/for-mac/issues/5020).

## Docker Desktop 3.0.3
2020-12-21

> [Download](https://desktop.docker.com/mac/stable/51017/Docker.dmg)

### Bug fixes and minor changes

- Fixed an issue that caused overlapping volume mounts to fail. Fixes [docker/for-mac#5157](https://github.com/docker/for-mac/issues/5157). However, the fixes for [docker/for-mac#4957](https://github.com/docker/for-mac/issues/4957) and [docker/for-mac#5124](https://github.com/docker/for-mac/issues/5124) have been reverted as a result of this change, so those issues are now present again.

### Known issues

- Some DNS addresses fail to resolve within containers based on Alpine Linux 3.13. See [docker/for-mac#5020](https://github.com/docker/for-mac/issues/5020).
- There can be timeouts during docker-compose up if there are several services being started. See [docker/for-mac#4957](https://github.com/docker/for-mac/issues/4957) and [docker/for-mac#5124](https://github.com/docker/for-mac/issues/5124).

## Docker Desktop 3.0.2
2020-12-18

> [Download](https://desktop.docker.com/mac/stable/50996/Docker.dmg)

### Bug fixes and minor changes

- Avoid timeouts during `docker-compose up` by making cache invalidation faster. Fixes [docker/for-mac#4957](https://github.com/docker/for-mac/issues/4957).
- Avoid generating a spurious filesystem DELETE event while invalidating caches. Fixes [docker/for-mac#5124](https://github.com/docker/for-mac/issues/5124).
- It is now possible to share directories in `~/Library` (except Docker Desktop data directories) with a container. Fixes [docker/for-mac#5115](https://github.com/docker/for-mac/issues/5115).
- You will now see a performance warning pop-up message if you create a container that shares the `Home` or  a user `Library` directory.

### Known issues

- Some DNS addresses fail to resolve within containers based on Alpine Linux 3.13. See [docker/for-mac#5020](https://github.com/docker/for-mac/issues/5020).

## Docker Desktop 3.0.1
2020-12-11

> [Download](https://desktop.docker.com/mac/stable/50773/Docker.dmg)

### Bug fixes and minor changes

- Fixed an issue that caused certain directories not to be mountable into containers. Fixes [docker/for-mac#5115](https://github.com/docker/for-mac/issues/5115). See Known issues below.

### Known issues

- It is currently not possible to bind mount files within `~/Libary` into a container. See [docker/for-mac#5115](https://github.com/docker/for-mac/issues/5115).
- Building an image with BuildKit from a git URL fails when using the form `github.com/org/repo`. To work around this issue, use the form `git://github.com/org/repo`.
- Some DNS addresses fail to resolve within containers based on Alpine Linux 3.13. See [docker/for-mac#5020](https://github.com/docker/for-mac/issues/5020).

## Docker Desktop 3.0.0
2020-12-10

> [Download](https://desktop.docker.com/mac/stable/50684/Docker.dmg)

### New

- Use of three-digit version number for Docker Desktop releases.
- Starting with Docker Desktop 3.0.0, updates are now much smaller as they will be applied using delta patches. For more information, see [Automatic updates](install.md#automatic-updates).
- First version of `docker compose` (as an alternative to the existing `docker-compose`). Supports some basic commands but not the complete functionality of `docker-compose` yet.

  - Supports the following subcommands: `up`, `down`, `logs`, `build`, `pull`, `push`, `ls`, `ps`
  - Supports basic volumes, bind mounts, networks, and environment variables

    Let us know your feedback by creating an issue in the [compose-cli](https://github.com/docker/compose-cli/issues){: target="blank" rel="noopener" class=“”} GitHub repository.
- [Docker Hub Tool v0.2.0](https://github.com/docker/roadmap/issues/117){: target="blank" rel="noopener" class=“”}

### Upgrades

- [Docker Engine 20.10.0](https://docs.docker.com/engine/release-notes/#20100)
- [Go 1.15.6](https://github.com/golang/go/issues?q=milestone%3AGo1.15.6+label%3ACherryPickApproved+)
- [Compose CLI v1.0.4](https://github.com/docker/compose-cli/releases/tag/v1.0.4)
- [Snyk v1.432.0](https://github.com/snyk/snyk/releases/tag/v1.432.0)

### Bug fixes and minor changes

- Downgraded the kernel to [4.19.121](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.121-2a1dbedf3f998dac347c499808d7c7e029fbc4d3-amd64/images/sha256-4e7d94522be4f25f1fbb626d5a0142cbb6e785f37e437f6fd4285e64a199883a?context=repo) to reduce the CPU usage of hyperkit. Fixes [docker/for-mac#5044](https://github.com/docker/for-mac/issues/5044)
- Avoid caching bad file sizes and modes when using `osxfs`. Fixes [docker/for-mac#5045](https://github.com/docker/for-mac/issues/5045).
- Fixed a possible file sharing error where a file may appear to have the wrong size in a container when it is modified on the host. This is a partial fix for [docker/for-mac#4999](https://github.com/docker/for-mac/issues/4999).
- Removed unnecessary log messages which slow down filesystem event injection.
- Re-enabled the experimental SOCKS proxy. Fixes [docker/for-mac#5048](https://github.com/docker/for-mac/issues/5048).
- Fixed an unexpected EOF error when trying to start a non-existing container with `-v /var/run/docker.sock:`. See [docker/for-mac#5025](https://github.com/docker/for-mac/issues/5025).
- Display an error message instead of crashing when the application needs write access on specific directories. See [docker/for-mac#5068](https://github.com/docker/for-mac/issues/5068)

### Known issues

- Building an image with BuildKit from a git URL fails when using the form `github.com/org/repo`. To work around this issue, use the form `git://github.com/org/repo`.
- Some DNS addresses fail to resolve within containers based on Alpine Linux 3.13.

## Docker Desktop Community 2.5.0.1
2020-11-10

> [Download](https://desktop.docker.com/mac/stable/49550/Docker.dmg)

### Upgrades

- [Compose CLI v1.0.2](https://github.com/docker/compose-cli/releases/tag/v1.0.2)
- [Snyk v1.424.4](https://github.com/snyk/snyk/releases/tag/v1.424.4)

### Bug fixes and minor changes
- Fixed an issue that caused Docker Desktop to crash on MacOS 11.0 (Big Sur) when VirtualBox was also installed. See [docker/for-mac#4997](https://github.com/docker/for-mac/issues/4997).

## Docker Desktop Community 2.5.0.0
2020-11-02

> [Download](https://desktop.docker.com/mac/stable/49427/Docker.dmg)

Docker Desktop 2.5.0.0 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### New

- Users subscribed to a Pro or a Team plan can now see the vulnerability scan report on the Remote repositories tab in Docker Desktop.
- Docker Desktop introduces a support option for users who have subscribed to a Pro or a Team Plan.

### Security

- Fixed local privilege escalation vulnerability caused by inadequate certificate checking. See [CVE-2021-3162](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-3162){:target="_blank" rel="noopener" class="_"}.

### Upgrades

- [Linux kernel 5.4.39](https://hub.docker.com/layers/linuxkit/kernel/5.4.39-f39f83d0d475b274938c86eaa796022bfc7063d2/images/sha256-8614670219aca0bb276d4749e479591b60cd348abc770ac9ecd09ee4c1575405?context=explore)
- [Docker Compose CLI 1.0.1](https://github.com/docker/compose-cli/releases/tag/v1.0.1)
- [Snyk v1.421.1](https://github.com/snyk/snyk/releases/tag/v1.421.1)
- [Go 1.15.2](https://github.com/golang/go/releases/tag/go1.15.2)
- [Kubernetes 1.19.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.19.3)

### Bug fixes and minor changes

- Renamed 'Run Diagnostics' to 'Get support'.
- Removed BlueStacks warning message. Fixes [docker/for-mac#4863](https://github.com/docker/for-mac/issues/4863).
- Made container start faster in cases where shared volumes have lots of files. Fixes [docker/for-mac#4957](https://github.com/docker/for-mac/issues/4957).
- File sharing: fixed changing ownership of read-only files. Fixes [docker/for-mac#4989](https://github.com/docker/for-mac/issues/4989), [docker/for-mac#4964](https://github.com/docker/for-mac/issues/4964).
- File sharing: generated `ATTRIB` inotify events as well as `MODIFY`. Fixes [docker/for-mac#4962](https://github.com/docker/for-mac/issues/4962).
- File sharing: returned `EOPNOTSUPP` from `fallocate` for unsupported modes. Fixes `minio`. Fixes [docker/for-mac#4964](https://github.com/docker/for-mac/issues/4964).
- File sharing: fixed a possible premature file handle close.
- When sharing Linux directories (`/var`, `/bin`, etc) with containers, Docker Desktop avoids watching paths in the host file system.
- When sharing a file into a container (e.g. `docker run -v ~/.gitconfig`) Docker Desktop does not watch the parent directory. Fixes [docker/for-mac#4981](https://github.com/docker/for-mac/issues/4981), [docker/for-mac#4975](https://github.com/docker/for-mac/issues/4975).
- Fixed an issue related to NFS mounting. Fixes [docker/for-mac#4958](https://github.com/docker/for-mac/issues/4958).
- Allow symlinks to point outside of shared volumes. Fixes [docker/for-mac#4862](https://github.com/docker/for-mac/issues/4862).
- Diagnostics: avoid hanging when Kubernetes is in a broken state.
- Docker Desktop now supports `S_ISUID`, `S_ISGID` and `S_ISVTX` in calls to `chmod(2)` on shared filesystems. Fixes [docker/for-mac#4943](https://github.com/docker/for-mac/issues/4943).

## Docker Desktop Community 2.4.0.0
2020-09-30

> [Download](https://desktop.docker.com/mac/stable/48506/Docker.dmg)

Docker Desktop 2.4.0.0 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### New

- [Docker Compose CLI - 0.1.18](https://github.com/docker/compose-cli), enabling use of volumes with Compose and the Cloud through ECS and ACI.
- Docker introduces the new Images view in the Docker Dashboard. The images view allows users to view the Hub images, pull them and manage their local images on disk including cleaning up unwanted and unused images. To access the new Images view, from the Docker menu, select **Dashboard** > **Images**.
- Docker Desktop now enables BuildKit by default after a reset to factory defaults. To revert to the old `docker build` experience, go to **Preferences** > **Docker Engine** and then disable the BuildKit feature.
- [Amazon ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper/releases/tag/v0.4.0)
- Docker Desktop now uses much less CPU when there are lots of file events on the host and when running Kubernetes, see [docker/roadmap#12](https://github.com/docker/roadmap/issues/12).
- Docker Desktop now uses gRPC-FUSE for file sharing by default. This uses much less CPU than osxfs, especially when there are lots of file events on the host. To switch back to `osxfs`, go to **Preferences** > **General** and disable gRPC-FUSE.

### Upgrades

- [Docker 19.03.13](https://github.com/docker/docker-ce/releases/tag/v19.03.13)
- [Docker Compose 1.27.4](https://github.com/docker/compose/releases/tag/1.27.4)
- [Go 1.14.7](https://github.com/golang/go/releases/tag/go1.14.7)
- [Alpine 3.12](https://alpinelinux.org/posts/Alpine-3.12.0-released.html)
- [Kubernetes 1.18.8](https://github.com/kubernetes/kubernetes/releases/tag/v1.18.8)
- [Qemu 4.2.0](https://git.qemu.org/?p=qemu.git;a=tag;h=1e4aa2dad329852aa6c3f59cefd65c2c2ef2062c)

### Bug fixes and minor changes

- Docker Desktop on macOS 10.13 is now deprecated.
- Removed the legacy Kubernetes context `docker-for-desktop`. The context `docker-desktop` should be used instead. Fixes [docker/for-win#5089](https://github.com/docker/for-win/issues/5089) and [docker/for-mac#4089](https://github.com/docker/for-mac/issues/5089).
- Adding the application to the dock and clicking on it will launch the container view if Docker is already running.
- Added support for emulating Risc-V via Qemu 4.2.0.
- Removed file descriptor limit (`setrlimit`) of `10240`. We now rely on the kernel to impose limits via `kern.maxfiles` and `kern.maxfilesperproc`.
- Fixed a Mac CPU usage bug by removing the serial console from `hyperkit`, see [docker/roadmap#12]( https://github.com/docker/roadmap/issues/12#issuecomment-663163280). To open a shell in the VM use either `nc -U ~/Library/Containers/com.docker.docker/Data/debug-shell.sock`.
- Copy container logs without ansi colors to clipboard. Fixes [docker/for-mac#4786](https://github.com/docker/for-mac/issues/4786).
- Fixed automatic start on log in. See [docker/for-mac#4877] and [docker/for-mac#4890].
- Fixed bug where the application won't start if the username is too long.
- Fixed a bug where adding directories like `/usr` to the filesharing list prevents Desktop from starting. Fixes [docker/for-mac#4488](https://github.com/docker/for-mac/issues/4488)
- Fixed application startup if `hosts` is specified inside the Docker `daemon.json`. See [docker/for-win#6895](https://github.com/docker/for-win/issues/6895#issuecomment-637429117)
- Docker Desktop always flushes filesystem caches synchronously on container start. See [docker/for-mac#4943](https://github.com/docker/for-mac/issues/4943).
- Compose-on-Kubernetes is no longer included in the Docker Desktop installer. You can download it separately from the compose-on-kubernetes [release page](https://github.com/docker/compose-on-kubernetes/releases).

### Known issues

-  There is a known issue when using `docker-compose` with named volumes and gRPC FUSE: second and subsequent calls to `docker-compose up` will fail due to the volume path having the prefix `/host_mnt`. To work around this issue, switch back to `osxfs` in Settings. See [docker/for-mac#4859](https://github.com/docker/for-mac/issues/4859).
- There is a known issue when enabling Kubernetes where the settings UI fails to update the Kubernetes state. To work around this issue, close and re-open the window.
- There is a rare known issue when switching users where the images view continues to show the repositories of the previous user. To work around this issue, close and re-open the window.

## Docker Desktop Community 2.3.0.5
2020-09-15

> [Download](https://desktop.docker.com/mac/stable/48029/Docker.dmg)

### New

- The new Cloud integration in Docker CLI makes it easy to run containers in the cloud using either Amazon ECS or Microsoft ACI. For more information, see [Deploying Docker containers on ECS](https://docs.docker.com/engine/context/ecs-integration/) and [Deploying Docker containers on Azure](https://docs.docker.com/engine/context/aci-integration/).

### Upgrades

- [Docker Compose 1.27.2](https://github.com/docker/compose/releases/tag/1.27.2)
- [Cloud integration v0.1.15](https://github.com/docker/aci-integration-beta/releases/tag/v0.1.15)

### Bug fixes and minor changes

- Fixed automatic start on log in. See [docker/for-mac#4877](https://github.com/docker/for-mac/issues/4877) and [docker/for-mac#4890](https://github.com/docker/for-mac/issues/4890)

### Known issues

- The `clock_gettime64` system call returns `EPERM` rather than `ENOSYS` 
in i386 images. To work around this issue, disable `seccomp` by using 
the `--privileged` flag. See [docker/for-win#8326](https://github.com/docker/for-win/issues/8326).

## Docker Desktop Community 2.3.0.4
2020-07-27

> [Download](https://desktop.docker.com/mac/stable/46911/Docker.dmg)

### Upgrades

- [Docker 19.03.12](https://github.com/docker/docker-ce/releases/tag/v19.03.12)
- [Docker Compose 1.26.2](https://github.com/docker/compose/releases/tag/1.26.2)
- [Go 1.13.14](https://github.com/golang/go/issues?q=milestone%3AGo1.13.14+label%3ACherryPickApproved)

### Bug fixes and minor changes

- Fixed a privilege escalation vulnerability in `com.docker.vmnetd`. See [CVE-2020-15360](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-15360)
- Fixed an issue with startup when the Kubernetes certificates have expired. See [docker/for-mac#4594](https://github.com/docker/for-mac/issues/4594)
- Fix an incompatibility between `hyperkit` and [osquery](https://osquery.io) which resulted in excessive `hyperkit` CPU usage. See [docker/for-mac#3499](https://github.com/docker/for-mac/issues/3499#issuecomment-619544836)
- Dashboard: Fixed containers logs which were sometimes truncated. Fixes [docker/for-win#5954](https://github.com/docker/for-win/issues/5954)

## Docker Desktop Community 2.3.0.3
2020-05-27

> [Download](https://desktop.docker.com/mac/stable/45519/Docker.dmg)

### Upgrades

- [Linux kernel 4.19.76](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.76-83885d3b4cff391813f4262099b36a529bca2df8-amd64/images/sha256-0214b82436af70054e013ea51cb1fea72bd943d0d6245b6521f1ff09a505c40f?context=repo)

### Bug fixes and minor changes

- Re-added device-mapper to the embedded Linux kernel. Fixes [docker/for-mac#4549](https://github.com/docker/for-mac/issues/4549).
- Fixed `hyperkit` on newer Macs and newer versions of `Hypervisor.framework`. Fixes [docker/for-mac#4562](https://github.com/docker/for-mac/issues/4562).

## Docker Desktop Community 2.3.0.2
2020-05-11

> [Download](https://download.docker.com/mac/stable/45183/Docker.dmg)

### New

Docker Desktop introduces a new onboarding tutorial upon first startup. The Quick Start tutorial guides users to get started with Docker in a few easy steps. It includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

### Upgrades

- [Docker Compose 1.25.5](https://github.com/docker/compose/releases/tag/1.25.5)
- [Go 1.13.10](https://github.com/golang/go/issues?q=milestone%3AGo1.13.10+label%3ACherryPickApproved)
- [Linux kernel 4.19.76](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.76-ce15f646db9b062dc947cfc0c1deab019fa63f96-amd64/images/sha256-6c252199aee548e4bdc8457e0a068e7d8e81c2649d4c1e26e4150daa253a85d8?context=repo)
- LinuxKit [init](https://hub.docker.com/layers/linuxkit/init/1a80a9907b35b9a808e7868ffb7b0da29ee64a95/images/sha256-64cc8fa50d63940dbaa9979a13c362c89ecb4439bcb3ab22c40d300b9c0b597e?context=explore), [runc](https://hub.docker.com/layers/linuxkit/runc/69b4a35eaa22eba4990ee52cccc8f48f6c08ed03/images/sha256-57e3c7cbd96790990cf87d7b0f30f459ea0b6f9768b03b32a89b832b73546280?context=explore) and [containerd](https://hub.docker.com/layers/linuxkit/containerd/09553963ed9da626c25cf8acdf6d62ec37645412/images/sha256-866be7edb0598430709f88d0e1c6ed7bfd4a397b5ed220e1f793ee9067255ff1?context=explore)

### Bug fixes and minor changes

- Reduced the size of the Docker Desktop installer from 708 MB to 456 MB.
- Fixed bug where containers disappeared from the UI when Kubernetes context is invalid. Fixes [docker/for-win#6037](https://github.com/docker/for-win/issues/6037).
- Fixed a file descriptor leak in `vpnkit-bridge`. Fixes [docker/for-win#5841](https://github.com/docker/for-win/issues/5841).
- Added a link to the Edge channel from the UI.
- Made the embedded terminal resizable.
- Fixed a bug where diagnostic upload would fail if the username contained spaces.
- Fixed a bug where the Docker UI could be started without the engine.
- Switched from `ahci-hd` to `virtio-blk` to avoid an AHCI deadlock, see [moby/hyperkit#94](https://github.com/moby/hyperkit/issues/94) and [docker/for-mac#1835](https://github.com/docker/for-mac/issues/1835).
- Fixed an issue where a container port could not be exposed on a specific host IP. See [docker/for-mac#4209](https://github.com/docker/for-mac/issues/4209).
- Removed port probing from dashboard, just unconditionally showing links to ports that should be available. Fixes [docker/for-mac#4264](https://github.com/docker/for-mac/issues/4264).
- Docker Desktop now shares `/var/folders` by default as it stores per-user temporary files and caches.
- Ceph support has been removed from Docker Desktop to save disk space.
- Fixed a performance regression when using shared volumes in 2.2.0.5. Fixes [docker/for-mac#4423].

## Docker Desktop Community 2.2.0.5
2020-04-02

> [Download](https://download.docker.com/mac/stable/43884/Docker.dmg)

### Bug fixes and minor changes

- Removed dangling `/usr/local/bin/docker-machine` symlinks which avoids custom installs of  Docker Machine being accidentally deleted in future upgrades. Note that if you have installed Docker Machine manually, then the install might have followed the symlink and installed Docker Machine in `/Applications/Docker.app`. In this case, you must manually reinstall Docker Machine after installing this version of Docker Desktop. Fixes [docker/for-mac#4208](https://github.com/docker/for-mac/issues/4208).

## Docker Desktop Community 2.2.0.4
2020-03-13

> [Download](https://download.docker.com/mac/stable/43472/Docker.dmg)

### Upgrades

- [Docker 19.03.8](https://github.com/docker/docker-ce/releases/tag/v19.03.8)

### Bug fixes and minor changes

- Kubernetes: Persistent volumes created by claims are now stored in the virtual machine. Fixes [docker/for-win#5665](https://github.com/docker/for-win/issues/5665).
- Fixed an issue which caused Docker Desktop Dashboard to attempt connecting to all exposed ports inside a container. Fixes [docker/for-mac#4264](https://github.com/docker/for-mac/issues/4264).

## Docker Desktop Community 2.2.0.3
2020-02-11

> [Download](https://download.docker.com/mac/stable/42716/Docker.dmg)

### Upgrades

- [Docker Compose 1.25.4](https://github.com/docker/compose/releases/tag/1.25.4)
- [Go 1.12.16](https://golang.org/doc/devel/release.html#go1.12)

## Docker Desktop Community 2.2.0.0
2020-01-21

> [Download](https://download.docker.com/mac/stable/42247/Docker.dmg)

Docker Desktop 2.2.0.0 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker Compose 1.25.2](https://github.com/docker/compose/releases/tag/1.25.2)
- [Kubernetes 1.15.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.15.5)
- Linux kernel 4.19.76
- [QEMU 4.0.1](https://github.com/docker/binfmt)

### New

- **Docker Desktop Dashboard:** The new Docker Desktop **Dashboard** provides a user-friendly interface which enables you to interact with containers and applications, and manage the lifecycle of your applications directly from the UI. In addition, it allows you to access the logs, view container details, and monitor resource utilization to explore the container behavior.
For detailed information about the new Dashboard UI, see [Docker Desktop Dashboard](../desktop/dashboard.md).

- Introduced a new user interface for the Docker Desktop **Preferences** menu.
- The Restart, Reset, and Uninstall options are now available on the **Troubleshoot** menu.
- Added the ability to start and stop existing Compose-based applications and view combined logs in the Docker Desktop **Dashboard** UI.

### Bug fixes and minor changes

- Added missing completions for the `fish` shell for Docker Compose. Fixes [docker/for-mac#3795](https://github.com/docker/for-mac/issues/3795).
- Fixed a bug that did not allow users to copy and paste text in the **Preferences** > **Daemon** window. Fixes [docker/for-mac#3798](https://github.com/docker/for-mac/issues/3798).
- Added support for `Expect: 100-continue` headers in the Docker API proxy. Some HTTP clients such as `curl` send this header when the payload is large, for example, when creating containers. Fixes [moby/moby#39693](https://github.com/moby/moby/issues/39693).
- Added a loading overlay to the **Settings** and **Troubleshoot** windows to prevent editing conflicts.
- Deactivated the **Reset Kubernetes** button when Kubernetes is not activated.
- Improved the navigation in **Settings** and **Troubleshoot** UI.
- Fixed a bug in the UEFI boot menu that sometimes caused Docker Desktop to hang during restart. Fixes [docker/for-mac#2655](https://github.com/docker/for-mac/issues/2655) and [docker/for-mac#3921](https://github.com/docker/for-mac/issues/3921).
- Docker Desktop now allows users to access the host’s SSH agent inside containers. Fixes [docker/for-mac#410](https://github.com/docker/for-mac/issues/410)
- Docker Machine is no longer included in the Docker Desktop installer. You can download it separately from the [Docker Machine releases](https://github.com/docker/machine/releases) page.
- Fixed an issue that caused VMs running on older hardware with macOS Catalina to fail on startup with the error `processor does not support desired secondary processor-based controls`.
- Fixed port forwarding when containers are using `overlay` networks.
- Fixed a container start error when a container has more than one port with an arbitrary or not-yet-configured external port number. For example, `docker run -p 80 -p 443 nginx`. Fixes [docker/for-win#4935](https://github.com/docker/for-win/issues/4935) and [docker/compose#6998](https://github.com/docker/compose/issues/6998).
- Fixed an issue that occurs when sharing overlapping directories.
- Fixed a bug that prevented users from changing the location of the VM disk image.
- Docker Desktop does not inject `inotify` events on directories anymore as these can cause mount points to disappear inside containers. Fixes [docker/for-mac#3976](https://github.com/docker/for-mac/issues/3976).
- Fixed an issue that caused Docker Desktop to fail on startup when there is an incomplete Kubernetes config file.
- Fixed an issue where attempts to log into Docker through Docker Desktop could sometimes fail with the `Incorrect authentication credentials` error. Fixes [docker/for-mac#4010](https://github.com/docker/for-mac/issues/4010).

### Known issues

- When you start a Docker Compose application and then start a Docker App which has the same name as the Compose application, Docker Desktop displays only one application on the Dashboard. However, when you expand the application, containers that belong to both applications are displayed on the Dashboard.

- When you deploy a Docker App with multiple containers on Kubernetes, Docker Desktop displays each Pod as an application on the Dashboard.

## Docker Desktop Community 2.1.0.5
2019-11-18

[Download](https://download.docker.com/mac/stable/40693/Docker.dmg)

Docker Desktop 2.1.0.5 contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker 19.03.5](https://github.com/docker/docker-ce/releases/tag/v19.03.5)
- [Kubernetes 1.14.8](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.8)
- [Go 1.12.13](https://golang.org/doc/devel/release.html#go1.12)

## Docker Desktop Community 2.1.0.4
2019-10-21

[Download](https://download.docker.com/mac/stable/39773/Docker.dmg)

### Upgrades

- [Docker 19.03.4](https://github.com/docker/docker-ce/releases/tag/v19.03.4)
- [Kubernetes 1.14.7](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.7)
- [Go 1.12.10](https://github.com/golang/go/issues?q=milestone%3AGo1.12.10+label%3ACherryPickApproved)
- [Kitematic 0.17.9](https://github.com/docker/kitematic/releases/tag/v0.17.9)

### New

Docker Desktop now enables you to sign into Docker Hub using two-factor authentication.
For more information, see [Two-factor authentication](index.md#docker-hub).

## Docker Desktop Community 2.1.0.3
2019-09-16

[Download](https://download.docker.com/mac/stable/38240/Docker.dmg)

### Upgrades

- [Kitematic 0.17.8](https://github.com/docker/kitematic/releases/tag/v0.17.8)

### Bug fixes and minor changes

- All binaries included in Docker Desktop are now notarized so that they can run on macOS Catalina. For more information, see [Notarization Requirement for Mac Software](https://developer.apple.com/news/?id=06032019i).

## Docker Desktop Community 2.1.0.2
2019-09-04

[Download](https://download.docker.com/mac/stable/37877/Docker.dmg)

Docker Desktop 2.1.0.2 contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker 19.03.2](https://github.com/docker/docker-ce/releases/tag/v19.03.2)
- [Kubernetes 1.14.6](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.6)
- [Go 1.12.9](https://github.com/golang/go/issues?q=milestone%3AGo1.12.9+label%3ACherryPickApproved)
- [Docker Machine 0.16.2](https://github.com/docker/machine/releases/tag/v0.16.2)

## Docker Desktop Community 2.1.0.1
2019-08-08

[Download](https://download.docker.com/mac/stable/37199/Docker.dmg)

Note that you must sign in and create a Docker ID in order to download Docker Desktop.

### Upgrades

* [Docker 19.03.1](https://github.com/docker/docker-ce/releases/tag/v19.03.1)
* [Docker Compose 1.24.1](https://github.com/docker/compose/releases/tag/1.24.1)
* [Kubernetes 1.14.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.3)
* [Compose on Kubernetes 0.4.23](https://github.com/docker/compose-on-kubernetes/releases/tag/v0.4.23)
* [Docker Machine 0.16.1](https://github.com/docker/machine/releases/tag/v0.16.1)
* [linuxkit v0.7](https://github.com/linuxkit/linuxkit/releases/tag/v0.7)
* Linux Kernel 4.9.184
* [Kitematic 0.17.6](https://github.com/docker/kitematic/releases/tag/v0.17.6)
* [Qemu 4.0.0](https://github.com/docker/binfmt) for cross compiling for ARM
* [Alpine 3.10](https://alpinelinux.org/posts/Alpine-3.10.0-released.html)
* [Docker Credential Helpers 0.6.3](https://github.com/docker/docker-credential-helpers/releases/tag/v0.6.3)
* [Hyperkit v0.20190802](https://github.com/moby/hyperkit/releases/tag/v0.20190802)

### New

* Selecting the ‘Experimental features’ checkbox in the Daemon **Preferences** menu turns on experimental features for Docker daemon and Docker CLI.
* Improved the reliability of `com.docker.osxfs trace` performance profiling command. Users can now run the `com.docker.osxfs trace --summary` option for a high-level summary of operations, instead of receiving a trace of all operations.
* Docker Desktop now supports large lists of DNS resource records on Mac.  Fixes [docker/for-mac#2160](https://github.com/docker/for-mac/issues/2160#issuecomment-431571031).

### Experimental

> Experimental features provide early access to future product functionality. These features are intended for testing and feedback only as they may change between releases without warning or can be removed entirely from a future release. Experimental features must not be used in production environments. Docker does not offer support for experimental features.

Docker Desktop Community 2.1.0.0 contains the following experimental features.

* Docker App: Docker App is a CLI plugin that helps configure, share, and install applications. For more information, see [Working with Docker App](/app/working-with-app/).
* Docker Buildx: Docker Buildx is a CLI plugin for extended build capabilities with BuildKit. For more information, see [Working with Docker Buildx](/buildx/working-with-buildx/).

### Bug fixes and minor changes

* Docker Desktop now allows users to expose privileged UDP ports. [docker/for-mac#3775](https://github.com/docker/for-mac/issues/3775)
* Fixed an issue where running some Docker commands can fail if you are not using Credential Helpers. [docker/for-mac#3785](https://github.com/docker/for-mac/issues/3785)
* Changed the host's kubernetes context so that `docker run -v .kube:kube ... kubectl` works.
* Restricted the `cluster-admin` role on local Kubernetes cluster to `kube-system` namespace.
* Reduced the VM startup time. swap is not created every time a virtual machine boots.
* Fixed Kubernetes installation with VPNkit subnet.
* Fixed a bug where the process output was not redirected to stdout when gathering diagnostics on Windows, which sometimes resulted in a crash.
* Added `/etc/machine-id` to the virtual machine. Fixes [docker/for-mac#3554](https://github.com/docker/for-mac/issues/3554).
* Docker Desktop does not send DNS queries for `docker-desktop.<domain>` every 10s. It now relies on the host’s DNS domain search order rather than trying to replicate it inside the VM.
* Removed the ability to log in using email address as a username as the Docker command line does not support this.
* Docker Desktop now allows running a Docker registry inside a container. Fixes [docker/for-mac#3611](https://github.com/docker/for-mac/issues/3611).
* Fixed a stability issue with the DNS resolver.
* Docker Desktop truncates UDP DNS responses which are over 512 bytes in sizes.
* Fixed port 8080 that was used on localhost when starting Kubernetes. Fixes [docker/for-mac#3522](https://github.com/docker/for-mac/issues/3522).
* Improved error messaging: Docker Desktop does not prompt users to run diagnostics, or to reset to factory default when it is not appropriate.
* Kubernetes: Docker Desktop now uses the default maximum number of pods for kubelet. [docker/for-mac#3453](https://github.com/docker/for-mac/issues/3453).
* Fixed DockerHelper crash. [docker/for-mac#3470](https://github.com/docker/for-mac/issues/3470).
* Fixed binding of privileged ports with specified IP. [docker/for-mac#3464](https://github.com/docker/for-mac/issues/3464)
* Fixed service log collection in diagnostics.
* Docker Desktop now gathers `/etc/hosts` to help with diagnostics.
* When two services have a common exposed port, Docker Desktop now exposes the available ports for the second service. [docker/for-mac#3438](https://github.com/docker/for-mac/issues/3438).
* Docker Desktop ensures localhost resolves to 127.0.0.1. This is related to [docker/for-mac#2990](https://github.com/docker/for-mac/issues/2990#issuecomment-443097942), [docker/for-mac#3383](https://github.com/docker/for-mac/issues/3383).

### Docker Community Edition 2.0.0.3 2019-02-15

[Download](https://download.docker.com/mac/stable/31259/Docker.dmg)

* Upgrades
  - [Docker 18.09.2](https://github.com/docker/docker-ce/releases/tag/v18.09.2), fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

### Docker Community Edition 2.0.0.2 2019-01-16

[Download](https://download.docker.com/mac/stable/30215/Docker.dmg)

* Upgrades
  - [Docker 18.09.1](https://github.com/docker/docker-ce/releases/tag/v18.09.1)
  - [Docker Machine 0.16.1](https://github.com/docker/machine/releases/tag/v0.16.1)
  - [Kubernetes 1.10.11](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.10.md#v11011), fixes [CVE-2018-1002105](https://github.com/kubernetes/kubernetes/issues/71411)
  - [Kitematic 0.17.6](https://github.com/docker/kitematic/releases/tag/v0.17.6)
  - Golang 1.10.6, fixes CVEs: [CVE-2018-16875](https://www.cvedetails.com/cve/CVE-2018-16875), [CVE-2018-16873](https://www.cvedetails.com/cve/CVE-2018-16873) and [CVE-2018-16874](https://www.cvedetails.com/cve/CVE-2018-16874)

* Bug fixes and minor changes
  - Add 18.09 missing daemon options

## Stable Releases of 2018

### Docker Community Edition 2.0.0.0-mac81 2018-12-07

[Download](https://download.docker.com/mac/stable/29211/Docker.dmg)

* Upgrades
  - [Docker compose 1.23.2](https://github.com/docker/compose/releases/tag/1.23.2)

### Docker Community Edition 2.0.0.0-mac78 2018-11-19

[Download](https://download.docker.com/mac/stable/28905/Docker.dmg)

* Upgrades
  - [Docker 18.09.0](https://github.com/docker/docker-ce-packaging/releases/tag/v18.09.0)
  - [Docker compose 1.23.1](https://github.com/docker/compose/releases/tag/1.23.1)
  - [Docker Machine 0.16.0](https://github.com/docker/machine/releases/tag/v0.16.0)
  - [Kitematic 0.17.5](https://github.com/docker/kitematic/releases/tag/v0.17.5)
  - Linux Kernel 4.9.125

* New
  - New version scheme

* Deprecation
  - Removed support of AUFS
  - Removed support of OSX 10.11

* Bug fixes and minor changes
  - Fix appearance in dark mode for OSX 10.14 (Mojave)
  - VPNKit: Improved scalability of port forwarding. Related to [docker/for-mac#2841](https://github.com/docker/for-mac/issues/2841)
  - VPNKit: Limit the size of the UDP NAT table. This ensures port forwarding and regular TCP traffic continue even when running very chatty UDP protocols.
  - Ensure Kubernetes can be installed when using a non-default internal IP subnet.
  - Fix panic in diagnose

### Docker Community Edition 18.06.1-ce-mac73 2018-08-29

[Download](https://download.docker.com/mac/stable/26764/Docker.dmg)

* Upgrades
  - [Docker 18.06.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.06.1-ce)

* Bug fixes and minor changes
  - Fix local DNS failing to resolve inside containers.

### Docker Community Edition 18.06.0-ce-mac70 2018-07-25

[Download](https://download.docker.com/mac/stable/26399/Docker.dmg)

* Upgrades
  - [Docker 18.06.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.06.0-ce)
  - [Docker Machine 0.15.0](https://github.com/docker/machine/releases/tag/v0.15.0)
  - [Docker compose 1.22.0](https://github.com/docker/compose/releases/tag/1.22.0)
  - [LinuxKit v0.5](https://github.com/linuxkit/linuxkit/releases/tag/v0.5)
  - Linux Kernel 4.9.93 with CEPH, DRBD, RBD, MPLS_ROUTING and MPLS_IPTUNNEL enabled

* New
  - Kubernetes Support. You can now run a single-node Kubernetes cluster from the "Kubernetes" Pane in Docker For Mac Preferences and use kubectl commands as well as docker commands. See [https://docs.docker.com/docker-for-mac/kubernetes/](kubernetes.md)
  - Add an experimental SOCKS server to allow access to container networks, see [docker/for-mac#2670](https://github.com/docker/for-mac/issues/2670#issuecomment-372365274). Also see [docker/for-mac#2721](https://github.com/docker/for-mac/issues/2721)
  - Re-enable raw as the default disk format for users running macOS 10.13.4 and higher. Note this change only takes effect after a "reset to factory defaults" or "remove all data" (from the Whale menu -> Preferences -> Reset). Related to [docker/for-mac#2625](https://github.com/docker/for-mac/issues/2625)

* Bug fixes and minor changes
  - AUFS storage driver is deprecated in Docker Desktop and AUFS support will be removed in the next major release. You can continue with AUFS in Docker Desktop 18.06.x, but you will need to reset disk image (in Preferences > Reset menu) before updating to the next major update. You can check documentation to [save images](https://docs.docker.com/engine/reference/commandline/save/#examples) and [backup volumes](https://docs.docker.com/storage/volumes/#backup-restore-or-migrate-data-volumes)
  - OS X El Captain 10.11 is deprecated in Docker Desktop. You will not be able to install updates after Docker Desktop 18.06.x. We recommend upgrading to the latest version of macOS.
  - Fix bug which would cause VM logs to be written to RAM rather than disk in some cases, and the VM to hang. See [docker/for-mac#2984](https://github.com/docker/for-mac/issues/2984)
  - Fix network connection leak triggered by haproxy TCP health-checks [docker/for-mac#1132](https://github.com/docker/for-mac/issues/1132)
  - Better message to reset vmnetd when it's disabled. See [docker/for-mac#3035](https://github.com/docker/for-mac/issues/3035)
  - Fix VPNKit memory leak. Fixes [moby/vpnkit#371](https://github.com/moby/vpnkit/issues/371)
  - Virtual Machine default disk path is stored relative to $HOME. Fixes [docker/for-mac#2928](https://github.com/docker/for-mac/issues/2928), [docker/for-mac#1209](https://github.com/docker/for-mac/issues/1209)
  - Use Simple NTP to minimise clock drift between the VM and the host. Fixes [docker/for-mac#2076](https://github.com/docker/for-mac/issues/2076)
  - Fix a race between calling stat on a file and calling close of a file descriptor referencing the file that could result in the stat failing with EBADF (often presented as "File not found"). Fixes [docker/for-mac#2870](https://github.com/docker/for-mac/issues/2870)
  - Do not allow install of Docker for Mac on macOS Yosemite 10.10, this version is not supported since Docker for Mac 17.09.0.
  - Fix button order in reset dialog windows. Fixes [docker/for-mac#2827](https://github.com/docker/for-mac/issues/2827)
  - Fix upgrade straight from pre-17.12 versions where Docker for Mac cannot restart once the upgrade has been performed. Fixes [docker/for-mac#2739](https://github.com/docker/for-mac/issues/2739)
  - Added log rotation for docker-ce logs inside the virtual machine.

### Docker Community Edition 18.03.1-ce-mac65 2018-04-30

[Download](https://download.docker.com/mac/stable/24312/Docker.dmg)

* Upgrades
  - [Docker 18.03.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.1-ce)
  - [Docker compose 1.21.1](https://github.com/docker/compose/releases/tag/1.21.1)
  - [Notary 0.6.1](https://github.com/docker/notary/releases/tag/v0.6.1)

* Bug fixes and minor changes
  - Fix Docker for Mac not starting due to socket file paths being too long (typically HOME folder path being too long). Fixes [docker/for-mac#2727](https://github.com/docker/for-mac/issues/2727), [docker/for-mac#2731](https://github.com/docker/for-mac/issues/2731).

### Docker Community Edition 18.03.1-ce-mac64 2018-04-26

[Download](https://download.docker.com/mac/stable/24245/Docker.dmg)

* Upgrades
  - [Docker 18.03.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.1-ce)
  - [Docker compose 1.21.0](https://github.com/docker/compose/releases/tag/1.21.0)
  - [Notary 0.6.1](https://github.com/docker/notary/releases/tag/v0.6.1)

* Bug fixes and minor changes
  - Fix Docker for Mac not starting due to socket file paths being too long (typically HOME folder path being too long). Fixes [docker/for-mac#2727](https://github.com/docker/for-mac/issues/2727), [docker/for-mac#2731](https://github.com/docker/for-mac/issues/2731).

### Docker Community Edition 18.03.0-ce-mac60 2018-03-30

[Download](https://download.docker.com/mac/stable/23751/Docker.dmg)

* Bug fixes and minor changes
  - Fix Upgrade straight from 17.09 versions where Docker for Mac cannot restart once the upgrade has been performed. Fixes [docker/for-mac#2739](https://github.com/docker/for-mac/issues/2739)

### Docker Community Edition 18.03.0-ce-mac59 2018-03-26

[Download](https://download.docker.com/mac/stable/23608/Docker.dmg)

* Upgrades
  - [Docker 18.03.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce)
  - [Docker Machine 0.14.0](https://github.com/docker/machine/releases/tag/v0.14.0)
  - [Docker compose 1.20.1](https://github.com/docker/compose/releases/tag/1.20.1)
  - [Notary 0.6.0](https://github.com/docker/notary/releases/tag/v0.6.0)
  - Linux Kernel 4.9.87
  - AUFS 20180312

* New
  - VM Swap size can be changed in settings. See [docker/for-mac#2566](https://github.com/docker/for-mac/issues/2566), [docker/for-mac#2389](https://github.com/docker/for-mac/issues/2389)
  - New menu item to restart Docker.
  - Support NFS Volume sharing.
  - The directory holding the disk images has been renamed (from `~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux` to ~/Library/Containers/com.docker.docker/Data/vms/0`).

* Bug fixes and minor changes
  - Fix daemon not starting properly when setting TLS-related options. Fixes [docker/for-mac#2663](https://github.com/docker/for-mac/issues/2663)
  - DNS name `host.docker.internal` should be used for host resolution from containers. Older aliases (still valid) are deprecated in favor of this one. (See https://tools.ietf.org/html/draft-west-let-localhost-be-localhost-06).
  - Fix for the HTTP/S transparent proxy when using "localhost" names (e.g. `host.docker.internal`).
  - Fix empty registry added by mistake in some cases in the Preference Daemon Pane. Fixes [docker/for-mac#2537](https://github.com/docker/for-mac/issues/2537)
  - Clearer error message when incompatible hardware is detected.
  - Fix some cases where selecting "Reset" after an error did not reset properly.
  - Fix incorrect NTP config. Fixes [docker/for-mac#2529](https://github.com/docker/for-mac/issues/2529)
  - Migration of Docker Toolbox images is not proposed anymore in Docker For Mac installer (still possible to [migrate Toolbox images manually](docker-toolbox.md#migrating-from-docker-toolbox-to-docker-desktop-on-mac) ).

### Docker Community Edition 17.12.0-ce-mac55 2018-02-27

[Download](https://download.docker.com/mac/stable/23011/Docker.dmg)

* Bug fixes and minor changes
  - Revert the default disk format to qcow2 for users running macOS 10.13 (High Sierra). There are confirmed reports of file corruption using the raw format which uses sparse files on APFS. Note this change only takes effect after a reset to factory defaults (from the Whale menu -> Preferences -> Reset). Related to [docker/for-mac#2625](https://github.com/docker/for-mac/issues/2625)
  - Fix VPNKit proxy for docker.for.mac.http.internal.

### Docker Community Edition 17.12.0-ce-mac49 2018-01-19

[Download](https://download.docker.com/mac/stable/21805/Docker.dmg)

* Bug fixes and minor changes
  - Fix error during resize/create Docker.raw disk image in some cases. Fixes [docker/for-mac#2383](https://github.com/docker/for-mac/issues/2383), [docker/for-mac#2447](https://github.com/docker/for-mac/issues/2447), [docker/for-mac#2453], (https://github.com/docker/for-mac/issues/2453), [docker/for-mac#2420](https://github.com/docker/for-mac/issues/2420)
  - Fix additional allocated disk space not available in containers. Fixes [docker/for-mac#2449](https://github.com/docker/for-mac/issues/2449)
  - Vpnkit port max idle time default restored to 300s. Fixes [docker/for-mac#2442](https://github.com/docker/for-mac/issues/2442)
  - Fix using an HTTP proxy with authentication. Fixes [docker/for-mac#2386](https://github.com/docker/for-mac/issues/2386)
  - Allow HTTP proxy excludes to be written as .docker.com as well as *.docker.com
  - Allow individual IP addresses to be added to HTTP proxy excludes.
  - Avoid hitting DNS timeouts when querying docker.for.mac.* when the upstream DNS servers are slow or missing.

### Docker Community Edition 17.12.0-ce-mac47 2018-01-12

[Download](https://download.docker.com/mac/stable/21698/Docker.dmg)

* Bug fixes and minor changes
  - Fix for `docker push` to an insecure registry. Fixes [docker/for-mac#2392](https://github.com/docker/for-mac/issues/2392)
  - Separate internal ports used to proxy HTTP and HTTPS content.

### Docker Community Edition 17.12.0-ce-mac46 2018-01-09

[Download](https://download.docker.com/mac/stable/21698/Docker.dmg)

* Upgrades
  - [Docker 17.12.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce)
  - [Docker compose 1.18.0](https://github.com/docker/compose/releases/tag/1.18.0)
  - [Docker Machine 0.13.0](https://github.com/docker/machine/releases/tag/v0.13.0)
  - Linux Kernel 4.9.60

* New
  - VM entirely built with Linuxkit
  - VM disk size can be changed in disk preferences. (See [docker/for-mac#1037](https://github.com/docker/for-mac/issues/1037))
  - For systems running APFS on SSD on High Sierra, use `raw` format VM disks by default. This improves disk throughput (from 320MiB/sec to 600MiB/sec in `dd` on a 2015 MacBook Pro) and disk space handling.
  Existing disks are kept in qcow format, if you want to switch to raw format you need to "Remove all data" or "Reset to factory defaults". See https://docs.docker.com/docker-for-mac/faqs/#disk-usage
  - DNS name `docker.for.mac.host.internal` should be used instead of `docker.for.mac.localhost` (still valid) for host resolution from containers, since since there is an RFC banning the use of subdomains of localhost. See  https://tools.ietf.org/html/draft-west-let-localhost-be-localhost-06.

* Bug fixes and minor changes
  - Display various component versions in About box.
  - Avoid VM reboot when changing host proxy settings.
  - Don't break HTTP traffic between containers by forwarding them via the external proxy. (See [docker/for-mac#981](https://github.com/docker/for-mac/issues/981))
  - Filesharing settings are now stored in settings.json.
  - Daemon restart button has been moved to settings / Reset Tab.
  - Better VM state handling & error messages in case of VM crashes.
  - Fix login into private repository with certificate issue. (See [docker/for-mac#2201](https://github.com/docker/for-mac/issues/2201))

## Stable Releases of 2017
### Docker Community Edition 17.09.1-ce-mac42 2017-12-11

[Download](https://download.docker.com/mac/stable/21090/Docker.dmg)

* Upgrades
  - [Docker 17.09.1-ce](https://github.com/docker/docker-ce/releases/tag/v17.09.1-ce)
  - [Docker compose 1.17.1](https://github.com/docker/compose/releases/tag/1.17.1)
  - [Docker Machine 0.13.0](https://github.com/docker/machine/releases/tag/v0.13.0)

* Bug fixes and minor changes
  - Fix bug not allowing to move qcow disk in some cases.

### Docker Community Edition 17.09.0-ce-mac35 2017-10-06

[Download](https://download.docker.com/mac/stable/19611/Docker.dmg)

* Bug fix
  - Fix Docker For Mac unable to start in some cases : removed use of libgmp sometimes causing the vpnkit process to die.

### Docker Community Edition 17.09.0-ce-mac33 2017-10-03

[Download](https://download.docker.com/mac/stable/19543/Docker.dmg)

* Bug fix
  - Do not show Toolbox migration assistant when there are existing Docker For Mac data.

### Docker Community Edition 17.09.0-ce-mac32 2017-10-02

[Download](https://download.docker.com/mac/stable/19506/Docker.dmg)

* Upgrades
  - [Docker 17.09.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce)
  - [Docker Compose 1.16.1](https://github.com/docker/compose/releases/tag/1.16.1)
  - [Docker Machine 0.12.2](https://github.com/docker/machine/releases/tag/v0.12.2)
  - [Docker Credential Helpers 0.6.0](https://github.com/docker/docker-credential-helpers/releases/tag/v0.6.0)
  - Linux Kernel 4.9.49
  - AUFS 20170911
  - DataKit update (fix instability on High Sierra)

* New
  - Add daemon options validation
  - VPNKit: add support for ping!
  - VPNKit: add slirp/port-max-idle-timeout to allow the timeout to be adjusted or even disabled
  - VPNKit: bridge mode is default everywhere now
  - Transparent proxy using macOS system proxies (if defined) directly
  - GUI settings are now stored in ~/Library/Group\ Containers/group.com.docker/settings.json. daemon.json in now a file in ~/.docker/
  - You can now change the default IP address used by Hyperkit if it collides with your network

* Bug fixes and minor changes
  - Fix instability on High Sierra (docker/for-mac#2069, docker/for-mac#2062, docker/for-mac#2052, docker/for-mac#2029, docker/for-mac#2024)
  - Fix password encoding/decoding (docker/for-mac#2008, docker/for-mac#2016, docker/for-mac#1919, docker/for-mac#712, docker/for-mac#1220).
  - Kernel: Enable TASK_XACCT and TASK_IO_ACCOUNTING (docker/for-mac#1608)
  - Rotate logs in the VM more often
  - VPNKit: change protocol to support error messages reported back from the server
  - VPNKit: fix a bug which causes a socket to leak if the corresponding TCP connection is idle
    for more than 5 minutes (related to [docker/for-mac#1374](https://github.com/docker/for-mac/issues/1374))
  - VPNKit: improve the logging around the Unix domain socket connections
  - VPNKit: automatically trim whitespace from int or bool database keys
  - Diagnose can be cancelled & Improved help information. Fixes docker/for-mac#1134, docker/for-mac#1474
  - Support paging of docker-cloud repositories & orgs. Fixes docker/for-mac#1538

### Docker Community Edition 17.06.2-ce-mac27 2017-09-06

[Download](https://download.docker.com/mac/stable/19124/Docker.dmg)

* Upgrades
  - [Docker 17.06.2-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.2-ce)
  - [Docker Machine 0.12.2](https://github.com/docker/machine/releases/tag/v0.12.2)

### Docker Community Edition 17.06.1-ce-mac24, 2017-08-21

[Download](https://download.docker.com/mac/stable/18950/Docker.dmg)

**Upgrades**
- [Docker 17.06.1-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v17.06.1-ce-rc1)
- Linux Kernel 4.9.36
- AUFS 20170703

**Bug fixes and minor changes**

- DNS Fixes. Fixes [docker/for-mac#1763](https://github.com/docker/for-mac/issues/176), [docker/for-mac#1811](https://github.com/docker/for-mac/issues/1811), [docker/for-mac#1803](https://github.com/docker/for-mac/issues/1803)

- Avoid unnecessary VM reboot (when changing proxy exclude, but no proxy set). Fixes [docker/for-mac#1809](https://github.com/docker/for-mac/issues/1809), [docker/for-mac#1801](https://github.com/docker/for-mac/issues/1801)

### Docker Community Edition 17.06.0-ce-mac18, 2017-06-28

[Download](https://download.docker.com/mac/stable/18433/Docker.dmg)

**Upgrades**

- [Docker 17.06.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce)
- [Docker Credential Helpers 0.5.2](https://github.com/docker/docker-credential-helpers/releases/tag/v0.5.2)
- [Docker Machine 0.12.0](https://github.com/docker/machine/releases/tag/v0.12.0)
- [Docker compose 1.14.0](https://github.com/docker/compose/releases/tag/1.14.0)
- qcow-tool v0.10.0 (improve the performance of `compact`: mirage/ocaml-qcow#94)
- OSX Yosemite 10.10 is marked as deprecated
- Linux Kernel 4.9.31

**New**

- Integration with Docker Cloud: control remote Swarms from the local CLI and view your repositories.
- GUI Option to opt out of credential store
- GUI option to reset Docker data without losing all settings (fixes [docker/for-mac#1309](https://github.com/docker/for-mac/issues/1309))
- Add an experimental DNS name for the host: `docker.for.mac.localhost`
- Support for client (i.e. "login") certificates for authenticating registry access (fixes [docker/for-mac#1320](https://github.com/docker/for-mac/issues/1320))
- OSXFS: support for `cached` mount flag to improve performance of macOS mounts when strict consistency is not necessary

**Bug fixes and minor changes**

- Resync HTTP(S) proxy settings on application start
- Interpret system proxy setting of `localhost` correctly (see [docker/for-mac#1511](https://github.com/docker/for-mac/issues/1511))
- All Docker binaries bundled with Docker for Mac are now signed
- Display all Docker Cloud organizations and repositories in the whale menu (fixes [docker/for-mac#1538 ](https://github.com/docker/for-mac/issues/1538))
- OSXFS: improved latency for many common operations, such as read and write, by approximately 25%
- Fixed GUI crash when text table view was selected and windows re-opened (fixes [docker/for-mac#1477](https://github.com/docker/for-mac/issues/1477))
- Reset to default / uninstall also remove `config.json` and `osxkeychain` credentials
- More detailed VirtualBox uninstall requirements ( [docker/for-mac#1343](https://github.com/docker/for-mac/issues/1343))
- Request time sync after waking up to improve [docker/for-mac#17](https://github.com/docker/for-mac/issues/17)
- VPNKit: Improved DNS timeout handling (fixes [docker/for-mac#202](https://github.com/docker/vpnkit/issues/202))
- VPNKit: Use DNSServiceRef API by default (only enabled for new installs or after factory reset)
- Add a reset to factory defaults button when application crashes
- Toolbox import dialog now defaults to "Skip"
- Buffered data should be treated correctly when Docker client requests are upgraded to raw streams
- Removed an error message from the output related to experimental features handling
- `vmnetd` should not crash when user home directory is on an external drive
- Improved settings database schema handling
- Disk trimming should work as expected
- Diagnostics now contains more settings data

### Docker Community Edition 17.03.1-ce-mac12, 2017-05-12

[Download](https://download.docker.com/mac/stable/17661/Docker.dmg)

**Upgrades**

- Security fix for CVE-2017-7308

### Docker Community Edition 17.03.1-ce-mac5, 2017-03-29

[Download](https://download.docker.com/mac/stable/16048/Docker.dmg)

**Upgrades**

- [Docker Credential Helpers 0.4.2](https://github.com/docker/docker-credential-helpers/releases/tag/v0.4.2)


### Docker Community Edition 17.03.0-ce-mac1, 2017-03-02

[Download](https://download.docker.com/mac/stable/15583/Docker.dmg)

**New**

- Renamed to Docker Community Edition
- Integration with Docker Cloud: control remote Swarms from the local CLI and view your repositories. This feature is going to be rolled out to all users progressively
- Docker will now securely store your IDs in the macOS keychain

**Upgrades**

- [Docker 17.03.0-ce](https://github.com/docker/docker/releases/tag/v17.03.0-ce)
- [Docker Compose 1.11.2](https://github.com/docker/compose/releases/tag/1.11.2)
- [Docker Machine 0.10.0](https://github.com/docker/machine/releases/tag/v0.10.0)
- Linux kernel 4.9.12

**Bug fixes and minor changes**

- Allow to reset faulty daemon.json through a link in advanced subpanel
- More options when moving disk image
- Add link to experimental features
- Filesharing and daemon table empty fields are editable
- Hide restart button in settings window
- Fix bug where update window hides when app not focused
- Don't use port 4222 inside the Linux VM
- Add page_poison=1 to boot args
- Add a new disk flushing option
- DNS forwarder ignores responses from malfunctioning servers (docker/for-mac#1025)
- DNS forwarder send all queries in parallel, process results in order
- DNS forwarder includes servers with zones in general searches (docker/for-mac#997)
- Parses aliases from /etc/hosts (docker/for-mac#983)
- Can resolve DNS requests via servers listed in the /etc/resolver directory on the host
- Limit vCPUs to 64
- Fixed for swap not being mounted
- Fixed aufs xattr delete issue (docker/docker#30245)
- osxfs: Catch EPERM when reading extended attributes of non-files
- VPNKit: Fixed unmarshalling of DNS packets containing pointers to pointers to labels
- VPNKit: Set the Recursion Available bit on DNS responses from the cache
- VPNKit: Avoid diagnostics to capture too much data
- VPNKit: Fixed a source of occasional packet loss (truncation) on the virtual ethernet link
- HyperKit: Dump guest physical and linear address from VMCS when dumping state
- HyperKit: Kernel boots with panic=1 arg

### Docker for Mac 1.13.1, 2017-02-09

[Download](https://download.docker.com/mac/stable/15353/Docker.dmg)

**Upgrades**

- [Docker 1.13.1](https://github.com/docker/docker/releases/tag/v1.13.1)
- [Docker Compose 1.11.1](https://github.com/docker/compose/releases/tag/1.11.1)
- Linux kernel 4.9.8

**Bug fixes and minor changes**

- Add link to experimental features
- New 1.13 cancellable operations should now be properly handled by the Docker for desktop
- `daemon.json` should render nicely in the UI
- Allow to reset faulty `daemon.json` through a link in advanced subpanel

### Docker for Mac 1.13.0, 2017-01-19

[Download](https://download.docker.com/mac/stable/15072/Docker.dmg)

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
- Dedicated preference pane for advanced configuration of the
docker daemon (edit `daemon.json`)
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
- Dedicated preference pane for HTTP proxy settings
- Dedicated preference pane for CPU & Memory computing resources
- Privacy settings moved to the general preference pane
- Fixed an issue where the preference pane disappeared when the welcome whale menu was closed
- HyperKit: code cleanup and minor fixes
- Improvements to Logging and Diagnostics
- osxfs: switch to libev/kqueue to improve latency
- VPNKit: improvements to DNS handling
- VPNKit: Improved diagnostics
- VPNKit: Forwarded UDP datagrams should have correct source port numbers
- VPNKit: add a local cache of DNS responses
- VPNKit: If one request fails, allow other concurrent requests to succeed.
  For example this allows IPv4 servers to work even if IPv6 is broken.
- VPNKit: Fix bug that could cause the connection tracking to
underestimate the number of active connections

## Stable Releases of 2016
### Docker for Mac 1.12.5, 2016-12-20

[Download](https://download.docker.com/mac/stable/14777/Docker.dmg)

**Upgrades**

- Docker 1.12.5
- Docker Compose 1.9.0

### Skipped Docker for Mac 1.12.4

We did not distribute a 1.12.4 stable release

### Docker for Mac 1.12.3, 2016-11-09

[Download](https://download.docker.com/mac/stable/13776/Docker.dmg)

**Upgrades**

- Docker 1.12.3
- Linux Kernel 4.4.27
- Notary 0.4.2
- Docker Machine 0.8.2
- Docker Compose 1.8.1
- Kernel vsock driver v7
- aufs 20160912

**Bug fixes and minor changes**

**General**

- Fixed an issue where the whale animation during setting change was inconsistent

- Fixed an issue where some windows stayed hidden behind another app

- Fixed an issue where the Docker status would continue to be yellow/animated after the VM had started correctly

- Fixed an issue where Docker for Mac was incorrectly reported as updated

- Channel is now displayed in About box

- Crash reports are sent over Bugsnag rather than HockeyApp

- Fixed an issue where some windows did not claim focus correctly

- Added UI when switching channel to prevent user losing containers and settings

- Check disk capacity before Toolbox import

- Import certificates in `etc/ssl/certs/ca-certificates.crt`

- disk: make the "flush" behaviour configurable for database-like workloads. This works around a performance regression in 1.12.1.

**Networking**

- Proxy: Fixed application of system or custom proxy settings over container restart

- DNS: reduce the number of UDP sockets consumed on the host

- VPNkit: improve the connection-limiting code to avoid running out of sockets on the host

- UDP: handle diagrams bigger than 2035, up to the configured macOS kernel limit

- UDP: make the forwarding more robust; drop packets and continue rather than stopping

**File sharing**

- osxfs: Fixed the prohibition of chown on read-only or mode 0 files, (fixes
  [docker/for-mac#117](https://github.com/docker/for-mac/issues/117),
  [docker/for-mac#263](https://github.com/docker/for-mac/issues/263),
  [docker/for-mac#633](https://github.com/docker/for-mac/issues/633))

- osxfs: Fixed race causing some reads to run forever

- osxfs: Fixed a simultaneous volume mount race which can result in a crash

**Moby**

- Increase default ulimit for memlock (fixes [docker/for-mac#801](https://github.com/docker/for-mac/issues/801))

### Docker for Mac 1.12.1, 2016-09-16

[Download](https://download.docker.com/mac/stable/1.12.1.12133/Docker.dmg)

**New**

* Support for macOS 10.12 Sierra

**Upgrades**

* Docker 1.12.1
* Docker machine 0.8.1
* Linux kernel 4.4.20
* aufs 20160905

**Bug fixes and minor changes**

**General**

 * Fixed communications glitch when UI talks to com.docker.vmnetd
 Fixes [docker/for-mac#90](https://github.com/docker/for-mac/issues/90)

 * `docker-diagnose`: display and record the time the diagnosis was captured

 * Don't compute the container folder in `com.docker.vmnetd`
   Fixes [docker/for-mac#47](https://github.com/docker/for-mac/issues/47)

 * Warn the user if BlueStacks is installed (potential kernel panic)

 * Automatic update interval changed from 1 hour to 24 hours

 * Include Zsh completions

 * UI Fixes

**Networking**

* VPNKit supports search domains

* slirp: support up to 8 external DNS servers

* slirp: reduce the number of sockets used by UDP NAT, reduce the probability that NAT rules will time out earlier than expected

* Entries from `/etc/hosts` should now resolve from within containers

* Allow ports to be bound on host addresses other than `0.0.0.0` and `127.0.0.1`
  Fixes issue reported in
  [docker/for-mac#68](https://github.com/docker/for-mac/issues/68)

* Use Mac System Configuration database to detect DNS

**File sharing (osxfs)**

* Fixed thread leak

* Fixed a malfunction of new directories that have the same name as an old directory that is still open

* Rename events now trigger DELETE and/or MODIFY `inotify` events
(saving with TextEdit works now)

* Fixed an issue that caused `inotify` failure and crashes

* Fixed a directory file descriptor leak

* Fixed socket `chowns`

**Moby**

* Use default `sysfs` settings, transparent huge pages disabled

* `cgroup` mount to support `systemd` in containers

* Increase Moby `fs.file-max` to 524288

* Fixed Moby Diagnostics and Update Kernel

**HyperKit**

* HyperKit updated with `dtrace` support and lock fixes

### Docker for Mac 2016-08-11 1.12.0-a

[Download](https://download.docker.com/mac/stable/11213/Docker.dmg)

This bug fix release contains osxfs improvements. The fixed issues may have
been seen as failures with apt-get and npm in containers, missed inotify
events or unexpected unmounts.

**Bug fixes**

* osxfs: fixed an issue causing access to children of renamed directories to fail (symptoms: npm failures, apt-get failures)

* osxfs: fixed an issue causing some ATTRIB and CREATE inotify events to fail delivery and other inotify events to stop

* osxfs: fixed an issue causing all inotify events to stop when an ancestor directory of a mounted directory was mounted

* osxfs: fixed an issue causing volumes mounted under other mounts to spontaneously unmount


### Docker for Mac 1.12.0, 2016-07-28

[Download](https://download.docker.com/mac/stable/10871/Docker.dmg)

* First stable release

**Components**

* Docker 1.12.0
* Docker Machine 0.8.0
* Docker Compose 1.8.0
