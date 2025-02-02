---
description: Get started with Docker for Windows. This guide covers system requirements,
  where to download, and instructions on how to install and update.
keywords: docker for windows, docker windows, docker desktop for windows, docker on
  windows, install docker windows, install docker on windows, docker windows 10, docker
  run on windows, installing docker for windows, windows containers, wsl, hyper-v
title: Install Docker Desktop on Windows
linkTitle: Windows
weight: 30
aliases:
- /desktop/windows/install/
- /docker-ee-for-windows/install/
- /docker-for-windows/install-windows-home/
- /docker-for-windows/install/
- /ee/docker-ee/windows/docker-ee/
- /engine/installation/windows/
- /engine/installation/windows/docker-ee/
- /install/windows/docker-ee/
- /install/windows/ee-preview/
- /installation/windows/
- /desktop/win/configuring-wsl/
- /desktop/install/windows-install/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains the download URL, information about system requirements, and instructions on how to install Docker Desktop for Windows.

{{< button text="Docker Desktop for Windows - x86_64" url="https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-win-amd64" >}}
{{< button text="Docker Desktop for Windows - Arm (Beta)" url="https://desktop.docker.com/win/main/arm64/Docker%20Desktop%20Installer.exe?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-win-arm64" >}}

_For checksums, see [Release notes](/manuals/desktop/release-notes.md)_

## System requirements

> [!TIP]
>
> **Should I use Hyper-V or WSL?**
>
> Docker Desktop's functionality remains consistent on both WSL and Hyper-V, without a preference for either architecture. Hyper-V and WSL have their own advantages and disadvantages, depending on your specific set up and your planned use case. 

{{< tabs >}}
{{< tab name="WSL 2 backend, x86_64" >}}

