---
title: Run Docker Hub images
keywords: get started, quick start, intro, concepts
description: Learn how to run Docker Hub images
---

You can share and store images in Docker Hub ([http://hub.docker.com](http://hub.docker.com){:target="_blank" rel="noopener" class="_"}). Docker hub has over 100,000 images created by developers that you can run locally. You can search for Docker Hub images and run them directly from Docker Desktop.

## Step 1: Search for the image

You can search for Docker Hub images on Docker Desktop. Use the search bar, or use the shortcut ⌘K on Mac or Ctrl + K on Windows. In the search bar, specify `welcome-to-docker` to find the `docker/welcome-to-docker` image used in this guide.

![Search Docker Desktop for the welcome-to-docker image](images/getting-started-search.png){:width="500px"}

## Step 2: Run the image

To run the `docker/welcome-to-docker` image, select **Run**. When the **Optional settings** appear, specify the **Host port** number `8090` and then select **Run**.

![Running the image in Docker Desktop](images/getting-started-run.gif){:width="500px"}

> **Note**
>
> Many images hosted on Docker Hub have a description that highlights what settings must be set in order to run them. You can read the description for the image on Docker Hub by selecting **View on Hub** or by searching for the image directly on [https://hub.docker.com](https://hub.docker.com){:target="_blank" rel="noopener" class="_"}.

## Step 3: Explore the container

That's it! The container is ready to use. Go to the **Containers** tab in Docker Desktop to view the container.

![Viewing the Containers tab in Docker Desktop](images/getting-started-view.png){:width="500px"}

## What's next

In this guide, you ran someone else's image from Docker Hub. Continue to the next guide to learn how you can publish your own image to Docker Hub.

[Publish your images](./publish-your-own-image.md){: .button .primary-btn}