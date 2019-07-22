---
title: Docker Desktop Enterprise release notes
description: Release notes for Docker Desktop Enterprise
keywords: Docker Desktop Enterprise, Windows, Mac, Docker Desktop, Enterprise,
toc_min: 1
toc_max: 2
---

This page contains information about the new features, improvements, known issues, and bug fixes in the Docker Desktop Enterprise release. Documentation is updated for each release. We suggest that you regularly visit this page to learn about updates.

For information on system requirements, installation, and download, see:

- [Install Docker Desktop Enterprise on Mac](/ee/desktop/admin/install/mac)
- [Install Docker Desktop Enterprise on Windows](/ee/desktop/admin/install/windows)

For Docker Enterprise Engine release notes, see [Docker Engine release notes](/engine/release-notes).

## Version 2.1.0.1
2019-07-22

Docker Desktop Enterprise 2.1.0.1 contains a Kubernetes upgrade. Note that your local Kubernetes cluster will be reset after installing this release.

### Upgrades

- [Docker 19.03.0](https://docs.docker.com/engine/release-notes/) in Version Pack Enterprise 3.0
- [Kubernetes 1.14.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.14.3) in Version Pack Enterprise 3.0
- [Compose on Kubernetes 0.4.23](https://github.com/docker/compose-on-kubernetes/releases/tag/v0.4.23) in Version Pack Enterprise 3.0
- [Docker Compose 1.24.1](https://github.com/docker/compose/releases/tag/1.24.1)
- [Docker 18.09.8](https://docs.docker.com/engine/release-notes/) in Version Pack Enterprise 2.1
- [Docker 17.06.2-ee-23](https://docs.docker.com/engine/release-notes/) in Version Pack Enterprise 2.0
- [Docker Credential Helpers 0.6.3](https://github.com/docker/docker-credential-helpers/releases/tag/v0.6.3)
- [Application Designer 0.1.4](/ee/desktop/app-designer/)

### Known issue

When you sign out of Windows while Docker Desktop is still running, and then sign in and start Docker Desktop, attempts to run Docker commands that rely on network connections can fail. For example, the command `docker pull alpine` returns `Error response from daemon: Get https://registry-1.docker.io/v2/: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)`.

Note that this issue occurs only when you sign out of Windows and not when you shut down or restart Windows.

Workaround: After signing back into Windows, when Docker Desktop has started, right-click the Docker menu from the systray and then select Quit Docker Desktop. When this is complete, start Docker Desktop normally.

## Version 2.0.0.6
2019-06-06

### Upgrades

- [Docker 19.03.0-rc2](https://docs.docker.com/engine/release-notes/) in Version Pack Enterprise 3.0
- Application Designer 0.1.3

### Bug fixes and minor changes

- Application Designer/ Docker Template CLI: The `DefaultServer` preference has been renamed to `DefaultRegistry`.
- Windows: Docker Desktop now allows administrators to specify authorized named pipes that can be mounted on Windows containers.
- Windows: Docker Desktop checks for stored credentials at startup before attempting to mount any shared drives. This prompts users to reenter the credentials if they are invalid.
- Mac: Docker Deskop does not send DNS queries for `docker-desktop.<domain>` every 10s. It relies on the host's DNS domain search order rather than trying to replicate it inside the VM.
- Mac: Docker Desktop now uses a separate thread to move its storage to a different disk. This allows the UI to remain responsive during the disk move. Fixes [docker/for-mac#3592](https://github.com/docker/for-mac/issues/3592)

## Version 2.0.0.5
2019-05-30

### Upgrades

- [Docker 19.03.0-rc1](https://docs.docker.com/engine/release-notes/) in Enterprise 3.0 version pack
- Application Designer 0.1.2
- [Qemu 4.0.0](https://github.com/docker/binfmt) to cross-compile ARM devices

### Bug fixes and minor changes

- Application Designer now allows users to copy and paste application logs.
- Users can browse the scaffold logs when scaffolding a new application using the Application Designer.
- Application Designer allows users to set preferences, including the default organization and registry.
- Docker Desktop admins can enforce Application Designer preferences using the `admin-settings.json` file.
- Security improvements: Docker Desktop now checks TLS certificates for the target endpoints when using `kubectl`.
- Fixed an issue where Visual Studio Code IDE path was not detected properly.

## Version 2.0.0.4
2019-05-16

### Upgrades

- [Docker 19.03.0-beta4](https://docs.docker.com/engine/release-notes/) in Enterprise 3.0 version pack
- [Docker 18.09.6](https://docs.docker.com/engine/release-notes/), [Kubernetes 1.11.10](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.11.md#v11110) in Enterprise 2.1 version pack
- [LinuxKit v0.7](https://github.com/linuxkit/linuxkit/releases/tag/v0.7)

### Bug fixes and minor changes

- Fixed a stability issue with the DNS resolver.
- Fixed a race condition where Kubernetes sometimes failed to start after restarting the application.
- Fixed a bug that causes Docker Compose to fail when a user logs out after logging in. See [docker/compose#6517](https://github.com/docker/compose/issues/6517)
- Improved the reliability of `com.docker.osxfs trace` performance profiling command.
- Docker Desktop now supports large lists of resource DNS records on Mac. See [docker/for-mac#2160](https://github.com/docker/for-mac/issues/2160#issuecomment-431571031).
- Users can now run a Docker registry in a container. See [docker/for-mac#3611](https://github.com/docker/for-mac/issues/3611).
- For Linux containers on Windows (LCOW), one physical computer system running Windows 10 Professional or Windows 10 Enterprise version 1809 or later is required.
- Added a dialog box during startup when a shared drive fails to mount. This allows users to retry mounting the drive or remove it from the shared drive list.
- Removed the ability to log in using an email address as a username as this is not supported by the Docker command line.

## Version 2.0.0.3
2019-04-26

### Upgrades

- [Docker Engine 19.03.0-beta2](https://docs.docker.com/engine/release-notes/) for Version Pack Enterprise 3.0.

## Version 2.0.0.2
2019-04-19

**WARNING:** You must upgrade the previously installed Version Packs to the latest revision.

### New

- Version Pack Enterprise 3.0 with [Docker Engine 19.03.0-beta1](https://docs.docker.com/engine/release-notes/) and [Kubernetes 1.14.1](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.14.md#changelog-since-v1141)

- Application Designer now includes new templates for AngularJS and VueJS.

### Upgrades

- [Docker Compose 1.24.0](https://github.com/docker/compose/releases/tag/1.24.0)
- [Docker Engine 18.09.5](https://docs.docker.com/engine/release-notes/), [Kubernetes 1.11.7](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.11.md#v1117) and [Compose on Kubernetes 0.4.22](https://github.com/docker/compose-on-kubernetes/releases/tag/v0.4.22) for Version Pack Enterprise 2.1
- [Docker Engine 17.06.2-ee-21](https://docs.docker.com/engine/release-notes/) for Version Pack Enterprise 2.0

### Bug fixes and minor changes

- For security, only administrators can install or upgrade Version Packs using the `dockerdesktop-admin` tool.
- Truncate UDP DNS responses which are over 512 bytes in size
- Fixed airgap install of kubernetes in version pack enterprise-2.0
- Reset to factory default now resets to admin defaults

### Known issues

- The Docker Template CLI plugin included in this version is an outdated version of the plugin and will fail when scaffolding templates. Note that the Application Designer is not affected by this outdated version of the CLI plugin.

## Version 2.0.0.1
2019-03-01

**WARNING:** You must upgrade the previously installed Version Packs to the latest revision.

### Windows

#### Upgrades

- Docker 18.09.3 for Version Pack Enterprise 2.1, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)
- Docker 17.06.2-ee-20 for Version Pack Enterprise 2.0, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

#### Bug fixes and minor changes

- Fixed port 8080 that was used on localhost when starting Kubernetes.
- Fixed Hub login through the desktop UI not sync with login through `docker login` command line.
- Fixed crash in system tray menu when the Hub login fails or Air gap mode.

### Mac

#### New features

- Added ability to list all installed version packs with the admin CLI command `dockerdesktop-admin version-pack list`.
- `dockerdesktop-admin app uninstall` will also remove Docker Desktop user files.

#### Upgrades

- Docker 18.09.3 for Version Pack Enterprise 2.1, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)
- Docker 17.06.2-ee-20 for Version Pack Enterprise 2.0, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

#### Bug fixes and minor changes

- Fixed port 8080 that was used on localhost when starting Kubernetes.
- Improved error messaging to suggest running diagnostics / resetting to factory default only when it is appropriate.

## Version 2.0.0.0
2019-01-31

### New features

- **Version selection**: Configurable version packs ensure the local
    instance of Docker Desktop Enterprise is a precise copy of the
    production environment where applications are deployed, and
    developers can switch between versions of Docker and
    Kubernetes with a single click.

- **Application Designer**: Application templates allow you to choose a
    technology and focus on business logic. Updates can be made with
    minimal syntax knowledge.

- **Device management**: The Docker Desktop Enterprise installer is available as standard MSI (Win) and PKG (Mac) downloads, which allows administrators to script an installation across many developer machines.

- **Administrative control**: IT organizations can specify and lock configuration parameters for creation of a standardized development environment, including disabling drive sharing and limiting version pack installations. Developers run commands in the command line without worrying about configuration settings.
