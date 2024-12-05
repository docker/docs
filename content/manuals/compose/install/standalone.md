---
title: Install the Docker Compose standalone
linkTitle: Standalone
description: How to install Docker Compose - Other Scenarios
keywords: compose, orchestration, install, installation, docker, documentation
toc_max: 3
weight: 20
---

This page contains instructions on how to install Docker Compose standalone on Linux or Windows Server, from the command line.

> [!WARNING]
>
> The Docker Compose standalone uses the `-compose` syntax instead of the current standard syntax `compose`.  
> For example, you must type `docker-compose up` when using Docker Compose standalone, instead of `docker compose up`.

## On Linux

1. To download and install the Docker Compose standalone, run:

   ```console
   $ curl -SL https://github.com/docker/compose/releases/download/{{% param "compose_version" %}}/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
   ```

2. Apply executable permissions to the standalone binary in the target path for the installation.

   ```console
   $ chmod +x /usr/local/bin/docker-compose
   ```

3. Test and execute Docker Compose commands using `docker-compose`.

> [!TIP]
>
> If the command `docker-compose` fails after installation, check your path.
> You can also create a symbolic link to `/usr/bin` or any other directory in your path.
> For example:
> ```console
> $ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
> ```

## On Windows Server

Follow these instructions if you are [running the Docker daemon directly
on Microsoft Windows Server](/manuals/engine/install/binaries.md#install-server-and-client-binaries-on-windows) and want to install Docker Compose.

1.  Run PowerShell as an administrator.
    In order to proceed with the installation, select **Yes** when asked if you want this app to make changes to your device.

2.  Optional. Ensure TLS1.2 is enabled. 
    GitHub requires TLS1.2 fore secure connections. If you’re using an older version of Windows Server, for example 2016, or suspect that TLS1.2 is not enabled, run the following command in PowerShell:

    ```powershell
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    ```

3. Download the latest release of Docker Compose ({{% param "compose_version" %}}). Run the following command:

    ```powershell
     Start-BitsTransfer -Source "https://github.com/docker/compose/releases/download/{{% param "compose_version" %}}/docker-compose-windows-x86_64.exe" -Destination $Env:ProgramFiles\Docker\docker-compose.exe
    ```

    To install a different version of Docker Compose, substitute `{{% param "compose_version" %}}` with the version of Compose you want to use.

    > [!NOTE]
    >
    > On Windows Server 2019 you can add the Compose executable to `$Env:ProgramFiles\Docker`.
     Because this directory is registered in the system `PATH`, you can run the `docker-compose --version` 
     command on the subsequent step with no additional configuration.

4.  Test the installation.

    ```console
    $ docker-compose.exe version
    Docker Compose version {{% param "compose_version" %}}
    ```
