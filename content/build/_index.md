---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
aliases:
  - /build/buildx/
  - /buildx/working-with-buildx/
  - /develop/develop-images/build_enhancements/
---

Docker Build is one of Docker Engine's most used features. Whenever you are
creating an image you are using Docker Build. Build is a key part of your
software development life cycle allowing you to package and bundle your code and
ship it anywhere.

The Docker Engine uses a client-server architecture and is composed of multiple components
and tools. The most common method of executing a build is by issuing a
[`docker build` command](../engine/reference/commandline/build.md). The CLI
sends the request to Docker Engine which, in turn, executes your build.

There are now two components in Engine that can be used to build an image.
Starting with the [18.09 release](../engine/release-notes/18.09.md#18090),
Engine is shipped with Moby [BuildKit](buildkit/index.md), the new component for
executing your builds by default.

The new client [Docker Buildx](https://github.com/docker/buildx){:target="blank" rel="noopener" class=""},
is a CLI plugin that extends the `docker` command with the full support of the
features provided by [BuildKit](buildkit/index.md) builder toolkit.
[`docker buildx build` command](../engine/reference/commandline/buildx_build.md)
provides the same user experience as `docker build` with many new features like
creating scoped [builder instances](drivers/index.md), building against
multiple nodes concurrently, outputs configuration, inline
[build caching](cache/index.md), and specifying target platform. In
addition, Buildx also supports new features that aren't yet available for
regular `docker build` like building manifest lists, distributed caching, and
exporting build results to OCI image tarballs.

Docker Build is more than a simple build command, and it's not only about
packaging your code. It's a whole ecosystem of tools and features that support
not only common workflow tasks but also provides support for more complex and
advanced scenarios.

<div class="component-container">
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/building/packaging/">
           <img src="/assets/images/build-packaging-software.svg" alt="Closed cardboard box" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/building/packaging/">Packaging your software</a></h2>
        <p>
          Build and package your application to run it anywhere: locally or in the cloud.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/building/multi-stage/">
           <img src="/assets/images/build-multi-stage.svg" alt="Staircase" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/building/multi-stage/">Multi-stage builds</a></h2>
        <p>
          Keep your images small and secure with minimal dependencies.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/building/multi-platform/">
           <img src="/assets/images/build-multi-platform.svg" alt="Stacked windows" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/building/multi-platform/">Multi-platform images</a></h2>
        <p>
          Build, push, pull, and run images seamlessly on different computer architectures.
        </p>
      </div>
    </div>
  </div>
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/drivers/">
           <img src="/assets/images/build-drivers.svg" alt="Silhouette of an engineer, with cogwheels in the background" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/drivers/">Build drivers</a></h2>
        <p>
          Configure where and how you run your builds.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/cache/">
           <img src="/assets/images/build-cache.svg" alt="Two arrows rotating in a circle" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/cache/">Build caching</a></h2>
        <p>
          Avoid unnecessary repetitions of costly operations, such as package installs.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/ci/">
           <img src="/assets/images/build-ci.svg" alt="Infinity loop" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/ci/">Continuous integration</a></h2>
        <p>
          Learn how to use Docker in your continuous integration pipelines.
        </p>
      </div>
    </div>
  </div>
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/exporters/">
           <img src="/assets/images/build-exporters.svg" alt="Arrow coming out of a box" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/exporters/">Exporters</a></h2>
        <p>
          Export any artifact you like, not just Docker images.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/bake/">
           <img src="/assets/images/build-bake.svg" alt="Cake silhouette" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/bake/">Bake</a></h2>
        <p>
          Orchestrate your builds with Bake.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/buildkit/dockerfile-frontend/">
           <img src="/assets/images/build-frontends.svg" alt="Pen writing on a document" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/buildkit/dockerfile-frontend/">Dockerfile frontend</a></h2>
        <p>
          Learn about the Dockerfile frontend for BuildKit.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/build/buildkit/configure/">
           <img src="/assets/images/build-configure-buildkit.svg" alt="Hammer and screwdriver" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/build/buildkit/configure/">Configure BuildKit</a></h2>
        <p>
          Take a deep dive into the internals of BuildKit to get the most out of
          your builds.
        </p>
      </div>
    </div>
  </div>
</div>
