---
title: Docker Compose v5 release notes
linkTitle: Release notes
weight: 10
description: Learn about the new features, bug fixes, and breaking changes for the newest Docker Compose release
keywords: release notes, compose
tags: [Release notes]
toc_max: 2
aliases:
- /release-notes/docker-compose/
- /compose/release-notes/
---

<!-- vale off -->

For more detailed information, see the [release notes in the Compose repo](https://github.com/docker/compose/releases/).

## 5.0.1

{{< release-date date="2025-12-18" >}}

### Bug fixes and enhancements

- Restored support for `COMPOSE_COMPATIBILITY`.
- Fixed grammatical errors and improve clarity in code..
- Fixed broken `run --quiet`.
- Fixed SDK example.
- Added a check buildx version is set before comparing.
- Fixed grammar: pluralize 'service' and remove apostrophes.
- Fixed progress UI not restoring terminal once operation.
- Fixed status alignment in progress UI.
- Restored image layer download progress details on pull.
- Added 'configured' event at the end of model configuration phase.
- Introduced a build tag to select watcher implementation.
- Removed mention for v2 on README.
- Fixed missing error handling in setEnvWithDotEnv.
- Adopted morikuni/aec library over raw ANSI sequences.
- Prevented incorrect progress metrics to break compose.
- Restored support for `BUILDKIT_PROGRESS`.
- Added check model plugin is successfully loaded.
- Added a warning when no service has been selected to build.

### Update

- Dependencies upgrade: bump golang 1.24.11
- Dependencies upgrade: bump osxcross
- Dependencies upgrade: bump golang.org/x/sys to 0.39.0
- Dependencies upgrade: bump github.com/docker/cli-docs-tool to 0.11.0
- Dependencies upgrade: bump golang.org/x/sync to 0.19.0
- Dependencies upgrade: bump tags.cncf.io/container-device-interface to 1.1.0
- Dependencies upgrade: bump github.com/moby/buildkit to 0.26.3

## 5.0.0

{{< release-date date="2025-12-02" >}}

Major changes in this release:

- Compose v5 introduces a new official Go SDK. This SDK provides a comprehensive API that lets you integrate Compose functionality directly into your applications, allowing you to load, validate, and manage multi-container environments without relying on the Compose CLI. For more information, see the [Compose SDK docs](/manuals/compose/compose-sdk.md).

> [!NOTE]
>
> To avoid confusion with the legacy Compose file formats labeled “v2” and “v3,” the versioning was advanced directly to v5.

### Bug fixes and enhancements

- Dropped support for internal buildkit builder.
- Introduced Compose SDK using functional parameters.
- Introduced abstractions to support SDK usage without requiring Docker CLI.
- Made progress Writer a configurable CLI component.
- Moved progress package to cmd as a CLI component.
- Added load project function to SDK.
- Introduced SDK docs.
- Documented support for OCI and Git remote resources.
- Added support for run hooks on restart.
- Fixed check for existing image to only select the target service in run command
- Introduced `--insecure-registry` flag, reserved for testing purpose.
- Added support for `build.no_cache_filter`.
- Added `--wait` option to the `docker compose start` command.
- Fixed OCI Compose override support.
- Fixed help output for `exec --no-tty` option.
- Fixed a typo in an error message.
- Disabled progress UI when build is ran with `--print`.
- Restores support for detach keys.
- Fixed images command to display the time the image was created.
- Fixed  support for port range.
- Fixed support for includes in publish command.
- Ignore runtime flags in model configuration.

### Update

- Dependencies upgrade: bump compose go to v2.10.0
- Dependencies upgrade: bump docker 28.5.2
- Dependencies upgrade: bump containerd to 2.2.0
- Dependencies upgrade: bump docker/cli to 28.5.2
- Dependencies upgrade: bump buildx v0.30.0, buildkit v0.26.0, otel v1.38.0, otel/contrib v0.63.0
- Dependencies upgrade: bump golang.org/x/sys 0.38.0
- Dependencies upgrade: bump golang.org/x/sync to 0.18.0
- Dependencies upgrade: bump github.com/hashicorp/go-version to 1.8.0
- Dependencies upgrade: bump golang.org/x/crypto v0.45.0
- Dockerfile: update golangci-lint to v2.6.2
