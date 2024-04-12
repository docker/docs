---
title: Contributing to Docker Official Images
description: |
  This article describes how Docker Official Images are created,
  and how you can contribute or leave feedback.
keywords: docker official images, doi, contributing, upstream, open source
---

Docker, Inc. sponsors a dedicated team that's responsible for reviewing and
publishing all content in Docker Official Images. This team works in
collaboration with upstream software maintainers, security experts, and the
broader Docker community.

While it's preferable to have upstream software authors maintaining their
Docker Official Images, this isn't a strict requirement. Creating
and maintaining images for Docker Official Images is a collaborative process.
It takes place [openly on GitHub](https://github.com/docker-library/official-images)
where participation is encouraged. Anyone can provide feedback, contribute
code, suggest process changes, or even propose a new Official Image.

## Creating a Docker Official Image

From a high level, an Official Image starts out as a proposal in the form
of a set of GitHub pull requests. The following GitHub repositories detail the proposal requirements:

- [Docker Official Images repository on GitHub](https://github.com/docker-library/official-images#readme)
- [Documentation for Docker Official Images](https://github.com/docker-library/docs#readme)

The Docker Official Images team, with help from community contributors, formally
review each proposal and provide feedback to the author. This initial review
process can be lengthy, often requiring a bit of back-and-forth before the proposal is accepted.

There are subjective considerations during the review process. These
subjective concerns boil down to the basic question: "is this image generally
useful?" For example, the [Python](https://hub.docker.com/_/python/)
Docker Official Image is "generally useful" to the larger Python developer
community, whereas an obscure text adventure game written in Python last week is
not.

Once a new proposal is accepted, the author is responsible for keeping their
images and documentation up-to-date and responding to user feedback. Docker is
responsible for building and publishing the images on Docker Hub. Updates to
Docker Official Images follow the same pull request process as for new images,
although the review process for updates is more streamlined. The Docker Official
Images team ultimately acts as a gatekeeper for all changes, which helps
ensures consistency, quality, and security.

## Submitting feedback for Docker Official Images

All Docker Official Images contain a **User Feedback** section in their
documentation which covers the details for that specific repository. In most
cases, the GitHub repository which contains the Dockerfiles for an Official
Image also has an active issue tracker.

General feedback and support questions about Docker Official Images
should be directed to the `#general` channel in the [Docker Community Slack](https://dockr.ly/comm-slack).

If you're a maintainer or contributor to Docker Official Images and you're
looking for help or advice, use the `#docker-library` channel on [Libera.Chat IRC](https://libera.chat).
