---
title: Docker Scout SBOMs
description: Use Docker Scout to extract the SBOM for your project.
keywords: scout, supply chain, sbom, software bill of material, spdx
aliases:
- /engine/sbom/
---

[Image analysis](./image-analysis.md) uses image SBOMs to understand what packages and versions an image contains.
Docker Scout uses SBOM attestations if available on the image (recommended).
If no SBOM attestation is available, Docker Scout creates one by indexing the image contents.

## View from CLI

To view the contents of the SBOM that Docker Scout generates, you can use the
`docker scout sbom` command.

```console
$ docker scout sbom [IMAGE]
```

By default, this prints the SBOM in a JSON format to stdout.
The default JSON format produced by `docker scout sbom` isn't SPDX-JSON.
To output SPDX, use the `--format spdx` flag:

```console
$ docker scout sbom --format spdx [IMAGE]
```

To generate a human-readable list, use the `--format list` flag:

```console
$ docker scout sbom --format list alpine

           Name             Version    Type
───────────────────────────────────────────────
  alpine-baselayout       3.4.3-r1     apk
  alpine-baselayout-data  3.4.3-r1     apk
  alpine-keys             2.4-r1       apk
  apk-tools               2.14.0-r2    apk
  busybox                 1.36.1-r2    apk
  busybox-binsh           1.36.1-r2    apk
  ca-certificates         20230506-r0  apk
  ca-certificates-bundle  20230506-r0  apk
  libc-dev                0.7.2-r5     apk
  libc-utils              0.7.2-r5     apk
  libcrypto3              3.1.2-r0     apk
  libssl3                 3.1.2-r0     apk
  musl                    1.2.4-r1     apk
  musl-utils              1.2.4-r1     apk
  openssl                 3.1.2-r0     apk
  pax-utils               1.3.7-r1     apk
  scanelf                 1.3.7-r1     apk
  ssl_client              1.36.1-r2    apk
  zlib                    1.2.13-r1    apk
```

For more information about the `docker scout sbom` command, refer to the [CLI
reference](../reference/cli/docker/scout/sbom.md).

## Attach as build attestation {#attest}

You can generate the SBOM and attach it to the image at build-time as an
[attestation](../build/attestations/_index.md). BuildKit provides a default
SBOM generator which is different from what Docker Scout uses.
You can configure BuildKit to use the Docker Scout SBOM generator
using the `--attest` flag for the `docker build` command.
The Docker Scout SBOM indexer provides richer results
and ensures better compatibility with the Docker Scout image analysis.

```console
$ docker build --tag <org>/<image> \
  --attest type=sbom,generator=docker/scout-sbom-indexer:latest \
  --push .
```

To build images with SBOM attestations, you must either turn on
the [containerd image store](../desktop/containerd.md) feature, or use a
`docker-container` builder together with the `--push` flag to push the image
(with attestations) directly to a registry.

## Extract to file

The command for extracting the SBOM of an image to an SPDX JSON file is
different depending on whether the image has been pushed to a registry or if
it's a local image.

### Remote image

To extract the SBOM of an image and save it to a file, you can use the `docker
buildx imagetools inspect` command. This command only works for images in a
registry.

```console
$ docker buildx imagetools inspect <image> --format "{{ json .SBOM }}" > sbom.spdx.json
```

### Local image

To extract the SPDX file for a local image, build the image with the `local`
exporter and use the `scout-sbom-indexer` SBOM generator plugin.

The following command saves the SBOM to a file at `build/sbom.spdx.json`.

```console
$ docker build --attest type=sbom,generator=docker/scout-sbom-indexer:latest \
  --output build .
```
