---
title: Docker Build Overview
weight: 10
description: Learn about Docker Build and its components.
keywords: build, buildkit, buildx, architecture
aliases:
- /build/install-buildx/
- /build/architecture/
---

Docker Build uses a client-server architecture with two main components:

- **Client**: Buildx serves as the user interface for running and managing builds
- **Server**: BuildKit handles the actual build execution

## How Docker Build works

When you invoke a build:

1. Buildx (client) sends a build request to BuildKit (backend).
2. BuildKit resolves build instructions and executes build steps.
3. Build output returns to the client or uploads to a registry like Docker Hub.

Both Buildx and BuildKit come pre-installed with Docker Desktop and Docker Engine.
The `docker build` command uses the default BuildKit bundled with Docker.

## Buildx

Buildx is the CLI tool for running builds. The `docker build` command
is a wrapper around Buildx functionality.

When you run `docker build`, Buildx:

- Interprets your build options
- Sends build requests to BuildKit backend
- Manages build execution

Beyond running builds, Buildx enables you to:

- Create and manage BuildKit backends (builders)
- Manage images in registries
- Run multiple builds concurrently
- Access advanced build features


Docker Buildx comes installed by default with Docker Desktop. For manual installation, you can:

- Build the CLI plugin from source.
- Download a binary from the [GitHub repository](https://github.com/docker/buildx#manual-download).

> [!NOTE]
> The `docker build` command differs slightly from `docker buildx build`.
> See [build command differences](../builders/_index.md#difference-between-docker-build-and-docker-buildx-build) for details.

## BuildKit

BuildKit is the daemon process that executes build workloads.

A build execution starts when you run a `docker build` command. Buildx interprets
your command and sends a build request to BuildKit. The build request includes:

- The Dockerfile
- Build arguments
- Export options
- Caching options

BuildKit resolves the build instructions and executes the build steps. While BuildKit
executes the build, Buildx monitors the build status and prints progress to the terminal.

If the build needs resources from the client, such as local files or build secrets,
BuildKit requests only the resources it needs from Buildx.

BuildKit is more efficient than the legacy builder used in earlier Docker versions.
BuildKit requests resources only when needed. The legacy builder always copies the
local filesystem.

Examples of resources that BuildKit can request from Buildx include:

- Local filesystem build contexts
- Build secrets
- SSH sockets
- Registry authentication tokens

For more information about BuildKit, see [BuildKit](/manuals/build/buildkit/_index.md).
