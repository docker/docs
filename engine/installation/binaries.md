---
description: Instructions for installing Docker as a binary. Mostly meant for hackers who want to try out Docker on a variety of environments.
keywords: binaries, installation, docker, documentation, linux
title: Install Docker from binaries
---

**This instruction set is meant for hackers who want to try out Docker
on a variety of environments.**

Before following these directions, you should really check if a packaged
version of Docker is already available for your distribution. We have
packages for many distributions, and more keep showing up all the time!

## Check runtime dependencies

To run properly, docker needs the following software to be installed at
runtime:

 - iptables version 1.4 or later
 - Git version 1.7 or later
 - procps (or similar provider of a "ps" executable)
 - XZ Utils 4.9 or later
 - a [properly mounted](
   https://github.com/tianon/cgroupfs-mount/blob/master/cgroupfs-mount)
   cgroupfs hierarchy (having a single, all-encompassing "cgroup" mount
   point [is](https://github.com/docker/docker/issues/2683)
   [not](https://github.com/docker/docker/issues/3485)
   [sufficient](https://github.com/docker/docker/issues/4568))

## Check kernel dependencies

Docker in daemon mode has specific kernel requirements. For details,
check your distribution in [*Installation*](index.md#on-linux).

A 3.10 Linux kernel is the minimum requirement for Docker.
Kernels older than 3.10 lack some of the features required to run Docker
containers. These older versions are known to have bugs which cause data loss
and frequently panic under certain conditions.

The latest minor version (3.x.y) of the 3.10 (or a newer maintained version)
Linux kernel is recommended. Keeping the kernel up to date with the latest
minor version will ensure critical kernel bugs get fixed.

> **Warning**:
> Installing custom kernels and kernel packages is probably not
> supported by your Linux distribution's vendor. Please make sure to
> ask your vendor about Docker support first before attempting to
> install custom kernels on your distribution.

> **Warning**:
> Installing a newer kernel might not be enough for some distributions
> which provide packages which are too old or incompatible with
> newer kernels.

Note that Docker also has a client mode, which can run on virtually any
Linux kernel (it even builds on macOS!).

## Enable AppArmor and SELinux when possible

Please use AppArmor or SELinux if your Linux distribution supports
either of the two. This helps improve security and blocks certain
types of exploits. Your distribution's documentation should provide
detailed steps on how to enable the recommended security mechanism.

Some Linux distributions enable AppArmor or SELinux by default and
they run a kernel which doesn't meet the minimum requirements (3.10
or newer). Updating the kernel to 3.10 or newer on such a system
might not be enough to start Docker and run containers.
Incompatibilities between the version of AppArmor/SELinux user
space utilities provided by the system and the kernel could prevent
Docker from running, from starting containers or, cause containers to
exhibit unexpected behaviour.

> **Warning**:
> If either of the security mechanisms is enabled, it should not be
> disabled to make Docker or its containers run. This will reduce
> security in that environment, lose support from the distribution's
> vendor for the system, and might break regulations and security
> policies in heavily regulated environments.

## Get the Docker Engine binaries

You can download either the latest release binaries or a specific version. View
the `docker/docker` [Releases page](https://github.com/docker/docker/releases).

A group of download links is included at the bottom of the release notes for
each version of Docker. You can use these links to download the source code
archive for that release, binaries for supported platforms, and static binaries
for unsupported Linux platforms. Use the links listed in the Downloads section
to download the appropriate binaries.

### Limitations of Windows and macOS binaries

For Windows, the `i386` download contains a 32-bit client-only binary. The
`x86_64` download contains both client and daemon binaries for 64-bit Windows
Server 2016 and Windows 10 systems.

The macOS binary only contains a client. You cannot use it to run the `dockerd`
daemon. If you need to run the daemon, install
[Docker for Mac](/docker-for-mac/index.md) instead.

### URL patterns for static binaries

The URLs for released binaries are stable and follow a predictable pattern.
Unfortunately, it is not possible to browse the releases in a directory
structure. If you do not want to get the links from the release notes for a
release, you can infer the URL for the binaries by using the following patterns:

| Description            | URL pattern                                                       |
|------------------------|-------------------------------------------------------------------|
| Latest Linux 64-bit    | `https://get.docker.com/builds/Linux/x86_64/docker-latest.tgz`    |
| Latest Linux 32-bit    | `https://get.docker.com/builds/Linux/i386/docker-latest.tgz`      |
| Specific version Linux 64-bit| `https://get.docker.com/builds/Linux/x86_64/docker-<version>.tgz` |
| Specific version Linux 32-bit| `https://get.docker.com/builds/Linux/i386/docker-<version>.tgz`   |
| Latest Windows 64-bit | `https://get.docker.com/builds/Windows/x86_64/docker-latest.zip`     |
| Latest Windows 32-bit | `https://get.docker.com/builds/Windows/i386/docker-latest.zip`      |
| Specific version Windows 64-bit | `https://get.docker.com/builds/Windows/x86_64/docker-<version>.zip` |
| Specific version Windows 32-bit | `https://get.docker.com/builds/Windows/i386/docker-<version>.zip` |
| Latest MacOS 64-bit   | `https://get.docker.com/builds/Darwin/x86_64/docker-latest.tgz` |
| Specific version MacOS 64-bit | `https://get.docker.com/builds/Darwin/x86_64/docker-<version>.tgz` |

For example, the stable URL for the Docker 1.11.0 64-bit static binary for Linux
is `https://get.docker.com/builds/Linux/x86_64/docker-1.11.0.tgz`.

> **Note** These instructions are for Docker Engine 1.11 and up. Engine 1.10 and
> under consists of a single binary, and instructions for those versions are
> different. To install version 1.10 or below, follow the instructions in the
> [1.10 documentation](/v1.10/engine/installation/binaries/){:target="_blank"}.

#### Verify downloaded files

To verify the integrity of downloaded files, you can get an MD5 or SHA256
checksum by adding `.md5` or `.sha256` to the end of the URL. For instance,
to verify the `docker-1.11.0.tgz` link above, use the URL
`https://get.docker.com/builds/Linux/x86_64/docker-1.11.0.tgz.md5` or
`https://get.docker.com/builds/Linux/x86_64/docker-1.11.0.tgz.sha256`.

## Install the Linux binaries

After downloading, you extract the archive, which puts the binaries in a
directory named `docker` in your current location.

```bash
$ tar -xvzf docker-latest.tgz

docker/
docker/docker
docker/docker-containerd
docker/docker-containerd-ctr
docker/docker-containerd-shim
docker/docker-proxy
docker/docker-runc
docker/dockerd
```

Engine requires these binaries to be installed in your host's `$PATH`.
For example, to install the binaries in `/usr/bin`:

```bash
$ mv docker/* /usr/bin/
```

> **Note**: If you already have Engine installed on your host, make sure you
> stop Engine before installing (`killall docker`), and install the binaries
> in the same location. You can find the location of the current installation
> with `dirname $(which docker)`.

### Run the Engine daemon on Linux

You can manually start the Engine in daemon mode using:

```bash
$ sudo dockerd &
```

The GitHub repository provides samples of init-scripts you can use to control
the daemon through a process manager, such as upstart or systemd. You can find
these scripts in the <a href="https://github.com/docker/docker/tree/master/contrib/init">
contrib directory</a>.

For additional information about running the Engine in daemon mode, refer to
the [daemon command](../reference/commandline/dockerd.md) in the Engine command
line reference.

## Install the macOS binaries

You can extract the downloaded archive either by double-clicking the downloaded
`.tgz` or on the command line, using `tar -xvzf docker-1.11.0.tgz`. You can run
the client binary from any location on your filesystem.

## Install the Windows binary

You can extract the downloaded archive by double-clicking the downloaded
`.zip`. You can run the client binary from any location on your filesystem.

## Run client commands without root access

On Linux, the `dockerd` daemon always runs as the root user and binds to a Unix
socket instead of a TCP port. By default that Unix socket is owned by the
`root` user. This means that by default, you need to use `sudo` to run `docker`
commands.

If you (or your Docker installer) create a Unix group called `docker` and add
users to it, the `dockerd` daemon will change the ownership of the Unix socket
to be readable and writable by members of the `docker` group when the daemon
starts. The `dockerd` daemon must always run as the root user, but you can run
`docker` client commands, such as `docker run`, as a non-privileged user.

> **Warning**:
> Membership in the *docker* group (or the group specified with `-G`) is equivalent
> to `root` access. See
> [*Docker Daemon Attack Surface*](../security/security.md#docker-daemon-attack-surface) details.

## Upgrade Docker Engine

Before you upgrade your manual installation of Docker Engine on Linux, first
stop the docker daemon:

```bash
$ killall dockerd
```

After the Docker daemon stops, move the old binaries out of the way and follow
the [regular installation steps](binaries.md#get-the-linux-binaries).

## Next steps

Continue with the [User Guide](../userguide/index.md).
