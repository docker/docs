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
