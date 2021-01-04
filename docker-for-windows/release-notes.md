---
description: Docker Desktop for Windows Release notes
keywords: Docker Desktop for Windows, release notes
redirect_from:
- /winkit/release-notes/
- /docker-for-windows/edge-release-notes/
title: Docker for Windows release notes
toc_min: 1
toc_max: 2
---

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

> **Important**
>
> Starting with Docker Desktop 3.0.0, Stable and Edge releases are combined into a single release stream for all users. Updates to Docker Desktop will now be available automatically as delta updates from the previous version. This means, when there is a newer version of Docker Desktop, it will be automatically downloaded to your machine. All you need to do is to click **Update and restart** from the Docker menu to install the latest update.
{: .important }

## Docker Desktop Community 3.0.0
2020-12-10

> [Download](https://desktop.docker.com/win/stable/50684/Docker%20Desktop%20Installer.exe)

### New

- Use of three-digit version number for Docker Desktop releases.
- Docker Desktop updates are now much smaller as they will be applied using delta patches. For more information, see [Automatic updates](install.md#automatic-updates).
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
- Fixed an unexpected EOF error when trying to start a non-existing container with `-v /var/run/docker.sock:`. See [docker/for-mac#5025](https://github.com/docker/for-mac/issues/5025).

### Known issues

- Building an image with BuildKit from a git URL fails when using the form `github.com/org/repo`. To work around this issue, use the form `git://github.com/org/repo`.
- Some DNS addresses fail to resolve within containers based on Alpine Linux 3.13.

## Docker Desktop Community 2.5.0.1
2020-11-10

> [Download](https://desktop.docker.com/win/stable/49550/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Compose CLI v1.0.2](https://github.com/docker/compose-cli/releases/tag/v1.0.2)
- [Snyk v1.424.4](https://github.com/snyk/snyk/releases/tag/v1.424.4)

## Docker Desktop Community 2.5.0.0
2020-11-02

> [Download](https://desktop.docker.com/win/stable/49427/Docker%20Desktop%20Installer.exe)

Docker Desktop 2.5.0.0 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### New

- Users subscribed to a Pro or a Team plan can now see the vulnerability scan report on the Remote repositories tab in Docker Desktop.
- Docker Desktop introduces a support option for users who have subscribed to a Pro or a Team Plan.

### Upgrades

- [Linux kernel 5.4.39](https://hub.docker.com/layers/linuxkit/kernel/5.4.39-f39f83d0d475b274938c86eaa796022bfc7063d2/images/sha256-8614670219aca0bb276d4749e479591b60cd348abc770ac9ecd09ee4c1575405?context=explore)
- [Docker Compose CLI 1.0.1](https://github.com/docker/compose-cli/releases/tag/v1.0.1)
- [Snyk v1.421.1](https://github.com/snyk/snyk/releases/tag/v1.421.1)
- [Go 1.15.2](https://github.com/golang/go/releases/tag/go1.15.2)
- [Kubernetes 1.19.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.19.3)

### Deprecation

- Docker Desktop cannot be installed on Windows 1703 (build 15063) anymore.

### Bug fixes and minor changes

- Renamed 'Run Diagnostics' to 'Get support'.
- Fixed an issue that intermittently made the WSL 2 backend fail to start.
- Fixed an issue related to NFS mounting. See [docker/for-mac#4958](https://github.com/docker/for-mac/issues/4958).
- Fixed an issue where bash in WSL would not find the docker CLI when it is started before Docker Desktop.
- Fixed an issue with HTTP proxy exclude lists containing entries such as `localhost` or `127.0.0.1`. Fixes [docker/for-win#8750](https://github.com/docker/for-win/issues/8750).
- When the WSL integration process unexpectedly stops, the user is now notified and can decide to restart it or not, instead of always trying to restart it in a loop. fixes [docker/for-win#8968](https://github.com/docker/for-win/issues/8968).
- Fixed an issue related to container logs lagging under heavy load. Fixes [docker/for-win#8216](https://github.com/docker/for-win/issues/8216).
- Diagnostics: avoid hanging when Kubernetes is in a broken state.
- Fixed the path to the installer log file default location when the username contains a space. Fixes [docker/for-win#7941](https://github.com/docker/for-win/issues/7941).
- Fixed an issue where some network plugins may fail to load, provoking a crash of the Docker daemon [docker/for-win#9282](https://github.com/docker/for-win/issues/9282).
- When sharing a file into a container (e.g. `docker run -v ~/.gitconfig`) Docker Desktop does not watch the parent directory. Fixes [docker/for-mac#4981](https://github.com/docker/for-mac/issues/4981).

## Docker Desktop Community 2.4.0.0
2020-09-30

> [Download](https://desktop.docker.com/win/stable/48506/Docker%20Desktop%20Installer.exe)

Docker Desktop 2.4.0.0 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### New

- [Docker Compose CLI - 0.1.18](https://github.com/docker/compose-cli), enabling use of volumes with Compose and the Cloud through ECS and ACI.
- Docker introduces the new Images view in the Docker Dashboard. The images view allows users to view the Hub images, pull them and manage their local images on disk including cleaning up unwanted and unused images. To access the new Images view, from the Docker menu, select **Dashboard** > **Images**.
- Docker Desktop now enables BuildKit by default after a reset to factory defaults. To revert to the old `docker build` experience, go to **Preferences** > **Docker Engine** and then disable the BuildKit feature.
- [Amazon ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper/releases/tag/v0.4.0)

### Upgrades

- [Docker 19.03.13](https://github.com/docker/docker-ce/releases/tag/v19.03.13)
- [Docker Compose 1.27.4](https://github.com/docker/compose/releases/tag/1.27.4)
- [Go 1.14.7](https://github.com/golang/go/releases/tag/go1.14.7)
- [Alpine 3.12](https://alpinelinux.org/posts/Alpine-3.12.0-released.html)
- [Kubernetes 1.18.8](https://github.com/kubernetes/kubernetes/releases/tag/v1.18.8)
- [Qemu 4.2.0](https://git.qemu.org/?p=qemu.git;a=tag;h=1e4aa2dad329852aa6c3f59cefd65c2c2ef2062c)

### Bug fixes and minor changes

- Removed the legacy Kubernetes context `docker-for-desktop`. The context `docker-desktop` should be used instead. Fixes [docker/for-win#5089](https://github.com/docker/for-win/issues/5089) and [docker/for-mac#4089](https://github.com/docker/for-mac/issues/4089).
- Removed the option to start with Windows containers from the installer
- Pinning the application to the taskbar and clicking on it will launch the container view if Docker is already running.
- Left-clicking on the whale in the system tray now launches the dashboard container view.
- Docker Desktop now uses a different systray icon for dark and light mode. Fixes [docker/for-win#4113](https://github.com/docker/for-win/issues/4113).
- Added support for emulating Risc-V via Qemu 4.2.0.
- Added a low-level debug shell accessible via `putty -serial \\.\pipe\dockerDebugShell`.
- Copy container logs without ansi colors to clipboard. Fixes [docker/for-mac#4786](https://github.com/docker/for-mac/issues/4786).
- Fix application startup if `hosts` is specified inside the Docker `daemon.json`. See [docker/for-win#6895](https://github.com/docker/for-win/issues/6895#issuecomment-637429117)
- Fixed DNS resolution of short names. See [docker/for-win#4425](https://github.com/docker/for-win/issues/4425).
- Switched from `chronyd` to `sntpcd` to work around host time synchronisation problems. Fixes [docker/for-win#4526](https://github.com/docker/for-win/issues/4526).
- Avoid blocking startup if "Expose daemon on tcp://localhost:2375 without TLS" is set and `localhost:2375` is in use by another program. See [docker/for-win#6929](https://github.com/docker/for-win/issues/6929) [docker/for-win#6961](https://github.com/docker/for-win/issues/6961).
- Fixed an issue where adding a folder on a non-existing drive in the settings would create an empty entry. See [docker/for-win#6797](https://github.com/docker/for-win/issues/6797).
- Avoid failing with "Function not implemented" during file I/O on shared volumes. Fixes [docker/for-win#5955](https://github.com/docker/for-win/issues/5955)
- Ensure that `docker run -v /var/run/docker.sock` rewrites Windows paths properly, see [docker/for-win#6628](https://github.com/docker/for-win/issues/6628).
- Fixed a crash which occurred when Docker Desktop loads a corrupted Docker CLI configuration file. Fixes [docker/for-win#6657](https://github.com/docker/for-win/issues/6657).
- Ensure `localhost` and `127.0.0.1` can both be used in the proxy settings to redirect to a proxy on the host. Fixes [docker/for-win#5715](https://github.com/docker/for-win/issues/5715) and [docker/for-win#6260](https://github.com/docker/for-win/issues/6260).
- Fixed a crash when failing to login with no Internet connection.
- Fixed bug in handling shared volume paths with ".." characters. Fixes [docker/for-win#5375](https://github.com/docker/for-win/issues/5375).
- Report check for updates errors in toast notification. Fixes [docker/for-win#6364](https://github.com/docker/for-win/issues/6364).
- Fixed an upgrade bug where users on versions still using the PowerShell based VM management can experience a silent un-installation crash leading to Docker Desktop being uninstalled instead of upgraded.
- Fixed the path to the installer log file default location when the username contains a space. Fixes [docker/for-win#7941](https://github.com/docker/for-win/issues/7941).
- Docker Desktop always flushes filesystem caches synchronously on container start. See [docker/for-mac#4943](https://github.com/docker/for-mac/issues/4943).
- Compose-on-Kubernetes is no longer included in the Docker Desktop installer. You can download it separately from the compose-on-kubernetes [release page](https://github.com/docker/compose-on-kubernetes/releases).

### WSL 2 changes

- Docker contexts are now synced between Windows and WSL distros.
- Fixed an issue that intermittently causes the backend to fail to start.
- Fixed a proxy crash that happened when glibc was not compatible. See [docker/for-win#8183](https://github.com/docker/for-win/issues/8183).
- Fixed an issue when removing a container mounting `/mnt/wsl` would break WSL integration. See [docker/for-win#7836](https://github.com/docker/for-win/issues/7836).
- Added support for mounting files from a distro using the Windows CLI (e.g. `docker run -v \\wsl$\Ubuntu\home\simon\web:/web ...`)
- Fixed the error message when attempting to use a shared volume with a relative path. Fixes [docker/for-win#6894](https://github.com/docker/for-win/issues/6894).
- Fixed an issue when updating Windows to a WSL 2 capable version, where a configuration file is locked by the old Hyper-V VM.
- Fixed a discrepancy in the version of Docker Compose inside and outside of WSL 2. Fixes [docker/for-win#6461](https://github.com/docker/for-win/issues/6461).
- Detect when the `docker-desktop` wsl distro is stopped and display a clearer error message.
- Fix a race condition when exposing ports.
- Enable dialog no longer blocks other windows.

### Known issues

-  There is a known issue when using `docker-compose` with named volumes and gRPC FUSE: second and subsequent calls to `docker-compose up` will fail due to the volume path having the prefix `/host_mnt`.
- There is a known issue when enabling Kubernetes where the settings UI sometimes fails to update the Kubernetes state. To work around this issue, close and re-open the Window.
- There is a rare known issue when switching users, where the images UI sometimes continues to show the repositories of the previous user. To work around this issue, close and re-open the Window.

## Docker Desktop Community 2.3.0.5
2020-09-15

> [Download](https://desktop.docker.com/win/stable/48029/Docker%20Desktop%20Installer.exe)

### New

- The new Cloud integration in Docker CLI makes it easy to run containers in the cloud using either Amazon ECS or Microsoft ACI. For more information, see [Deploying Docker containers on ECS](https://docs.docker.com/engine/context/ecs-integration/) and [Deploying Docker containers on Azure](https://docs.docker.com/engine/context/aci-integration/).

### Upgrades

- [Docker Compose 1.27.2](https://github.com/docker/compose/releases/tag/1.27.2)
- [Cloud integration v0.1.15](https://github.com/docker/aci-integration-beta/releases/tag/v0.1.15)

### Bug fixes and minor changes

- WSL2: Fixed a crash when using an incompatible glibc. See [docker/for-win#8183](https://github.com/docker/for-win/issues/8183).

### Known issues

- The `clock_gettime64` system call returns `EPERM` rather than `ENOSYS`
in i386 images. To work around this issue, disable `seccomp` by using
the `--privileged` flag. See [docker/for-win#8326](https://github.com/docker/for-win/issues/8326).

## Docker Desktop Community 2.3.0.4
2020-07-27

> [Download](https://desktop.docker.com/win/stable/46911/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker 19.03.12](https://github.com/docker/docker-ce/releases/tag/v19.03.12)
- [Docker Compose 1.26.2](https://github.com/docker/compose/releases/tag/1.26.2)
- [Go 1.13.14](https://github.com/golang/go/issues?q=milestone%3AGo1.13.14+label%3ACherryPickApproved)

### Bug fixes and minor changes

- Docker Desktop now prompts users to share directories for paths like `////c/Users/foo`, rather than only paths like `C:\Users\foo` and `C:/Users/foo`.
- The installer now returns a non-zero exit code on error, or when the user cancels the installation.
- Fixed the path to the installer log file default location when the username contains a space. Fixes [docker/for-win#6552](https://github.com/docker/for-win/issues/6552)
- Dashboard: Fixed opening a CLI for Windows containers. See [docker/for-win#7079](https://github.com/docker/for-win/issues/7079)
- Dashboard: Fixed containers logs which were sometimes truncated. Fixes [docker/for-win#5954](https://github.com/docker/for-win/issues/5954)
- WSL 2: Fixed an issue where changing the user's default shell blocked WSL integration. Fixes [docker/for-win#7653](https://github.com/docker/for-win/issues/7653)
- WSL 2: Fixed an issue to recover WSL distributions stuck in the "Installing" state.

## Docker Desktop Community 2.3.0.3
2020-05-27

> [Download](https://desktop.docker.com/win/stable/45519/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Linux kernel 4.19.76](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.76-83885d3b4cff391813f4262099b36a529bca2df8-amd64/images/sha256-0214b82436af70054e013ea51cb1fea72bd943d0d6245b6521f1ff09a505c40f?context=repo)

### Bug fixes and minor changes

- Fixed a crash when resetting to factory defaults with Hyper-V disabled. Fixes [docker/for-win#6738](https://github.com/docker/for-win/issues/6738).
- Fixed opening an application in VS Code from a WSL 2 running app. Fixes [docker/for-win#6472](https://github.com/docker/for-win/issues/6472).
- Fixed Swarm mounts in WSL 2. Fixes [docker/for-win#6507](https://github.com/docker/for-win/issues/6507).
- Fixed a bug when using the Microsoft `mssql` image. Fixes [docker/for-win#6646](https://github.com/docker/for-win/issues/6646)
- Implemented `fallocate` for shared file systems. See [docker/for-win#6658](https://github.com/docker/for-win/issues/6658#issuecomment-627736820)
- Fixed an installer crash when an old and/or partially uninstalled version of Docker Desktop was present on the system. Fixes [docker/for-win#6536](https://github.com/docker/for-win/issues/6536).
- Fixed a delay caused when opening the systray menu. Fixes [docker/for-win#1011](https://github.com/docker/for-win/issues/1011).
- Fixed a regression where containers couldn’t anymore refer to mounted folders using the host drive letters. Fixes [docker/for-win#6628](https://github.com/docker/for-win/issues/6628).
- Fixed a regression where sharing a folder using a double leading slash notation would fail. Fixes [docker/for-win#6668](https://github.com/docker/for-win/issues/6668).
- Re-added device-mapper to the embedded Linux kernel. Fixes [docker/for-mac#4549](https://github.com/docker/for-mac/issues/4549).
- Fixed an issue which caused the `:z` attribute in bind mounts to fail. Fixes [docker/for-win#6634](https://github.com/docker/for-win/issues/6634).

## Docker Desktop Community 2.3.0.2
2020-05-11

> [Download](https://download.docker.com/win/stable/45183/Docker%20Desktop%20Installer.exe)

### New

- Windows 10 Home users can now use Docker Desktop through WSL 2. This requires Windows 10, version 2004 or higher. For more information, see [Install Docker Desktop on Windows Home](install-windows-home.md).
- Docker Desktop introduces a new onboarding tutorial upon first startup. The Quick Start tutorial guides users to get started with Docker in a few easy steps. It includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.
- Docker Desktop now allows sharing individual folders, rather than whole drives, giving more control to users over what is being shared.

### Upgrades

- [Docker Compose 1.25.5](https://github.com/docker/compose/releases/tag/1.25.5)
- [Go 1.13.10](https://github.com/golang/go/issues?q=milestone%3AGo1.13.10+label%3ACherryPickApproved)
- [Linux kernel 4.19.76](https://hub.docker.com/layers/docker/for-desktop-kernel/4.19.76-ce15f646db9b062dc947cfc0c1deab019fa63f96-amd64/images/sha256-6c252199aee548e4bdc8457e0a068e7d8e81c2649d4c1e26e4150daa253a85d8?context=repo)
- LinuxKit [init](https://hub.docker.com/layers/linuxkit/init/1a80a9907b35b9a808e7868ffb7b0da29ee64a95/images/sha256-64cc8fa50d63940dbaa9979a13c362c89ecb4439bcb3ab22c40d300b9c0b597e?context=explore), [runc](https://hub.docker.com/layers/linuxkit/runc/69b4a35eaa22eba4990ee52cccc8f48f6c08ed03/images/sha256-57e3c7cbd96790990cf87d7b0f30f459ea0b6f9768b03b32a89b832b73546280?context=explore) and [containerd](https://hub.docker.com/layers/linuxkit/containerd/09553963ed9da626c25cf8acdf6d62ec37645412/images/sha256-866be7edb0598430709f88d0e1c6ed7bfd4a397b5ed220e1f793ee9067255ff1?context=explore)

### Bug fixes and minor changes

**WSL 2**

- Docker Desktop only exposes host ports in Linux if the Windows port is available.
- Docker Desktop now allows users to refresh the list of distros.
- Docker Desktop defaults to WSL 2 on install on compatible OS versions.
- Docker Desktop detects missing Linux kernel and adds a pointer to the Microsoft documentation to download the kernel.
- Detect when the WSL 2 backend stops and allow the user to restart it.
- Added support for `chmod/chown` on WSL 2 bind mounts. Fixes [docker/for-win#6284](https://github.com/docker/for-win/issues/6284).
- Added a check that virtualization is enabled in the BIOS.

**Filesharing**

- Fixed a bug where renaming parent directories of shared folders led to spurious `File not found` errors. Fixes [docker/for-win#6200](https://github.com/docker/for-win/issues/6200).
- Fixed bug where a letter and two digit root folder name causes docker-compose to fail creating directories inside the volume. Fixes [docker/for-win#6248](https://github.com/docker/for-win/issues/6248).
- Fixed a bug which prevented containers from seeing file updates on shared volumes, when under load and a container restarts. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530#issuecomment-608804192).
- Fixed a bug where host paths are incorrectly translated into VM paths. Fixes [docker/for-win#6209](https://github.com/docker/for-win/issues/6209).
- Fixed bug receiving file events on long paths (> 260 characters). Fixes [docker/for-win#6337](https://github.com/docker/for-win/issues/6337).
- Docker Desktop represents valid directory junctions as directories (rather than symlinks) and handles cache invalidation and event injection properly. Fixes [docker/for-win#5582](https://github.com/docker/for-win/issues/5582).

**Other fixes**

- Reduced the size of the Docker Desktop installer from 960 MB to 409 MB.
- Added an option to delete data from the Troubleshoot screen.
- Fixed a bug where containers disappeared from the UI when Kubernetes context is invalid. Fixed [docker/for-win#6037](https://github.com/docker/for-win/issues/6037).
- Fixed Windows event logs filtering when copying them to the Docker Desktop log files. Fixed [docker/for-win#6258](https://github.com/docker/for-win/issues/6258).
- Fixed a handle leak in `vpnkit-bridge`. Fixed [docker/for-win#5841](https://github.com/docker/for-win/issues/5841)
- Fixed bug when removing Docker Desktop virtual switch.
- Added a link to the Edge channel from the UI.
- Made the embedded terminal resizable.
- Fixed a bug which only respected the `expose on TCP` Docker engine API setting on application restart. Setting will now take effect when the Apply button is clicked.
- Fixed a bug where diagnostic upload would fail if the username contained spaces.
- Docker Desktop now reads the Hyper-V VM disk max size on startup and uses it as the value to display in the Settings.
- Fixed a bug where the Docker Desktop UI could be started without the engine. Fixes [docker/for-win#5376](https://github.com/docker/for-win/issues/5376).
- Docker Desktop now uses the least possible privilege for querying the Server service. Fixes [docker/for-win#5150](https://github.com/docker/for-win/issues/5150).
- Fixed an issue where a container port could not be exposed on a specific host IP. See [docker/for-win#5546](https://github.com/docker/for-mac/issues/5546).
- Removed port probing from dashboard, just unconditionally showing links to ports that should be available. Fixes [docker/for-win#5903](https://github.com/docker/for-win/issues/5903).
- Ceph support has been removed from Docker Desktop to save disk space.
- Fixed an issue that caused the Windows log file archive to grow indefinitely. Fixes [docker/for-win#5113](https://github.com/docker/for-win/issues/5113).
- Added a prerequisite check in the installer for the LanmanServer service. Fixes [docker/for-win#5150](https://github.com/docker/for-win/issues/5150)

### Known issues

- Some CLI commands fail if you are running Docker Desktop in the experimental Linux Containers on Windows (LCOW) mode. As alternatives, we recommend running either traditional Linux containers, or the [WSL 2 backend](wsl.md).

**WSL 2**

- Swarm service bind mounts are not always restored correctly.
- Bind mounts of files living outside of the root mount point (files within `/mnt/c`, `/tmp`, `/run...`) don't work correctly when multiple containers mount them.

## Docker Desktop Community 2.2.0.5
2020-04-02

> [Download](https://download.docker.com/win/stable/43884/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

- Docker Desktop generates `fsnotify.WRITE` events in Linux containers when files are changed on the host. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530#issuecomment-585572414).
- Fixed a race condition in `readlink` on shared volumes using "mfsymlinks". Fixes [docker/for-win#5793](https://github.com/docker/for-win/issues/5793).
- Made the VM time sync more reliable. See [docker/for-win#4526](https://github.com/docker/for-win/issues/4526).
- Fixed bug where files did not change in shared volumes when using `volumes_from` in `docker-compose.yml`. See [docker/for-win#5530](https://github.com/docker/for-win/issues/5530).
- Fixed a bug where opening a read-only file would fail with an `Operation not permitted` error. Fixes [docker/for-win#6016](https://github.com/docker/for-win/issues/6016) and [docker/for-win#6017](https://github.com/docker/for-win/issues/6017).

## Docker Desktop Community 2.2.0.4
2020-03-13

> [Download](https://download.docker.com/win/stable/43472/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker 19.03.8](https://github.com/docker/docker-ce/releases/tag/v19.03.8)

### Bug fixes and minor changes

- Security: Diagnostics were collected with Administrator privileges, which led to a possible privilege escalation on systems where the Docker Desktop user was not an administrator.
- Docker Desktop now displays hidden files in shared volumes. Fixes [docker/for-win#5808](https://github.com/docker/for-win/issues/5808).
- Docker Desktop now generates the inotify `MODIFY` events on shared file systems for Windows file sharing. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530).
- Attempts to create files in a shared volume with an identical filename but a different case (uppercase/lowercase) on Windows filesystem will now fail with the error `EEXIST`. Fixes [docker/for-win#5894](https://github.com/docker/for-win/issues/5894).
- Fixed cache invalidation and event injection in shared volumes with host paths longer than 260 characters.
- Docker Desktop now allows users to rename open files in shared volumes. Fixes [docker/for-win#5565](https://github.com/docker/for-win/issues/5565).
- Fixed an issue which caused Docker Desktop Dashboard to attempt connecting to all exposed ports inside a container. Fixes [docker/for-win#5903](https://github.com/docker/for-win/issues/5903).
- Kubernetes: Persistent volumes created by claims are now stored in the virtual machine. Fixes [docker/for-win#5665](https://github.com/docker/for-win/issues/5665).
- Fixed an issue which caused Docker Desktop to hang when users tried to reset to factory defaults.
- Fixed a file sharing issue that caused Docker Desktop to lock random files. Fixes [docker/for-win#5624](https://github.com/docker/for-win/issues/5624) and [docker/for-win#5575](https://github.com/docker/for-win/issues/5575).

### Known issues

- Some CLI commands fail if you are running Docker Desktop in the experimental Linux Containers on Windows (LCOW) mode. As alternatives, we recommend running either traditional Linux containers, or the experimental [WSL backend](wsl.md).
- It is not possible to resize the disk image using the Docker Desktop **Settings** UI. If you would like to update the size of the disk image (for example, to 128 GB), run the following command in PowerShell:

  ```powershell
  Resize-VHD -Path 'C:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx' -SizeBytes 128gb
  ```

## Docker Desktop Community 2.2.0.3
2020-02-11

> [Download](https://download.docker.com/win/stable/42716/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Docker Compose 1.25.4](https://github.com/docker/compose/releases/tag/1.25.4)
- [Go 1.12.16](https://golang.org/doc/devel/release.html#go1.12)

### Bug fixes and minor changes

- Fixed an issue that prevented users from creating files with special characters in the filenames within a shared volume. Fixes [docker/for-win#5520](https://github.com/docker/for-win/issues/5520).
- Fixed handling of shared volumes with relative paths in `docker-compose.yml`. Fixes [docker/for-win#5516](https://github.com/docker/for-win/issues/5516).
- Fixed handling of shared volumes where the path case (uppercase / lowercase) did not exactly match the host. Fixes [docker/for-win#5516](https://github.com/docker/for-win/issues/5516).
- Fixed an issue where changing a file in the Windows file system did not update the file within the container. Fixes [docker/for-win#5530](https://github.com/docker/for-win/issues/5530) and [docker/for-win#5550](https://github.com/docker/for-win/issues/5550).
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
- Starting with Docker Desktop 2.2.0.3, you must access all shared files using their original case. For example, if you have created a file called `test`, you must open it as `test`. Attempts to open the file as `Test` will fail with the error `No such file or directory`. For more information, see _Tips on shared drives, permissions, and volume mounts_ in [File sharing](index.md#file-sharing).

### Known issues

- DockerNAT has been removed from Docker Desktop 2.2.0.0 as using an IP address to communicate from the host to a container is not a supported feature. To communicate from a container to the host, you must use the special DNS name `host.docker.internal`. We also recommend using ports to communicate from the host to a container. For more information, see [Networking](/docker-for-win/networking/#use-cases-and-workarounds).

  However, if your current setup relies on IP addresses for communication, you can use a temporary workaround to reinstate DockerNAT. To do this, open
`C:\Program Files\Docker\Docker\resources\MobyLinux.ps1` and add `$SwitchName = "DockerNAT"` between line 175 and 176. Note that the temporary workaround to reinstate DockerNAT may be removed from future releases.
- Directory junctions in shared volumes pointing outside the shared volume do not work.
- It is currently not possible to rename open files within shared volumes.
- In some cases, Docker containers do not synchronize time with the host when the host machine is put to sleep.
- It is not possible to use `localhost` in the proxy settings. You can instead use `host.docker.internal` to work around this issue.

## Docker Desktop Community 2.2.0.0
2020-01-21

Docker Desktop 2.2.0.0 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker Compose 1.25.2](https://github.com/docker/compose/releases/tag/1.25.2)
- [Kubernetes 1.15.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.15.5)
- Linux kernel 4.19.76
- [QEMU 4.0.1](https://github.com/docker/binfmt)

### New

- **Docker Desktop Dashboard:** The new Docker Desktop **Dashboard** provides a user-friendly interface which enables you to interact with containers and applications, and manage the lifecycle of your applications directly from the UI. In addition, it allows you to access the logs, view container details, and monitor resource utilization to explore the container behavior.
For detailed information about the new Dashboard UI, see [Docker Desktop Dashboard](../desktop/dashboard.md).

- **WSL 2 backend:** The experimental Docker Desktop WSL 2 backend architecture introduces support for Kubernetes, provides an updated Docker daemon, offers VPN-friendly networking, and additional features. For more information, see [Docker Desktop WSL 2 backend](https://docs.docker.com/docker-for-windows/wsl/).

- **New file sharing implementation:** Docker Desktop introduces a new file sharing implementation that replaces Samba, CIFS, and Hyper-V networking. The new implementation  offers improved I/O performance. Additionally, when using the new file system:

  - Users don't have to expose the Samba port, and therefore do not experience issues related to IT firewall or drive-sharing policy.
  - There is no need to provide user credentials to Docker Desktop. File access rights are automatically enforced when accessing mounted folders through containers.

    For more information, see the blog post [ New file sharing implementation in Docker Desktop Windows](https://www.docker.com/blog/new-filesharing-implementation-in-docker-desktop-windows/).

- Added the **WSL Integration** option in **Settings** > **Resources** to select WSL 2 distributions. This option will be visible only after enabling WSL 2 on your machine.
- Introduced a new user interface for the Docker Desktop **Settings** menu.
- The Restart, Reset, and Uninstall options are now available on the **Troubleshoot** menu.
- Added the ability to start and stop existing Compose-based applications and view combined logs in the Docker Desktop **Dashboard** UI.

### Bug fixes and minor changes

- Enabled Windows features such as Hyper-V and Containers during installation, thereby reducing the need for another restart after installation.
- Added support for `Expect: 100-continue` headers in the Docker API proxy. Some HTTP clients such as `curl` send this header when the payload is large, for example, when creating containers. Fixes [moby/moby#39693](https://github.com/moby/moby/issues/39693).
- Added a loading overlay to the **Settings** and **Troubleshoot** windows to prevent editing conflicts.
- Deactivated the **Reset Kubernetes** button when Kubernetes is not activated.
- Improved the navigation in **Settings** and **Troubleshoot** UI.
- Docker Machine is no longer included in the Docker Desktop installer. You can download it separately from the [Docker Machine releases](https://github.com/docker/machine/releases) page.
- Docker Desktop now automatically restarts after an update.
- Fixed an issue where Docker Desktop auto-start was not being disabled properly on some machines.
- Fixed a container start error when a container has more than one port with an arbitrary or not-yet-configured external port number. For example, `docker run -p 80 -p 443 nginx`). Fixes [docker/for-win#4935](https://github.com/docker/for-win/issues/4935) and [docker/compose#6998](https://github.com/docker/compose/issues/6998).
- Fixed an issue which caused Docker Desktop to crash when resetting to factory defaults while running Windows containers.
- Fixed multiple issues related to Fast Startup.
- Docker Desktop now supports `inotify` events on shared filesystems.
- Docker Desktop startup is now more reliable and does not clash with host firewall software.
- Fixed a rare issue that caused to Docker Desktop to crash with the error `Unable to stop Hyper-V VM: Cannot validate argument on parameter 'SwitchName'. The argument is null or empty.`
- Fixed a bug that caused a rare crash when uninstalling Docker Desktop.
- Fixed an issue that caused Docker Desktop to fail on startup when there is an incomplete Kubernetes config file.
- Fixed various stability issues.

### Known issues

- When you start a Docker Compose application and then start a Docker App which has the same name as the Compose application, Docker Desktop displays only one application on the Dashboard. However, when you expand the application, containers that belong to both applications are displayed on the Dashboard.
- When you deploy a Docker App with multiple containers on Kubernetes, Docker Desktop displays each Pod as an application on the Dashboard.
- WSL 2 requires Windows 10 Insider Preview build 19018 or greater.
- The Dashboard stops updating when you switch the container mode between Linux and Windows. To work around this issue, close and reopen the Dashboard.
- The new gRPC FUSE file sharing implementation does not support connecting to new drives (for example, USB drives) added after Docker Desktop starts. If you would like to share the new drive in Docker Desktop, you must quit Docker Desktop and then start the application to refresh the list of drives in Settings.
- Sharing files from virtual drives created by `ImDisk Toolkit` does not work.

## Docker Desktop Community 2.1.0.5
2019-11-18

[Download](https://download.docker.com/win/stable/40693/Docker%20Desktop%20Installer.exe)

Docker Desktop 2.1.0.5 contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker 19.03.5](https://github.com/docker/docker-ce/releases/tag/v19.03.5)
- [Kubernetes 1.14.8](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.8)
- [Go 1.12.13](https://golang.org/doc/devel/release.html#go1.12)

## Docker Desktop Community 2.1.0.4
2019-10-21

[Download](https://download.docker.com/win/stable/39773/Docker%20Desktop%20Installer.exe)

Docker Desktop 2.1.0.4 contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker 19.03.4](https://github.com/docker/docker-ce/releases/tag/v19.03.4)
- [Kubernetes 1.14.7](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.7)
- [Go 1.12.10](https://github.com/golang/go/issues?q=milestone%3AGo1.12.10+label%3ACherryPickApproved)
- [Kitematic 0.17.9](https://github.com/docker/kitematic/releases/tag/v0.17.9)

### New

Docker Desktop now enables you to sign into Docker Hub using two-factor authentication. For more information, see [Two-factor authentication](index/#docker-hub).

## Docker Desktop Community 2.1.0.3
2019-09-16

[Download](https://download.docker.com/win/stable/38240/Docker%20Desktop%20Installer.exe)

### Upgrades

- [Kitematic 0.17.8](https://github.com/docker/kitematic/releases/tag/v0.17.8)

## Docker Desktop Community 2.1.0.2
2019-09-04

[Download](https://download.docker.com/win/stable/37877/Docker%20Desktop%20Installer.exe)

Docker Desktop 2.1.0.2 contains a Kubernetes upgrade. Your local Kubernetes cluster will be reset after installing this version.

### Upgrades

- [Docker 19.03.2](https://github.com/docker/docker-ce/releases/tag/v19.03.2)
- [Kubernetes 1.14.6](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.6)
- [Go 1.12.9](https://github.com/golang/go/issues?q=milestone%3AGo1.12.9+label%3ACherryPickApproved)
- [Docker Machine 0.16.2](https://github.com/docker/machine/releases/tag/v0.16.2)

## Docker Desktop Community 2.1.0.1
2019-08-08

[Download](https://download.docker.com/win/stable/37199/Docker%20Desktop%20Installer.exe)

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

### New

* Selecting the ‘Experimental features’ checkbox in the Daemon, **Settings** menu turns on experimental features for Docker daemon and Docker CLI.
* Docker Desktop now checks for stored credentials at startup before attempting to mount any shared drives. This prompts users to reenter the credentials if they are invalid.

### Experimental

> Experimental features provide early access to future product functionality. These features are intended for testing and feedback only as they may change between releases without warning or can be removed entirely from a future release. Experimental features must not be used in production environments. Docker does not offer support for experimental features.

Docker Desktop Community 2.1.0.0 contains the following experimental features:

* Docker App: Docker App is a CLI plugin that helps configure, share, and install applications. For more information, see [Working with Docker App](/app/working-with-app/).
* Docker Buildx: Docker Buildx is a CLI plugin for extended build capabilities with BuildKit. For more information, see [Working with Docker Buildx](/buildx/working-with-buildx/).

### Bug fixes and minor changes

* Fixed PowerShell script signing issue that caused `AuthorizationManager check failed` errors on machines with strict group policies on PowerShell script signing. [docker/for-win#4376](https://github.com/docker/for-win/issues/4376)
* Fixed an issue where users were unable to start Docker Desktop after upgrading to version 2.1.0.0. [docker/for-win#4390](https://github.com/docker/for-win/issues/4390)
* Fixed an issue where attempts to upgrade Docker Desktop to version 2.1.0.0 sometimes failed with the error `value cannot be null`. [docker/for-win#4343](https://github.com/docker/for-win/issues/4343)
* Fixed an issue that caused the installer to hang when upgrading Docker Desktop to version 2.1.0.0. [docker/for-win#4387](https://github.com/docker/for-win/issues/4387)
* Improved the error messages displayed during VM lifecycle operations. [docker/for-win#4348](https://github.com/docker/for-win/issues/4348)
* Docker Desktop now supports a configurable user timeout for VMs on slower machines. [docker/for-win#4393](https://github.com/docker/for-win/issues/4393)
* Restricted the `cluster-admin` role on local Kubernetes cluster to `kube-system` namespace.
* Reduced the VM startup time. swap is not created every time a virtual machine boots.
* Fixed a bug which caused Windows to crash when a user cancels switching the version using Windows User Account Control (UAC) settings.
* Fixed an issue where Docker Desktop restarts when a user logs out of Windows and logs back in, which results in retaining the exported ports on containers.
* Fixed Kubernetes installation with VPNkit subnet.
* Fixed a bug where the process output was not redirected to stdout when gathering diagnostics on Windows, which sometimes resulted in a crash.
* Fixed race condition where Kubernetes sometimes fails to start after the app is restarted.
* The system tray icon now opens the Docker Desktop menu with left and right mouse button.
* When displaying the crash report window, Docker Desktop does not send a bugsnag crash report unless the user needs the report to upload diagnostics.
* Docker Desktop has removed the ability to log in using email address as a username as the Docker command line does not support this.
* For Linux containers on Windows (LCOW), the host must run Windows 10 Professional or Windows 10 Enterprise version 1809 or later.
* The **Send usage statistics** checkbox is selected by default in Docker Desktop for Windows Community. This option cannot be modified.
* Docker Desktop has added a new dialog box during startup which allows users to retry mounting a shared drive or remove it from the shared drives list after a failed attempt.
* Fixed the Kubernetes Delete operation of persistent volume claims on Windows.
* Docker Desktop now truncates UDP DNS responses which are over 512 bytes in size.
* Fixed a bug where Docker was not added to PATH after installation in some cases.
* Fixed port 8080 that was used on localhost when starting Kubernetes.
* Renamed the product from **Docker for Windows** to **Docker Desktop**.
* Fixed the **Create issue** link in the Diagnostics window.
* Fixed a bug related to service log collection in diagnostics.
* Docker Desktop now gathers `/etc/hosts` to help with diagnostics.
* Added the missing daemon options in Docker Engine 18.09.
* When two services have a common exposed port, Docker Desktop exposes the available ports for the second service.
* Fixed the UI lock when changing the Kubernetes state.

### Docker Community Edition 2.0.0.3 2019-02-15

[Download](https://download.docker.com/win/stable/31259/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.09.2](https://github.com/docker/docker-ce/releases/tag/v18.09.2), fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

* Bug fix
  - Fix crash in system tray menu when the Hub login fails or Air gap mode

### Docker Community Edition 2.0.0.2 2019-01-16

[Download](https://download.docker.com/win/stable/30215/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.09.1](https://github.com/docker/docker-ce/releases/tag/v18.09.1)
  - [Docker Machine 0.16.1](https://github.com/docker/machine/releases/tag/v0.16.1)
  - [Kubernetes 1.10.11](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.10.md#v11011), fixes [CVE-2018-1002105](https://github.com/kubernetes/kubernetes/issues/71411)
  - [Kitematic 0.17.6](https://github.com/docker/kitematic/releases/tag/v0.17.6)
  - Golang 1.10.6, fixes CVEs: [CVE-2018-16875](https://www.cvedetails.com/cve/CVE-2018-16875), [CVE-2018-16873](https://www.cvedetails.com/cve/CVE-2018-16873) and [CVE-2018-16874](https://www.cvedetails.com/cve/CVE-2018-16874)
  - Windows 14393 is marked as deprecated ; it will not be supported anymore in the next major stable release (2.1.0.0 and further)

* Bug fixes and minor changes
  - Rename Docker for Windows to Docker Desktop
  - Add 18.09 missing daemon options

## Stable Releases of 2018

### Docker Community Edition 2.0.0.0-win81 2018-12-07

[Download](https://download.docker.com/win/stable/29211/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker compose 1.23.2](https://github.com/docker/compose/releases/tag/1.23.2)

* Bug fixes and minor changes
  - Compose: Fixed a bug where build context URLs would fail to build on Windows. Fixes [docker/for-win#2918](https://github.com/docker/for-win/issues/2918)

### Docker Community Edition 2.0.0.0-win78 2018-11-19

[Download](https://download.docker.com/win/stable/28905/Docker%20for%20Windows%20Installer.exe)

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

* Bug fixes and minor changes
  - Windows Containers: Fix group daemon option settings. Fixes [docker/for-win#2647](https://github.com/docker/for-win/issues/2647)
  - Windows Containers: Better host.docker.internal resolution, don't rewrite it if not modified. From [docker/for-win#1976](https://github.com/docker/for-win/issues/1976)
  - Do not try to update samba share mounts when using Windows Containers
  - Improved dns update too verbose in logs
  - VPNKit: Improved scalability of port forwarding.
  - VPNKit: Limit the size of the UDP NAT table. This ensures port forwarding and regular TCP traffic continue even when running very chatty UDP protocols.
  - LCOW does not need --platform flag on multi-arch images any more
  - Disk size can now be ajusted from 32GiB to drive space.
  - Fix dns update too verbose in logs
  - Fix panic in diagnose

### Docker Community Edition 18.06.1-ce-win73 2018-08-29

[Download](https://download.docker.com/win/stable/19507/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.06.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.06.1-ce)

* Bug fixes and minor changes
  - Fix bug in VM activity detection, preventing Docker Desktop from starting. Fixes [docker/for-win#2404](https://github.com/docker/for-win/issues/2404)
  - Fix bug in detection when Windows service is not running, and propose to restart the service.
  - Fix local DNS failing to resolve inside containers. Fixes [docker/for-win#2301](https://github.com/docker/for-win/issues/2301), [docker/for-win#2304](https://github.com/docker/for-win/issues/2304)
  - Fix Kubernetes status display after reset to factory default
  - Fix bug where `host.docker.internal` is not resolved in some cases. Fixes [docker/for-win#2402](https://github.com/docker/for-win/issues/2402)
  - Use a 1MB vhdx blocksize instead of the default 32MB. See [docker/for-win#244](https://github.com/docker/for-win/issues/244). Also see [Microsoft Best Practices for running Linux on Hyper-V](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/best-practices-for-running-linux-on-hyper-v)
  - Fix diagnostics in specific cases when the Windows service is not started.
  - Changed the samba default file permission to avoid an issue with too open rights. Fixes [docker/for-win#2170](https://github.com/docker/for-win/issues/2170)
  - On RS5 insider, fixed wrong detection of missing feature "Containers" requiring you to install the feature and then reboot.

### Docker Community Edition 18.06.0-ce-win72 2018-07-26

[Download](https://download.docker.com/win/stable/19098/Docker%20for%20Windows%20Installer.exe)

* New
  - Renewed signing certificate. The installer may show a Windows Defender popup until the renewed certificate is whitelisted. Click "More info" to view the app is published by "Docker Inc" and run it.

* Bug fixes and minor changes
  - Fix bug in automatic feature enabling if "Hyper-V" and "Containers" Windows features are not enabled already when starting Docker Desktop.

### Docker Community Edition 18.06.0-ce-win70 2018-07-25

[Download](https://download.docker.com/win/stable/19075/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.06.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.06.0-ce)
  - [Docker Machine 0.15.0](https://github.com/docker/machine/releases/tag/v0.15.0)
  - [Docker compose 1.22.0](https://github.com/docker/compose/releases/tag/1.22.0)
  - [LinuxKit v0.4](https://github.com/linuxkit/linuxkit/releases/tag/v0.4)
  - Linux Kernel 4.9.93 with CEPH, DRBD, RBD, MPLS_ROUTING and MPLS_IPTUNNEL enabled

* New
  - Kubernetes Support. You can now run a single-node Kubernetes cluster from the "Kubernetes" Pane in Docker for Windows settings and use kubectl commands as well as Docker commands. See https://docs.docker.com/docker-for-windows/kubernetes/

* Bug fixes and minor changes
  - AUFS storage driver is deprecated in Docker Desktop and AUFS support will be removed in the next major release. You can continue with AUFS in Docker Desktop 18.06.x, but you will need to reset the disk image (in Settings > Reset menu) before updating to the next major update. You can check documentation to [save images](https://docs.docker.com/engine/reference/commandline/save/#examples) and [backup volumes](https://docs.docker.com/storage/volumes/#backup-restore-or-migrate-data-volumes)
  - Fix bug which would in some cases cause virtual machine logs to be written to RAM rather than disk, and the virtual machine to hang.
  - Fix security issue with named pipe connection to docker service.
  - Fix VPNKit memory leak. Fixes [docker/for-win#2087](https://github.com/docker/for-win/issues/2087), [moby/vpnkit#371](https://github.com/moby/vpnkit/issues/371)
  - Fix restart issue when using Windows fast startup on latest 1709 Windows updates. Fixes [docker/for-win#1741](https://github.com/docker/for-win/issues/1741), [docker/for-win#1741](https://github.com/docker/for-win/issues/1741)
  - DNS name `host.docker.internal` can be used for host resolution from Windows Containers.  Fixes [docker/for-win#1976](https://github.com/docker/for-win/issues/1976)
  - Fix broken link in diagnostics window.
  - Added log rotation for docker-ce logs inside the virtual machine.
  - Changed smb permission to avoid issue when trying to manipulate files with different users in containers. Fixes [docker/for-win#2170](https://github.com/docker/for-win/issues/2170)

### Docker Community Edition 18.03.1-ce-win65 2018-04-30

[Download](https://download.docker.com/win/stable/17513/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.1-ce)
  - [Docker compose 1.21.1](https://github.com/docker/compose/releases/tag/1.21.1)
  - [Notary 0.6.1](https://github.com/docker/notary/releases/tag/v0.6.1)

* Bug fixes and minor changes
  - Fix startup failure when the HOME environment variable is already defined (typically started from the command line). Fixes [docker/for-win#1880](https://github.com/docker/for-win/issues/1880)
  - Fix startup failure due to incompatibility with other programs (like Razer Synapse 3). Fixes [docker/for-win#1723](https://github.com/docker/for-win/issues/1723)

### Docker Community Edition 18.03.1-ce-win64 2018-04-26

[Download](https://download.docker.com/win/stable/17438/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.1-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.1-ce)
  - [Docker compose 1.21.0](https://github.com/docker/compose/releases/tag/1.21.0)
  - [Notary 0.6.1](https://github.com/docker/notary/releases/tag/v0.6.1)

* Bug fixes and minor changes
  - Fix startup failure when the HOME environment variable is already defined (typically started from the command line). Fixes [docker/for-win#1880](https://github.com/docker/for-win/issues/1880)
  - Fix startup failure due to incompatibility with other programs (like Razer Synapse 3). Fixes [docker/for-win#1723](https://github.com/docker/for-win/issues/1723)

### Docker Community Edition 18.03.0-ce-win59 2018-03-26

[Download](https://download.docker.com/win/stable/16762/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 18.03.0-ce](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce)
  - [Docker Machine 0.14.0](https://github.com/docker/machine/releases/tag/v0.14.0)
  - [Docker compose 1.20.1](https://github.com/docker/compose/releases/tag/1.20.1)
  - [Notary 0.6.0](https://github.com/docker/notary/releases/tag/v0.6.0)
  - Linux Kernel 4.9.87
  - AUFS 20180312

* New
  - Virtual machine disk size can be changed in settings. Fixes [docker/for-win#105](https://github.com/docker/for-win/issues/105)
  - Virtual machine swap size can be changed in settings.
  - New menu item to restart Docker.
  - Support NFS Volume sharing. See [docker/for-win#1700](https://github.com/docker/for-win/issues/1700)
  - Allow to activate Windows Containers during installation (avoid virtual machine disk creation and virtual machine boot when working only on Windows containers). See [docker/for-win#217](https://github.com/docker/for-win/issues/217).
  - Experimental feature: LCOW containers can now be run next to Windows containers (on Windows RS3 build 16299 and later). Use `--platform=linux` in Windows container mode to run Linux Containers On Windows. Note that LCOW is experimental; it requires the daemon `experimental` option.

* Bug fixes and minor changes
  - Fix Windows Containers port forwarding on Windows 10 build 16299 post KB4074588. Fixes [docker/for-win#1707](https://github.com/docker/for-win/issues/1707), [docker/for-win#1737](https://github.com/docker/for-win/issues/1737)
  - Fix daemon not starting properly when setting TLS-related options.
  - DNS name `host.docker.internal` shoud be used for host resolution from containers. Older aliases (still valid) are deprecated in favor of this one. (See https://tools.ietf.org/html/draft-west-let-localhost-be-localhost-06).
  - Fix for the HTTP/S transparent proxy when using "localhost" names (for example, `host.docker.internal`). Fixes [docker/for-win#1130](https://github.com/docker/for-win/issues/1130)
  - Fix Linuxkit start on Windows RS4 Insider. Fixes [docker/for-win#1458](https://github.com/docker/for-win/issues/1458), [docker/for-win#1514](https://github.com/docker/for-win/issues/1514), [docker/for-win#1640](https://github.com/docker/for-win/issues/1640)
  - Fix risk of privilege escalation. (https://www.tenable.com/sc-report-templates/microsoft-windows-unquoted-service-path-vulnerability)
  - All users present in the docker-users group are now able to use Docker. Fixes [docker/for-win#1732](https://github.com/docker/for-win/issues/1732)
  - Migration of Docker Toolbox images is not proposed in Docker For Windows installer (still possible to [migrate Toolbox images manually](https://docs.docker.com/docker-for-windows/docker-toolbox/) ).
  - Better cleanup for Windows containers and images on reset/uninstall. Fixes [docker/for-win#1580](https://github.com/docker/for-win/issues/1580), [docker/for-win#1544](https://github.com/docker/for-win/issues/1544), [docker/for-win#191](https://github.com/docker/for-win/issues/191)
  - Desktop icon creation is optional in installer; do not recreate Desktop icon on upgrade (effective on next upgrade). Fixes [docker/for-win#246](https://github.com/docker/for-win/issues/246), [docker/for-win#925](https://github.com/docker/for-win/issues/925), [docker/for-win#1551](https://github.com/docker/for-win/issues/1551)

### Docker Community Edition 17.12.0-ce-win47 2018-01-12

[Download](https://download.docker.com/win/stable/15139/Docker%20for%20Windows%20Installer.exe)

* Bug fixes and minor changes
  - Fix linuxkit port-forwarder sometimes not being able to start. Fixes [docker/for-win#1506](https://github.com/docker/for-win/issues/1506)
  - Fix certificate management when connecting to a private registry. Fixes [docker/for-win#1512](https://github.com/docker/for-win/issues/1512)
  - Fix Mount compatibility when mounting drives with `-v //c/...`, now mounted in /host_mnt/c in the linuxkit virtual machine. Fixes [docker/for-win#1509](https://github.com/docker/for-win/issues/1509), [docker/for-win#1516](https://github.com/docker/for-win/issues/1516), [docker/for-win#1497](https://github.com/docker/for-win/issues/1497)
  - Fix icon displaying edge. Fixes [docker/for-win#1508](https://github.com/docker/for-win/issues/1508)

### Docker Community Edition 17.12.0-ce-win46 2018-01-09

[Download](https://download.docker.com/win/stable/15048/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 17.12.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce)
  - [Docker compose 1.18.0](https://github.com/docker/compose/releases/tag/1.18.0)
  - [Docker Machine 0.13.0](https://github.com/docker/machine/releases/tag/v0.13.0)
  - Linux Kernel 4.9.60

* New
  - Virtual machine entirely built with linuxkit
  - Add localhost port forwarder for Windows (thanks @simonferquel). Use Microsoft localhost port forwarder when it is available (insider build RS4).

* Bug fixes and minor changes
  - Display various component versions in About box.
  - Fix vpnkit issue when username has spaces. See [docker/for-win#1429](https://github.com/docker/for-win/issues/1429)
  - Diagnostic improvements to get VM logs before VM shutdown.
  - Fix installer check for not supported Windows `CoreCountrySpecific` Edition.
  - Fix a class of startup failures where the database fails to start. See [docker/for-win#498](https://github.com/docker/for-win/issues/498)
  - Links in update changelog now open the default browser instead of IE. (fixes [docker/for-win#1311](https://github.com/docker/for-win/issues/1311))

## Stable Releases of 2017

### Docker Community Edition 17.09.1-ce-win42 2017-12-11

[Download](https://download.docker.com/win/stable/14687/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 17.09.1-ce](https://github.com/docker/docker-ce/releases/tag/v17.09.1-ce)
  - [Docker compose 1.17.1](https://github.com/docker/compose/releases/tag/1.17.1)
  - [Docker Machine 0.13.0](https://github.com/docker/machine/releases/tag/v0.13.0)

* Bug fixes and minor changes
  - Fix bug during Windows fast-startup process. Fixes [for-win/#953](https://github.com/docker/for-win/issues/953)
  - Fix uninstaller issue (in some specific cases dockerd process was not killed properly)
  - Fix Net Promoter Score Gui bug. Fixes [for-win/#1277](https://github.com/docker/for-win/issues/1277)
  - Fix `docker.for.win.localhost` not working in proxy settings. Fixes [for-win/#1130](https://github.com/docker/for-win/issues/1130)
  - Increased timeout for virtual machine boot startup to 2 minutes.


### Docker Community Edition 17.09.0-ce-win33 2017-10-06

[Download](https://download.docker.com/win/stable/13620/Docker%20for%20Windows%20Installer.exe)

* Bug fixes
  - Fix Docker For Windows unable to start in some cases: removed use of libgmp sometimes causing the vpnkit process to die.

### Docker Community Edition 17.09.0-ce-win32 2017-10-02

[Download](https://download.docker.com/win/stable/13529/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 17.09.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce)
  - [Docker Compose 1.16.1](https://github.com/docker/compose/releases/tag/1.16.1)
  - [Docker Machine 0.12.2](https://github.com/docker/machine/releases/tag/v0.12.2)
  - [Docker Credential Helpers 0.6.0](https://github.com/docker/docker-credential-helpers/releases/tag/v0.6.0)
  - Linux Kernel 4.9.49
  - AUFS 20170911

* New
  - Windows Docker daemon is now started as service for better lifecycle management
  - Store Linux daemon configuration in ~\.docker\daemon.json instead of settings file
  - Store Windows daemon configuration in C:\ProgramData\Docker\config\daemon.json instead of settings file
  - VPNKit: add support for ping!
  - VPNKit: add slirp/port-max-idle-timeout to allow the timeout to be adjusted or even disabled
  - VPNKit: bridge mode is default everywhere now
  - Add `Skip This Version` button in update window

* Security fixes
  - VPNKit: security fix to reduce the risk of DNS cache poisoning attack (reported by Hannes Mehnert https://hannes.nqsb.io/)

* Bug fixes and minor changes
  - Kernel: Enable TASK_XACCT and TASK_IO_ACCOUNTING
  - Rotate logs in the virtual machine more often (docker/for-win#244)
  - Reset to default stops all engines and removes settings including all daemon.json files
  - Better backend service checks (related to https://github.com/docker/for-win/issues/953)
  - Fix auto updates checkbox, no need to restart the application
  - Fix check for updates menu when auto updates was disable
  - VPNKit: do not block startup when ICMP permission is denied. (Fixes docker/for-win#1036, docker/for-win#1035, docker/for-win#1040)
  - VPNKit: change protocol to support error messages reported back from the server
  - VPNKit: fix a bug which causes a socket to leak if the corresponding TCP connection is idle
    for more than 5 minutes (related to [docker/for-mac#1374](https://github.com/docker/for-mac/issues/1374))
  - VPNKit: improve the logging around the UNIX domain socket connections
  - VPNKit: automatically trim whitespace from int or bool database keys
  - Do not move credentials into the credential store at startup

### Docker Community Edition 17.06.2-ce-win27 2017-09-06

[Download](https://download.docker.com/win/stable/13194/Docker%20for%20Windows%20Installer.exe)

* Upgrades
  - [Docker 17.06.2-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.2-ce)
  - [Docker Machine 0.12.2](https://github.com/docker/machine/releases/tag/v0.12.2)

### Docker Community Edition 17.06.1-ce-rc1-win24 2017-08-24

[Download](https://download.docker.com/win/stable/13025/Docker%20for%20Windows%20Installer.exe)

**Upgrades**

- [Docker 17.06.1-ce-rc1](https://github.com/docker/docker-ce/releases/tag/v17.06.1-ce-rc1)
- Linux Kernel 4.9.36
- AUFS 20170703

**Bug fixes and minor**

- Fix locked container id file (Fixes [docker/for-win#818](https://github.com/docker/for-win/issues/818))
- Avoid expanding variables in PATH env variable (Fixes [docker/for-win#859](https://github.com/docker/for-win/issues/859))

### Docker Community Edition 17.06.0-ce-win18 2017-06-28

[Download](https://download.docker.com/win/stable/12627/Docker%20for%20Windows%20Installer.exe)

**Upgrades**

- [Docker 17.06.0-ce](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce)
- [Docker Credential Helpers 0.5.2](https://github.com/docker/docker-credential-helpers/releases/tag/v0.5.2)
- [Docker Machine 0.12.0](https://github.com/docker/machine/releases/tag/v0.12.0)
- [Docker compose 1.14.0](https://github.com/docker/compose/releases/tag/1.14.0)
- Linux Kernel 4.9.31

**New**

- Windows Server 2016 support
- Windows 10586 is marked as deprecated; it is not supported going forward in stable releases
- Integration with Docker Cloud, with the ability to control remote Swarms from the local command line interface (CLI) and view your repositories
- Unified login between the Docker CLI and Docker Hub, Docker Cloud.
- Sharing a drive can be done on demand, the first time a mount is requested
- Add an experimental DNS name for the host: docker.for.win.localhost
- Support for client (i.e. "login") certificates for authenticating registry access (fixes [docker/for-win#569](https://github.com/docker/for-win/issues/569))
- New installer experience

**Bug fixes and minor changes**

- Fixed group access check for users logged in with Active Directory (fixes [docker/for-win#785](https://github.com/docker/for-win/issues/785))
- Checked environment variables and add some warnings in logs if they can cause docker to fail
- Many processes that were running in admin mode now run within the user identity
- Cloud federation command lines now open in the user home directory
- Named pipes are now created with more constrained security descriptors to improve security
- Security fix: Users must be part of the specific group "docker-users" to run Docker for Windows
- Reset to default / uninstall also reset Docker cli settings and logout user from Docker Cloud and registries
- Detected a bitlocker policy preventing windows containers to work
- Fixed an issue on filesharing when explicitly disabled on vmswitch interface
- Fixed virtual machine not starting when the machine had a very long name
- Fixed a bug where Windows daemon.json file was not written (fixes [docker/for-win#670](https://github.com/docker/for-win/issues/670))
- Added patches to the kernel to fix VMBus crash
- Named pipe client connection should not trigger dead locks on `docker run` with data in stdin any more
- Buffered data should be treated correctly when docker client requests are upgraded to raw streams

### Docker Community Edition 17.03.1-ce-win12  2017-05-12

[Download](https://download.docker.com/win/stable/12058/Docker%20for%20Windows%20Installer.exe)

**Upgrades**

- Security fix for CVE-2017-7308

### Docker Community Edition 17.03.0, 2017-03-02

[Download](https://download.docker.com/win/stable/10743/Docker%20for%20Windows%20Installer.exe)

**New**

- Renamed to Docker Community Edition
- Integration with Docker Cloud: control remote Swarms from the local CLI and view your repositories. This feature is going to be rolled out to all users
progressively

**Upgrades**

- [Docker 17.03.0-ce](https://github.com/docker/docker/releases/tag/v17.03.0-ce)
- [Docker Compose 1.11.2](https://github.com/docker/compose/releases/tag/1.11.2)
- [Docker Machine 0.10.0](https://github.com/docker/machine/releases/tag/v0.10.0)
- Linux kernel 4.9.12

**Bug fixes and minor changes**

- Matched Hyper-V Integration Services by ID, not name
- Don't consume 100% CPU when the service is stopped
- Log the diagnostic id when uploading
- Improved firewall handling: stop listing the rules since it can take a lot of time
- Don't rollback to the previous engine when the desired engine fails to start
- Don't use port 4222 inside the Linux virtual machine
- Fixed startup error of ObjectNotFound in Set-VMFirmware
- Added detailed logs when firewall is configured
- Added a link to the Experimental Features documentation
- Fixed the copyright in the About dialog
- VPNKit: fix unmarshalling of DNS packets containing pointers to pointers to labels
- VPNKit: set the Recursion Available bit on DNS responses from the cache
- VPNKit: Avoid diagnostics to capture too much data
- VPNKit: fix a source of occasional packet loss (truncation) on the virtual ethernet link
- Fixed negotiation of TimeSync protocol version (through kernel update)

### Docker for Windows 1.13.1, 2017-02-09

[Download](https://download.docker.com/win/stable/1.13.1.10072/InstallDocker.msi)

**Upgrades**

- [Docker 1.13.1](https://github.com/docker/docker/releases/tag/v1.13.1)
- [Docker Compose 1.11.1](https://github.com/docker/compose/releases/tag/1.11.1)
- Linux kernel 4.9.8

**Bug fixes and minor changes**

- Add link to experimental features
- New 1.13 cancellable operations should now be properly handled by the Docker for desktop
- Various typos fixes
- Fix in Hyper-V VM setup (should fix `ObjectNotFound` errors)

### Docker for Windows 1.13.0, 2017-01-19

[Download](https://download.docker.com/win/stable/1.13.0.9795/InstallDocker.msi)

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
- Virtual machine's time synchronization is forced after the host wakes from sleep mode
- Docker Experimental mode can be toggled

**Bug fixes and minor changes**

- Improved Proxy UI
- Improvements to Logging and Diagnostics
- About Box is now copy/paste enabled
- Improvements in drive sharing code
- Optimized boot process
- Trend Micro Office Scan made the API proxy think no drive was shared
- Show a link to the virtualization documentation
- Always remove the disk on factory reset
- VPNKit: Improved diagnostics
- VPNKit: Forwarded UDP datagrams have correct source port numbers
- VPNKit: If one request fails, allow other concurrent requests to succeed.
  For example this allows IPv4 servers to work even if IPv6 is broken.
- VPNKit: Fix bug which could cause the connection tracking to
  underestimate the number of active connections
- VPNKit: add a local cache of DNS responses

## Stable Releases of 2016

### Docker for Windows 1.12.5, 2016-12-20

[Download](https://download.docker.com/win/stable/1.12.5.9503/InstallDocker.msi)

**Upgrades**

- Docker 1.12.5
- Docker Compose 1.9.0

### Skipped Docker for Windows 1.12.4

We did not distribute a 1.12.4 stable release

### Docker for Windows 1.12.3, 2016-11-09

[Download](https://download.docker.com/win/stable/1.12.3.8488/InstallDocker.msi)

**New**

- Restore the virtual machine's configuration after user changes

- Detect firewall configurations that might block file sharing

- Send more GUI usage statistics to help us improve the product

- The path to HyperV disks is not hardcoded anymore, making the Toolbox import work with non-standard paths

- Verify that all HyperV features are enabled

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

- Made sure we didn't use an older Nlog library from the GAC

- Fixed a password escaping regression

- Supports writing large values to the database, specially for trusted CAs

- Preserves the Powershell stacktraces

- Writes OS and Application versions at the top of each log file

- Don't recreate the virtual machine if only the DNS server is set

- The uninstaller now kills the service if it failed to stop it properly

- Improved debug information

**Networking**

- VpnKit is now restarted if it stops

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

### Docker for Windows 1.12.1, 2016-09-16

[Download](https://download.docker.com/win/stable/1.12.1.7135/InstallDocker.msi)

>**Important Note**:
>
> The auto-update function in Beta 21 cannot install this update. To install the latest beta manually if you are still on Beta 21, download the installer here:

> [https://download.docker.com/win/beta/InstallDocker.msi](https://download.docker.com/win/beta/InstallDocker.msi)

> This problem is fixed as of Beta 23 for subsequent auto-updates.

**New**

* To support trusted registry transparently, all trusted CAs (root or intermediate) on the Windows host are automatically copied to Moby

* `Reset Credentials` also unshares the shared drives

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

* Docker stops asking to import from the Toolbox after an upgrade

* Docker can now import from the Toolbox just after HyperV is activated

* Added more debug information to the diagnostics

* Sending anonymous statistics doesn't hang anymore when Mixpanel is not available

* Support newlines in release notes

* Improve the error message when the Docker daemon is not responding

* The configuration database is now stored in memory

* Preserve the stacktrace of PowerShell errors

* Display service stacktrace in error windows

**Networking**

* Improve name servers discovery
* VpnKit supports search domains
* VpnKit is now compiled with OCaml 4.03 rather than 4.02.3

**File sharing**

* Set `cifs` version to 3.02

* VnpKit: reduce the number of sockets used by UDP NAT, reduce the probability

* `slirp`: reduce the number of sockets used by UDP NAT, reduce the probability that NAT rules time out earlier than expected

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

* Docker automatically disables lingering net adapters. The only way to remove them is manually using `devmgmt.msc` as documented in [Remove stale network adapters](troubleshoot.md#4-remove-stale-network-adapters) under [Networking issues](troubleshoot.md#networking-issues) in Troubleshooting.

### Docker for Windows 1.12.0, 2016-07-28

[Download](https://download.docker.com/win/stable/1.12.0.5968/InstallDocker.msi)

* First stable release

**Components**

* Docker 1.12.0
* Docker Machine 0.8.0
* Docker Compose 1.8.0
