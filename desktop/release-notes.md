---
description: Release notes for Docker Desktop for Mac, Linux, Windows
keywords: Docker desktop, release notes, linux, mac, windows
title: Docker Desktop release notes
toc_max: 2
redirect_from:
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

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank" rel="noopener" class="_"} to see what's coming next.

For frequently asked questions about Docker Desktop releases, see [FAQs](faqs/general.md/#releases)

## Docker Desktop 4.13.0
2022-10-19

> Download Docker Desktop
>
> {% include desktop-install.html %}

### New

- Two new security features have been introduced for Docker Business users, Settings Management and Enhanced Container Isolation. Read more about Docker Desktop’s new [Hardened Desktop security model](hardened-desktop/index.md).
- Added the new Dev Environments CLI `docker dev`, so you can create, list, and run Dev Envs via command line. Now it's easier to integrate Dev Envs into custom scripts.
- Docker Desktop can now be installed to any drive and folder using the `--installation-dir`. Partially addresses [docker/roadmap#94](https://github.com/docker/roadmap/issues/94).

### Upgrades

- [Docker Scan v0.21.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.21.0)
- [Go 1.19.2](https://github.com/golang/go/releases/tag/go1.19.2) to address [CVE-2022-2879](https://www.cve.org/CVERecord?id=CVE-2022-2879){: target="_blank" rel="noopener"}, [CVE-2022-2880](https://www.cve.org/CVERecord?id=CVE-2022-2880){: target="_blank" rel="noopener"} and  [CVE-2022-41715](https://www.cve.org/CVERecord?id= CVE-2022-41715){: target="_blank" rel="noopener"}
- Updated Docker Engine and Docker CLI to [v20.10.20](https://docs.docker.com/engine/release-notes/#201020),
  which contain mitigations against a Git vulnerability, tracked in [CVE-2022-39253](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-39253){:target="_blank" rel="noopener"},
  and updated handling of `image:tag@digest` image references, as well as a fix for [CVE-2022-36109](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-36109).
- [Docker Credential Helpers v0.7.0](https://github.com/docker/docker-credential-helpers/releases/tag/v0.7.0){: target="blank" rel="noopener" class=""}
- [Docker Compose v2.12.0](https://github.com/docker/compose/releases/tag/v2.12.0)
- [Kubernetes v1.25.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.25.2)
- [Qemu 7.0.0](https://wiki.qemu.org/ChangeLog/7.0) used for cpu emulation, inside the Docker Desktop VM.
- [Linux kernel 5.15.49](https://hub.docker.com/layers/docker/for-desktop-kernel/5.15.49-13422a825f833d125942948cf8a8688cef721ead/images/sha256-ebf1f6f0cb58c70eaa260e9d55df7c43968874d62daced966ef6a5c5cd96b493?context=explore)

### Bug fixes and minor changes

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
- Docker Desktop now installs credential helpers from Github releases. See [docker/for-win#10247](https://github.com/docker/for-win/issues/10247), [docker/for-win#12995](https://github.com/docker/for-win/issues/12995), [docker/for-mac#12399](https://github.com/docker/for-mac/issues/12399).
- Fixed an issue where users were logged out of Docker Desktop after 7 days.


#### For Mac

- Added **Hide**, **Hide others**, **Show all** menu items for Docker Desktop. See [docker/for-mac#6446](https://github.com/docker/for-mac/issues/6446).
- Fixed a bug which caused the application to be deleted when running the install utility from the installed application. Fixes [docker/for-mac#6442](https://github.com/docker/for-mac/issues/6442).
- By default Docker will not create the /var/run/docker.sock symlink on the host and use the docker-desktop CLI context instead.

#### For Linux
- Fixed a bug that prevented pushing images from the Dashboard

## Docker Desktop 4.12.0
2022-09-01

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/85629/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/85629/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/85629/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/85629/docker-desktop-4.12.0-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/85629/docker-desktop-4.12.0-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/85629/docker-desktop-4.12.0-x86_64.pkg.tar.zst)

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
  <div class="panel panel-default">
    <div class="panel-heading" role="tab" id="headingSeven">
      <h5 class="panel-title">
        <a role="button" data-toggle="collapse" data-parent="#accordion" href="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
          Checksums
          <i class="fa fa-chevron-down"></i>
        </a>
      </h5>
    </div>
    <div id="collapseSeven" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingSeven">
      <div class="panel-body">
      <li><b>Windows:</b> SHA-256 996a4c5fff5b80b707ecfc0121d7ebe70d96c0bd568f058fd96f32cdec0c10cf</li>
      <li><b>Mac Intel:</b> SHA-256 41085009458ba1741c6a86c414190780ff3b288879aa27821fc4a985d229653c</li>
      <li><b>Mac Arm:</b> SHA-256 7eb63b4819cd1f87c61d5e8f54613692e07fb203d81bcf8d66f5de55489d3b81</li>
      <li><b>Linux DEB:</b> SHA-256 4407023db032219d6ac6031f81da6389ab192d3d06084ee6dad1ba4f4c64a4fe</li>
      <li><b>Linux RPM:</b> SHA-256 05e91f2a9763089acdfe710140893cb096bec955bcd99279bbe3aea035d09bc5</li>
      <li><b>Linux Arch:</b> SHA-256 7c6b43c8ab140c755e6c8ce4ec494b3f5c4f3b0c1ab3cee8bfd0b6864f795d8a</li>
      </div>
    </div>
  </div>
</div>

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

### Upgrades

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

### Bug fixes and minor changes

#### For all platforms

- Compose V2 is now enabled after factory reset.
- Compose V2 is now enabled by default on new installations of Docker Desktop.
- Precedence order of environment variables in Compose is more consistent, and clearly [documented](../compose/envvars-precedence.md).
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
- Fixed a bug that caused Docker Desktop to run in the background even after a user had quit the application.  Fixes [https://github.com/docker/for-mac/issues/6440]
- Disabled both Virtualization Framework and VirtioFS for users running macOS < 12.5

#### For Windows

- Fixed a bug where versions displayed during an update could be incorrect. Fixes [for-win/issues#12822](https://github.com/docker/for-win/issues/12822).

### Security 

#### For all platforms
- Fix RCE via query parameters in the message-box route in the Electron client.
- Fix RCE via extension description/changelog which could be abused by a malicious extension.

#### For Windows
- Fixed a bypass for the `--no-windows-containers` installation flag which was introduced in version 4.11. This flag allows administrators to disable the use of Windows containers.
- Fixed the argument injection to the Docker Desktop installer which may result in local privilege escalation.

## Docker Desktop 4.11.1
2022-08-05

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/84025/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/84025/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/84025/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/84025/docker-desktop-4.11.1-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/84025/docker-desktop-4.11.1-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/84025/docker-desktop-4.11.1-x86_64.pkg.tar.zst)

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
  <div class="panel panel-default">
    <div class="panel-heading" role="tab" id="headingSeven">
      <h5 class="panel-title">
        <a role="button" data-toggle="collapse" data-parent="#accordion" href="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
          Checksums
          <i class="fa fa-chevron-down"></i>
        </a>
      </h5>
    </div>
    <div id="collapseSeven" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingSeven">
      <div class="panel-body">
      <li><b>Windows:</b> SHA-256 8af32948447ddab655455542f6a12c8d752642a2bd451e2a48f76398cfd872b0</li>
      <li><b>Mac Intel:</b> SHA-256 b2f4ad8fea37dfb7d9147f169a9ceab71d7d0d12ff912057c60b58c0e91aed35</li>
      <li><b>Mac Arm:</b> SHA-256 a7d84117bef83764cb9bf275cd01b8ba0c43f08dbfe4d4a7d4f05549cdd81f54</li>
      <li><b>Linux DEB:</b> SHA-256 8877443ded0dee19b1bacaa608bd81d4bb216b59ff5fc12c89489e9ac5b00e0f</li>
      <li><b>Linux RPM:</b> SHA-256 a4a12071cdb4c3a845711eec13b97b838ae088f85f81cb5dd0db51aa6b050ed5</li>
      <li><b>Linux Arch:</b> SHA-256 66bdf3b4eb3cd29e190cf660ede53d3e854a4ec823c2ea04a4a02a175203f880</li>
      </div>
    </div>
  </div>
</div>

### Bug fixes and minor changes

#### For all platforms

- Fixed regression preventing VM system locations (e.g. /var/lib/docker) from being bind mounted [for-mac/issues#6433](https://github.com/docker/for-mac/issues/6433)

#### For Windows

- Fixed `docker login` to private registries from WSL2 distro [docker/for-win#12871](https://github.com/docker/for-win/issues/12871)

## Docker Desktop 4.11.0
2022-07-28

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/83626/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/83626/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/83626/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/83626/docker-desktop-4.11.0-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/83626/docker-desktop-4.11.0-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/83626/docker-desktop-4.11.0-x86_64.pkg.tar.zst)

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
    <div class="panel panel-default">
      <div class="panel-heading" role="tab" id="headingSeven">
        <h5 class="panel-title">
          <a role="button" data-toggle="collapse" data-parent="#accordion" href="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
            Checksums
            <i class="fa fa-chevron-down"></i>
          </a>
        </h5>
      </div>
      <div id="collapseSeven" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingSeven">
        <div class="panel-body">
        <li><b>Windows:</b> SHA-256 48ca8cabe67aee94a934b4c0f97a5001e89cb66bbbf824924fbc8bed6a8c90d3</li>
        <li><b>Mac Intel:</b> SHA-256 295694d7c2df05e37ac0d27fe8be5af6295b1edc6fa00a00a47134a14d5d0b34</li>
        <li><b>Mac Arm:</b> SHA-256 9824103e3d5a7d01a4d7d8086210157e1cc02217cb9edd82fe4bf2d16c138c44</li>
        <li><b>Linux DEB:</b> SHA-256 a0dc8ac97cc21e5a13a9e316cac11d85b7c248fd0c166b22a2ab239d17d43d9f</li>
        <li><b>Linux RPM:</b> SHA-256 eb077737298827092b283d3c85edacd128ecd993e987aa30d8081e2306401774</li>
        <li><b>Linux Arch:</b> SHA-256 a85fd5e83d5b613ef43d335c0ab0af4600aeb8a92921b617cb7a555826e361de</li>
        </div>
      </div>
    </div>
  </div>

### New

- Docker Desktop is now fully supported for Docker Business customers inside VMware ESXi and Azure VMs. For more information, see [Run Docker Desktop inside a VM or VDI environment](../desktop/vm-vdi.md)
- Added two new extensions ([vcluster](https://hub.docker.com/extensions/loftsh/vcluster-dd-extension) and [PGAdmin4](https://hub.docker.com/extensions/mochoa/pgadmin4-docker-extension)) to the Extensions Marketplace.
- The ability to sort extensions has been added to the Extensions Marketplace.
- Fixed a bug that caused some users to be asked for feedback too frequently. You'll now only be asked for feedback twice a year.
- Added custom theme settings for Docker Desktop. This allows you to specify dark or light mode for Docker Desktop independent of your device settings. Fixes [docker/for-win#12747](https://github.com/docker/for-win/issues/12747)
- Added a new flag for Windows installer. `--no-windows-containers` disables the Windows containers integration.
- Added a new flag for Mac install command. `--user <username>` sets up Docker Desktop for a specific user, preventing them from needing an admin password on first run.

### Upgrades
- [Docker Compose v2.7.0](https://github.com/docker/compose/releases/tag/v2.7.0)
- [Docker Compose "Cloud Integrations" v1.0.28](https://github.com/docker/compose-cli/releases/tag/v1.0.28)
- [Kubernetes v1.24.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.2)
- [Go 1.18.4](https://github.com/golang/go/releases/tag/go1.18.4)

### Bug fixes and minor changes

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


## Docker Desktop 4.10.1
2022-07-05

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/82475/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/82475/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/82475/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/82475/docker-desktop-4.10.1-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/82475/docker-desktop-4.10.1-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/82475/docker-desktop-4.10.1-x86_64.pkg.tar.zst)

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
    <div class="panel panel-default">
      <div class="panel-heading" role="tab" id="headingSeven">
        <h5 class="panel-title">
          <a role="button" data-toggle="collapse" data-parent="#accordion" href="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
            Checksums
            <i class="fa fa-chevron-down"></i>
          </a>
        </h5>
      </div>
      <div id="collapseSeven" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingSeven">
        <div class="panel-body">
        <li><b>Windows:</b> SHA-256 fe430d19d41cc56fd9a4cd2e22fc0e3522bed910c219208345918c77bbbd2a65</li>
        <li><b>Mac Intel:</b> SHA-256 8be8e5245d6a8dbf7b8cb580fb7d99f04cc143c95323695c0d9be4f85dd60b0e</li>
        <li><b>Mac Arm:</b> SHA-256 b3d4ef222325bde321045f3b8d946c849cd2812e9ad52a801000a95edb8af57b</li>
        <li><b>Linux DEB:</b> SHA-256 9363bc584478c5c7654004bacb51429c275b58a868ef43c3bc6249d5844ec5be</li>
        <li><b>Linux RPM:</b> SHA-256 92371d1a1ae4b57921721da95dc0252aefa4c79eb12208760c800ac07c0ae1d2</li>
        <li><b>Linux Arch:</b> SHA-256 799af244b05e8b08f03b6e0dbbc1dfcc027ff49f15506b3c460e0f9bae06ca5d</li>
        </div>
      </div>
    </div>
  </div>

### Bug fixes and minor changes

#### For Windows

- Fixed a bug where actions in the UI failed with Compose apps that were created from WSL. Fixes [docker/for-win#12806](https://github.com/docker/for-win/issues/12806).

#### For Mac
- Fixed a bug where the install command failed because paths were not initialized. Fixes [docker/for-mac#6384](https://github.com/docker/for-mac/issues/6384).


## Docker Desktop 4.10.0
2022-06-30

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/82025/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/82025/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/82025/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/82025/docker-desktop-4.10.0-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/82025/docker-desktop-4.10.0-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/82025/docker-desktop-4.10.0-x86_64.pkg.tar.zst)

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
    <div class="panel panel-default">
      <div class="panel-heading" role="tab" id="headingSeven">
        <h5 class="panel-title">
          <a role="button" data-toggle="collapse" data-parent="#accordion" href="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
            Checksums
            <i class="fa fa-chevron-down"></i>
          </a>
        </h5>
      </div>
      <div id="collapseSeven" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingSeven">
        <div class="panel-body">
        <li><b>Windows:</b> SHA-256 10615f4425e59eef7a22ce79ec13e41057df278547aa81c9fe4d623a848e80d8</li>
        <li><b>Mac Intel:</b> SHA-256 07bfe00296b724e4e772e268217bc8169a8b23ad98e6da419b13ebfe31b54643</li>
        <li><b>Mac Arm:</b> SHA-256 c9d2e72e5438726ab5a94c227d9130a65719f8fd09b877860ca2dcd86cfc188e</li>
        <li><b>Linux DEB:</b> SHA-256 c5f10b3d902b4ea10c8f75c17ba174e8838fc75889f76bc27abcab6afaf1969c</li>
        <li><b>Linux RPM:</b> SHA-256 a8ad3f8d4e93dfb6f28559f7dc84b7652e651fd6a49506e18958f1e69b51d9be</li>
        <li><b>Linux Arch:</b> SHA-256 37131c48df6436c1066c41ec0beda039e726e33bee689f751648c473f4abd96e</li>
        </div>
      </div>
    </div>
  </div>

### New

- You can now add environment variables before running an image in Docker Desktop.
- Added features to make it easier to work with a container's logs, such as regular expression search and the ability to clear container logs while the container is still running.
- Implemented feedback on the containers table. Added ports and separated container and image names.
- Added two new extensions, Ddosify and Lacework, to the Extensions Marketplace.

### Removed

- Removed Homepage while working on a new design. You can provide [feedback here](https://docs.google.com/forms/d/e/1FAIpQLSfYueBkJHdgxqsWcQn4VzBn2swu4u_rMQRIMa8LExYb_72mmQ/viewform?entry.1237514594=4.10).

### Upgrades
- [Docker Engine v20.10.17](../engine/release-notes/index.md#201017)
- [Docker Compose v2.6.1](https://github.com/docker/compose/releases/tag/v2.6.1)
- [Kubernetes v1.24.1](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.1)
- [cri-dockerd to v0.2.1](https://github.com/Mirantis/cri-dockerd/releases/tag/v0.2.1)
- [CNI plugins to v1.1.1](https://github.com/containernetworking/plugins/releases/tag/v1.1.1)
- [containerd to v1.6.6](https://github.com/containerd/containerd/releases/tag/v1.6.6)
- [runc to v1.1.2](https://github.com/opencontainers/runc/releases/tag/v1.1.2)
- [Go 1.18.3](https://github.com/golang/go/releases/tag/go1.18.3)

### Bug fixes and minor changes

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

## Docker Desktop 4.9.1
2022-06-16

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/81317/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/81317/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/81317/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/81317/docker-desktop-4.9.1-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/81317/docker-desktop-4.9.1-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/81317/docker-desktop-4.9.1-x86_64.pkg.tar.zst)

### Bug fixes and minor changes

#### For all platforms

- Fixed blank dashboard screen. Fixes [docker/for-win#12759](https://github.com/docker/for-win/issues/12759).

## Docker Desktop 4.9.0
2022-06-02

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/80466/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/80466/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/80466/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/80466/docker-desktop-4.9.0-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/80466/docker-desktop-4.9.0-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/80466/docker-desktop-4.9.0-x86_64.pkg.tar.zst)

### New

- Added additional guides on the homepage for: Elasticsearch, MariaDB, Memcached, MySQL, RabbitMQ and Ubuntu.
- Added a footer to the Docker Dashboard with general information about the Docker Desktop update status and Docker Engine statistics
- Re-designed the containers table, adding:
  - A button to copy a container ID to the clipboard
  - A pause button for each container
  - Column resizing for the containers table
  - Persistence of sorting and resizing for the containers table
  - Bulk deletion for the containers table

### Upgrades

- [Compose v2.6.0](https://github.com/docker/compose/releases/tag/v2.6.0)
- [Docker Engine v20.10.16](../engine/release-notes/index.md#201016)
- [containerd v1.6.4](https://github.com/containerd/containerd/releases/tag/v1.6.4)
- [runc v1.1.1](https://github.com/opencontainers/runc/releases/tag/v1.1.1)
- [Go 1.18.2](https://github.com/golang/go/releases/tag/go1.18.2)

### Bug fixes and minor changes

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

## Docker Desktop 4.8.2
2022-05-18

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/79419/Docker%20Desktop%20Installer.exe)|
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/79419/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/79419/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/79419/docker-desktop-4.8.2-amd64.deb) |
> [RPM](https://desktop.docker.com/linux/main/amd64/79419/docker-desktop-4.8.2-x86_64.rpm) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/79419/docker-desktop-4.8.2-x86_64.pkg.tar.zst)

### Upgrades

- [Compose v2.5.1](https://github.com/docker/compose/releases/tag/v2.5.1)

### Bug fixes and minor changes

- Fixed an issue with manual proxy settings which caused problems when pulling images. Fixes [docker/for-win#12714](https://github.com/docker/for-win/issues/12714) and [docker/for-mac#6315](https://github.com/docker/for-mac/issues/6315).
- Fixed high CPU usage when extensions are disabled. Fixes [docker/for-mac#6310](https://github.com/docker/for-mac/issues/6310).
- Docker Desktop now redacts HTTP proxy passwords in log files and diagnostics.

### Known issues

#### For Linux

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## Docker Desktop 4.8.1
2022-05-09

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/78998/Docker%20Desktop%20Installer.exe)|
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/78998/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/78998/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.1-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [RPM](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.1-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.1-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64)

### New

- Released [Docker Desktop for Linux](install/linux-install.md).
- Beta release of [Docker Extensions](/extensions/index.md) and Extensions SDK.
- Created a Docker Homepage where you can run popular images and discover how to use them.
- [Compose V2 is now GA](https://www.docker.com/blog/announcing-compose-v2-general-availability/)

### Bug fixes and minor changes

- Fixed a bug that caused the Kubernetes cluster to be deleted when updating Docker Desktop.

### Known issues

#### For Linux

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## Docker Desktop 4.8.0
2022-05-06

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/78933/Docker%20Desktop%20Installer.exe)|
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/78933/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/78933/Docker.dmg) |
> [Debian](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [RPM](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64)

### New

- Released [Docker Desktop for Linux](install/linux-install.md).
- Beta release of [Docker Extensions](/extensions/index.md) and Extensions SDK.
- Created a Docker Homepage where you can run popular images and discover how to use them.
- [Compose V2 is now GA](https://www.docker.com/blog/announcing-compose-v2-general-availability/)

### Upgrades

- [Compose v2.5.0](https://github.com/docker/compose/releases/tag/v2.5.0)
- [Go 1.18.1](https://github.com/golang/go/releases/tag/go1.18.1)
- [Kubernetes 1.24](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.0)

### Bug fixes and minor changes

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

## Docker Desktop 4.7.1
2022-04-19

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/77678/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/77678/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/77678/Docker.dmg)

### Bug fixes and minor changes

#### For all platforms 

 - Fixed a crash on the Quick Start Guide final screen.

#### For Windows

 - Fixed a bug where update was failing with a symlink error. Fixes [docker/for-win#12650](https://github.com/docker/for-win/issues/12650).
 - Fixed a bug that prevented using Windows container mode. Fixes [docker/for-win#12652](https://github.com/docker/for-win/issues/12652).


## Docker Desktop 4.7.0
2022-04-07

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/77141/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/77141/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/77141/Docker.dmg)

### Security

- Update Docker Engine to v20.10.14 to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769){: target="_blank" rel="noopener" class="_"}
- Update containerd to v1.5.11 to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769)

### New

- IT Administrators can now install Docker Desktop remotely using the command line.
- Add  the Docker Software Bill of Materials (SBOM) CLI plugin. The new CLI plugin enables users to generate SBOMs for Docker images. For more information, see [Docker SBOM](../engine/sbom/index.md).
- Use [cri-dockerd](https://github.com/Mirantis/cri-dockerd){: target="_blank" rel="noopener" class="_"}  for new Kubernetes clusters instead of `dockershim`. The change is transparent from the user's point of view and Kubernetes containers run on the Docker Engine as before. `cri-dockerd` allows Kubernetes to manage Docker containers using the standard [Container Runtime Interface](https://github.com/kubernetes/cri-api#readme){: target="_blank" rel="noopener" class="_"}, the same interface used to control other container runtimes. For more information, see [The Future of Dockershim is cri-dockerd](https://www.mirantis.com/blog/the-future-of-dockershim-is-cri-dockerd/){: target="_blank" rel="noopener" class="_"}.

### Upgrades

- [Docker Engine v20.10.14](../engine/release-notes/index.md#201014)
- [Compose v2.4.1](https://github.com/docker/compose/releases/tag/v2.4.1)
- [Buildx 0.8.2](https://github.com/docker/buildx/releases/tag/v0.8.2)
- [containerd v1.5.11](https://github.com/containerd/containerd/releases/tag/v1.5.11)
- [Go 1.18](https://golang.org/doc/go1.18)

### Bug fixes and minor changes

#### For all platforms 
 - Fixed a bug where the Registry Access Management policy was never refreshed after a failure.
 - Logs and terminals in the UI now respect your OS theme in light and dark mode.
 - Easily clean up many volumes at once via multi-select checkboxes.
 - Improved login feedback.

#### For Mac 

- Fixed an issue that sometimes caused Docker Desktop to display a blank white screen. Fixes [docker/for-mac#6134](https://github.com/docker/for-mac/issues/6134).
- Fixed a problem where gettimeofday() performance drops after waking from sleep when using Hyperkit. Fixes [docker/for-mac#3455](https://github.com/docker/for-mac/issues/3455).
- Fixed an issue that caused Docker Desktop to become unresponsive during startup when osxfs is used for file sharing.

#### For Windows

 - Fixed volume title. Fixes [docker/for-win#12616](https://github.com/docker/for-win/issues/12616).
 - Fixed a bug in the WSL 2 integration that caused Docker commands to stop working after restarting Docker Desktop or after switching to Windows containers.

## Docker Desktop 4.6.1
2022-03-22

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/76265/Docker%20Desktop%20Installer.exe)|
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/76265/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/76265/Docker.dmg)


### Upgrades

- [Buildx 0.8.1](https://github.com/docker/buildx/releases/tag/v0.8.1)

### Bug fixes and minor changes

- Prevented spinning in vpnkit-forwarder filling the logs with error messages.
- Fixed diagnostics upload when there is no HTTP proxy set. Fixes [docker/for-mac#6234](https://github.com/docker/for-mac/issues/6234).
- Removed a false positive "vm is not running" error from self-diagnose. Fixes [docker/for-mac#6233](https://github.com/docker/for-mac/issues/6233).

## Docker Desktop 4.6.0
2022-03-14

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/75818/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/75818/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/75818/Docker.dmg)

### Security

#### For all platforms 
- Fixed [CVE-2022-0847](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-0847){: target="_blank" rel="noopener" class="_"}, aka “Dirty Pipe”, an issue that could enable attackers to modify files in container images on the host, from inside a container.
  If using the WSL 2 backend, you must update WSL 2 by running `wsl --update`.

#### For Windows 

- Fixed [CVE-2022-26659](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-26659){: target="_blank" rel="noopener" class="_"}, which could allow an attacker to overwrite any administrator writable file on the system during the installation or the update of Docker Desktop.

### New

#### For all platforms 

- The Docker Dashboard Volume Management feature now offers the ability to efficiently clean up volumes using multi-select checkboxes.

#### For Mac

- Docker Desktop 4.6.0 gives macOS users the option of enabling a new experimental file sharing technology called VirtioFS. During testing VirtioFS has been shown to drastically reduce the time taken to sync changes between the host and VM, leading to substantial performance improvements. For more information, see [VirtioFS](settings/mac.md#beta-features).

### Upgrades

#### For all platforms 

- [Docker Engine v20.10.13](../engine/release-notes/index.md#201013)
- [Compose v2.3.3](https://github.com/docker/compose/releases/tag/v2.3.3)
- [Buildx 0.8.0](https://github.com/docker/buildx/releases/tag/v0.8.0)
- [containerd v1.4.13](https://github.com/containerd/containerd/releases/tag/v1.4.13)
- [runc v1.0.3](https://github.com/opencontainers/runc/releases/tag/v1.0.3)
- [Go 1.17.8](https://golang.org/doc/go1.17)
- [Linux kernel 5.10.104](https://hub.docker.com/layers/docker/for-desktop-kernel/5.10.104-379cadd2e08e8b25f932380e9fdaab97755357b3/images/sha256-7753b60f4544e5c5eed629d12151a49c8a4b48d98b4fb30e4e65cecc20da484d?context=explore)

#### For Mac

- [Qemu 6.2.0](https://wiki.qemu.org/ChangeLog/6.2)

### Bug fixes and minor changes

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

## Docker Desktop 4.5.1
2022-02-15

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/74721/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

#### For Windows 
- Fixed an issue that caused new installations to default to the Hyper-V backend instead of WSL 2.
- Fixed a crash in the Docker Dashboard which would make the systray menu disappear.

If you are running Docker Desktop on Windows Home, installing 4.5.1 will switch it back to WSL 2 automatically. If you are running another version of Windows, and you want Docker Desktop to use the WSL 2 backend, you must manually switch by enabling the **Use the WSL 2 based engine** option in the **Settings > General** section.
Alternatively, you can edit the Docker Desktop settings file located at `%APPDATA%\Docker\settings.json` and manually switch the value of the `wslEngineEnabled` field to `true`.

## Docker Desktop 4.5.0
2022-02-10

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/74594/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/74594/Docker.dmg)

### Security

#### For Mac

- Fixed [CVE-2021-44719](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44719){: target="_blank" rel="noopener" class="_"} where Docker Desktop could be used to access any user file on the host from a container, bypassing the allowed list of shared folders.

#### For Windows

- Fixed [CVE-2022-23774](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-23774){: target="_blank" rel="noopener" class="_"} where Docker Desktop allows attackers to move arbitrary files.

### New

- Docker Desktop 4.5.0 introduces a new version of the Docker menu which creates a consistent user experience across all operating systems. For more information, see the blog post [New Docker Menu & Improved Release Highlights with Docker Desktop 4.5](https://www.docker.com/blog/new-docker-menu-improved-release-highlights-with-docker-desktop-4-5/){: target="_blank" rel="noopener" class="_"}
- The 'docker version' output now displays the version of Docker Desktop installed on the machine.

### Upgrades

- [Amazon ECR Credential Helper v0.6.0](https://github.com/awslabs/amazon-ecr-credential-helper/releases/tag/v0.6.0){: target="blank" rel="noopener" class=""}

### Bug fixes and minor changes

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

## Docker Desktop 4.4.4 
2022-01-24

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/73704/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

#### For Windows 

- Fixed logging in from WSL 2. Fixes [docker/for-win#12500](https://github.com/docker/for-win/issues/12500).

### Known issues

#### For Windows

- Clicking **Proceed to Desktop** after signing in through the browser, sometimes does not bring the Dashboard to the front.
- After logging in, when the Dashboard receives focus, it sometimes stays in the foreground even when clicking a background window. As a workaround you need to click the Dashboard before clicking another application window.
- The tips of the week show on top of the mandatory login dialog when an organization restriction is enabled via a `registry.json` file.

## Docker Desktop 4.4.3
2022-01-14

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/73365/Docker%20Desktop%20Installer.exe)

### Bug fixes and minor changes

#### For Windows

- Disabled Dashboard shortcuts to prevent capturing them even when minimized or un-focussed. Fixes [docker/for-win#12495](https://github.com/docker/for-win/issues/12495).

### Known issues

#### For Windows

- Clicking **Proceed to Desktop** after signing in through the browser, sometimes does not bring the Dashboard to the front.
- After logging in, when the Dashboard receives focus, it sometimes stays in the foreground even when clicking a background window. As a workaround you need to click the Dashboard before clicking another application window.
- The tips of the week show on top of the mandatory login dialog when an organization restriction is enabled via a `registry.json` file.

## Docker Desktop 4.4.2
2022-01-13

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/73305/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/73305/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/73305/Docker.dmg)

### Security

- Fixed [CVE-2021-45449](../security/index.md#cve-2021-45449) that affects users currently on Docker Desktop version 4.3.0 or 4.3.1.

Docker Desktop version 4.3.0 and 4.3.1 has a bug that may log sensitive information (access token or password) on the user's machine during login.
This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files.

### New

- Easy, Secure sign in with Auth0 and Single Sign-on
  - Single Sign-on: Users with a Docker Business subscription can now configure SSO to authenticate using their identity providers (IdPs) to access Docker. For more information, see [Single Sign-on](../single-sign-on/index.md).
  - Signing in to Docker Desktop now takes you through the browser so that you get all the benefits of auto-filling from password managers.

### Upgrades

- [Docker Engine v20.10.12](../engine/release-notes/index.md#201012)
- [Compose v2.2.3](https://github.com/docker/compose/releases/tag/v2.2.3)
- [Kubernetes 1.22.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.22.5)
- [docker scan v0.16.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.16.0){: target="_blank" rel="noopener" class="_"}

### Bug fixes and minor changes

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

## Docker Desktop 4.3.2
2021-12-21

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/72729/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/72729/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/72729/Docker.dmg)

### Security

- Fixed [CVE-2021-45449](../security/index.md#cve-2021-45449) that affects users currently on Docker Desktop version 4.3.0 or 4.3.1.

Docker Desktop version 4.3.0 and 4.3.1 has a bug that may log sensitive information (access token or password) on the user's machine during login.
This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files.

### Upgrades

[docker scan v0.14.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.14.0){: target="_blank" rel="noopener" class="_"}

### Security

**Log4j 2 CVE-2021-44228**: We have updated the `docker scan` CLI plugin.
This new version of `docker scan` is able to detect [Log4j 2
CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228){:
target="_blank" rel="noopener" class="_"} and [Log4j 2
CVE-2021-45046](https://nvd.nist.gov/vuln/detail/CVE-2021-45046)

For more information, read the blog post [Apache Log4j 2
CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/){: target="_blank" rel="noopener" class="_"}.

## Docker Desktop 4.3.1
2021-12-11

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/72247/Docker%20Desktop%20Installer.exe) |
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
blog post [Apache Log4j 2
CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/){: target="_blank" rel="noopener" class="_"}.

## Docker Desktop 4.3.0
2021-12-02

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/71786/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/71786/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/71786/Docker.dmg)


### Upgrades

- [Docker Engine v20.10.11](../engine/release-notes/index.md#201011)
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
  - Containers running `systemd` need the following options: [`--privileged
    --cgroupns=host -v
    /sys/fs/cgroup:/sys/fs/cgroup:rw`](https://serverfault.com/questions/1053187/systemd-fails-to-run-in-a-docker-container-when-using-cgroupv2-cgroupns-priva).

#### For Mac

- Docker Desktop on Apple silicon no longer requires Rosetta 2, with the exception of [three optional command line tools](mac/apple-silicon.md#known-issues).

#### For Windows

- Fixed an issue that caused Docker Desktop to fail during startup if the home directory path contains a character used in regular expressions. Fixes [docker/for-win#12374](https://github.com/docker/for-win/issues/12374).

### Known issue

Docker Dashboard incorrectly displays the container memory usage as zero on
Hyper-V based machines.
You can use the [`docker stats`](../engine/reference/commandline/stats.md)
command on the command line as a workaround to view the
actual memory usage. See
[docker/for-mac#6076](https://github.com/docker/for-mac/issues/6076).

### Deprecation

- The following internal DNS names are deprecated and will be removed from a future release: `docker-for-desktop`, `docker-desktop`, `docker.for.mac.host.internal`, `docker.for.mac.localhost`, `docker.for.mac.gateway.internal`. You must now use `host.docker.internal`, `vm.docker.internal`, and `gateway.docker.internal`.
- Removed: Custom RBAC rules have been removed from Docker Desktop as it gives `cluster-admin` privileges to all Service Accounts. Fixes [docker/for-mac/#4774](https://github.com/docker/for-mac/issues/4774).

## Docker Desktop 4.2.0
2021-11-09

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/70708/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/70708/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/70708/Docker.dmg)

### New

**Pause/Resume**: You can now pause your Docker Desktop session when you are not actively using it and save CPU resources on your machine.

- Ships [Docker Public Roadmap#226](https://github.com/docker/roadmap/issues/226){: target="_blank" rel="noopener" class="_"}

**Software Updates**: The option to turn off automatic check for updates is now available for users on all Docker subscriptions, including Docker Personal and Docker Pro. All update-related settings have been moved to the **Software Updates** section. 

- Ships [Docker Public Roadmap#228](https://github.com/docker/roadmap/issues/228){: target="_blank" rel="noopener" class="_"}

**Window management**: The Docker Dashboard window size and position persists when you close and reopen Docker Desktop.

### Upgrades

- [Docker Engine v20.10.10](../engine/release-notes/index.md#201010)
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


## Docker Desktop 4.1.1
2021-10-12

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/69879/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/69879/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/69879/Docker.dmg)

### Bug fixes and minor changes

#### For Mac

> When upgrading from 4.1.0, the Docker menu does not change to **Update and restart** so you can just wait for the download to complete (icon changes) and then select **Restart**. This bug is fixed in 4.1.1, for future upgrades.

- Fixed a bug where a `Docker.qcow2` file would be ignored on upgrade and a fresh `Docker.raw` used instead, resulting in containers and images disappearing. If a system has both files (due to the previous bug), then the most recently modified file will be used to avoid recent containers and images disappearing again. To force the use of the old `Docker.qcow2`, delete the newer `Docker.raw` file. Fixes [docker/for-mac#5998](https://github.com/docker/for-mac/issues/5998).
- Fixed the update notification overlay sometimes getting out of sync between the **Settings** button and the **Software update** button in the Docker Dashboard.
- Fixed the menu entry to install a newly downloaded Docker Desktop update. When an update is ready to install, the **Restart** option changes to **Update and restart**.

#### For Windows

- Fixed a regression in WSL 2 integrations for some distros (e.g. Arch or Alpine). Fixes [docker/for-win#12229](https://github.com/docker/for-win/issues/12229)
- Fixed update notification overlay sometimes getting out of sync between the Settings button and the Software update button in the Dashboard.

## Docker Desktop 4.1.0
2021-09-30

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/69386/Docker%20Desktop%20Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/69386/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/69386/Docker.dmg)

### New

- **Software Updates**: The Settings tab now includes a new section to help you manage Docker Desktop updates. The **Software Updates** section notifies you whenever there's a new update and allows you to download the update or view information on what's included in the newer version. 
- **Compose V2** You can now specify whether to use Docker Compose V2 in the General settings.
- **Volume Management**: Volume management is now available for users on any subscription, including Docker Personal. Ships [Docker Public Roadmap#215](https://github.com/docker/roadmap/issues/215){: target="_blank" rel="noopener" class="_"}

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

## Docker Desktop 4.0.1
2021-09-13

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/68347/Docker Desktop Installer.exe) |
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/68347/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/68347/Docker.dmg)

### Upgrades

- [Compose V2 RC3](https://github.com/docker/compose/releases/tag/v2.0.0-rc.3)
  - Compose v2 is now hosted on github.com/docker/compose.
  - Fixed go panic on downscale using `compose up --scale`.
  - Fixed  a race condition in `compose run --rm` while capturing exit code.

### Bug fixes and minor changes

#### For all platforms

- Fixed a bug where copy-paste was not available in the Docker Dashboard.

#### For Windows

- Fixed a bug where Docker Desktop would not start correctly with the Hyper-V engine. See [docker/for-win#11963](https://github.com/docker/for-win/issues/11963)

## Docker Desktop 4.0.0
2021-08-31

> Download Docker Desktop
>
> [Windows](https://desktop.docker.com/win/main/amd64/67817/Docker Desktop Installer.exe) |
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

To understand how these changes affect you, read the [FAQs](https://www.docker.com/pricing/faq){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"}.
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
