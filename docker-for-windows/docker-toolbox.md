---
description: Docker Desktop for Windows and Docker Toolbox
keywords: windows, alpha, beta, toolbox, docker-machine, tutorial
title: Migrate Docker Toolbox
---

This page explains how to migrate your Docker Toolbox disk image, or images if
you have them, to Docker Desktop for Windows.

## How to migrate Docker Toolbox disk images to Docker Desktop

Docker Desktop does not propose Toolbox image migration as part of its
installer since version 18.01.0. You can migrate existing Docker
Toolbox images with the steps described below.

In a terminal, while running Toolbox, use `docker commit` to create an image snapshot
from a container, for each container you wish to preserve:

```
> docker commit nginx
sha256:1bc0ee792d144f0f9a1b926b862dc88b0206364b0931be700a313111025df022
```

Next, export each of these images (and any other images you wish to keep):

```
> docker save -o nginx.tar sha256:1bc0ee792d144f0f9a1b926b862dc88b0206364b0931be700a313111025df022
```

Next, when running Docker Desktop on Windows, reload all these images:

```
> docker load -i nginx.tar
Loaded image ID: sha256:1bc0ee792d144f0f9a1b926b862dc88b0206364b0931be700a313111025df022
```

Note these steps will not migrate any `docker volume` contents: these must
be copied across manually.

## How to uninstall Docker Toolbox

Whether or not you migrate your Docker Toolbox images, you may decide to
uninstall it. For details on how to perform a clean uninstall of Toolbox,
see [How to uninstall Toolbox](../toolbox/toolbox_install_windows.md#how-to-uninstall-toolbox){: target="_blank" class="_"}.
