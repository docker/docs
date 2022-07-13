---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
---

Docker Build is one of Docker Engine’s most used features. Whenever you are creating an image you are using Docker Build. Build is a key part of your software development life cycle allowing you to package and bundle your code and ship it anywhere.

Engine uses a client-server architecture and is composed of multiple components and tools. The most common method of executing a build is by issuing a docker build command from the Docker CLI. The CLI sends the request to Docker Engine which, in turn, executes your build.

There are now two components in Engine that can be used to create the build. Starting with the 18.09 release, Engine is shipped with BuildKit, a new component for executing your builds and the one that is now used by default.

The previous component, which we are calling the Legacy Builder, still exists in Engine to cover some functionality not yet supported by BuildKit.
With BuildKit, a new client, Docker Buildx, becomes available as a CLI plugin. 
Docker Buildx extends the docker build command - namely through the additional docker buildx build command - and fully supports the new features BuildKit offers. 

BuildKit is the backend evolution from the Legacy Builder, it comes with new and much improved functionality that can be powerful tools in helping improve performance of your builds, reusability of your Dockerfiles and introduces support for complex scenarios.
