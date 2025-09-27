---
title: Docker Scout image analysis
description:
  Docker Scout image analysis provides a detailed view into the composition of
  your images and the vulnerabilities that they contain
keywords: scout, scanning, vulnerabilities, supply chain, security, analysis
aliases:
  - /scout/advanced-image-analysis/
  - /scout/image-analysis/
---

When you activate image analysis for a repository,
Docker Scout automatically analyzes new images that you push to that repository.

Image analysis extracts the Software Bill of Material (SBOM)
and other image metadata,and evaluates it against vulnerability data from
[security advisories](/manuals/scout/deep-dive/advisory-db-sources.md).

If you run image analysis as a one-off task using the CLI or Docker Desktop,
Docker Scout won't store any data about your image.
If you enable Docker Scout for your container image repositories however,
Docker Scout saves a metadata snapshot of your images after the analysis.
As new vulnerability data becomes available, Docker Scout recalibrates the analysis using the metadata snapshot, which means your security status for images is updated in real-time.
This dynamic evaluation means there's no need to re-analyze images when new CVE information is disclosed.

Docker Scout image analysis is available by default for Docker Hub repositories.
You can also integrate third-party registries and other services. To learn more,
see [Integrating Docker Scout with other systems](/manuals/scout/integrations/_index.md).

## Activate Docker Scout on a repository

Docker Personal comes with 1 Scout-enabled repository. You can upgrade your
Docker subscription if you need additional repositories.
See [Subscriptions and features](../../subscription/details.md)
to learn how many Scout-enabled
repositories come with each subscription tier.

