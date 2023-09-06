---
description: High-level overview of the Registry
keywords: registry, on-prem, images, tags, repository, distribution
title: Docker Registry
aliases:
- /registry/overview/
---

{{< include "registry.md" >}}

## What it is

The Registry is a stateless, highly scalable server side application that stores
and lets you distribute Docker images. The Registry is open-source, under the
permissive [Apache license](https://en.wikipedia.org/wiki/Apache_License).
You can find the source code on
[GitHub](https://github.com/distribution/distribution).

## Why use it

You should use the Registry if you want to:

 * Tightly control where your images are being stored
 * Fully own your images distribution pipeline
 * Integrate image storage and distribution tightly into your in-house development workflow

## Alternatives

If you're looking for a zero maintenance, ready-to-go solution, [Docker Hub](https://hub.docker.com), provides a free-to-use, hosted Registry, plus additional features (organization accounts,
Automated builds, and more).

## Requirements

The Registry is compatible with Docker engine version 1.6.0 or later.

## Basic commands

Start your registry:

```console
$ docker run -d -p 5000:5000 --name registry registry:2
```

Pull (or build) an image from the hub:
    
```console
$ docker pull ubuntu
```

Tag the image so that it points to your registry:

```console
$ docker image tag ubuntu localhost:5000/myfirstimage
```

Push it:

```console
$ docker push localhost:5000/myfirstimage
```

Pull it back:

```console
$ docker pull localhost:5000/myfirstimage
```

Now stop your registry and remove all data:

```console
$ docker container stop registry && docker container rm -v registry
 ```

## Next

Read the [detailed introduction about the registry](introduction.md) or jump directly to [deployment instructions](deploying.md).