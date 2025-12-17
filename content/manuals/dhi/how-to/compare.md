---
title: Compare Docker Hardened Images
linktitle: Compare images
description: Learn how to compare Docker Hardened Images with other container images to evaluate security improvements and differences.
keywords: compare docker images, docker scout compare, image comparison, vulnerability comparison, security comparison
weight: 40
---

Docker Hardened Images (DHIs) are designed to provide enhanced security,
minimized attack surfaces, and production-ready foundations for your
applications. Comparing a DHI to a standard image helps you understand the
security improvements, package differences, and overall benefits of adopting
hardened images.

This page explains how to use Docker Scout to compare a Docker Hardened Image
with another image, such as a Docker Official Image (DOI) or a custom image, to
evaluate differences in vulnerabilities, packages, and configurations.

## Compare images using Docker Scout

Docker Scout provides a built-in comparison feature that lets you analyze the
differences between two images. This is useful for:

- Evaluating the security improvements when migrating from a standard image to a
  DHI
- Understanding package and vulnerability differences between image variants
- Assessing the impact of customizations or updates

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

This command provides a detailed comparison including:

- Vulnerability differences (CVEs added, removed, or changed)
- Package differences (packages added, removed, or updated)
- Overall security posture improvements

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
the two images, making it easier to identify the security improvements and
changes.

### Show overview only

For a concise overview of the comparison results, you can extract just the
overview section using standard shell tools:

```console
$ docker scout compare dhi.io/node:22-debian13 \
    --to node:22 \
    --platform linux/amd64 \
    --ignore-unchanged \
    2>/dev/null | sed -n '/## Overview/,/^  ## /p' | head -n -1
```

The result is a clean summary showing the key differences between the two
images. Example output:

```console
  ## Overview
  
                      │                    Analyzed Image                     │              Comparison Image
  ────────────────────┼───────────────────────────────────────────────────────┼─────────────────────────────────────────────
    Target            │  dhi.io/node:22-debian13                              │  node:22
      digest          │  55d471f61608                                         │  9ee3220f602f
      tag             │  22-debian13                                          │  22
      platform        │ linux/amd64                                           │ linux/amd64
      provenance      │ https://github.com/docker-hardened-images/definitions │ https://github.com/nodejs/docker-node.git
                      │  9fe491f53122b84eebba81e13f20157c18c10de2             │  bf78d7603fbea92cd3652edb3b2edadd6f5a3fe8
      vulnerabilities │    0C     0H     0M     0L                            │    0C     1H     3M   153L     4?
                      │           -1     -3   -153     -4                     │
      size            │ 41 MB (-367 MB)                                       │ 408 MB
      packages        │ 19 (-726)                                             │ 745
                      │                                                       │
```

## Interpret comparison results

The comparison output includes the following sections.

### Overview

The overview section provides high-level statistics about both images:

- Target and comparison image details (digest, tag, platform, provenance)
- Vulnerability counts for each image
- Size comparison
- Package counts

Look for:

- Vulnerability reductions (negative numbers in the delta row)
- Size reductions showing storage efficiency
- Package count reductions indicating a minimal attack surface

### Environment Variables

The environment variables section shows environment variables that differ between
the two images, prefixed with `+` for added or `-` for removed.

Look for:

- Removed environment variables that may have been necessary for your specific use-case

### Labels

The labels section displays labels that differ between the two images, prefixed
with `+` for added or `-` for removed.

### Packages and Vulnerabilities

The packages and vulnerabilities section lists all package differences and their
associated security vulnerabilities. Packages are prefixed with:

- `-` for packages removed from the target image (not present in the compared image)
- `+` for packages added to the target image (not present in the base image)
- `↑` for packages upgraded in the target image
- `↓` for packages downgraded in the target image

For packages with associated vulnerabilities, the CVEs are listed with their
severity levels and identifiers.

Look for:

- Removed packages and vulnerabilities: Indicates a reduced attack surface in the DHI
- Added packages: May indicate DHI-specific tooling or dependencies
- Upgraded packages: Shows version updates that may include security fixes

## When to compare images

### Evaluate migration benefits

Before migrating from a Docker Official Image to a DHI, compare them to
understand the security improvements. For example:

```console
$ docker scout compare dhi.io/python:3.13 \
    --to python:3.13 \
    --platform linux/amd64 \
    --ignore-unchanged
```

This helps justify the migration by showing concrete vulnerability reductions
and package minimization.

### Assess customization impact

After customizing a DHI, compare the customized version with the original to
ensure you haven't introduced new vulnerabilities. For example:

```console
$ docker scout compare <your-namespace>/dhi-python:3.13-custom \
    --to dhi.io/python:3.13 \
    --platform linux/amd64
```

### Track updates over time

Compare different versions of the same DHI to see what changed between releases. For example:

```console
$ docker scout compare dhi.io/node:22-debian13 \
    --to dhi.io/node:20-debian12 \
    --platform linux/amd64 \
    --ignore-unchanged
```
