---
description: Hints, tips and guidelines for writing clean, reliable Dockerfiles
keywords: parent image, images, dockerfile, best practices, hub, official image
title: General best practices for writing Dockerfiles 
---

## Create ephemeral containers

The image defined by your Dockerfile should generate containers that are as
ephemeral as possible. Ephemeral means that the container can be stopped
and destroyed, then rebuilt and replaced with an absolute minimum set up and
configuration.

Refer to [Processes](https://12factor.net/processes) under _The Twelve-factor App_
methodology to get a feel for the motivations of running containers in such a
stateless fashion.

## Exclude with .dockerignore

To exclude files not relevant to the build, without restructuring your source
repository, use a `.dockerignore` file. This file supports exclusion patterns
similar to `.gitignore` files. For information on creating one, see 
[.dockerignore file](../../engine/reference/builder.md#dockerignore-file).

## Don't install unnecessary packages

Avoid installing extra or unnecessary packages just because they might be nice to have. For example, you don’t need to include a text editor in a database image.

When you avoid installing extra or unnecessary packages, your images have reduced complexity, reduced dependencies, reduced file sizes, and reduced build times.

## Decouple applications

Each container should have only one concern. Decoupling applications into
multiple containers makes it easier to scale horizontally and reuse containers.
For instance, a web application stack might consist of three separate
containers, each with its own unique image, to manage the web application,
database, and an in-memory cache in a decoupled manner.

Limiting each container to one process is a good rule of thumb, but it's not a
hard and fast rule. For example, not only can containers be
[spawned with an init process](../../engine/reference/run.md#specify-an-init-process),
some programs might spawn additional processes of their own accord. For
instance, [Celery](https://docs.celeryproject.org/) can spawn multiple worker
processes, and [Apache](https://httpd.apache.org/) can create one process per
request.

Use your best judgment to keep containers as clean and modular as possible. If
containers depend on each other, you can use [Docker container networks](../../network/index.md)
to ensure that these containers can communicate.

## Minimize the number of layers

In older versions of Docker, it was important that you minimized the number of
layers in your images to ensure they were performant. The following features
were added to reduce this limitation:

- Only the instructions `RUN`, `COPY`, and `ADD` create layers. Other instructions
  create temporary intermediate images, and don't increase the size of the build.

- Where possible, use [multi-stage builds](../../build/building/multi-stage.md),
  and only copy the artifacts you need into the final image. This allows you to
  include tools and debug information in your intermediate build stages without
  increasing the size of the final image.

## Sort multi-line arguments

Whenever possible, sort multi-line arguments alphanumerically to make maintenance easier. 
This helps to avoid duplication of packages and make the
list much easier to update. This also makes PRs a lot easier to read and
review. Adding a space before a backslash (`\`) helps as well.

Here’s an example from the [buildpack-deps image](https://github.com/docker-library/buildpack-deps):

```dockerfile
RUN apt-get update && apt-get install -y \
  bzr \
  cvs \
  git \
  mercurial \
  subversion \
  && rm -rf /var/lib/apt/lists/*
```

## Understand build context

See [Build context](../../build/building/context.md) for more information.

## Pipe a Dockerfile through stdin

Docker has the ability to build images by piping a Dockerfile through stdin
with a local or remote build context. Piping a Dockerfile through stdin
can be useful to perform one-off builds without writing a Dockerfile to disk,
or in situations where the Dockerfile is generated, and should not persist
afterward.

> **Note**
>
> The examples in the following sections use [here documents](https://tldp.org/LDP/abs/html/here-docs.html)
> for convenience, but any method to provide the Dockerfile on stdin can be
> used.
> 
> For example, the following commands are equal: 
> 
> ```bash
> echo -e 'FROM busybox\nRUN echo "hello world"' | docker build -
> ```
> 
> ```bash
> docker build -<<EOF
> FROM busybox
> RUN echo "hello world"
> EOF
> ```
> 
> You can substitute the examples with your preferred approach, or the approach
> that best fits your use case.

### Build an image using a Dockerfile from stdin, without sending build context

Use this syntax to build an image using a Dockerfile from stdin, without
sending additional files as build context. The hyphen (`-`) takes the position
of the `PATH`, and instructs Docker to read the build context, which only
contains a Dockerfile, from stdin instead of a directory:

```bash
docker build [OPTIONS] -
```

The following example builds an image using a Dockerfile that is passed through
stdin. No files are sent as build context to the daemon.

```bash
docker build -t myimage:latest -<<EOF
FROM busybox
RUN echo "hello world"
EOF
```

Omitting the build context can be useful in situations where your Dockerfile
doesn't require files to be copied into the image, and improves the build speed,
as no files are sent to the daemon.

If you want to improve the build speed by excluding some files from the build
context, refer to [exclude with .dockerignore](#exclude-with-dockerignore).

> **Note**
>
> If you attempt to build an image using a Dockerfile from stdin, without sending a build context, then the build will fail if you use `COPY` or `ADD`.
> The following example illustrates this:
> 
> ```bash
> # create a directory to work in
> mkdir example
> cd example
> 
> # create an example file
> touch somefile.txt
> 
> docker build -t myimage:latest -<<EOF
> FROM busybox
> COPY somefile.txt ./
> RUN cat /somefile.txt
> EOF
> 
> # observe that the build fails
> ...
> Step 2/3 : COPY somefile.txt ./
> COPY failed: stat /var/lib/docker/tmp/docker-builder249218248/somefile.txt: no such file or directory
> ```

### Build from a local build context, using a Dockerfile from stdin

Use this syntax to build an image using files on your local filesystem, but using
a Dockerfile from stdin. The syntax uses the `-f` (or `--file`) option to
specify the Dockerfile to use, and it uses a hyphen (`-`) as filename to instruct
Docker to read the Dockerfile from stdin:

```bash
docker build [OPTIONS] -f- PATH
```

The following example uses the current directory (`.`) as the build context, and builds
an image using a Dockerfile that is passed through stdin using a [here
document](https://tldp.org/LDP/abs/html/here-docs.html).

```bash
# create a directory to work in
mkdir example
cd example

# create an example file
touch somefile.txt

# build an image using the current directory as context, and a Dockerfile passed through stdin
docker build -t myimage:latest -f- . <<EOF
FROM busybox
COPY somefile.txt ./
RUN cat /somefile.txt
EOF
```

### Build from a remote build context, using a Dockerfile from stdin

Use this syntax to build an image using files from a remote Git repository, 
using a Dockerfile from stdin. The syntax uses the `-f` (or `--file`) option to
specify the Dockerfile to use, using a hyphen (`-`) as filename to instruct
Docker to read the Dockerfile from stdin:

```bash
docker build [OPTIONS] -f- PATH
```

This syntax can be useful in situations where you want to build an image from a
repository that doesn't contain a Dockerfile, or if you want to build with a custom
Dockerfile, without maintaining your own fork of the repository.

The following example builds an image using a Dockerfile from stdin, and adds
the `hello.c` file from the [hello-world](https://github.com/docker-library/hello-world) repository on GitHub.

```bash
docker build -t myimage:latest -f- https://github.com/docker-library/hello-world.git <<EOF
FROM busybox
COPY hello.c ./
EOF
```

> **Note**
>
> When building an image using a remote Git repository as the build context, Docker
> performs a `git clone` of the repository on the local machine, and sends
> those files as the build context to the daemon. This feature requires you to
> install Git on the host where you run the `docker build` command.

## Use multi-stage builds

[Multi-stage builds](../../build/building/multi-stage.md) allow you to
drastically reduce the size of your final image, without struggling to reduce
the number of intermediate layers and files.

Because an image is built during the final stage of the build process, you can
minimize image layers by [leveraging build cache](#leverage-build-cache).

For example, if your build contains several layers and you want to ensure the build cache is reusable, you can order them from the less frequently changed to the more frequently changed. The following list is an example of the order of instructions:

1. Install tools you need to build your application

2. Install or update library dependencies

3. Generate your application

A Dockerfile for a Go application could look like:

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:{{% param "example_go_version" %}}-alpine AS build

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache git

# List project dependencies with go.mod and go.sum
# These layers are only re-built when Gopkg files are updated
WORKDIR /go/src/project/
COPY go.mod go.sum /go/src/project/
# Install library dependencies
RUN go mod download

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
COPY . /go/src/project/
RUN go build -o /bin/project

# This results in a single layer image
FROM scratch
COPY --from=build /bin/project /bin/project
ENTRYPOINT ["/bin/project"]
CMD ["--help"]
```

### Leverage build cache

When building an image, Docker steps through the instructions in your
Dockerfile, executing each in the order specified. As each instruction is
examined, Docker looks for an existing image in its cache,
rather than creating a new, duplicate image.

If you don't want to use the cache at all, you can use the `--no-cache=true`
option on the `docker build` command. However, if you do let Docker use its
cache, it's important to understand when it can, and can't, find a matching
image. The basic rules that Docker follows are outlined below:

- Starting with a parent image that's already in the cache, the next
  instruction is compared against all child images derived from that base
  image to see if one of them was built using the exact same instruction. If
  not, the cache is invalidated.

- In most cases, simply comparing the instruction in the Dockerfile with one
  of the child images is sufficient. However, certain instructions require more
  examination and explanation.

- For the `ADD` and `COPY` instructions, the contents of each file
  in the image are examined and a checksum is calculated for each file.
  The last-modified and last-accessed times of each file aren't considered in
  these checksums. During the cache lookup, the checksum is compared against the
  checksum in the existing images. If anything has changed in any file, such
  as the contents and metadata, then the cache is invalidated.

- Aside from the `ADD` and `COPY` commands, cache checking doesn't look at the
  files in the container to determine a cache match. For example, when processing
  a `RUN apt-get -y update` command the files updated in the container
  aren't examined to determine if a cache hit exists. In that case just
  the command string itself is used to find a match.

Once the cache is invalidated, all subsequent Dockerfile commands generate new
images and the cache isn't used.
