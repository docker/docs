---
description: Turn on the Docker WSL 2 backend and get to work using best practices,
  GPU support, and more in this thorough guide.
keywords: wsl, wsl2, installing wsl2, wsl installation, docker wsl2, wsl docker, wsl2
  tech preview, wsl install docker, install docker wsl, how to install docker in wsl
title: Docker Desktop WSL 2 backend on Windows
linkTitle: WSL
weight: 120
aliases:
- /docker-for-windows/wsl/
- /docker-for-windows/wsl-tech-preview/
- /desktop/windows/wsl/
- /desktop/wsl/
---

Windows Subsystem for Linux (WSL) 2 is a full Linux kernel built by Microsoft that lets Linux distributions run without managing virtual machines. With Docker Desktop running on WSL 2, users can leverage Linux workspaces and avoid maintaining both Linux and Windows build scripts. In addition, WSL 2 provides improvements to file system sharing, faster cold-start times, and dynamic resource allocation.

Because WSL 2 uses dynamic memory allocation, Docker Desktop requests only the CPU and memory it actually needs — freeing resources for the rest of your system, while still letting memory-intensive tasks such as multi-stage image builds run at full speed.

## Prerequisites

Before you turn on the Docker Desktop WSL 2 feature, ensure you have:

- At a minimum WSL version 2.1.5, but ideally the latest version of WSL to [avoid Docker Desktop not working as expected](best-practices.md).
- Met the Docker Desktop for Windows' [system requirements](/manuals/desktop/setup/install/windows-install.md#system-requirements).
- Installed the WSL 2 feature on Windows. For detailed instructions, refer to the [Microsoft documentation](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

> [!TIP]
>
> Consider enabling the WSL [autoMemoryReclaim](https://learn.microsoft.com/en-us/windows/wsl/wsl-config#experimental-settings) setting, available since WSL 1.3.10 (experimental).
>This setting allows Windows to reclaim unused memory from the WSL virtual machine, preventing the Linux kernel's page cache from holding onto large amounts of RAM after container image builds complete. The result is better memory availability for other applications on the host. 

## Turn on Docker Desktop WSL 2

Before installing Docker Desktop, uninstall any version of Docker Engine or the Docker CLI
that was installed directly inside a WSL Linux distribution. Running both can cause conflicts.

1. Download and install the latest version of [Docker Desktop for Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-windows).
2. Follow the usual installation instructions to install Docker Desktop. Depending on which version of Windows you are using, Docker Desktop may prompt you to turn on WSL 2 during installation. Read the information displayed on the screen and turn on the WSL 2 feature to continue.
3. Start Docker Desktop from the **Windows Start** menu.
4. Navigate to **Settings**.
5. From the **General** tab, select **Use WSL 2 based engine**.

    If you have installed Docker Desktop on a system that supports WSL 2, this option is turned on by default.
6. Select **Apply**.

`docker` commands are now available from any Windows terminal using the WSL 2 engine.

> [!TIP]
>
> By default, Docker Desktop stores the data for the WSL 2 engine at `C:\Users\[USERNAME]\AppData\Local\Docker\wsl`.
> If you want to change the location, go to `Settings -> Resources -> Advanced` page from the Docker Dashboard.
> Read more about this and other Windows settings at [Changing settings](/manuals/desktop/settings-and-maintenance/settings.md)

## Enable Docker in a WSL 2 distribution

WSL 2 lets multiple Linux distributions run side-by-side on a single shared kernel. Docker Desktop doesn't require a particular distribution to be installed, and `docker` commands work from Windows without one. However, enabling WSL integration for a distribution gives you direct access to `docker` commands from that distribution's terminal — which is useful for Linux-native development workflows.

1. Ensure the distribution runs in WSL 2 mode. WSL can run distributions in both v1 or v2 mode.

    To check the WSL mode, run:

     ```console
     $ wsl.exe -l -v
     ```

    To upgrade the Linux distribution to v2, run:

    ```console
    $ wsl.exe --set-version (distribution name) 2
    ```

    To set v2 as the default version for future installations, run:

    ```console
    $ wsl.exe --set-default-version 2
    ```

2. When Docker Desktop starts, go to **Settings** > **Resources** > **WSL Integration**.

    The Docker-WSL integration is enabled on the default WSL distribution, which is [Ubuntu](https://learn.microsoft.com/en-us/windows/wsl/install). To change your default WSL distribution, run:
     ```console
    $ wsl.exe --set-default <distribution name>
    ```
   If **WSL integrations** isn't available under **Resources**, Docker may be in Windows container mode. In your taskbar, select the Docker menu and then **Switch to Linux containers**.

3. Select **Apply**.

## WSL 2 security in Docker Desktop

Docker Desktop's WSL 2 integration works within WSL's existing security model and does not introduce security risks beyond standard WSL behavior.

Docker Desktop runs inside its own `docker-desktop` WSL distribution, isolated from other distributions in the same way any two WSL distributions are isolated from each other. Interaction between Docker Desktop and other distributions only occurs when you explicitly enable WSL integration for those distributions. This feature allows easy access to the Docker CLI from integrated distributions. 

WSL is designed to aid interoperability between Windows and Linux environments. Its file system is accessible from the Windows host `\\wsl$`, meaning Windows processes can read and modify files within WSL. This behavior is not specific to Docker Desktop, but rather a core aspect of WSL itself.

For environments that require stricter isolation:

- Run Docker Desktop in Hyper-V mode instead of WSL 2 to avoid the shared-kernel model entirely.
- Enable [Enhanced Container Isolation](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/_index.md) to add an additional layer of protection around container workloads regardless of backend.

## Additional resources

- [Explore best practices](best-practices.md)
- [Understand how to develop with Docker and WSL 2](use-wsl.md)
- [Learn about GPU support with WSL 2](/manuals/desktop/features/gpu.md)
- [Custom kernels on WSL](custom-kernels.md)
