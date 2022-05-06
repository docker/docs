---
description: Docker Desktop for Linux vs Docker Engine
keywords: linux, desktop, docker engine, docker desktop, docker desktop for linux, dd4l
title: Docker Desktop for Linux overview
---

This page provides an overview of Docker Desktop for Linuxand how it compares to
other products, such as the open source [Docker Engine for Linux](../../engine/index.md),
and [Docker Desktop for Mac](../mac/index.md) and [Docker Desktop for Windows](../windows/index.md).

## Differences between Docker Desktop for Linux and Docker Engine

Docker Desktop is an easy-to-install desktop application that enables you to
build and share containerized applications and microservices. The installation
includes Docker Engine, the Docker CLI client, Docker Compose, Kubernetes, and
a Credential Helper. Docker Desktop also comes with a graphical user interface
(Dashboard) to manage your containers, apps, volumes (and more), and a settings
panel, which allows you to configure your Docker Desktop installation, restrict
and manage its [resources](index.md#resources), enable [Kubernetes integration](index.md#kubernetes),
and [check for updates](index.md#software-updates).

## Docker Desktop and Docker Engine on the same machine

Docker Desktop for Linux and Docker Engine can be installed side-by-side on the
same machine. Docker Desktop for Linux stores containers and images in an isolated
storage location [within a VM](#why-docker-desktop-for-linux-runs-a-vm) and offers
controls to restrict [its resources](index.md#resources). Using a dedicated storage
location for Docker Desktop prevents it from interfering with a Docker Engine
installation on the same machine.

While it's possible to run both Docker Desktop and Docker Engine simultaneously,
there may be situations where running both at the same time can cause issues.
For example, when mapping network ports (`-p` / `--publish`) for containers, both
Docker Desktop and Docker Engine may attempt to reserve the same port on your
machine, which can lead to conflicts ("port already in use").

We generally recommend stopping the Docker Engine while you're using Docker Desktop
to prevent the Docker Engine from consuming resources and to prevent conflicts
as described above.
 
Use the following command to stop the Docker Engine service:

```console
$ sudo systemctl stop docker docker.socket containerd
```

Depending on your installation, the Docker Engine may be configured to automatically
start as a system service when your machine starts. Use the following command to
disable the Docker Engine service, and to prevent it from starting automatically:

```console
$ sudo systemctl disable docker docker.socket containerd
```

## Switch between Docker Desktop and Docker Engine
{: id="context" }

The Docker CLI can be used to interact with multiple Docker Engines. For example,
you can use the same Docker CLI to control a local Docker Engine and to control
a remote Docker Engine instance running in the cloud. [Docker Contexts](../../engine/context/working-with-contexts.md)
allow you to switch between Docker Engines instances.

When installing Docker Desktop, a dedicated "desktop-linux" context is created to
interact with Docker Desktop. On startup, Docker Desktop automatically sets its
own context (`desktop-linux`) as the current context. This means that subsequent
Docker CLI commands target Docker Desktop. On shutdown, Docker Desktop resets
the current context to the `default` context.

Use the `docker context ls` command to view what contexts are available on your
machine. The current context is indicated with an asterisk (`*`);

```console
$ docker context ls
NAME            DESCRIPTION                               DOCKER ENDPOINT                                  ...
default         Current DOCKER_HOST based configuration   unix:///var/run/docker.sock                      ...
desktop-linux *                                           unix:///home/<user>/.docker/desktop/docker.sock  ...
staging-server  Remote engine (staging)                   tcp://staging.mycorp.example.com:2376            ...
```

If you have both Docker Desktop and Docker Engine installed on the same machine,
you can run the `docker context use` command to switch between the Docker Desktop
and Docker Engine contexts. For example, use the "default" context to interact
with the Docker Engine;

```console
$ docker context use default
default
Current context is now "default"
```

And use the `docker-desktop` context to interact with Docker Desktop:

```console
$ docker context use docker-desktop
docker-desktop
Current context is now "docker-desktop"
```

Refer to the [Docker Context documentation](../../engine/context/working-with-contexts.md)
for more details.

## Why Docker Desktop for Linux runs a VM

Docker Desktop for Linux runs a Virtual Machine (VM) for the following reasons:

1. **To ensure that Docker Desktop provides a consistent experience across platforms**.

   During research, the most frequently cited reason for users wanting Docker
   Desktop for Linux (DD4L) was to ensure a consistent Docker Desktop
   experience with feature parity across all major operating systems. Utilizing
   a VM ensures that the Docker Desktop experience for Linux users will closely
   match that of Windows and macOS.

   This need to deliver a consistent experience across all major OSs becomes
   increasingly important as we look towards adding exciting new features, such
   as Docker Extensions, to Docker Desktop that benefits users across all tiers.
   We'll provide more details on these at [DockerCon22](https://www.docker.com/dockercon/){: target="_blank" rel="noopener" class="_"}.
   Watch this space.

3. **To make use of new kernel features**

   Sometimes we want to make use of new operating system features. Because we
   control the kernel and the OS inside the VM, we can roll these out to all
   users immediately, even to users who are intentionally sticking on an LTS
   version of their machine OS.

4. **To enhance security**

   Container image vulnerabilities pose a security risk for the host environment.
   There is a large number of unofficial images that are not guaranteed to be
   verified for known vulnerabilities. Malicious users can push images to public
   registries and use different methods to trick users into pulling and running
   them. The VM approach mitigates this threat as any malware that gains root
   privileges is restricted to the VM environment without access to the host.

   Why not run rootless Docker? Although this has the benefit of superficially
   limiting access to the root user so everything looks safer in "top", it
   allows unprivileged users to gain `CAP_SYS_ADMIN` in their own user namespace
   and access kernel APIs which are not expecting to be used by unprivileged
   users, resulting in [vulnerabilities](https://www.openwall.com/lists/oss-security/2022/01/18/7){:target="_blank" rel="noopener" class="_"}.

5. **To provide the benefits of feature parity and enhanced security, with minimal impact on performance**

   The VM utilized by DD4L uses [`virtiofs`](https://virtio-fs.gitlab.io){:target="_blank" rel="noopener" class="_"},
   a shared file system that allows virtual machines to access a directory tree
   located on the host. Our internal benchmarking shows that with the right
   resource allocation to the VM, near-native file system performance can be
   achieved with virtiofs.

   As such, we have adjusted the default memory available to the VM in DD4L. You
   can tweak this setting to your specific needs by using the **Memory** slider
   within the **Settings** > **Resources** tab of Docker Desktop.
