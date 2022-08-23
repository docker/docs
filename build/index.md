---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
---

Docker Build is one of Docker Engine’s most used features. Whenever you are
creating an image you are using Docker Build. Build is a key part of your
software development life cycle allowing you to package and bundle your code
and ship it anywhere.

Engine uses a client-server architecture and is composed of multiple components
and tools. The most common method of executing a build is by issuing a
`docker build` command from the Docker CLI. The CLI sends the request to Docker
Engine which, in turn, executes your build.

There are now two components in Engine that can be used to create the build.
Starting with the 18.09 release, Engine is shipped with [Moby BuildKit](https://github.com/moby/buildkit){:target="_blank" rel="noopener" class="_"},
the new component for executing your builds by default.

With BuildKit, the new client [Docker Buildx](buildx/index.md), becomes
available as a CLI plugin. Docker Buildx extends the docker build command -
namely through the additional `docker buildx build` command - and fully
supports the new features BuildKit offers.

BuildKit is the backend evolution from the Legacy Builder, it comes with new
and much improved functionality that can be powerful tools for improving your
builds' performance or reusability of your Dockerfiles, and it also introduces
support for complex scenarios.

## Docker Build features

Docker Build is way more than your `docker build` command and is not only about packaging your code, it’s a whole ecosystem of tools and features that support you not only with common workflow tasks but also provides you with support for more complex and advanced scenarios.
Here’s an overview of all the use cases with which Build can support you:

### Building your images

* **Packaging your software**  
Bundle and package your code to run anywhere, from your local Docker Desktop, to Docker Engine and Kubernetes on the cloud.  
To get started with Build, see the [Hello Build](hellobuild.md) page.

* **Choosing a build driver**  
Run Buildx with different configurations depending on the scenario you are
working on, regardless of whether you are using your local machine or a remote
compute cluster, all from the comfort of your local working environment.
For more information on drivers, see the [drivers guide](buildx/drivers/index.md).

* **Optimizing builds with cache management**  
Improve build performance by using a persistent shared build cache to avoid repeating costly operations such as package installations, downloading files from the internet, or code build steps.

* **Creating build-once, run-anywhere with multi-platform builds**
Collaborate across platforms with one build artifact.  
See [Build multi-platform images](buildx/multiplatform-images.md).

### Automating your builds

* **Integrating with GitHub**  
Automate your image builds to run in GitHub actions using the official docker build actions. See:  
    * [GitHub Action to build and push Docker images with Buildx](https://github.com/docker/build-push-action).
    * [GitHub Action to extract metadata from Git reference and GitHub events](https://github.com/docker/metadata-action/).

* **Orchestrating builds across complex projects together**  
Connect your builds together and easily parameterize your images using buildx bake.  
See [High-level build options with Bake](bake/index.md).

### Customizing your Builds

* **Select your build output format**  
Choose from a variety of available output formats, to export any artifact you like from BuildKit, not just docker images.  
See [Set the export action for the build result](../engine/reference/commandline/buildx_build.md/#output).

* **Managing build secrets**  
Securely access protected repositories and resources at build time without leaking data into the final build or the cache.

### Extending BuildKit

* **Custom syntax on Dockerfile**  
Use experimental versions of the Dockerfile frontend, or even just bring your own to BuildKit using the power of custom frontends.  
See also the [Syntax directive](../engine/reference/builder/#syntax).

* **Configure BuildKit**  
Take a deep dive into the internal BuildKit configuration to get the most out of your builds.  
See also [`buildkitd.toml`](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md), the configuration file for `buildkitd`.
