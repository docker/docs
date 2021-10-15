---
description: Instructions for installing Docker as a binary. Mostly meant for hackers who want to try out Docker on a variety of environments.
keywords: binaries, installation, docker, documentation, linux
title: Install Docker Engine from binaries
redirect_from:
- /engine/installation/binaries/
- /engine/installation/linux/docker-ce/binaries/
- /install/linux/docker-ce/binaries/
- /installation/binaries/
---

> **Important**
>
> This page contains information on how to install Docker using binaries. These
> instructions are mostly suitable for testing purposes. We do not recommend
> installing Docker using binaries in production environments as they will not be
> updated automatically with security updates. The Linux binaries described on this
> page are statically linked, which means that vulnerabilities in build-time
> dependencies are not automatically patched by security updates of your Linux
> distribution.
>
> Updating binaries is also slightly more involved when compared to Docker packages
> installed using a package manager or through Docker Desktop, as it requires
> (manually) updating the installed version whenever there is a new release of
> Docker.
>
> Also, static binaries may not include all functionalities provided by the dynamic
> packages.
>
> On Windows and Mac, we recommend that you install [Docker Desktop](../../desktop/index.md)
> instead. For Linux, we recommend that you follow the instructions specific for
> your distribution.
{: .important}

If you want to try Docker or use it in a testing environment, but you're not on
a supported platform, you can try installing from static binaries. If possible,
you should use packages built for your operating system, and use your operating
system's package management system to manage Docker installation and upgrades.

Static binaries for the Docker daemon binary are only available for Linux (as
`dockerd`) and Windows (as `dockerd.exe`).
Static binaries for the Docker client are available for Linux, Windows, and macOS (as `docker`).

This topic discusses binary installation for Linux, Windows, and macOS:

