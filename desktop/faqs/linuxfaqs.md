---
description: Frequently asked questions
keywords: desktop, linux, faqs
title: Frequently asked questions for Linux
redirect_from:
- /desktop/linux/space/
---

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
