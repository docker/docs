---
description: Fun stuff to do with your registry
keywords: registry, on-prem, images, tags, repository, distribution, recipes, advanced
title: Recipes overview
---

{{< include "registry.md" >}}

This section provides end-to-end scenarios for exotic or otherwise advanced use-cases.
These examples are not useful for most standard set-ups.

## Requirements

Before following these steps, work through the [deployment guide](../deploying.md).

You must meet the following requirements:

 * Make sure you understand Docker security requirements, and how to configure your Docker Engine properly
 * You have installed Docker Compose
 * You have a certificate from a known CA instead of self-signed certificates. This is highly recommended.
 * Inside the current directory, you have a X509 `domain.crt` and `domain.key`, for the CN `myregistrydomain.com`
 * You have stopped and removed any previously running registry (typically `docker container stop registry && docker container rm -v registry`)

## Recipes

 * [Using Apache as an authenticating proxy](apache.md)
 * [Using Nginx as an authenticating proxy](nginx.md)
 * [Running a Registry on macOS](osx-setup-guide.md)
 * [Mirror the Docker Hub](mirror.md)