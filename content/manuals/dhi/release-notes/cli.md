---
title: DHI CLI release notes
linkTitle: CLI release notes
description: New features, bug fixes, and changes in the DHI CLI
keywords: docker hardened images, dhi, dhictl, cli, release notes, changelog
toc_min: 1
toc_max: 2
tags:
  - Release notes
---

This page lists changes in recent stable releases of the DHI CLI (`docker dhi`). For
the full release history, including pre-releases and downloads, see the
[dhictl releases on GitHub](https://github.com/docker-hardened-images/dhictl/releases).

<!-- BEGIN GENERATED RELEASES -->

## 0.0.13

{{< release-date date="2026-11-05" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.13)

Maintenance release with dependency updates.

## 0.0.12

{{< release-date date="2026-11-05" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.12)

Maintenance release with dependency updates.

## 0.0.11

{{< release-date date="2026-11-04" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.11)

This release adds platform validation for image customizations in `dhictl customization create`.

### What's New

- Adds platforms validation to `dhictl customization create`, requiring platforms to be specified for image customizations and checking that every target's tag definition supports them — missing or unsupported platforms now fail with an error

## 0.0.10

{{< release-date date="2026-10-29" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.10)

This release adds support for Helm chart customization variable mappings in dhictl.

### What's New

- Adds `variable_mappings` support for Helm chart customizations, allowing each chart variable to point at a DHI image, a mirrored image tag definition, or a customization

## 0.0.9

{{< release-date date="2026-10-21" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.9)

This release adds working directory support for image customizations, improves customization inheritance behavior, and enhances the helm chart `get` command output.

### What's New

- Adds `working-directory` field support for image customizations (single customizations only; bulk customizations display a "cannot be customized" comment and `create`/`edit` reject it with an error)
- Displays tag definition ID in the output of the `catalog get` command for helm charts

### Bug Fixes

- Fixes `customization prepare` to write base image values as commented-out fields instead of live settings, allowing proper inheritance from tag definitions rather than freezing values at scaffold time

## 0.0.8

{{< release-date date="2026-10-14" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.8)

This release adds symlink support to customization files in dhictl.

### What's New

- Adds a `symlinks` section to customization YAML files, enabling `create`, `edit`, and `get` commands to manage symlinks with `path`, `source`, and optional `uid`/`gid` fields

## 0.0.7

{{< release-date date="2026-07-21" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.7)

### Bug Fixes

- Fixes JSON output incorrectly escaping ampersands as `\u0026` — values such as catalog categories containing `&` now appear as-is in output

## 0.0.6

{{< release-date date="2026-07-13" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.6)

This release fixes an issue that prevented customizations from being created via the CLI.

### Bug Fixes

- Fixes `dhictl customization create` failing when GraphQL rejected mutation inputs that incorrectly included the output-only `__typename` field from OCI artifact data

## 0.0.5

{{< release-date date="2026-06-29" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.5)

Maintenance release with dependency updates.

## 0.0.4

{{< release-date date="2026-05-25" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.4)

### What's New

- Adds `deb` subcommand for DHI DEB repositories that emits netrc-style credentials for authenticating against DHI DEB repositories

## 0.0.3

{{< release-date date="2026-04-22" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.3)

### What's New

- Adds attestation list and get commands for managing attestations
- Adds SBOM subcommand for software bill of materials attestation
- Adds bulk support to prepare command for customizations
- Adds compression field support for customizations
- Adds tag-definition-id column to catalog get output

### Breaking change

We removed the `--output` flags from the few commands that had it (`customization prepare` and `customization get`) in favor of stdout redirections.
```console
# before
dhictl customization prepare --org my-org golang 1.25 --output my-customization.yaml

# after 
dhictl customization prepare --org my-org golang 1.25 > my-customization.yaml
```

## 0.0.2

{{< release-date date="2026-03-19" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.2)

This is a maintenance release focused on build system improvements.

### Technical Changes

- Disables CGO globally to fix macOS 16 dyld crash and simplify build process

## 0.0.1

{{< release-date date="2026-03-12" >}}

[GitHub release](https://github.com/docker-hardened-images/dhictl/releases/tag/v0.0.1)

This release improves the mirroring functionality in dhictl by allowing command arguments.

### Improvements

- Mirror start command now accepts arguments for more flexible mirroring operations

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[dhictl releases on GitHub](https://github.com/docker-hardened-images/dhictl/releases).
