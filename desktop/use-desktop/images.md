---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore Images
---

The **Images**  view is a simple interface that lets you manage Docker images without having to use the CLI. By default, it displays a list of all Docker images on your local disk. 

You can also view images in remote repositories, once you have signed in to Docker Hub. This allows you to collaborate with your team and manage your images directly through Docker Desktop.

The **Images** view allows you to perform core operations such as running an image as a container, pulling the latest version of an image from Docker Hub, pushing the image to Docker Hub, and inspecting images.

The **Images** view displays metadata about the image such as the:
- Tag
- Image ID
- Date created
- Size of the image.

It also displays **In Use** tags next to images used by running and stopped containers.

The **Images on disk** status bar displays the number of images and the total disk space used by the images.

## Manage your images

Use the **Search** field to search for any specific image.

You can sort images by:

- Name
- Date created
- Size


## Run an image as a container

From the **Images view**, hover over an image and click **Run**.

When prompted you can either:

- Click the **Optional settings** drop-down to specify a name, port, volumes, environment variables and click **Run**
- Click **Run** without specifying any optional settings.

## Inspect an image

Inspecting an image displays detailed information about the image such as the:

- Image history
- Image ID
- Date the image was created
- Size of the image

To inspect an image, hover over an image, select the **More options** button and then select **Inspect** from the dropdown menu. 

## Pull the latest image from Docker Hub

Select the image from the list, click the **More options** button and click **Pull**.

> **Note**
>
> The repository must exist on Docker Hub in order to pull the latest version of an image. You must be logged in to pull private images.

## Push an image to Docker Hub

Select the image from the list, click the **More options** button and click **Push to Hub**.

> **Note**
>
> You can only push an image to Docker Hub if the image belongs to your Docker ID or your organization. That is, the image must contain the correct username/organization in its tag to be able to push it to Docker Hub.

## Remove an image

> **Note**
>
> To remove an image used by a running or a stopped container, you must first remove the associated container.

You can remove individual images or use the **Clean up** option to delete unused and dangling images.

An unused image is an image which is not used by any running or stopped containers. An image becomes dangling when you build a new version of the image with the same tag.

To remove individual images, select the image from the list, click the **More options** button and click **Remove**

To remove an unused or a dangling image:

1. Select the **Clean up** option from the **Images on disk** status bar.
2. Use the **Unused** or **Dangling** check boxes to select the type of images you would like to remove.

    The **Clean up** images status bar displays the total space you can reclaim by removing the selected images.
3.. Select **Remove** to confirm.

## Interact with remote repositories

The **Images** view also allows you to manage and interact with images in remote repositories and lets you switch between organizations. Select an organization from the drop-down to view a list of repositories in your organization.

> **Note**
>
> If you have a paid Docker subscription and enabled [Vulnerability Scanning](../../docker-hub/vulnerability-scanning.md) in Docker Hub, the scan results appear on the Remote repositories tab.
The **Pull** option allows you to pull the latest version of the image from Docker Hub. The **View in Hub** option opens the Docker Hub page and displays detailed information about the image, such as the OS architecture, size of the image, the date when the image was pushed, and a list of the image layers.

To interact with remote repositories:

1. Click the **Remote repositories** tab.
2. Select an organization from the drop-down list. This displays a list of repositories in your organization.
3. Hover over an image from the list and then select **Pull** to pull the latest image from the remote repository.

To view a detailed information about the image in Docker Hub, select the image and then click **View in Hub**.

The **View in Hub** option opens the Docker Hub page and displays detailed information about the image, such as the OS architecture, size of the image, the date when the image was pushed, and a list of the image layers.

If you have a paid Docker subscription and have enabled [Vulnerability Scanning](../../docker-hub/vulnerability-scanning.md) the Docker Hub page also displays a summary of the vulnerability scan report and provides detailed information about the vulnerabilities identified.
