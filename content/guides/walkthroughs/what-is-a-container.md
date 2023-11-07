---
title: What is a container?
keywords: get started, quick start, intro, concepts
description: Learn what a container is by seeing and inspecting a running container.
aliases:
- /get-started/what-is-a-container/
---

A container is an isolated environment for your code. This means that a
container has no knowledge of your operating system, or your files. It runs on
the environment provided to you by Docker Desktop. Containers have everything
that your code needs in order to run, down to a base operating system. You can
use Docker Desktop to manage and explore your containers.

In this walkthrough, you'll view and explore an actual container in Docker
Desktop.

Before you start, [get Docker Desktop](../../get-docker.md).

## Step 1: Set up the walkthrough

The first thing you need is a running container. For this guide, use the
pre-made `welcome-to-docker` image. An image is like a blueprint for a
container. Do the following to set up the walkthrough.

1. Open Docker Desktop and select the search.
2. Specify `docker/welcome-to-docker` in the search and then select **Run**.
3. Expand the **Optional settings**.
4. In **Container name**, specify `welcome-to-docker`.
5. In **Host port**, specify `8088`.
6. Select **Run**.

![Set up the walkthrough](images/getting-started-setup.gif?w=400&border=true)

## Step 2: View containers on Docker Desktop

You just ran a container! You can view it in the **Containers** tab of Docker
Desktop. This container runs a simple web server that displays a simple website.
When working with more complex projects, you'll run different parts in different
containers. For example, a different container for the frontend, backend, and
database. In this walkthrough, you only have a simple frontend container.

![Docker Desktop with get-started container running](images/getting-started-container.png?w=400)

## Step 3: View the frontend

The frontend is accessible on port 8088 of your local host. Select the link in
the **Port(s)** column of your container, or visit
[http://localhost:8088](http://localhost:8088) in your browser to view it.

![Accessing container frontend from Docker Desktop](images/getting-started-frontend.png?w=300&border=true)

## Step 4: Explore your container

Docker Desktop lets you easily view and interact with different aspects of your
container. Try it out yourself. Select your container and then select **Files**
to explore your container's isolated file system.

![Viewing container details in Docker Desktop](images/getting-started-explore-container.png?w=300&border=true)

## Step 5: Stop your container

The `welcome-to-docker` container continues to run until you stop it. To stop
the container in Docker Desktop, go to the **Containers** tab and select the
**Stop** icon in the **Actions** column of your container.

![Stopping a container in Docker Desktop](images/getting-started-stop.png?w=400)

## Summary

In this walkthrough, you ran a pre-made image and explored a container. In addition to running pre-made images, you can build and run your own application as container.

Related information:

- Read more about containers in [Use containers to Build, Share and Run your applications](https://www.docker.com/resources/what-container/)
- Deep dive in Liz Rice's [Containers from Scratch](https://www.youtube.com/watch?v=8fi7uSYlOdc&t=1s) video presentation

## Next steps

Continue to the next walkthrough to learn what you need to create your own image
and run it as container.

{{< button url="./run-a-container.md" text="How do I run a container?" >}}