---
title: Docker Scout image analysis
description:
  Docker Scout image analysis provides a detailed view into the composition of
  your images and the vulnerabilities that they contain
keywords: scanning, vulnerabilities, supply chain, security, analysis
aliases:
  - /scout/advanced-image-analysis/
---

When you activate image analysis for a repository, Docker Scout analyzes new
images automatically when you push to that repository. Docker Scout image
analysis is more than point-in-time scanning, the analysis gets reevaluated
continuously, meaning you don't need to re-scan the image to see an updated
vulnerability report.

Docker Scout image analysis is available by default for Docker Hub
repositories. You can also integrate third-party registries, such as Amazon ECR
and JFrog Artifactory, and even run image analysis locally on your development
machine.

The following video shows how to activate Docker Scout image analysis on your
repositories.

<iframe class="border-0 w-full aspect-video mb-8" allow="fullscreen" src="https://www.loom.com/embed/a6fb14ede0a94d0d984edf6cf16604e0?sid=ba34f694-32a6-4b74-b3f8-9cc6b80ef66f"></iframe>

## Activate image analysis

The free tier of Docker Scout lets you use Docker Scout for up to 3
repositories per Docker organization. You can update your Docker Scout plan if
you need additional repositories, see [Docker Scout
billing](../billing/scout-billing.md).

Before you can activate image analysis for a repository, ensure that the
registry is integrated with Docker Scout. Docker Hub is integrated by default.
For information about integrating Docker Scout with registries and other
systems, see [Integrating Docker Scout](./integrations/_index.md)

> **Note**
>
> You must have the **Editor** or **Owner** role in the Docker organization to
> activate image analysis on a repository.

To activate image analysis:

1. Go to the [Docker Scout Dashboard](https://scout.docker.com/)
2. Sign in with your Docker ID.
3. Make sure that the correct Docker organization is selected.
4. Open the settings menu and select **Repository settings**.
5. Select the repositories that you want to enable.
6. Select **Enable image analysis**.

If your repositories already contain images, Docker Scout pulls and analyzes
the latest images automatically.

## Analyze registry images

To trigger image analysis for an image in a registry, push the image to a
registry that's integrated with Docker Scout, to a repository where image
analysis is activated.

1. Sign in with your Docker ID, either using the `docker login` command or the
   **Sign in** button in Docker Desktop.
2. Build and push the image that you want to analyze.

   ```console
   $ docker build --push --tag <org>/<image:tag> --provenance=true --sbom=true .
   ```

   Building with the `--provenance=true` and `--sbom=true` flags attaches
   [build attestations](../build/attestations/_index.md) to the image. Docker
   Scout uses attestations to provide more fine-grained analysis results.

   The default `docker` driver only supports build attestations if you use the
   [containerd image store](../desktop/containerd.md).

3. Go to the [Docker Scout Dashboard](https://scout.docker.com/)
4. Sign in with your Docker ID.
5. Select the Docker organization that contains the image you just pushed.
6. Go to the **Images** tab. The image appears in the list shortly after you
   push it to the registry.

   It may take a few minutes for the analysis report to appear. If the analysis
   report is not available, wait a moment and then refresh the page.

## Analyze images locally

You can analyze local images with Docker Scout using Docker Desktop or the
`docker scout quickview` and `docker scout cves` commands for the Docker CLI.

### Docker Desktop

> **Note**
>
> There is a 3 GB size limit on images analyzed by Docker Scout in Docker
> Desktop.

To analyze an image locally using the Docker Desktop GUI:

1. Pull or build the image that you want to analyze.
2. Go to the **Images** view in the Docker Dashboard.
3. Select one of your local images in the list.

   This opens the [Image details view](./image-details-view.md), showing a
   breakdown of packages and vulnerabilities found by the Docker Scout analysis
   for the image you selected.

### CLI

The `docker scout` CLI commands provide a terminal interface for using Docker
Scout with local and remote images.

Using the `docker scout quickview` and `docker scout cves` CLI commands, you
can analyze images locally and view the analysis report in text format. You can
print the results directly to stdout, or export them to a file using a
structured format, such as Static Analysis Results Interchange Format (SARIF).

#### Install

The Docker Scout CLI plugin is available in Docker Desktop starting with
version 4.17 and available as a standalone binary.

To install the latest version of the plugin manually, run the following
commands:

```console
$ curl -fsSL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh -o install-scout.sh
$ sh install-scout.sh
```

> **Note**
>
> Always examine scripts downloaded from the internet before running them
> locally. Before installing, make yourself familiar with potential risks and
> limitations of the convenience script.

If you want to install the plugin manually, you can find full instructions in
the [plugin's repository](https://github.com/docker/scout-cli).

The plugin is also available as [a container image](https://hub.docker.com/r/docker/scout-cli)
and as [a GitHub action](https://github.com/docker/scout-action).

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

- [`docker scout quickview`](../engine/reference/commandline/scout_quickview.md)
- [`docker scout cves`](../engine/reference/commandline/scout_cves.md)
