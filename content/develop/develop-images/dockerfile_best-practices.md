---
description: Overview of a Dockerfile and introduction to best practices
keywords: parent image, images, dockerfile, best practices, hub, official image
title: Overview of best practices for writing Dockerfiles
aliases:
- /articles/dockerfile_best-practices/
- /engine/articles/dockerfile_best-practices/
- /docker-cloud/getting-started/intermediate/optimize-dockerfiles/
- /docker-cloud/tutorials/optimize-dockerfiles/
- /engine/userguide/eng-image/dockerfile_best-practices/
---

This topic covers recommended best practices and methods for building
efficient images. It provides [general guidelines for your Dockerfiles](guidelines.md) and more [specific best practices for each Dockerfile instruction](instructions.md). 

## What is a Dockerfile?

Docker builds images automatically by reading the instructions from a
Dockerfile which is a text file that contains all commands, in order, needed to
build a given image. A Dockerfile adheres to a specific format and set of
instructions which you can find at [Dockerfile reference](../../reference/dockerfile.md).

A Docker image consists of read-only layers each of which represents a
Dockerfile instruction. The layers are stacked and each one is a delta of the
changes from the previous layer. 

```dockerfile
# syntax=docker/dockerfile:1

FROM ubuntu:22.04
COPY . /app
RUN make /app
CMD python /app/app.py
```

In the example above, each instruction creates one layer:

- `FROM` creates a layer from the `ubuntu:22.04` Docker image.
- `COPY` adds files from your Docker client's current directory.
- `RUN` builds your application with `make`.
- `CMD` specifies what command to run within the container.

When you run an image and generate a container, you add a new writable layer, also called the container layer, on top of the underlying layers. All changes made to
the running container, such as writing new files, modifying existing files, and
deleting files, are written to this writable container layer.

## Additional resources

* [Dockerfile reference](../../reference/dockerfile.md)
* [More about Automated builds](../../docker-hub/builds/index.md)
* [Guidelines for creating Docker Official Images](../../trusted-content/official-images/_index.md)
* [Best practices to containerize Node.js web applications with Docker](https://snyk.io/blog/10-best-practices-to-containerize-nodejs-web-applications-with-docker)
* [More about base images](../../build/building/base-images.md)
* [More on image layers and how Docker builds and stores images](../../storage/storagedriver/index.md).

## Examples of Docker Official Images

These Official Images have exemplary Dockerfiles:

* [Go](https://hub.docker.com/_/golang/)
* [Perl](https://hub.docker.com/_/perl/)
* [Hy](https://hub.docker.com/_/hylang/)
* [Ruby](https://hub.docker.com/_/ruby/)