- WSL version 1.1.3.0 or later.
- Windows 11 64-bit: Home or Pro version 22H2 or higher, or Enterprise or Education version 22H2 or higher.
- Windows 10 64-bit: Minimum required is Home or Pro 22H2 (build 19045) or higher, or Enterprise or Education 22H2 (build 19045) or higher. 
- Turn on the WSL 2 feature on Windows. For detailed instructions, refer to the
  [Microsoft documentation](https://docs.microsoft.com/en-us/windows/wsl/install-win10).
- The following hardware prerequisites are required to successfully run
  WSL 2 on Windows 10 or Windows 11:
  - 64-bit processor with [Second Level Address Translation (SLAT)](https://en.wikipedia.org/wiki/Second_Level_Address_Translation)
  - 4GB system RAM
  - Enable hardware virtualization in BIOS. For more information, see
    [Virtualization](/manuals/desktop/troubleshoot-and-support/troubleshoot/topics.md#virtualization).

For more information on setting up WSL 2 with Docker Desktop, see [WSL](/manuals/desktop/features/wsl/_index.md).

> [!NOTE]
>
> Docker only supports Docker Desktop on Windows for those versions of Windows that are still within [Microsoft’s servicing timeline](https://support.microsoft.com/en-us/help/13853/windows-lifecycle-fact-sheet). Docker Desktop is not supported on server versions of Windows, such as Windows Server 2019 or Windows Server 2022. For more information on how to run containers on Windows Server, see [Microsoft's official documentation](https://learn.microsoft.com/virtualization/windowscontainers/quick-start/set-up-environment).

> [!IMPORTANT]
>
> To run Windows containers, you need Windows 10 or Windows 11 Professional or Enterprise edition.
> Windows Home or Education editions only allow you to run Linux containers.

{{< /tab >}}
{{< tab name="Hyper-V backend, x86_64" >}}

- Windows 11 64-bit: Enterprise, Pro, or Education version 22H2 or higher.
- Windows 10 64-bit: Enterprise, Pro, or Education version 22H2 (build 19045) or higher.
- Turn on Hyper-V and Containers Windows features.
- The following hardware prerequisites are required to successfully run Client
  Hyper-V on Windows 10:

  - 64 bit processor with [Second Level Address Translation (SLAT)](https://en.wikipedia.org/wiki/Second_Level_Address_Translation)
  - 4GB system RAM
  - Turn on BIOS-level hardware virtualization support in the
    BIOS settings. For more information, see
    [Virtualization](/manuals/desktop/troubleshoot-and-support/troubleshoot/topics.md#virtualization).

> [!NOTE]
>
> Docker only supports Docker Desktop on Windows for those versions of Windows that are still within [Microsoft’s servicing timeline](https://support.microsoft.com/en-us/help/13853/windows-lifecycle-fact-sheet). Docker Desktop is not supported on server versions of Windows, such as Windows Server 2019 or Windows Server 2022. For more information on how to run containers on Windows Server, see [Microsoft's official documentation](https://learn.microsoft.com/virtualization/windowscontainers/quick-start/set-up-environment).

> [!IMPORTANT]
>
> To run Windows containers, you need Windows 10 or Windows 11 Professional or Enterprise edition.
> Windows Home or Education editions only allow you to run Linux containers.

{{< /tab >}}
{{< tab name="WSL 2 backend, Arm (Beta)" >}}

- WSL version 1.1.3.0 or later.
- Windows 11 64-bit: Home or Pro version 22H2 or higher, or Enterprise or Education version 22H2 or higher.
- Windows 10 64-bit: Minimum required is Home or Pro 22H2 (build 19045) or higher, or Enterprise or Education 22H2 (build 19045) or higher.
- Turn on the WSL 2 feature on Windows. For detailed instructions, refer to the
  [Microsoft documentation](https://docs.microsoft.com/en-us/windows/wsl/install-win10).
- The following hardware prerequisites are required to successfully run
  WSL 2 on Windows 10 or Windows 11:
  - 64-bit processor with [Second Level Address Translation (SLAT)](https://en.wikipedia.org/wiki/Second_Level_Address_Translation)
  - 4GB system RAM
  - Enable hardware virtualization in BIOS. For more information, see
    [Virtualization](/manuals/desktop/troubleshoot-and-support/troubleshoot/topics.md#virtualization).

> [!IMPORTANT]
>
> Windows containers are not supported.

{{< /tab >}}
{{< /tabs >}}

Containers and images created with Docker Desktop are shared between all
user accounts on machines where it is installed. This is because all Windows
accounts use the same VM to build and run containers. Note that it is not possible to share containers and images between user accounts when using the Docker Desktop WSL 2 backend.

Running Docker Desktop inside a VMware ESXi or Azure VM is supported for Docker Business customers.
It requires enabling nested virtualization on the hypervisor first.
For more information, see [Running Docker Desktop in a VM or VDI environment](/manuals/desktop/setup/vm-vdi.md).

{{< accordion title="How do I switch between Windows and Linux containers?" >}}

From the Docker Desktop menu, you can toggle which daemon (Linux or Windows)
the Docker CLI talks to. Select **Switch to Windows containers** to use Windows
containers, or select **Switch to Linux containers** to use Linux containers
(the default).

For more information on Windows containers, refer to the following documentation:

- Microsoft documentation on [Windows containers](https://docs.microsoft.com/en-us/virtualization/windowscontainers/about/index).

- [Build and Run Your First Windows Server Container (Blog Post)](https://www.docker.com/blog/build-your-first-docker-windows-server-container/)
  gives a quick tour of how to build and run native Docker Windows containers on Windows 10 and Windows Server 2016 evaluation releases.

- [Getting Started with Windows Containers (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
  shows you how to use the [MusicStore](https://github.com/aspnet/MusicStore/)
  application with Windows containers. The MusicStore is a standard .NET application and,
  [forked here to use containers](https://github.com/friism/MusicStore), is a good example of a multi-container application.

- To understand how to connect to Windows containers from the local host, see
  [I want to connect to a container from Windows](/manuals/desktop/features/networking.md#i-want-to-connect-to-a-container-from-the-host)

> [!NOTE]
>
> When you switch to Windows containers, **Settings** only shows those tabs that are active and apply to your Windows containers.

If you set proxies or daemon configuration in Windows containers mode, these
apply only on Windows containers. If you switch back to Linux containers,
proxies and daemon configurations return to what you had set for Linux
containers. Your Windows container settings are retained and become available
again when you switch back.

{{< /accordion >}}

## Install Docker Desktop on Windows

> [!TIP]
>
> See the [FAQs](/manuals/desktop/troubleshoot-and-support/faqs/general.md#how-do-i-run-docker-desktop-without-administrator-privileges) on how to install and run Docker Desktop without needing administrator privileges.

### Install interactively

1. Download the installer using the download button at the top of the page, or from the [release notes](/manuals/desktop/release-notes.md).

2. Double-click `Docker Desktop Installer.exe` to run the installer. By default, Docker Desktop is installed at `C:\Program Files\Docker\Docker`.

3. When prompted, ensure the **Use WSL 2 instead of Hyper-V** option on the Configuration page is selected or not depending on your choice of backend.

   If your system only supports one of the two options, you won't be able to select which backend to use.

4. Follow the instructions on the installation wizard to authorize the installer and proceed with the install.

5. When the installation is successful, select **Close** to complete the installation process.

6. [Start Docker Desktop](#start-docker-desktop).

If your administrator account is different to your user account, you must add the user to the **docker-users** group:
1. Run **Computer Management** as an **administrator**.
2. Navigate to **Local Users and Groups** > **Groups** > **docker-users**. 
3. Right-click to add the user to the group.
4. Sign out and sign back in for the changes to take effect.

### Install from the command line

After downloading `Docker Desktop Installer.exe`, run the following command in a terminal to install Docker Desktop:

```console
$ "Docker Desktop Installer.exe" install
```

If you’re using PowerShell you should run it as:

```powershell
Start-Process 'Docker Desktop Installer.exe' -Wait install
```

If using the Windows Command Prompt:

```sh
start /w "" "Docker Desktop Installer.exe" install
```

By default, Docker Desktop is installed at `C:\Program Files\Docker\Docker`.

The `install` command accepts the following flags:
- `--quiet`: Suppresses information output when running the installer 
- `--accept-license`: Accepts the [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement) now, rather than requiring it to be accepted when the application is first run
- `--no-windows-containers`: Disables the Windows containers integration. This can improve security. For more information, see [Windows containers](/manuals/desktop/setup/install/windows-permission-requirements.md#windows-containers).
- `--allowed-org=<org name>`: Requires the user to sign in and be part of the specified Docker Hub organization when running the application
- `--backend=<backend name>`: Selects the default backend to use for Docker Desktop, `hyper-v`, `windows` or `wsl-2` (default)
- `--installation-dir=<path>`: Changes the default installation location (`C:\Program Files\Docker\Docker`)
- `--admin-settings`: Automatically creates an `admin-settings.json` file which is used by admins to control certain Docker Desktop settings on client machines within their organization. For more information, see [Settings Management](/manuals/security/for-admins/hardened-desktop/settings-management/_index.md).
  - It must be used together with the `--allowed-org=<org name>` flag. 
  - For example:`--allowed-org=<org name> --admin-settings="{'configurationFileVersion': 2, 'enhancedContainerIsolation': {'value': true, 'locked': false}}"`

- `--proxy-http-mode=<mode>`: Sets the HTTP Proxy mode, `system` (default) or `manual`
- `--override-proxy-http=<URL>`: Sets the URL of the HTTP proxy that must be used for outgoing HTTP requests, requires `--proxy-http-mode` to be `manual`
- `--override-proxy-https=<URL>`: Sets the URL of the HTTP proxy that must be used for outgoing HTTPS requests, requires `--proxy-http-mode` to be `manual`
- `--override-proxy-exclude=<hosts/domains>`: Bypasses proxy settings for the hosts and domains. Uses a comma-separated list.
- `--proxy-enable-kerberosntlm`: Enables Kerberos and NTLM proxy authentication. If you are enabling this, ensure your proxy server is properly configured for Kerberos/NTLM authentication. Available with Docker Desktop 4.32 and later.
- `--hyper-v-default-data-root=<path>`: Specifies the default location for the Hyper-V VM disk. 
- `--windows-containers-default-data-root=<path>`: Specifies the default location for the Windows containers.
- `--wsl-default-data-root=<path>`: Specifies the default location for the WSL distribution disk.
- `--always-run-service`: After installation completes, starts `com.docker.service` and sets the service startup type to Automatic. This circumvents the need for administrator privileges, which are otherwise necessary to start `com.docker.service`. `com.docker.service` is required by Windows containers and Hyper-V backend.

> [!NOTE]
>
> If you're using PowerShell, you need to use the `ArgumentList` parameter before any flags. 
> For example:
> ```powershell
> Start-Process 'Docker Desktop Installer.exe' -Wait -ArgumentList 'install', '--accept-license'
> ```

If your admin account is different to your user account, you must add the user to the **docker-users** group:

```console
$ net localgroup docker-users <user> /add
```

## Start Docker Desktop

Docker Desktop does not start automatically after installation. To start Docker Desktop:

1. Search for Docker, and select **Docker Desktop** in the search results.

2. The Docker menu ({{< inline-image src="images/whale-x.svg" alt="whale menu" >}}) displays the Docker Subscription Service Agreement.

   {{% include "desktop-license-update.md" %}}

3. Select **Accept** to continue. Docker Desktop starts after you accept the terms.

   Note that Docker Desktop won't run if you do not agree to the terms. You can choose to accept the terms at a later date by opening Docker Desktop.

   For more information, see [Docker Desktop Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement/). It is recommended that you read the [FAQs](https://www.docker.com/pricing/faq).

> [!TIP]
>
> As an IT administrator, you can use endpoint management (MDM) software to identify the number of Docker Desktop instances and their versions within your environment. This can provide accurate license reporting, help ensure your machines use the latest version of Docker Desktop, and enable you to [enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md).
> - [Intune](https://learn.microsoft.com/en-us/mem/intune/apps/app-discovered-apps)
> - [Jamf](https://docs.jamf.com/10.25.0/jamf-pro/administrator-guide/Application_Usage.html)
> - [Kandji](https://support.kandji.io/support/solutions/articles/72000559793-view-a-device-application-list)
> - [Kolide](https://www.kolide.com/features/device-inventory/properties/mac-apps)
> - [Workspace One](https://blogs.vmware.com/euc/2022/11/how-to-use-workspace-one-intelligence-to-manage-app-licenses-and-reduce-costs.html)

## Where to go next

- Explore [Docker's subscriptions](https://www.docker.com/pricing/) to see what Docker can offer you.
- [Get started with Docker](/get-started/introduction/_index.md).
- [Explore Docker Desktop](/manuals/desktop/use-desktop/_index.md) and all its features.
- [Troubleshooting](/manuals/desktop/troubleshoot-and-support/troubleshoot/_index.md) describes common problems, workarounds, and
  how to get support.
- [FAQs](/manuals/desktop/troubleshoot-and-support/faqs/general.md) provide answers to frequently asked questions.
- [Release notes](/manuals/desktop/release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Back up and restore data](/manuals/desktop/settings-and-maintenance/backup-and-restore.md) provides instructions on backing up and restoring data related to Docker.
