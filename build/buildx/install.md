---
title: Install Docker Buildx
description: How to install Docker Buildx
keywords: build, buildx, buildkit
---

## Docker Desktop

Docker Buildx is included in [Docker Desktop](../../desktop/index.md) for Windows, macOS, and Linux.

## Linux packages

Docker Linux packages also include Docker Buildx when installed using the [DEB or RPM packages](../../engine/install/index.md).

## Manual download

> **Important**
>
> This section is for unattended installation of the buildx component. These
> instructions are mostly suitable for testing purposes. We do not recommend
> installing buildx using manual download in production environments as they
> will not be updated automatically with security updates.
>
> On Windows, macOS, and Linux workstations we recommend that you install
> [Docker Desktop](../../desktop/index.md) instead. For Linux servers, we recommend
> that you follow the [instructions specific for your distribution](#linux-packages).
{: .important}

You can also download the latest binary from the [releases page on GitHub](https://github.com/docker/buildx/releases/latest){:target="_blank" rel="noopener" class="_"}.

Rename the relevant binary and copy it to the destination matching your OS:

| OS       | Binary name          | Destination folder                       |
| -------- | -------------------- | -----------------------------------------|
| Linux    | `docker-buildx`      | `$HOME/.docker/cli-plugins`              |
| macOS    | `docker-buildx`      | `$HOME/.docker/cli-plugins`              |
| Windows  | `docker-buildx.exe`  | `%USERPROFILE%\.docker\cli-plugins`      |

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

## Dockerfile

Here is how to install and use Buildx inside a Dockerfile through the
[`docker/buildx-bin`](https://hub.docker.com/r/docker/buildx-bin) image:

```dockerfile
FROM docker
COPY --from=docker/buildx-bin:latest /buildx /usr/libexec/docker/cli-plugins/docker-buildx
RUN docker buildx version
```

## Set buildx as the default builder

Running the command [`docker buildx install`](../../engine/reference/commandline/buildx_install.md)
sets up docker builder command as an alias to `docker buildx`. This results in
the ability to have [`docker build`](../../engine/reference/commandline/build.md)
use the current buildx builder.

To remove this alias, run [`docker buildx uninstall`](../../engine/reference/commandline/buildx_uninstall.md).
