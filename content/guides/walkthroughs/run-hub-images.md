---
title: Run Docker Hub images
keywords: get started, quick start, intro, concepts
description: Learn how to run Docker Hub images
aliases:
- /get-started/run-docker-hub-images/
---

You can share and store images in Docker Hub
([http://hub.docker.com](http://hub.docker.com)). Docker Hub has over 100,000
images created by developers that you can run locally. You can search for Docker
Hub images and run them directly from Docker Desktop.

Before you start, [get Docker Desktop](../../get-docker.md).

## Step 1: Search for the image

You can search for Docker Hub images on Docker Desktop. To search for the image used in this walkthrough, do the following:

1. Open Docker Desktop and select the search.
2. Specify `docker/welcome-to-docker` in the search.

![Search Docker Desktop for the welcome-to-docker image](images/getting-started-search.png?w=400)

## Step 2: Run the image

To run the `docker/welcome-to-docker` image, do the following:

1. After finding the image using search, select **Run**.
2. Expand the **Optional settings**.
3. In **Host port**, specify `8090`.
4. Select **Run**.

![Running the image in Docker Desktop](images/getting-started-run.gif?w=400&border=true)

> **Note**
>
> Many images hosted on Docker Hub have a description that highlights what
> settings must be set in order to run them. You can read the description for
> the image on Docker Hub by selecting the image name in the search or by
> searching for the image directly on
> [https://hub.docker.com](https://hub.docker.com).

## Step 3: Explore the container

That's it! The container is ready to use. Go to the **Containers** tab in Docker Desktop to view the container.

![Viewing the Containers tab in Docker Desktop](images/getting-started-view.png?w=400)

## Summary

In this walkthrough, you searched for an image on Docker Hub and ran it as a container. Docker Hub has over 100,000 more images that you can use to help build your own application.

Related information:

- Deep dive into the [Docker Hub manual](../../docker-hub/_index.md)
- Explore more images on [Docker Hub](https://hub.docker.com)

## Next steps

Continue to the next walkthrough to learn how you can use Docker to run
multi-container applications.

{{< button url="./multi-container-apps.md" text="Run multi-container apps" >}}