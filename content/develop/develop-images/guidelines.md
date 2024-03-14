---
description: Hints, tips and guidelines for writing clean, reliable Dockerfiles
keywords: parent image, images, dockerfile, best practices, hub, official image
title: General best practices for writing Dockerfiles
---

## Use multi-stage builds

Multi-stage builds let you reduce the size of your final image, by creating a
cleaner separation between the building of your image and the final output.
Split your Dockerfile instructions into distinct stages to make sure that the
resulting output only contains the files that's needed to run the application.

Using multiple stages can also let you build more efficiently by executing
build steps in parallel.

See [Multi-stage builds](../../build/building/multi-stage.md) for more
information.

## Exclude with .dockerignore

To exclude files not relevant to the build, without restructuring your source
repository, use a `.dockerignore` file. This file supports exclusion patterns
similar to `.gitignore` files. For information on creating one, see
[Dockerignore file](../../build/building/context.md#dockerignore-files).

## Create ephemeral containers

The image defined by your Dockerfile should generate containers that are as
ephemeral as possible. Ephemeral means that the container can be stopped
and destroyed, then rebuilt and replaced with an absolute minimum set up and
configuration.

Refer to [Processes](https://12factor.net/processes) under _The Twelve-factor App_
methodology to get a feel for the motivations of running containers in such a
stateless fashion.

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

## Leverage build cache

When building an image, Docker steps through the instructions in your
Dockerfile, executing each in the order specified. For each instruction, Docker
checks whether it can reuse the instruction from the build cache.

Understanding how the build cache works, and how cache invalidation occurs,
is critical for ensuring faster builds.
For more information about the Docker build cache and how to optimize your builds,
see [Docker build cache](../../build/cache/_index.md).

## Pin base image versions

Image tags are mutable, meaning a publisher can update a tag to point to a new
image. This is useful because it lets publishers update tags to point to
newer versions of an image. And as an image consumer, it means you
automatically get the new version when you re-build your image.

For example, if you specify `FROM alpine:3.19` in your Dockerfile, `3.19`
resolves to the latest patch version for `3.19`.

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine:3.19
```

At one point in time, the `3.19` tag might point to version 3.19.1 of the
image. If you rebuild the image 3 months later, the same tag might point to a
different version, such as 3.19.4. This publishing workflow is best practice,
and most publishers use this tagging strategy, but it isn't enforced.

The downside with this is that you're not guaranteed to get the same for every
build. This could result in breaking changes, and it means you also don't have
an audit trail of the exact image versions that you're using.

To fully secure your supply chain integrity, you can pin the image version to a
specific digest. By pinning your images to a digest, you're guaranteed to
always use the same image version, even if a publisher replaces the tag with a
new image. For example, the following Dockerfile pins the Alpine image to the
same tag as earlier, `3.19`, but this time with a digest reference as well.

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine:3.19@sha256:13b7e62e8df80264dbb747995705a986aa530415763a6c58f84a3ca8af9a5bcd
```

With this Dockerfile, even if the publisher updates the `3.19` tag, your builds
would still use the pinned image version:
`13b7e62e8df80264dbb747995705a986aa530415763a6c58f84a3ca8af9a5bcd`.

While this helps you avoid unexpected changes, it's also more tedious to have
to look up and include the image digest for base image versions manually each
time you want to update it. And you're opting out of automated security fixes,
which is likely something you want to get.

Docker Scout has a built-in [**Outdated base images**
policy](../../scout/policy/_index.md#outdated-base-images) that checks for
whether the base image version you're using is in fact the latest version. This
policy also checks if pinned digests in your Dockerfile correspond to the
correct version. If a publisher updates an image that you've pinned, the policy
evaluation returns a non-compliant status, indicating that you should update
your image.

Docker Scout also supports an automated remediation workflow for keeping your
base images up-to-date. When a new image digest is available, Docker Scout can
automatically raise a pull request on your repository to update your
Dockerfiles to use the latest version. This is better than using a tag that
changes the version automatically, because you're in control and you have an
audit trail of when and how the change occurred.

For more information about automatically updating your base images with Docker
Scout, see
[Remediation](../../scout/policy/remediation.md#automatic-base-image-updates)
