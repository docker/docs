---
description: Guidelines for Official Images on Docker Hub
keywords: Docker, docker, registry, accounts, plans, Dockerfile, Docker Hub, docs, official,image, documentation
title: Docker Official Images
redirect_from:
- /docker-hub/official_repos/
---

The [Docker Official Images](https://hub.docker.com/search?q=&type=image&image_filter=official){:target="_blank" rel="noopener" class="_"} are a
curated set of Docker repositories hosted on Docker Hub. They are
designed to:

* Provide essential base OS repositories (for example,
  [ubuntu](https://hub.docker.com/_/ubuntu/){:target="_blank" rel="noopener" class="_"},
  [centos](https://hub.docker.com/_/centos/){:target="_blank" rel="noopener" class="_"}) that serve as the
  starting point for the majority of users.

* Provide drop-in solutions for popular programming language runtimes, data
  stores, and other services, similar to what a Platform as a Service (PAAS)
  would offer.

* Exemplify [`Dockerfile` best practices](/engine/userguide/eng-image/dockerfile_best-practices/)
  and provide clear documentation to serve as a reference for other `Dockerfile`
  authors.

* Ensure that security updates are applied in a timely manner. This is
  particularly important as Docker Official Images are some of the most
  popular on Docker Hub.

Docker, Inc. sponsors a dedicated team that is responsible for reviewing and
publishing all content in the Docker Official Images. This team works in
collaboration with upstream software maintainers, security experts, and the
broader Docker community.

While it is preferable to have upstream software authors maintaining their
corresponding Docker Official Images, this is not a strict requirement. Creating
and maintaining images for Docker Official Images is a collaborative process. It takes
place openly on GitHub where participation is encouraged. Anyone can provide
feedback, contribute code, suggest process changes, or even propose a new
Official Image.

> **Note**
>
> Docker Official Images are an intellectual property of Docker.

## When to use Docker Official Images

If you are new to Docker, we recommend that you use the Docker Official Images in your
projects. These images have clear documentation, promote best practices,
and are designed for the most common use cases. Advanced users can
review Docker Official Images as part of your `Dockerfile` learning process.

A common rationale for diverging from Docker Official Images is to optimize for
image size. For instance, many of the programming language stack images contain
a complete build toolchain to support installation of modules that depend on
optimized code. An advanced user could build a custom image with just the
necessary pre-compiled libraries to save space.

A number of language stacks such as
[python](https://hub.docker.com/_/python/) and
[ruby](https://hub.docker.com/_/ruby/) have `-slim` tag variants
designed to fill the need for optimization. Even when these "slim" variants are
insufficient, it is still recommended to inherit from an Official Image
base OS image to leverage the ongoing maintenance work, rather than duplicating
these efforts.

## Submitting Feedback for Docker Official Images

All Docker Official Images contain a **User Feedback** section in their
documentation which covers the details for that specific repository. In most
cases, the GitHub repository which contains the Dockerfiles for an Official
Repository also has an active issue tracker. General feedback and support
questions should be directed to `#docker-library` on [Libera.Chat IRC](https://libera.chat).

## Creating a Docker Official Image

From a high level, an Official Image starts out as a proposal in the form
of a set of GitHub pull requests. Detailed and objective proposal
requirements are documented in the following GitHub repositories:

* [docker-library/official-images](https://github.com/docker-library/official-images){:target="_blank" rel="noopener" class="_"}

* [docker-library/docs](https://github.com/docker-library/docs){:target="_blank" rel="noopener" class="_"}

The Docker Official Images team, with help from community contributors, formally
review each proposal and provide feedback to the author. This initial review
process may require a bit of back-and-forth before the proposal is accepted.

There are also subjective considerations during the review process. These
subjective concerns boil down to the basic question: "is this image generally
useful?" For example, the [python](https://hub.docker.com/_/python/){:target="_blank" rel="noopener" class="_"}
Docker Official Image is "generally useful" to the larger Python developer
community, whereas an obscure text adventure game written in Python last week is
not.

Once a new proposal is accepted, the author is responsible for keeping
their images up-to-date and responding to user feedback. The Official
Repositories team becomes responsible for publishing the images and
documentation on Docker Hub. Updates to the Docker Official Image follow the same pull request process, though with less review. The Docker Official Images team ultimately acts as a gatekeeper for all changes, which helps mitigate the risk of quality and security issues from being introduced.
