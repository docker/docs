---
title: How do I run a container?
keywords: get started, quick start, intro, concepts
description: Learn how to build your own image and run it as a container
aliases:
- /get-started/run-your-own-container/
---

In this walkthrough, you'll learn the basic steps of building an image and running your own container. This walkthrough uses a sample Node.js application, but it's not necessary to know Node.js.

![Running an image in Docker Desktop](images/getting-started-run-intro.png?w=400)

Before you start, [get Docker Desktop](../../get-docker.md).

## Step 1: Get the sample application

If you have git, you can clone the repository for the sample application. Otherwise, you can download the sample application. Choose one of the following options.

{{< tabs >}}
{{< tab name="Clone with git" >}}

Use the following command in a terminal to clone the sample application repository.

```console
$ git clone https://github.com/docker/welcome-to-docker
```

{{< /tab >}}
{{< tab name="Download" >}}

Download the source and extract it.

{{< button url="https://github.com/docker/welcome-to-docker/archive/refs/heads/main.zip" text="Download the source" >}}

{{< /tab >}}
{{< /tabs >}}

## Step 2: View the Dockerfile in your project folder

To run your code in a container, the most fundamental thing you need is a
Dockerfile. A Dockerfile describes what goes into a container. This sample already contains a `Dockerfile`. For your own projects, you'll need to create your own `Dockerfile`. You can open the `Dockerfile` in a code or text editor and explore its contents.

![Viewing Dockefile contents](images/getting-started-dockerfile.png?w=400)

## Step 3: Build your first image

You always need an image to run a container. In a terminal, run the following commands to build the image. Replace `/path/to/welcome-to-docker/` with the path to your `welcome-to-docker` directory.

{{< include "open-terminal.md" >}}

```console
$ cd /path/to/welcome-to-docker/
```
```console
$ docker build -t welcome-to-docker .
```

Building the image may take some time. After your image is built, you can view your image in the **Images** tab in Docker Desktop.

## Step 4: Run your container

To run your image as a container, do the following:

1. In Docker Desktop, go to the **Images** tab.
2. Next to your image, select **Run**.
3. Expand the **Optional settings**.
4. In **Host port**, specify `8089`.
5. Select **Run**.

![Running an image in Docker Desktop](images/getting-started-run-image.gif?w=400&border=true)

## Step 5: View the frontend

You can use Docker Desktop to access your running container. Select the link next to your container in Docker Desktop or go to [http://localhost:8089](http://localhost:8089) to view the frontend.

![Selecting the container link](images/getting-started-frontend-2.png?w=300&border=true)

## Summary

In this walkthrough, you built your own image and ran it as a container. In addition to building and running your own images, you can run images from Docker Hub.

Related information:

- Deep dive into building images in the [Build with Docker guide](../../build/guide/_index.md)

## Next steps

Continue to the next walkthrough to learn how you can run one of over 100,000 pre-made images from Docker Hub.

{{< button url="./run-hub-images.md" text="Run Docker Hub images" >}}