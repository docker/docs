---
description: Change log / release notes per edge release
keywords: Docker Desktop for Windows, edge, release notes
title: Docker Desktop for Windows Edge Release notes
toc_min: 1
toc_max: 2
---

This page contains information about Docker Desktop Edge releases. Edge releases give you early access to our newest features. Note that some of the features may be experimental, and some of them may not ever reach the Stable release.

For information about Stable releases, see the [Stable release notes](release-notes.md). For Docker Desktop system requirements, see
[What to know before you install](install.md#what-to-know-before-you-install).

## Docker Desktop Community 2.3.4.0
2020-07-28

> [Download](https://desktop.docker.com/win/edge/46980/Docker%20Desktop%20Installer.exe)

### New

- Docker Desktop introduces the new **Images** view to the Docker Dashboard. The Images view allows users to view a list of Docker images on the disk, run an image as a container, pull the latest version of an image from Docker Hub, inspect images, and remove any unwanted images from the disk.

  To access the new Images view, from the Docker menu, select **Dashboard** > **Images**.

### Upgrades

- [Docker ECS integration v1.0.0-beta.4](https://github.com/docker/ecs-plugin/releases/tag/v1.0.0-beta.4){: target="_blank" class="_”}
- [Kubernetes 1.18.6](https://github.com/kubernetes/kubernetes/releases/tag/v1.18.6){: target="_blank" class="_”}

### Bug fixes and minor changes

- UI changes to dark or light theme depending on Windows default app mode.
- Copying the container logs from the dashboard does not copy the ANSI color codes to the clipboard anymore.

## Docker Desktop Community 2.3.3.2
2020-07-21

> [Download](https://desktop.docker.com/win/edge/46784/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker ECS integration v1.0.0-beta.2](https://github.com/docker/ecs-plugin/releases/tag/v1.0.0-beta.2){: target="_blank" class="_”}
- [Docker ACI integration 0.1.10](https://github.com/docker/aci-integration-beta/releases/tag/v0.1.10){: target="_blank" class="_”}

## Docker Desktop Community 2.3.3.1
2020-07-10

> [Download](https://desktop.docker.com/win/edge/46608/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

- Fixed a bug where the ECS plugin was not visible from WSL 2.

## Docker Desktop Community 2.3.3.0
2020-07-09

> [Download](https://desktop.docker.com/win/edge/46574/Docker%20Desktop%20Installer.exe)

### Upgrades

- Beta release of [Docker ECS integration v1.0.0-beta.1.](https://docs.docker.com/engine/context/ecs-integration/)
- [Docker ACI integration v0.1.7](https://github.com/docker/aci-integration-beta/releases/tag/v0.1.7)
- [Docker Compose 1.26.2](https://github.com/docker/compose/releases/tag/1.26.2)

### Bug fixes and minor changes

- Compose-on-Kubernetes is no longer included in the Docker Desktop installer. You can download it separately from the compose-on-kubernetes [release page](https://github.com/docker/compose-on-kubernetes/releases).

## Docker Desktop Community 2.3.2.1
2020-06-29

> [Download](https://desktop.docker.com/win/edge/46329/Docker%20Desktop%20Installer.exe)

### Security

- Fixed a regression in Docker Desktop Edge 2.3.2.0 on Windows where the Docker API was exposed by default on a random port on all interfaces.

## Docker Desktop Community 2.3.2.0
2020-06-25

### Upgrades

- [Docker 19.03.12](https://github.com/docker/docker-ce/releases/tag/v19.03.12)
- [Docker Compose 1.26.0](https://github.com/docker/compose/releases/tag/1.26.0)
- [Kubernetes 1.18.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.18.3)
- Beta release of the [Docker ACI integration](https://docs.docker.com/engine/context/aci-integration/)

### Bug fixes and minor changes

- Fixed application startup if `hosts` is specified inside the Docker `daemon.json`. See [docker/for-win#6895](https://github.com/docker/for-win/issues/6895#issuecomment-637429117).
- Fixed opening a CLI from the dashboard for Windows containers. See [docker/for-win#7079](https://github.com/docker/for-win/issues/7079).
- Fixed an issue where adding a folder on a non-existing drive in the settings would create an empty entry. See [docker/for-win#6797](docker/for-win#6797).
- Fixed an issue when updating Windows to a WSL 2 capable version, where a configuration file is locked by the old Hyper-V VM.
- Pinning the application to the taskbar and clicking on it will launch the container view if Docker is already running.
- Left-clicking on the Docker icon in the system tray now launches the container view Dashboard. See [docker/for-win#6650](https://github.com/docker/for-win/issues/6650).
- Fixed an issue which blocks Docker Desktop startup if `localhost:2375` is used by another program. See [docker/for-win#6929](https://github.com/docker/for-win/issues/6929) and [docker/for-win#6961](https://github.com/docker/for-win/issues/6961).
- Docker contexts are now synced between Windows and WSL 2 distributions.
- Docker Desktop now prompts users to share drives for paths like `////c/Users/foo`, rather than only paths like `C:\Users\foo` and `C:/Users/foo`.
- The installer now returns a non-zero exit code on error or when the user cancels the installation.
- Avoid failing withFixed an issue which caused Docker Desktop to fail with `Function not implemented` error during file I/O on shared volumes. See [docker/for-win#5955](https://github.com/docker/for-win/issues/5955).
- Fixed an issue when bind mounting a mounted host directory inside a container. See [docker/for-win#5089](https://github.com/docker/for-win/issues/5089).
- Removed the legacy Kubernetes context `docker-for-desktop`. The context `docker-desktop` should be used instead. See [docker/for-win#5089](https://github.com/docker/for-win/issues/5089).

## Docker Desktop Community 2.3.1.0
2020-05-20

> [Download](https://desktop.docker.com/win/edge/45408/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker Compose 1.26.0-rc4](https://github.com/docker/compose/releases/tag/1.26.0-rc4)
- Upgrade to Qemu 4.2.0, add Risc-V support

### Bug fixes and minor changes

**Hyper-V**

- Create drive symlinks in VM at startup to avoid breaking setup that was using them. Fixes [docker/for-win#6628](https://github.com/docker/for-win/issues/6628).
- Implement `fallocate` for shared filesystems. See [docker/for-win#6658](https://github.com/docker/for-win/issues/6658#issuecomment-627736820).

**WSL 2**

- Configure CLI to use Docker Desktop credential store.
- Adds a restart button to the popup prompting the user to install the Linux kernel.
- More reliable bootstrap, do not rely on `wslpath` for path translation, retry mounts on errors at startup.

**Dashboard**

- Fixed containers logs which were sometimes truncated. Fixes [docker/for-win#5954](https://github.com/docker/for-win/issues/5954)
- Fixed `open with vs code` button for compose app deployed from a WSL 2 Linux Workspace.

**Other fixes**

- Fixed an installer crash when an old and/or partially uninstalled version of Docker Desktop was present on the system. [Fixes docker/for-win/6536](https://github.com/docker/for-win/issues/6536).
- Fixed home expansion when using tilde in bind mounts source (for example, `-v ~/dir:/vm-dir`)
- `localhost` and `127.0.0.1` can both be used in the proxy settings to redirect to a proxy on the host. Fixes [docker/for-win#5715](https://github.com/docker/for-win/issues/5715).
- Fixed a typo in the backend destroy notification. Fixes [docker/for-win#6739](https://github.com/docker/for-win/issues/6739).
- Fixed a crash which sometimes occurred when Docker Desktop loads a corrupted Docker CLI configuration file. Fixes [docker/for-win#6657](https://github.com/docker/for-win/issues/6657).
- Fixed a delay when opening the systray menu. Fixes [docker/for-win#1011](https://github.com/docker/for-win/issues/1011).

### Known issues

**WSL 2**

- Swarm service bind mounts are not always restored correctly.
- Bind mounts of files living outside of the root mount point (files within `/mnt/c`, `/tmp`, `/run...`) don't work correctly when multiple containers mount them.

## Docker Desktop Community 2.3.0.1
2020-04-28

> [Download](https://download.docker.com/win/edge/44875/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

- Fixed an upgrade bug where users on versions still using the PowerShell-based VM management experience a silent un-installation crash leading to Docker Desktop being uninstalled instead of upgraded.
- Fixed a crash when failing to login due to no Internet connection.
- Fixed bug in handling shared volume paths with `..` characters. Fixes [docker/for-win#5375](https://github.com/docker/for-win/issues/5375).
- WSL 2: Docker Desktop detects when the wsl-distro is stopped and displays an actionable error message.
- Fixed a bug in `ftruncate` which prevented [libgit2sharp](https://github.com/libgit2/libgit2sharp) from cloning repositories on shared volumes. See [docker/for-win#5808](https://github.com/docker/for-win/issues/5808#issuecomment-610996272).
- Fixed a bug that caused starting and stopping of a Compose application from the UI to fail when the path contains whitespace.

## Docker Desktop Community 2.3.0.0
2020-04-20

> [Download](https://download.docker.com/win/edge/44472/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker Compose 1.25.5](https://github.com/docker/compose/releases/tag/1.25.5)
- [Go 1.13.10](https://github.com/golang/go/issues?q=milestone%3AGo1.13.10+label%3ACherryPickApproved)
- [Linux kernel 4.19.76](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.76-ce15f646db9b062dc947cfc0c1deab019fa63f96-amd64/images/sha256-6c252199aee548e4bdc8457e0a068e7d8e81c2649d4c1e26e4150daa253a85d8?context=repo)
- LinuxKit [init](https://hub.docker.com/layers/linuxkit/init/1a80a9907b35b9a808e7868ffb7b0da29ee64a95/images/sha256-64cc8fa50d63940dbaa9979a13c362c89ecb4439bcb3ab22c40d300b9c0b597e?context=explore), [runc](https://hub.docker.com/layers/linuxkit/runc/69b4a35eaa22eba4990ee52cccc8f48f6c08ed03/images/sha256-57e3c7cbd96790990cf87d7b0f30f459ea0b6f9768b03b32a89b832b73546280?context=explore), and [containerd](https://hub.docker.com/layers/linuxkit/containerd/09553963ed9da626c25cf8acdf6d62ec37645412/images/sha256-866be7edb0598430709f88d0e1c6ed7bfd4a397b5ed220e1f793ee9067255ff1?context=explore)

### Bug fixes and minor changes

> Docker Desktop Edge 2.3.0.0 fixes 10 issues reported on the [docker/for-win](https://github.com/docker/for-win/issues) GitHub repository.

**WSL 2**

- Detect when the WSL 2 backend stops and allow the user to restart it.
- Added support for chmod/chown on WSL 2 bind mounts. Fixed [docker/for-win#6284](https://github.com/docker/for-win/issues/6284).
- Fixed a race condition when exposing ports.
- Prevent the WSL 2 dialog from blocking other windows.
- Added a check that virtualization is enabled in the BIOS.

**Filesharing**

- Fixed a bug when parent directories of shared folders are renamed, leading to spurious `File not found` errors. Fixes [docker/for-win#6200](https://github.com/docker/for-win/issues/6200).
- Fixed bug where a letter and two digit root folder name causes Docker Compose to fail creating directories inside the volume. Fixes [docker/for-win#6248](https://github.com/docker/for-win/issues/6248).
- Fixed a bug which prevented containers from seeing file updates on shared volumes, when under load and a container restarts. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530#issuecomment-608804192)
- Fixed a bug where host paths are incorrectly translated into VM paths. Fixes [docker/for-win#6209](https://github.com/docker/for-win/issues/6209).
- Fixed bug receiving file events on long paths (> 260 characters). Fixes [docker/for-win#6337](https://github.com/docker/for-win/issues/6337).

**Other fixes**

- Fixed a bug where containers disappeared from the UI when Kubernetes context is invalid. Fixes [docker/for-win#6037](https://github.com/docker/for-win/issues/6037).
- Fixed Windows event logs filtering when copying them to the Docker Desktop log files. Fixes [docker/for-win#6258](https://github.com/docker/for-win/issues/6258).
- Fixed a handle leak in `vpnkit-bridge`. Fixed [docker/for-win#5841](https://github.com/docker/for-win/issues/5841)
- Fixed a bug when removing Docker Desktop virtual switch.
- Added a link to the Stable channel from the Docker Desktop UI.
- IPv6 has been re-enabled in the embedded Linux kernel, so listening on IPv6 addresses works again. Fixed [docker/for-win#6206](https://github.com/docker/for-win/issues/6206) and [docker/for-mac#4415](https://github.com/docker/for-mac/issues/4415).
- Made the embedded terminal resizable.
- Fixed a bug which only respected the "expose on TCP" Docker engine API setting on application restart. Setting will now take effect when the Apply button is clicked.
- Fixed a bug where diagnostic upload would fail if the username contained spaces.

## Docker Desktop Community 2.2.3.0
2020-04-02

> [Download](https://download.docker.com/win/edge/43965/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker 19.03.8](https://github.com/docker/docker-ce/releases/tag/v19.03.8)
- [Docker Compose 1.26.0-rc3](https://github.com/docker/compose/releases/tag/1.26.0-rc3)
- [Linux 4.19.76](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.76-4e5d9e5f3bde0abf236f97e4a81b029ae0f5f6e7-amd64/images/sha256-11dc0f6ee3187088219ba1463ebb378f5093a7d98f176ddfd62dd6b741c2dd2d?context=repo)

### New

- Docker Desktop introduces a new onboarding tutorial upon first startup. The Quick Start tutorial guides users to get started with Docker in a few easy steps. It includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

### Bug fixes and minor changes

> Docker Desktop Edge 2.2.3.0 fixes 28 issues reported on the [docker/for-win](https://github.com/docker/for-win/issues) GitHub repository.

**WSL 2**
- Docker Desktop only exposes host ports in Linux if the Windows port is available.
- Docker Desktop now allows users to refresh the list of Linux distros.
- Docker Desktop defaults to WSL 2 on install on compatible OS versions.
- Docker Desktop detects missing Linux kernel and adds a pointer to the Microsoft documentation to download the kernel.

**File sharing**
- Kubernetes: Persistent volumes created by claims are now stored in the virtual machine. Fixes [for-win/issues/5665](https://github.com/docker/for-win/issues/5665).
- Docker Desktop ensures that host paths accessed by containers are within the shared folders list.
- Fixed a bug where opening a read-only file would fail with an `Operation not permitted error`. Fixes [docker/for-win#6016](https://github.com/docker/for-win/issues/6016) and [docker/for-win#6017](https://github.com/docker/for-win/issues/6017).
- Fixed path handling for `docker volume create -o type=none -o o=bind -o device=C:\Some\Windows\path`.
- Fixed a bug which prevented open files from being deleted. Fixes [docker/for-win#5565](https://github.com/docker/for-win/issues/5565).
- Docker Desktop now avoids locking files on the host which are open in containers.
- Docker Desktop generates `fsnotify.WRITE` events in Linux containers when files are changed on the host. [docker/for-win#5530](https://github.com/docker/for-win/issues/5530#issuecomment-585572414)
- Docker Desktop now displays hidden files on shared volumes. Fixes [docker/for-win#5808](https://github.com/docker/for-win/issues/5808).
- Fixed cache invalidation and event injection for shared volumes with long paths.
- Docker Desktop now handles case-insensitivity correctly during file creation.
- Docker Desktop represents valid directory junctions as directories (rather than symlinks) and handles cache invalidation and event injection properly. Fixes [docker/for-win#5582](https://github.com/docker/for-win/issues/5582).
- Fixed a race condition in `readlink` on shared volumes using "mfsymlinks". Fixes [docker/for-win#5793](https://github.com/docker/for-win/issues/5793)
- Fixed a bug where files did not change in shared volumes when using `volumes_from` in docker-compose.yml. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530).

**Other fixes**
- Reduced the size of the Docker Desktop installer from 960 MB to 400 MB.
- Added an option to delete container and image data from the **Troubleshoot** screen.
- Docker Desktop now reads the Hyper-V VM disk max size on startup and uses it as the value to display in the Settings.
- Fixed the Hyper-V used disk image size being reported incorrectly.
- Fixed the Hyper-V VM disk size increase. Fixes [docker/for-win#5881](https://github.com/docker/for-win/issues/5881).
- Fixed container time drift when system hibernates. Fixes [docker/for-win#4526](https://github.com/docker/for-win/issues/4526).
- Fixed a bug where the Docker Desktop UI could be started without the engine. Fixes [docker/for-win#5376](https://github.com/docker/for-win/issues/5376).
- Docker Desktop now uses the least possible privilege for querying the `Server service`. Fixes [docker/for-win#5150](https://github.com/docker/for-win/issues/5150).
- Fixed a bug where diagnostics upload could fail silently.
- Capturing diagnostics is now faster and easier.
- Fixed an issue where a container port could not be exposed on a specific host IP. See [docker/for-win#5546](https://github.com/docker/for-mac/issues/5546).
- Removed port probing from dashboard, just unconditionally showing links to ports that should be available. Fixes [docker/for-win#5903](https://github.com/docker/for-win/issues/5903).

### Known issues

- Loopback and unspecified IPv6 addresses (`::` and `::1`) within a container do not currently work. Some web servers and other programs may be using these addresses in their configuration files.

## Docker Desktop Community 2.2.2.0
2020-03-02

> [Download](https://download.docker.com/win/edge/43066/Docker%20Desktop%20Installer.exe)

This release contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installing Docker Desktop.

### Upgrades

- [Kubernetes 1.16.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.16.5)
- [Go 1.13.8](https://golang.org/doc/devel/release.html#go1.13)

### New

- Windows 10 Home users can now use Docker Desktop through the experimental WSL 2 support. This requires Windows Insider Preview Build 19018 or later.
- Docker Desktop now allows sharing individual folders, rather than whole drives, giving more control to users over what is being shared.

### Bug fixes and minor changes

- Added a prerequisite check in the installer for the `LanmanServer` service. Fixes [docker/for-win#5150](https://github.com/docker/for-win/issues/5150).
- DockerNAT has been removed from Docker Desktop as using an IP address to communicate from the host to a container is not a supported feature. To communicate from a container to the host, you must use the special DNS name `host.docker.internal`. We also recommend using ports to communicate from the host to a container. For more information, see [Networking](/docker-for-win/networking/#use-cases-and-workarounds).
- Ceph support has been removed from Docker Desktop to save disk space.
- Fixed an issue that caused the Windows log file archive to grow indefinitely. Fixes [docker/for-win#5113](https://github.com/docker/for-win/issues/5113).

## Docker Desktop Community 2.2.1.0
2020-02-12

> [Download](https://download.docker.com/win/edge/42746/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker Compose 1.25.4](https://github.com/docker/compose/releases/tag/1.25.4)
- [Go 1.12.16](https://golang.org/doc/devel/release.html#go1.12)

### Bug fixes and minor changes

- Fixed an issue that prevented users from creating files with special characters in the filenames within a shared volume. Fixes [docker/for-win#5520](https://github.com/docker/for-win/issues/5520).
- Fixed handling of shared volumes with relative paths in `docker-compose.yml`. Fixes [docker/for-win#5516](https://github.com/docker/for-win/issues/5516).
- Fixed handling of shared volumes where the path case (uppercase / lowercase) did not exactly match the host. Fixes [docker/for-win#5516](https://github.com/docker/for-win/issues/5516).
- Fixed an issue where changing a file in the Windows file system did not update the file within the container. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530), [docker/for-win#5550](https://github.com/docker/for-win/issues/5550), and [docker/for-win#5318](https://github.com/docker/for-win/issues/5318).
- Fixed an issue that prevented users from sharing drives and sometimes incorrectly prompted users to enter file system credentials. Fixes [docker/for-win#5567](https://github.com/docker/for-win/issues/5567).
- Fixed an issue which prevented users from mounting nested volumes. Fixes [docker/for-win#5540](https://github.com/docker/for-win/issues/5540).
- Fixed a file synchronization issue on bind mounts. Fixes [docker/for-win#5533](https://github.com/docker/for-win/issues/5533).
- Fixed an issue where timestamps were reset to zero on files within shared volumes. Fixes [docker/for-win#5528](https://github.com/docker/for-win/issues/5528) and [docker/for-win#5543](https://github.com/docker/for-win/issues/5543).
- Fixed a bug which caused Docker Desktop to fail when sharing file paths longer than 260 characters. Fixes [docker/for-win#5572](https://github.com/docker/for-win/issues/5572).
- Fixed an issue where some users were unable to start a container if there is a symlink in a shared volume. Fixes [docker/for-win#5582](https://github.com/docker/for-win/issues/5582).
- Fixed a bug where users were unable to modify the **Manual proxy configuration** settings through the Docker Desktop UI. Fixes [docker/for-win#5606](https://github.com/docker/for-win/issues/5606) and [docker/for-win#5548](https://github.com/docker/for-win/issues/5548).
- Fixed an issue where the Docker Desktop UI failed to start if the `no_proxy` environment variable has an entry that starts with a dot (.). Fixes [docker/for-win/5551](https://github.com/docker/for-win/issues/5551).
- Fixed a bug that did not allow users to access the Docker Desktop UI on machines with a policy preventing access to the Windows registry. Fixes [docker/for-win#5536](https://github.com/docker/for-win/issues/5536).
- Fixed an issue where disabling the WSL2 integration and re-enabling it later resulted in broken symlinks. Fixes [docker/for-win#5613](https://github.com/docker/for-win/issues/5613).
- Fixed a race condition when starting the WSL engine which caused Docker Desktop to incorrectly report that the containers have exited. Fixes [docker/for-win#5607](https://github.com/docker/for-win/issues/5607).
- Fixed an issue where editing code inside a container resulted in an error. Fixes [docker/for-win#5528](https://github.com/docker/for-win/issues/5528).
- Fixed a bug where running the command `DockerCli.exe -SharedDrives` failed to display a list of drives that are shared. Fixes [docker/for-win#5625](https://github.com/docker/for-win/issues/5625).

### Known issues

- DockerNAT has been removed from Docker Desktop 2.2.0.0 as using an IP address to communicate from the host to a container is not a supported feature. To communicate from a container to the host, you must use the special DNS name `host.docker.internal`. We also recommend using ports to communicate from the host to a container. For more information, see [Networking](/docker-for-win/networking/#use-cases-and-workarounds).

  However, if your current setup relies on IP addresses for communication, you can use a temporary workaround to reinstate DockerNAT. To do this, open
`C:\Program Files\Docker\Docker\resources\MobyLinux.ps1` and add `$SwitchName = "DockerNAT"` between line 175 and 176. Note that the temporary workaround to reinstate DockerNAT may be removed from future releases.
- Directory junctions in shared volumes pointing outside the shared volume do not work.
- It is currently not possible to rename open files within shared volumes.
- In some cases, Docker containers do not synchronize time with the host when the host machine is put to sleep.
- It is not possible to use `localhost` in the proxy settings. You can instead use `host.docker.internal` to work around this issue.

## Docker Desktop Community 2.1.7.0
2019-12-11

[Download](https://download.docker.com/win/edge/41561/Docker%20Desktop%20Installer.exe)

> **Note:** Docker Desktop Edge 2.1.7.0 is the release candidate for the upcoming major Stable release. Please help us test this version before the wider release and report any issues in the [docker/for-win](https://github.com/docker/for-win/issues) GitHub repository.

### Upgrades

- [Docker Compose 1.25.1-rc1](https://github.com/docker/compose/releases/tag/1.25.1-rc1)

### New

- Added the **WSL Integration** option in **Settings** > **Resources** to select WSL 2 distributions. This option will be visible only after enabling WSL 2 on your machine.
- File sharing: The new gRPC FUSE system is no longer experimental and is enabled by default even when experimental features are turned off.

### Bug fixes and minor changes

- The Docker Desktop Dashboard now displays port information inline with the container status.
- Fixed an issue which caused the 'back' button on the Dashboard UI to behave inconsistently when repeatedly switching between container details and the Settings window.
- Various minor improvements to the Dashboard UI.
- Fixed a bug that caused a rare crash when uninstalling Docker Desktop.
- Fixed an issue when overwriting existing files with the WSL 2 CLI bundle, such as kubectl.
- Fixed an issue where WSL 2 failed to start when the distribution directory is compressed or encrypted. Fixes [docker/for-win#5256](https://github.com/docker/for-win/issues/5256).
- Fixed an issue that caused Docker Desktop to fail on startup when there is an incomplete Kubernetes config file.
- Fixed various stability issues.

### Known issues

- The Dashboard stops updating when you switch the container mode between Linux and Windows. To work around this issue, close and reopen the Dashboard.
- The new gRPC FUSE file sharing implementation does not support connecting to new drives after Docker Desktop starts (for example, a USB drive). If you would like to share the new drive in Docker Desktop, you must quit Docker Desktop and then start the application to refresh the list of drives in Settings.

## Docker Desktop Community 2.1.6.1
2019-11-20

[Download](https://download.docker.com/win/edge/40920/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

- Fixed an issue that prevented Kubernetes to start with WSL 2 on machines with multiple CPU cores.
- Fixed a rare issue that caused to Docker Desktop to crash with the error `Unable to stop Hyper-V VM: Cannot validate argument on parameter 'SwitchName'. The argument is null or empty.`

### Known issue

Windows Insider Preview Slow Ring users running OS builds older than 19025 cannot run WSL 2 following upgrade to Docker Desktop Edge 2.1.6.1. WSL 2 requires Windows 10 Insider Preview build 19018 or greater.

### Docker Desktop Community 2.1.6.0
2019-11-18

[Download](https://download.docker.com/win/edge/40807/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker 19.03.5](https://github.com/docker/docker-ce/releases/tag/v19.03.5)
- [Go 1.12.13](https://golang.org/doc/devel/release.html#go1.12)

### New

Added the ability to start and stop Compose-based applications and view combined logs in the Docker Desktop **Dashboard** UI.

### Bug fixes and minor changes

- Docker Desktop now automatically restarts after an update.
- Fixed an issue where Docker Desktop auto-start was not being disabled properly on some machines.
- Fixed a container start error when a container has more than one port with an arbitrary or not-yet-configured external port number. For example, `docker run -p 80 -p 443 nginx`). Fixes [docker/for-win#4935](https://github.com/docker/for-win/issues/4935) and [docker/compose#6998](https://github.com/docker/compose/issues/6998).
- Fixed an issue which caused Docker Desktop to crash when resetting to factory defaults while running Windows containers.
- Fixed multiple issues related to Fast Startup.
- Injected Docker CLI, CLI plugins, Docker Compose, Notary, and kubectl into WSL distros when Docker Desktop WSL integration is enabled.
- Fixed an issue where bind mounts created with Docker Compose from a WSL distro were incorrectly translated. Fixes [docker/for-win#5084](https://github.com/docker/for-win/issues/5084).
- Docker Desktop now supports inotify events on shared filesystems for Windows file sharing.
- Fixed a cache invalidation bug when a file in a shared volume is renamed on the host for Windows file sharing.
- Fixed a handle leak when calling `Mknod` on a shared volume for Windows file sharing.
- To make VM startup more reliable, Docker Desktop now avoids adding a Hyper-V NIC to the Windows VM when using Hypervisor sockets for Windows file sharing (rather than Samba).

### Known issue

Windows Insider Preview Slow Ring users running OS builds older than 19025 cannot run WSL 2 following upgrade to Docker Desktop Edge 2.1.6.0. WSL 2 requires Windows 10 Insider Preview build 19018 or greater.

## Docker Desktop Community 2.1.5.0
2019-11-04

[Download](https://download.docker.com/win/edge/40323/Docker%20Desktop%20Installer.exe)

This release contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installation.

### Upgrades

- [Kubernetes 1.15.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.15.5)
- [Docker Compose 1.25.0-rc4](https://github.com/docker/compose/releases/tag/1.25.0-rc4)
- Linux kernel 4.19.76

### New

- **Docker Desktop Dashboard:** The new Docker Desktop **Dashboard** provides a user-friendly interface which enables you to interact with containers and applications, and manage the lifecycle of your applications directly from the UI. In addition, it allows you to access the logs, view container details, and monitor resource utilization to explore the container behavior.

  To access the Dashboard UI, select the Docker menu from the system tray and then click **Dashboard**.

- **WSL 2 backend:** The new Docker Desktop WSL 2 backend replaces the Docker Desktop WSL 2 Tech Preview. The WSL 2 backend architecture introduces support for Kubernetes, provides an updated Docker daemon, offers VPN-friendly networking, and additional features. For more information, see [Docker Desktop WSL 2 backend](https://docs.docker.com/docker-for-windows/wsl/).

- **New file sharing implementation:** Docker Desktop introduces a new file sharing implementation which uses gRPC, FUSE, and Hypervisor sockets instead of Samba, CIFS, and Hyper-V networking. The new implementation  offers improved I/O performance. Additionally, when using the new file system:

  - Users don't have to expose the Samba port, and therefore do not experience issues related to IT firewall or the drive-sharing policy.
  - There is no need to provide user credentials to Docker Desktop. File access rights are automatically enforced when accessing mounted folders through containers.

  > To switch back to the Samba file system:
  >
  > From the Docker menu, select **Settings** > **Command line**, and then clear the **Enable experimental features** check box. Alternatively, install a Docker Desktop Stable release.

### Known issues

- When you start a Docker Compose application and then start a Docker App which has the same name as the Compose application, Docker Desktop displays only one application on the Dashboard. However, when you expand the application, containers that belong to both applications are displayed on the Dashboard.

- When you deploy a Docker App with multiple containers on Kubernetes, Docker Desktop displays each Pod as an application on the Dashboard.

## Docker Desktop Community 2.1.4.0
2019-10-15

[Download](https://download.docker.com/win/edge/39357/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker 19.03.3](https://github.com/docker/docker-ce/releases/tag/v19.03.3)
- [Kubernetes 1.15.4](https://github.com/kubernetes/kubernetes/releases/tag/v1.15.4)
- [Go 1.12.10](https://github.com/golang/go/issues?q=milestone%3AGo1.12.10+label%3ACherryPickApproved) for [CVE-2019-16276](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-16276)
- [Kitematic 0.17.9](https://github.com/docker/kitematic/releases/tag/v0.17.9)

### Bug fixes and minor changes

- Improved the navigation in **Settings** and **Troubleshoot** UI.
- Fixed a bug that prevented users from accessing WSL 2 Tech Preview. Fixes [docker/for-win#4734](https://github.com/docker/for-win/issues/4734).
- Docker Machine is no longer included in the Docker Desktop installer. You can download it separately from the [Docker Machine releases](https://github.com/docker/machine/releases) page.

## Docker Desktop Community 2.1.3.0
2019-09-16

[Download](https://download.docker.com/win/edge/38275/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

- Added a loading overlay to the **Settings** and **Troubleshoot** windows to prevent editing conflicts.
- Deactivated the **Reset Kubernetes** button when Kubernetes is not activated.

## Docker Desktop Community 2.1.2.0
2019-09-09

[Download](https://download.docker.com/win/edge/38030/Docker%20Desktop%20Installer.exe)

#### Upgrades

- [Docker 19.03.2](https://github.com/docker/docker-ce/releases/tag/v19.03.2)
- [Kubernetes 1.14.6](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.6)
- [Go 1.12.9](https://github.com/golang/go/issues?q=milestone%3AGo1.12.9+label%3ACherryPickApproved)
- [Qemu 4.0.1](https://github.com/docker/binfmt)
- [Docker Machine 0.16.2](https://github.com/docker/machine/releases/tag/v0.16.2)
- [Kitematic 0.17.8](https://github.com/docker/kitematic/releases/tag/v0.17.8)

#### Bug fixes and minor changes

- Reduced the Virtual Machine (VM) startup time.
- Added support for `Expect: 100-continue` headers in the Docker API proxy. Some HTTP clients such as `curl` send this header when the payload is large, for example, when creating containers. Fixes [moby/moby#39693](https://github.com/moby/moby/issues/39693).

## Docker Desktop Community 2.1.1.0
2019-08-12

[Download](https://download.docker.com/win/edge/37260/Docker%20Desktop%20Installer.exe)

#### Upgrades

- Linux Kernel 4.14.131

#### New

- [Docker Desktop WSL 2 Tech Preview](https://docs.docker.com/docker-for-windows/wsl/)

#### Bug fixes and minor changes

- Fixed a PowerShell script signing issue that caused `AuthorizationManager check failed` errors on machines with strict group policies on PowerShell script signing. [docker/for-win#4376](https://github.com/docker/for-win/issues/4376)
- Fixed an issue where attempts to upgrade Docker Desktop to version 2.1.0.0 sometimes failed with the error value cannot be null. [docker/for-win#3375](https://github.com/docker/for-win/issues/3375), [docker/for-win#4381](https://github.com/docker/for-win/issues/4381), [docker/for-win#4343](https://github.com/docker/for-win/issues/4343), [docker/for-win#4359](https://github.com/docker/for-win/issues/4359)
- Fixed an issue to ensure Docker Desktop Community installer does not prompt for UAC when UAC is disabled.
- Improved the error messages displayed during VM lifecycle operations. [docker/for-win#4348](https://github.com/docker/for-win/issues/4348)
- Docker Desktop now supports a configurable user timeout for VMs on slower machines. [docker/for-win#4393](https://github.com/docker/for-win/issues/4393)
- Enabled Windows features such as Hyper-V and Containers during installation, thereby reducing the need for another restart after installation.

## Docker Desktop Community 2.1.0.0 
2019-07-30

[Download](https://download.docker.com/win/edge/36873/Docker%20Desktop%20Installer.exe)

This release contains Kubernetes security improvements. Note that your local Kubernetes PKI and cluster will be reset after installation.

#### Upgrades

 - [Docker 19.03.1](https://github.com/docker/docker-ce/releases/tag/v19.03.1)
 - [Docker Compose 1.24.1](https://github.com/docker/compose/releases/tag/1.24.1)
 - [Alpine 3.10](https://alpinelinux.org/posts/Alpine-3.10.0-released.html)
 - Linux Kernel 4.9.184
 - [Docker Credential Helpers 0.6.3](https://github.com/docker/docker-credential-helpers/releases/tag/v0.6.3)

#### New

 - Introduced a new user interface for the Docker Desktop **Settings** menu.
 - The **Restart** and **Reset** options are now available on the **Troubleshoot** menu. 

#### Bug fixes and minor changes

 - Changed the host's kubernetes context to ensure `docker run -v .kube:kube ... kubectl` works.
 - Restricted the `cluster-admin` role on local Kubernetes cluster to `kube-system` namespace.
 - Fixed Kubernetes installation with VPNkit subnet.
 - Fixed an issue where Docker Desktop restarts when a user logs out of Windows and logs back in, which results in retaining the    
   exported ports on containers.
 - Reduced the VM startup time. `swap` is not created every time a virtual machine boots.
 - Fixed a bug which caused Docker Desktop to crash when a user cancels switching the version using Windows User Account Control (UAC)    settings.
 - Fixed a bug where the process output was not redirected to stdout when gathering diagnostics on Windows, which sometimes resulted in    a crash.


## Docker Community Edition 2.0.5.0 2019-06-12

[Download](https://download.docker.com/win/edge/35318/Docker%20Desktop%20Installer.exe)

This is the Edge channel, which gives you early access to our newest features. Be aware that some of them may be experimental, and some of them may not ever reach the Stable release.

This release contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installation.

* Upgrades
  - [Docker 19.03.0-rc2](https://github.com/docker/docker-ce/releases/tag/v19.03.0-rc2)
  - [Kubernetes 1.14.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.3)
  - [Compose on Kubernetes 0.4.23](https://github.com/docker/compose-on-kubernetes/releases/tag/v0.4.23)
  - [linuxkit v0.7](https://github.com/linuxkit/linuxkit/releases/tag/v0.7)
  - [Qemu 4.0.0](https://github.com/docker/binfmt) for cross compiling for ARM

* New
  - Docker Desktop includes the `buildx` plugin (currently experimental).
  - Selecting the `Experimental features` checkbox on the Docker Desktop settings Daemon page enables experimental features in the Docker daemon and the Docker CLI.
  - Docker Desktop now checks for stored credentials at startup before attempting to mount any shared drives. This prompts users to reenter the credentials if they are invalid.

* Bug fixes and minor changes
  - Fixed race condition where Kubernetes sometimes fails to start after the app is restarted.
  - The system tray icon now opens the Docker Desktop menu with left or right mouse button.
  - When displaying the crash report window, Docker Desktop does not send a bugsnag crash report unless the user needs the report to upload diagnostics.
  - Docker Desktop has removed the ability to log in using email address as a username as the Docker command line does not support this.
  - For Linux containers on Windows (LCOW), at least one physical computer running Windows 10 Professional or Windows 10 Enterprise version 1809 or later is required.
  - The `Send usage statistics` checkbox is selected by default in Docker Desktop for Windows (Community). This option cannot be modified.
  - Docker Desktop has added a new dialog box during startup which allows users to retry mounting a shared drive or remove it from the shared drives list after a failed attempt.


* Known issues
  - Windows containers networking does not work properly on Windows 1903.

## Docker Community Edition 2.0.4.1 2019-05-07

[Download](https://download.docker.com/win/edge/34207/Docker%20Desktop%20Installer.exe)

* Bug fixes and minor changes
  - Upgrade QEMU from 2.8.0 to 3.1.0 to fix an emulation issue when building and running Java applications on Arm64 devices.

## Docker Community Edition 2.0.4.0 2019-04-30

[Download](https://download.docker.com/win/edge/33772/Docker%20Desktop%20Installer.exe)

* Upgrades
  - [Docker 19.03.0-beta3](https://github.com/docker/docker-ce/releases/tag/v19.03.0-beta3)
  - [Docker Compose 1.24.0](https://github.com/docker/compose/releases/tag/1.24.0)
  - [Compose on Kubernetes 0.4.22](https://github.com/docker/compose-on-kubernetes/releases/tag/v0.4.22)
  - [Kubernetes 1.14.1](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.14.md#changelog-since-v1141)

* New

  - App: Docker CLI plugin to configure, share, and install applications
  
    - Extend Compose files with metadata and parameters
    - Reuse the same application across multiple environments (Development/QA/Staging/Production)
    - Multi-orchestrator installation (Swarm or Kubernetes)
    - Push/Pull/Promotion/Signing supported for application, with the same workflow as images
    - Fully CNAB compliant
    - Full support for Docker Contexts
    
  - Buildx (Tech Preview): Docker CLI plugin for extended build capabilities with BuildKit
  
    - Familiar UI from docker build
    - Full BuildKit capabilities with container driver
    - Multiple builder instance support
    - Multi-node builds for cross-platform images (out-of-the-box support for linux/arm/v7 and linux/arm64)
    - Parallel building of compose files
    - High-level build constructs with `bake`

* Bug fixes and minor changes

  - Fix `Delete` of persistent volume claims
  - Truncate UDP DNS responses which are over 512 bytes in size

## Docker Community Edition 2.0.3.0 2019-03-05

[Download](https://download.docker.com/win/edge/31778/Docker%20Desktop%20Installer.exe)

* Upgrades
  - [Docker 18.09.3](https://github.com/docker/docker-ce/releases/tag/v18.09.3)

* Bug fixes and minor changes
  - Fixed docker not added to PATH after install in some cases
  - Fixed port 8080 that was used on localhost when starting Kubernetes.
  - Fixed "create issue" link in diagnostics windows.

## Docker Community Edition 2.0.2.1 2019-02-15

[Download](https://download.docker.com/win/edge/31274/Docker%20Desktop%20Installer.exe)

* Upgrades
  - [Docker 18.09.2](https://github.com/docker/docker-ce/releases/tag/v18.09.2), fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

## Docker Community Edition 2.0.2.0 2019-02-06

[Download](https://download.docker.com/win/edge/30972/Docker%20Desktop%20Installer.exe)

* Upgrades
  - [Docker Compose 1.24.0-rc1](https://github.com/docker/compose/releases/tag/1.24.0-rc1)
  - [Docker Machine 0.16.1](https://github.com/docker/machine/releases/tag/v0.16.1)
  - [Compose on Kubernetes 0.4.18](https://github.com/docker/compose-on-kubernetes/releases/tag/v0.4.18)

* New
  - Rebranded UI
  
* Bug fixes and minor changes
  - Kubernetes: use default maximum number of pods for kubelet. [docker/for-mac#3453](https://github.com/docker/for-mac/issues/3453)

## Docker Community Edition 2.0.1.0 2019-01-11

[Download](https://download.docker.com/win/edge/30090/Docker%20Desktop%20Installer.exe)

* Upgrades
  - [Docker 18.09.1](https://github.com/docker/docker-ce/releases/tag/v18.09.1)
  - [Kubernetes 1.13.0](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.13.md#v1130)
  - [Kitematic 0.17.6](https://github.com/docker/kitematic/releases/tag/v0.17.6)
  - Golang 1.10.6, fixes CVEs: [CVE-2018-16875](https://www.cvedetails.com/cve/CVE-2018-16875), [CVE-2018-16873](https://www.cvedetails.com/cve/CVE-2018-16873) and [CVE-2018-16874](https://www.cvedetails.com/cve/CVE-2018-16874)
  
  WARNING: If you have an existing Kubernetes cluster created with Docker Desktop, this upgrade will reset the cluster. If you need to back up your Kubernetes cluster or persistent volumes you can use [Ark](https://github.com/heptio/ark).

* Bug fixes and minor changes
  - Fix service log collection in diagnostics
  - Gather /etc/hosts to help diagnostics
  - Add 18.09 missing daemon options
  - Rename Docker for Windows to Docker Desktop
  - Partially open services ports if possibles
  - Quit will not check if service is running anymore
  - Fix UI lock when changing kubernetes state

## Edge Releases of 2018

### Docker Community Edition 2.0.0.0-win82 2018-12-07

[Download](https://download.docker.com/win/edge/29268/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker compose 1.23.2](https://github.com/docker/compose/releases/tag/1.23.2)

* Bug fixes and minor changes
  - Compose: Fixed a bug where build context URLs would fail to build on Windows. Fixes [docker/for-win#2918](https://github.com/docker/for-win/issues/2918) 

### Docker Community Edition 2.0.0.0-win77 2018-11-14

[Download](https://download.docker.com/win/edge/28777/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.09.0](https://github.com/docker/docker-ce-packaging/releases/tag/v18.09.0)
  - [Docker compose 1.23.1](https://github.com/docker/compose/releases/tag/1.23.1)
  - [Docker Machine 0.16.0](https://github.com/docker/machine/releases/tag/v0.16.0)
  - [Kitematic 0.17.5](https://github.com/docker/kitematic/releases/tag/v0.17.5)

* Bug fixes and minor changes
  - Windows Containers: Fix group daemon option settings. Fixes [docker/for-win#2647](https://github.com/docker/for-win/issues/2647) 
  - Windows Containers: Improve host.docker.internal ip resolution
  - Do not try to update samba share mounts when using Windows containers
  - Improved dns update too verbose in logs
  - VPNKit: Improved scalability of port forwarding.
  - VPNKit: Limit the size of the UDP NAT table. This ensures port forwarding and regular TCP traffic continue even when running very chatty UDP protocols.

### Docker Community Edition 2.0.0.0-beta1-win75 2018-09-14

[Download](https://download.docker.com/win/edge/19925/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.09.0-ce-beta1](https://github.com/docker/docker-ce/releases/tag/v18.09.0-ce-beta1)
  - Linux Kernel 4.9.125

* New
  - New version scheme

* Deprecation
  - Removed support of AUFS
  
* Bug fixes and minor changes
  - LCOW does not anymore need --platform flag on multi-arch images
  - Better WCOW host.docker.internal resolution on host, don't rewrite it if not modified. From [docker/for-win#1976](https://github.com/docker/for-win/issues/1976)
  - Disk size can now be ajusted from 32GiB to drive space.
  - Fix dns update too verbose in logs
  - Fix panic in diagnose

### Docker Community Edition 18.06.1-ce-win74 2018-08-29

[Download](https://download.docker.com/win/edge/19508/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.06.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.06.1-ce)

* Bug fixes and minor changes
  - Fix bug in virtual machine activity detection, preventing Docker Desktop from starting. Fixes [docker/for-win#2404](https://github.com/docker/for-win/issues/2404)
  - Fix bug in detection when Windows service is not running, and proposes to restart the service.
  - Fix local DNS failing to resolve inside containers. Fixes [docker/for-win#2301](https://github.com/docker/for-win/issues/2301), [docker/for-win#2304](https://github.com/docker/for-win/issues/2304)
  - Fix Kubernetes status display after reset to factory default
  - Fix bug where `host.docker.internal` is not resolved in some cases. Fixes [docker/for-win#2402](https://github.com/docker/for-win/issues/2402)
  - Use a 1MB vhdx blocksize instead of the default 32MB. See [docker/for-win#244](https://github.com/docker/for-win/issues/244). Also see [Microsoft Best Practices for running Linux on Hyper-V](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/best-practices-for-running-linux-on-hyper-v)
  - Fix diagnostics in specific cases when Windows service is not started.
  - Changed samba default file permissions back to avoid issue with too-open rights. Fixes [docker/for-win#2170](https://github.com/docker/for-win/issues/2170)
  - On RS5 insider, fix wrong detection of missing feature "Containers" requiring to install feature + reboot.

### Docker Community Edition 18.06.0-ce-win71 2018-07-27

[Download](https://download.docker.com/win/edge/19101/Docker%20for%20Windows%20Installer.exe)

* Bug fixes and minor changes
  - Fix bug in automatic feature enabling if "Hyper-V" and "Containers" Windows features are not enabled already when starting Docker Desktop.

### Docker Community Edition 18.06.0-ce-win69 2018-07-25

[Download](https://download.docker.com/win/edge/19070/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.06.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.06.0-ce)

### Docker Community Edition 18.06.0-ce-rc3-win68 2018-07-19

[Download](https://download.docker.com/win/edge/18994/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.06.0-ce-rc3](https://github.com/docker/docker-ce/releases/tag/v18.06.0-ce-rc3)
  - [Docker Machine 0.15.0](https://github.com/docker/machine/releases/tag/v0.15.0)
  - [Docker compose 1.22.0](https://github.com/docker/compose/releases/tag/1.22.0)

* New
  - New Kubernetes menu item allowing to switch Kubernetes context & connect to clusters other than the local one.

* Bug fixes and minor changes
  - AUFS storage driver is deprecated in Docker Desktop and AUFS support will be removed in the next major release. You can continue with AUFS in Docker Desktop 18.06.x, but you will need to reset disk image (in Settings > Reset menu) before updating to the next major update. You can check documentation to [save images](https://docs.docker.com/engine/reference/commandline/save/#examples) and [backup volumes](https://docs.docker.com/storage/volumes/#backup-restore-or-migrate-data-volumes)
  - Fix startup issue with AUFS
  - Fix status bug which could prevent the kubernetes cluster from starting.
  - Fix bug which would cause VM logs to be written to RAM rather than disk in some cases, and the VM to hang.
  - Fix security issue with named pipe connection to docker service.

### Docker Community Edition 18.05.0-ce-win67 2018-06-07

[Download](https://download.docker.com/win/edge/18263/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [LinuxKit v0.4](https://github.com/linuxkit/linuxkit/releases/tag/v0.4)
  - Linux Kernel 4.9.93 with CEPH, DRBD, RBD, MPLS_ROUTING and MPLS_IPTUNNEL enabled
  - [Kubernetes 1.10.3](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.10.md#v1103). If Kubernetes is enabled, the upgrade will be performed automatically when starting Docker for Windows.

* Bug fixes and minor changes
  - Fix VPNKit memory leak. Fixes [docker/for-win#2087](https://github.com/docker/for-win/issues/2087), [moby/vpnkit#371](https://github.com/moby/vpnkit/issues/371)
  - Fix potential loss of Kubernetes local cluster certificates on reboot.

### Docker Community Edition 18.05.0-ce-win66 2018-05-17

[Download](https://download.docker.com/win/edge/17760/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.05.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.05.0-ce)
  - [Docker compose 1.21.2](https://github.com/docker/compose/releases/tag/1.21.2)

* New 
  - Allow orchestrator selection from the UI in the "Kubernetes" pane, to allow "docker stack" commands to deploy to swarm clusters, even if Kubernetes is enabled in Docker for Windows.

* Bug fixes and minor changes
  - Fix restart issue when using Windows fast startup on latest 1709 Windows updates. Fixes [docker/for-win#1741](https://github.com/docker/for-win/issues/1741), [docker/for-win#1741](https://github.com/docker/for-win/issues/1741)
  - DNS name `host.docker.internal` can be used for host resolution from Windows containers.  Fixes [docker/for-win#1976](https://github.com/docker/for-win/issues/1976)
  - Fix broken link in diagnostics window. 

### Docker Community Edition 18.05.0-ce-rc1-win63 2018-04-26

[Download](https://download.docker.com/win/edge/17439/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.05.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v18.05.0-ce-rc1)
  - [Notary 0.6.1](https://github.com/docker/notary/releases/tag/v0.6.1)

* Bug fixes and minor changes
  - Fix startup issue due to incompatibility with other programs (like Razer Synapse 3). Fixes [docker/for-win#1723](https://github.com/docker/for-win/issues/1723)
  - Fix Kubernetes hostPath translation for PersistentVolumeClaim (PVC). Previously failing PVCs must be deleted and recreated. Fixes [docker/for-win#1758](https://github.com/docker/for-win/issues/1758)
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
  - Fix for the HTTP/S transparent proxy when using "localhost" names (for example "host.docker.internal", "docker.for.win.host.internal", "docker.for.win.localhost").
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
  - Virtual machine disk size can be changed in settings. Fixes [docker/for-win#105](https://github.com/docker/for-win/issues/105)
  - New menu item to restart Docker.

* Bug fixes and minor changes
  - Migration of Docker Toolbox images is not proposed anymore in Docker For Windows installer (still possible to migrate Toolbox images manually).

### Docker Community Edition 18.02.0-ce-rc1-win50 2018-01-26

* Upgrades
  - [Docker 18.02.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v18.02.0-ce-rc1)

* New
  - Experimental Kubernetes Support. You can now run a single-node Kubernetes cluster from the "Kubernetes" Pane in Docker for Windows settings and use kubectl commands as well as docker commands. See https://docs.docker.com/docker-for-windows/kubernetes/
  - LCOW containers can now be run next to Windows containers (on Windows RS3 build 16299 and later). Use `--platform=linux` in Windows container mode to run Linux Containers On Windows. Note that LCOW is still experimental, and requires daemon `experimental` option.

* Bug fixes and minor changes
  - Better cleanup for Windows containers and images on reset/uninstall. Fixes [docker/for-win#1580](https://github.com/docker/for-win/issues/1580), [docker/for-win#1544](https://github.com/docker/for-win/issues/1544), [docker/for-win#191](https://github.com/docker/for-win/issues/191)
  - Do not recreate Desktop icon on upgrade (effective on next upgrade). Fixes [docker/for-win#246](https://github.com/docker/for-win/issues/246), [docker/for-win#925](https://github.com/docker/for-win/issues/925), [docker/for-win#1551](https://github.com/docker/for-win/issues/1551)
  - Fix proxy for docker.for.win.localhost and docker.for.win.host.internal. Fixes [docker/for-win#1130](https://github.com/docker/for-win/issues/1130)

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
  - Fix VPNkit issue when username has spaces. See [docker/for-win#1429](https://github.com/docker/for-win/issues/1429).
  - Diagnostic improvements to get virtual machine logs before virtual machine shutdown.

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
  - Better removal of LCOW images and containers when uninstalling Docker
  - Links in update changelog open the default browser instead of IE (fixes [docker/for-win#1311](https://github.com/docker/for-win/issues/1311))

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
  - Increased timeout for virtual machine boot startup to 2 minutes.

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
  - Fix centos:5 and centos:6 images not starting properly with LinuxKit virtual machine (fixes [docker/for-win#1245](https://github.com/docker/for-win/issues/1245)).

### Docker Community Edition 17.10.0-ce-win36 2017-10-24

* Upgrades
  - [Docker 17.10.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.10.0-ce)
  - [Docker Machine 0.13.0](https://github.com/docker/machine/releases/tag/v0.13.0)
  - [Docker compose 1.17.0-rc1](https://github.com/docker/compose/releases/tag/1.17.0-rc1)

* New
  - Virtual machine entirely built with Linuxkit
  - Experimental support for Microsoft Linux Containers On Windows, on Windows 10 RS3.


### Docker Community Edition 17.09.0-ce-win34 2017-10-06

* Bug fixes
  - Fix Docker For Windows unable to start in some cases: removed use of libgmp sometimes causing the VPNkit process to die.


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
  - VPNkit: do not block startup when ICMP permission is denied. (Fixes docker/for-win#1036, docker/for-win#1035, docker/for-win#1040)
  - Fix minor bug on update checks

### Docker Community Edition 17.09.0-ce-rc1-win28 2017-09-07

* Upgrades
  - [Docker 17.09.0-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce-rc1)
  - [Docker compose 1.16.1](https://github.com/docker/compose/releases/tag/1.16.1)
  - Linux Kernel 4.9.46

* New
  - Add `Skip this version` button in update window

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

- VPNKit: Improved the logging around the UNIX domain socket connections
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
- Security fix: Users must be part of a specific group "docker-users" to run Docker for Windows

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

- Reset to default / uninstall and also reset docker cli settings and logout user from Docker Cloud and registries
- Detected a bitlocker policy preventing windows containers from working
- Fixed an issue on filesharing when explicitly disabled on VMswitch interface
- Fixed virtual machine not starting when machine had very long name
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
- Fix negotiation of TimeSync protocol version (through kernel update)

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

- You cannot edit settings while Docker is restarting
- Support Copy/Paste in About box
- Auto update polling every 24h
- Kernel boots with vsyscall=emulate arg and CONFIG_LEGACY_VSYSCALL is set to NONE in Moby
- Fixed vsock deadlock under heavy write load
- If you opt-out of analytics, you're prompted for approval before a bug report is sent
- Fixed bug where search domain could be read as `DomainName`
- Dedicated preference pane for HTTP proxy settings.
- Dedicated preference pane for CPU and Memory computing resources.
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

- Windows containers settings panel and options are working. In previous releases, settings were not implemented for
[Windows containers mode](index.md#switch-between-windows-and-linux-containers). (See
[About the Docker Windows containers specific dialogs](index.md#about-the-docker-windows-containers-specific-dialogs).)
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
- Virtual machine's time synchronization is forced after the host wakes from sleep mode

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
requires update `KB3192366` (soon to be released through Windows Update) to fully
work. Some insider builds may not work.

**New**

- Restore the virtual machine's configuration when it was changed by the user
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
requires update `KB3192366` (soon to be released through Windows Update) to fully
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
- VPNKit is now restarted if it dies
- Make sure invalid "DockerNat" switches are not used
- Preserve the PowerShell stacktraces
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
* Kernel 4.4.21
* AUFS 20160912

**Bug fixes and minor changes**

* Improve the switching between Linux and Windows containers: better errors, more reliable, deal with more edge cases
* Kill lingering dockerd that users might have still around because they played with Windows Containers before
* Don't recreate the virtual machine if only the DNS server is set
* The uninstaller now kills the service if it failed to stop it properly
* Restart VPNKit and DataKit when the processes die
* VPNKit: impose a connection limit to avoid exhausting file descriptors
* VPNKit: handle UDP datagrams larger than 2035 bytes
* VPNKit: reduce the number of file descriptors consumed by DNS
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
* AUFS 20160905

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
* VPNKit supports search domains
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
* AUFS 20160822

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
* AUFS 20160808

**Bug fixes and minor changes**

* Fixed password handling for host file system sharing
* Automatically disable lingering net adapters that prevent Docker from starting or using the network
* Automatically delete duplicated MobyLinuxVMs on a `reset to factory defaults`
* Docker stops asking to import from toolbox after an upgrade
* Docker can now import from toolbox just after Hyper-V is activated
* Fixed Moby Diagnostics and Update Kernel
* Added more debug information to the diagnostics
* Sending anonymous statistics shouldn't hang anymore when Mixpanel is not available
* Improved the Hyper-V detection and activation mechanism
* VPNKit is now compiled with OCaml 4.03 rather than 4.02.3
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

* Docker for Windows is now available from two channels: **stable** and **beta**. New features and bug fixes go out first in auto-updates to users in the beta channel. Updates to the stable channel are much less frequent and happen in sync with major and minor releases of the Docker engine. Only features that are well-tested and ready for production are added to the stable channel releases. For downloads of both and more information, see the [Getting Started](index.md#download-docker-for-windows).

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
* When killing the virtual machine process, ignore when the process is already stopped
* When stopping the virtual machine, always stop the docker proxy
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
* Removed the console from Hyper-V manager
* Improved diagnostic for VPN connection and addedlogs for the service port openers
* Improve Moby's boot sequence to adapt to longer boot time when swarm services are running
* Forcefully turn off a virtual machine that won't shut down
* Clicking on a link from the changelog opens a browser
* Fix links to the documentation
* Fix the url to download Kitematic
* Renewed the signing certificates
* Fixed errors with the firewall and the network switch
* Fixed parsing errors in the PowerShell script

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
* Linux 4.4.12, AUFS 20160530

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
* Removed DHCP for virtual machine network configuration
* User configurable NAT prefix and DNS server
* New feedback window to upload diagnostics dialog
* New status indicator in **Settings** window
* Virtual machine logs are uploaded with a crash report
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
* Enabling Hyper-V (was broken in Beta12)
* Fixed error message when invalid labels are passed to `docker run`
* Mixpanel no longer uses roaming App Data
* UI improvements
* Support was added for virtual machines with other IP addresses out of the `10.0.75.0/24` range
* Improved FAQ

**Known issues**

* Due to limitation in the Windows NAT implementation, co-existence with other NAT prefixes needs to be carefully managed. See [NAT Configuration](troubleshoot.md#nat-configuration) in [Troubleshooting](troubleshoot.md) for more details.

### Beta 12 Release (2016-17-10 1.11.1-beta12)

**New**

* The application is now separated in two parts. A back-end service and a front-end GUI. The front-end GUI no longer asks for elevated access.

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
* Allocate virtual machine memory by 256 MB increments, instead of 1 GB
* Show a meaningful error when the user has an empty password
* Improved [Troubleshooting](troubleshoot.md) page

**Upgrades**

* docker-compose 1.7.1  (see <a href="https://github.com/docker/compose/releases/tag/1.7.1" target="_blank"> changelog</a>)
* Kernel 4.4.9

**Bug fixes and minor changes**

* Report the virtual machine's IP in `docker port`
* Handle passwords with spaces
* Show a clear error message when trying to install on Home editions
* Slower whale animation in the system tray
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

* Improved Settings panel, allows you to configure the virtual machine’s memory and CPUs
* Co-exists with multiple internal Hyper-V switches and improved DHCP handling
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

* If multiple internal Hyper-V switches exist the Moby virtual machine
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
* Configured MobyLinux virtual machine with a fixed memory of 2GB
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
- Support Windows paths with -v
- Updated Docker CLI and daemon to 1.11.0-rc3
- Updated docker-machine to 0.7.0-rc2
- Updated docker-compose to 1.7.0-rc1
- Now install docker-credential-wincred
- Allow non-root users in containers to create files on volume mounts
- Automatically install Hyper-V
- The application is now 64 bits
- Improved wording in all dialog boxes and error messages
- Removed exit confirmation
- Show clickable URL in the Install HyperV message box
- Dashboard link to Kitematic (as on Mac)
- Moby kernel updated to 4.4.6
- The registry key was changed to HKLM\SOFTWARE\Docker Inc.\Docker\1.0

**Known issues**

- Migration from Docker Toolbox can fail sometimes. If this happens, the workaround is to restart the application.

- Docker needs to open ports on the firewall, which can activate a firewall alert dialog. Users should allow the ports to be opened.

- The application was upgraded to 64 bits. The installation path changed to `C:\Program Files\Docker\Docker`. If users have PowerShell/Cmd windows already open before the update, they might need to close them to catch the new PATH. In some cases, users need to log off and on again.

**Bug Fixes**

- Kill virtual machines that cannot be shut down properly

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
* Added CloudFront for downloads from Europe
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
- The application refuses to start if Hyper-V is not enabled
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

  - retrieving credentials from user

**Filesystem**

  - hot-mounting host filesystem with credential

**General**

  - state management
  - stability, logging
  - bugfixes, eye candies

### Alpha 0 Release (2016-02-09 1.10.0.0-0)

**Hypervision**

  - Hyper-V-backed virtual machines
  - boots Moby in a few seconds
  - installs CLI in `PATH`
  - proxies Docker commands to Moby

**Filesystem**

  - mounts host filesystem to support `--volume`
  - Samba client with a hardcoded password
  - allows live reload

**Networking**

  - live debugging Node.js application
