---
title: Build your Python image
keywords: python, build, images, dockerfile
description: Learn how to build an image of your Python application
---

## Prerequisites

* You have installed the latest version of [Docker Desktop](../../get-docker.md).
* You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

## Overview

This guide walks you through building your first Python image. An image
includes everything needed to run an application - the code or binary, runtime,
dependencies, and any other file system objects required.

## Sample application

The sample application uses the popular [Flask](https://flask.palletsprojects.com/) framework.

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/docker/python-docker
```

## Test the application without Docker (optional)

You can test the application locally without Docker before you continue building and running the application with Docker. This section requires you to have Python 3.11 or later installed on your machine. Download and install [Python](https://www.python.org/downloads/).

Open your terminal and navigate to the working directory you created. Create an environment, install the dependencies, and start the application to make sure it’s running.

```console
$ cd /path/to/python-docker
$ python3 -m venv .venv
$ source .venv/Scripts/activate
(.venv) $ python3 -m pip install -r requirements.txt
(.venv) $ python3 -m flask run
```

To test that the application is working, open a new browser and navigate to `http://localhost:5000`.

Switch back to the terminal where the server is running and you should see the following requests in the server logs. The data and timestamp will be different on your machine.

```shell
127.0.0.1 - - [22/Sep/2020 11:07:41] "GET / HTTP/1.1" 200 -
```

## Create a Dockerfile for Python

Now that you have an application, you can use `docker init` to create a Dockerfile for it. Inside the `python-docker` directory, run the `docker init` command. Refer to the following example to answer the prompts from `docker init`.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml

Let's get started!

? What application platform does your project use? Python
? What version of Python do you want to use? 3.11.4
? What port do you want your app to listen on? 5000
? What is the command to run your app? python3 -m flask run --host=0.0.0.0
```

You should now have the following 3 new files in your `python-docker`
directory:
 - Dockerfile
 - `.dockerignore`
 - `compose.yaml`

The Dockerfile is used to build the image. Open the Dockerfile
in your favorite IDE or text editor and see what it contains. To learn more
about Dockerfiles, see the [Dockerfile reference](../../engine/reference/builder.md).

## .dockerignore file

When you run `docker init`, it also creates a [`.dockerignore`](../../engine/reference/builder.md#dockerignore-file) file. Use the `.dockerignore` file to specify patterns and paths that you don't want copied into the image in order to keep the image as small as possible. Open up the `.dockerignore` file in your favorite IDE or text editor and see what's inside already.

## Build an image

Now that you’ve created the Dockerfile, you can build the image. To do this, use the `docker build` command. The `docker build` command builds Docker images from a Dockerfile and a build context. A build context is the set of files that the build has access to.

The build command optionally takes a `--tag` flag. The tag sets the name of the image and an optional tag in the format `name:tag`. Leave off the optional `tag` for now to help simplify things. If you don't pass a tag, Docker uses “latest” as its default tag.

Build the Docker image.

```console
$ docker compose up --build
```

You should get output similar to the following.

```console
[+] Building 1.3s (12/12) FINISHED
 => [internal] load .dockerignore
 => => transferring context: 680B
 => [internal] load build definition from Dockerfile
 => => transferring dockerfile: 1.59kB
 => resolve image config for docker.io/docker/dockerfile:1
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:39b85bbfa7536a5feceb7372a0817649ecb2724562a38360f4
 => [internal] load metadata for docker.io/library/python:3.11.4-slim
 => [base 1/5] FROM docker.io/library/python:3.11.4-slim@sha256:36b544be6e796eb5caa0bf1ab75a17d2e20211cad7f66f04f
 => [internal] load build context
 => => transferring context: 63B
 => CACHED [base 2/5] WORKDIR /app
 => CACHED [base 3/5] RUN adduser     --disabled-password     --gecos ""     --home "/nonexistent"     --shell
 => CACHED [base 4/5] RUN --mount=type=cache,target=/root/.cache/pip     --mount=type=bind,source=requirements.tx
 => CACHED [base 5/5] COPY . .
 => exporting to image
 => => exporting layers
 => => writing image sha256:37f9294069a95e5b34bb9e9035c6ea6ad16657818207c9d0dc73594f70144ce4
 => => naming to docker.io/library/python-docker
```

## View local images

To see a list of images you have on your local machine, you have two options. One is to use the Docker CLI and the other is to use [Docker Desktop](../../desktop/use-desktop/images.md). As you are working in the terminal already, take a look at listing images using the CLI.

To list images, run the `docker images` command.

```console
$ docker images
REPOSITORY      TAG               IMAGE ID       CREATED         SIZE
python-docker   latest            8cae92a8fbd6   3 minutes ago   123MB
```

You should see at least one image listed, including the image you just built `python-docker:latest`.

## Tag images

As mentioned earlier, an image name is made up of slash-separated name components. Name components may contain lowercase letters, digits, and separators. A separator can include a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.

An image is made up of a manifest and a list of layers. Don't worry too much about manifests and layers at this point other than a “tag” points to a combination of these artifacts. You can have multiple tags for an image. Create a second tag for the image you built and take a look at its layers.

To create a new tag for the image you built, run the following command.

```console
$ docker tag python-docker:latest python-docker:v1.0.0
```

The `docker tag` command creates a new tag for an image. It doesn't create a new image. The tag points to the same image and is just another way to reference the image.

Now, run the `docker images` command to see a list of the local images.

```console
$ docker images
REPOSITORY      TAG               IMAGE ID       CREATED         SIZE
python-docker   latest            8cae92a8fbd6   4 minutes ago   123MB
python-docker   v1.0.0            8cae92a8fbd6   4 minutes ago   123MB
...
```

You can see that two images start with `python-docker`. You know they're the same image because if you take a look at the `IMAGE ID` column, you can see that the values are the same for the two images.

Remove the tag you just created. To do this, use the `rmi` command. The `rmi` command stands for remove image.

```console
$ docker rmi python-docker:v1.0.0
Untagged: python-docker:v1.0.0
```

Note that the response from Docker tells you that Docker didn't remove the image, but only “untagged” it. You can check this by running the `docker images` command.

```console
$ docker images
REPOSITORY      TAG               IMAGE ID       CREATED         SIZE
python-docker   latest            8cae92a8fbd6   6 minutes ago   123MB
...
```

Docker removed the image tagged with `:v1.0.0`, but the `python-docker:latest` tag is available on your machine.

## Summary

This section showed how you can use `docker init` to create a Dockerfile and .dockerignore file for a Python application. It then showed you how to build an image. And finally, it showed you how to tag an image and list all images.

Related information:
 - [Dockerfile reference](../../engine/reference/builder.md)
 - [.dockerignore file reference](../../engine/reference/builder.md#dockerignore-file)
 - [docker init CLI reference](../../engine/reference/commandline/init.md)
 - [docker build CLI reference](../../engine/reference/commandline/build.md)
 - [Build with Docker guide](../../build/guide/index.md)

## Next steps

In the next section learn how to run your image as a container.

{{< button text="Run the image as a container" url="run-containers.md" >}}
