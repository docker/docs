---
title: Dockerfile reference release notes
description: Release notes for Dockerfile reference
keywords: build, dockerfile, reference, release notes
toc_max: 2
---

This page contains information about the new features, improvements, known
issues, and bug fixes in [Dockerfile reference](index.md).

## Dockerfile 1.4.2

* Release date: **2022-05-06**
* Usage: [`docker/dockerfile:1.4.2`](index.md#syntax)

### Notable changes

* Fix loading certain environment variables from an image passed with built
  context

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.4.2){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.4.1

* Release date: **2022-04-08**
* Usage: [`docker/dockerfile:1.4.1`](index.md#syntax)

### Notable changes

* Fix named context resolution for cross-compilation cases from input when input
  is built for a different platform

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.4.1){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.4.0

* Release date: **2022-03-09**
* Usage: [`docker/dockerfile:1.4.0`](index.md#syntax)

### Notable changes

* [`COPY --link` and `ADD --link`](index.md#copy---link) allow copying
  files with increased cache efficiency and rebase images without requiring them
  to be rebuilt. `--link` copies files to a separate layer and then uses new LLB
  MergeOp implementation to chain independent layers together
* [Heredocs](index.md#here-documents) support have been promoted from labs
  channel to stable. This feature allows writing multiline inline scripts and
  files
* Additional [named build contexts](https://docs.docker.com/engine/reference/commandline/buildx_build/#build-context)
  can be passed to build to add or overwrite a stage or an image inside the
  build. A source for the context can be a local source, image, Git, or HTTP URL
* When using a cross-compilation stage, the target platform for a step is now
  seen on progress output
* [`BUILDKIT_SANDBOX_HOSTNAME` build-arg](index.md#buildkit-built-in-build-args)
  can be used to set the default hostname for the `RUN` steps
* Fix some cases where Heredocs incorrectly removed quotes from content

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.4.0){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.3.1

* Release date: **2022-10-04**
* Usage: [`docker/dockerfile:1.3.1`](index.md#syntax)

### Notable changes

* Fix parsing "required" mount key without a value

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.3.1){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.3.0 (labs)

* Release date: **2021-07-16**
* Usage: [`docker/dockerfile:1.3.0-labs`](index.md#syntax)

> **Note**
>
> The "labs" channel provides early access to Dockerfile features that are not
> yet available in the stable channel.

### Notable changes

* `RUN` and `COPY` commands now support [Here-document syntax](index.md#here-documents)
  allowing writing multiline inline scripts and files

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.3.0-labs){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.3.0

* Release date: **2021-07-16**
* Usage: [`docker/dockerfile:1.3.0`](index.md#syntax)

### Notable changes

* `RUN` command allows [`--network` flag](index.md#run---network) for
  requesting a specific type of network conditions. `--network=host` requires
  allowing `network.host` entitlement. This feature was previously only available on labs channel
* `ADD` command with a remote URL input now correctly handles the `--chmod` flag
* Values for [`RUN --mount` flag](index.md#run---mount) now support variable
  expansion, except for the `from` field
* Allow [`BUILDKIT_MULTI_PLATFORM` build arg](index.md#buildkit-built-in-build-args)
  to force always creating multi-platform image, even if only contains single
  platform

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.3.0){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.2.1 (labs)

* Release date: **2020-12-12**
* Usage: [`docker/dockerfile:1.2.1-labs`](index.md#syntax)

> **Note**
>
> The "labs" channel provides early access to Dockerfile features that are not
> yet available in the stable channel.

### Notable changes

* `RUN` command allows [`--network` flag](index.md#run---network) for
  requesting a specific type of network conditions. `--network=host` requires
  allowing `network.host` entitlement

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.2.1-labs){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.2.1

* Release date: **2020-12-12**
* Usage: [`docker/dockerfile:1.2.1`](index.md#syntax)

### Notable changes

* Revert "Ensure ENTRYPOINT command has at least one argument"
* Optimize processing `COPY` calls on multi-platform cross-compilation builds

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.2.1){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.2.0 (labs)

* Release date: **2020-12-03**
* Usage: [`docker/dockerfile:1.2.0-labs`](index.md#syntax)

> **Note**
>
> The "labs" channel provides early access to Dockerfile features that are not
> yet available in the stable channel.

### Notable changes

* Experimental channel has been renamed to *labs*

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.2.0-labs){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.2.0

* Release date: **2020-12-03**
* Usage: [`docker/dockerfile:1.2.0`](index.md#syntax)

### Notable changes

* [`RUN --mount` syntax](index.md#run---mount) for creating secret, ssh,
  bind, and cache mounts have been moved to mainline channel
* Metadata load errors are now handled as fatal to avoid incorrect build results
* [`ARG` command](index.md#arg) now supports defining multiple build args
  on the same line similarly to `ENV`
* `--chown` flag in `ADD` now allows parameter expansion
* Allow lowercase Dockerfile name
* `ENTRYPOINT` requires at least one argument to avoid creating broken images

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.2.0){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.1.7

* Release date: **2020-04-18**
* Usage: [`docker/dockerfile:1.1.7`](index.md#syntax)

### Notable changes

* Forward `FrontendInputs` to the gateway

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.1.7){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.1.2 (experimental)

* Release date: **2019-07-31**
* Usage: [`docker/dockerfile-upstream:1.1.2-experimental`](index.md#syntax)

> **Note**
>
> The "experimental" channel provides early access to Dockerfile features that
> are not yet available in the stable channel.

### Notable changes

* Allow setting security mode for a process with `RUN --security=sandbox|insecure`
* Allow setting uid/gid for [cache mounts](index.md#run---mounttypecache)
* Avoid requesting internally linked paths to be pulled to build context
* Ensure missing cache IDs default to target paths
* Allow setting namespace for cache mounts with [`BUILDKIT_CACHE_MOUNT_NS` build arg](index.md#buildkit-built-in-build-args)

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.1.2-experimental){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.1.2

* Release date: **2019-07-31**
* Usage: [`docker/dockerfile:1.1.1`](index.md#syntax)

### Notable changes

* Fix workdir creation with correct user and don't reset custom ownership
* Fix handling empty build args also used as `ENV`
* Detect circular dependencies

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.1.2){:target="_blank" rel="noopener" class="_"} for more details.

## Dockerfile 1.1.0

* Release date: **2019-04-27**
* Usage: [`docker/dockerfile:1.1.0`](index.md#syntax)

### Notable changes

* `ADD/COPY` commands now support implementation based on `llb.FileOp` and do
  not require helper image if builtin file operations support is available
* To find the files ignored from the build context Dockerfile frontend will
  first look for a file `<path/to/Dockerfile>.dockerignore` and if it is not
  found `.dockerignore` file will be looked up from the root of the build
  context. This allows projects with multiple Dockerfiles to use different
  `.dockerignore` definitions
* `--chown` flag for `COPY` command now supports variable expansion

> See [full release notes](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.1.0){:target="_blank" rel="noopener" class="_"} for more details.