Before you can activate image analysis on a repository in a third-party registry,
the registry must be integrated with Docker Scout for your Docker organization.
Docker Hub is integrated by default. For more information, see
See [Container registry integrations](/manuals/scout/integrations/_index.md#container-registries)

> [!NOTE]
>
> You must have the **Editor** or **Owner** role in the Docker organization to
> activate image analysis on a repository.

To activate image analysis:

1. Go to [Repository settings](https://scout.docker.com/settings/repos) in the Docker Scout Dashboard.
2. Select the repositories that you want to enable.
3. Select **Enable image analysis**.

If your repositories already contain images,
Docker Scout pulls and analyzes the latest images automatically.

## Analyze registry images

To trigger image analysis for an image in a registry, push the image to a
registry that's integrated with Docker Scout, to a repository where image
analysis is activated.

> [!NOTE]
>
> Image analysis on the Docker Scout platform has a maximum image file size
> limit of 10 GB, unless the image has an SBOM attestation.
> See [Maximum image size](#maximum-image-size).

1. Sign in with your Docker ID, either using the `docker login` command or the
   **Sign in** button in Docker Desktop.
2. Build and push the image that you want to analyze.

   ```console
   $ docker build --push --tag <org>/<image:tag> --provenance=true --sbom=true .
   ```

   Building with the `--provenance=true` and `--sbom=true` flags attaches
   [build attestations](/manuals/build/metadata/attestations/_index.md) to the image. Docker
   Scout uses attestations to provide more fine-grained analysis results.

   > [!NOTE]
   >
   > The default `docker` driver only supports build attestations if you use the
   > [containerd image store](/manuals/desktop/features/containerd.md).

3. Go to the [Images page](https://scout.docker.com/reports/images) in the Docker Scout Dashboard.

   The image appears in the list shortly after you push it to the registry.
   It may take a few minutes for the analysis results to appear.

## Analyze images locally

You can analyze local images with Docker Scout using Docker Desktop or the
`docker scout` commands for the Docker CLI.

### Docker Desktop

> [!NOTE]
>
> Docker Desktop background indexing supports images up to 10 GB in size.
> See [Maximum image size](#maximum-image-size).

To analyze an image locally using the Docker Desktop GUI:

1. Pull or build the image that you want to analyze.
2. Go to the **Images** view in the Docker Dashboard.
3. Select one of your local images in the list.

   This opens the [Image details view](./image-details-view.md), showing a
   breakdown of packages and vulnerabilities found by the Docker Scout analysis
   for the image you selected.

### CLI

The `docker scout` CLI commands provide a command line interface for using Docker
Scout from your terminal.

- `docker scout quickview`: summary of the specified image, see [Quickview](#quickview)
- `docker scout cves`: local analysis of the specified image, see [CVEs](#cves)
- `docker scout compare`: analyzes and compares two images

By default, the results are printed to standard output.
You can also export results to a file in a structured format,
such as Static Analysis Results Interchange Format (SARIF).

#### Quickview

The `docker scout quickview` command provides an overview of the
vulnerabilities found in a given image and its base image.

```console
$ docker scout quickview traefik:latest
    ✓ SBOM of image already cached, 311 packages indexed

  Your image  traefik:latest  │    0C     2H     8M     1L
  Base image  alpine:3        │    0C     0H     0M     0L
```

If your the base image is out of date, the `quickview` command also shows how
updating your base image would change the vulnerability exposure of your image.

```console
$ docker scout quickview postgres:13.1
    ✓ Pulled
    ✓ Image stored for indexing
    ✓ Indexed 187 packages

  Your image  postgres:13.1                 │   17C    32H    35M    33L
  Base image  debian:buster-slim            │    9C    14H     9M    23L
  Refreshed base image  debian:buster-slim  │    0C     1H     6M    29L
                                            │    -9    -13     -3     +6
  Updated base image  debian:stable-slim    │    0C     0H     0M    17L
                                            │    -9    -14     -9     -6
```

#### CVEs

The `docker scout cves` command gives you a complete view of all the
vulnerabilities in the image. This command supports several flags that lets you
specify more precisely which vulnerabilities you're interested in, for example,
by severity or package type:

```console
$ docker scout cves --format only-packages --only-vuln-packages \
  --only-severity critical postgres:13.1
    ✓ SBOM of image already cached, 187 packages indexed
    ✗ Detected 10 vulnerable packages with a total of 17 vulnerabilities

     Name            Version         Type        Vulnerabilities
───────────────────────────────────────────────────────────────────────────
  dpkg        1.19.7                 deb      1C     0H     0M     0L
  glibc       2.28-10                deb      4C     0H     0M     0L
  gnutls28    3.6.7-4+deb10u6        deb      2C     0H     0M     0L
  libbsd      0.9.1-2                deb      1C     0H     0M     0L
  libksba     1.3.5-2                deb      2C     0H     0M     0L
  libtasn1-6  4.13-3                 deb      1C     0H     0M     0L
  lz4         1.8.3-1                deb      1C     0H     0M     0L
  openldap    2.4.47+dfsg-3+deb10u5  deb      1C     0H     0M     0L
  openssl     1.1.1d-0+deb10u4       deb      3C     0H     0M     0L
  zlib        1:1.2.11.dfsg-1        deb      1C     0H     0M     0L
```

For more information about these commands and how to use them, refer to the CLI
reference documentation:

- [`docker scout quickview`](/reference/cli/docker/scout/quickview.md)
- [`docker scout cves`](/reference/cli/docker/scout/cves.md)

## Vulnerability severity assessment

Docker Scout assigns a severity rating to vulnerabilities based on
vulnerability data from [advisory sources](/manuals/scout/deep-dive/advisory-db-sources.md).
Advisories are ranked and prioritized depending on the type of package that's
affected by a vulnerability. For example, if a vulnerability affects an OS
package, the severity level assigned by the distribution maintainer is
prioritized.

If the preferred advisory source has assigned a severity rating to a CVE, but
not a CVSS score, Docker Scout falls back to displaying a CVSS score from
another source. The severity rating from the preferred advisory and the CVSS
score from the fallback advisory are displayed together. This means a
vulnerability can have a severity rating of `LOW` with a CVSS score of 9.8, if
the preferred advisory assigns a `LOW` rating but no CVSS score, and a fallback
advisory assigns a CVSS score of 9.8.

Vulnerabilities that haven't been assigned a CVSS score in any source are
categorized as **Unspecified** (U).

Docker Scout doesn't implement a proprietary vulnerability metrics system. All
metrics are inherited from security advisories that Docker Scout integrates
with. Advisories may use different thresholds for classifying vulnerabilities,
but most of them adhere to the CVSS v3.0 specification, which maps CVSS scores
to severity ratings according to the following table:

| CVSS score | Severity rating  |
| ---------- | ---------------- |
| 0.1 – 3.9  | **Low** (L)      |
| 4.0 – 6.9  | **Medium** (M)   |
| 7.0 – 8.9  | **High** (H)     |
| 9.0 – 10.0 | **Critical** (C) |

For more information, see [Vulnerability Metrics (NIST)](https://nvd.nist.gov/vuln-metrics/cvss).

Note that, given the advisory prioritization and fallback mechanism described
earlier, severity ratings displayed in Docker Scout may deviate from this
rating system.

## Maximum image size

Image analysis on the Docker Scout platform, and analysis triggered by background
indexing in Docker Desktop, has an image file size limit of 10 GB (uncompressed).
To analyze images larger than that:

- Attach an [SBOM attestation](/manuals/build/metadata/attestations/sbom.md) at build-time. When an image includes an SBOM attestation, Docker Scout uses it instead of generating one, so the 10 GB limit doesn’t apply.
- Alternatively, you can use the [CLI](#cli) to analyze the image locally. The 10 GB limit doesn’t apply when using the CLI. If the image includes an SBOM attestation, the CLI uses it to complete the analysis faster.

