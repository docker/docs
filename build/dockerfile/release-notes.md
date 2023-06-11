---
title: Dockerfile release notes
description: Release notes for Dockerfile frontend
keywords: build, dockerfile, frontend, release notes
toc_max: 2
---

This page contains information about the new features, improvements, known
issues, and bug fixes in [Dockerfile reference](../../engine/reference/builder.md).

For usage, see the [Dockerfile frontend syntax](frontend.md) page.

## 1.5.2

{% include release-date.html date="2023-02-14" %}

### Bug fixes and enhancements

* Fix building from Git reference that is missing branch name but contains a
  subdir
* 386 platform image is now included in the release

## 1.5.1

{% include release-date.html date="2023-01-18" %}

### Bug fixes and enhancements

* Fix possible panic when warning conditions appear in multi-platform builds

## 1.5.0 (labs)

{% include release-date.html date="2023-01-10" %}

{% include dockerfile-labs-channel.md %}

### New

* `ADD` command now supports [`--checksum` flag](../../engine/reference/builder.md#verifying-a-remote-file-checksum-add---checksumchecksum-http-src-dest)
  to validate the contents of the remote URL contents

## 1.5.0

{% include release-date.html date="2023-01-10" %}

### New

* `ADD` command can now [import files directly from Git URLs](../../engine/reference/builder.md#adding-a-git-repository-add-git-ref-dir)

### Bug fixes and enhancements

* Named contexts now support `oci-layout://` protocol for including images from
  local OCI layout structure
* Dockerfile now supports secondary requests for listing all build targets or
  printing outline of accepted parameters for a specific build target
* Dockerfile `#syntax` directive that redirects to an external frontend image
  now allows the directive to be also set with `//` comments or JSON. The file
  may also contain a shebang header
* Named context can now be initialized with an empty scratch image
* Named contexts can now be initialized with an SSH Git URL
* Fix handling of `ONBUILD` when importing Schema1 images

## 1.4.3

{% include release-date.html date="2022-08-23" %}

### Bug fixes and enhancements

* Fix creation timestamp not getting reset when building image from
  `docker-image://` named context
* Fix passing `--platform` flag of `FROM` command when loading
  `docker-image://` named context

## 1.4.2

{% include release-date.html date="2022-05-06" %}

### Bug fixes and enhancements

* Fix loading certain environment variables from an image passed with built
  context

## 1.4.1

{% include release-date.html date="2022-04-08" %}

### Bug fixes and enhancements

* Fix named context resolution for cross-compilation cases from input when input
  is built for a different platform

## 1.4.0

{% include release-date.html date="2022-03-09" %}

### New

* [`COPY --link` and `ADD --link`](../../engine/reference/builder.md#copy---link)
  allow copying files with increased cache efficiency and rebase images without
  requiring them to be rebuilt. `--link` copies files to a separate layer and
  then uses new LLB MergeOp implementation to chain independent layers together
* [Heredocs](../../engine/reference/builder.md#here-documents) support have
  been promoted from labs channel to stable. This feature allows writing
  multiline inline scripts and files
* Additional [named build contexts](../../engine/reference/commandline/buildx_build.md#build-context)
  can be passed to build to add or overwrite a stage or an image inside the
  build. A source for the context can be a local source, image, Git, or HTTP URL
* [`BUILDKIT_SANDBOX_HOSTNAME` build-arg](../../engine/reference/builder.md#buildkit-built-in-build-args)
  can be used to set the default hostname for the `RUN` steps

### Bug fixes and enhancements

* When using a cross-compilation stage, the target platform for a step is now
  seen on progress output
* Fix some cases where Heredocs incorrectly removed quotes from content

## 1.3.1

{% include release-date.html date="2021-10-04" %}

### Bug fixes and enhancements

* Fix parsing "required" mount key without a value

## 1.3.0 (labs)

{% include release-date.html date="2021-07-16" %}

{% include dockerfile-labs-channel.md %}

### New

* `RUN` and `COPY` commands now support [Here-document syntax](../../engine/reference/builder.md#here-documents)
  allowing writing multiline inline scripts and files

## 1.3.0

{% include release-date.html date="2021-07-16" %}

### New

* `RUN` command allows [`--network` flag](../../engine/reference/builder.md#run---network)
  for requesting a specific type of network conditions. `--network=host`
  requires allowing `network.host` entitlement. This feature was previously
  only available on labs channel

### Bug fixes and enhancements

* `ADD` command with a remote URL input now correctly handles the `--chmod` flag
* Values for [`RUN --mount` flag](../../engine/reference/builder.md#run---mount)
  now support variable expansion, except for the `from` field
* Allow [`BUILDKIT_MULTI_PLATFORM` build arg](../../engine/reference/builder.md#buildkit-built-in-build-args)
  to force always creating multi-platform image, even if only contains single
  platform

## 1.2.1 (labs)

{% include release-date.html date="2020-12-12" %}

{% include dockerfile-labs-channel.md %}

### Bug fixes and enhancements

* `RUN` command allows [`--network` flag](../../engine/reference/builder.md#run---network)
  for requesting a specific type of network conditions. `--network=host`
  requires allowing `network.host` entitlement

## 1.2.1

{% include release-date.html date="2020-12-12" %}

### Bug fixes and enhancements

* Revert "Ensure ENTRYPOINT command has at least one argument"
* Optimize processing `COPY` calls on multi-platform cross-compilation builds

## 1.2.0 (labs)

{% include release-date.html date="2020-12-03" %}

{% include dockerfile-labs-channel.md %}

### Bug fixes and enhancements

* Experimental channel has been renamed to *labs*

## 1.2.0

{% include release-date.html date="2020-12-03" %}

### New

* [`RUN --mount` syntax](../../engine/reference/builder.md#run---mount) for
  creating secret, ssh, bind, and cache mounts have been moved to mainline
  channel
* [`ARG` command](../../engine/reference/builder.md#arg) now supports defining
  multiple build args on the same line similarly to `ENV`

### Bug fixes and enhancements

* Metadata load errors are now handled as fatal to avoid incorrect build results
* Allow lowercase Dockerfile name
* `--chown` flag in `ADD` now allows parameter expansion
* `ENTRYPOINT` requires at least one argument to avoid creating broken images

## 1.1.7

{% include release-date.html date="2020-04-18" %}

### Bug fixes and enhancements

* Forward `FrontendInputs` to the gateway

## 1.1.2 (experimental)

{% include release-date.html date="2019-07-31" %}

{% include dockerfile-labs-channel.md %}

### Bug fixes and enhancements

* Allow setting security mode for a process with `RUN --security=sandbox|insecure`
* Allow setting uid/gid for [cache mounts](../../engine/reference/builder.md#run---mounttypecache)
* Avoid requesting internally linked paths to be pulled to build context
* Ensure missing cache IDs default to target paths
* Allow setting namespace for cache mounts with [`BUILDKIT_CACHE_MOUNT_NS` build arg](../../engine/reference/builder.md#buildkit-built-in-build-args)

## 1.1.2

{% include release-date.html date="2019-07-31" %}

### Bug fixes and enhancements

* Fix workdir creation with correct user and don't reset custom ownership
* Fix handling empty build args also used as `ENV`
* Detect circular dependencies

## 1.1.0

{% include release-date.html date="2019-04-27" %}

### New

* `ADD/COPY` commands now support implementation based on `llb.FileOp` and do
  not require helper image if builtin file operations support is available
* `--chown` flag for `COPY` command now supports variable expansion

### Bug fixes and enhancements

* To find the files ignored from the build context Dockerfile frontend will
  first look for a file `<path/to/Dockerfile>.dockerignore` and if it is not
  found `.dockerignore` file will be looked up from the root of the build
  context. This allows projects with multiple Dockerfiles to use different
  `.dockerignore` definitions
