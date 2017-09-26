---
description: How to install Docker EE for Windows
keywords: Windows, Windows Server, install, download, ucp, Docker EE
title: Install Docker Enterprise Edition for Windows Server 2016
redirect_from:
- /docker-ee-for-windows/install/
---

Docker Enterprise Edition for Windows Server 2016 (*Docker EE*) enables native
Docker containers on Windows Server 2016. The Docker EE installation package
includes everything you need to run Docker on Windows Server 2016.
This topic describes pre-install considerations, and how to download and
install Docker EE.

>**Looking for Release Notes?** [Get release notes for all
versions here](https://docs.docker.com/release-notes/) or subscribe to the
[releases feed on the Docker Blog](http://blog.docker.com/category/engineering/docker-releases/).

## Docker Universal Control Plane and Windows

With Docker EE, your Windows nodes can join swarms that are managed
by Docker Universal Control Plane (UCP). When you have Docker EE installed
on Windows Server 2016 and you have a
[UCP manager node provisioned](/datacenter/ucp/2.2/guides/admin/install/), you can 
[join your Windows worker nodes to a swarm](/datacenter/ucp/2.2/guides/admin/configure/join-windows-worker-nodes/). 

## Install Docker EE

Docker EE for Windows requires Windows Server 2016. See
[What to know before you install](#what-to-know-before-you-install) for a
full list of prerequisites.

1.  Open a PowerShell command prompt, and type the following commands.

    ```ps
    Install-Module DockerProvider -Force
    Install-Package Docker -ProviderName DockerProvider -Force
    ```
     
2.  Test your Docker EE installation by running the `hello-world` container.

    ```ps
    PS> docker container run hello-world:nanoserver

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
    <snip>
    ```

### (optional) Make sure you have all required updates

Some advanced Docker features (like Swarm) require that Windows is updated to include the fixes in [KB4015217](https://support.microsoft.com/en-us/help/4015217/windows-10-update-kb4015217) (or a later cumulative patch).

    ```ps
    PS> sconfig
    ```
    
Select option `6) Download and Install Updates`.

## Use a script to install Docker EE

Use the following steps when you want to install manually, script automated
installs, or install on air-gapped systems.

1.  In a PowerShell command prompt, download the installer archive on a machine
    that has a connection.

    ```ps
    # On an online machine, download the zip file.
    PS> invoke-webrequest -UseBasicparsing -Outfile docker.zip
    {{ page.win_server_zip_url }}
    ```

2.  Copy the zip file to the machine where you want to install Docker. In a
    PowerShell command prompt, use the following commands to extract the archive,
    register, and start the Docker service.

    ```ps
    # Extract the archive.
    PS> Expand-Archive docker.zip -DestinationPath $Env:ProgramFiles

    # Clean up the zip file.
    PS> Remove-Item -Force docker.zip

    # Install Docker. This will require rebooting.
    $null = Install-WindowsFeature containers

    # Add Docker to the path for the current session.
    PS> $env:path += ";$env:ProgramFiles\docker"

    # Optionally, modify PATH to persist across sessions.
    PS> $newPath = "$env:ProgramFiles\docker;" +
    [Environment]::GetEnvironmentVariable("PATH",
    [EnvironmentVariableTarget]::Machine)

    PS> [Environment]::SetEnvironmentVariable("PATH", $newPath,
    [EnvironmentVariableTarget]::Machine)

    # Register the Docker daemon as a service.
    PS> dockerd --register-service

    # Start the Docker service.
    PS> Start-Service docker
    ```

3.  Test your Docker EE installation by running the `hello-world` container.

    ```ps
    PS> docker container run hello-world:nanoserver
    ```

## Install a specific version

To install a specific Docker version, you can use the `MaximumVersion` and `MinimumVersion` flags. For example:

```ps
Install-Package -Name docker -ProviderName DockerMsftProvider -Source Docker -Force -MaximumVersion 17.03
...
Name                           Version          Source           Summary
----                           -------          ------           -------
Docker                         17.03.0-ee       Docker           Contains Docker EE for use with Windows Server 2016...
```

## Update Docker EE

To update Docker EE on Windows Server 2016:

```ps
PS> Unregister-PackageSource -ProviderName DockerMsftProvider -Name DockerDefault -Erroraction Ignore
PS> Register-PackageSource -ProviderName DockerMsftProvider -Name Docker -Erroraction Ignore -Location https://download.docker.com/components/engine/windows-server/index.json
PS> Install-Package -Name docker -ProviderName DockerMsftProvider -Update -Force

# Start the Docker service.
PS> Start-Service Docker
```

## What to know before you install

* **What the Docker EE for Windows install includes**: The installation
provides [Docker Engine](/engine/userguide/intro.md) and the
[Docker CLI client](https://docs.docker.com/engine/reference/commandline/cli/).

## About Docker EE containers and Windows Server 2016

Looking for information on using Docker EE containers?

* [Getting Started with Windows Containers (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
provides a tutorial on how to set up and run Windows containers on Windows 10
or Windows Server 2016. It shows you how to use a MusicStore application with
Windows containers.

* [Setup - Windows Server 2016 (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/Setup-Server2016.md)
describes environment setup in detail.

* Docker Container Platform for Windows Server 2016 [articles and blog
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
