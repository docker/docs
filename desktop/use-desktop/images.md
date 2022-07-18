---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore Images
---

The **Images**  view is a simple interface that lets you manage Docker images without having to use the CLI. By default, it displays a list of all Docker images on your local disk. To view images in remote repositories, click **Sign in** and connect to Docker Hub. This allows you to collaborate with your team and manage your images directly through Docker Desktop.

The Images view allows you to perform core operations such as running an image as a container, pulling the latest version of an image from Docker Hub, pushing the image to Docker Hub, and inspecting images.

In addition, the Images view displays metadata about the image such as the tag, image ID, date when the image was created, and the size of the image. It also displays **In Use** tags next to images used by running and stopped containers. This allows you to review the list of images and use the **Clean up images** option to remove any unwanted images from the disk to reclaim space.

The Images view also allows you to search images on your local disk and sort them using various options.

Let's explore the various options in the **Images** view.

If you don’t have any images on your disk, run the command `docker pull redis` in a terminal to pull the latest Redis image. This command pulls the latest Redis image from Docker Hub.

Select **Dashboard** > **Images** to see the Redis image.

### Run an image as a container

Now that you have a Redis image on your disk, let’s run this image as a container:

1. From Docker Dashboard, select **Images**. A list of images on your local disk displays.
2. Select the Redis image from the list and click **Run**.
3. Optional step: When prompted, click the **Optional settings** drop-down to specify a name, port, volumes, environment variables and click **Run**.

    To use the defaults, click **Run** without specifying any optional settings. This creates a new container from the Redis image and opens it on the **Containers** view.

### Pull the latest image from Docker Hub

To pull the latest image from Docker Hub:

1. From the Docker menu, select **Dashboard** > **Images**. This displays a list of images on your local disk.
2. Select the image from the list and click the more options button.
3. Click **Pull**. This pulls the latest version of the image from Docker Hub.

> **Note**
>
> The repository must exist on Docker Hub in order to pull the latest version of an image. You must be logged in to pull private images.
### Push an image to Docker Hub

To push an image to Docker Hub:

1. From the Docker menu, select **Dashboard** > **Images**. This displays a list of images on your local disk.
2. Select the image from the list and click the more options button.
3. Click **Push to Hub.**

> **Note**
>
> You can only push an image to Docker Hub if the image belongs to your Docker ID or your organization. That is, the image must contain the correct username/organization in its tag to be able to push it to Docker Hub.
### Inspect an image

Inspecting an image displays detailed information about the image such as the image history, image ID, the date the image was created, size of the image, etc. To inspect an image:

1. From the Docker menu, select **Dashboard** > **Images**. This displays a list of images on your local disk.
2. Select the image from the list and click the more options button.
3. Click **Inspect**.
4. The image inspect view also provides options to pull the latest image, push image to Hub, remove the image, or run the image as a container.

### Remove an image

The **Images** view allows you to remove unwanted images from the disk. The Images on disk status bar displays the number of images and the total disk space used by the images.

You can remove individual images or use the **Clean up** option to delete unused and dangling images.

To remove individual images:

1. From the Docker menu, select **Dashboard** > **Images**. This displays a list of images on your local disk.
2. Select the image from the list and click the more options button.
3. Click **Remove**. This removes the image from your disk.

> **Note**
>
> To remove an image used by a running or a stopped container, you must first remove the associated container.
**To remove unused and dangling images:**

An **unused** image is an image which is not used by any running or stopped containers. An image becomes **dangling** when you build a new version of the image with the same tag.

**To remove an unused or a dangling image:**

1. From the Docker menu, select **Dashboard** > **Images**. This displays a list of images on your disk.
2. Select the **Clean up** option from the **Images on disk** status bar.
3. Use the **Unused** and **Dangling** check boxes to select the type of images you would like to remove.

    The **Clean up** images status bar displays the total space you can reclaim by removing the selected images.
4. Click **Remove** to confirm.

### Interact with remote repositories

The Images view also allows you to manage and interact with images in remote repositories and lets you switch between organizations. Select an organization from the drop-down to view a list of repositories in your organization.

> **Note**
>
> If you have a paid Docker subscription and enabled [Vulnerability Scanning](../../docker-hub/vulnerability-scanning.md) in Docker Hub, the scan results will appear on the Remote repositories tab.
The **Pull** option allows you to pull the latest version of the image from Docker Hub. The **View in Hub** option opens the Docker Hub page and displays detailed information about the image, such as the OS architecture, size of the image, the date when the image was pushed, and a list of the image layers.

To interact with remote repositories:

1. Click the **Remote repositories** tab.
2. Select an organization from the drop-down list. This displays a list of repositories in your organization.
3. Click on an image from the list and then select **Pull** to pull the latest image from the remote repository.
4. To view a detailed information about the image in Docker Hub, select the image and then click **View in Hub**.

    The **View in Hub** option opens the Docker Hub page and displays detailed information about the image, such as the OS architecture, size of the image, the date when the image was pushed, and a list of the image layers.

    If you have a paid Docker subscription and have enabled [Vulnerability Scanning](../../docker-hub/vulnerability-scanning.md) the Docker Hub page also displays a summary of the vulnerability scan report and provides detailed information about the vulnerabilities identified.