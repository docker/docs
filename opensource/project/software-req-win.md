---
description: How to set up a server to test Docker Engine on Windows
keywords: development, inception, container, image Dockerfile, dependencies, Go, artifacts, windows
title: Build and test Docker on Windows
---

This page explains how to get the software you need to build, test, and run the Docker source code for Windows and setup the required software and services:

- Windows containers
- GitHub account
- Git

## 1. Docker Windows containers

To test and run the Windows Docker daemon, you need a system that supports Windows Containers:

 * Windows 10 Anniversary Edition
 * Windows Server 2016 running in a VM, on bare metal or in the cloud

Check out the [getting started documentation](https://github.com/docker/labs/blob/master/windows/windows-containers/Setup.md) for details.

## 2. GitHub account

To contribute to the Docker project, you need a <a href="https://github.com" target="_blank">GitHub account</a>. A free account is fine. All the Docker project repositories are public and visible to everyone.

This guide assumes that you have basic familiarity with Git and Github terminology and usage. Refer to [GitHub For Beginners: Don’t Get Scared, Get Started](http://readwrite.com/2013/09/30/understanding-github-a-journey-for-beginners-part-1/) to get up to speed on Github.

## 3. Git

In PowerShell, run:

    Invoke-Webrequest "https://github.com/git-for-windows/git/releases/download/v2.7.2.windows.1/Git-2.7.2-64-bit.exe" -OutFile git.exe -UseBasicParsing
    Start-Process git.exe -ArgumentList '/VERYSILENT /SUPPRESSMSGBOXES /CLOSEAPPLICATIONS /DIR=c:\git\' -Wait
    setx /M PATH "$env:Path;c:\git\cmd"

You are now ready clone and build the Docker source code.

## 4. Clone Docker

In a new (to pick up the path change) PowerShell prompt, run:

    git clone https://github.com/moby/moby
    cd moby

This clones the main Docker repository. Check out [Docker on GitHub](https://github.com/moby/moby) to learn about the other software that powers the Docker platform.

## 5. Build and run

Create a builder-container with the Docker source code. You can change the source code on your system and rebuild any time:

    docker  build -t nativebuildimage -f .\Dockerfile.windows .

To build Docker, run:

    docker run --name out nativebuildimage sh -c 'cd /c/go/src/github.com/moby/moby; hack/make.sh binary'

Copy out the resulting Windows Docker daemon binary to dockerd.exe in the current directory:

    docker cp out:C:\go\src\github.com\docker\docker\bundles\$(cat VERSION)\binary-daemon\dockerd-$(cat VERSION).exe dockerd.exe

To test it, stop the system Docker daemon and start the one you just built:

    Stop-Service Docker
    .\dockerd-1.13.0-dev.exe -D

The other make targets work too, to run unit tests try: `docker run --rm docker-builder sh -c 'cd /c/go/src/github.com/moby/moby; hack/make.sh test-unit'`.


## Where to go next

In the next section, you'll [learn how to set up and configure Git for
contributing to Docker](set-up-git.md).
