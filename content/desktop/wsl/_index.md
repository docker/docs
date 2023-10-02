---
description: Turn on the Docker WSL 2 backend and get to work using best practices,
  GPU support, and more in this thorough guide.
keywords: wsl, wsl2, installing wsl2, wsl installation, docker wsl2, wsl docker, wsl2
  tech preview, wsl install docker, install docker wsl, how to install docker in wsl
title: Docker Desktop WSL 2 backend on Windows
aliases:
- /docker-for-windows/wsl/
- /docker-for-windows/wsl-tech-preview/
- /desktop/windows/wsl/
---

Windows Subsystem for Linux (WSL) 2 is a full Linux kernel built by Microsoft, which lets Linux distributions run without managing virtual machines. With Docker Desktop running on WSL 2, users can leverage Linux workspaces and avoid maintaining both Linux and Windows build scripts. In addition, WSL 2 provides improvements to file system sharing and boot time.

Docker Desktop uses the dynamic memory allocation feature in WSL 2 to improve the resource consumption. This means Docker Desktop only uses the required amount of CPU and memory resources it needs, while allowing CPU and memory-intensive tasks such as building a container, to run much faster.

Additionally, with WSL 2, the time required to start a Docker daemon after a cold start is significantly faster.

## Prerequisites

Before you turn on the Docker Desktop WSL 2 feature, ensure you have:

- At a minimum WSL version 1.1.3.0., but ideally the latest version of WSL to [avoid Docker Desktop not working as expected](best-practices.md).
- Windows 10, version 21H2 or later, or Windows 11, version 21H2 or later. For more information, see [System requirements](https://docs.docker.com/desktop/install/windows-install/#system-requirements).
- Installed the WSL 2 feature on Windows. For detailed instructions, refer to the [Microsoft documentation](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

>**Tip**
>
> For a better experience on WSL, consider enabling the WSL
> [autoMemoryReclaim](https://learn.microsoft.com/en-us/windows/wsl/wsl-config)
> setting available since WSL 1.3.10 (experimental).
>
> This feature causes the Windows host to better reclaim unused memory inside
> the WSL virtual machine, thereby resulting in better memory availability to
> other host applications. This is particularly helpful with Docker Desktop,
> since otherwise the WSL VM may consume large amounts (GBs) of memory in the
> Linux kernel's page cache as Docker builds container images, without ever
> returning that memory to the host when it becomes unused inside the VM.
{ .tip }

## Turn on Docker Desktop WSL 2

> **Important**
>
> To avoid any potential conflicts with using WSL 2 on Docker Desktop, you must uninstall any previous versions of Docker Engine and CLI installed directly through Linux distributions before installing Docker Desktop.
{ .important }

1. Download and install the latest version of [Docker Desktop for Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe).
2. Follow the usual installation instructions to install Docker Desktop. Depending on which version of Windows you are using, Docker Desktop may prompt you to turn on WSL 2 during installation. Read the information displayed on the screen and turn on the WSL 2 feature to continue.
3. Start Docker Desktop from the **Windows Start** menu.
4. Navigate to **Settings**.
5. From the **General** tab, select **Use WSL 2 based engine**..

    If you have installed Docker Desktop on a system that supports WSL 2, this option is turned on by default.
6. Select **Apply & Restart**.

Now `docker` commands work from Windows using the new WSL 2 engine.

## Enabling Docker support in WSL 2 distros

WSL 2 adds support for "Linux distros" to Windows, where each distro behaves like a VM except they all run on top of a single shared Linux kernel.

Docker Desktop does not require any particular Linux distros to be installed. The `docker` CLI and UI all work fine from Windows without any additional Linux distros. However for the best developer experience, we recommend installing at least one additional distro and enable Docker support:

1. Ensure the distribution runs in WSL 2 mode. WSL can run distributions in both v1 or v2 mode.

    To check the WSL mode, run:

     ```console
     $ wsl.exe -l -v
     ```

    To upgrade the Linux distro to v2, run:

    ```console
    $ wsl.exe --set-version (distro name) 2
    ```

    To set v2 as the default version for future installations, run:

    ```console
    $ wsl.exe --set-default-version 2
    ```

2. When Docker Desktop starts, go to **Settings** > **Resources** > **WSL Integration**.

    The Docker-WSL integration is enabled on the default WSL distribution, which is [Ubuntu](https://learn.microsoft.com/en-us/windows/wsl/install). To change your default WSL distro, run:
     ```console
    $ wsl --set-default <distro name>
    ```

3. Select **Apply & Restart**.

> **Note**
>
> Docker Desktop installs two special-purpose internal Linux distros `docker-desktop` and `docker-desktop-data`. The first (`docker-desktop`) is used to run the Docker engine (`dockerd`) while the second (`docker-desktop-data`) stores containers and images. Neither can be used for general development.

## Additional resources

- [Explore best practices](best-practices.md)
- [Understand how to develop with Docker and WSL 2 and GPU support for WSL](use-wsl.md)
