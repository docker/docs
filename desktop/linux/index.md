---
description: Docker Desktop for Linux Tech Preview
keywords: docker, linux, tech preview
title: Docker Desktop for Linux (Tech Preview)
toc_min: 1
toc_max: 2
---

Welcome to the Docker Desktop for Linux Tech Preview. This Tech Preview is aimed at early adopters who would like to try an experimental build of Docker Desktop for Linux and provide feedback.

Docker Desktop is an easy-to-install application that enables you to build and
share containerized applications and microservices. Docker Desktop for Linux
(DD4L) is the second-most popular feature request in our [public
roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank"
rel="noopener" class="_"}.

## Download and install

Docker Desktop for Linux is currently available on Ubuntu 21.04, 21.10 and
Debian distributions.

To install Docker Desktop for Linux:

1. Set up the [Docker repository](../../engine/install/ubuntu.md#install-using-the-repository).
2. Download and install the Tech Preview Debian package:
    ```console
    $ curl https://desktop-stage.docker.com/linux/main/amd64/73772/docker-desktop.deb --output docker-desktop.deb
    $ sudo apt install ./docker-desktop.deb
    ```
3. Check whether the user belongs to `docker` and `kvm` groups. You may need to restart the host to load the group configuration.

## Check the shared memory

Before you run Docker Desktop for Linux, verify whether the shared memory available on the host is **higher** than the memory allocated to the VM. By default, Docker Desktop allocates half of the memory and CPU from the host. The **available shared memory** should be higher than this.


```console
$ df -h /dev/shm
Filesystem      Size  Used Avail Use% Mounted on
tmpfs            16G  200M   16G   2% /dev/shm
```

To set the shared memory size, run:

```console
$ sudo mount -o remount,size=<the-size-you-want-in-GB> /dev/shm
```

To ensure this setting persists after a reboot, add the following entry to the `/etc/fstab`:

```console
none    /dev/shm    tmpfs   defaults,size=<the-size-you-want-in-GB>   0   0
```

For example:

```console
none    /dev/shm    tmpfs   defaults,size=8G    0   0
```

## Launch Docker Desktop

> **Note:** 
> 
> You may need to restart the host to load the group configuration.

To start Docker Desktop for Linux, search **Docker Desktop** on the
**Applications** menu and open it. This launches the whale menu icon and opens
the Docker Dashboard, reporting the status of Docker Desktop.

Alternatively, open a terminal and run:

```console
$ systemctl --user start docker-desktop
```

When Docker Desktop starts, it creates a dedicated context that the Docker CLI can use as a target. This is to avoid a clash with a local Docker Engine that may be running on the Linux host and using the default context.

Run the following command to switch to the desktop-linux context.

```console
 $ docker context use desktop-linux
```

The Docker Desktop installer updates Docker Compose and the Docker CLI binaries
on the host. It installs Docker Compose V2 as the default Docker Compose. It
also replaces the default Docker CLI with a new Docker CLI binary that includes
cloud-integration capabilities.

After you’ve successfully installed Docker Desktop, you can check the versions
of these binaries by running the following command:

```console
$ docker-compose version
Docker Compose version v2.2.3

$ docker --version
Docker version 20.10.12, build e91ed57

$ docker version
Client: Docker Engine - Community
Cloud integration: 1.0.17
Version:           20.10.12
API version:       1.41
...
```

To enable Docker Desktop to start on login, from the Docker menu, select
**Settings** > **General** > **Start Docker Desktop when you log in**.

Alternatively, open a terminal and run:

```console
$ systemctl --user enable docker-desktop
```

To stop Docker Desktop, click on the whale menu tray icon to open the Docker menu and select **Quit Docker Desktop**.

Alternatively, open a terminal and run:

```console
$ systemctl --user stop docker-desktop
```

## Logs

If you experience any issues, you can access Docker Desktop logs by running the following command:

```console
$ journalctl --user --unit=docker-desktop
```

You can also find additional logs for the internal components included in Docker
Desktop at `$HOME/.docker/desktop/log/`.

## Uninstall

To remove Docker Desktop for Linux, run:

```console
$ sudo apt remove docker-desktop
```


## Known issues

 - `Reset to factory defaults` currently does not work
 - At the end of installation, apt produces an error due to installing a downloaded package. This should be ignored.
  ```
  N: Download is performed unsandboxed as root, as file '/home/user/Downloads/docker-desktop.deb' couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)
  ```
  

