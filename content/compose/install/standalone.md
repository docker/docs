---
title: Install Compose standalone
description: How to install Docker Compose - Other Scenarios
keywords: compose, orchestration, install, installation, docker, documentation
toc_max: 3
---

On this page you can find instructions on how to install Compose standalone on Linux or Windows Server, from the command line.

### On Linux

> **Compose standalone**
>
> Note that Compose standalone uses the `-compose` syntax instead of the current standard syntax `compose`.  
>For example type `docker-compose up` when using Compose standalone, instead of `docker compose up`.

1. To download and install Compose standalone, run:
   ```console
   $ curl -SL https://github.com/docker/compose/releases/download/{{% param "compose_version" %}}/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
   ```
2. Apply executable permissions to the standalone binary in the target path for the installation.
3. Test and execute compose commands using `docker-compose`.

   > [!TIP]
   >
   > If the command `docker-compose` fails after installation, check your path.
   > You can also create a symbolic link to `/usr/bin` or any other directory in your path.
   > For example:
   > ```console
   > $ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
   > ```

### On Windows Server

Follow these instructions if you are running the Docker daemon and client directly
on Microsoft Windows Server and want to install Docker Compose.

1.  Run PowerShell as an administrator.
    When asked if you want to allow this app to make changes to your device, select **Yes** in order to proceed with the installation.

2.  GitHub now requires TLS1.2. In PowerShell, run the following:

    ```powershell
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    ```
3. Run the following command to download the latest release of Compose ({{% param "compose_version" %}}):

    ```powershell
     Start-BitsTransfer -Source "https://github.com/docker/compose/releases/download/{{% param "compose_version" %}}/docker-compose-windows-x86_64.exe" -Destination $Env:ProgramFiles\Docker\docker-compose.exe
    ```

    > [!NOTE]
    >
    > On Windows Server 2019 you can add the Compose executable to `$Env:ProgramFiles\Docker`.
     Because this directory is registered in the system `PATH`, you can run the `docker-compose --version` 
     command on the subsequent step with no additional configuration.

    > To install a different version of Compose, substitute `{{% param "compose_version" %}}`
    > with the version of Compose you want to use.

4.  Test the installation.

    ```console
    $ docker-compose.exe version
    Docker Compose version {{% param "compose_version" %}}
    ```
