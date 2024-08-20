---
title: Build your Rust image
keywords: rust, build, images, dockerfile
description: Learn how to build your first Rust Docker image
---

## Prerequisites

* You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
* You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

## Overview

This guide walks you through building your first Rust image. An image
includes everything needed to run an application - the code or binary, runtime,
dependencies, and any other file system objects required.

## Get the sample application

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/docker/docker-rust-hello
```

## Create a Dockerfile for Rust

Now that you have an application, you can use `docker init` to create a
Dockerfile for it. Inside the `docker-rust-hello` directory, run the `docker
init` command. `docker init` provides some default configuration, but you'll
need to answer a few questions about your application. Refer to the following
example to answer the prompts from `docker init` and use the same answers for
your prompts.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

? What application platform does your project use? Rust
? What version of Rust do you want to use? 1.70.0
? What port does your server listen on? 8000
```

You should now have the following new files in your `docker-rust-hello`
directory:
 - Dockerfile
 - .dockerignore
 - compose.yaml
 - README.Docker.md

For building an image, only the Dockerfile is necessary. Open the Dockerfile
in your favorite IDE or text editor and see what it contains. To learn more
about Dockerfiles, see the [Dockerfile reference](../../reference/dockerfile.md).

## .dockerignore file

When you run `docker init`, it also creates a [`.dockerignore`](../../reference/dockerfile.md#dockerignore-file) file. Use the `.dockerignore` file to specify patterns and paths that you don't want copied into the image in order to keep the image as small as possible. Open up the `.dockerignore` file in your favorite IDE or text editor and see what's inside already.

## Build an image

Now that you’ve created the Dockerfile, you can build the image. To do this, use
the `docker build` command. The `docker build` command builds Docker images from
a Dockerfile and a context. A build's context is the set of files located in
the specified PATH or URL. The Docker build process can access any of the files
located in this context.

The build command optionally takes a `--tag` flag. The tag sets the name of the
image and an optional tag in the format `name:tag`. If you don't pass a tag,
Docker uses "latest" as its default tag.

Build the Docker image.

```console
$ docker build --tag docker-rust-image .
```

You should see output like the following.

```console
[+] Building 62.6s (14/14) FINISHED
 => [internal] load .dockerignore                                                                                                    0.1s
 => => transferring context: 2B                                                                                                      0.0s 
 => [internal] load build definition from Dockerfile                                                                                 0.1s
 => => transferring dockerfile: 2.70kB                                                                                               0.0s 
 => resolve image config for docker.io/docker/dockerfile:1                                                                           2.3s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:39b85bbfa7536a5feceb7372a0817649ecb2724562a38360f4d6a7782a409b14      0.0s
 => [internal] load metadata for docker.io/library/debian:bullseye-slim                                                              1.9s
 => [internal] load metadata for docker.io/library/rust:1.70.0-slim-bullseye                                                         1.7s 
 => [build 1/3] FROM docker.io/library/rust:1.70.0-slim-bullseye@sha256:585eeddab1ec712dade54381e115f676bba239b1c79198832ddda397c1f  0.0s
 => [internal] load build context                                                                                                    0.0s 
 => => transferring context: 35.29kB                                                                                                 0.0s 
 => [final 1/3] FROM docker.io/library/debian:bullseye-slim@sha256:7606bef5684b393434f06a50a3d1a09808fee5a0240d37da5d181b1b121e7637  0.0s 
 => CACHED [build 2/3] WORKDIR /app                                                                                                  0.0s
 => [build 3/3] RUN --mount=type=bind,source=src,target=src     --mount=type=bind,source=Cargo.toml,target=Cargo.toml     --mount=  57.7s 
 => CACHED [final 2/3] RUN adduser     --disabled-password     --gecos ""     --home "/nonexistent"     --shell "/sbin/nologin"      0.0s
 => CACHED [final 3/3] COPY --from=build /bin/server /bin/                                                                           0.0s
 => exporting to image                                                                                                               0.0s
 => => exporting layers                                                                                                              0.0s
 => => writing image sha256:f1aa4a9f58d2ecf73b0c2b7f28a6646d9849b32c3921e42adc3ab75e12a3de14                                         0.0s
 => => naming to docker.io/library/docker-rust-image
```

## View local images

To see a list of images you have on your local machine, you have two options. One is to use the Docker CLI and the other is to use [Docker Desktop](../../desktop/use-desktop/images.md). As you are working in the terminal already, take a look at listing images using the CLI.

To list images, run the `docker images` command.

```console
$ docker images
REPOSITORY                TAG               IMAGE ID       CREATED         SIZE
docker-rust-image         latest            8cae92a8fbd6   3 minutes ago   123MB
```

You should see at least one image listed, including the image you just built `docker-rust-image:latest`.

## Tag images

As mentioned earlier, an image name is made up of slash-separated name components. Name components may contain lowercase letters, digits, and separators. A separator can include a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.

An image is made up of a manifest and a list of layers. Don't worry too much about manifests and layers at this point other than a "tag" points to a combination of these artifacts. You can have multiple tags for an image. Create a second tag for the image you built and take a look at its layers.

To create a new tag for the image you built, run the following command.

```console
$ docker tag docker-rust-image:latest docker-rust-image:v1.0.0
```

The `docker tag` command creates a new tag for an image. It doesn't create a new image. The tag points to the same image and is just another way to reference the image.

Now, run the `docker images` command to see a list of the local images.

```console
$ docker images
REPOSITORY                TAG               IMAGE ID       CREATED         SIZE
docker-rust-image         latest            8cae92a8fbd6   4 minutes ago   123MB
docker-rust-image         v1.0.0            8cae92a8fbd6   4 minutes ago   123MB
rust                      latest            be5d294735c6   4 minutes ago   113MB
```

You can see that two images start with `docker-rust-image`. You know they're the same image because if you take a look at the `IMAGE ID` column, you can see that the values are the same for the two images.

Remove the tag you just created. To do this, use the `rmi` command. The `rmi` command stands for remove image.

```console
$ docker rmi docker-rust-image:v1.0.0
Untagged: docker-rust-image:v1.0.0
```

Note that the response from Docker tells you that Docker didn't remove the image, but only "untagged" it. You can check this by running the `docker images` command.

```console
$ docker images
REPOSITORY               TAG               IMAGE ID       CREATED         SIZE
docker-rust-image        latest            8cae92a8fbd6   6 minutes ago   123MB
rust                     latest            be5d294735c6   6 minutes ago   113MB
```

Docker removed the image tagged with `:v1.0.0`, but the `docker-rust-image:latest` tag is available on your machine.

## Summary

This section showed how you can use `docker init` to create a Dockerfile and .dockerignore file for a Rust application. It then showed you how to build an image. And finally, it showed you how to tag an image and list all images.

Related information:
 - [Dockerfile reference](../../reference/dockerfile.md)
 - [.dockerignore file](../../reference/dockerfile.md#dockerignore-file)
 - [docker init CLI reference](../../reference/cli/docker/init.md)
 - [docker build CLI reference](../../reference/cli/docker/buildx/build.md)


## Next steps

In the next section learn how to run your image as a container.

{{< button text="Run the image as a container" url="run-containers.md" >}}
