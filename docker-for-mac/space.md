---
description: Disk utilization
keywords: mac, disk
title: Disk utilization in Docker for Mac
---

Docker Desktop stores Linux containers and images in a single, large "disk image" file in the Mac filesystem. This is different from Docker on Linux, which usually stores containers and images in the `/var/lib/docker` directory.

## Where is the disk image file?

To locate the disk image file, select the Docker icon and then
**Preferences** > **Resources** > **Advanced**.

![Disk preferences](images/menu/prefs-advanced.png){:width="750px"}

The **Advanced** tab displays the location of the disk image. It also displays the maximum size of the disk image and the actual space the disk image is consuming. Note that other tools might display space usage of the file in terms of the maximum file size, and not the actual file size.

## If the file is too big

If the disk image file is too big, you can:

- move it to a bigger drive,
- delete unnecessary containers and images, or
- reduce the maximum allowable size of the file.

### Move the file to a bigger drive

To move the disk image file to a different location:

1. Select **Preferences** > **Resources** > **Advanced**.

2. In the **Disk image location** section, click **Browse** and choose a new location for the disk image.

3. Click **Apply & Restart** for the changes to take effect.

Do not move the file directly in Finder as this can cause Docker Desktop to lose track of the file.

### Delete unnecessary containers and images

Check whether you have any unnecessary containers and images. If your client and daemon API are running version 1.25 or later (use the `docker version` command on the client to check your client and daemon API versions), you can see the detailed space usage information by running:

```
docker system df -v
```

Alternatively, to list images, run:

```bash
$ docker image ls
```

and then, to list containers, run:

```bash
$ docker container ls -a
```

If there are lots of redundant objects, run the command:

```bash
$ docker system prune
```

This command removes all stopped containers, unused networks, dangling images, and build cache.

It might take a few minutes to reclaim space on the host depending on the format of the disk image file:

- If the file is named `Docker.raw`: space on the host should be reclaimed within a few seconds.
- If the file is named `Docker.qcow2`: space will be freed by a background process after a few minutes.

Space is only freed when images are deleted. Space is not freed automatically when files are deleted inside running containers. To trigger a space reclamation at any point, run the command:

```
$ docker run --privileged --pid=host docker/desktop-reclaim-space
```

Note that many tools report the maximum file size, not the actual file size.
To query the actual size of the file on the host from a terminal, run:

```bash
$ cd ~/Library/Containers/com.docker.docker/Data
$ cd vms/0/data
$ ls -klsh Docker.raw
2333548 -rw-r--r--@ 1 username  staff    64G Dec 13 17:42 Docker.raw
```

In this example, the actual size of the disk is `2333548` KB, whereas the maximum size of the disk is `64` GB.

### Reduce the maximum size of the file

To reduce the maximum size of the disk image file:

1. Select the Docker icon and then select **Preferences** > **Resources** > **Advanced**.

2. The **Disk image size** section contains a slider that allows you to change the maximum size of the disk image. Adjust the slider to set a lower limit.

3. Click **Apply & Restart**.

When you reduce the maximum size, the current disk image file is deleted, and therefore, all containers and images will be lost.
