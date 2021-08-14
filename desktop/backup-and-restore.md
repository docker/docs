---
title: Back up and restore data
keywords: Docker Desktop, backup, restore, migration, reinstall, containers, images, volumes
---

You can use the following procedure to save and restore images and container data.
For example to reset your VM disk or to move your Docker environment to a new
computer.

## Save your data

1. If you have containers that contain data that must be backed up, commit those
   containers to an image with [`docker container commit`](../engine/reference/commandline/commit.md).

   Committing a container stores the container filesystem changes and some of the
   container's configuration (labels, environment-variables, command/entrypoint)
   as a local image. Be aware that environment variables may contain sensitive
   information such as passwords or proxy-authentication, so care should be taken
   when pushing the resulting image to a registry.

   Also note that filesystem changes in volume that are attached to the
   container are not included in the image, and must be backed up separately
   (see step 3 below).

   Refer to the [`docker container commit` page](../engine/reference/commandline/commit.md)
   in the Docker Engine command line reference section for details on using this
   command.

   > Should I back up my containers?
   >
   > If you use volumes or bind-mounts to store your container data, backing up
   > your containers may not be needed, but make sure to remember the options that
   > were used when creating the container or use a [Docker Compose file](../compose/compose-file/index.md)
   > if you want to re-create your containers with the same configuration after
   > re-installation.

2. Use [`docker push`](../engine/reference/commandline/push.md) to push any
   images you have built locally and want to keep to the [Docker Hub registry](../docker-hub/index.md).
   Make sure to configure the [repository's visibility as "private"](../docker-hub/repos.md#private-repositories)
   for images that should not be publicly accessible. Refer to the [`docker push` page](../engine/reference/commandline/push.md)
   in the Docker Engine command line reference section for details on using this
   command.

   Alternatively, use [`docker image save -o images.tar image1 [image2 ...]`](../engine/reference/commandline/save.md)
   to save any images you want to keep to a local tar file. Refer to the
   [`docker image  save` page](../engine/reference/commandline/save.md) in the
   Docker Engine command line reference section for details on using this command.

3. If you use [named volume](../storage/index.md#more-details-about-mount-types)
   to store container data, such as databases, refer to the
   [backup, restore, or migrate data volumes](../storage/volumes.md#backup-restore-or-migrate-data-volumes)
   page in the storage section.

After backing up your data, you can uninstall the current version of Docker Desktop
and install a different version ([Windows](../docker-for-windows/install.md)
[macOS](mac/install.md), or reset Docker Desktop to factory defaults.

## Restore your data

1. Use [`docker pull`]((../engine/reference/commandline/load.md)) to restore images
   you pushed to Docker Hub in "step 2." in the [save your data section](#save-your-data)

   If you backed up your images to a local tar file, use [`docker image load -i images.tar`](../engine/reference/commandline/load.md)
   to restore previously saved images.

   Refer to the [`docker image load` page](../engine/reference/commandline/load.md)
   in the Docker Engine command line reference section for details on using this
   command.
2. Refer to the [backup, restore, or migrate data volumes](../storage/volumes.md#backup-restore-or-migrate-data-volumes)
   page in the storage section to restore volume data.
3. Re-create your containers if needed, using [`docker run`](../engine/reference/commandline/load.md),
   or [Docker Compose](../compose/index.md).
