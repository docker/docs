---
description: Answers to common questions about Docker Desktop release frequency,
  supported versions, installing older versions, and download errors
keywords: desktop, mac, windows, faqs, releases, support, install, 403, version,
  older version, old version, previous version, downgrade
title: FAQs on Docker Desktop releases
linkTitle: Releases
tags: [FAQ]
---

### How often will new releases be?

New releases are available every week, unless there are critical fixes that
need to be released sooner.

The **Automatically check for updates** setting in the **Software updates**
tab is turned on by default. This means you receive notifications in the
Docker menu and a notification badge on the Docker Desktop Dashboard when a
new version is available.

On Mac and Windows, you can also let Docker Desktop automatically download
new updates in the background by selecting the **Always download updates**
checkbox.

Sometimes new versions are rolled out gradually over a few days. Therefore,
if you wait, it will turn up soon. Alternatively, you can select **Check for
updates** in the Docker menu to get the latest version immediately.

### How do I ensure that all users in my organization are using the same version?

This is managed through your IT administrator's endpoint management software.

### Which Docker Desktop versions are supported?

> [!IMPORTANT]
>
> Docker Desktop versions older than six months from the latest release aren't
> available for download. Previous release notes are available in the
> [documentation repository](https://github.com/docker/docs/tree/main/content/manuals/desktop/previous-versions).

Docker Support covers:

- Docker Business: Versions up to six months older than the latest version.
  Fixes apply to the latest version only.
- Docker Pro and Team: The latest version only.

For details, see [Supported versions](/manuals/support/_index.md#supported-versions).

Docker recommends running the latest Docker Desktop release. See the
[release notes](/manuals/desktop/release-notes.md) for the latest version and
download links.

### How do I install an older version?

Installers are limited to the
[six-month download window](#which-docker-desktop-versions-are-supported).
Choose a version that still includes **Download Docker Desktop** links in the
[release notes](/manuals/desktop/release-notes.md).

The install commands don't take a version number. Choose the version by
downloading that release's installer, then run the install steps for your
operating system.

#### Find a specific version

1. Open the [release notes](/manuals/desktop/release-notes.md).
1. Go to the version you want. Under the version heading and date, look for
   a **Download Docker Desktop** section with links for Windows, Mac, and
   Linux.
1. If that section is missing, Docker doesn't provide an installer for that
   entry from the docs. Choose another version that includes **Download Docker
   Desktop**.
1. Download the installer for your platform.

If a newer version of Docker Desktop is already installed,
[uninstall](/manuals/desktop/uninstall.md) it before installing the older
version.

#### Install the version you downloaded

{{< tabs >}}
{{< tab name="Windows" >}}

1. Download `Docker Desktop Installer.exe` for the target version from that
   release's **Download Docker Desktop** links.
1. In PowerShell, run the installer for a per-user install:

   ```powershell
   Start-Process 'Docker Desktop Installer.exe' -Wait -ArgumentList 'install', '--user'
   ```

   For an all-users install (requires administrator privileges):

   ```powershell
   Start-Process 'Docker Desktop Installer.exe' -Wait -ArgumentList 'install'
   ```

For more options, see
[Install from the command line](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line).

{{< /tab >}}
{{< tab name="macOS" >}}

1. Download `Docker.dmg` for the target version (Apple chip or Intel chip)
   from that release's **Download Docker Desktop** links.
1. In a terminal, install from the DMG you downloaded:

   ```console
   $ sudo hdiutil attach Docker.dmg
   $ sudo /Volumes/Docker/Docker.app/Contents/MacOS/install
   $ sudo hdiutil detach /Volumes/Docker
   ```

For more options, see
[Install from the command line](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line).

{{< /tab >}}
{{< tab name="Linux" >}}

1. Download the `.deb` or `.rpm` package for the target version from that
   release's **Download Docker Desktop** links.
1. Install the package file you downloaded.

   Ubuntu / Debian:

   ```console
   $ sudo apt-get install ./docker-desktop-amd64.deb
   ```

   Fedora:

   ```console
   $ sudo dnf install ./docker-desktop-x86_64.rpm
   ```

Package filenames vary by distribution. For distribution-specific steps, see
[Install Docker Desktop on Linux](/manuals/desktop/setup/install/linux/_index.md).

{{< /tab >}}
{{< /tabs >}}

### Why do I get a 403 download error?

A `403 Forbidden` response means the installer for that build isn't available
to download. This usually means the build is outside the
[six-month download window](#which-docker-desktop-versions-are-supported).
This applies on Windows, Mac, and Linux, whether you use a browser, a package
manager, or another download tool. The failure happens when fetching the
installer, not when running the install command.

Download a version that includes a **Download Docker Desktop** section in the
[release notes](/manuals/desktop/release-notes.md).
