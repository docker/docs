---
description: How to install Docker EE for Windows
keywords: windows, windows server, install, download
title: Install Docker Enterprise Edition for Windows Server 2016
---

Docker Enterprise Edition for Windows Server 2016 (*Docker EE*) enables native
Docker containers on Windows Server 2016. The Docker EE installation package
includes everything you need to run Docker on Windows Server 2016.
This topic describes pre-install considerations, and how to download and
install Docker EE.

> **Already have Docker EE for Windows?** If you already have Docker EE for
Windows installed, and you're ready to get started, skip to
[Get started with Docker for Windows](/docker-for-windows/index.md) for a quick
tour of the command line, settings, and tools.
>
>**Looking for Release Notes?** [Get release notes for all
versions here](https://docs.docker.com/release-notes/) or subscribe to the
[releases feed on the Docker Blog](http://blog.docker.com/category/engineering/docker-releases/).

## Install Docker EE

Docker EE for Windows requires Windows Server 2016. See [What to know before
you install](/docker-ee-for-windows/#what-to-know-before-you-install) for a
full list of prerequisites. To install Docker Community Edition (*Docker CE*)
on a Windows 10 machine, see [Install Docker for Windows](/docker-for-windows/install.md).

1.  Open a PowerShell command prompt, and type the following commands.

    ```ps
    PS> Install-Module -Name DockerMsftProvider -Force
    PS> Install-Package -Name docker -ProviderName DockerMsftProvider -Force
    PS> Restart-Computer -Force
    ```

2.  Test your Docker EE installation by running the `hello-world` container.

    ```ps
    PS> docker run hello-world:nanoserver

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


## Using a script to install Docker EE

Use the following steps when you want to install manually, script automated
installs, or install on air-gapped systems.

1.  In a PowerShell command prompt, download the installer archive on a machine
that has a connection.

    ```ps
    # On an online machine, download the zip file.
    PS> invoke-webrequest -UseBasicparsing -Outfile docker.zip https://download.docker.com/components/engine/windows-server/17.03/docker-17.03.0-ee.zip
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
    PS> $env:path += "$env:ProgramFiles\docker"

    # Optionally, modify PATH to persist across sessions.
    PS> $newPath = "$env:ProgramFiles\docker;" +
    [Environment]::GetEnvironmentVariable("PATH",
    [EnvironmentVariableTarget]::Machine)

    PS> [Environment]::SetEnvironmentVariable("PATH", $newPath,
    [EnvironmentVariableTarget]::Machine)

    # Register the Docker daemon as a service.
    PS> dockerd --register-service

    # Start the daemon.
    PS> Start-Service docker
    ```

3.  Test your Docker EE installation by running the `hello-world` container.

    ```ps
    PS> docker run hello-world:nanoserver
    ```


## Install Docker EE using OneGet

If you want to install Docker EE by using [OneGet](https://github.com/oneget/oneget),
follow the steps described in [Windows Containers on Windows
Server](https://docs.microsoft.com/en-us/virtualization/windowscontainers/quick-start/quick-start-windows-server).

##  What to know before you install

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
