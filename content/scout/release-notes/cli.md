---
title: Docker Scout CLI release notes
description: Learn about the latest features of the Docker Scout CLI plugin
keywords: docker scout, release notes, changelog, cli, features, changes, delta, new, releases, github actions
---

This page contains information about the new features, improvements, known
issues, and bug fixes in the Docker Scout [CLI plugin](https://github.com/docker/scout-cli/)
and the `docker/scout-action` [GitHub Action](https://github.com/docker/scout-action).

## 1.7.0

{{< release-date date="2024-04-15" >}}

### New

- The [`docker scout push` command](/reference/cli/docker/scout/push/) is now fully available: analyze images locally and push the SBOM to Docker Scout.

### Bug fixes and enhancements

- Fix adding attestations with `docker scout attestation add` to images in private repositories
- Fix image processing for images based on the empty `scratch` base image
- A new `sbom://` protocol for Docker Scout CLI commands let you read a Docker Scout SBOM from standard input.

  ```console
  $ docker scout sbom IMAGE | docker scout qv sbom://
  ```

- Add classifier for Joomla packages

## 1.6.4

{{< release-date date="2024-03-26" >}}

### Bug fixes and enhancements

- Fix epoch handling for RPM-based Linux distributions

## 1.6.3

{{< release-date date="2024-03-22" >}}

### Bug fixes and enhancements

- Improve package detection to ignore referenced but not installed packages.

## 1.6.2

{{< release-date date="2024-03-22" >}}

### Bug fixes and enhancements

- EPSS data is now fetched via the backend, as opposed to via the CLI client.
- Fix an issue when rendering markdown output using the `sbom://` prefix.

### Removed

- The `docker scout cves --epss-date` and `docker scout cache prune --epss` flags have been removed.

## 1.6.1

{{< release-date date="2024-03-20" >}}

> **Note**
>
> This release only affects the `docker/scout-action` GitHub Action.

### New

- Add support for passing in SBOM files in SDPX or in-toto SDPX format

  ```yaml
  uses: docker/scout-action@v1
  with:
      command: cves
      image: sbom://alpine.spdx.json
  ```

- Add support for SBOM files in `syft-json` format

  ```yaml
  uses: docker/scout-action@v1
  with:
      command: cves
      image: sbom://alpine.syft.json
  ```

## 1.6.0

{{< release-date date="2024-03-19" >}}

> **Note**
>
> This release only affects the CLI plugin, not the GitHub Action

### New

- Add support for passing in SBOM files in SDPX or in-toto SDPX format

  ```console
  $ docker scout cves sbom://path/to/sbom.spdx.json
  ```

- Add support for SBOM files in `syft-json` format

  ```console
  $ docker scout cves sbom://path/to/sbom.syft.json
  ```

- Reads SBOM files from standard input

  ```console
  $ syft -o json alpine | docker scout cves sbom://
  ```

- Prioritize CVEs by EPSS score

  - `--epss` to display and prioritise the CVEs
  - `--epss-score` and `--epss-percentile` to filter by score and percentile
  - Prune cached EPSS files with `docker scout cache prune --epss`

### Bug fixes and enhancements

- Use Windows cache from WSL2
  
  When inside WSL2 with Docker Desktop running, the Docker Scout CLI plugin now
  uses the cache from Windows. That way, if an image has been indexed for
  instance by Docker Desktop there's no need anymore to re-index it on WSL2
  side.

- Fix a panic that would occur when analyzing a single-image `oci-dir` input
- Improve local attestation support with the containerd image store

## Earlier versions

Release notes for earlier versions of the Docker Scout CLI plugin are available
on [GitHub](https://github.com/docker/scout-cli/releases).
