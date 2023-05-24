---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore Images
---

The **Images**  view is a simple interface that lets you manage Docker images without having to use the CLI. By default, it displays a list of all Docker images on your local disk. 

You can also view Hub images once you have signed in to Docker Hub. This allows you to collaborate with your team and manage your images directly through Docker Desktop.

The **Images** view allows you to perform core operations such as running an image as a container, pulling the latest version of an image from Docker Hub, pushing the image to Docker Hub, and inspecting images.

The **Images** view displays metadata about the image such as the:
- Tag
- Image ID
- Date created
- Size of the image.

It also displays **In Use** tags next to images used by running and stopped containers. You can choose what information you want displayed by selecting the **More options** menu to the right of the search bar, and then use the toggle switches according to your preferences. 

The **Images on disk** status bar displays the number of images and the total disk space used by the images and when this information was last refreshed.

## Manage your images

Use the **Search** field to search for any specific image.

You can sort images by:

- In use
- Unused
- Dangling

## Run an image as a container

From the **Images view**, hover over an image and select **Run**.

When prompted you can either:

- Select the **Optional settings** drop-down to specify a name, port, volumes, environment variables and select **Run**
- Select **Run** without specifying any optional settings.

## Inspect an image

To inspect an image, simply select the image row. Inspecting an image displays detailed information about the image such as the:

- Image history
- Image ID
- Date the image was created
- Size of the image
- Layers making up the image
- Base images used
- Vulnerabilities found
- Packages inside the image

The image view is powered by [Docker Scout](../../scout/index.md).
For more information about this view, see [Image details view](../../scout/image-details-view.md)

## Pull the latest image from Docker Hub

Select the image from the list, select the **More options** button and select **Pull**.

> **Note**
>
> The repository must exist on Docker Hub in order to pull the latest version of an image. You must be logged in to pull private images.

## Push an image to Docker Hub

Select the image from the list, select the **More options** button and select **Push to Hub**.

> **Note**
>
> You can only push an image to Docker Hub if the image belongs to your Docker ID or your organization. That is, the image must contain the correct username/organization in its tag to be able to push it to Docker Hub.

## Remove an image

> **Note**
>
> To remove an image used by a running or a stopped container, you must first remove the associated container.

You can remove individual images or use the **Clean up** option to delete unused and dangling images.

An unused image is an image which is not used by any running or stopped containers. An image becomes dangling when you build a new version of the image with the same tag.

To remove individual images, select the image from the list, select the **More options** button and select **Remove**

To remove an unused or a dangling image:

1. Select the **Clean up** option from the **Images on disk** status bar.
2. Use the **Unused** or **Dangling** check boxes to select the type of images you would like to remove.

    The **Clean up** images status bar displays the total space you can reclaim by removing the selected images.
3. Select **Remove** to confirm.

## Interact with remote repositories

The **Images** view also allows you to manage and interact with images in remote repositories.
By default, when you go to **Images** in Docker Desktop, you see a list of images that exist in your local image store.
The **Local** and **Hub** tabs near the top toggles between viewing images in your local image store,
and images in remote Docker Hub repositories that you have access to.

You can also [connect JFrog Artifactory registries](#connect-an-artifactory-registry),
and browse images in JFrog repositories directly in Docker Desktop.

### Hub

Switching to the **Hub** tab prompts you to sign in to your Docker ID, if you're not already signed in.
When signed in, it shows you a list of images in Docker Hub organizations and repositories that you have access to.

Select an organization from the drop-down to view a list of repositories for that organization.

If you have enabled [Vulnerability Scanning](../../docker-hub/vulnerability-scanning.md) in Docker Hub, the scan results appear next to the image tags.

Hovering over an image tag reveals two options:

- **Pull**: pulls the latest version of the image from Docker Hub.
- **View in Hub**: opens the Docker Hub page and displays detailed information about the image.

### Artifactory

The Artifactory integration lets you interact with images in JFrog Artifactory,
and JFrog container registry, directly in the **Images** view of Docker Desktop.
The integration described here connects your local Docker Desktop client with Artifactory.
You can browse, filter, save, and pull images in the Artifactory instance you configure.

You may also want to consider activating automatic image analysis for your Artifactory repositories.
Learn more about [Artifactory integration with Docker Scout](../../scout/artifactory.md).

#### Connect an Artifactory registry

To connect a new Artifactory registry to Docker Desktop:

1. Sign in to an Artifactory registry using the `docker login` command:

   ```console
   $ cat ./password.txt | docker login -u <username> --password-stdin <hostname>
   ```

   - `password.txt`: text file containing your Artifactory password.
   - `username`: your Artifactory username.
   - `hostname`: hostname for your Artifactory instance.

2. Open the **Images** view in Docker Desktop.
3. Select the **Artifactory** tab near the top of the image view to see Artifactory images.

When signed in, a new **Artifactory** tab appears in the **Images** view.
By default, the image list shows images sorted by push date: the newest images appear higher in the list.

