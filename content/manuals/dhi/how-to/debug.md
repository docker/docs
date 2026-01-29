---
title: Debug a Docker Hardened Image container
linkTitle: Debug a container
weight: 60
keywords: debug, hardened images, DHI, troubleshooting, ephemeral container, docker debug, non-root containers, hardened container image, debug secure container
description: Learn how to use Docker Debug to troubleshoot Docker Hardened Images (DHI) locally or in production.
---

Docker Hardened Images (DHI) prioritize minimalism and security, which means
they intentionally leave out many common debugging tools (like shells or package
managers). This makes direct troubleshooting difficult without introducing risk.
To address this, you can use [Docker
Debug](../../../reference/cli/docker/debug.md), a secure workflow that
temporarily attaches an ephemeral debug container to a running service or image
without modifying the original image.

This guide shows how to debug Docker Hardened Images locally during development.
With Docker Debug, you can also debug containers remotely using the `--host`
option.

## Use Docker Debug

### Step 1: Run a container from a Hardened Image

Start with a DHI-based container that simulates an issue:

```console
$ docker run -d --name myapp dhi.io/python:3.13 python -c "import time; time.sleep(300)"
```

This container doesn't include a shell or tools like `ps`, `top`, or `cat`.

If you try:

```console
$ docker exec -it myapp sh
```

You'll see:

```console
exec: "sh": executable file not found in $PATH
```

### Step 2: Use Docker Debug to inspect the container

Use the `docker debug` command to attach a temporary, tool-rich debug container to the running instance.

```console
$ docker debug myapp
```

From here, you can inspect running processes, network status, or mounted files.

For example, to check running processes:

```console
$ ps aux
```

Type `exit` to leave the container when done.

## Alternative debugging approaches

In addition to using Docker Debug, you can also use the following approaches for
debugging DHI containers.

### Use the -dev variant

Docker Hardened Images offer a `-dev` variant that includes a such as a shell
and a package manager to install debugging tools. Simply replace the image tag
with `-dev`:

```console
$ docker run -it --rm dhi.io/python:3.13-dev sh
```

Type `exit` to leave the container when done. Note that using the `-dev` variant
increases the attack surface and it is not recommended as a runtime for
production environments.

### Mount debugging tools with image mounts

You can use the image mount feature to mount debugging tools into your container
without modifying the base image.

#### Step 1: Run a container from a Hardened Image

Start with a DHI-based container that simulates an issue:

```console
$ docker run -d --name myapp dhi.io/python:3.13 python -c "import time; time.sleep(300)"
```

#### Step 2: Mount debugging tools into the container

Run a new container that mounts a tool-rich image (like `busybox`) into
the running container's namespace:

```console
$ docker run --rm -it --pid container:myapp \
  --mount type=image,source=busybox,destination=/dbg,ro \
  dhi.io/python:3.13 /dbg/bin/sh
```

This mounts the BusyBox image at `/dbg`, giving you access to its tools while
keeping your original container image unchanged. Since the hardened Python image
doesn't include standard utilities, you need to use the full path to the mounted
tools:

```console
$ /dbg/bin/ls /
$ /dbg/bin/ps aux
$ /dbg/bin/cat /etc/os-release
```

Type `exit` to leave the container when done.

## What's next

This guide covered three approaches for debugging Docker Hardened Images:

- Docker Debug: Attach an ephemeral debug container without modifying the original image
- `-dev` variants: Use development images that include debugging tools
- Image mount: Mount tool-rich images like BusyBox to access debugging utilities

Each method helps you troubleshoot hardened containers while maintaining
security. Docker Debug and image mounts avoid modifying your production images,
while `-dev` variants provide convenience during development.

If you encounter issues related to permissions, ports, missing shells, or
package managers, see [Troubleshoot Docker Hardened Images](../troubleshoot.md)
for recommended solutions and workarounds.