---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
---


Docker Build is one of Docker Engine's most used features. Whenever you are creating an image you are using Docker Build.

Engine uses a [client-server architecture](../get-started/overview.md#docker-architecture) and is composed of multiple components and tools. 
The most common method to execute a build is by issuing a [`docker build` command](../engine/reference/commandline/build.md) in the Docker CLI. The CLI sends the request to Docker Engine which, in turn, execute your build. 

There are now two components in Engine that can be used to create the build. Starting with the 18.09 release, Engine is shipped with [BuildKit](https://github.com/moby/buildkit), a new component for executing your builds.

The previous component, which we are calling the Legacy Builder, still exists in Engine to cover some functionality not yet supported by BuildKit. 
BuildKit is the backend evolution from the Legacy Builder and which provides new and much improved functionality. 

With BuildKit, a new client, [Docker Buildx](../buildx/working-with-buildx.md), becomes available as a CLI plugin that fully supports the new features BuildKit offers that. This extends the `docker build` command, namely through the additional `docker buildx build` command. 

