---
title: Install Docker Buildx
description: How to install Docker Buildx
keywords: build, buildx, buildkit
aliases:
  - /build/buildx/install/
---

## Docker Desktop

Docker Buildx is included by default in Docker Desktop.

## Docker Engine via package manager

Docker Linux packages also include Docker Buildx when installed using the
`.deb` or `.rpm` packages.

## Install using a Dockerfile

Here is how to install and use Buildx inside a Dockerfile through the
[`docker/buildx-bin`](https://hub.docker.com/r/docker/buildx-bin){:target="blank" rel="noopener" class=""}
image:

```dockerfile
# syntax=docker/dockerfile:1
FROM docker
COPY --from=docker/buildx-bin:latest /buildx /usr/libexec/docker/cli-plugins/docker-buildx
RUN docker buildx version
```

## Download manually

> **Important**
>
> This section is for unattended installation of the Buildx component. These
> instructions are mostly suitable for testing purposes. We do not recommend
> installing Buildx using manual download in production environments as they
> will not be updated automatically with security updates.
>
> On Windows, macOS, and Linux workstations we recommend that you install
> Docker Desktop instead. For Linux servers, we recommend that you follow the
> instructions specific for your distribution.
{: .important}

You can also download the latest binary from the [releases page on GitHub](https://github.com/docker/buildx/releases/latest){:target="blank" rel="noopener" class="_"}.

Rename the relevant binary and copy it to the destination matching your OS:

| OS       | Binary name          | Destination folder                         |
|----------|----------------------|--------------------------------------------|
| Linux    | `docker-buildx`      | `$HOME/.docker/cli-plugins`                |
| macOS    | `docker-buildx`      | `$HOME/.docker/cli-plugins`                |
| Windows  | `docker-buildx.exe`  | `%USERPROFILE%\.docker\cli-plugins`        |

Or copy it into one of these folders for installing it system-wide.

On Unix environments:

* `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
* `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

On Windows:

* `C:\ProgramData\Docker\cli-plugins`
* `C:\Program Files\Docker\cli-plugins`

> **Note**
>
> On Unix environments, it may also be necessary to make it executable with `chmod +x`:
> ```shell
> $ chmod +x ~/.docker/cli-plugins/docker-buildx
> ```

## Set Buildx as the default builder

Running the command [`docker buildx install`](../engine/reference/commandline/buildx_install.md)
sets up the `docker build` command as an alias to `docker buildx`. This results in
the ability to have [`docker build`](../engine/reference/commandline/build.md)
use the current Buildx builder.

To remove this alias, run [`docker buildx uninstall`](../engine/reference/commandline/buildx_uninstall.md).
