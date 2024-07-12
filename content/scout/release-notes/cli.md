---
title: Docker Scout CLI release notes
description: Learn about the latest features of the Docker Scout CLI plugin
keywords: docker scout, release notes, changelog, cli, features, changes, delta, new, releases, github actions
---

This page contains information about the new features, improvements, known
issues, and bug fixes in the Docker Scout [CLI plugin](https://github.com/docker/scout-cli/)
and the `docker/scout-action` [GitHub Action](https://github.com/docker/scout-action).

## 1.10.0

{{< release-date date="2024-06-26" >}}

### Bug fixes and enhancements

- Add new classifiers:
  - `irssi`
  - `Backdrop`
  - `CrateDB CLI (Crash)`
  - `monica`
  - `Openliberty`
  - `dumb-init`
  - `friendica`
  - `redmine`
- Fix whitespace-only originator on package breaking BuildKit exporters
- Fix parsing image references in SPDX statement for images with a digest
- Support `sbom://` prefix for image comparison:

  ```console {title="CLI"}
  $ docker scout compare sbom://image1.json --to sbom://image2.json
  ```

  ```yaml {title="GitHub Action"}
  uses: docker/scout-action@v1
  with:
    command: compare
    image: sbom://image1.json
    to: sbom://image2.json
  ```

## 1.9.3

{{< release-date date="2024-05-28" >}}

### Bug fix

- Fix a panic while retrieving cached SBOMs.

## 1.9.1

{{< release-date date="2024-05-27" >}}

### New

- Add support for the [GitLab container scanning file format](https://docs.gitlab.com/ee/development/integrations/secure.html#container-scanning) with `--format gitlab` on `docker scout cves` command.

  Here is an example pipeline:

  ```yaml
     docker-build:
    # Use the official docker image.
    image: docker:cli
    stage: build
    services:
      - docker:dind
    variables:
      DOCKER_IMAGE_NAME: $CI_REGISTRY_IMAGE:$CI_COMMIT_REF_SLUG
    before_script:
      - docker login -u "$CI_REGISTRY_USER" -p "$CI_REGISTRY_PASSWORD" $CI_REGISTRY
         
      # Install curl and the Docker Scout CLI
      - |
        apk add --update curl
        curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s -- 
        apk del curl 
        rm -rf /var/cache/apk/* 
      # Login to Docker Hub required for Docker Scout CLI
      - echo "$DOCKER_HUB_PAT" | docker login --username "$DOCKER_HUB_USER" --password-stdin
  
    # All branches are tagged with $DOCKER_IMAGE_NAME (defaults to commit ref slug)
    # Default branch is also tagged with `latest`
    script:
      - docker buildx b --pull -t "$DOCKER_IMAGE_NAME" .
      - docker scout cves "$DOCKER_IMAGE_NAME" --format gitlab --output gl-container-scanning-report.json
      - docker push "$DOCKER_IMAGE_NAME"
      - |
        if [[ "$CI_COMMIT_BRANCH" == "$CI_DEFAULT_BRANCH" ]]; then
          docker tag "$DOCKER_IMAGE_NAME" "$CI_REGISTRY_IMAGE:latest"
          docker push "$CI_REGISTRY_IMAGE:latest"
        fi
    # Run this job in a branch where a Dockerfile exists
    rules:
      - if: $CI_COMMIT_BRANCH
        exists:
          - Dockerfile
    artifacts:
      reports:
        container_scanning: gl-container-scanning-report.json
  ```

### Bug fixes and enhancements

- Support single-architecture images for `docker scout attest add` command
- Indicate on the `docker scout quickview` and `docker scout recommendations` commands if image provenance was not created using `mode=max`.
  Without `mode=max`, base images may be incorrectly detected, resulting in less accurate results.

## 1.9.0

{{< release-date date="2024-05-24" >}}

Discarded in favor of [1.9.1](#191).

## 1.8.0

{{< release-date date="2024-04-25" >}}

### Bug fixes and enhancements

- Improve format of EPSS score and percentile.
  
  Before:

  ```text
  EPSS Score      : 0.000440
  EPSS Percentile : 0.092510
  ```

  After:

  ```text
  EPSS Score      : 0.04%
  EPSS Percentile : 9th percentile
  ```

- Fix markdown output of the `docker scout cves` command when analyzing local filesystem. [docker/scout-cli#113](https://github.com/docker/scout-cli/issues/113)

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
- Indexing is now blocked in the CLI if it has been disabled using
  [Settings Management](../../desktop/hardened-desktop/settings-management/configure.md) feature.

- Fix a panic that would occur when analyzing a single-image `oci-dir` input
- Improve local attestation support with the containerd image store

## Earlier versions

Release notes for earlier versions of the Docker Scout CLI plugin are available
on [GitHub](https://github.com/docker/scout-cli/releases).
