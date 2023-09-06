---
description: Using Docker for continuous integration
keywords: ci, build
title: Continuous integration with Docker
aliases:
  - /ci-cd/best-practices/
---

Continuous Integration (CI) is the part of the development process where you're
looking to get your code changes merged with the main branch of the project. At
this point, development teams run tests and builds to vet that the code changes
don't cause any unwanted or unexpected behaviors.

![Git branches about to get merged](./images/continuous-integration.svg)

There are several uses for Docker at this stage of development, even if you
don't end up packaging your application as a container image.

## Docker as a build environment

Containers are reproducible, isolated environments that yield predictable
results. Building and testing your application in a Docker container makes it
easier to prevent unexpected behaviors from occurring. Using a Dockerfile, you
define the exact requirements for the build environment, including programming
runtimes, operating system, binaries, and more.

Using Docker to manage your build environment also eases maintenance. For
example, updating to a new version of a programming runtime can be as simple as
changing a tag or digest in a Dockerfile. No need to SSH into a pet VM to
manually reinstall a newer version and update the related configuration files.

Additionally, just as you expect third-party open source packages to be secure,
the same should go for your build environment. You can scan and index a builder
image, just like you would for any other containerized application.

The following links provide instructions for how you can get started using
Docker for building your applications in CI:

- [GitHub Actions](https://docs.github.com/en/actions/creating-actions/creating-a-docker-container-action)
- [GitLab](https://docs.gitlab.com/runner/executors/docker.html)
- [Circle CI](https://circleci.com/docs/using-docker/)
- [Render](https://render.com/docs/docker)

### Docker in Docker

You can also use a Dockerized build environment to build container images using
Docker. That is, your build environment runs inside a container which itself is
equipped to run Docker builds. This method is referred to as "Docker in Docker".

Docker provides an official [Docker image](https://hub.docker.com/_/docker)
that you can use for this purpose.

## What's next

Docker maintains a set of official GitHub Actions that you can use to build,
annotate, and push container images on the GitHub Actions platform. See
[Introduction to GitHub Actions](./github-actions/index.md) to learn more and
get started.
