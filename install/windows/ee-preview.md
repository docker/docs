---
description: How to install Docker on Windows Server 1709 and Windows 10 Fall Creators Update
keywords: Windows, Windows Server, install, download, Docker EE, preview, Windows 1709
title: Preview Docker for Windows Server 1709 and Windows 10 Fall Creators Update
redirect_from:
- /engine/installation/windows/ee-preview/
---

Windows versions released in the fall of 2017 (Windows Server 1709 and Windows
10 Fall Creators Update for server and client respectively) contain improvements
that improve the experience of running Docker on Windows. Enhancements include:

 * Docker ingress mode service publishing and virtual-IP based service discovery
 * Mounting container-host named pipes into Windows containers
 * Smaller Windows container base images
 * Linux containers on Windows

Check out the
[blog post announcement for details](https://blog.docker.com/2017/09/docker-windows-server-1709/)
on use-cases for these new features.

## Windows 10 Fall Creators Update

Windows 10 users can take advantage of the new Windows features by
[installing Docker CE _Edge_ for Windows 10](/docker-for-windows/install/). Edge
comes with the latest Docker platform builds that support the new Windows features.

## Windows Server 1709

The latest version of Docker Enterprise Edition Basic is 17.06 and Docker EE
17.06 is the recommended way to run Docker containers on Windows Server 1709 in
production. Unfortunately, Docker Universal Control plane does not support
managing Windows Server 1709 workers yet. Only Docker EE Basic is supported on
Windows Server 1709.

Docker EE 17.06 does not support improved networking (ingress and virtual-IP
based service discovery), named pipe mounting nor Linux containers on Windows
Server 1709 (it does benefit from the smaller base images). For users that want
to test these new Windows features with Docker, Docker Inc. makes available
preview builds of Docker EE. These can be installed from the `preview` channel
using the standard PowerShell installation method:

    Install-Module DockerProvider
    Install-Package Docker -ProviderName DockerProvider -RequiredVersion preview

Again, Docker Universal Control Plane does not work on Windows Server 1709
whether running Docker EE 17.06 or Docker EE preview.

## FAQ

### Image compatibility

Windows container base images are not compatible between Windows Server 2016 and
Windows Server 1709 (except that Windows Server 1709 can run Windows Server 2016
based images when using Hyper-V isolation - but not the other way around). Check
the following Microsoft compatibility docs for details.

- [Version compatibility](https://docs.microsoft.com/en-us/virtualization/windowscontainers/deploy-containers/version-compatibility)
- [System requirements](https://docs.microsoft.com/en-us/virtualization/windowscontainers/deploy-containers/system-requirements)
