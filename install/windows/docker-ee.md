---
description: How to install Docker Engine - Enterprise for Windows Server
keywords: Windows, Windows Server, install, download, ucp, Docker Engine - Enterprise
title: Install Docker Engine - Enterprise on Windows Servers
redirect_from:
- /docker-ee-for-windows/install/
- /engine/installation/windows/docker-ee/
---

{% capture filename %}{{ page.win_latest_build }}.zip{% endcapture %} {% capture download_url %}https://download.docker.com/components/engine/windows-server/{{ site.docker_ee_version }}/{{ filename }}{% endcapture %}

Docker Engine - Enterprise enables native Docker containers on Windows Server. Windows Server 2016 and later versions are supported. The Docker Engine - Enterprise installation package includes everything you need to run Docker on Windows Server.  This topic describes pre-install considerations, and how to download and install Docker Engine - Enterprise.

> Release notes
>
> [Release notes for all versions](/release-notes/)

## System requirements

Windows OS requirements around specific CPU and RAM requirements also need to be met as specified 
in the [Windows Server Requirements](https://docs.microsoft.com/en-us/windows-server/get-started/system-requirements).
This provides information for specific CPU and memory specs and capabilities (instruction sets like CMPXCHG16b, 
LAHF/SAHF, and PrefetchW, security: DEP/NX, etc.).

* OS Versions: Server 2016 (Core and GUI), 1709 and 1803
* RAM: 4GB
* Disk space: [32 GB minimum recommendation for Windows](https://docs.microsoft.com/en-us/windows-server/get-started/system
  requirements). An additional 32 GB of Space is recommended for base images for ServerCore and NanoServer along with buffer
  space for workload containers running IIS, SQL Server and .Net apps.

## Install Docker Engine - Enterprise

Docker Engine - Enterprise requires Windows Server 2016, 1703, or 1803. See
[What to know before you install](#what-to-know-before-you-install) for a
full list of prerequisites.

1.  Open a PowerShell command prompt, and type the following commands.

    ```PowerShell
    Install-Module DockerMsftProvider -Force
    Install-Package Docker -ProviderName DockerMsftProvider -Force
    ```

2.  Check if a reboot is required, and if yes, restart your instance:

    ```PowerShell
    (Install-WindowsFeature Containers).RestartNeeded
    ```
    If the output of this command is **Yes**, then restart the server with:

    ```PowerShell
    Restart-Computer
    ```

3.  Test your Docker Engine - Enterprise installation by running the 
    `hello-world` container.

    ```PowerShell
    docker run hello-world:nanoserver

    Unable to find image 'hello-world:nanoserver' locally
    nanoserver: Pulling from library/hello-world
    bce2fbc256ea: Pull complete
    3ac17e2e6106: Pull complete
    8cac44e17f16: Pull complete
    5e160e4d8db3: Pull complete
    Digest: sha256:25eac12ba40f7591969085ab3fb9772e8a4307553c14ea72d0e6f98b2c8ced9d
    Status: Downloaded newer image for hello-world:nanoserver

    Hello from Docker!
    This message shows that your installation appears to be working correctly.
    ```

### (optional) Make sure you have all required updates

Some advanced Docker features, such as swarm mode, require the fixes included in
[KB4015217](https://support.microsoft.com/en-us/help/4015217/windows-10-update-kb4015217)
(or a later cumulative patch).

```PowerShell
sconfig
```

Select option `6) Download and Install Updates`.


### FIPS 140-2 cryptographic module support

[Federal Information Processing Standards (FIPS) Publication 140-2](https://csrc.nist.gov/csrc/media/publications/fips/140/2/final/documents/fips1402.pdf) is a United States Federal security requirement for cryptographic modules.

With Docker EE Basic license for versions 18.09 and later, Docker provides FIPS 140-2 support in Windows Server 2016. This includes a FIPS supported cryptographic module. If the Windows implementation already has FIPS support enabled, FIPS is automatically enabled in the Docker engine.

**NOTE:** FIPS 140-2 is only supported in the Docker EE engine. UCP and DTR currently do not have support for FIPS 140-2.

To enable FIPS 140-2 compliance on a system that is not in FIPS 140-2 mode, do the following in PowerShell:

```
[System.Environment]::SetEnvironmentVariable("DOCKER_FIPS", "1", "Machine") 
```

Restart the Docker service by running the following command.

```
net stop docker
net start docker
```

To confirm Docker is running with FIPS-140-2 enabled, run the `docker info` command:

```
Labels:    
 com.docker.security.fips=enabled 
```

**NOTE:** If the system has the FIPS-140-2 cryptographic module installed on the operating system, it is possible to disable FIPS-140-2 compliance. To disable FIPS-140-2 in Docker but not the operating system, set the value `"DOCKER_FIPS","0"` in the `[System.Environment]`.`

## Use a script to install Docker EE

Use the following steps when you want to install manually, script automated
installs, or install on air-gapped systems.

1.  In a PowerShell command prompt, download the installer archive on a machine
    that has a connection.

    ```PowerShell
    # On an online machine, download the zip file.
    invoke-webrequest -UseBasicparsing -Outfile {{ filename }} {{ download_url }}
    ```

2.  Copy the zip file to the machine where you want to install Docker. In a
    PowerShell command prompt, use the following commands to extract the archive,
    register, and start the Docker service.

    ```PowerShell
    #Stop Docker service
    Stop-Service docker
    
    # Extract the archive.
    Expand-Archive {{ filename }} -DestinationPath $Env:ProgramFiles -Force

    # Clean up the zip file.
    Remove-Item -Force {{ filename }}

    # Install Docker. This requires rebooting.
    $null = Install-WindowsFeature containers

    # Add Docker to the path for the current session.
    $env:path += ";$env:ProgramFiles\docker"

    # Optionally, modify PATH to persist across sessions.
    $newPath = "$env:ProgramFiles\docker;" +
    [Environment]::GetEnvironmentVariable("PATH",
    [EnvironmentVariableTarget]::Machine)

    [Environment]::SetEnvironmentVariable("PATH", $newPath,
    [EnvironmentVariableTarget]::Machine)

    # Register the Docker daemon as a service.
    dockerd --register-service

    # Start the Docker service.
    Start-Service docker
    ```

3.  Test your Docker EE installation by running the `hello-world` container.

    ```PowerShell
    docker container run hello-world:nanoserver
    ```

## Install a specific version

To install a specific version, use the `RequiredVersion` flag:

```PowerShell
Install-Package -Name docker -ProviderName DockerMsftProvider -Force -RequiredVersion 18.09
...
Name                      Version               Source           Summary
----                      -------               ------           -------
Docker                    18.09.0               Docker           Contains Docker Engine - Enterprise for use with Windows Server...
```

### Updating the DockerMsftProvider
Installing specific Docker EE versions may require an update to previously installed DockerMsftProvider modules. To update:

```PowerShell
Update-Module DockerMsftProvider
```

Then open a new Powershell session for the update to take effect.

## Update Docker Engine - Enterprise

To update Docker Engine - Enterprise to the most recent release, specify the `-RequiredVersion` and `-Update` flags:

```PowerShell
Install-Package -Name docker -ProviderName DockerMsftProvider -RequiredVersion 18.09 -Update -Force
```
The required version must match any of the versions available in this json file: https://dockermsft.blob.core.windows.net/dockercontainer/DockerMsftIndex.json


## Preparing a Docker EE Engine for use with UCP

Run the
[UCP installation script for Windows](/datacenter/ucp/3.0/guides/admin/configure/join-windows-worker-nodes/#run-the-windows-node-setup-script).

Start the Docker service:

```PowerShell
Start-Service Docker
```

* **What the Docker Engine - Enterprise  install includes**: The installation
provides [Docker Engine](/engine/userguide/intro.md) and the
[Docker CLI client](/engine/reference/commandline/cli.md).

## About Docker Engine - Enterprise containers and Windows Server

Looking for information on using Docker Engine - Enterprise containers?

* [Getting Started with Windows Containers (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
provides a tutorial on how to set up and run Windows containers on Windows 10
or Windows Server 2016. It shows you how to use a MusicStore application with
Windows containers.

* [Setup - Windows Server 2016 (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/Setup-Server2016.md)
describes environment setup in detail.

* Docker Container Platform for Windows Server [articles and blog
posts](https://www.docker.com/microsoft/) on the Docker website.

## Where to go next

* [Getting started](/docker-for-windows/index.md) provides an overview of
Docker for Windows, basic Docker command examples, how to get help or give
feedback, and links to all topics in the Docker for Windows guide.

* [FAQs](/docker-for-windows/faqs.md) provides answers to frequently asked
questions.

* [Release Notes](/docker-for-windows/release-notes.md) lists component
updates, new features, and improvements associated with Stable and Edge
releases.

* [Learn Docker](/learn.md) provides general Docker tutorials.

* [Windows Containers on Windows Server](https://docs.microsoft.com/en-us/virtualization/windowscontainers/quick-start/quick-start-windows-server)
is the official Microsoft documentation.
