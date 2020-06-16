---
description: Engine
keywords: Engine
redirect_from:
- /engine/misc/
- /engine/ce-ee-node-activate/
- /linux/
- /edge/
- /manuals/ # TODO remove this redirect after we've created a landing page for the product manuals section
title: Docker Engine overview
---

Docker Engine is an open source containerization technology for building and
containerizing your applications. Docker Engine acts as a client-server
application with:

* A server with a long-running daemon process [`dockerd`](/engine/reference/commandline/dockerd).
* APIs which specify interfaces that programs can use to talk to and
  instruct the Docker daemon.
* A command line interface (CLI) client [`docker`](/engine/reference/commandline/cli/).

The CLI uses [Docker APIs](api/index.md) to control or interact with the Docker
daemon through scripting or direct CLI commands. Many other Docker applications
use the underlying API and CLI. The daemon creates and manage Docker objects,
such as images, containers, networks, and volumes.

For more details, see [Docker Architecture](../get-started/overview.md#docker-architecture).

## Docker user guide

To learn about Docker in more detail and to answer questions about usage and
implementation, check out the [overview page in "get started"](../get-started/overview.md).

## Installation guides

The [installation section](install/index.md) shows you how to install Docker
on a variety of platforms.

## Release notes

A summary of the changes in each release in the current series can now be found
on the separate [Release Notes page](release-notes/index.md)

## Feature Deprecation Policy

As changes are made to Docker there may be times when existing features
need to be removed or replaced with newer features. Before an existing
feature is removed it is labeled as "deprecated" within the documentation
and remains in Docker for at least 3 stable releases unless specified
explicitly otherwise. After that time it may be removed.

Users are expected to take note of the list of deprecated features each
release and plan their migration away from those features, and (if applicable)
towards the replacement features as soon as possible.

The complete list of deprecated features can be found on the
[Deprecated Features page](deprecated.md).

## Licensing

Docker is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full
license text.
