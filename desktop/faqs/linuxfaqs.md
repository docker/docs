---
description: Frequently asked questions
keywords: desktop, linux, faqs
title: Frequently asked questions for Linux
redirect_from:
- /desktop/linux/space/
---

## What is the difference between Docker Desktop for Linux and Docker Engine

Docker Desktop for Linux and Docker Engine can be installed side-by-side on the
same machine. Docker Desktop for Linux stores containers and images in an isolated
storage location within a VM and offers
controls to restrict [its resources](../settings/linux.md#resources). Using a dedicated storage
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

### How do I switch between Docker Desktop and Docker Engine
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
default *       Current DOCKER_HOST based configuration   unix:///var/run/docker.sock                      ...
desktop-linux                                             unix:///home/<user>/.docker/desktop/docker.sock  ...        
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

And use the `desktop-linux` context to interact with Docker Desktop:

```console
$ docker context use desktop-linux
desktop-linux
Current context is now "desktop-linux"
```

Refer to the [Docker Context documentation](../../engine/context/working-with-contexts.md) for more details.

## Why does Docker Desktop for Linux run a VM?

Docker Desktop for Linux runs a Virtual Machine (VM) for the following reasons:

1. **To ensure that Docker Desktop provides a consistent experience across platforms**.

    During research, the most frequently cited reason for users wanting Docker Desktop for Linux (DD4L) was to ensure a consistent Docker Desktop
    experience with feature parity across all major operating systems. Utilizing
    a VM ensures that the Docker Desktop experience for Linux users will closely
    match that of Windows and macOS.

2. **To make use of new kernel features**

    Sometimes we want to make use of new operating system features. Because we control the kernel and the OS inside the VM, we can roll these out to all users immediately, even to users who are intentionally sticking on an LTS version of their machine OS.

3. **To enhance security**

    Container image vulnerabilities pose a security risk for the host environment. There is a large number of unofficial images that are not guaranteed to be verified for known vulnerabilities. Malicious users can push images to public registries and use different methods to trick users into pulling and running them. The VM approach mitigates this threat as any malware that gains root privileges is restricted to the VM environment without access to the host.

    Why not run rootless Docker? Although this has the benefit of superficially limiting access to the root user so everything looks safer in "top", it allows unprivileged users to gain `CAP_SYS_ADMIN` in their own user namespace and access kernel APIs which are not expecting to be used by unprivileged users, resulting in [vulnerabilities](https://www.openwall.com/lists/oss-security/2022/01/18/7){: target="_blank" rel="noopener" class="_"}.

4. **To provide the benefits of feature parity and enhanced security, with minimal impact on performance**

    The VM utilized by DD4L uses [`virtiofs`](https://virtio-fs.gitlab.io){:target="_blank" rel="noopener" class="_"}, a shared file system that allows virtual machines to access a directory tree located on the host. Our internal benchmarking shows that with the right resource allocation to the VM, near native file system performance can be achieved with virtiofs.

    As such, we have adjusted the default memory available to the VM in DD4L. You can tweak this setting to your specific needs by using the **Memory** slider within the **Settings** > **Resources** tab of Docker Desktop.

## How do I enable file sharing?

Docker Desktop for Linux uses [virtiofs](https://virtio-fs.gitlab.io/){:target="_blank" rel="noopener"}{:target="_blank" rel="noopener"} as the
default (and currently only) mechanism to enable file sharing between the host
and Docker Desktop VM. In order not to require elevated privileges, without
unnecessarily restricting operations on the shared files, Docker Desktop runs
the file sharing service (`virtiofsd`) inside a user namespace (see
`user_namespaces(7)`) with UID and GID mapping configured. As a result Docker
Desktop relies on the host being configured to enable the current user to use
subordinate ID delegation. For this to be true `/etc/subuid` (see `subuid(5)`)
and `/etc/subgid` (see `subgid(5)`) must be present. Docker Desktop only
supports subordinate ID delegation configured via files. Docker Desktop maps the
current user ID and GID to 0 in the containers. It uses the first entry
corresponding to the current user in `/etc/subuid` and `/etc/subgid` to set up
mappings for IDs above 0 in the containers.

| ID in container | ID on host                                                                       |
| --------------- | -------------------------------------------------------------------------------- |
| 0 (root)        | ID of the user running DD (e.g. 1000)                                            |
| 1               | 0 + beginning of ID range specified in `/etc/subuid`/`/etc/subgid` (e.g. 100000) |
| 2               | 1 + beginning of ID range specified in `/etc/subuid`/`/etc/subgid` (e.g. 100001) |
| 3               | 2 + beginning of ID range specified in `/etc/subuid`/`/etc/subgid` (e.g. 100002) |
| ...             | ...                                                                              |

If `/etc/subuid` and `/etc/subgid` are missing, they need to be created.
Both should contain entries in the form -
`<username>:<start of id range>:<id range size>`. For example, to allow the current user
to use IDs from 100000 to 165535:

```console
$ grep "$USER" /etc/subuid >> /dev/null 2&>1 || (echo "$USER:100000:65536" | sudo tee -a /etc/subuid)
$ grep "$USER" /etc/subgid >> /dev/null 2&>1 || (echo "$USER:100000:65536" | sudo tee -a /etc/subgid)
```

To verify the configs have been created correctly, inspect their contents:

```console
$ echo $USER
exampleuser
$ cat /etc/subuid
exampleuser:100000:65536
$ cat /etc/subgid
exampleuser:100000:65536
```

In this scenario if a shared file is `chown`ed inside a Docker Desktop container
owned by a user with a UID of 1000, it shows up on the host as owned by
a user with a UID of 100999. This has the unfortunate side effect of preventing
easy access to such a file on the host. The problem is resolved by creating
a group with the new GID and adding our user to it, or by setting a recursive
ACL (see `setfacl(1)`) for folders shared with the Docker Desktop VM.

## Where does Docker Desktop store Linux containers?
Docker Desktop stores Linux containers and images in a single, large "disk image" file in the Linux filesystem. This is different from Docker on Linux, which usually stores containers and images in the `/var/lib/docker` directory on the host's filesystem.

### Where is the disk image file?

To locate the disk image file, select **Preferences** from the Docker Dashboard then **Advanced** from the **Resources** tab.

The **Advanced** tab displays the location of the disk image. It also displays the maximum size of the disk image and the actual space the disk image is consuming. Note that other tools might display space usage of the file in terms of the maximum file size, and not the actual file size.

#### What if the file is too large?

If the disk image file is too large, you can:

- Move it to a bigger drive
- Delete unnecessary containers and images
- Reduce the maximum allowable size of the file

#### How do I move the file to a bigger drive?

To move the disk image file to a different location:

1. Select **Preferences** then  **Advanced** from the **Resources** tab.

2. In the **Disk image location** section, click **Browse** and choose a new location for the disk image.

3. Click **Apply & Restart** for the changes to take effect.

Do not move the file directly in Finder as this can cause Docker Desktop to lose track of the file.

#### How do I delete unnecessary containers and images?

Check whether you have any unnecessary containers and images. If your client and daemon API are running version 1.25 or later (use the `docker version` command on the client to check your client and daemon API versions), you can see the detailed space usage information by running:

```console
$ docker system df -v
```

Alternatively, to list images, run:

```console
$ docker image ls
```

and then, to list containers, run:

```console
$ docker container ls -a
```

If there are lots of redundant objects, run the command:

```console
$ docker system prune
```

This command removes all stopped containers, unused networks, dangling images, and build cache.

It might take a few minutes to reclaim space on the host depending on the format of the disk image file:

- If the file is named `Docker.raw`: space on the host should be reclaimed within a few seconds.
- If the file is named `Docker.qcow2`: space will be freed by a background process after a few minutes.

Space is only freed when images are deleted. Space is not freed automatically when files are deleted inside running containers. To trigger a space reclamation at any point, run the command:

```console
$ docker run --privileged --pid=host docker/desktop-reclaim-space
```

Note that many tools report the maximum file size, not the actual file size.
To query the actual size of the file on the host from a terminal, run:

```console
$ cd ~/.docker/desktop/vms/0/data
$ ls -klsh Docker.raw
2333548 -rw-r--r--@ 1 username  staff    64G Dec 13 17:42 Docker.raw
```

In this example, the actual size of the disk is `2333548` KB, whereas the maximum size of the disk is `64` GB.

#### How do I reduce the maximum size of the file?

To reduce the maximum size of the disk image file:

1. From Docker Dashboard select **Preferences** then **Advanced** from the **Resources** tab.

2. The **Disk image size** section contains a slider that allows you to change the maximum size of the disk image. Adjust the slider to set a lower limit.

3. Click **Apply & Restart**.

When you reduce the maximum size, the current disk image file is deleted, and therefore, all containers and images will be lost.
