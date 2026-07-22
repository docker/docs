---
title: Search and evaluate Docker Hardened Images
linktitle: Search and evaluate
description: Learn how to find, compare, and evaluate Docker Hardened Images using the Docker Hub catalog, DHI CLI, and Docker Scout.
keywords: search docker images, image variants, docker hub catalog, compare docker images, docker scout compare, image comparison, vulnerability comparison
weight: 10
aliases:
  - /dhi/how-to/compare/
  - /dhi/how-to/explore/
---

## Search the catalog

You can browse, search, or filter images by category in the [Docker Hub
catalog](https://hub.docker.com/hardened-images/catalog). For details about
the catalog interface, see [Docker Hub](/dhi/tools/hub/).

Alternatively, use the [DHI MCP server](/dhi/tools/mcp/) to search and
inspect the catalog directly from your AI assistant, or use the [DHI
CLI](/dhi/tools/cli/) to browse the catalog from the command line:

```console
$ docker dhi catalog list
```

Filter by image type, name, or compliance requirements:

```console
$ docker dhi catalog list --type image
$ docker dhi catalog list --filter python
$ docker dhi catalog list --fips
$ docker dhi catalog list --stig
```

To view repository details, including available tags and CVE counts:

```console
$ docker dhi catalog get python
```

## Compare and evaluate images

Docker Scout lets you analyze the differences between two images. Comparing a
DHI to a standard image helps you understand the security improvements, package
differences, and overall benefits of adopting hardened images.

Comparison is useful for:

- Evaluating the security improvements when migrating from a standard image to a DHI
- Understanding package and vulnerability differences between image variants
- Assessing the impact of customizations or updates

### Prerequisites

Before comparing images:

- Install [Docker Desktop](/desktop/) to use Docker Scout comparison features.
- Sign in to `dhi.io` for Docker Hardened Images:

  ```console
  $ docker login dhi.io
  ```

### Basic comparison

To compare a Docker Hardened Image with another image, use the [`docker scout
compare`](/reference/cli/docker/scout/compare/) command:

```console
$ docker scout compare dhi.io/<image>:<tag> \
    --to <comparison-image>:<tag> \
    --platform <platform>
```

For example, to compare a DHI Node.js image with the official Node.js image:

```console
$ docker scout compare dhi.io/node:22-debian13 \
    --to node:22 \
    --platform linux/amd64
```

The output shows an overview at the top with key comparison metrics, followed by
detailed package and vulnerability information. Example overview:

```console
  ## Overview

                      │                    Analyzed Image                     │              Comparison Image
  ────────────────────┼───────────────────────────────────────────────────────┼─────────────────────────────────────────────
    Target            │  dhi.io/node:22-debian13                              │  node:22
      digest          │  55d471f61608                                         │  9ee3220f602f
      platform        │ linux/amd64                                           │ linux/amd64
      vulnerabilities │    0C     0H     0M     0L                            │    0C     1H     3M   153L     4?
                      │           -1     -3   -153     -4                     │
      size            │ 41 MB (-367 MB)                                       │ 408 MB
      packages        │ 19 (-726)                                             │ 745
```

### Filter unchanged packages

To focus only on the differences and ignore unchanged packages, use the
`--ignore-unchanged` flag:

```console
$ docker scout compare dhi.io/node:22-debian13 \
    --to node:22 \
    --platform linux/amd64 \
    --ignore-unchanged
```

This output highlights only the packages and vulnerabilities that differ between
the two images.
