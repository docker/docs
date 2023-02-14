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

### Image Hierarchy

The image you inspect may have one or more base images listed under **Image hierarchy**. This means the author of the image used another image as a starting point when building the image. Often these base images are either operating system images such as Debian, Ubuntu, and Alpine, or programming language images such as PHP, Python, and Java. 

A base image may have its own parent base image so there is a chain of base images represented in **Image hierarchy**. Selecting each image in the chain lets you see which layers originate from each base image. Selecting the **ALL** row reselects all the layers and base images for the entire image.

One or more of the base images may have updates available, which may include updated security patches that remove vulnerabilities from your image. Any base images with available updates are noted to the right of **Image hierarchy**.

### Layers

A Docker image consists of layers. Image layers are listed from top to bottom, with the earliest layer at the top and the most recent layer at the bottom. Often, the layers at the top of the list originate from a base image, and the layers towards the bottom are layers added by the image author, often by adding commands to a Dockerfile. To see which layers originate from a base image, simply select a base image under **Image hierarchy** and the relevant layers are highlighted.

Selecting individual or multiple layers filters the packages and vulnerabilities on the right-hand side to see what has been added by the selected layers.

### Vulnerabilities

Images may be exposed to vulnerabilities and exploits. These are detected and listed on the right-hand side, grouped by package, and sorted in order of severity. Further information on whether the vulnerability has an available fix, for example, can be examined by expanding the sections. For even more details, you can visit dso.docker.com by selecting the link. 

If you have any feedback about this feature you can use the **Give feedback** link to contact the development team.

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
3.. Select **Remove** to confirm.

## Interact with remote repositories

The **Images** view also allows you to manage and interact with images in remote repositories.
By default when you go to **Images** in Docker Desktop, you see a list of images that exist in your local image store.
Use the **Local** and **Hub** tabs near the top to toggle between viewing images in your local image store,
and images in remote Docker Hub repositories that you have access to.

You can also [connect JFrog Artifactory registries](#connect-an-artifactory-registry),
and browse images in JFrog repositories directly in Docker Desktop.

### Hub

Switching to the **Hub** tab prompts you to sign in to your Docker ID, if you're not already signed in.
When signed in, it shows you a list of images in Docker Hub organizations and repositories that you have access to.

Select an organization from the drop-down to view a list of repositories in that organization.

If you have enabled [Vulnerability Scanning](../../docker-hub/vulnerability-scanning.md) in Docker Hub, the scan results appear next to the image tags.

Hovering over an image tag reveals two options:

- **Pull**: pulls the latest version of the image from Docker Hub.
- **View in Hub**: opens the Docker Hub page and displays detailed information about the image.

### Artifactory

Integrating Artifactory lets you interact with images in Artifactory, directly in the **Images** view.
The integration described here connects your local Docker Desktop client with Artifactory,
similar to running the `docker login` command.
You can browse, filter, save, and pull images in the Artifactory instance you configure.

You may also want to consider enabling automatic image analysis for your Artifactory repositories.
Learn more about [Artifactory integration with Docker Scout]().

#### Connect an Artifactory registry

To connect a new Artifactory registry to Docker Desktop:

1. In the **Images** view, select the plus ("+") button in the tab row near the top.
2. If you aren't already signed in to an Artifactory registry, you're prompted to enter the sign-in details for your Artifactory instance:
   
   - **Host**: the base hostname for your Artifactory instance.
   - **Username**: your Artifactory username.
   - **Password**: your Artifactory password, or API token.

   This is the same as using the `docker login` command.

3. When signed in, a new **Artifactory** tab appears and opens, listing the most recent images in your Artifactory instance.

By default, the image list sorts images by most recent push, where the newest images appear higher in the list.
