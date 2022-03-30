---
description: Docker Desktop for Linux (Beta)
keywords: docker, Desktop for linux, beta, tech preview
title: Docker Desktop for Linux (Beta)
toc_min: 1
toc_max: 2
---

Welcome to the Docker Desktop for Linux (Beta). The Beta program is aimed at early adopters who would like to try Docker Desktop for Linux and provide feedback.

Docker Desktop is an easy-to-install application that enables you to build and
share containerized applications and microservices. Docker Desktop for Linux
(DD4L) is the second-most popular feature request in our [public
roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank"
rel="noopener" class="_"}.

## Download and install

Docker Desktop for Linux (Beta) is currently available on Ubuntu 21.04, 21.10
and Debian distributions.

To install Docker Desktop for Linux:

1. Set up the [Docker repository](../../engine/install/ubuntu.md#install-using-the-repository).
2. Download and install the Debian package, if you have previously installed one of the preview releases it is a good idea to run `sudo apt remove docker-desktop`:
    ```console
    $ curl https://desktop-stage.docker.com/linux/main/amd64/76677/docker-desktop.deb --output docker-desktop.deb
    $ sudo apt install ./docker-desktop.deb
    ```

  There are a few post-install configuration steps done through the maintainers' scripts (post-install script contained
  in the deb package.

  The post-install script:

  - sets the capability on the Docker Desktop binary to map privileged ports and set resource limits
  - adds a DNS name for Kubernetes to `/etc/hosts`
  - creates a link from `/usr/bin/docker` to `/usr/local/bin/com.docker.cli`
  - installs systemd units for each user

## Launch Docker Desktop

To start Docker Desktop for Linux, search **Docker Desktop** on the
**Applications** menu and open it. This launches the whale menu icon and opens
the Docker Dashboard, reporting the status of Docker Desktop.

Alternatively, open a terminal and run:

```console
$ systemctl --user start docker-desktop
```

When Docker Desktop starts, it creates a dedicated context that the Docker CLI
can use as a target and sets it as the current context in use. This is to avoid
a clash with a local Docker Engine that may be running on the Linux host and
using the default context. On shutdown, Docker Desktop resets the current
context to the previous one.

The Docker Desktop installer updates Docker Compose and the Docker CLI binaries
on the host. It installs Docker Compose V2 and gives users the choice to
link it as docker-compose from the Settings panel. Docker Desktop installs
the new Docker CLI binary that includes cloud-integration capabilities in `/usr/local/bin`
and creates a symlink to the classic Docker CLi at `/usr/local/bin/com.docker.cli`.

After you’ve successfully installed Docker Desktop, you can check the versions
of these binaries by running the following command:

```console
$ docker compose version
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

> **Note:**
>
> Docker Desktop relies on `pass` to store credentials. Before signing in to
> Docker Hub from the Docker Dashboard or the Docker menu, you must initialize `pass`.
> Docker Desktop displays a warning if you've not initialized `pass`.

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

To create and upload a diagnostics bundle:

1. From the Docker menu, select **Troubleshoot** > **Get support**.
2. When the diagnostics are available, click **Upload to get a Diagnostic ID**.
3. Make a note of the Diagnostic ID displayed on the Support page. You can send
   this ID with your bug report to investigate any issues. Wait for a bundle to
   be generated, once uploaded, it displays a diagnostics ID that can be sent to
   us for investigation.

Or, if you prefer to investigate the issue, you can access Docker Desktop logs by running the following command:

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

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
$ sudo apt purge docker-desktop
```

## Known issue

At the end of the installation process, `apt` displays an error due to installing a downloaded package. You can ignore this error message.

  ```
  N: Download is performed unsandboxed as root, as file '/home/user/Downloads/docker-desktop.deb' couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)
  ```

If you have installed one of the previous releases and install the new package over it (as opposed to removing the old package explicitly), you will have to make sure that `~/.config/systemd/user/docker-desktop.service` and `~/.local/share/systemd/user/docker-desktop.service` are removed.

## Why Docker Desktop for Linux runs a VM

Docker Desktop for Linux runs a Virtual Machine (VM) for the following reasons:

1. **To ensure that Docker Desktop provides a consistent experience across platforms**.

    During research, the most frequently cited reason for users wanting Docker
    Desktop for Linux (DD4L) was to ensure a consistent Docker Desktop
    experience with feature parity across all major operating systems. Utilizing
    a VM ensures that the Docker Desktop experience for Linux users will closely
    match that of Windows and macOS.

    This need to deliver a consistent experience across all major OSs will become increasingly important as we look towards adding exciting new features, such as Docker Extensions, to Docker Desktop that will benefit users across all tiers.  We’ll provide more details on these at [DockerCon22](https://www.docker.com/dockercon/){: target="_blank" rel="noopener" class="_"}. Watch this space.

2. **To make use of new kernel features**

    Sometimes we want to make use of new operating system features. Because we control the kernel and the OS inside the VM, we can roll these out to all users immediately, even to users who are intentionally sticking on an LTS version of their machine OS.

3. **To enhance security**

    Container image vulnerabilities pose a security risk for the host environment. There is a large number of unofficial images that are not guaranteed to be verified for known vulnerabilities. Malicious users can push images to public registries and use different methods to trick users into pulling and running them. The VM approach mitigates this threat as any malware that gains root privileges is restricted to the VM environment without access to the host.

    Why not run rootless Docker? Although this has the benefit of superficially limiting access to the root user so everything looks safer in "top", it allows unprivileged users to gain `CAP_SYS_ADMIN` in their own user namespace and access kernel APIs which are not expecting to be used by unprivileged users, resulting in vulnerabilities like [this](https://www.openwall.com/lists/oss-security/2022/01/18/7){: target="_blank" rel="noopener" class="_"}.

4. **To provide the benefits of feature parity and enhanced security, with minimal impact on performance**

    The VM utilized by DD4L uses `virtiofs`, a shared file system that allows virtual machines to access a directory tree located on the host. Our internal benchmarking shows that with the right resource allocation to the VM, near native file system performance can be achieved with virtiofs.

    As such, we have adjusted the default memory available to the VM in DD4L. You can tweak this setting to your specific needs by using the **Memory** slider within the **Settings** > **Resources** tab of Docker Desktop.
