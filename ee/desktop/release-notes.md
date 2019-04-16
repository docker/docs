---
title: Docker Desktop Enterprise release notes
description: Release notes for Docker Desktop Enterprise
keywords: Docker Desktop Enterprise, Windows, Mac, Docker Desktop, Enterprise,
---

This topic contains information about the main improvements and issues, starting with the
current release. The documentation is updated for each release.

For information on system requirements, installation, and download, see:

- [Install Docker Desktop Enterprise on Mac](/ee/desktop/admin/install/mac)
- [Install Docker Desktop Enterprise on Windows](/ee/desktop/admin/install/windows)

For Docker Enterprise Engine release notes, see [Docker Engine release notes](/engine/release-notes).

# Docker Desktop Enterprise Releases of 2019

## Docker Desktop Enterprise 2.0.0.1

2019-03-01

**WARNING:** You must upgrade the previously installed Version Packs to the latest revision.

### Windows

Upgrades:

- Docker 18.09.3 for Version Pack Enterprise 2.1, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

- Docker 17.06.2-ee-20 for Version Pack Enterprise 2.0, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

Bug fixes and minor changes:

- Fixed port 8080 that was used on localhost when starting Kubernetes.
- Fixed Hub login through the desktop UI not sync with login through `docker login` command line.
- Fixed crash in system tray menu when the Hub login fails or Air gap mode.

### Mac

New features:

- Added ability to list all installed version packs with the admin CLI command `dockerdesktop-admin version-pack list`.
- `dockerdesktop-admin app uninstall` will also remove Docker Desktop user files.

 Upgrades:

- Docker 18.09.3 for Version Pack Enterprise 2.1, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

- Docker 17.06.2-ee-20 for Version Pack Enterprise 2.0, fixes [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

Bug fixes and minor changes:

- Fixed port 8080 that was used on localhost when starting Kubernetes.
- Improved error messaging to suggest running diagnostics / resetting to factory default only when it is appropriate.

## Docker Desktop Enterprise 2.0.0.0

2019-01-31

New features:

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
