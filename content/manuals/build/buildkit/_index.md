---
title: BuildKit
weight: 100
description: Introduction and overview of BuildKit
keywords: build, buildkit
---

## Overview

[BuildKit](https://github.com/moby/buildkit)
is an improved backend to replace the legacy builder. BuildKit is the default builder
for users on Docker Desktop, and Docker Engine as of version 23.0.

BuildKit provides new functionality and improves your builds' performance.
It also introduces support for handling more complex scenarios:

- Detect and skip executing unused build stages
- Parallelize building independent build stages
- Incrementally transfer only the changed files in your
  [build context](../concepts/context.md) between builds
- Detect and skip transferring unused files in your
  [build context](../concepts/context.md)
- Use [Dockerfile frontend](frontend.md) implementations with many
  new features
- Avoid side effects with rest of the API (intermediate images and containers)
- Prioritize your build cache for automatic pruning

Apart from many new features, the main areas BuildKit improves on the current
experience are performance, storage management, and extensibility. From the
performance side, a significant update is a new fully concurrent build graph
solver. It can run build steps in parallel when possible and optimize out
commands that don't have an impact on the final result.
The access to the local source files has also been optimized. By tracking
only the updates made to these
files between repeated build invocations, there is no need to wait for local
files to be read or uploaded before the work can begin.

## LLB

At the core of BuildKit is a
[Low-Level Build (LLB)](https://github.com/moby/buildkit#exploring-llb) definition format. LLB is an intermediate binary format
that allows developers to extend BuildKit. LLB defines a content-addressable
dependency graph that can be used to put together complex build
definitions. It also supports features not exposed in Dockerfiles, like direct
data mounting and nested invocation.

{{< figure src="../images/buildkit-dag.svg" class="invertible" >}}

Everything about execution and caching of your builds is defined in LLB. The
caching model is entirely rewritten compared to the legacy builder. Rather than
using heuristics to compare images, LLB directly tracks the checksums of build
graphs and content mounted to specific operations. This makes it much faster,
more precise, and portable. The build cache can even be exported to a registry,
where it can be pulled on-demand by subsequent invocations on any host.

LLB can be generated directly using a
[golang client package](https://pkg.go.dev/github.com/moby/buildkit/client/llb) that allows defining the relationships between your
build operations using Go language primitives. This gives you full power to run
anything you can imagine, but will probably not be how most people will define
their builds. Instead, most users would use a frontend component, or LLB nested
invocation, to run a prepared set of build steps.

## Frontend

A frontend is a component that takes a human-readable build format and converts
it to LLB so BuildKit can execute it. Frontends can be distributed as images,
and the user can target a specific version of a frontend that is guaranteed to
work for the features used by their definition.

For example, to build a [Dockerfile](/reference/dockerfile.md) with
BuildKit, you would
[use an external Dockerfile frontend](frontend.md).

## Getting started

BuildKit is the default builder for users on Docker Desktop and Docker Engine
v23.0 and later.

If you have installed Docker Desktop, you don't need to enable BuildKit. If you
are running a version of Docker Engine version earlier than 23.0, you can enable
BuildKit either by setting an environment variable, or by making BuildKit the
default setting in the daemon configuration.

To set the BuildKit environment variable when running the `docker build`
command, run:

```console
$ DOCKER_BUILDKIT=1 docker build .
```

> [!NOTE]
>
> Buildx always uses BuildKit.

To use Docker BuildKit by default, edit the Docker daemon configuration in
`/etc/docker/daemon.json` as follows, and restart the daemon.

```json
{
  "features": {
    "buildkit": true
  }
}
```

If the `/etc/docker/daemon.json` file doesn't exist, create new file called
`daemon.json` and then add the following to the file. And restart the Docker
daemon.

## BuildKit on Windows

> [!WARNING]
>
> BuildKit only fully supports building Linux containers. Windows container
> support is experimental.

BuildKit has experimental support for Windows containers (WCOW) as of version 0.13.
This section walks you through the steps for trying it out.
To share feedback, [open an issue in the repository](https://github.com/moby/buildkit/issues/new), especially `buildkitd.exe`.

### Known limitations

For information about open bugs and limitations related to BuildKit on Windows,
see [GitHub issues](https://github.com/moby/buildkit/issues?q=is%3Aissue%20state%3Aopen%20label%3Aarea%2Fwindows-wcow).

### Prerequisites

- Architecture: `amd64`, `arm64` (binaries available but not officially tested yet).
- Supported OS: Windows Server 2019, Windows Server 2022, Windows 11.
- Base images: `ServerCore:ltsc2019`, `ServerCore:ltsc2022`, `NanoServer:ltsc2022`.
  See the [compatibility map here](https://learn.microsoft.com/en-us/virtualization/windowscontainers/deploy-containers/version-compatibility?tabs=windows-server-2019%2Cwindows-11#windows-server-host-os-compatibility).
- Docker Desktop version 4.29 or later

### Steps

> [!NOTE]
>
> The following commands require administrator (elevated) privileges in a PowerShell terminal.

1. Enable the **Hyper-V** and **Containers** Windows features.

   ```console
   > Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V, Containers -All
   ```

   If you see `RestartNeeded` as `True`, restart your machine and re-open a PowerShell terminal as an administrator.
   Otherwise, continue with the next step.

2. Switch to Windows containers in Docker Desktop.

   Select the Docker icon in the taskbar, and then **Switch to Windows containers...**.

3. Install containerd version 1.7.7 or later following the setup instructions [here](https://github.com/containerd/containerd/blob/main/docs/getting-started.md#installing-containerd-on-windows).

4. Download and extract the latest BuildKit release.

   ```powershell
   $version = "v0.22.0" # specify the release version, v0.13+
   $arch = "amd64" # arm64 binary available too
   curl.exe -LO https://github.com/moby/buildkit/releases/download/$version/buildkit-$version.windows-$arch.tar.gz
   # there could be another `.\bin` directory from containerd instructions
   # you can move those
   mv bin bin2
   tar.exe xvf .\buildkit-$version.windows-$arch.tar.gz
   ## x bin/
   ## x bin/buildctl.exe
   ## x bin/buildkitd.exe
   ```

5. Install BuildKit binaries on `PATH`.

   ```powershell
   # after the binaries are extracted in the bin directory
   # move them to an appropriate path in your $Env:PATH directories or:
   Copy-Item -Path ".\bin" -Destination "$Env:ProgramFiles\buildkit" -Recurse -Force
   # add `buildkitd.exe` and `buildctl.exe` binaries in the $Env:PATH
   $Path = [Environment]::GetEnvironmentVariable("PATH", "Machine") + `
       [IO.Path]::PathSeparator + "$Env:ProgramFiles\buildkit"
   [Environment]::SetEnvironmentVariable( "Path", $Path, "Machine")
   $Env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + `
       [System.Environment]::GetEnvironmentVariable("Path","User")
   ```
6. Start the BuildKit daemon.

   ```console
   > buildkitd.exe
   ```
   > [!NOTE]
   > If you are running a _dockerd-managed_ `containerd` process, use that instead, by supplying the address:
   > `buildkitd.exe --containerd-worker-addr "npipe:////./pipe/docker-containerd"`

7. In another terminal with administrator privileges, create a remote builder that uses the local BuildKit daemon.

   > [!NOTE]
   >
   > This requires Docker Desktop version 4.29 or later.

   ```console
   > docker buildx create --name buildkit-exp --use --driver=remote npipe:////./pipe/buildkitd
   buildkit-exp
   ```

8. Verify the builder connection by running `docker buildx inspect`.

   ```console
   > docker buildx inspect
   ```

   The output should indicate that the builder platform is Windows,
   and that the endpoint of the builder is a named pipe.

   ```text
   Name:          buildkit-exp
    Driver:        remote
    Last Activity: 2024-04-15 17:51:58 +0000 UTC
    Nodes:
    Name:             buildkit-exp0
    Endpoint:         npipe:////./pipe/buildkitd
    Status:           running
    BuildKit version: v0.13.1
    Platforms:        windows/amd64
   ...
   ```

9. Create a Dockerfile and build a `hello-buildkit` image.

   ```console
   > mkdir sample_dockerfile
   > cd sample_dockerfile
   > Set-Content Dockerfile @"
   FROM mcr.microsoft.com/windows/nanoserver:ltsc2022
   USER ContainerAdministrator
   COPY hello.txt C:/
   RUN echo "Goodbye!" >> hello.txt
   CMD ["cmd", "/C", "type C:\\hello.txt"]
   "@
   Set-Content hello.txt @"
   Hello from BuildKit!
   This message shows that your installation appears to be working correctly.
   "@
   ```

10. Build and push the image to a registry.

    ```console
    > docker buildx build --push -t <username>/hello-buildkit .
    ```

11. After pushing to the registry, run the image with `docker run`.

    ```console
    > docker run <username>/hello-buildkit
    ```
