---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore Volumes
---
## Explore volumes

You can use [volumes](../../storage/volumes.md) to store files and share them among containers. Volumes are created and are directly managed by Docker. They are also the preferred mechanism to persist data in Docker containers and services.

> **Volume Management is now available for all subscriptions**
>
> Starting with Docker Desktop 4.1.0 release, Volume management is available for users on any subscription, including Docker Personal. Update Docker Desktop to 4.1.0 to start managing your volumes for free.
{: .important}

The **Volumes** view in Docker Dashboard enables you to easily create and delete volumes and see which ones are being used. You can also see which container is using a specific volume and explore the files and folders in your volumes.

### Manage volumes

By default, the **Volumes** view displays a list of all the volumes. Volumes that are currently used by a container display the **In Use** badge.

Use the **Search** field to search for any specific volumes. You can also sort volumes by the name, the date created, and the size of the volume.

To explore the details of a specific volume, select a volume from the list. This opens the detailed view.

The **In Use** tab displays the name of the container using the volume, the image name, the port number used by the container, and the target. A target is a path inside a container that gives access to the files in the volume.

The **Data** tab displays the files and folders in the volume and their file size. To save a file or a folder, hover over the file or folder and click on the more options menu. Select **Save As** and then specify a location to download the file.

To delete a file or a folder from the volume, select **Remove** from the more options menu.

### Remove a volume

Removing a volume deletes the volume and all its data. To remove a volume, hover over the volume and then click the **Delete** button. Alternatively, select the volume from the list and then click the **Delete** button.