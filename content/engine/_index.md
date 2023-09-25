---
description: Engine
keywords: Engine
title: Docker Engine overview
grid:
- title: Install Docker Engine
  description: Learn how to install the open source Docker Engine for your distribution.
  icon: download
  link: /engine/install
- title: Storage
  description: Use persistent data with Docker containers.
  icon: database
  link: /storage
- title: Networking
  description: Manage network connections between containers.
  icon: network_node
  link: /network
- title: Container logs
  description: Learn how to view and read container logs.
  icon: feed
  link: /config/containers/logging/
- title: Prune
  description: Tidy up unused resources.
  icon: cut
  link: /config/pruning
- title: Configure the daemon
  description: Delve into the configuration options of the Docker daemon.
  icon: tune
  link: /config/daemon
- title: Rootless mode
  description: Run Docker without root privileges.
  icon: security
  link: /engine/security/rootless
- title: Deprecated features
  description: Find out what features of Docker Engine you should stop using.
  icon: folder_delete
  link: /engine/deprecated/
- title: Release notes
  description: Read the release notes for the latest version.
  icon: note_add
  link: /engine/release-notes
aliases:
- /edge/
- /engine/ce-ee-node-activate/
- /engine/misc/
- /linux/
- /manuals/
---

Docker Engine is an open source containerization technology for building and
containerizing your applications. Docker Engine acts as a client-server
application with:

- A server with a long-running daemon process
  [`dockerd`](/engine/reference/commandline/dockerd).
- APIs which specify interfaces that programs can use to talk to and instruct
  the Docker daemon.
- A command line interface (CLI) client
  [`docker`](/engine/reference/commandline/cli/).

The CLI uses [Docker APIs](api/index.md) to control or interact with the Docker
daemon through scripting or direct CLI commands. Many other Docker applications
use the underlying API and CLI. The daemon creates and manage Docker objects,
such as images, containers, networks, and volumes.

For more details, see
[Docker Architecture](../get-started/overview.md#docker-architecture).

{{< grid >}}

## Licensing

The Docker Engine is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full license
text.
