---
title: Troubleshoot
description: Resolve common issues when building, running, or debugging Docker Hardened Images, such as non-root behavior, missing shells, and port access.
weight: 40
tags: [Troubleshooting]
keywords: troubleshoot hardened image, docker debug container, non-root permission issue, missing shell error, no package manager, debug, hardened images, DHI, troubleshooting, ephemeral container, docker debug, non-root containers, hardened container image, debug secure container
aliases:
- /dhi/how-to/debug/
---

This page covers debugging techniques and common issues you may encounter while
migrating to or using Docker Hardened Images (DHIs).

## General debugging

Docker Hardened Images prioritize minimalism and security, which means
they intentionally leave out many common debugging tools (like shells or package
managers). This makes direct troubleshooting difficult without introducing risk.
To address this, you can use [Docker
Debug](/reference/cli/docker/debug/), a secure workflow that
temporarily attaches an ephemeral debug container to a running service or image
without modifying the original image.

This section shows how to debug Docker Hardened Images locally during development.
With Docker Debug, you can also debug containers remotely using the `--host`
option.

### Use Docker Debug

#### Step 1: Run a container from a Hardened Image

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

#### Step 2: Use Docker Debug to inspect the container

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

### Alternative debugging approaches

In addition to using Docker Debug, you can also use the following approaches for
debugging DHI containers.

#### Use the -dev variant

Docker Hardened Images offer a `-dev` variant that includes a shell
and a package manager to install debugging tools. Simply replace the image tag
with `-dev`:

```console
$ docker run -it --rm dhi.io/python:3.13-dev sh
```

Type `exit` to leave the container when done. Note that using the `-dev` variant
increases the attack surface and it is not recommended as a runtime for
production environments.

#### Mount debugging tools with image mounts

You can use the image mount feature to mount debugging tools into your container
without modifying the base image.

##### Step 1: Run a container from a hardened image

Start with a DHI-based container that simulates an issue:

```console
$ docker run -d --name myapp dhi.io/python:3.13 python -c "import time; time.sleep(300)"
```

##### Step 2: Mount debugging tools into the container

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

## Common issues

The following are specific issues you may encounter when working with Docker
Hardened Images, along with recommended solutions.

### Permissions

DHIs run as a nonroot user by default for enhanced security. This can result in
permission issues when accessing files or directories. Ensure your application
files and runtime directories are owned by the expected UID/GID or have
appropriate permissions.

To find out which user a DHI runs as, check the repository page for the image on
Docker Hub. See [View image variant
details](./how-to/explore.md#image-variant-details) for more information.

### Privileged ports

Nonroot containers cannot bind to ports below 1024 by default. This is enforced
by both the container runtime and the kernel (especially in Kubernetes and
Docker Engine < 20.10).

Inside the container, configure your application to listen on an unprivileged
port (1025 or higher). For example `docker run -p 80:8080 my-image` maps
port 8080 in the container to port 80 on the host, allowing you to access it
without needing root privileges.

### No shell

Runtime DHIs omit interactive shells like `sh` or `bash`. If your build or
tooling assumes a shell is present (e.g., for `RUN` instructions), use a `dev`
variant of the image in an earlier build stage and copy the final artifact into
the runtime image.

To find out which shell, if any, a DHI has, check the repository page for the
image on Docker Hub. See [View image variant
details](./how-to/explore.md#image-variant-details) for more information.

Also, use Docker Debug when you need shell access to a running container. For
more details, see [General debugging](#general-debugging).

### Entry point differences

DHIs may define different entry points compared to Docker Official Images (DOIs)
or other community images.

To find out the ENTRYPOINT or CMD for a DHI, check the repository page for the
image on Docker Hub. See [View image variant
details](./how-to/explore.md#image-variant-details) for more information.

### No package manager

Runtime Docker Hardened Images are stripped down for security and minimal attack
surface. As a result, they don't include a package manager such as `apk` or
`apt`. This means you can't install additional software directly in the runtime
image.

If your build or application setup requires installing packages (for example, to
compile code, install runtime dependencies, or add diagnostic tools), use a `dev`
variant of the image in a build stage. Then, copy only the necessary artifacts
into the final runtime image.