- [Install daemon and client binaries on Linux](#install-daemon-and-client-binaries-on-linux)
- [Install client binaries on macOS](#install-client-binaries-on-macos)
- [Install server and client binaries on Windows](#install-server-and-client-binaries-on-windows)

## Install daemon and client binaries on Linux

### Prerequisites

Before attempting to install Docker from binaries, be sure your host machine
meets the prerequisites:

- A 64-bit installation
- Version 3.10 or higher of the Linux kernel. The latest version of the kernel
  available for your platform is recommended.
- `iptables` version 1.4 or higher
- `git` version 1.7 or higher
- A `ps` executable, usually provided by `procps` or a similar package.
- [XZ Utils](https://tukaani.org/xz/) 4.9 or higher
- A [properly mounted](
  https://github.com/tianon/cgroupfs-mount/blob/master/cgroupfs-mount)
  `cgroupfs` hierarchy; a single, all-encompassing `cgroup` mount
  point is not sufficient. See Github issues
  [#2683](https://github.com/moby/moby/issues/2683),
  [#3485](https://github.com/moby/moby/issues/3485),
  [#4568](https://github.com/moby/moby/issues/4568)).

#### Secure your environment as much as possible

##### OS considerations

Enable SELinux or AppArmor if possible.

It is recommended to use AppArmor or SELinux if your Linux distribution supports
either of the two. This helps improve security and blocks certain
types of exploits. Review the documentation for your Linux distribution for
instructions for enabling and configuring AppArmor or SELinux.

> Security Warning
>
> If either of the security mechanisms is enabled, do not disable it as a
> work-around to make Docker or its containers run. Instead, configure it
> correctly to fix any problems.
{:.warning}

##### Docker daemon considerations

- Enable `seccomp` security profiles if possible. See
  [Enabling `seccomp` for Docker](../security/seccomp.md).

- Enable user namespaces if possible. See the
  [Daemon user namespace options](/engine/reference/commandline/dockerd/#daemon-user-namespace-options).

### Install static binaries

1.  Download the static binary archive. Go to
    [https://download.docker.com/linux/static/stable/](https://download.docker.com/linux/static/stable/),
    choose your hardware platform, and download the `.tgz` file relating to the
    version of Docker Engine you want to install.

2.  Extract the archive using the `tar` utility. The `dockerd` and `docker`
    binaries are extracted.

    ```console
    $ tar xzvf /path/to/<FILE>.tar.gz
    ```

3.  **Optional**: Move the binaries to a directory on your executable path, such
    as `/usr/bin/`. If you skip this step, you must provide the path to the
    executable when you invoke `docker` or `dockerd` commands.

    ```console
    $ sudo cp docker/* /usr/bin/
    ```

4.  Start the Docker daemon:

    ```console
    $ sudo dockerd &
    ```

    If you need to start the daemon with additional options, modify the above
    command accordingly or create and edit the file `/etc/docker/daemon.json`
    to add the custom configuration options.

5.  Verify that Docker is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

## Install client binaries on macOS

> **Note**
>
> The following instructions are mostly suitable for testing purposes. The macOS
> binary includes the Docker client only. It does not include the `dockerd` daemon
> which is required to run containers. Therefore, we recommend that you install
> [Docker Desktop](../../desktop/index.md) instead.

The binaries for Mac also do not contain:

- A runtime environment. You must set up a functional engine either in a Virtual Machine, or on a remote Linux machine.
- Docker components such as `buildx`, `docker scan`, and `docker compose`.

To install client binaries, perform the following steps:

1.  Download the static binary archive. Go to
    [https://download.docker.com/mac/static/stable/](https://download.docker.com/mac/static/stable/) and select `x86_64` (for Mac on Intel chip) or `aarch64` (for Mac on Apple silicon),
    and then download the `.tgz` file relating to the version of Docker Engine you want
    to install.

2.  Extract the archive using the `tar` utility. The `docker` binary is
    extracted.

    ```console
    $ tar xzvf /path/to/<FILE>.tar.gz
    ```

3.  Clear the extended attributes to allow it run.

    ```console
    $ sudo xattr -rc docker
    ```

    Now, when you run the following command, you can see the Docker CLI usage instructions:

    ```console
    $ docker/docker
    ```

4.  **Optional**: Move the binary to a directory on your executable path, such
    as `/usr/local/bin/`. If you skip this step, you must provide the path to the
    executable when you invoke `docker` or `dockerd` commands.

    ```console
    $ sudo cp docker/docker /usr/local/bin/
    ```

5.  Verify that Docker is installed correctly by running the `hello-world`
    image. The value of `<hostname>` is a hostname or IP address running the
    Docker daemon and accessible to the client.

    ```console
    $ sudo docker -H <hostname> run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

## Install server and client binaries on Windows

> **Note**
>
> The following section describes how to install the Docker daemon on Windows
> Server which allows you to run Windows containers only. The binaries for
> Windows do not contain Docker components such as `buildx`, `docker scan`, and
> `docker compose`. If you are running Windows 10 or 11, we recommend that you
> install [Docker Desktop](../../desktop/index.md) instead.

Binary packages on Windows include both `dockerd.exe` and `docker.exe`. On Windows,
these binaries only provide the ability to run native Windows containers (not
Linux containers).

To install server and client binaries, perform the following steps:

1. Download the static binary archive. Go to
    [https://download.docker.com/win/static/stable/x86_64](https://download.docker.com/win/static/stable/x86_64) and select the latest version from the list.

2. Run the following PowerShell commands to install and extract the archive to your program files:

    ```powershell
    PS C:\> Expand-Archive /path/to/<FILE>.zip -DestinationPath $Env:ProgramFiles
    ```

3. Register the service and start the Docker Engine:

    ```powershell
    PS C:\> $Env:ProgramFiles\Docker\dockerd --register-service
    PS C:\> Start-Service docker
    ```

4.  Verify that Docker is installed correctly by running the `hello-world`
    image.

    ```powershell
    PS C:\> $Env:ProgramFiles\Docker\docker run hello-world:nanoserver
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

## Upgrade static binaries

To upgrade your manual installation of Docker Engine, first stop any
`dockerd` or `dockerd.exe`  processes running locally, then follow the
regular installation steps to install the new version on top of the existing
version.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
- Take a look at the [Get started](../../get-started/index.md) training modules to learn  how to build an image and run it as a containerized application.
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
