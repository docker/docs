---
description: Release notes for Docker Desktop for Mac, Linux, and Windows
keywords: Docker desktop, release notes, linux, mac, windows
title: Docker Desktop release notes
toc_max: 2
aliases:
- /docker-for-mac/release-notes/
- /docker-for-mac/edge-release-notes/
- /desktop/mac/release-notes/
- /docker-for-windows/edge-release-notes/
- /docker-for-windows/release-notes/
- /desktop/windows/release-notes/
- /desktop/linux/release-notes/
- /mackit/release-notes/
---

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

> **Note**
>
> The information below is applicable to all platforms, unless stated otherwise.

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1) to see what's coming next.

For frequently asked questions about Docker Desktop releases, see [FAQs](faqs/releases.md).

## 4.26.1

{{< release-date date="2023-12-14" >}}

{{< desktop-install all=true version="4.26.1" build_path="/131620/" >}}

### Bug fixes and enhancements

#### For all platforms

- Updated feedback links inside Docker Desktop to ensure they continue to work correctly

#### For Windows

- Switch the CLI binaries to a version compatible with older versions of glibc, such as used in Ubuntu 20.04 fixes [docker/for-win#13824](https://github.com/docker/for-win/issues/13824)

## 4.26.0

{{< release-date date="2023-12-04" >}}

{{< desktop-install all=true version="4.26.0" build_path="/130397/" >}}

### New

- Administrators can now control access to beta and experimental features in the **Features in development** tab with [Settings Management](hardened-desktop/settings-management/configure.md).
- Introduced four new version update states in the footer.
- `docker init` (Beta) now supports PHP with Apache + Composer.
- The [**Builds** view](use-desktop/builds.md) is now GA. You can now inspect builds, troubleshoot errors, and optimize build speed.

### Upgrades

- [Compose v2.23.3](https://github.com/docker/compose/releases/tag/v2.23.3)
- [Docker Scout CLI v1.2.0](https://github.com/docker/scout-cli/releases/tag/v1.2.0).
- [Buildx v0.12.0](https://github.com/docker/buildx/releases/tag/v0.12.0)
- [Wasm](../desktop/wasm/_index.md) runtimes:
  - wasmtime, wasmedge and wasmer `v0.3.1`.
  - lunatic, slight, spin, and wws `v0.10.0`.
  - Wasmtime is now based on wasmtime `v14.0` and supports wasi preview-2 components
  - Wasmedge is now based on WasmEdge `v0.13.5`
  - Spin is now based on Spin `v2.0.1`
  - wws is now based on wws `v1.7.0`
- [Docker Engine v24.0.7](https://docs.docker.com/engine/release-notes/24.0/#2407)
- [Containerd v1.6.25](https://github.com/containerd/containerd/releases/tag/v1.6.25)
- [runc v1.1.10](https://github.com/opencontainers/runc/releases/tag/v1.1.10)

### Bug fixes and enhancements

#### For all platforms

- You can now provide feedback from the commandline by using `docker feedback`.
- Improved the text and position of the startup options in the **General** settings tab.
- Redesigned the dashboard's header bar, added links to other Docker resources, improved display of account information.
- Fixed a bug  where enabling the containerd image store and Wasm simultaneously would not enable Wasm.
- containerd integration:
  - Fixed `docker push/pull` authentication not being sent to non-DockerHub registries in cases where `ServerAddress` is not provided.
  - Fixed `docker history` reporting wrong IDs and tags.
  - Fixed `docker tag` not preserving internal metadata.
  - Fixed `docker commit` when the daemon configured with `--userns-remap`.
  - Fixed `docker image list` to show real image creation date.
  - Added support for `-a` flag to `docker pull` (pull all remote repository tags).
  - Added support for `--group-add` flag to `docker run` (append extra groups).
  - Adjusted some errors reported by `docker push/pull`.
- Docker Init:
  - Improved cross-compilation in Dockerfiles for Golang and Rust.
  - Improved caching in Dockerfile for ASP.NET Core.
- Docker Desktop now gives more detailed information about pending updates in the dashboard footer.
- Fixed a bug in Enhanced Container Isolation mode where `docker run --init` was failing.
- Fixed a bug where a notification prompting the user to download a new version of Docker Desktop remained visible after the user started downloading the new version.
- Added a notification that indicates when Docker Desktop is installing a new version.
- Fixed a bug where the cursor changed to a pointer when the user hovered over a notification that has no call to action.

#### For Mac

- Fixed an issue where Rosetta would not work with PHP. Fixes [docker/for-mac#6773](https://github.com/docker/for-mac/issues/6773)  and [docker/for-mac#7037](https://github.com/docker/for-mac/issues/7037).
- Fixed several issues related to Rosetta not working. Fixed [[docker/for-mac#6973](https://github.com/docker/for-mac/issues/6973), [[docker/for-mac#7009](https://github.com/docker/for-mac/issues/7009), [[docker/for-mac#7068](https://github.com/docker/for-mac/issues/7068) and [[docker/for-mac#7075](https://github.com/docker/for-mac/issues/7075)
- Improved the performance of NodeJS under Rosetta.
- Fixed the **Unable to open /proc/self/exe** Rosetta errors.
- Fixed a bug were the setting **Start Docker Desktop when you sign in** would not work. Fixes [docker/for-mac#7052](https://github.com/docker/for-mac/issues/7052).
- You can now enable the use of Kernel networking path for UDP through the UI. Fixes [docker/for-mac#7008](https://github.com/docker/for-mac/issues/7008).
- Fixed a regression where the `uninstall` CLI tool was missing.
- Addressed an issue which caused Docker Desktop to become unresponsive when analytics were disabled with Settings Management. 

#### For Windows

- Added support for WSL mirrored mode networking (requires WSL `v2.0.4` and up).
- Added missing signatures on DLL and VBS files.

### Known issues

#### For Windows

- Docker CLI doesn’t work when using WSL 2 integration on an older Linux distribution (for example, Ubuntu 20.04) which uses a `glibc` version older than `2.32`. This will be fixed in future releases. See [docker/for-win#13824](https://github.com/docker/for-win/issues/13824).
  
## 4.25.2

{{< release-date date="2023-11-21" >}}

{{< desktop-install all=true version="4.25.2" build_path="/129061/" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed a bug where a blank UI would appear after submitting a response in the **Welcome Survey**.

#### For Windows

- Fixed a bug where Docker Desktop on WSL 2 would shut down dockerd unexpectedly when idle. Fixes [docker/for-win#13789](https://github.com/docker/for-win/issues/13789)

## 4.25.1

{{< release-date date="2023-11-13" >}}

{{< desktop-install all=true version="4.25.1" build_path="/128006/" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed a regression in 4.25 where Docker would not start if the swap file was corrupt. Corrupt swap files will be re-created on next boot.
- Fixed a bug when swap is disabled. Fixes [docker/for-mac#7045](https://github.com/docker/for-mac/issues/7045), [docker/for-mac#7044](https://github.com/docker/for-mac/issues/7044) and [docker/for-win#13779](https://github.com/docker/for-win/issues/13779).
- The `sysctl vm.max_map_count` is now set to 262144. See [docker/for-mac#7047](https://github.com/docker/for-mac/issues/7047)

#### For Windows

- Fixed an issue where **Switch to Windows Containers** would not appear on the tray menu for some users. See [docker/for-win#13761](https://github.com/docker/for-win/issues/13761).
- Fixed a bug where the WSL integration would not work for users using a shell other than `sh`. See [docker/for-win#13764](https://github.com/docker/for-win/issues/13764).
- Re-added `DockerCli.exe`.

## 4.25.0

{{< release-date date="2023-10-26" >}}

{{< desktop-install all=true version="4.25.0" build_path="/126437/" >}}

### New

- Rosetta is now Generally Available for all users on macOS 13 or later. It provides faster emulation of Intel-based images on Apple Silicon. To use Rosetta, see [Settings](settings/mac.md). Rosetta is enabled by default on macOS 14.1 and later.
- Docker Desktop now detects if a WSL version is out of date. If an out dated version of WSL is detected, you can allow Docker Desktop to automatically update the installation or you can manually update WSL outside of Docker Desktop.
- New installations of Docker Desktop for Windows now require a Windows version of 19044 or later.
- Administrators now have the ability to control Docker Scout image analysis  in [Settings Management](hardened-desktop/settings-management/configure.md).

### Upgrades

- [Compose v2.23.0](https://github.com/docker/compose/releases/tag/v2.23.0)
- [Docker Scout CLI v1.0.9](https://github.com/docker/scout-cli/releases/tag/v1.0.9).
- [Kubernetes v1.28.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.28.2)
  - [cri-dockerd v0.3.4](https://github.com/Mirantis/cri-dockerd/releases/tag/v0.3.4)
  - [CNI plugins v1.3.0](https://github.com/containernetworking/plugins/releases/tag/v1.3.0)
  - [cri-tools v1.28.0](https://github.com/kubernetes-sigs/cri-tools/releases/tag/v1.28.0)

### Bug fixes and enhancements

#### For all platforms

- Fixed a spacing problem in the `Accept License` pop-up.
- Fixed a bug where the **Notifications drawer** changed size when navigating between **Notifications list** and **Notification details** view.
- containerd integration:
  - `docker push` now supports `Layer already exists` and `Mounted from` progress statuses.
  - `docker save` is now able to export images from all tags of the repository.
  - Hide push upload progress of manifests, configs and indexes (small json blobs) to match the original push behavior.
  - Fixed `docker diff` containing extra differences.
  - Fixed `docker history` not showing intermediate image IDs for images built with the classic builder.
  - Fixed `docker load` not being able to load images from compressed tar archives.
  - Fixed registry mirrors not working.
  - Fixed `docker diff` not working correctly when called multiple times concurrently for the same container.
  - Fixed `docker push` not reusing layers when pushing layers to different repositories on the same registry.
- Docker Init:
  - Fixed outdated links to Docker documentation included in generated files
  - Add support for ASP.NET Core 8 (in addition to 6 and 7)
- Fixed a bug that caused a failure when installing Wasm shims.
- Fixed a bug where Docker Desktop exits the [Resource Saver mode](https://docs.docker.com/desktop/use-desktop/resource-saver/) every 15 minutes, or, if the timer is set above 15 minutes, the resource saver mode never kicks in.
- Promoted the **Enable background SBOM indexing** option to **General settings**.

#### For Mac

- Minimum OS version to install or update Docker Desktop on macOS is now macOS Monterey (version 12) or later.
- Enhanced error messaging when an update cannot be completed if the user doesn't match the owner of `Docker.app`. Fixes [docker/for-mac#7000](https://github.com/docker/for-mac/issues/7000).
- Fixed a bug where **Re-apply configuration** might not work when `/var/run/docker.sock` is mis-configured.
- Docker Desktop doesn't overwrite `ECRCredentialHelper` if already present in `/usr/local/bin`.

#### For Windows

- Fixed an issue where **Switch to Windows Containers** would show in the tray menu on Windows Home Editions. Fixes [docker/for-win#13715](https://github.com/docker/for-win/issues/13715)

#### For Linux

- Fixed a bug in `docker login`. Fixes  [docker/docker-credential-helpers#299](https://github.com/docker/docker-credential-helpers/issues/299)

### Known Issues

#### For Mac

- Upgrading to MacOS 14 can cause Docker Desktop to also update to a latest version even if the auto update option is disabled.
- Uninstalling Docker Desktop from the command line is not available. As a workaround, you can [uninstall Docker Desktop from the Dashboard](https://docs.docker.com/desktop/uninstall/).

#### For Windows

- **Switch to Windows containers** option in the tray menu may not show up on Windows. As a workaround, edit the [`settings.json` file](https://docs.docker.com/desktop/settings/windows/) and set `"displaySwitchWinLinContainers": true`.

#### For all platforms
- Docker operations, such as pulling images or logging in, fail with 'connection refused' or 'timeout' errors if the Swap file size is set to 0MB. As a workaround, configure the swap file size to a non-zero value in the **Resources** tab in **Settings**. 

## 4.24.2

{{< release-date date="2023-10-12" >}}

{{< desktop-install all=true version="4.24.2" build_path="/124339/" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed a bug where Docker Desktop would send multiple requests to `notify.bugsnag.com`. Fixes [docker/for-win#13722](https://github.com/docker/for-win/issues/13722).
- Fixed a performance regression for PyTorch.

## 4.24.1

{{< release-date date="2023-10-04" >}}

{{< desktop-install win=true version="4.24.1" build_path="/123237/" >}}

### Bug fixes and enhancements

#### For Windows

- Fixed a bug on Docker Desktop for Windows where the Docker Dashboard wouldn't display container logs correctly. Fixes [docker/for-win#13714](https://github.com/docker/for-win/issues/13714).

## 4.24.0

{{< release-date date="2023-09-28" >}}

{{< desktop-install all=true version="4.24.0" build_path="/122432/" >}}

### New

- The new Notification center is now available to all users so you can be notified of new releases, installation progress updates, and more. Select the bell icon in the bottom-right corner of the Docker Dashboard to access the notification center.
- Compose Watch is now available to all users. For more information, see [Use Compose Watch](../compose/file-watch.md).
- Resource Saver is now available to all users and is enabled by default. To configure this feature, navigate to the **Resources** tab in **Settings**. For more information see [Docker Desktop's Resource Saver mode](use-desktop/resource-saver.md).
- You can now view and manage the Docker Engine state, with pause, stop, and resume, directly from the Docker Dashboard.

### Upgrades

- [Compose v2.22.0](https://github.com/docker/compose/releases/tag/v2.22.0)
- [Go 1.21.1](https://github.com/golang/go/releases/tag/go1.21.1)
- [Wasm](../desktop/wasm.md) runtimes:
  - wasmtime, wasmedge `v0.2.0`.
  - lunatic, slight, spin, and wws`v0.9.1`.
  - Added wasmer wasm shims.


### Bug fixes and enhancements

#### For all platforms

- Docker Init:
  - Fixed an issue formatting Dockerfile file paths for ASP.NET projects on Windows.
  - Improved performance on language detection for large directories with lots of files.
- Added a timeout to polling for resource usage stats used by the **Containers** view. Fixes [docker/for-mac#6962](https://github.com/docker/for-mac/issues/6962).
- containerd integration:
  - Implemented push/pull/save image events.
  - Implemented pulling legacy schema1 images.
  - Implemented `docker push --all-tags`.
  - Implemented counting containers using a specific image (visible for example in `docker system df -v`).
  - Validated pulled image names are not reserved.
  - Handle `userns-remap` daemon setting.
  - Fixed legacy builder build errors when multiple COPY/ADD instructions are used.
  - Fixed `docker load` causing pool corruption which could some subsequent image related operations.
  - Fixed not being able to reference images via truncated digest with a `sha256:` prefix.
  - Fixed `docker images` (without `--all`) showing intermediate layers (created by the legacy classic builder).
  - Fixed `docker diff` containing extra differences.
  - Changed `docker pull` output to match the output with containerd integration disabled.
- Fixed a grammatical error in Kubernetes status message. See [docker/for-mac#6971](https://github.com/docker/for-mac/issues/6971).
- Docker containers now use all host CPU cores by default.
- Improved inter-process security in dashboard UI.

#### For Mac

- Fixed a kernel panic on Apple silicon Macs with macOS version below 12.5. Fixes [docker/for-mac#6975](https://github.com/docker/for-mac/issues/6975).
- Fixed a bug where Docker Desktop failed to start if invalid directories were included in `filesharingDirectories`. Fixes [docker/for-mac#6980](https://github.com/docker/for-mac/issues/6980).
- Fixed a bug where installer is creating root-owned directories. Fixes [docker/for-mac#6984](https://github.com/docker/for-mac/issues/6984).
- Fixed a bug where installer is failing on setting up the docker socket when missing `/Library/LaunchDaemons`. Fixes [docker/for-mac#6967](https://github.com/docker/for-mac/issues/6967).
- Fixed a permission denied error when binding a privileged port to a non-localhost IP on macOS. Fixes [docker/for-mac#697](https://github.com/docker/for-mac/issues/6977).
- Fixed a resource leak introduced in 4.23. Related to [docker/for-mac#6953](https://github.com/docker/for-mac/issues/6953).

#### For Windows

- Fixed a bug where a "Docker Desktop service not running" popup appeared when service is already running. See [docker/for-win#13679](https://github.com/docker/for-win/issues/13679).
- Fixed a bug that caused Docker Desktop fail to start on Windows hosts. Fixes [docker/for-win#13662](https://github.com/docker/for-win/issues/13662).
- Modified the Docker Desktop resource saver feature to skip reducing kernel memory on WSL when no containers are running, as this was causing timeouts in some cases. Instead, users are encouraged to enable "autoMemoryReclaim" on WSL directly via the .wslconfig file (available since WSL 1.3.10).

### Known issues

#### For Mac

- Creating a container with the port 53 fails with the error address `already in use`. As a workaround, deactivate network acceleration by adding `"kernelForUDP": false`, in the `settings.json` file located at `~/Library/Group Containers/group.com.docker/settings.json`.
## 4.23.0

{{< release-date date="2023-09-11" >}}

{{< desktop-install all=true version="4.23.0" build_path="/120376/" >}}

### Upgrades
- [Compose v2.21.0](https://github.com/docker/compose/releases/tag/v2.21.0)
- [Docker Engine v24.0.6](https://docs.docker.com/engine/release-notes/24.0/#2406)
- [Docker Scout CLI v0.24.1](https://github.com/docker/scout-cli/releases/tag/v0.24.1).
- [Wasm](../desktop/wasm.md) runtimes:
  - wasmtime, wasmedge revision `d0a1a1cd`.
  - slight and spin wasm `v0.9.0`.

### New

- Added support for new Wasm runtimes: wws and lunatic.
- [`docker init`](../engine/reference/commandline/init.md) now supports ASP.NET
- Increased performance of exposed ports on macOS, for example with `docker run -p`.

### Bug fixes and enhancements

#### For all platforms

- With [Docker Scout](../scout/_index.md), you can now:
  - Manage temporary and cached files with `docker scout cache`.
  - Manage environments with `docker scout environment`.
  - Configure the default organization with `docker scout config`.
  - List packages of an image with their vulnerabilities with `docker scout cves --format only-packages`.
  - Enroll an organization with Docker scout with `docker scout enroll`.
  - Stop, analyze, and compare local file systems with `docker scout cves --type fs`.
- Fixed a bug where `docker stats` would hang when Docker Desktop was in Resource Saver mode.
- Fixed a bug where turning off experimental features via **Settings** in the Docker Dashboard would not fully turn off Resource Saver mode.
- Fixed a bug where the **Containers list** action button was clipped.
- containerd image store:
  - Fixed `failed to read config content` error when interacting with some images.
  - Fixed building Dockerfiles with `FROM scratch` instruction when using the legacy classic builder (`DOCKER_BUILDKIT=0`).
  - Fixed `mismatched image rootfs errors` when building images with legacy classic builder (`DOCKER_BUILDKIT=0`).
  - Fixed `ONBUILD` and `MAINTAINER` Dockerfile instruction
  - Fixed healthchecks.

#### For Mac

- All users on macOS 12.5 or greater now have VirtioFS turned on by default. You can revert this in **Settings** in the **General** tab.
- Improved single-stream TCP throughput.
- Reinstated the health check for macOS that notifies you if there has been a change on your system which might cause problems running Docker binaries.

#### For Linux

- Fixed a bug where the GUI is killed when opening the Docker Desktop app twice. See [docker/desktop-linux#148](https://github.com/docker/desktop-linux/issues/148).

#### For Windows

- Fixed a bug where non-admin users would get prompted for credentials when switching to Windows Containers or after disabling WSL and switching to the Hyper-V engine.
  This issue would occur after an OS restart, or on a cold start of Docker Desktop.

### Security

#### For all platforms

- Fixed [CVE-2023-5165](https://www.cve.org/cverecord?id=CVE-2023-5165) which allows Enhanced Container Isolation bypass via debug shell. The affected functionality is available for Docker Business customers only and assumes an environment where users are not granted local root or Administrator privileges.
- Fixed [CVE-2023-5166](https://www.cve.org/cverecord?id=CVE-2023-5166) which allows Access Token theft via a crafted extension icon URL.

### Known Issues

- Binding a priviledged port on Docker Desktop does not work on macOS. As a workaround you can expose the port on all interfaces (using `0.0.0.0`) or using localhost (using `127.0.0.1`).

## 4.22.1

{{< release-date date="2023-08-24" >}}

{{< desktop-install all=true version="4.22.1" build_path="/118664/" >}}

### Bug fixes and enhancements

#### For all platforms

- Mitigated several issues impacting Docker Desktop startup and Resource Saver mode. [docker/for-mac#6933](https://github.com/docker/for-mac/issues/6933)

#### For Windows

- Fixed `Clean / Purge data` troubleshoot option on Windows. [docker/for-win#13630](https://github.com/docker/for-win/issues/13630)

## 4.22.0

{{< release-date date="2023-08-03" >}}

{{< desktop-install all=true version="4.22.0" build_path="/117440/" >}}

### Upgrades

- [Buildx v0.11.2](https://github.com/docker/buildx/releases/tag/v0.11.2)
- [Compose v2.20.2](https://github.com/docker/compose/releases/tag/v2.20.2)
- [Docker Engine v24.0.5](https://docs.docker.com/engine/release-notes/24.0/#2405)

> **Note**
>
> In this release, the bundled Docker Compose and Buildx binaries show a different version string. This relates to our efforts to test new features without causing backwards compatibility issues.
>
> For example, `docker buildx version` outputs `buildx v0.11.2-desktop.1`.

### New

- [Resource Usage](use-desktop/container.md) has moved from experimental to GA.
- You can now split large Compose projects into multiple sub-projects with [`include`](../compose/multiple-compose-files/include.md).

### Bug fixes and enhancements

#### For all platforms

- [Settings Management](hardened-desktop/settings-management/index.md) now lets you turn off Docker Extensions for your organisation.
- Fixed a bug where turning on Kubernetes from the UI failed when the system was paused.
- Fixed a bug where turning on Wasm from the UI failed when the system was paused.
- Bind mounts are now shown when you [inspect a container](use-desktop/container.md).
- You can now download Wasm runtimes when the containerd image store is enabled.
- With [Quick Search](../desktop/use-desktop/_index.md/#quick-search), you can now:
  - Find any container or Compose app residing on your local system. In
    addition, you can access environment variables and perform essential actions
    such as starting, stopping, or deleting containers.
  - Find public Docker Hub images, local images, or images from remote
    repositories.
  - Discover more about specific extensions and install them.
  - Navigate through your volumes and gain insights about the associated
    containers.
  - Search and access Docker's documentation.

#### For Mac

- Fixed a bug that prevented Docker Desktop from starting. [docker/for-mac#6890](https://github.com/docker/for-mac/issues/6890)
- Resource Saver is now available on Mac. It optimises Docker Desktop's usage of your system resources when no containers are running. To access this feature, make sure you have [turned on access to experimental features](settings/mac.md) in settings.

#### For Windows

- Fixed a bug where the self-diagnose tool showed a false-positive failure when vpnkit is expected to be not running. Fixes [docker/for-win#13479](https://github.com/docker/for-win/issues/13479).
- Fixed a bug where an invalid regular expression in the search bar caused an error. Fixes [docker/for-win#13592](https://github.com/docker/for-win/issues/13592).
- Resource Saver is now available on Windows Hyper-V. It optimises Docker Desktop's usage of your system resources when no containers are running. To access this feature, make sure you have [turned on access to experimental features](settings/windows.md) in settings.

## 4.21.1

{{< release-date date="2023-07-03" >}}

{{< desktop-install all=true version="4.21.1" build_path="/114176/" >}}

#### For all platforms

- Fixed connection leak for Docker contexts using SSH ([docker/for-mac#6834](https://github.com/docker/for-mac/issues/6834) and [docker/for-win#13564](https://github.com/docker/for-win/issues/13564))

#### For Mac

- Removed configuration health check for further investigation and addressing specific setups.

## 4.21.0

{{< release-date date="2023-06-29" >}}

{{< desktop-install all=true version="4.21.0" build_path="/113844/" >}}

### New

- Added support for new Wasm runtimes: slight, spin, and wasmtime. Users can download Wasm runtimes on demand when the containerd image store is enabled.
- Added Rust server support to Docker init.
- Beta release of the [**Builds** view](use-desktop/builds.md) that lets you inspect builds and manage builders. This can be found in the **Features in Development** tab in **Settings**.

### Upgrades

- [Buildx v0.11.0](https://github.com/docker/buildx/releases/tag/v0.11.0)
- [Compose v2.19.0](https://github.com/docker/compose/releases/tag/v2.19.0)
- [Kubernetes v1.27.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.27.2)
- [cri-tools v1.27.0](https://github.com/kubernetes-sigs/cri-tools/releases/tag/v1.27.0)
- [cri-dockerd v0.3.2](https://github.com/Mirantis/cri-dockerd/releases/tag/v0.3.2)
- [coredns v1.10.1](https://github.com/coredns/coredns/releases/tag/v1.10.1)
- [cni v1.2.0](https://github.com/containernetworking/plugins/releases/tag/v1.2.0)
- [etcd v3.5.7](https://github.com/etcd-io/etcd/releases/tag/v3.5.7)

### Bug fixes and enhancements

#### For all platforms

- Docker Desktop now automatically pauses the Docker Engine when it is not in use and wakes up again on demand.
- VirtioFS is now the default file sharing implementation for new installations of Docker Desktop on macOS 12.5 and higher.
- Improved product usage reporting using OpenTelemetry (experimental).
- Fixed Docker socket permissions. Fixes [docker/for-win#13447](https://github.com/docker/for-win/issues/13447) and [docker/for-mac#6823](https://github.com/docker/for-mac/issues/6823).
- Fixed an issue which caused Docker Desktop to hang when quitting the application whilst paused.
- Fixed a bug which caused the **Logs** and **Terminal** tab content in the **Container** view to be covered by a fixed toolbar [docker/for-mac#6814](https://github.com/docker/for-mac/issues/6814).
- Fixed a bug which caused input labels to overlap with input values on the container run dialog. Fixes [docker/for-win#13304](https://github.com/docker/for-win/issues/13304).
- Fixed a bug which meant users couldn't select the Docker Extension menu. Fixes [docker/for-mac#6840](https://github.com/docker/for-mac/issues/6840) and [docker/for-mac#6855](https://github.com/docker/for-mac/issues/6855)

#### For Mac

- Added a health check for macOS that notifies users if there has been a change on their system which might cause problems running Docker binaries.

#### For Windows

- Fixed a bug on WSL 2 where if Desktop is paused, killed, and then restarted, the startup hangs unless WSL is shut down first with `wsl --shutdown`.
- Fixed the WSL engine in cases where wsl.exe is not on the PATH [docker/for-win#13547](https://github.com/docker/for-win/issues/13547).
- Fixed the WSL engine's ability to detect cases where one of the Docker Desktop distros' drive is missing [docker/for-win#13554](https://github.com/docker/for-win/issues/13554).
- A slow or unresponsive WSL integration no longer prevents Docker Desktop from starting. Fixes [docker/for-win#13549](https://github.com/docker/for-win/issues/13549).
- Fixed a bug that caused Docker Desktop to crash on startup [docker/for-win#6890](https://github.com/docker/for-mac/issues/6890).
- Added the following installer flags:
  - `--hyper-v-default-data-root` which specifies the default location for Hyper-V VM disk.
  - `--windows-containers-default-data-root` which specifies the default data root for Windows Containers.
  - `--wsl-default-data-root` which specifies the default location for WSL distro disks.

## 4.20.1

{{< release-date date="2023-06-05" >}}

{{< desktop-install all=true version="4.20.1" build_path="/110738/" >}}

### Bug fixes and enhancements

#### For all platforms

- containerd image store: Fixed a bug that caused `docker load` to fail when loading an image that contains attestations.
- containerd image store: Fixed the default image exporter during build.

#### For Windows

- Fixed a bug that made it difficult to parse the WSL version on the host in non-western locales. Fixes [docker/for-win#13518](https://github.com/docker/for-win/issues/13518) and [docker/for-win#13524](https://github.com/docker/for-win/issues/13524).

## 4.20.0

{{< release-date date="2023-05-30" >}}

{{< desktop-install all=true version="4.20.0" build_path="/109717/" >}}

### Upgrades

- [Buildx v0.10.5](https://github.com/docker/buildx/releases/tag/v0.10.5)
- [Compose v2.18.1](https://github.com/docker/compose/releases/tag/v2.18.1)
- [Docker Engine v24.0.2](https://docs.docker.com/engine/release-notes/24.0/#2402)
- [Containerd v1.6.21](https://github.com/containerd/containerd/releases/tag/v1.6.21)
- [runc v1.1.7](https://github.com/opencontainers/runc/releases/tag/v1.1.5)

### Bug fixes and enhancements

#### For all platforms

- [Docker Scout CLI](https://docs.docker.com/scout/#docker-scout-cli) now finds the most recently built image if it is not provided as an argument.
- Improved the [Docker Scout CLI](https://docs.docker.com/scout/#docker-scout-cli) `compare` command.
- Added a warning about the [retirement of Docker Compose ECS/ACS integrations in November 2023](https://docs.docker.com/go/compose-ecs-eol/). Can be suppressed with `COMPOSE_CLOUD_EOL_SILENT=1`.
- Fixed an HTTP proxy bug where an HTTP 1.0 client could receive an HTTP 1.1 response.
- Enabled Docker Desktop's Enhanced Container Isolation (ECI) feature on WSL-2. This is available with a Docker Business subscription.
- Fixed a bug on the **Containers** table where previously hidden columns were displayed again after a fresh installation of Docker Desktop.

#### For Mac

- You can now reclaim disk space more quickly when files are deleted in containers. Related to [docker/for-mac#371](https://github.com/docker/for-mac/issues/371).
- Fixed a bug that prevented containers accessing 169.254.0.0/16 IPs. Fixes [docker/for-mac#6825](https://github.com/docker/for-mac/issues/6825).
- Fixed a bug in `com.docker.diagnose check` where it would complain about a missing vpnkit even when vpnkit is not expected to be running. Related to [docker/for-mac#6825](https://github.com/docker/for-mac/issues/6825).

#### For Windows

- Fixed a bug that meant WSL data could not be moved to a different disk. Fixes [docker/for-win#13269](https://github.com/docker/for-win/issues/13269).
- Fixed a bug where Docker Desktop was not stopping its WSL distros (docker-desktop and docker-desktop-data) when it was shutdown, consuming host memory unnecessarily.
- Added a new setting that allows the Windows Docker daemon to use Docker Desktop's internal proxy when running Windows containers. See [Windows proxy settings](settings/windows.md#proxies).

#### For Linux

- Fixed an issue with the Docker Compose V1/V2 compatibility setting.

## 4.19.0

{{< release-date date="2023-04-27" >}}

{{< desktop-install all=true version="4.19.0" build_path="/106363/" >}}

### New

- Docker Engine and CLI updated to [Moby 23.0](https://github.com/moby/moby/releases/tag/v23.0.0).
- The **Learning Center** now supports in-product walkthroughs.
- Docker init (Beta) now supports Node.js and Python.
- Faster networking between VM and host on macOS.
- You can now inspect and analyze remote images from Docker Desktop without pulling them.
- Usability and performance improvements to the **Artifactory images** view.

### Removed

- Removed `docker scan` command. To continue learning about the vulnerabilities of your images, and many other features, use the new `docker scout` command. Run `docker scout --help`, or [read the docs to learn more](../engine/reference/commandline/scout.md).

### Upgrades

- [Docker Engine v23.0.5](https://docs.docker.com/engine/release-notes/23.0/#2305)
- [Compose 2.17.3](https://github.com/docker/compose/releases/tag/v2.17.3)
- [Containerd v1.6.20](https://github.com/containerd/containerd/releases/tag/v1.6.20)
- [Kubernetes v1.25.9](https://github.com/kubernetes/kubernetes/releases/tag/v1.25.9)
- [runc v1.1.5](https://github.com/opencontainers/runc/releases/tag/v1.1.5)
- [Go v1.20.3](https://github.com/golang/go/releases/tag/go1.20.3)

### Bug fixes and enhancements

#### For all platforms

- Improved `docker scout compare` command to compare two images, now also aliased under `docker scout diff`.
- Added more details to dashboard errors when a `docker-compose` action fails ([docker/for-win#13378](https://github.com/docker/for-win/issues/13378)).
- Added support for setting HTTP proxy configuration during installation. This can be done via the `--proxy-http-mode`, `--overrider-proxy-http`, `--override-proxy-https` and `--override-proxy-exclude` installer flags in the case of installation from the CLI on [Mac](install/mac-install.md#install-from-the-command-line) and [Windows](install/windows-install.md#install-from-the-command-line), or alternatively by setting the values in the `install-settings.json` file.
- Docker Desktop now stops overriding .docker/config.json `credsStore` keys on application start. Note that if you use a custom credential helper then the CLI `docker login` and `docker logout` does not affect whether the UI is signed in to Docker or not. In general, it is better to sign into Docker via the UI since the UI supports multi-factor authentication.
- Added a warning about the [forthcoming removal of Compose V1 from Docker Desktop](../compose/migrate.md). Can be suppressed with `COMPOSE_V1_EOL_SILENT=1`.
- In the Compose config, boolean fields in YAML should be either `true` or `false`. Deprecated YAML 1.1 values such as “on” or “no” now produce a warning.
- Improved UI for image table, allowing rows to use more available space.
- Fixed various bugs in port-forwarding.
- Fixed a HTTP proxy bug where an HTTP request without a Server Name Indication record would be rejected with an error.

#### For Windows

- Reverted to fully patching etc/hosts on Windows (includes `host.docker.internal` and `gateway.docker.internal` again). For WSL, this behavior is controlled by a new setting in the **General** tab. Fixes [docker/for-win#13388](https://github.com/docker/for-win/issues/13388) and [docker/for-win#13398](https://github.com/docker/for-win/issues/13398).
- Fixed a spurious `courgette.log` file appearing on the Desktop when updating Docker Desktop. Fixes [docker/for-win#12468](https://github.com/docker/for-win/issues/12468).
- Fixed the "zoom in" shortcut (ctrl+=). Fixes [docker/for-win#13392](https://github.com/docker/for-win/issues/13392).
- Fixed a bug where the tray menu would not correctly update after second container type switch. Fixes [docker/for-win#13379](https://github.com/docker/for-win/issues/13379).

#### For Mac

- Increased the performance of VM networking when using the Virtualization framework on macOS Ventura and above. Docker Desktop for Mac now uses gVisor instead of VPNKit. To continue using VPNKit, add `"networkType":"vpnkit"` to your `settings.json` file located at `~/Library/Group Containers/group.com.docker/settings.json`.
- Fixed a bug where an error window is displayed on uninstall.
- Fixed a bug where the setting `deprecatedCgroupv1` was ignored. Fixes [docker/for-mac#6801](https://github.com/docker/for-mac/issues/6801).
- Fixed cases where `docker pull` would return `EOF`.

#### For Linux

- Fixed a bug where the VM networking crashes after 24h. Fixes [docker/desktop-linux#131](https://github.com/docker/desktop-linux/issues/131).

### Security

#### For all platforms

- Fixed a security issue allowing users to bypass Image Access Management (IAM) restrictions configured by their organisation by avoiding `registry.json` enforced login via deleting the `credsStore` key from their Docker CLI configuration file. Only affects Docker Business customers.
- Fixed [CVE-2023-24532](https://github.com/advisories/GHSA-x2w5-7wp4-5qff).
- Fixed [CVE-2023-25809](https://github.com/advisories/GHSA-m8cg-xc2p-r3fc).
- Fixed [CVE-2023-27561](https://github.com/advisories/GHSA-vpvm-3wq2-2wvm).
- Fixed [CVE-2023-28642](https://github.com/advisories/GHSA-g2j6-57v7-gm8c).
- Fixed [CVE-2023-28840](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-28840).
- Fixed [CVE-2023-28841](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-28841).
- Fixed [CVE-2023-28842](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-28842).

## 4.18.0

{{< release-date date="2023-04-03" >}}

{{< desktop-install all=true version="4.18.0" build_path="/104112/" >}}

### New

- Initial beta release of `docker init` as per [the roadmap](https://github.com/docker/roadmap/issues/453).
- Added a new **Learning Center** tab to help users get started with Docker.
- Added an experimental file-watch command to Docker Compose that automatically updates your running Compose services as you edit and save your code.

### Upgrades

- [Buildx v0.10.4](https://github.com/docker/buildx/releases/tag/v0.10.4)
- [Compose 2.17.2](https://github.com/docker/compose/releases/tag/v2.17.2)
- [Containerd v1.6.18](https://github.com/containerd/containerd/releases/tag/v1.6.18), which includes fixes for [CVE-2023-25153](https://github.com/advisories/GHSA-259w-8hf6-59c2) and [CVE-2023-25173](https://github.com/advisories/GHSA-hmfx-3pcx-653p).
- [Docker Engine v20.10.24](https://docs.docker.com/engine/release-notes/20.10/#201024), which contains fixes for [CVE-2023-28841](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-28841),
  [CVE-2023-28840](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-28840), and
  [CVE-2023-28842](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-28842).

### Bug fixes and enhancements

#### For all platforms

- [Docker Scout CLI](../scout/index.md#docker-scout-cli) can now compare two images and display packages and vulnerabilities differences. This command is in [Early Access](../release-lifecycle.md) and might change in the future.
- [Docker Scout CLI](../scout/index.md#docker-scout-cli) now displays base image update and remediation recommendations using `docker scout recommendations`. It also displays a short overview of an image using `docker scout quickview` commands.
- You can now search for extensions direct from the Marketplace, as well as using **Global Search**.
- Fixed a bug where `docker buildx` container builders would lose access to the network after 24hrs.
- Reduced how often users are prompted for feedback on Docker Desktop.
- Removed minimum VM swap size.
- Added support for subdomain match, CIDR match, `.` and `_.` in HTTP proxy exclude lists.
- Fixed a bug in the transparent TLS proxy when the Server Name Indication field is not set.
- Fixed a grammatical error in Docker Desktop engine status message.

### For Windows

- Fixed a bug where `docker run --gpus=all` hangs. Fixes [docker/for-win#13324](https://github.com/docker/for-win/issues/13324).
- Fixed a bug where Registry Access Management policy updates were not downloaded.
- Docker Desktop now allows Windows containers to work when BitLocker is enabled on `C:`.
- Docker Desktop with the WSL backend no longer requires the `com.docker.service` privileged service to run permanently. For more information see [Permission requirements for Windows](https://docs.docker.com/desktop/windows/permission-requirements/).

### For Mac

- Fixed a performance issue where attributes stored on the host would not be cached for VirtioFS users.
- The first time Docker Desktop for Mac is launched, the user is presented with an installation window to confirm or adjust the configuration that requires privileged access. For more information see [Permission requirements for Mac](https://docs.docker.com/desktop/mac/permission-requirements/).
- Added the **Advanced** tab in **Settings**, where users can adjust the settings which require privileged access.

### For Linux

- Fixed a bug where the VM networking crashes after 24h. [docker/for-linux#131](https://github.com/docker/desktop-linux/issues/131)

### Security

#### For all platforms

- Fixed [CVE-2023-1802](https://www.cve.org/cverecord?id=CVE-2023-1802) where a security issue with the Artifactory Integration would cause it to fall back to sending registry credentials over plain HTTP if HTTPS check failed. Only users who have `Access experimental features` enabled are affected. Fixes [docker/for-win#13344](https://github.com/docker/for-win/issues/13344).

#### For Mac

- Removed the `com.apple.security.cs.allow-dyld-environment-variables` and `com.apple.security.cs.disable-library-validation` entitlements which allow an arbitrary dynamic library to be loaded with Docker Desktop via the `DYLD_INSERT_LIBRARIES` environment variable.

### Known Issues

- Uninstalling Docker Desktop on Mac from the **Troubleshoot** page might trigger an unexpected fatal error popup.

## 4.17.1

{{< release-date date="2023-03-20" >}}

### Bug fixes and enhancements

#### For Windows

- Docker Desktop now allows Windows containers to work when BitLocker is enabled on C:
- Fixed a bug where `docker buildx` container builders would lose access to the network after 24hrs.
- Fixed a bug where Registry Access Management policy updates were not downloaded.
- Improved debug information to better characterise failures under WSL 2.

### Known Issues

- Running containers with `--gpus` on Windows with the WSL 2 backend does not work. This will be fixed in future releases. See [docker/for-win/13324](https://github.com/docker/for-win/issues/13324).

## 4.17.0

{{< release-date date="2023-02-27" >}}

### New

- Docker Desktop now ships with Docker Scout. Pull and view analysis for images from Docker Hub and Artifactory repositories, get base image updates and recommended tags and digests, and filter your images on vulnerability information. To learn more, see [Docker Scout](../scout/index.md).
- `docker scan` has been replaced by `docker scout`. See [Docker Scout CLI](../scout/index.md#docker-scout-cli), for more information.
- You can now discover extensions that have been autonomously published in the Extensions Marketplace. For more information on self-published extensions, see [Marketplace Extensions](extensions/marketplace.md).
- **Container File Explorer** is available as an experimental feature. Debug the filesystem within your containers straight from the GUI.
- You can now search for volumes in **Global Search**.

### Upgrades

- [Containerd v1.6.18](https://github.com/containerd/containerd/releases/tag/v1.6.18), which includes fixes for [CVE-2023-25153](https://github.com/advisories/GHSA-259w-8hf6-59c2) and [CVE-2023-25173](https://github.com/advisories/GHSA-hmfx-3pcx-653p).
- [Docker Engine v20.10.23](https://docs.docker.com/engine/release-notes/20.10/#201023).
- [Go 1.19.5](https://github.com/golang/go/releases/tag/go1.19.5)

### Bug fixes and enhancements

#### For all platforms

- Fixed a bug where diagnostic gathering could hang waiting for a subprocess to exit.
- Prevented the transparent HTTP proxy from mangling requests too much. Fixes Tailscale extension login, see [tailscale/docker-extension#49](https://github.com/tailscale/docker-extension/issues/49).
- Fixed a bug in the transparent TLS proxy where the Server Name Indication field is not set.
- Added support for subdomain match, CIDR match, `.` and `*.` in HTTP proxy exclude lists.
- Ensured HTTP proxy settings are respected when uploading diagnostics.
- Fixed fatal error when fetching credentials from the credential helper.
- Fixed fatal error related to concurrent logging.
- Improved the UI for Extension actions in the Marketplace.
- Added new filters in the Extensions Marketplace. You can now filter extensions by category and reviewed status.
- Added a way to report a malicious extension to Docker.
- Updated Dev Environments to v0.2.2 with initial set up reliability & security fixes.
- Added a whalecome survey for new users only.
- The confirmation dialogs on the troubleshooting page are now consistent in style with other similar dialogs.
- Fixed fatal error caused by resetting the Kubernetes cluster before it has started.
- Implemented `docker import` for the containerd integration.
- Fixed image tagging with an existing tag with the containerd integration.
- Implemented the dangling filter on images for the containerd integration.
- Fixed `docker ps` failing with containers whose images are no longer present with the containerd integration.

#### For Mac

- Fixed download of Registry Access Management policy on systems where the privileged helper tool `com.docker.vmnetd` is not installed.
- Fixed a bug where `com.docker.vmnetd` could not be installed if `/Library/PrivilegedHelperTools` does not exist.
- Fixed a bug where the "system" proxy would not handle "autoproxy" / "pac file" configurations.
- Fixed a bug where vmnetd installation fails to read `Info.Plist` on case-sensitive file systems. The actual filename is `Info.plist`. Fixes [docker/for-mac#6677](https://github.com/docker/for-mac/issues/6677).
- Fixed a bug where user is prompted to create the docker socket symlink on every startup. Fixes [docker/for-mac#6634](https://github.com/docker/for-mac/issues/6634).
- Fixed a bug that caused the **Start Docker Desktop when you log in** setting not to work. Fixes [docker/for-mac#6723](https://github.com/docker/for-mac/issues/6723).
- Fixed UDP connection tracking and `host.docker.internal`. Fixes [docker/for-mac#6699](https://github.com/docker/for-mac/issues/6699).
- Improved kubectl symlink logic to respect existing binaries in `/usr/local/bin`. Fixes [docker/for-mac#6328](https://github.com/docker/for-mac/issues/6328).
- Docker Desktop now automatically installs Rosetta when you opt-in to use it but have not already installed it.

### For Windows

- Added statical linking of WSL integration tools against `musl` so there is no need to install `alpine-pkg-glibc` in user distros.
- Added support for running under cgroupv2 on WSL 2. This is activated by adding `kernelCommandLine = systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all` to your `%USERPROFILE%\.wslconfig` file in the `[wsl2]` section.
- Fixed an issue that caused Docker Desktop to get stuck in the "starting" phase when in WSL 2 mode (introduced in 4.16).
- Fixed Docker Desktop failing to start the WSL 2 backend when file system compression or encryption is enabled on `%LOCALAPPDATA%`.
- Fixed Docker Desktop failing to report a missing or outdated (incapable of running WSL version 2 distros) WSL installation when starting.
- Fixed a bug where opening in Visual Studio Code fails if the target path has a space.
- Fixed a bug that causes `~/.docker/context` corruption and the error message "unexpected end of JSON input". You can also remove `~/.docker/context` to work around this problem.
- Ensured the credential helper used in WSL 2 is properly signed. Related to [docker/for-win#10247](https://github.com/docker/for-win/issues/10247).
- Fixed an issue that caused WSL integration agents to be terminated erroneously. Related to [docker/for-win#13202](https://github.com/docker/for-win/issues/13202).
- Fixed corrupt contexts on start. Fixes [docker/for-win#13180](https://github.com/docker/for-win/issues/13180) and [docker/for-win#12561](https://github.com/docker/for-win/issues/12561).

### For Linux

- Added Docker Buildx plugin for Docker Desktop for Linux.
- Changed compression algorithm to `xz` for RPM and Arch Linux distribution.
- Fixed a bug that caused leftover files to be left in the root directory of the Debian package. Fixes [docker/for-linux#123](https://github.com/docker/desktop-linux/issues/123).

### Security

#### For all platforms

- Fixed [CVE-2023-0628](https://www.cve.org/cverecord?id=CVE-2023-0628), which allows an attacker to execute an arbitrary command inside a Dev Environments container during initialization by tricking a user to open a crafted malicious `docker-desktop://` URL.
- Fixed [CVE-2023-0629](https://www.cve.org/cverecord?id=CVE-2023-0629), which allows an unprivileged user to bypass Enhanced Container Isolation (ECI) restrictions by setting the Docker host to `docker.raw.sock`, or `npipe:////.pipe/docker_engine_linux` on Windows, via the `-H` (`--host`) CLI flag or the `DOCKER_HOST` environment variable and launch containers without the additional hardening features provided by ECI. This does not affect already running containers, nor containers launched through the usual approach (without Docker's raw socket).

## 4.16.3

{{< release-date date="2023-01-30" >}}

### Bug fixes and enhancements

#### For Windows

- Fixed Docker Desktop failing to start the WSL 2 backend when file system compression or encryption is enabled on `%LOCALAPPDATA%`. Fixes [docker/for-win#13184](https://github.com/docker/for-win/issues/13184).
- Fixed Docker Desktop failing to report a missing or outdated WSL installation when starting. Fixes [docker/for-win#13184](https://github.com/docker/for-win/issues/13184).

## 4.16.2

{{< release-date date="2023-01-19" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed an issue where `docker build` and `docker tag` commands produced an `image already exists` error if the containerd integration feature is enabled.
- Fixed a regression introduced with Docker Desktop 4.16 breaking networking from containers with target platform linux/386 on amd64 systems. Fixes [docker/for-mac/6689](https://github.com/docker/for-mac/issues/6689).

#### For Mac

- Fixed the capitalization of `Info.plist` which caused `vmnetd` to break on case-sensitive file systems. Fixes [docker/for-mac/6677](https://github.com/docker/for-mac/issues/6677).

#### For Windows

- Fixed a regression introduced with Docker Desktop 4.16 causing it to get stuck in the "starting" phase when in WSL2 mode. Fixes [docker/for-win/13165](https://github.com/docker/for-win/issues/13165)

## 4.16.1

{{< release-date date="2023-01-13" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed `sudo` inside a container failing with a security related error for some images. Fixes [docker/for-mac/6675](https://github.com/docker/for-mac/issues/6675) and [docker/for-win/13161](https://github.com/docker/for-win/issues/13161).

## 4.16.0

{{< release-date date="2023-01-12" >}}

### New

- Extensions have moved from Beta to GA.
- Quick Search has moved from experimental to GA.
- Extensions are now included in Quick Search.
- Analyzing large images is now up to 4x faster.
- New local images view has moved from experimental to GA.
- New Beta feature for MacOS 13, Rosetta for Linux, has been added for faster emulation of Intel-based images on Apple Silicon.

### Upgrades

- [Compose v2.15.1](https://github.com/docker/compose/releases/tag/v2.15.1)
- [Containerd v1.6.14](https://github.com/containerd/containerd/releases/tag/v1.6.14)
- [Docker Engine v20.10.22](https://docs.docker.com/engine/release-notes/20.10/#201022)
- [Buildx v0.10.0](https://github.com/docker/buildx/releases/tag/v0.10.0)
- [Docker Scan v0.23.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.23.0)
- [Go 1.19.4](https://github.com/golang/go/releases/tag/go1.19.4)

### Bug fixes and enhancements

#### For all platforms

- Fixed `docker build --quiet` not outputting the image identifier with the `containerd` integration.
- Fixed image inspect not showing image labels with the `containerd` integration.
- Increased the contrast between running and stopped container icons to make it easier for colorblind people to scan the containers list.
- Fixed a bug where the user is prompted for new HTTP proxy credentials repeatedly until Docker Desktop is restarted.
- Added a diagnostics command `com.docker.diagnose login` to check HTTP proxy configuration.
- Fixed actions on compose stack not working properly. Fixes [docker/for-mac#6566](https://github.com/docker/for-mac/issues/6566).
- Fixed the Docker dashboard trying at startup to get disk usage information and display an error banner before the engine was running.
- Added an informational banner with instructions on how to opt-out of experimental feature access next to all experimental features.
- Docker Desktop now supports downloading Kubernetes images via an HTTP proxy.
- Fixed tooltips to not block action buttons. Fixes [docker/for-mac#6516](https://github.com/docker/for-mac/issues/6516).
- Fixed the blank "An error occurred" container list on the **Container** view.

#### For Mac

- Minimum OS version to install or update Docker Desktop on macOS is now macOS Big Sur (version 11) or later.
- Fixed the Docker engine not starting when Enhanced Container Isolation is enabled if the legacy `osxfs` implementation is used for file sharing.
- Fixed files created on VirtioFS having the executable bit set. Fixes [docker/for-mac#6614](https://github.com/docker/for-mac/issues/6614).
- Added back a way to uninstall Docker Desktop from the command line. Fixes [docker/for-mac#6598](https://github.com/docker/for-mac/issues/6598).
- Fixed hardcoded `/usr/bin/kill`. Fixes [docker/for-mac#6589](https://github.com/docker/for-mac/issues/6589).
- Fixed truncation (for example with the `truncate` command) of very large files (> 38GB) shared on VirtioFS with an incorrect size.
- Changed the disk image size in **Settings** to use the decimal system (base 10) to coincide with how Finder displays disk capacity.
- Fixed Docker crash under network load. Fixes [docker/for-mac#6530](https://github.com/docker/for-mac/issues/6530).
- Fixed an issue causing Docker to prompt the user to install the `/var/run/docker.sock` symlink after every reboot.
- Ensured the Login Item which installs the `/var/run/docker.sock` symlink is signed.
- Fixed bug where `$HOME/.docker` was removed on factory reset.

### For Windows

- Fixed `docker build` hanging while printing "load metadata for". Fixes [docker/for-win#10247](https://github.com/docker/for-win/issues/10247).
- Fixed typo in diagnose.exe output Fixes [docker/for-win#13107](https://github.com/docker/for-win/issues/13107).
- Added support for running under cgroupv2 on WSL 2. This is activated by adding `kernelCommandLine = systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all` to your `%USERPROFILE%\.wslconfig` file in the `[wsl2]` section.

### Known Issues

- Calling `sudo` inside a container fails with a security related error for some images. See [docker/for-mac/6675](https://github.com/docker/for-mac/issues/6675) and [docker/for-win/13161](https://github.com/docker/for-win/issues/13161).

## 4.15.0

{{< release-date date="2022-12-01" >}}

### New

- Substantial performance improvements for macOS users with the option of enabling the new VirtioFS file sharing technology. Available for macOS 12.5 and above.
- Docker Desktop for Mac no longer needs to install the privileged helper process `com.docker.vmnetd` on install or on the first run. For more information see [Permission requirements for Mac](https://docs.docker.com/desktop/mac/permission-requirements/).
- Added [WebAssembly capabilities](wasm/index.md). Use with the [containerd integration](containerd/index.md).
- Improved the descriptions for beta and experimental settings to clearly explain the differences and how people can access them.
- Available disk space of VM now displays in the footer of Docker Dashboard for Mac and Linux.
- A disk space warning now displays in the footer if available space is below 3GB.
- Changes to Docker Desktop's interface as we become more ADA accessible and visually unified.
- Added a **Build** tab inside **Extensions** which contains all the necessary resources to build an extension.
- Added the ability to share extensions more easily, either with `docker extension share` CLI or with the share button in the extensions **Manage** tab.
- Extensions in the Marketplace now display the number of installs. You can also sort extensions by the number of installs.
- Dev Environments allow cloning a Git repository to a local bind mount, so you can use any local editor or IDE.
- More Dev Environments improvements: custom names, better private repo support, improved port handling.

### Upgrades

- [Compose v2.13.0](https://github.com/docker/compose/releases/tag/v2.13.0)
- [Containerd v1.6.10](https://github.com/containerd/containerd/releases/tag/v1.6.10)
- [Docker Hub Tool v0.4.5](https://github.com/docker/hub-tool/releases/tag/v0.4.5)
- [Docker Scan v0.22.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.22.0)

### Bug fixes and enhancements

#### For all platforms

- Containers are now restored on restart with the containerd integration.
- Fixed listing multi-platform images with the containerd integration.
- Better handling of dangling images with the containerd integration.
- Implement "reference" filter for images with the containerd integration.
- Added support for selecting upstream HTTP/HTTPS proxies automatically via `proxy.pac` in containers, `docker pull` etc.
- Fixed regressions when parsing image references on pull. Fixes [docker/for-win#13053](https://github.com/docker/for-win/issues/13053), [docker/for-mac#6560](https://github.com/docker/for-mac/issues/6560), and [docker/for-mac#6540](https://github.com/docker/for-mac/issues/6540).

#### For Mac

- Improved the performance of `docker pull`.

#### For Windows

- Fixed an issue where the system HTTP proxies were not used when Docker starts and the developer logs in.
- When Docker Desktop is using "system" proxies and if the Windows settings change, Docker Desktop now uses the new Windows settings without a restart.

#### For Linux

- Fixed hot-reload issue on Linux. Fixes [docker/desktop-linux#30](https://github.com/docker/desktop-linux/issues/30).
- Disabled tray icon animations on Linux which fixes crashes for some users.

## 4.14.1

{{< release-date date="2022-11-17" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed container DNS lookups when using Registry Access Management.

#### For Mac

- Fixed an issue preventing the **Analyze Image** button on the **Images** tab from working.
- Fixed a bug causing symlinks to not be created for the user if `/usr/local/lib` doesn't already exist. Fixes [docker/for-mac#6569](https://github.com/docker/for-mac/issues/6569)

## 4.14.0

{{< release-date date="2022-11-10" >}}

### New

- Set Virtualization framework as the default hypervisor for macOS >= 12.5.
- Migrate previous install to Virtualization framework hypervisor for macOS >= 12.5.
- The Enhanced Container Isolation feature, available to Docker Business users, can now be enabled from the General Settings.

### Updates

- [Docker Engine v20.10.21](../engine/release-notes/20.10.md#201021),
  which contains mitigations against a Git vulnerability, tracked in [CVE-2022-39253](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-39253),
  and updates the handling of `image:tag@digest` image references.
- [Docker Compose v2.12.2](https://github.com/docker/compose/releases/tag/v2.12.2)
- [Containerd v1.6.9](https://github.com/containerd/containerd/releases/tag/v1.6.9)
- [Go 1.19.3](https://github.com/golang/go/releases/tag/go1.19.3)

### Bug fixes and enhancements

#### For all platforms

- Docker Desktop now requires an internal network subnet of size /24. If you were previously using a /28, it is automatically expanded to /24. If you experience networking issues, check to see if you have a clash between the Docker subnet and your infrastructure. Fixes [docker/for-win#13025](https://github.com/docker/for-win/issues/13025).
- Fixed an issue that prevents users from creating Dev Environments when the Git URL has upper-case characters.
- Fix the `vpnkit.exe is not running` error reported in diagnostics.
- Reverted qemu to 6.2.0 to fix errors like `PR_SET_CHILD_SUBREAPER is unavailable` when running emulated amd64 code.
- Enabled [contextIsolation](https://www.electronjs.org/docs/latest/tutorial/context-isolation) and [sandbox](https://www.electronjs.org/docs/latest/tutorial/sandbox) mode inside Extensions. Now Extensions run in a separate context and this limits the harm that malicious code can cause by limiting access to most system resources.
- Included `unpigz` to allow parallel decompression of pulled images.
- Fixed issues related to performing actions on selected containers. [Fixes https://github.com/docker/for-win/issues/13005](https://github.com/docker/for-win/issues/13005)
- Added functionality that allows you to display timestamps for your container or project view.
- Fixed a possible segfault when interrupting `docker pull` with Control+C.
- Increased the default DHCP lease time to avoid the VM's network glitching and dropping connections every two hours.
- Removed the infinite spinner on the containers list. [Fixes https://github.com/docker/for-mac/issues/6486](https://github.com/docker/for-mac/issues/6486)
- Fixed bug which showed incorrect values on used space in **Settings**.
- Fixed a bug that caused Kubernetes not to start with the containerd integration.
- Fixed a bug that caused `kind` not to start with the containerd integration.
- Fixed a bug that caused Dev Environments to not work with the containerd integration.
- Implemented `docker diff` in the containerd integration.
- Implemented `docker run —-platform` in the containerd integration.
- Fixed a bug that caused insecure registries not to work with the containerd integration.

#### For Mac

- Fixed a startup failure for users of Virtualization framework.
- Re-added the `/var/run/docker.sock` on Mac by default, to increase compatibility with tooling like `tilt` and `docker-py.`
- Fixed an issue that prevented the creation of Dev Environments on new Mac installs (error "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?").

#### For Windows

- Re-added `DockerCli.exe -SharedDrives`. Fixes [docker/for-win#5625](https://github.com/docker/for-win#5625).
- Docker Desktop now allows Docker to function on machines where PowerShell is disabled.
- Fixed an issue where Compose v2 was not always enabled by default on Windows.
- Docker Desktop now deletes the `C:\Program Files\Docker` folder at uninstall.

### Known Issues

- For some users on Mac OS there is a known issue with the installer that prevents the installation of a new helper tool needed for the experimental vulnerability and package discovery feature in Docker Desktop. To fix this, a symlink is needed that can be created with the following command: `sudo ln -s /Applications/Docker.app/Contents/Resources/bin/docker-index /usr/local/bin/docker-index`

## 4.13.1

{{< release-date date="2022-10-31" >}}

### Updates

- [Docker Compose v2.12.1](https://github.com/docker/compose/releases/tag/v2.12.1)

### Bug fixes and enhancements

#### For all platforms

- Fixed a possible segfault when interrupting `docker pull` with `Control+C` or `CMD+C`.
- Increased the default DHCP lease time to avoid the VM's network glitching and dropping connections every two hours.
- Reverted `Qemu` to `6.2.0` to fix errors like `PR_SET_CHILD_SUBREAPER is unavailable` when running emulated amd64 code.

#### For Mac

- Added back the `/var/run/docker.sock` symlink on Mac by default, to increase compatibility with tooling like `tilt` and `docker-py`. Fixes [docker/for-mac#6529](https://github.com/docker/for-mac/issues/6529).
- Fixed an issue preventing the creation of Dev Environments on new Mac installs and causing `error "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?")`

#### For Windows

- Docker Desktop now functions on machines where PowerShell is disabled.

## 4.13.0

{{< release-date date="2022-10-19" >}}

### New

- Two new security features have been introduced for Docker Business users, Settings Management and Enhanced Container Isolation. Read more about Docker Desktop’s new [Hardened Docker Desktop security model](hardened-desktop/index.md).
- Added the new Dev Environments CLI `docker dev`, so you can create, list, and run Dev Envs via command line. Now it's easier to integrate Dev Envs into custom scripts.
- Docker Desktop can now be installed to any drive and folder using the `--installation-dir`. Partially addresses [docker/roadmap#94](https://github.com/docker/roadmap/issues/94).

### Updates

- [Docker Scan v0.21.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.21.0)
- [Go 1.19.2](https://github.com/golang/go/releases/tag/go1.19.2) to address [CVE-2022-2879](https://www.cve.org/CVERecord?id=CVE-2022-2879), [CVE-2022-2880](https://www.cve.org/CVERecord?id=CVE-2022-2880) and [CVE-2022-41715](https://www.cve.org/CVERecord?id=CVE-2022-41715)
- Updated Docker Engine and Docker CLI to [v20.10.20](../engine/release-notes/20.10.md#201020),
  which contain mitigations against a Git vulnerability, tracked in [CVE-2022-39253](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-39253),
  and updated handling of `image:tag@digest` image references, as well as a fix for [CVE-2022-36109](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-36109).
- [Docker Credential Helpers v0.7.0](https://github.com/docker/docker-credential-helpers/releases/tag/v0.7.0)
- [Docker Compose v2.12.0](https://github.com/docker/compose/releases/tag/v2.12.0)
- [Kubernetes v1.25.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.25.2)
- [Qemu 7.0.0](https://wiki.qemu.org/ChangeLog/7.0) used for cpu emulation, inside the Docker Desktop VM.
- [Linux kernel 5.15.49](https://hub.docker.com/layers/docker/for-desktop-kernel/5.15.49-13422a825f833d125942948cf8a8688cef721ead/images/sha256-ebf1f6f0cb58c70eaa260e9d55df7c43968874d62daced966ef6a5c5cd96b493?context=explore)

### Bug fixes and enhancements

#### For all platforms

- Docker Desktop now allows the use of TLS when talking to HTTP and HTTPS proxies to encrypt proxy usernames and passwords.
- Docker Desktop now stores HTTP and HTTPS proxy passwords in the OS credential store.
- If Docker Desktop detects that the HTTP or HTTPS proxy password has changed then it will prompt developers for the new password.
- The **Bypass proxy settings for these hosts and domains** setting now handles domain names correctly for HTTPS.
- The **Remote Repositories** view and Tip of the Day now works with HTTP and HTTPS proxies which require authentication
- We’ve introduced dark launch for features that are in early stages of the product development lifecycle. Users that are opted in can opt out at any time in the settings under the “beta features” section.
- Added categories to the Extensions Marketplace.
- Added an indicator in the whale menu and on the **Extension** tab on when extension updates are available.
- Fixed failing uninstalls of extensions with image names that do not have a namespace, as in 'my-extension'.
- Show port mapping explicitly in the **Container** tab.
- Changed the refresh rate for disk usage information for images to happen automatically once a day.
- Made the tab style consistent for the **Container** and **Volume** tabs.
- Fixed Grpcfuse filesharing mode enablement in **Settings**. Fixes [docker/for-mac#6467](https://github.com/docker/for-mac/issues/6467)
- Virtualization Framework and VirtioFS are disabled for users running macOS < 12.5.
- Ports on the **Containers** tab are now clickable.
- The Extensions SDK now allows `ddClient.extension.vm.cli.exec`, `ddClient.extension.host.cli.exec`, `ddClient.docker.cli.exec` to accept a different working directory and pass environment variables through the options parameters.
- Added a small improvement to navigate to the Extensions Marketplace when clicking on **Extensions** in the sidebar.
- Added a badge to identify new extensions in the Marketplace.
- Fixed kubernetes not starting with the containerd integration.
- Fixed `kind` not starting with the containerd integration.
- Fixed dev environments not working with the containerd integration.
- Implemented `docker diff` in the containerd integration.
- Implemented `docker run —-platform` in the containerd integration.
- Fixed insecure registries not working with the containerd integration.
- Fixed a bug that showed incorrect values on used space in **Settings**.
- Docker Desktop now installs credential helpers from Github releases. See [docker/for-win#10247](https://github.com/docker/for-win/issues/10247), [docker/for-win#12995](https://github.com/docker/for-win/issues/12995).
- Fixed an issue where users were logged out of Docker Desktop after 7 days.

#### For Mac

- Added **Hide**, **Hide others**, **Show all** menu items for Docker Desktop. See [docker/for-mac#6446](https://github.com/docker/for-mac/issues/6446).
- Fixed a bug which caused the application to be deleted when running the install utility from the installed application. Fixes [docker/for-mac#6442](https://github.com/docker/for-mac/issues/6442).
- By default Docker will not create the /var/run/docker.sock symlink on the host and use the docker-desktop CLI context instead.

#### For Linux

- Fixed a bug that prevented pushing images from the Dashboard

## 4.12.0

{{< release-date date="2022-09-01" >}}

### New

- Added the ability to use containerd for pulling and storing images. This is an experimental feature.
- Docker Desktop now runs untagged images. Fixes [docker/for-mac#6425](https://github.com/docker/for-mac/issues/6425).
- Added search capabilities to Docker Extension's Marketplace. Fixes [docker/roadmap#346](https://github.com/docker/roadmap/issues/346).
- Added the ability to zoom in, out or set Docker Desktop to Actual Size. This is done by using keyboard shortcuts ⌘ + / CTRL +, ⌘ - / CTRL -, ⌘ 0 / CTRL 0 on Mac and Windows respectively, or through the View menu on Mac.
- Added compose stop button if any related container is stoppable.
- Individual compose containers are now deletable from the **Container** view.
- Removed the workaround for virtiofsd <-> qemu protocol mismatch on Fedora 35, as it is no longer needed. Fedora 35 users should upgrade the qemu package to the most recent version (qemu-6.1.0-15.fc35 as of the time of writing).
- Implemented an integrated terminal for containers.
- Added a tooltip to display the link address for all external links by default.

### Updates

- [Docker Compose v2.10.2](https://github.com/docker/compose/releases/tag/v2.10.2)
- [Docker Scan v0.19.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.19.0)
- [Kubernetes v1.25.0](https://github.com/kubernetes/kubernetes/releases/tag/v1.25.0)
- [Go 1.19](https://github.com/golang/go/releases/tag/go1.19)
- [cri-dockerd v0.2.5](https://github.com/Mirantis/cri-dockerd/releases/tag/v0.2.5)
- [Buildx v0.9.1](https://github.com/docker/buildx/releases/tag/v0.9.1)
- [containerd v1.6.8](https://github.com/containerd/containerd/releases/tag/v1.6.8)
- [containerd v1.6.7](https://github.com/containerd/containerd/releases/tag/v1.6.7)
- [runc v1.1.4](https://github.com/opencontainers/runc/releases/tag/v1.1.4)
- [runc v1.1.3](https://github.com/opencontainers/runc/releases/tag/v1.1.3)

### Security

#### For all platforms

- Fixed [CVE-2023-0626](https://www.cve.org/cverecord?id=CVE-2023-0626) which allows RCE via query parameters in the message-box route in the Electron client.
- Fixed [CVE-2023-0625](https://www.cve.org/cverecord?id=CVE-2023-0625) which allows RCE via extension description/changelog which could be abused by a malicious extension.

#### For Windows

- Fixed [CVE-2023-0627](https://www.cve.org/cverecord?id=CVE-2023-0627) which allows to bypass for the `--no-windows-containers` installation flag which was introduced in version 4.11. This flag allows administrators to disable the use of Windows containers.
- Fixed [CVE-2023-0633](https://www.cve.org/cverecord?id=CVE-2023-0633) in which an argument injection to the Docker Desktop installer which may result in local privilege escalation.

### Bug fixes and minor enhancements

#### For all platforms

- Compose V2 is now enabled after factory reset.
- Compose V2 is now enabled by default on new installations of Docker Desktop.
- Precedence order of environment variables in Compose is more consistent, and clearly [documented](../compose/environment-variables/envvars-precedence.md).
- Upgraded kernel to 5.10.124.
- Improved overall performance issues caused by calculating disk size. Related to [docker/for-win#9401](https://github.com/docker/for-win/issues/9401).
- Docker Desktop now prevents users on ARM macs without Rosetta installed from switching back to Compose V1, which has only intel binaries.
- Changed the default sort order to descending for volume size and the **Created** column, along with the container's **Started** column.
- Re-organized container row actions by keeping only the start/stop and delete actions visible at all times, while allowing access to the rest via the row menu item.
- The Quickstart guide now runs every command immediately.
- Defined the sort order for container/compose **Status** column to running > some running > paused > some paused > exited > some exited > created.
- Fixed issues with the image list appearing empty in Docker Desktop even though there are images. Related to [docker/for-win#12693](https://github.com/docker/for-win/issues/12693) and [docker/for-mac#6347](https://github.com/docker/for-mac/issues/6347).
- Defined what images are "in use" based on whether or not system containers are displayed. If system containers related to Kubernetes and Extensions are not displayed, the related images are not defined as "in use."
- Fixed a bug that made Docker clients in some languages hang on `docker exec`. Fixes [https://github.com/apocas/dockerode/issues/534](https://github.com/apocas/dockerode/issues/534).
- A failed spawned command when building an extension no longer causes Docker Desktop to unexpectedly quit.
- Fixed a bug that caused extensions to be displayed as disabled in the left menu when they are not.
- Fixed `docker login` to private registries when Registry Access Management is enabled and access to Docker Hub is blocked.
- Fixed a bug where Docker Desktop fails to start the Kubernetes cluster if the current cluster metadata is not stored in the `.kube/config` file.
- Updated the tooltips in Docker Desktop and MUI theme package to align with the overall system design.
- Copied terminal contents do not contain non-breaking spaces anymore.

#### For Mac

- Minimum version to install or update Docker Desktop on macOS is now 10.15. Fixes [docker/for-mac#6007](https://github.com/docker/for-mac/issues/6007).
- Fixed a bug where the Tray menu incorrectly displays "Download will start soon..." after downloading the update. Fixes some issue reported in [for-mac/issues#5677](https://github.com/docker/for-mac/issues/5677)
- Fixed a bug that didn't restart Docker Desktop after applying an update.
- Fixed a bug that caused the connection to Docker to be lost when the computer sleeps if a user is using virtualization.framework and restrictive firewall software.
- Fixed a bug that caused Docker Desktop to run in the background even after a user had quit the application. Fixes [docker/for-mac##6440](https://github.com/docker/for-mac/issues/6440)
- Disabled both Virtualization Framework and VirtioFS for users running macOS < 12.5

#### For Windows

- Fixed a bug where versions displayed during an update could be incorrect. Fixes [for-win/issues#12822](https://github.com/docker/for-win/issues/12822).

## 4.11.1

{{< release-date date="2022-08-05" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed regression preventing VM system locations (e.g. /var/lib/docker) from being bind mounted [for-mac/issues#6433](https://github.com/docker/for-mac/issues/6433)

#### For Windows

- Fixed `docker login` to private registries from WSL2 distro [docker/for-win#12871](https://github.com/docker/for-win/issues/12871)

## 4.11.0

{{< release-date date="2022-07-28" >}}

### New

- Docker Desktop is now fully supported for Docker Business customers inside VMware ESXi and Azure VMs. For more information, see [Run Docker Desktop inside a VM or VDI environment](../desktop/vm-vdi.md)
- Added two new extensions ([vcluster](https://hub.docker.com/extensions/loftsh/vcluster-dd-extension) and [PGAdmin4](https://hub.docker.com/extensions/mochoa/pgadmin4-docker-extension)) to the Extensions Marketplace.
- The ability to sort extensions has been added to the Extensions Marketplace.
- Fixed a bug that caused some users to be asked for feedback too frequently. You'll now only be asked for feedback twice a year.
- Added custom theme settings for Docker Desktop. This allows you to specify dark or light mode for Docker Desktop independent of your device settings. Fixes [docker/for-win#12747](https://github.com/docker/for-win/issues/12747)
- Added a new flag for Windows installer. `--no-windows-containers` disables the Windows containers integration.
- Added a new flag for Mac install command. `--user <username>` sets up Docker Desktop for a specific user, preventing them from needing an admin password on first run.

### Updates

- [Docker Compose v2.7.0](https://github.com/docker/compose/releases/tag/v2.7.0)
- [Docker Compose "Cloud Integrations" v1.0.28](https://github.com/docker/compose-cli/releases/tag/v1.0.28)
- [Kubernetes v1.24.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.2)
- [Go 1.18.4](https://github.com/golang/go/releases/tag/go1.18.4)

### Bug fixes and enhancements

#### For all platforms

- Added the Container / Compose icon as well as the exposed port(s) / exit code to the Containers screen.
- Updated the Docker theme palette colour values to match our design system.
- Improved an error message from `docker login` if Registry Access Management is blocking the Docker engine's access to Docker Hub.
- Increased throughput between the Host and Docker. For example increasing performance of `docker cp`.
- Collecting diagnostics takes less time to complete.
- Selecting or deselecting a compose app on the containers overview now selects/deselects all its containers.
- Tag names on the container overview image column are visible.
- Added search decorations to the terminal's scrollbar so that matches outside the viewport are visible.
- Fixed an issue with search which doesn't work well on containers page [docker/for-win#12828](https://github.com/docker/for-win/issues/12828).
- Fixed an issue which caused infinite loading on the **Volume** screen [docker/for-win#12789](https://github.com/docker/for-win/issues/12789).
- Fixed a problem in the Container UI where resizing or hiding columns didn't work. Fixes [docker/for-mac#6391](https://github.com/docker/for-mac/issues/6391).
- Fixed a bug where the state of installing, updating, or uninstalling multiple extensions at once was lost when leaving the Marketplace screen.
- Fixed an issue where the compose version in the about page would only get updated from v2 to v1 after restarting Docker Desktop.
- Fixed an issue where users cannot see the log view because their underlying hardware didn't support WebGL2 rendering. Fixes [docker/for-win#12825](https://github.com/docker/for-win/issues/12825).
- Fixed a bug where the UI for Containers and Images got out of sync.
- Fixed a startup race when the experimental virtualization framework is enabled.

#### For Mac

- Fixed an issue executing Compose commands from the UI. Fixes [docker/for-mac#6400](https://github.com/docker/for-mac/issues/6400).

#### For Windows

- Fixed horizontal resizing issue. Fixes [docker/for-win#12816](https://github.com/docker/for-win/issues/12816).
- If an HTTP/HTTPS proxy is configured in the UI, then it automatically sends traffic from image builds and running containers to the proxy. This avoids the need to separately configure environment variables in each container or build.
- Added the `--backend=windows` installer option to set Windows containers as the default backend.

#### For Linux

- Fixed bug related to setting up file shares with spaces in their path.

## 4.10.1

{{< release-date date="2022-07-05" >}}

### Bug fixes and enhancements

#### For Windows

- Fixed a bug where actions in the UI failed with Compose apps that were created from WSL. Fixes [docker/for-win#12806](https://github.com/docker/for-win/issues/12806).

#### For Mac

- Fixed a bug where the install command failed because paths were not initialized. Fixes [docker/for-mac#6384](https://github.com/docker/for-mac/issues/6384).

## 4.10.0

{{< release-date date="2022-06-30" >}}

### New

- You can now add environment variables before running an image in Docker Desktop.
- Added features to make it easier to work with a container's logs, such as regular expression search and the ability to clear container logs while the container is still running.
- Implemented feedback on the containers table. Added ports and separated container and image names.
- Added two new extensions, Ddosify and Lacework, to the Extensions Marketplace.

### Removed

- Removed Homepage while working on a new design. You can provide [feedback here](https://docs.google.com/forms/d/e/1FAIpQLSfYueBkJHdgxqsWcQn4VzBn2swu4u_rMQRIMa8LExYb_72mmQ/viewform?entry.1237514594=4.10).

### Updates

- [Docker Engine v20.10.17](../engine/release-notes/20.10.md#201017)
- [Docker Compose v2.6.1](https://github.com/docker/compose/releases/tag/v2.6.1)
- [Kubernetes v1.24.1](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.1)
- [cri-dockerd to v0.2.1](https://github.com/Mirantis/cri-dockerd/releases/tag/v0.2.1)
- [CNI plugins to v1.1.1](https://github.com/containernetworking/plugins/releases/tag/v1.1.1)
- [containerd to v1.6.6](https://github.com/containerd/containerd/releases/tag/v1.6.6)
- [runc to v1.1.2](https://github.com/opencontainers/runc/releases/tag/v1.1.2)
- [Go 1.18.3](https://github.com/golang/go/releases/tag/go1.18.3)

### Bug fixes and enhancements

#### For all platforms

- Added additional bulk actions for starting/pausing/stopping selected containers in the **Containers** tab.
- Added pause and restart actions for compose projects in the **Containers** tab.
- Added icons and exposed ports or exit code information in the **Containers** tab.
- External URLs can now refer to extension details in the Extension Marketplace using links such as `docker-desktop://extensions/marketplace?extensionId=docker/logs-explorer-extension`.
- The expanded or collapsed state of the Compose apps is now persisted.
- `docker extension` CLI commands are available with Docker Desktop by default.
- Increased the size of the screenshots displayed in the Extension marketplace.
- Fixed a bug where a Docker extension fails to load if its backend container(s) are stopped. Fixes [docker/extensions-sdk#16](https://github.com/docker/extensions-sdk/issues/162).
- Fixed a bug where the image search field is cleared without a reason. Fixes [docker/for-win#12738](https://github.com/docker/for-win/issues/12738).
- Fixed a bug where the license agreement does not display and silently blocks Docker Desktop startup.
- Fixed the displayed image and tag for unpublished extensions to actually display the ones from the installed unpublished extension.
- Fixed the duplicate footer on the Support screen.
- Dev Environments can be created from a subdirectory in a GitHub repository.
- Removed the error message if the tips of the day cannot be loaded when using Docker Desktop offline. Fixes [docker/for-mac#6366](https://github.com/docker/for-mac/issues/6366).

#### For Mac

- Fixed a bug with location of bash completion files on macOS. Fixes [docker/for-mac#6343](https://github.com/docker/for-mac/issues/6343).
- Fixed a bug where Docker Desktop does not start if the username is longer than 25 characters. Fixes [docker/for-mac#6122](https://github.com/docker/for-mac/issues/6122).
- Fixed a bug where Docker Desktop was not starting due to invalid system proxy configuration. Fixes some issues reported in [docker/for-mac#6289](https://github.com/docker/for-mac/issues/6289).
- Fixed a bug where Docker Desktop failed to start when the experimental virtualization framework is enabled.
- Fixed a bug where the tray icon still displayed after uninstalling Docker Desktop.

#### For Windows

- Fixed a bug which caused high CPU usage on Hyper-V. Fixes [docker/for-win#12780](https://github.com/docker/for-win/issues/12780).
- Fixed a bug where Docker Desktop for Windows would fail to start. Fixes [docker/for-win#12784](https://github.com/docker/for-win/issues/12784).
- Fixed the `--backend=wsl-2` installer flag which did not set the backend to WSL 2. Fixes [docker/for-win#12746](https://github.com/docker/for-win/issues/12746).

#### For Linux

- Fixed a bug when settings cannot be applied more than once.
- Fixed Compose version displayed in the `About` screen.

### Known Issues

- Occasionally the Docker engine will restart during a `docker system prune`. This is a [known issue](https://github.com/moby/buildkit/pull/2177) in the version of buildkit used in the current engine and will be fixed in future releases.

## 4.9.1

{{< release-date date="2022-06-16" >}}

{{< desktop-install all=true version="4.9.1" build_path="/81317/" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed blank dashboard screen. Fixes [docker/for-win#12759](https://github.com/docker/for-win/issues/12759).

## 4.9.0

{{< release-date date="2022-06-02" >}}

### New

- Added additional guides on the homepage for: Elasticsearch, MariaDB, Memcached, MySQL, RabbitMQ and Ubuntu.
- Added a footer to the Docker Dashboard with general information about the Docker Desktop update status and Docker Engine statistics
- Re-designed the containers table, adding:
  - A button to copy a container ID to the clipboard
  - A pause button for each container
  - Column resizing for the containers table
  - Persistence of sorting and resizing for the containers table
  - Bulk deletion for the containers table

### Updates

- [Compose v2.6.0](https://github.com/docker/compose/releases/tag/v2.6.0)
- [Docker Engine v20.10.16](../engine/release-notes/20.10.md#201016)
- [containerd v1.6.4](https://github.com/containerd/containerd/releases/tag/v1.6.4)
- [runc v1.1.1](https://github.com/opencontainers/runc/releases/tag/v1.1.1)
- [Go 1.18.2](https://github.com/golang/go/releases/tag/go1.18.2)

### Bug fixes and enhancements

#### For all platforms

- Fixed an issue which caused Docker Desktop to hang if you quit the app whilst Docker Desktop was paused.
- Fixed the Kubernetes cluster not resetting properly after the PKI expires.
- Fixed an issue where the Extensions Marketplace was not using the defined http proxies.
- Improved the logs search functionality in Docker Dashboard to allow spaces.
- Middle-button mouse clicks on buttons in the Dashboard now behave as a left-button click instead of opening a blank window.

#### For Mac

- Fixed an issue to avoid creating `/opt/containerd/bin` and `/opt/containerd/lib` on the host if `/opt` has been added to the file sharing directories list.

#### For Windows

- Fixed a bug in the WSL 2 integration where if a file or directory is bind-mounted to a container, and the container exits, then the file or directory is replaced with the other type of object with the same name. For example, if a file is replaced with a directory or a directory with a file, any attempts to bind-mount the new object fails.
- Fixed a bug where the Tray icon and Dashboard UI didn't show up and Docker Desktop didn't fully start. Fixes [docker/for-win#12622](https://github.com/docker/for-win/issues/12622).

### Known issues

#### For Linux

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## 4.8.2

{{< release-date date="2022-05-18" >}}

### Updates

- [Compose v2.5.1](https://github.com/docker/compose/releases/tag/v2.5.1)

### Bug fixes and minor enahancements

- Fixed an issue with manual proxy settings which caused problems when pulling images. Fixes [docker/for-win#12714](https://github.com/docker/for-win/issues/12714) and [docker/for-mac#6315](https://github.com/docker/for-mac/issues/6315).
- Fixed high CPU usage when extensions are disabled. Fixes [docker/for-mac#6310](https://github.com/docker/for-mac/issues/6310).
- Docker Desktop now redacts HTTP proxy passwords in log files and diagnostics.

### Known issues

#### For Linux

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## 4.8.1

{{< release-date date="2022-05-09" >}}

### New

- Released [Docker Desktop for Linux](install/linux-install.md).
- Beta release of [Docker Extensions](extensions/index.md) and Extensions SDK.
- Created a Docker Homepage where you can run popular images and discover how to use them.
- [Compose V2 is now GA](https://www.docker.com/blog/announcing-compose-v2-general-availability/)

### Bug fixes and enhancements

- Fixed a bug that caused the Kubernetes cluster to be deleted when updating Docker Desktop.

### Known issues

#### For Linux

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## 4.8.0

{{< release-date date="2022-05-06" >}}

### New

- Released [Docker Desktop for Linux](install/linux-install.md).
- Beta release of [Docker Extensions](extensions/index.md) and Extensions SDK.
- Created a Docker Homepage where you can run popular images and discover how to use them.
- [Compose V2 is now GA](https://www.docker.com/blog/announcing-compose-v2-general-availability/)

### Updates

- [Compose v2.5.0](https://github.com/docker/compose/releases/tag/v2.5.0)
- [Go 1.18.1](https://github.com/golang/go/releases/tag/go1.18.1)
- [Kubernetes 1.24](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.0)

### Bug fixes and minor enhancements

#### For all platforms

- Introduced reading system proxy. You no longer need to manually configure proxies unless it differs from your OS level proxy.
- Fixed a bug that showed Remote Repositories in the Dashboard when running behind a proxy.
- Fixed vpnkit establishing and blocking the client connection even if the server is gone. See [docker/for-mac#6235](https://github.com/docker/for-mac/issues/6235)
- Made improvements on the Volume tab in Docker Desktop:
  - Volume size is displayed.
  - Columns can be resized, hidden and reordered.
  - A columns sort order and hidden state is persisted, even after Docker Desktop restarts.
  - Row selection is persisted when switching between tabs, even after Docker Desktop restarts.
- Fixed a bug in the Dev Environments tab that did not add a scroll when more items were added to the screen.
- Standardised the header title and action in the Dashboard.
- Added support for downloading Registry Access Management policies through HTTP proxies.
- Fixed an issue related to empty remote repositories when the machine is in sleep mode for an extended period of time.
- Fixed a bug where dangling images were not selected in the cleanup process if their name was not marked as "&lt;none>" but their tag is.
- Improved the error message when `docker pull` fails because an HTTP proxy is required.
- Added the ability to clear the search bar easily in Docker Desktop.
- Renamed the "Containers / Apps" tab to "Containers".
- Fixed a silent crash in the Docker Desktop installer when `C:\ProgramData\DockerDesktop` is a file or a symlink.
- Fixed a bug where an image with no namespace, for example `docker pull <private registry>/image`, would be erroneously blocked by Registry Access Management unless access to Docker Hub was enabled in settings.

#### For Mac

- Docker Desktop's icon now matches Big Sur Style guide. See [docker/for-mac#5536](https://github.com/docker/for-mac/issues/5536)
- Fixed a problem with duplicate Dock icons and Dock icon not working as expected. Fixes [docker/for-mac#6189](https://github.com/docker/for-mac/issues/6189).
- Improved support for the `Cmd+Q` shortcut.

#### For Windows

- Improved support for the `Ctrl+W` shortcut.

### Known issues

#### For all platforms

- Currently, if you are running a Kubernetes cluster, it will be deleted when you upgrade to Docker Desktop 4.8.0. We aim to fix this in the next release.

#### For Linux

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## 4.7.1

{{< release-date date="2022-04-19" >}}

### Bug fixes and enhancements

#### For all platforms

- Fixed a crash on the Quick Start Guide final screen.

#### For Windows

- Fixed a bug where update was failing with a symlink error. Fixes [docker/for-win#12650](https://github.com/docker/for-win/issues/12650).
- Fixed a bug that prevented using Windows container mode. Fixes [docker/for-win#12652](https://github.com/docker/for-win/issues/12652).

## 4.7.0

{{< release-date date="2022-04-07" >}}

### New

- IT Administrators can now install Docker Desktop remotely using the command line.
- Add the Docker Software Bill of Materials (SBOM) CLI plugin. The new CLI plugin enables users to generate SBOMs for Docker images. For more information, see [Docker SBOM](../engine/sbom/index.md).
- Use [cri-dockerd](https://github.com/Mirantis/cri-dockerd) for new Kubernetes clusters instead of `dockershim`. The change is transparent from the user's point of view and Kubernetes containers run on the Docker Engine as before. `cri-dockerd` allows Kubernetes to manage Docker containers using the standard [Container Runtime Interface](https://github.com/kubernetes/cri-api#readme), the same interface used to control other container runtimes. For more information, see [The Future of Dockershim is cri-dockerd](https://www.mirantis.com/blog/the-future-of-dockershim-is-cri-dockerd/).

### Updates

- [Docker Engine v20.10.14](../engine/release-notes/20.10.md#201014)
- [Compose v2.4.1](https://github.com/docker/compose/releases/tag/v2.4.1)
- [Buildx 0.8.2](https://github.com/docker/buildx/releases/tag/v0.8.2)
- [containerd v1.5.11](https://github.com/containerd/containerd/releases/tag/v1.5.11)
- [Go 1.18](https://golang.org/doc/go1.18)

### Security

- Update Docker Engine to v20.10.14 to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769)
- Update containerd to v1.5.11 to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769)

### Bug fixes and enahncements

#### For all platforms

- Fixed a bug where the Registry Access Management policy was never refreshed after a failure.
- Logs and terminals in the UI now respect your OS theme in light and dark mode.
- Easily clean up many volumes at once via multi-select checkboxes.
- Improved login feedback.

#### For Mac

- Fixed an issue that sometimes caused Docker Desktop to display a blank white screen. Fixes [docker/for-mac#6134](https://github.com/docker/for-mac/issues/6134).
- Fixed a problem where gettimeofday() performance drops after waking from sleep when using Hyperkit. Fixes [docker/for-mac#3455](https://github.com/docker/for-mac/issues/3455).
- Fixed an issue that caused Docker Desktop to become unresponsive during startup when `osxfs` is used for file sharing.

#### For Windows

- Fixed volume title. Fixes [docker/for-win#12616](https://github.com/docker/for-win/issues/12616).
- Fixed a bug in the WSL 2 integration that caused Docker commands to stop working after restarting Docker Desktop or after switching to Windows containers.

## 4.6.1

{{< release-date date="2022-03-22" >}}

### Updates

- [Buildx 0.8.1](https://github.com/docker/buildx/releases/tag/v0.8.1)

### Bug fixes and enahncements

- Prevented spinning in vpnkit-forwarder filling the logs with error messages.
- Fixed diagnostics upload when there is no HTTP proxy set. Fixes [docker/for-mac#6234](https://github.com/docker/for-mac/issues/6234).
- Removed a false positive "vm is not running" error from self-diagnose. Fixes [docker/for-mac#6233](https://github.com/docker/for-mac/issues/6233).

## 4.6.0

{{< release-date date="2022-03-14" >}}

### New

#### For all platforms

- The Docker Dashboard Volume Management feature now offers the ability to efficiently clean up volumes using multi-select checkboxes.

#### For Mac

- Docker Desktop 4.6.0 gives macOS users the option of enabling a new experimental file sharing technology called VirtioFS. During testing VirtioFS has been shown to drastically reduce the time taken to sync changes between the host and VM, leading to substantial performance improvements. For more information, see [VirtioFS](settings/mac.md#beta-features).

### Updates

#### For all platforms

- [Docker Engine v20.10.13](../engine/release-notes/20.10.md#201013)
- [Compose v2.3.3](https://github.com/docker/compose/releases/tag/v2.3.3)
- [Buildx 0.8.0](https://github.com/docker/buildx/releases/tag/v0.8.0)
- [containerd v1.4.13](https://github.com/containerd/containerd/releases/tag/v1.4.13)
- [runc v1.0.3](https://github.com/opencontainers/runc/releases/tag/v1.0.3)
- [Go 1.17.8](https://golang.org/doc/go1.17)
- [Linux kernel 5.10.104](https://hub.docker.com/layers/docker/for-desktop-kernel/5.10.104-379cadd2e08e8b25f932380e9fdaab97755357b3/images/sha256-7753b60f4544e5c5eed629d12151a49c8a4b48d98b4fb30e4e65cecc20da484d?context=explore)

### Security

#### For all platforms

- Fixed [CVE-2022-0847](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-0847), aka “Dirty Pipe”, an issue that could enable attackers to modify files in container images on the host, from inside a container.
  If using the WSL 2 backend, you must update WSL 2 by running `wsl --update`.

#### For Windows

- Fixed [CVE-2022-26659](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-26659), which could allow an attacker to overwrite any administrator writable file on the system during the installation or the update of Docker Desktop.

#### For Mac

- [Qemu 6.2.0](https://wiki.qemu.org/ChangeLog/6.2)

### Bug fixes and enhancements

#### For all platforms

- Fixed uploading diagnostics when an HTTPS proxy is set.
- Made checking for updates from the systray menu open the Software updates settings section.

#### For Mac

- Fixed the systray menu not displaying all menu items after starting Docker Desktop. Fixes [docker/for-mac#6192](https://github.com/docker/for-mac/issues/6192).
- Fixed a regression about Docker Desktop not starting in background anymore. Fixes [docker/for-mac#6167](https://github.com/docker/for-mac/issues/6167).
- Fixed missing Docker Desktop Dock icon. Fixes [docker/for-mac#6173](https://github.com/docker/for-mac/issues/6173).
- Used speed up block device access when using the experimental `virtualization.framework`. See [benchmarks](https://github.com/docker/roadmap/issues/7#issuecomment-1050626886).
- Increased default VM memory allocation to half of physical memory (min 2 GB, max 8 GB) for better out-of-the-box performances.

#### For Windows

- Fixed the UI stuck in `starting` state forever although Docker Desktop is working fine from the command line.
- Fixed missing Docker Desktop systray icon [docker/for-win#12573](https://github.com/docker/for-win/issues/12573)
- Fixed Registry Access Management under WSL 2 with latest 5.10.60.1 kernel.
- Fixed a UI crash when selecting the containers of a Compose application started from a WSL 2 environment. Fixes [docker/for-win#12567](https://github.com/docker/for-win/issues/12567).
- Fixed copying text from terminal in Quick Start Guide. Fixes [docker/for-win#12444](https://github.com/docker/for-win/issues/12444).

### Known issues

#### For Mac

- After enabling VirtioFS, containers with processes running with different Unix user IDs may experience caching issues. For example if a process running as `root` queries a file and another process running as user `nginx` tries to access the same file immediately, the `nginx` process will get a "Permission Denied" error.

## 4.5.1

{{< release-date date="2022-02-15" >}}

### Bug fixes and enhancements

#### For Windows

- Fixed an issue that caused new installations to default to the Hyper-V backend instead of WSL 2.
- Fixed a crash in the Docker Dashboard which would make the systray menu disappear.

If you are running Docker Desktop on Windows Home, installing 4.5.1 will switch it back to WSL 2 automatically. If you are running another version of Windows, and you want Docker Desktop to use the WSL 2 backend, you must manually switch by enabling the **Use the WSL 2 based engine** option in the **Settings > General** section.
Alternatively, you can edit the Docker Desktop settings file located at `%APPDATA%\Docker\settings.json` and manually switch the value of the `wslEngineEnabled` field to `true`.

## 4.5.0

{{< release-date date="2022-02-10" >}}

### New

- Docker Desktop 4.5.0 introduces a new version of the Docker menu which creates a consistent user experience across all operating systems. For more information, see the blog post [New Docker Menu & Improved Release Highlights with Docker Desktop 4.5](https://www.docker.com/blog/new-docker-menu-improved-release-highlights-with-docker-desktop-4-5/)
- The 'docker version' output now displays the version of Docker Desktop installed on the machine.

### Updates

- [Amazon ECR Credential Helper v0.6.0](https://github.com/awslabs/amazon-ecr-credential-helper/releases/tag/v0.6.0)

### Security

#### For Mac

- Fixed [CVE-2021-44719](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44719) where Docker Desktop could be used to access any user file on the host from a container, bypassing the allowed list of shared folders.

#### For Windows

- Fixed [CVE-2022-23774](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-23774) where Docker Desktop allows attackers to move arbitrary files.

### Bug fixes and enhancements

#### For all platforms

- Fixed an issue where Docker Desktop incorrectly prompted users to sign in after they quit Docker Desktop and start the application.
- Increased the filesystem watch (inotify) limits by setting `fs.inotify.max_user_watches=1048576` and `fs.inotify.max_user_instances=8192` in Linux. Fixes [docker/for-mac#6071](https://github.com/docker/for-mac/issues/6071).

#### For Mac

- Fixed an issue that caused the VM to become unresponsive during startup when using `osxfs` and when no host directories are shared with the VM.
- Fixed an issue that didn't allow users to stop a Docker Compose application using Docker Dashboard if the application was started in a different version of Docker Compose. For example, if the user started a Docker Compose application in V1 and then switched to Docker Compose V2, attempts to stop the Docker Compose application would fail.
- Fixed an issue where Docker Desktop incorrectly prompted users to sign in after they quit Docker Desktop and start the application.
- Fixed an issue where the **About Docker Desktop** window wasn't working anymore.
- Limit the number of CPUs to 8 on Mac M1 to fix the startup problem. Fixes [docker/for-mac#6063](https://github.com/docker/for-mac/issues/6063).

#### For Windows

- Fixed an issue related to compose app started with version 2, but the dashboard only deals with version 1

### Known issues

#### For Windows

Installing Docker Desktop 4.5.0 from scratch has a bug which defaults Docker Desktop to use the Hyper-V backend instead of WSL 2. This means, Windows Home users will not be able to start Docker Desktop as WSL 2 is the only supported backend. To work around this issue, you must uninstall 4.5.0 from your machine and then download and install Docker Desktop 4.5.1 or a higher version. Alternatively, you can edit the Docker Desktop settings.json file located at `%APPDATA%\Docker\settings.json` and manually switch the value of the `wslEngineEnabled` field to `true`.

## 4.4.4

{{< release-date date="2022-01-24" >}}

### Bug fixes and enhancements

#### For Windows

- Fixed logging in from WSL 2. Fixes [docker/for-win#12500](https://github.com/docker/for-win/issues/12500).

### Known issues

#### For Windows

- Clicking **Proceed to Desktop** after signing in through the browser, sometimes does not bring the Dashboard to the front.
- After logging in, when the Dashboard receives focus, it sometimes stays in the foreground even when clicking a background window. As a workaround you need to click the Dashboard before clicking another application window.
- The tips of the week show on top of the mandatory login dialog when an organization restriction is enabled via a `registry.json` file.

## 4.4.3

{{< release-date date="2022-01-14" >}}

### Bug fixes and enhancements

#### For Windows

- Disabled Dashboard shortcuts to prevent capturing them even when minimized or un-focussed. Fixes [docker/for-win#12495](https://github.com/docker/for-win/issues/12495).

### Known issues

#### For Windows

- Clicking **Proceed to Desktop** after signing in through the browser, sometimes does not bring the Dashboard to the front.
- After logging in, when the Dashboard receives focus, it sometimes stays in the foreground even when clicking a background window. As a workaround you need to click the Dashboard before clicking another application window.
- The tips of the week show on top of the mandatory login dialog when an organization restriction is enabled via a `registry.json` file.

## 4.4.2

{{< release-date date="22-01-13" >}}

### New

- Easy, Secure sign in with Auth0 and Single Sign-on
  - Single Sign-on: Users with a Docker Business subscription can now configure SSO to authenticate using their identity providers (IdPs) to access Docker. For more information, see [Single Sign-on](../security/for-admins/single-sign-on/index.md).
  - Signing in to Docker Desktop now takes you through the browser so that you get all the benefits of auto-filling from password managers.

### Upgrades

- [Docker Engine v20.10.12](../engine/release-notes/20.10.md#201012)
- [Compose v2.2.3](https://github.com/docker/compose/releases/tag/v2.2.3)
- [Kubernetes 1.22.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.22.5)
- [docker scan v0.16.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.16.0)

### Security

- Fixed [CVE-2021-45449](../security/index.md#cve-2021-45449) that affects users currently on Docker Desktop version 4.3.0 or 4.3.1.

Docker Desktop version 4.3.0 and 4.3.1 has a bug that may log sensitive information (access token or password) on the user's machine during login.
This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files.

### Bug fixes and enhancements

#### For all platforms

- Docker Desktop displays an error if `registry.json` contains more than one organization in the `allowedOrgs` field. If you are using multiple organizations for different groups of developers, you must provision a separate `registry.json` file for each group.
- Fixed a regression in Compose that reverted the container name separator from `-` to `_`. Fixes [docker/compose-switch](https://github.com/docker/compose-switch/issues/24).

#### For Mac

- Fixed the memory statistics for containers in the Dashboard. Fixes [docker/for-mac/#4774](https://github.com/docker/for-mac/issues/6076).
- Added a deprecated option to `settings.json`: `"deprecatedCgroupv1": true`, which switches the Linux environment back to cgroups v1. If your software requires cgroups v1, you should update it to be compatible with cgroups v2. Although cgroups v1 should continue to work, it is likely that some future features will depend on cgroups v2. It is also possible that some Linux kernel bugs will only be fixed with cgroups v2.
- Fixed an issue where putting the machine to Sleep mode after pausing Docker Desktop results in Docker Desktop not being able to resume from pause after the machine comes out of Sleep mode. Fixes [for-mac#6058](https://github.com/docker/for-mac/issues/6058).

#### For Windows

- Doing a `Reset to factory defaults` no longer shuts down Docker Desktop.

### Known issues

#### For all platforms

- The tips of the week show on top of the mandatory login dialog when an organization restriction is enabled via a `registry.json` file.

#### For Windows

- Clicking **Proceed to Desktop** after logging in in the browser, sometimes does not bring the Dashboard to the front.
- After logging in, when the Dashboard receives focus, it sometimes stays in the foreground even when clicking a background window. As a workaround you need to click the Dashboard before clicking another application window.
- When the Dashboard is open, even if it does not have focus or is minimized, it will still catch keyboard shortcuts (e.g. ctrl-r for Restart)

## 4.3.2

{{< release-date date="2021-12-21" >}}

### Security

- Fixed [CVE-2021-45449](../security/index.md#cve-2021-45449) that affects users currently on Docker Desktop version 4.3.0 or 4.3.1.

Docker Desktop version 4.3.0 and 4.3.1 has a bug that may log sensitive information (access token or password) on the user's machine during login.
This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files.

### Upgrades

[docker scan v0.14.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.14.0)

### Security

**Log4j 2 CVE-2021-44228**: We have updated the `docker scan` CLI plugin.
This new version of `docker scan` is able to detect [Log4j 2
CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228) and [Log4j 2
CVE-2021-45046](https://nvd.nist.gov/vuln/detail/CVE-2021-45046)

For more information, read the blog post [Apache Log4j 2
CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/).

## 4.3.1

{{< release-date date="2021-12-11" >}}

### Upgrades

[docker scan v0.11.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.11.0)

### Security

**Log4j 2 CVE-2021-44228**: We have updated the `docker scan` CLI plugin for you.
Older versions of `docker scan` in Docker Desktop 4.3.0 and earlier versions are
not able to detect [Log4j 2
CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228).

For more information, read the
blog post [Apache Log4j 2
CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/).

## 4.3.0

{{< release-date date="2021-12-02" >}}

### Upgrades

- [Docker Engine v20.10.11](../engine/release-notes/20.10.md#201011)
- [containerd v1.4.12](https://github.com/containerd/containerd/releases/tag/v1.4.12)
- [Buildx 0.7.1](https://github.com/docker/buildx/releases/tag/v0.7.1)
- [Compose v2.2.1](https://github.com/docker/compose/releases/tag/v2.2.1)
- [Kubernetes 1.22.4](https://github.com/kubernetes/kubernetes/releases/tag/v1.22.4)
- [Docker Hub Tool v0.4.4](https://github.com/docker/hub-tool/releases/tag/v0.4.4)
- [Go 1.17.3](https://golang.org/doc/go1.17)

### Bug fixes and minor changes

#### For all platforms

- Added a self-diagnose warning if the host lacks Internet connectivity.
- Fixed an issue which prevented users from saving files from a volume using the Save As option in the Volumes UI. Fixes [docker/for-win#12407](https://github.com/docker/for-win/issues/12407).
- Docker Desktop now uses cgroupv2. If you need to run `systemd` in a container then:
  - Ensure your version of `systemd` supports cgroupv2. [It must be at least `systemd` 247](https://github.com/systemd/systemd/issues/19760#issuecomment-851565075). Consider upgrading any `centos:7` images to `centos:8`.
  - Containers running `systemd` need the following options: [`--privileged --cgroupns=host -v /sys/fs/cgroup:/sys/fs/cgroup:rw`](https://serverfault.com/questions/1053187/systemd-fails-to-run-in-a-docker-container-when-using-cgroupv2-cgroupns-priva).

#### For Mac

- Docker Desktop on Apple silicon no longer requires Rosetta 2, with the exception of [three optional command line tools](troubleshoot/known-issues.md).

#### For Windows

- Fixed an issue that caused Docker Desktop to fail during startup if the home directory path contains a character used in regular expressions. Fixes [docker/for-win#12374](https://github.com/docker/for-win/issues/12374).

### Known issue

Docker Dashboard incorrectly displays the container memory usage as zero on
Hyper-V based machines.
You can use the [`docker stats`](../engine/reference/commandline/container_stats.md)
command on the command line as a workaround to view the
actual memory usage. See
[docker/for-mac#6076](https://github.com/docker/for-mac/issues/6076).

### Deprecation

- The following internal DNS names are deprecated and will be removed from a future release: `docker-for-desktop`, `docker-desktop`, `docker.for.mac.host.internal`, `docker.for.mac.localhost`, `docker.for.mac.gateway.internal`. You must now use `host.docker.internal`, `vm.docker.internal`, and `gateway.docker.internal`.
- Removed: Custom RBAC rules have been removed from Docker Desktop as it gives `cluster-admin` privileges to all Service Accounts. Fixes [docker/for-mac/#4774](https://github.com/docker/for-mac/issues/4774).

## 4.2.0

{{< release-date date="2021-11-09" >}}

### New

**Pause/Resume**: You can now pause your Docker Desktop session when you are not actively using it and save CPU resources on your machine.

- Ships [Docker Public Roadmap#226](https://github.com/docker/roadmap/issues/226)

**Software Updates**: The option to turn off automatic check for updates is now available for users on all Docker subscriptions, including Docker Personal and Docker Pro. All update-related settings have been moved to the **Software Updates** section.

- Ships [Docker Public Roadmap#228](https://github.com/docker/roadmap/issues/228)

**Window management**: The Docker Dashboard window size and position persists when you close and reopen Docker Desktop.

### Upgrades

- [Docker Engine v20.10.10](../engine/release-notes/20.10.md#201010)
- [containerd v1.4.11](https://github.com/containerd/containerd/releases/tag/v1.4.11)
- [runc v1.0.2](https://github.com/opencontainers/runc/releases/tag/v1.0.2)
- [Go 1.17.2](https://golang.org/doc/go1.17)
- [Compose v2.1.1](https://github.com/docker/compose/releases/tag/v2.1.1)
- [docker-scan 0.9.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.9.0)

### Bug fixes and minor changes

#### For all platforms

- Improved: Self-diagnose now also checks for overlap between host IPs and `docker networks`.
- Fixed the position of the indicator that displays the availability of an update on the Docker Dashboard.

#### For Mac

- Fixed an issue that caused Docker Desktop to stop responding upon clicking **Exit** on the fatal error dialog.
- Fixed a rare startup failure affecting users having a `docker volume` bind-mounted on top of a directory from the host. If existing, this fix will also remove manually user added `DENY DELETE` ACL entries on the corresponding host directory.
- Fixed a bug where a `Docker.qcow2` file would be ignored on upgrade and a fresh `Docker.raw` used instead, resulting in containers and images disappearing. Note that if a system has both files (due to the previous bug) then the most recently modified file will be used, to avoid recent containers and images disappearing again. To force the use of the old `Docker.qcow2`, delete the newer `Docker.raw` file. Fixes [docker/for-mac#5998](https://github.com/docker/for-mac/issues/5998).
- Fixed a bug where subprocesses could fail unexpectedly during shutdown, triggering an unexpected fatal error popup. Fixes [docker/for-mac#5834](https://github.com/docker/for-mac/issues/5834).

#### For Windows

- Fixed Docker Desktop sometimes hanging when clicking Exit in the fatal error dialog.
- Fixed an issue that frequently displayed the **Download update** popup when an update has been downloaded but hasn't been applied yet [docker/for-win#12188](https://github.com/docker/for-win/issues/12188).
- Fixed installing a new update killing the application before it has time to shut down.
- Fixed: Installation of Docker Desktop now works even with group policies preventing users to start prerequisite services (e.g. LanmanServer) [docker/for-win#12291](https://github.com/docker/for-win/issues/12291).

## 4.1.1

{{< release-date date="2021-10-12" >}}

### Bug fixes and minor changes

#### For Mac

> When upgrading from 4.1.0, the Docker menu does not change to **Update and restart** so you can just wait for the download to complete (icon changes) and then select **Restart**. This bug is fixed in 4.1.1, for future upgrades.

- Fixed a bug where a `Docker.qcow2` file would be ignored on upgrade and a fresh `Docker.raw` used instead, resulting in containers and images disappearing. If a system has both files (due to the previous bug), then the most recently modified file will be used to avoid recent containers and images disappearing again. To force the use of the old `Docker.qcow2`, delete the newer `Docker.raw` file. Fixes [docker/for-mac#5998](https://github.com/docker/for-mac/issues/5998).
- Fixed the update notification overlay sometimes getting out of sync between the **Settings** button and the **Software update** button in the Docker Dashboard.
- Fixed the menu entry to install a newly downloaded Docker Desktop update. When an update is ready to install, the **Restart** option changes to **Update and restart**.

#### For Windows

- Fixed a regression in WSL 2 integrations for some distros (e.g. Arch or Alpine). Fixes [docker/for-win#12229](https://github.com/docker/for-win/issues/12229)
- Fixed update notification overlay sometimes getting out of sync between the Settings button and the Software update button in the Dashboard.

## 4.1.0

{{< release-date date="2021-09-30" >}}

### New

- **Software Updates**: The Settings tab now includes a new section to help you manage Docker Desktop updates. The **Software Updates** section notifies you whenever there's a new update and allows you to download the update or view information on what's included in the newer version.
- **Compose V2** You can now specify whether to use Docker Compose V2 in the General settings.
- **Volume Management**: Volume management is now available for users on any subscription, including Docker Personal. Ships [Docker Public Roadmap#215](https://github.com/docker/roadmap/issues/215)

### Upgrades

- [Compose V2](https://github.com/docker/compose/releases/tag/v2.0.0)
- [Buildx 0.6.3](https://github.com/docker/buildx/releases/tag/v0.6.3)
- [Kubernetes 1.21.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.5)
- [Go 1.17.1](https://github.com/golang/go/releases/tag/go1.17.1)
- [Alpine 3.14](https://alpinelinux.org/posts/Alpine-3.14.0-released.html)
- [Qemu 6.1.0](https://wiki.qemu.org/ChangeLog/6.1)
- Base distro to debian:bullseye

### Bug fixes and minor changes

#### For Windows

- Fixed a bug related to anti-malware software triggering, self-diagnose avoids calling the `net.exe` utility.
- Fixed filesystem corruption in the WSL 2 Linux VM in self-diagnose. This can be caused by [microsoft/WSL#5895](https://github.com/microsoft/WSL/issues/5895).
- Fixed `SeSecurityPrivilege` requirement issue. See [docker/for-win#12037](https://github.com/docker/for-win/issues/12037).
- Fixed CLI context switch sync with UI. See [docker/for-win#11721](https://github.com/docker/for-win/issues/11721).
- Added the key `vpnKitMaxPortIdleTime` to `settings.json` to allow the idle network connection timeout to be disabled or extended.
- Fixed a crash on exit. See [docker/for-win#12128](https://github.com/docker/for-win/issues/12128).
- Fixed a bug where the CLI tools would not be available in WSL 2 distros.
- Fixed switching from Linux to Windows containers that was stuck because access rights on panic.log. See [for-win#11899](https://github.com/docker/for-win/issues/11899).

### Known Issues

#### For Windows

Docker Desktop may fail to start when upgrading to 4.1.0 on some WSL-based distributions such as ArchWSL. See [docker/for-win#12229](https://github.com/docker/for-win/issues/12229)

## 4.0.1

{{< release-date date="2021-09-13" >}}

### Upgrades

- [Compose V2 RC3](https://github.com/docker/compose/releases/tag/v2.0.0-rc.3)
  - Compose v2 is now hosted on github.com/docker/compose.
  - Fixed go panic on downscale using `compose up --scale`.
  - Fixed a race condition in `compose run --rm` while capturing exit code.

### Bug fixes and minor changes

#### For all platforms

- Fixed a bug where copy-paste was not available in the Docker Dashboard.

#### For Windows

- Fixed a bug where Docker Desktop would not start correctly with the Hyper-V engine. See [docker/for-win#11963](https://github.com/docker/for-win/issues/11963)

## 4.0.0

{{< release-date date="2021-08-31" >}}

### New

Docker has [announced](https://www.docker.com/blog/updating-product-subscriptions/) updates and extensions to the product subscriptions to increase productivity, collaboration, and added security for our developers and businesses.

The updated [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement) includes a change to the terms for **Docker Desktop**.

- Docker Desktop **remains free** for small businesses (fewer than 250 employees AND less than $10 million in annual revenue), personal use, education, and non-commercial open source projects.
- It requires a paid subscription (**Pro, Team, or Business**), for as little as $5 a month, for professional use in larger enterprises.
- The effective date of these terms is August 31, 2021. There is a grace period until January 31, 2022 for those that will require a paid subscription to use Docker Desktop.
- The Docker Pro and Docker Team subscriptions now **include commercial use** of Docker Desktop.
- The existing Docker Free subscription has been renamed **Docker Personal**.
- **No changes** to Docker Engine or any other upstream **open source** Docker or Moby project.

To understand how these changes affect you, read the [FAQs](https://www.docker.com/pricing/faq).
For more information, see [Docker subscription overview](../subscription/index.md).

### Upgrades

- [Compose V2 RC2](https://github.com/docker/compose-cli/releases/tag/v2.0.0-rc.2)
  - Fixed project name to be case-insensitive for `compose down`. See [docker/compose-cli#2023](https://github.com/docker/compose-cli/issues/2023)
  - Fixed non-normalized project name.
  - Fixed port merging on partial reference.
- [Kubernetes 1.21.4](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.4)

### Bug fixes and minor changes

#### For Mac

- Fixed a bug where SSH was not available for builds from git URL. Fixes [for-mac#5902](https://github.com/docker/for-mac/issues/5902)

#### For Windows

- Fixed a bug where the CLI tools would not be available in WSL 2 distros.
- Fixed a bug when switching from Linux to Windows containers due to access rights on `panic.log`. [for-win#11899](https://github.com/docker/for-win/issues/11899)
