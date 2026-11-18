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
- /desktop/windows/wsl/
- /desktop/wsl/
- /docker-for-windows/wsl-tech-preview/
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

1. [Install Docker Desktop for Windows](/manuals/desktop/setup/install/windows-install.md).
2. During installation, Docker Desktop may prompt you to turn on WSL 2. Read the information displayed on the screen and turn on the WSL 2 feature to continue.
3. Start Docker Desktop from the **Windows Start** menu.
4. Navigate to **Settings**.
5. From the **General** tab, select **Use WSL 2 based engine**.

    If you have installed Docker Desktop on a system that supports WSL 2, this option is turned on by default and the setting is not visible. 
6. Select **Apply**.

`docker` commands are now available from any Windows terminal using the WSL 2 engine.

> [!TIP]
>
> By default, Docker Desktop stores the data for the WSL 2 engine at `C:\Users\[USERNAME]\AppData\Local\Docker\wsl`.
> If you want to change the location, go to `Settings -> Resources -> Advanced` page from the Docker Dashboard.
> Read more about this and other Windows settings at [Changing settings](/manuals/desktop/settings-and-maintenance/settings.md)

## Enable Docker in a WSL 2 distribution

WSL 2 lets multiple Linux distributions run side-by-side on a single shared kernel. Docker Desktop doesn't require a particular distribution to be installed, and `docker` commands work from Windows without one. However, enabling WSL integration for a distribution gives you direct access to `docker` commands from that distribution's terminal — which is useful for Linux-native development workflows.

> [!NOTE]
>
> Full WSL integration is available only when Docker Desktop uses the WSL 2 engine. With Docker VMM, only the Docker CLI is available in your WSL 2 distributions. For more information, see [Use Docker from WSL distributions with Docker VMM](#use-docker-from-wsl-distributions-with-docker-vmm).

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

## Use Docker from WSL distributions with Docker VMM

When Docker Desktop uses [Docker VMM](/manuals/desktop/features/vmm.md), you can use the Docker CLI from your WSL 2 distributions, but not full WSL integration. `docker` and `docker compose` work without any setup, but there's no Docker socket in the distribution, and you can only bind-mount files that are stored on Windows drives.

In a WSL 2 distribution, the `docker` and `docker compose` commands run the Windows `docker.exe` that comes with Docker Desktop, through WSL interop. Docker Desktop doesn't install anything into your distribution.

### Requirements

- Docker Desktop uses Docker VMM.
- The distribution runs in WSL 2 mode. WSL 1 distributions aren't supported.
- WSL interop is turned on, and the distribution adds the Windows `PATH` to its own `PATH`. WSL does both by default.

If WSL was already running when you installed Docker Desktop, run `wsl --shutdown` so that your distributions pick up the updated Windows `PATH`.

### Bind mounts

You can bind-mount files that are stored on a Windows drive, using the path from inside your distribution. For example, `/mnt/c/projects/app`. Docker Desktop translates the path to the Windows path, `C:\projects\app`. This also works if you changed the mount location of Windows drives with the `root` setting in the `[automount]` section of `/etc/wsl.conf`. The first time you bind-mount a directory, Docker Desktop asks whether you want to share it.

You can't bind-mount files that are stored inside your distribution, such as paths under `/home` or `/tmp`. If you try, the command fails with an error like the following:

```text
Mount denied: /home/user/app is inside the WSL distro "Ubuntu". With Docker VMM, only files on Windows drives (for example under /mnt/c) can be bind-mounted from WSL.
```

To use bind mounts, keep your projects on a Windows drive. You can still build images from a directory inside your distribution, but `docker.exe` reads the build context through WSL file sharing. This is slower, and file permissions might differ from what a Linux `docker` CLI sees.

### Limitations

- There's no Docker socket inside your distribution. Tools that connect to the Docker Engine API directly can't connect, such as Docker SDKs, Testcontainers, and VS Code Dev Containers. A Linux `docker` CLI that you installed in the distribution can't connect either.
- Stopping `docker run` with `kill` or `timeout` stops the CLI, but the container keeps running. Use `docker stop` instead. Pressing `Ctrl+C` works as usual.
- Docker Desktop passes these environment variables from your distribution to `docker.exe`: `DOCKER_CONTEXT`, `DOCKER_BUILDKIT`, `DOCKER_DEFAULT_PLATFORM`, `BUILDKIT_PROGRESS`, `COMPOSE_PROJECT_NAME`, and `COMPOSE_PROFILES`. It doesn't pass `DOCKER_HOST`, `DOCKER_CONFIG`, or the TLS variables, so `docker.exe` uses the Docker CLI configuration and contexts from Windows.

## WSL 2 security in Docker Desktop

Docker Desktop's WSL 2 integration works within WSL's existing security model and does not introduce security risks beyond standard WSL behavior.

Docker Desktop runs inside its own `docker-desktop` WSL distribution, isolated from other distributions in the same way any two WSL distributions are isolated from each other. Interaction between Docker Desktop and other distributions only occurs when you explicitly enable WSL integration for those distributions. This feature allows easy access to the Docker CLI from integrated distributions. 

WSL is designed to aid interoperability between Windows and Linux environments. Its file system is accessible from the Windows host `\\wsl$`, meaning Windows processes can read and modify files within WSL. This behavior is not specific to Docker Desktop, but rather a core aspect of WSL itself.

For environments that require stricter isolation:

- Run Docker Desktop with Hyper-V or [Docker VMM](/manuals/desktop/features/vmm.md) instead of WSL 2 to avoid the shared-kernel model entirely.
- Enable [Enhanced Container Isolation](/manuals/desktop/enterprise/hardened-desktop/enhanced-container-isolation/_index.md) to add an additional layer of protection around container workloads regardless of backend.

## Additional resources

- [Explore best practices](best-practices.md)
- [Understand how to develop with Docker and WSL 2](use-wsl.md)
- [Learn about GPU support with WSL 2](/manuals/desktop/features/gpu.md)
- [Custom kernels on WSL](custom-kernels.md)
