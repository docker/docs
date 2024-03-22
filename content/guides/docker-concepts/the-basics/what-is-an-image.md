---
title: What is an image
keywords: concepts, build, images, container, docker desktop
description: What is an image
---

<iframe width="650" height="365" src="https://www.youtube.com/embed/nsWWQ1xoEy0?rel=0" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Explanation

Remember those magic boxes we introduced for each component of your web application (`React` frontend, `Python` API backend, `PostgreSQL` database)? What makes them work? Well, those boxes are powered by container images. You can think of a container image as a self-contained package. It’s a complete bundle that includes everything your container needs to run perfectly - files and programs, configuration settings, libraries, binaries and all the dependencies. 

So if you think of a `container` as an isolated sandbox environment using all the files in that environment without even installing it on the host, then the question is how do I get those files into that environment? Where do those files come from? How do we share those environments? That’s where container images come in. A container image is a standardized package that includes all of the files, binaries, libraries, and configurations to run a container. 

Two important aspects of images:



1.  **Images are immutable**

Consider a PostgreSQL Docker image. You can’t edit it directly. They are unchanging once created. Instead, you can create a new image with the desired changes.



2. **Container images are composed of Layers**

Container images are built using layers, and layers represent file system changes. Each layer adds, removes, or modifies layers that might have been there previously. For example, if you have a `Python` Docker Image then you can use the image and add additional layers on the top of the Python image and run it in the form of a `container` without even managing Python itself. If you need to make a change to a Python container image, you either need to build a new image or extend an existing base image, adding your layers on top of it.

So how are these images distributed? The answer is `Docker Hub`.

[Docker Hub](https://hub.docker.com) is the default global marketplace for storing and distributing images. It has over 100,000 images created by developers that you can run locally. You can search for Docker Hub images and run them directly from Docker Desktop.

Docker officially maintains Docker Hub. You can find a variety of [DOI](https://docs.docker.com/trusted-content/official-images/)(Docker Official Images), [DVP](https://docs.docker.com/trusted-content/dvp-program/) (Docker Verified Publisher), and [DSOS](https://docs.docker.com/trusted-content/dsos-program/) (Docker Sponsored Open-source) images that provide a good starting point to start with the images in the Docker Hub. 

For example, `Redis` and `Memcached` are a few popular ready-to-go DOI Docker images. You can download these images and have these services up and running pretty quickly. There are also base images, like the `Node.js` Docker image, that you can use as a starting point and add your own files and configurations. Whenever you download a Docker image to start a Docker container, you’re downloading everything that’s needed to run a Docker container.


## Try it now

{{< tabs >}}
{{< tab name="Using Docker Desktop" >}}

## Search and pull a container image

1. Open the Docker Desktop dashboard and select the **Images** tab in the left-hand navigation menu.

![A screenshot of the Docker Desktop dashboard showing the image tab on the left sidebar](images/click-image.webp?border=true&w=1050&h=400)

2. Select the **Search images to run** button. If you don't see it, select the `global search bar` at the top of the screen.

![A screenshot of the Docker Desktop dashboard showing the search ta](images/search-image.webp?border)

3. In the **Search** field, enter "welcome-to-docker". Once the search has completed, select the `docker/welcome-to-docker` image.

 ![A screenshot of the Docker Desktop dashboard showing the search results for the docker/welcome-to-docker image](images/select-image.webp?border=true&w=1050&h=400)


Select **Pull** to download the image to your local system.


## View your image

Once you have an image downloaded, you can view quite a few details about the image either through the GUI or the CLI.

1. In the Docker Desktop Dashboard, select the **Images** tab.

2. Select the **docker/welcome-to-docker** image to open details about the image.

![A screenshot of the Docker Desktop dashboard showing the images tab with an arrow pointing to the docker/welcome-to-docker image](images/pulled-image.webp?border=true&w=1050&h=400)


3. The image details page presents you with information regarding the layers of the image, the packages and libraries installed in the image, and any discovered vulnerabilities.

![A screenshot of the image details view for the docker/welcome-to-docker image](images/image-layers.webp?border=true&w=1050&h=400)


In this walkthrough, you searched and pulled a Docker image. In addition to pulling a Docker image, you also learned about the layers of a Docker Image.

{{< /tab >}}

{{< tab name="CLI" >}}

Follow the instructions to search and pull a Docker image using CLI to view its layers.

1. [Download and install](https://www.docker.com/products/docker-desktop/) Docker Desktop.

2. Open a terminal and run the following command:

```console
 docker search docker/welcome-to-docker
```

You will see output like the following:

```console
 docker search docker/welcome-to-docker
 NAME                       DESCRIPTION                                     STARS     OFFICIAL
 docker/welcome-to-docker   Docker image for new users getting started w…   20
```

This output is the result of running the docker search command. It shows you information about the image available on the Docker Hub.


3. Pull the image

```console
 docker pull docker/welcome-to-docker
```

You will see output like the following:

```console
 Using default tag: latest
 latest: Pulling from docker/welcome-to-docker
 579b34f0a95b: Download complete
 d11a451e6399: Download complete
 1c2214f9937c: Download complete
 b42a2f288f4d: Download complete
 54b19e12c655: Download complete
 1fb28e078240: Download complete
 94be7e780731: Download complete
 89578ce72c35: Download complete
 Digest: sha256:eedaff45e3c78538087bdd9dc7afafac7e110061bbdd836af4104b10f10ab693
 Status: Downloaded newer image for docker/welcome-to-docker:latest
 docker.io/docker/welcome-to-docker:latest
```

Each of these lines represents different layers of the image being downloaded and completed successfully. Each layer contributes to the overall functionality of the image. The image is being pulled from the official Docker repository named `docker` and the image name is `welcome-to-docker`.

4. View your image

```console
 docker images
 REPOSITORY                 TAG       IMAGE ID       CREATED        SIZE
 docker/welcome-to-docker   latest    eedaff45e3c7   4 months ago   29.7MB
```

The command shows a list of Docker images currently available on your system. The `docker/welcome-to-docker` has a total size of approximately 29.7MB.


In this walkthrough, you searched and pulled a Docker image. In addition to pulling a Docker image, you also learned about the layers of a Docker Image.

{{< /tab >}}
{{< /tabs >}}


## Additional resources

* [Explore the Image view in Docker Desktop](https://docs.docker.com/desktop/use-desktop/images/)
* [Packaging your software](https://docs.docker.com/build/building/packaging/)
* [Explore the Containers view in Docker Desktop](https://docs.docker.com/desktop/use-desktop/container/)
* [Overview of Docker Hub](https://hub.docker.com)

{{< button text="What is a registry?" url="what-is-a-registry" >}}
