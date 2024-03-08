---
title: What is a container?
keywords: concepts, build, images, container, docker desktop
description: What is a container? This concept page will teach you about containers and provide a quick hands-on where you will run your first container.
---

<iframe width="650" height="365" src="https://www.youtube.com/embed/nsWWQ1xoEy0?rel=0" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Explanation

A `container` is an isolated sandbox environment on your machine. Imagine containers like shipping containers for cargo ships. Each `container` holds its own cargo (application code, libraries, configurations) and is sealed off from other containers. This ensures that even if one container gets damaged or malfunctions, it doesn't affect the contents or functionality of other containers on the same ship (your machine).

While containers are isolated from each other, they share the underlying operating system kernel. This means they are much lighter and more efficient compared to virtual machines (VMs) which require their own complete operating system. This allows you to run more containers on a single machine compared to VMs.

Containers are self-contained and include everything an application needs to run. This makes them highly portable. You can move a container from one machine to another (e.g., from development to production) without worrying about compatibility issues as long as the underlying system supports containerization.

Consider a web application with separate components: frontend, backend, and database. Each component can be packaged in its own container. They can all run on the same machine, but they are isolated from each other. This simplifies development, testing, and deployment as you manage each component independently.

## Try it now

In this hands-on, you will see how to run a simple `Nginx` container using Docker Desktop.

### Run a container

Use the following instructions to run a container.

1. Open Docker Desktop and select the search.
2. Specify `nginx` in the search and then select **Run**.

    ![Screenshot of the global search menu with nginx typed in as the search term](images/nginx.webp)

3. Expand the **Optional settings**.
4. In **Container name**, specify `nginx`.
5. In **Host port**, specify `80`.

   ![Specifying host port 80](images/nginx-parameters.webp?border=true&w=350&h=400)

6. Click **Run**.

###  View your container

You just ran a container! You can view it in the **Containers** tab of the Docker Desktop GUI.

![nginx container running](images/nginx-running.webp?border=true)


This container runs a simple web server that displays a simple website.
When working with more complex projects, you'll run different parts in different
containers. For example, a different container for the frontend, backend, and
database. In this walkthrough, you only have a simple frontend container.

###  Access the frontend

The frontend is accessible on port 80 of your local host. Select the link in
the **Port(s)** column of your container, or visit
[http://localhost:80](http://localhost:80) in your browser to view it.

![Accessing container frontend from Docker Desktop](images/nginx-success.webp?border=true)

### Explore your container

Docker Desktop lets you easily view and interact with different aspects of your
container. Try it out yourself. Select your container and then select **Files**
to explore your container's isolated file system.

![Viewing container details in Docker Desktop](images/nginx-files.webp?border=true&w=600&h=403)

### Stop your container

The `nginx` container continues to run until you stop it. To stop
the container in Docker Desktop, go to the **Containers** tab and select the
**Stop** icon in the **Actions** column of your container.

![Stopping a container in Docker Desktop](images/nginx-stopped.webp?border=true)

In this walkthrough, you ran a pre-made image and explored a container. In addition to running pre-made images, you can build and run your own application as container.

## Additional resources

- [What is a container](https://docs.docker.com/guides/walkthroughs/what-is-a-container/)
- [Overview of container](https://www.docker.com/resources/what-container/)
- [Why Docker?](https://www.docker.com/why-docker/)

{{< button text="What is an image" url="what-is-an-image" >}}
