---
description: Change log / release notes for Docker Desktop for Mac
keywords: Docker Desktop for Mac, release notes
title: Docker Desktop for Mac release notes
toc_min: 1
toc_max: 2
redirect_from:
- /docker-for-mac/release-notes/
- /docker-for-mac/edge-release-notes/
- /mackit/release-notes/
---

> **Update to the Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) now requires a paid
> subscription. The grace period for those that will require a paid subscription
> ends on January 31, 2022. [Learn more](https://www.docker.com/blog/the-grace-period-for-the-docker-subscription-service-agreement-ends-soon-heres-what-you-need-to-know/){:
 target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}.
{: .important}

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank" rel="noopener" class="_"} to see what's coming next.

## Docker Desktop 4.8.1
2022-05-09

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-amd64){: .button .primary-btn }
> [Mac with Apple
> chip](https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64){:
> .button .primary-btn }

### New

- Released [Docker Desktop for Linux](https://docs.docker.com/desktop/linux/).
- Beta release of [Docker Extensions](https://docs.docker.com/desktop/extensions) and Extensions SDK.
- Created a Docker Homepage where you can run popular images and discover how to use them.
- [Compose V2 is now GA](https://www.docker.com/blog/announcing-compose-v2-general-availability/)

### Bug fixes and minor changes

- Fixed a bug that caused the Kubernetes cluster to be deleted when updating Docker Desktop.

## Docker Desktop 4.8.0
2022-05-06

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/78933/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/78933/Docker.dmg)

### New

- Released [Docker Desktop for Linux](https://docs.docker.com/desktop/linux/).
- Beta release of [Docker Extensions](https://docs.docker.com/desktop/extensions) and Extensions SDK.
- Created a Docker Homepage where you can run popular images and discover how to use them.
- [Compose V2 is now GA](https://www.docker.com/blog/announcing-compose-v2-general-availability/)

### Upgrades

- [Compose v2.5.0](https://github.com/docker/compose/releases/tag/v2.5.0)
- [Go 1.18.1](https://golang.org/doc/go1.18)
- [Kubernetes 1.24](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.0)

### Bug fixes and minor changes

- Introduced reading system proxy. You no longer need to manually configure proxies unless it differs from your OS level proxy.
- Fixed a bug that showed Remote Repositories in the Dashboard when running behind a proxy.
- Fixed vpnkit establishing and blocking the client connection even if the server is gone. See [docker/for-mac#6235](https://github.com/docker/for-mac/issues/6235).
- Made improvements on the Volume tab in Docker Desktop:
  - Volume size is displayed.
  - Columns can be resized, hidden and reordered.
  - A columns sort order and hidden state is persisted, even after Docker Desktop restarts.
  - Row selection is persisted when switching between tabs, even after Docker Desktop restarts.
- Fixed a bug in the Dev Environments tab that did not add a scroll when more items were added to the screen.
- Standardised the header title and action in the Dashboard.
- Added support for downloading Registry Access Management policies through HTTP proxies.
- Fixed an issue related to empty remote repositories when the machine is in sleep mode for an extended period of time.
- Fixed a bug where dangling images were not selected in the cleanup process if their name was not marked as "<none>" but their tag is.
- Docker Desktop's icon now matches Big Sur Style guide. See [docker/for-mac#5536](https://github.com/docker/for-mac/issues/5536).
- Fixed a problem with duplicate Dock icons and Dock icon not working as expected. Fixes [docker/for-mac#6189](https://github.com/docker/for-mac/issues/6189).
- Improved the error message when `docker pull` fails because an HTTP proxy is required.
- Added the ability to clear the search bar easily in Docker Desktop. 
- Renamed the "Containers / Apps" tab to "Containers".
- Improved support for the `Cmd+Q` shortcut.

### Known issues

- Currently, if you are running a Kubernetes cluster, it will be deleted when you upgrade to Docker Desktop 4.8.0. We aim to fix this in the next release.


## Docker Desktop 4.7.1
2022-04-19

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/77678/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/77678/Docker.dmg)
### Bug fixes and minor changes
 - Fixed a crash on the Quick Start Guide final screen.

## Docker Desktop 4.7.0
2022-04-07

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/77141/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/77141/Docker.dmg)

### Security

- Update Docker Engine to v20.10.14 to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769){: target="_blank" rel="noopener" class="_"}
- Update containerd to v1.5.11 to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769)

### New

- IT Administrators can now install Docker Desktop remotely using the command line. For more information, see [Install Docker Desktop using the command line](../install.md#install-from-the-command-line).
- Add  the Docker Software Bill of Materials (SBOM) CLI plugin. The new CLI plugin enables users to generate SBOMs for Docker images. For more information, see [Docker SBOM](../../../engine/sbom/index.md).
- Use [cri-dockerd](https://github.com/Mirantis/cri-dockerd){: target="_blank" rel="noopener" class="_"}  for new Kubernetes clusters instead of `dockershim`. The change is transparent from the user's point of view and Kubernetes containers run on the Docker Engine as before. `cri-dockerd` allows Kubernetes to manage Docker containers using the standard [Container Runtime Interface](https://github.com/kubernetes/cri-api#readme){: target="_blank" rel="noopener" class="_"}, the same interface used to control other container runtimes. For more information, see [The Future of Dockershim is cri-dockerd](https://www.mirantis.com/blog/the-future-of-dockershim-is-cri-dockerd/){: target="_blank" rel="noopener" class="_"}.

### Upgrades

- [Docker Engine v20.10.14](https://docs.docker.com/engine/release-notes/#201014)
- [Compose v2.4.1](https://github.com/docker/compose/releases/tag/v2.4.1)
- [Buildx 0.8.2](https://github.com/docker/buildx/releases/tag/v0.8.2)
- [containerd v1.5.11](https://github.com/containerd/containerd/releases/tag/v1.5.11)
- [Go 1.18](https://golang.org/doc/go1.18)

### Bug fixes and minor changes
 - Fixed an issue that sometimes caused Docker Desktop to display a blank white screen. Fixes [docker/for-mac#6134](https://github.com/docker/for-mac/issues/6134).
 - Fixed a bug where the Registry Access Management policy was never refreshed after a failure.
 - Fixed an issue that caused Docker Desktop to become unresponsive during startup when osxfs is used for file sharing.
 - Fixed a problem where gettimeofday() performance drops after waking from sleep when using Hyperkit. Fixes [docker/for-mac#3455](https://github.com/docker/for-mac/issues/3455).
 - Logs and terminals in the UI now respect your OS theme in light and dark mode.
 - Easily clean up many volumes at once via multi-select checkboxes.
 - Improved login feedback.

### Deprecation
- Deprecate macOS 10.14. For macOS system requirements, see [Docker Desktop for Mac system requirements](../install.md#system-requirements).

## Docker Desktop 4.6.1
2022-03-22

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/76265/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/76265/Docker.dmg)

### Upgrades

- [Buildx 0.8.1](https://github.com/docker/buildx/releases/tag/v0.8.1)

### Bug fixes and minor changes

- Fixed diagnostics upload when there is no HTTP proxy set. Fixes [docker/for-mac#6234](https://github.com/docker/for-mac/issues/6234).
- Removed a false positive "vm is not running" error from self-diagnose. Fixes [docker/for-mac#6233](https://github.com/docker/for-mac/issues/6233).
- Prevented spinning in vpnkit-forwarder filling the logs with error messages.

## Docker Desktop 4.6.0
2022-03-14

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/75818/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/75818/Docker.dmg)

### Security

- Fixed [CVE-2022-0847](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-0847){: target="_blank" rel="noopener" class="_"}, aka “Dirty Pipe”, an issue that could enable attackers to modify files in container images on the host, from inside a container.

### New

- Docker Desktop 4.6.0 gives macOS users the option of enabling a new experimental file sharing technology called VirtioFS. During testing VirtioFS has been shown to drastically reduce the time taken to sync changes between the host and VM, leading to substantial performance improvements. For more information, see [VirtioFS](../index.md#experimental-features).
- The Docker Dashboard Volume Management feature now offers the ability to efficiently clean up volumes using multi-select checkboxes.

### Upgrades

- [Docker Engine v20.10.13](https://docs.docker.com/engine/release-notes/#201013)
- [Compose v2.3.3](https://github.com/docker/compose/releases/tag/v2.3.3)
- [Buildx 0.8.0](https://github.com/docker/buildx/releases/tag/v0.8.0)
- [containerd v1.4.13](https://github.com/containerd/containerd/releases/tag/v1.4.13)
- [runc v1.0.3](https://github.com/opencontainers/runc/releases/tag/v1.0.3)
- [Go 1.17.8](https://golang.org/doc/go1.17)
- [Linux kernel 5.10.104](https://hub.docker.com/layers/docker/for-desktop-kernel/5.10.104-379cadd2e08e8b25f932380e9fdaab97755357b3/images/sha256-7753b60f4544e5c5eed629d12151a49c8a4b48d98b4fb30e4e65cecc20da484d?context=explore)
- [Qemu 6.2.0](https://wiki.qemu.org/ChangeLog/6.2)

### Bug fixes and minor changes

- Fixed uploading diagnostics when an HTTPS proxy is set.
- Fixed the systray menu not displaying all menu items after starting Docker Desktop. Fixes [docker/for-mac#6192](https://github.com/docker/for-mac/issues/6192).
- Fixed a regression about Docker Desktop not starting in background anymore. Fixes [docker/for-mac#6167](https://github.com/docker/for-mac/issues/6167).
- Fixed missing Docker Desktop Dock icon. Fixes [docker/for-mac#6173](https://github.com/docker/for-mac/issues/6173).
- Made checking for updates from the systray menu open the Software updates settings section.
- Used speed up block device access when using the experimental `virtualization.framework`. See [benchmarks](https://github.com/docker/roadmap/issues/7#issuecomment-1050626886).
- Increased default VM memory allocation to half of physical memory (min 2 GB, max 8 GB) for better out-of-the-box performances.
- Fixed copying text from terminal in Quick Start Guide. Fixes [docker/for-win#12444](https://github.com/docker/for-win/issues/12444).

### Known issues

- After enabling VirtioFS, containers with processes running with different Unix user IDs may experience caching issues. For example if a process running as `root` queries a file and another process running as user `nginx` tries to access the same file immediately, the `nginx` process will get a "Permission Denied" error.

## Docker Desktop 4.5.0
2022-02-10

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/74594/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/74594/Docker.dmg)

### Security

- Fixed [CVE-2021-44719](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44719){: target="_blank" rel="noopener" class="_"} where Docker Desktop could be used to access any user file on the host from a container, bypassing the allowed list of shared folders.

### New

- Docker Desktop 4.5.0 introduces a new version of the Docker menu which creates a consistent user experience across all operating systems. For more information, see the blog post [New Docker Menu & Improved Release Highlights with Docker Desktop 4.5](https://www.docker.com/blog/new-docker-menu-improved-release-highlights-with-docker-desktop-4-5/){: target="_blank" rel="noopener" class="_"}
- The 'docker version' output now displays the version of Docker Desktop installed on the machine.

### Upgrades

- [Amazon ECR Credential Helper v0.6.0](https://github.com/awslabs/amazon-ecr-credential-helper/releases/tag/v0.6.0){: target="blank" rel="noopener" class=""}

### Bug fixes and minor changes

- Fixed an issue that caused the VM to become unresponsive during startup when using `osxfs` and when no host directories are shared with the VM.
- Increased the filesystem watch (inotify) limits by setting `fs.inotify.max_user_watches=1048576` and `fs.inotify.max_user_instances=8192` in Linux. Fixes [docker/for-mac#6071](https://github.com/docker/for-mac/issues/6071).
- Fixed an issue that didn't allow users to stop a Docker Compose application using Docker Dashboard if the application was started in a different version of Docker Compose. For example, if the user started a Docker Compose application in V1 and then switched to Docker Compose V2, attempts to stop the Docker Compose application would fail.
- Fixed an issue where Docker Desktop incorrectly prompted users to sign in after they quit Docker Desktop and start the application.
- Fixed an issue where the **About Docker Desktop** window wasn't working anymore.
- Limit the number of CPUs to 8 on Mac M1 to fix the startup problem. Fixes [docker/for-mac#6063](https://github.com/docker/for-mac/issues/6063).

## Docker Desktop 4.4.2
2022-01-13

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/73305/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/73305/Docker.dmg)

### Security

- Fixed [CVE-2021-45449](https://docs.docker.com/security/#cve-2021-45449) that affects users currently on Docker Desktop version 4.3.0 or 4.3.1.

Docker Desktop version 4.3.0 and 4.3.1 has a bug that may log sensitive information (access token or password) on the user's machine during login.
This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files.

### New

- Easy, Secure sign in with Auth0 and Single Sign-on
  - Single Sign-on: Users with a Docker Business subscription can now configure SSO to authenticate using their identity providers (IdPs) to access Docker. For more information, see [Single Sign-on](../../../single-sign-on/index.md).
  - Signing in to Docker Desktop now takes you through the browser so that you get all the benefits of auto-filling from password managers.

### Upgrades

- [Docker Engine v20.10.12](https://docs.docker.com/engine/release-notes/#201012)
- [Compose v2.2.3](https://github.com/docker/compose/releases/tag/v2.2.3)
- [Kubernetes 1.22.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.22.5)
- [docker scan v0.16.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.16.0){: target="_blank" rel="noopener" class="_"}

### Bug fixes and minor changes

- Docker Desktop displays an error if `registry.json` contains more than one organization in the `allowedOrgs` field. If you are using multiple organizations for different groups of developers, you must provision a separate `registry.json` file for each group.
- Fixed the memory statistics for containers in the Dashboard. Fixes [docker/for-mac/#4774](https://github.com/docker/for-mac/issues/6076).
- Added a deprecated option to `settings.json`: `"deprecatedCgroupv1": true`, which switches the Linux environment back to cgroups v1. If your software requires cgroups v1, you should update it to be compatible with cgroups v2. Although cgroups v1 should continue to work, it is likely that some future features will depend on cgroups v2. It is also possible that some Linux kernel bugs will only be fixed with cgroups v2.
- Fixed a regression in Compose that reverted the container name separator from `-` to `_`. Fixes [docker/compose-switch](https://github.com/docker/compose-switch/issues/24).
- Fixed an issue where putting the machine to Sleep mode after pausing Docker Desktop results in Docker Desktop not being able to resume from pause after the machine comes out of Sleep mode. Fixes [for-mac#6058](https://github.com/docker/for-mac/issues/6058).

### Known issues

- The tips of the week show on top of the mandatory login dialog when an organization restriction is enabled via a `registry.json` file.

## Docker Desktop 4.3.2
2021-12-21

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/72729/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/72729/Docker.dmg)

### Security

- Fixed [CVE-2021-45449](https://docs.docker.com/security/#cve-2021-45449) that affects users currently on Docker Desktop version 4.3.0 or 4.3.1.

Docker Desktop version 4.3.0 and 4.3.1 has a bug that may log sensitive information (access token or password) on the user's machine during login.
This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files.

### Upgrades

[docker scan v0.14.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.14.0){: target="_blank" rel="noopener" class="_"}

### Security

**Log4j 2 CVE-2021-44228**: We have updated the `docker scan` CLI plugin.
This new version of `docker scan` is able to detect [Log4j 2
CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228){:
target="_blank" rel="noopener" class="_"} and [Log4j 2
CVE-2021-45046](https://nvd.nist.gov/vuln/detail/CVE-2021-45046).

For more information, read the blog post [Apache Log4j 2
CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/){: target="_blank" rel="noopener" class="_"}.

## Docker Desktop 4.3.1
2021-12-11

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/72247/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/72247/Docker.dmg)

### Upgrades

[docker scan v0.11.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.11.0){: target="_blank" rel="noopener" class="_"}

### Security

**Log4j 2 CVE-2021-44228**: We have updated the `docker scan` CLI plugin for you.
Older versions of `docker scan` in Docker Desktop 4.3.0 and earlier versions are
not able to detect [Log4j 2
CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228){:
target="_blank" rel="noopener" class="_"}.

For more information, read the
blog post [Apache Log4j 2 CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/){: target="_blank" rel="noopener" class="_"}.

## Docker Desktop 4.3.0
2021-12-02

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/71786/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/71786/Docker.dmg)

### Upgrades

- [Docker Engine v20.10.11](https://docs.docker.com/engine/release-notes/#201011)
- [containerd v1.4.12](https://github.com/containerd/containerd/releases/tag/v1.4.12)
- [Buildx 0.7.1](https://github.com/docker/buildx/releases/tag/v0.7.1)
- [Compose v2.2.1](https://github.com/docker/compose/releases/tag/v2.2.1)
- [Kubernetes 1.22.4](https://github.com/kubernetes/kubernetes/releases/tag/v1.22.4)
- [Docker Hub Tool v0.4.4](https://github.com/docker/hub-tool/releases/tag/v0.4.4)
- [Go 1.17.3](https://golang.org/doc/go1.17)

### Bug fixes and minor changes

- Docker Desktop on Apple silicon no longer requires Rosetta 2, with the exception of [three optional command line tools](../apple-silicon/#known-issues).
- Fixed an issue which prevented users from saving files from a volume using the Save As option in the Volumes UI. Fixes [docker/for-win#12407](https://github.com/docker/for-win/issues/12407).
- Fixed an issue that sometimes launched the existing version of Docker Desktop even after updating to a newer version.
- Added a self-diagnose warning if the host lacks Internet connectivity.
- Docker Desktop now uses cgroupv2. If you need to run `systemd` in a container then:
  - Ensure your version of `systemd` supports cgroupv2. [It must be at least `systemd` 247](https://github.com/systemd/systemd/issues/19760#issuecomment-851565075). Consider upgrading any `centos:7` images to `centos:8`.
  - Containers running `systemd` need the following options: [`--privileged --cgroupns=host -v /sys/fs/cgroup:/sys/fs/cgroup:rw`](https://serverfault.com/questions/1053187/systemd-fails-to-run-in-a-docker-container-when-using-cgroupv2-cgroupns-priva).

### Known issue

Docker Dashboard incorrectly displays the container memory usage as zero.
You can use the [`docker stats`](../../../engine/reference/commandline/stats.md)
command on the command line as a workaround to view the
actual memory usage. See
[docker/for-mac#6076](https://github.com/docker/for-mac/issues/6076).

### Deprecation

- The following internal DNS names are deprecated and will be removed from a future release: `docker-for-desktop`, `docker-desktop`, `docker.for.mac.host.internal`, `docker.for.mac.localhost`, `docker.for.mac.gateway.internal`. You must now use `host.docker.internal`, `vm.docker.internal`, and `gateway.docker.internal`.
- Custom RBAC rules have been removed from Docker Desktop as it gives `cluster-admin` privileges to all Service Accounts. Fixes [docker/for-mac/#4774](https://github.com/docker/for-mac/issues/4774).

## Docker Desktop 4.2.0
2021-11-09

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/70708/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/70708/Docker.dmg)

### New

**Pause/Resume**: You can now pause your Docker Desktop session when you are not actively using it and save CPU resources on your machine. For more information, see [Pause/Resume](../index.md#pauseresume).

- Ships [Docker Public Roadmap#226](https://github.com/docker/roadmap/issues/226){: target="_blank" rel="noopener" class="_"}

**Software Updates**: The option to turn off automatic check for updates is now available for users on all Docker subscriptions, including Docker Personal and Docker Pro. All update-related settings have been moved to the **Software Updates** section. For more information, see [Software updates](../index.md#software-updates).

- Ships [Docker Public Roadmap#228](https://github.com/docker/roadmap/issues/228){: target="_blank" rel="noopener" class="_"}

**Window management**: The Docker Dashboard window size and position persists when you close and reopen Docker Desktop.

### Upgrades

- [Docker Engine v20.10.10](https://docs.docker.com/engine/release-notes/#201010)
- [containerd v1.4.11](https://github.com/containerd/containerd/releases/tag/v1.4.11)
- [runc v1.0.2](https://github.com/opencontainers/runc/releases/tag/v1.0.2)
- [Go 1.17.2](https://golang.org/doc/go1.17)
- [Compose v2.1.1](https://github.com/docker/compose/releases/tag/v2.1.1)
- [docker-scan 0.9.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.9.0)

### Bug fixes and minor changes

- Improved: Self-diagnose now also checks for overlap between host IPs and `docker networks`.
- Fixed the position of the indicator that displays the availability of an update on the Docker Dashboard.
- Fixed an issue that caused Docker Desktop to stop responding upon clicking **Exit** on the fatal error dialog.
- Fixed a rare startup failure affecting users having a `docker volume` bind-mounted on top of a directory from the host. If existing, this fix will also remove manually user added `DENY DELETE` ACL entries on the corresponding host directory.
- Fixed a bug where a `Docker.qcow2` file would be ignored on upgrade and a fresh `Docker.raw` used instead, resulting in containers and images disappearing. Note that if a system has both files (due to the previous bug) then the most recently modified file will be used, to avoid recent containers and images disappearing again. To force the use of the old `Docker.qcow2`, delete the newer `Docker.raw` file. Fixes [docker/for-mac#5998](https://github.com/docker/for-mac/issues/5998).
- Fixed a bug where subprocesses could fail unexpectedly during shutdown, triggering an unexpected fatal error popup. Fixes [docker/for-mac#5834](https://github.com/docker/for-mac/issues/5834).


## Docker Desktop 4.1.1
2021-10-12

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/69879/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/69879/Docker.dmg)

### Bug fixes and minor changes

> When upgrading from 4.1.0, the Docker menu does not change to **Update and restart** so you can just wait for the download to complete (icon changes) and then select **Restart**. This bug is fixed in 4.1.1, for future upgrades.

- Fixed a bug where a `Docker.qcow2` file would be ignored on upgrade and a fresh `Docker.raw` used instead, resulting in containers and images disappearing. If a system has both files (due to the previous bug), then the most recently modified file will be used to avoid recent containers and images disappearing again. To force the use of the old `Docker.qcow2`, delete the newer `Docker.raw` file. Fixes [docker/for-mac#5998](https://github.com/docker/for-mac/issues/5998).
- Fixed the update notification overlay sometimes getting out of sync between the **Settings** button and the **Software update** button in the Docker Dashboard.
- Fixed the menu entry to install a newly downloaded Docker Desktop update. When an update is ready to install, the **Restart** option changes to **Update and restart**.

## Docker Desktop 4.1.0
2021-09-30

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/69386/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/69386/Docker.dmg)

### New

- **Software Updates**: The Settings tab now includes a new section to help you manage Docker Desktop updates. The **Software Updates** section notifies you whenever there's a new update and allows you to download the update or view information on what's included in the newer version. For more information, see [Software Updates](../index.md#software-updates).
- **Compose V2** You can now specify whether to use [Docker Compose V2](../../../compose/index.md#compose-v2-and-the-new-docker-compose-command) in the General settings.
- **Volume Management**: Volume management is now available for users on any subscription, including Docker Personal. For more information, see [Explore volumes](../../dashboard.md#explore-volumes). Ships [Docker Public Roadmap#215](https://github.com/docker/roadmap/issues/215){: target="_blank" rel="noopener" class="_"}

### Upgrades

- [Compose V2](https://github.com/docker/compose/releases/tag/v2.0.0)
- [Buildx 0.6.3](https://github.com/docker/buildx/releases/tag/v0.6.3)
- [Kubernetes 1.21.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.5)
- [Go 1.17.1](https://github.com/golang/go/releases/tag/go1.17.1)
- [Alpine 3.14](https://alpinelinux.org/posts/Alpine-3.14.0-released.html)
- [Qemu 6.1.0](https://wiki.qemu.org/ChangeLog/6.1)
- Base distro to debian:bullseye

## Docker Desktop 4.0.1
2021-09-13

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/68347/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/68347/Docker.dmg)

### Upgrades

- [Compose V2 RC3](https://github.com/docker/compose/releases/tag/v2.0.0-rc.3)
  - Compose v2 is now hosted on github.com/docker/compose.
  - Fixed go panic on downscale using `compose up --scale`.
  - Fixed  a race condition in `compose run --rm` while capturing exit code.

### Bug fixes and minor changes

- Fixed a bug where copy-paste was not available in the Docker Dashboard.

## Docker Desktop 4.0.0
2021-08-31

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/67817/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/67817/Docker.dmg)

### New

Docker has [announced](https://www.docker.com/blog/updating-product-subscriptions/){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"} updates and extensions to the product subscriptions to increase productivity, collaboration, and added security for our developers and businesses.

The updated [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement) includes a change to the terms for **Docker Desktop**.

- Docker Desktop **remains free** for small businesses (fewer than 250 employees AND less than $10 million in annual revenue), personal use, education, and non-commercial open source projects.
- It requires a paid subscription (**Pro, Team, or Business**), for as little as $5 a month, for professional use in larger enterprises.
- The effective date of these terms is August 31, 2021. There is a grace period until January 31, 2022 for those that will require a paid subscription to use Docker Desktop.
- The Docker Pro and Docker Team subscriptions now **include commercial use** of Docker Desktop.
- The existing Docker Free subscription has been renamed **Docker Personal**.
- **No changes** to Docker Engine or any other upstream **open source** Docker or Moby project.

To understand how these changes affect you, read the [FAQs](https://www.docker.com/pricing/faq){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"}. For more information, see [Docker subscription overview](../../../subscription/index.md).

### Upgrades

- [Compose V2 RC2](https://github.com/docker/compose-cli/releases/tag/v2.0.0-rc.2)
  - Fixed project name to be case-insensitive for `compose down`. See [docker/compose-cli#2023](https://github.com/docker/compose-cli/issues/2023)
  - Fixed non-normalized project name.
  - Fixed port merging on partial reference.
- [Kubernetes 1.21.4](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.4)

### Bug fixes and minor changes

- Fixed a bug where SSH was not available for builds from git URL. Fixes [for-mac#5902](https://github.com/docker/for-mac/issues/5902)
