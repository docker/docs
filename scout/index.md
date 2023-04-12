---
title: Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves
description: >
  Docker Scout helps you understand your images and their dependencies
redirect_from:
  - /atomist/
  - /atomist/try-atomist/
  - /atomist/get-started/
  - /atomist/configure/settings/
  - /atomist/configure/advisories/
  - /atomist/integrate/github/
  - /atomist/integrate/deploys/
---

{% include scout-early-access.md %}

Docker Scout is a collection of software supply chain features that appear
throughout Docker user interfaces and the command line interface (CLI). These features provide detailed
insights into the composition and security of container images.

Docker Scout analyzes image contents and generates a detailed report of
packages and vulnerabilities that it detects. Docker Scout can also help provide
you with suggestions for how you can remediate issues discovered by the image
analysis.

The [image details view](./image-details-view.md) in Docker Desktop and the tag
details pages on Docker Hub are both powered by Docker Scout.

You can view and interact with Docker Scout from your terminal through the
`docker scout`
[plugin for Docker CLI](../engine/reference/commandline/scout_cves.md).

There's also a [web UI](https://dso.docker.com/explore){: target="\_blank"
rel="noopener" } that you can use to explore additional information about
images, packages, and CVEs.

## Get started with Docker Scout

To start using Docker Scout, turn on any of the following features:

- [Advanced image analysis in Docker Hub](./advanced-image-analysis.md)
- [Artifactory integration](./artifactory.md)

Enabling either feature requires a paid Docker subscription. Once enabled,
Docker Scout is activated for your organization automatically.

### Advanced image analysis

Advanced image analysis is a feature in Docker Hub which, when enabled, triggers
a Docker Scout analysis every time you push an image. The analysis updates
continuously, meaning that the vulnerability report for an image is always up to
date as new CVEs are discovered. No need to re-scan an image.

For more information, see
[Advanced image analysis](./advanced-image-analysis.md).

### Artifactory integration

Users of JFrog Artifactory, or JFrog Container Registry, can integrate Docker
Scout to enable automatic analysis of images locally and remotely. For more information, see
[Artifactory integration](./artifactory.md).

## Docker Scout CLI

The `docker scout` CLI plugin provides a terminal interface for Docker Scout.

Using the CLI, you can analyze images and view the analysis report in text
format. You can print the results directly to stdout, or export them to a file
using a structured format, such as Static Analysis Results Interchange Format
(SARIF). For more information about how to use the `docker scout` CLI, see the
[reference documentation](../engine/reference/commandline/scout_cves.md).
