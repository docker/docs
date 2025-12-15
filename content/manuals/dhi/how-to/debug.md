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

This guide shows how to debug Docker Hardened Images locally during
development. You can also debug containers remotely using the `--host` option.

The following example uses a mirrored `python:3.13` image, but the same steps apply to any image.

## Step 1: Run a container from a Hardened Image

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

## Step 2: Use Docker Debug to inspect the container

Use the `docker debug` command to attach a temporary, tool-rich debug container to the running instance.

```console
$ docker debug myapp
```

From here, you can inspect running processes, network status, or mounted files.

For example, to check running processes:

```console
$ ps aux
```

Exit the debug session with:

```console
$ exit
```

## What's next

Docker Debug helps you troubleshoot hardened containers without compromising the
integrity of the original image. Because the debug container is ephemeral and
separate, it avoids introducing security risks into production environments.

If you encounter issues related to permissions, ports, missing shells, or
package managers, see [Troubleshoot Docker Hardened Images](../troubleshoot.md)
for recommended solutions and workarounds.