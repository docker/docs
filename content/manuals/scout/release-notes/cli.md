---
title: Docker Scout CLI release notes
linkTitle: CLI release notes
description: Learn about the latest features of the Docker Scout CLI plugin
keywords: docker scout, release notes, changelog, cli, features, changes, delta, new, releases, github actions
---

This page contains information about the new features, improvements, known
issues, and bug fixes in the Docker Scout [CLI plugin](https://github.com/docker/scout-cli/)
and the `docker/scout-action` [GitHub Action](https://github.com/docker/scout-action).

## 1.18.3

{{< release-date date="2025-08-13" >}}

### New

- Add `docker scout vex get` command to retrieve a merged VEX document from all VEX attestations.

### Bug fixes

- Minor fixes for Docker Hardened Images (DHI).

## 1.18.2

{{< release-date date="2025-07-21" >}}

### New

- Add `--skip-tlog` flag to `docker scout attest get` to skip signature verification against the transparency log.

### Enhancements

- Add predicate type human-readable names for DHI FIPS and STIG attestations.

### Bug fixes

- Do not filter CVEs that are marked with a VEX `under_investigation` statement.
- Minor fixes for Docker Hardened Images (DHI).

## 1.18.1

{{< release-date date="2025-05-26" >}}

### Bug fixes

- Fix issues with `docker scout attest list` and `docker scout attest get` for local images.

## 1.18.0

{{< release-date date="2025-05-13" >}}

### New

- Add `docker scout attest list` and `docker scout attest get` commands to list attestations.
- Add support for Docker Hardened Images (DHI) VEX documents.

## 1.16.1

{{< release-date date="2024-12-13" >}}

### Bug fixes

- Fix in-toto subject digest for the `docker scout attestation add` command.

## 1.16.0

{{< release-date date="2024-12-12" >}}

### New

- Add secret scanning to the `docker scout sbom` command.
- Add support for attestations for images from Tanzu Application Catalog.

### Enhancements

- Normalize licenses using the SPDX license list.
- Make licenses unique.
- Print platform in markdown output.
- Keep original pattern to find nested matches.
- Updates to make SPDX output spec-compliant.
- Update Go, crypto module, and Alpine dependencies.

### Bug fixes

- Fix behavior with multiple images in the `docker scout attest` command.
- Check directory existence before creating temporary file.

## 1.15.0

{{< release-date date="2024-10-31" >}}

### New

- New `--format=cyclonedx` flag for the `docker scout sbom` to output the SBOM in CycloneDX format.

### Enhancements

- Use high-to-low sort order for CVE summary.
- Support for enabling and disabling repositories that enabled by `docker scout push` or `docker scout watch`.

### Bug fixes

- Improve messaging when analyzing `oci` directories without attestations.
  Only single-platform images and multi-platform image _with attestations_ are supported.
  Multi-platform images without attestations are not supported.
- Improve classifiers and SBOM indexer:
  - Add classifier for Liquibase `lpm`.
  - Add Rakudo Star/MoarVM binary classifier.
  - Add binary classifiers for silverpeas utilities.
- Improve reading and caching of attestations with the containerd image store.

## 1.14.0

{{< release-date date="2024-09-24" >}}

### New

- Add suppression information at the CVE level in the `docker scout cves` command.

### Bug fixes

- Fix listing CVEs for dangling images, for example: `local://sha256:...`
- Fix panic when analysing a file system input, for instance with `docker scout cves fs://.`

## 1.13.0

{{< release-date date="2024-08-05" >}}

### New

- Add `--only-policy` filter option to the `docker scout quickview`, `docker scout policy` and `docker scout compare` commands.
- Add `--ignore-suppressed` filter option to `docker scout cves` and `docker scout quickview`  commands to filter out CVEs affected by [exceptions](/scout/explore/exceptions/).

### Bug fixes and enhancements

- Use conditional policy name in checks.
- Add support for detecting the version of a Go project set using linker flags,
  for example:

  ```console
  $ go build -ldflags "-X main.Version=1.2.3"
  ```

## 1.12.0

{{< release-date date="2024-07-31" >}}

### New

- Only display vulnerabilities from the base image:

  ```console {title="CLI"}
  $ docker scout cves --only-base IMAGE
  ```

  ```yaml {title="GitHub Action"}
  uses: docker/scout-action@v1
  with:
    command: cves
    image: [IMAGE]
    only-base: true
  ```

- Account for VEX in `quickview` command.

  ```console {title="CLI"}
  $ docker scout quickview IMAGE --only-vex-affected --vex-location ./path/to/my.vex.json
  ```

  ```yaml {title="GitHub Action"}
  uses: docker/scout-action@v1
  with:
    command: quickview
    image: [IMAGE]
    only-vex-affected: true
    vex-location: ./path/to/my.vex.json
  ```

- Account for VEX in `cves` command (GitHub Actions).

  ```yaml {title="GitHub Action"}
  uses: docker/scout-action@v1
  with:
    command: cves
    image: [IMAGE]
    only-vex-affected: true
    vex-location: ./path/to/my.vex.json
  ```

### Bug fixes and enhancements

- Update `github.com/docker/docker` to `v26.1.5+incompatible` to fix CVE-2024-41110.
- Update Syft to 1.10.0.

## 1.11.0

{{< release-date date="2024-07-25" >}}

### New

- Filter CVEs listed in the CISA Known Exploited Vulnerabilities catalog.

  ```console {title="CLI"}
  $ docker scout cves [IMAGE] --only-cisa-kev

  ... (cropped output) ...
  ## Packages and Vulnerabilities

  0C     1H     0M     0L  io.netty/netty-codec-http2 4.1.97.Final
  pkg:maven/io.netty/netty-codec-http2@4.1.97.Final

  ✗ HIGH CVE-2023-44487  CISA KEV  [OWASP Top Ten 2017 Category A9 - Using Components with Known Vulnerabilities]
    https://scout.docker.com/v/CVE-2023-44487
    Affected range  : <4.1.100
    Fixed version   : 4.1.100.Final
    CVSS Score      : 7.5
    CVSS Vector     : CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H
  ... (cropped output) ...
  ```

  ```yaml {title="GitHub Action"}
  uses: docker/scout-action@v1
  with:
    command: cves
    image: [IMAGE]
    only-cisa-kev: true
  ```

- Add new classifiers:
  - `spiped`
  - `swift`
  - `eclipse-mosquitto`
  - `znc`

### Bug fixes and enhancements

- Allow VEX matching when no subcomponents.
- Fix panic when attaching an invalid VEX document.
- Fix SPDX document root.
- Fix base image detection when image uses SCRATCH as the base image.

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

> [!NOTE]
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

> [!NOTE]
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
  [Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md) feature.

- Fix a panic that would occur when analyzing a single-image `oci-dir` input
- Improve local attestation support with the containerd image store

## Earlier versions

Release notes for earlier versions of the Docker Scout CLI plugin are available
on [GitHub](https://github.com/docker/scout-cli/releases).
