---
description: Understand permission requirements for Docker Desktop for Mac and the
  differences between versions
keywords: Docker Desktop, mac, security, install, permissions
title: Understand permission requirements for Docker Desktop on Mac
linkTitle: Mac permission requirements
aliases:
- /desktop/mac/permission-requirements/
weight: 20
---

This page contains information about the permission requirements for running and installing Docker Desktop on Mac.

It also provides clarity on running containers as `root` as opposed to having `root` access on the host.

Docker Desktop on Mac is designed with security in mind. Administrative rights are only required when absolutely necessary.

## Permission requirements

Docker Desktop for Mac is run as an unprivileged user. However, Docker Desktop requires certain functionalities to perform a limited set of privileged configurations such as:
 - [Installing symlinks](#installing-symlinks) in`/usr/local/bin`.
 - [Ensuring `localhost` and `kubernetes.docker.internal` are defined](#ensuring-localhost-and-kubernetesdockerinternal-are-defined) in `/etc/hosts`. Some old macOS installs don't have `localhost` in `/etc/hosts`, which causes Docker to fail. Defining the DNS name `kubernetes.docker.internal` allows Docker to share Kubernetes contexts with containers.

Privileged access, when needed, is granted through the **Advanced** page in **Settings**, rather than at install time.

From the **Advanced** page, you can enable:
- The location of the Docker CLI tools either in the system or user directory
- The default Docker socket

If you work in an environment with elevated security requirements, for instance where local administrative access is prohibited, then you can leave these settings at their defaults to avoid granting privileged access.

Depending on which advanced settings you configure, you must enter your password to confirm.

### Installing symlinks

The Docker binaries are installed by default in `/Applications/Docker.app/Contents/Resources/bin`. By default, Docker Desktop creates symlinks for the binaries in `$HOME/.docker/bin`, which is automatically added to your terminal's `PATH`.

You can change the symlink location to `/usr/local/bin` from the **Advanced** page in **Settings**. If `/usr/local/bin` is not writable by unprivileged users, Docker Desktop requires authorization to confirm this choice before the symlinks to Docker binaries are created there.

You are also given the option to enable the installation of the `/var/run/docker.sock` symlink from the **Advanced** page in **Settings**. Creating this symlink ensures various Docker clients relying on the default Docker socket path work without additional changes.

As the `/var/run` is mounted as a tmpfs, its content is deleted on restart, symlink to the Docker socket included. To ensure the Docker socket exists after restart, Docker Desktop sets up a `launchd` startup task that creates the symlink by running `ln -s -f /Users/<user>/.docker/run/docker.sock /var/run/docker.sock`. This ensures the you aren't prompted on each startup to create the symlink. If you don't enable this option, the symlink and the startup task is not created and you may have to explicitly set the `DOCKER_HOST` environment variable to `/Users/<user>/.docker/run/docker.sock` in the clients it is using. The Docker CLI relies on the current context to retrieve the socket path, the current context is set to `desktop-linux` on Docker Desktop startup.

### Ensuring `localhost` and `kubernetes.docker.internal` are defined

It is your responsibility to ensure that localhost is resolved to `127.0.0.1` and if Kubernetes is used, that `kubernetes.docker.internal` is resolved to `127.0.0.1`.

## Installing from the command line

Privileged configurations are applied during the installation with the `--user` flag on the [install command](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line). In this case, you are not prompted to grant root privileges on the first run of Docker Desktop. Specifically, the `--user` flag:
- Uninstalls the previous `com.docker.vmnetd` if present
- Sets up symlinks
- Ensures that `localhost` is resolved to `127.0.0.1`

The limitation of this approach is that Docker Desktop can only be run by one user-account per machine, namely the one specified in the `-–user` flag.

## Backend helper socket

The Docker Desktop backend process (`com.docker.backend`) uses an internal helper socket
(`~/Library/Containers/com.docker.docker/Data/forkexecd.sock`) to fork and execute
helper processes as part of running Docker Desktop.

This socket does not run as `root` and grants no elevated privileges. It is owned by,
and accessible only to, the same macOS user running Docker Desktop, and is contained
in Docker Desktop's application container.

## Containers running as root within the Linux VM

With Docker Desktop, the Docker daemon and containers run in a lightweight Linux
VM managed by Docker. This means that although containers run by default as
`root`, this doesn't grant `root` access to the Mac host machine. The Linux VM
serves as a security boundary and limits what resources can be accessed from the
host. Any directories from the host bind mounted into Docker containers still
retain their original permissions.

## Enhanced Container Isolation

In addition, Docker Desktop supports [Enhanced Container Isolation
mode](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/_index.md) (ECI),
available to Business customers only, which further secures containers without
impacting developer workflows.

ECI automatically runs all containers within a Linux user-namespace, such that
root in the container is mapped to an unprivileged user inside the Docker
Desktop VM. ECI uses this and other advanced techniques to further secure
containers within the Docker Desktop Linux VM, such that they are further
isolated from the Docker daemon and other services running inside the VM.
