---
description: Manage Builds and Images in Docker Cloud
keywords: builds, images, Cloud
title: Builds and images overview
notoc: true
---

Docker Cloud provides a hosted registry service where you can create
repositories to store your Docker images. You can choose to push images to the
repositories, or link to your source code and build them directly in Docker
Cloud.

You can build images manually, or set up automated builds to rebuild your Docker
image on each `git push` to the source code. You can also create automated
tests, and when the tests pass use autoredeploy to automatically update your
running services when a build passes its tests.

* [Repositories in Docker Cloud](repos.md)
* [Push images to Docker Cloud](push-images.md)
* [Link to a source code repository](link-source.md)
* [Automated builds](automated-build.md)
* [Automated repository tests](automated-testing.md)
* [Advanced options for Autobuild and Autotest](advanced.md)

![Docker Cloud repository General view](images/repo-general.png){:width="650px"}
