---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore Volumes
---

The **Volumes** view in Docker Dashboard enables you to easily create and delete volumes and see which ones are being used. You can also see which container is using a specific volume and explore the files and folders in your volumes.

For more information about volumes, see [Use volumes](../../storage/volumes.md)

By default, the **Volumes** view displays a list of all the volumes. Volumes that are currently used by a container display the **In Use** badge.

## Manage your volumes

Use the **Search** field to search for any specific volume. 

You can sort volumes by:
- Name
- Date created
- Size

From the **More options** menu to the right of the search bar, you can choose what volume information to display.

## Inspect a volume

To explore the details of a specific volume, select a volume from the list. This opens the detailed view.

The **In Use** tab displays the name of the container using the volume, the image name, the port number used by the container, and the target. A target is a path inside a container that gives access to the files in the volume.

The **Data** tab displays the files and folders in the volume and the file size. To save a file or a folder, hover over the file or folder and click on the more options menu. Select **Save As** and then specify a location to download the file.

To delete a file or a folder from the volume, select **Remove** from the **More options** menu.

## Remove a volume

Removing a volume deletes the volume and all its data. 

To remove a volume, hover over the volume and then click the **Delete** icon. Alternatively, select the volume from the list and then click the **Delete** button.