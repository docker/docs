---
description: Instructions for installing Docker as a binary. Mostly meant for hackers who want to try out Docker on a variety of environments.
keywords: binaries, installation, docker, documentation, linux
title: Install Docker from binaries
---

> **Note**: You may have been redirected to this page because there is no longer
> a dynamically-linked Docker package for your Linux distribution.

If you want to try Docker or use it in a testing environment, but you're not on
a supported platform, you can try installing from static binaries. If possible,
you should use packages built for your operating system, and use your operating
system's package management system to manage Docker installation and upgrades.
Be aware that 32-bit static binary archives do not include the Docker daemon.

Static binaries for the Docker daemon binary are only available for Linux (as
`dockerd`) and Windows Server 2016 or Windows 10 (as `dockerd.exe`). Static
binaries for the Docker client are available for Linux and macOS (as `docker`),
and Windows Server 2016 or Windows 10 (as `docker.exe`).

## Install daemon and client binaries on Linux

### Prerequisites

Before attempting to install Docker from binaries, be sure your host machine
meets the prerequisites:

- A 64-bit installation
- Version 3.10 or higher of the Linux kernel. The latest version of the kernel
  available for you platform is recommended.
- `iptables` version 1.4 or higher
- `git` version 1.7 or higher
- A `ps` executable, usually provided by `procps` or a similar package.
- [XZ Utils](http://tukaani.org/xz/) 4.9 or higher
- a [properly mounted](
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

> **Warning**:
> If either of the security mechanisms is enabled, do not disable it as a
> work-around to make Docker or its containers run. Instead, configure it
> correctly to fix any problems.
{:.warning}

##### Docker daemon considerations

- Enable `seccomp` security profiles if possible. See
  [Enabling `seccomp` for Docker](/engine/security/seccomp.md).

- Enable user namespaces if possible. See the
  [Daemon user namespace options](/engine/reference/commandline/dockerd/#/daemon-user-namespace-options).

### Install static binaries

1.  Download the static binary archive. You can download either the latest
    release binaries or a specific version. To find the download link, see the
    [release notes](https://github.com/moby/moby/releases) for the version
    of Docker you want to install. You can choose a `tar.gz` archive or `zip`
    archive.

2.  Extract the archive using `tar` or `unzip`, depending on the format you
    downloaded. The `dockerd` and `docker` binaries are extracted.

    ```bash
    $ tar xzvf /path/to/<FILE>.tar.gz
    ```

    ```bash
    $ unzip /path/to/<FILE>.zip
    ```

3.  **Optional**: Move the binaries to a directory on your executable path, such
    as `/usr/bin/`. If you skip this step, you must provide the path to the
    executable when you invoke `docker` or `dockerd` commands.

    ```bash
    $ sudo cp docker/* /usr/bin/
    ```

4.  Start the Docker daemon:

    ```bash
    $ sudo dockerd &
    ```

    If you need to start the daemon with additional options, modify the above
    command accordingly.

5.  Verify that Docker is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

### Next steps

- Continue to [Post-installation steps for Linux](/engine/installation/linux/linux-postinstall.md)

- Continue with the [User Guide](/engine/userguide/index.md).

## Install client binaries on macOS

The macOS binary includes the Docker client only. It does not include the
`dockerd` daemon.

1.  Download the static binary archive. You can download either the latest
    release binaries or a specific version. To find the download link, see the
    [release notes](https://github.com/moby/moby/releases) for the version
    of Docker you want to install. You can choose a `tar.gz` archive or
    `zip` archive.

2.  Extract the archive using `tar` or `unzip`, depending on the format you
    downloaded. The `docker` binary is extracted.

    ```bash
    $ tar xzvf /path/to/<FILE>.tar.gz
    ```

    ```bash
    $ unzip /path/to/<FILE>.zip
    ```
3.  **Optional**: Move the binaries to a directory on your executable path, such
    as `/usr/local/bin/`. If you skip this step, you must provide the path to the
    executable when you invoke `docker` or `dockerd` commands.

    ```bash
    $ sudo cp docker/docker /usr/local/bin/
    ```

4.  Verify that Docker is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker -H <hostname> run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.


## Install server and client binaries on Windows

You can install Docker from binaries on Windows Server 2016 or Windows 10.

- To install both client and server binaries, download the 64-bit binary. The
  archive includes `x86_64` in the path.

- To install the client only, download the 32-bit binary. The archive includes
  `i386` in the path.

1.  Use the following PowerShell commands to install and start Docker:

    ```none
    Invoke-WebRequest https://get.docker.com/builds/Windows/x86_64/docker-17.03.0-ce.zip -UseBasicParsing -OutFile docker.zip
    Expand-Archive docker.zip -DestinationPath $Env:ProgramFiles
    Remove-Item -Force docker.zip

    dockerd --register-service

    Start-Service docker
    ```

2.  Verify that Docker is installed correctly by running the `hello-world`
    image.


    ```none
    docker run hello-world:nanoserver
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

## Upgrade static binaries

To upgrade your manual installation of Docker Engine on Linux, first stop any
`dockerd` processes running locally, then follow the
[regular installation steps](#get-the-linux-binaries), overwriting any existing
`dockerd` or `docker` binaries with the newer versions.

## Next steps

Continue with the [User Guide](../userguide/index.md).
