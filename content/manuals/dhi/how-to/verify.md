---
title: Verify Docker Hardened Images and charts
linktitle: Verify images and charts
description: Use Docker Scout or cosign to verify signed attestations like SBOMs, provenance, and vulnerability data for Docker Hardened Images and Helm charts.
weight: 40
keywords: verify container image, docker scout attest, cosign verify, sbom validation, signed container attestations, helm chart verification
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

Docker Hardened Images (DHI) and DHI Helm charts include signed attestations that verify the
build process, contents, and security posture. These attestations are available
for each image variant and chart, and can be verified using
[cosign](https://docs.sigstore.dev/) or the Docker Scout CLI.

Docker's public key for DHI images and charts is published at:

- https://registry.scout.docker.com/keyring/dhi/latest.pub
- https://github.com/docker-hardened-images/keyring

## Available DHI attestations

See [available image attestations](../core-concepts/attestations.md#image-attestations) for a list
of attestations available for each DHI, and [available Helm chart attestations](../core-concepts/attestations.md#helm-chart-attestations)
for attestations included with DHI charts.

### Explore image attestations on Docker Hub

You can also browse attestations visually when [exploring an image
variant](./explore.md#view-image-variant-details). The **Attestations** section
lists each available attestation with its:

- Type (for example, SBOM or VEX)
- Predicate type URI
- Digest reference for use with `cosign`

These attestations are generated and signed automatically as part of the Docker
Hardened Image build process.


## Verify image attestations with Docker Scout

You can use the [Docker Scout](/scout/) CLI to list and retrieve attestations for Docker
Hardened Images, including images mirrored into your organization's namespace.

> [!NOTE]
>
> Before you run `docker scout attest` commands, ensure any image that you have
> pulled locally is up to date with the remote image. You can do this by running
> `docker pull`. If you don't do this, you may see `No attestation found`.

### Why use Docker Scout instead of cosign directly?

While you can use cosign to verify attestations manually, the Docker Scout CLI
offers several key advantages when working with Docker Hardened Images:

- Purpose-built experience: Docker Scout understands the structure of DHI
  attestations and image naming conventions, so you don't have to construct full
  image digests or URIs manually.

- Automatic platform resolution: With Scout, you can specify the platform (e.g.,
  `--platform linux/amd64`), and it automatically verifies the correct image
  variant. Cosign requires you to look up the digest yourself.

- Human-readable summaries: Scout returns summaries of attestation contents
  (e.g., package counts, provenance steps), whereas cosign only returns raw
  signature validation output.

- One-step validation: The `--verify` flag in `docker scout attest get` validates
  the attestation and shows the equivalent cosign command, making it easier to
  understand what's happening behind the scenes.

- Integrated with Docker Hub and DHI trust model: Docker Scout is tightly
  integrated with Docker’s attestation infrastructure and public keyring,
  ensuring compatibility and simplifying verification for users within the
  Docker ecosystem.

In short, Docker Scout streamlines the verification process and reduces the chances of human error, while still giving
you full visibility and the option to fall back to cosign when needed.

### List available attestations

To list attestations for a mirrored DHI:

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://docs/dhi-python:3.13` instead of `docs/dhi-python:3.13`.

```console
$ docker scout attest list <your-org-namespace>/dhi-<image>:<tag>
```

This command shows all available attestations, including SBOMs, provenance, vulnerability reports, and more.

### Retrieve a specific attestation

To retrieve a specific attestation, use the `--predicate-type` flag with the full predicate type URI:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  <your-org-namespace>/dhi-<image>:<tag>
```

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://docs/dhi-python:3.13` instead of `docs/dhi-python:3.13`.

For example:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  docs/dhi-python:3.13
```

To retrieve only the predicate body:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --predicate \
  <your-org-namespace>/dhi-<image>:<tag>
```

For example:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --predicate \
  docs/dhi-python:3.13
```

### Validate the attestation with Docker Scout

To validate the attestation using Docker Scout, you can use the `--verify` flag:

```console
$ docker scout attest get <image-name>:<tag> \
   --predicate-type https://scout.docker.com/sbom/v0.1 --verify
```

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://docs/dhi-node:20.19-debian12-fips-20250701182639` instead of
> `docs/dhi-node:20.19-debian12-fips-20250701182639`.


For example, to verify the SBOM attestation for the `dhi/node:20.19-debian12-fips-20250701182639` image:

```console
$ docker scout attest get docs/dhi-node:20.19-debian12-fips-20250701182639 \
   --predicate-type https://scout.docker.com/sbom/v0.1 --verify
```

#### Handle missing transparency log entries

When using `--verify`, you may sometimes see an error like:

```text
ERROR no matching signatures: signature not found in transparency log
```

This occurs because Docker Hardened Images don't always record attestations in
the public [Rekor](https://docs.sigstore.dev/logging/overview/) transparency
log. In cases where an attestation would contain private user information (for
example, your organization's namespace in the image reference), writing it to
Rekor would expose that information publicly.

Even if the Rekor entry is missing, the attestation is still signed with
Docker's public key and can be verified offline by skipping the Rekor
transparency log check.

To skip the transparency log check and validate against Docker's key, use the
`--skip-tlog` flag:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  <your-org-namespace>/dhi-<image>:<tag> \
  --verify --skip-tlog
```

> [!NOTE]
>
> The `--skip-tlog` flag is only available in Docker Scout CLI version 1.18.2 and
> later.
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://docs/dhi-python:3.13` instead of `docs/dhi-python:3.13`.


This is equivalent to using `cosign` with the `--insecure-ignore-tlog=true`
flag, which validates the signature against Docker's published public key, but
ignores the transparency log check.

### Show the equivalent cosign command

When using the `--verify` flag, it also prints the corresponding
[cosign](https://docs.sigstore.dev/) command to verify the image signature:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --verify \
  <your-org-namespace>/dhi-<image>:<tag>
```

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://docs/dhi-python:3.13` instead of `docs/dhi-python:3.13`.

For example:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --verify \
  docs/dhi-python:3.13
```

If verification succeeds, Docker Scout prints the full `cosign verify` command.

Example output:

```console
    v SBOM obtained from attestation, 101 packages found
    v Provenance obtained from attestation
    v cosign verify registry.scout.docker.com/docker/dhi-python@sha256:b5418da893ada6272add2268573a3d5f595b5c486fb7ec58370a93217a9785ae \
        --key https://registry.scout.docker.com/keyring/dhi/latest.pub --experimental-oci11
    ...
```

> [!IMPORTANT]
>
> When using cosign, you must first authenticate to both the Docker Hub registry
> and the Docker Scout registry.
>
> For example:
>
> ```console
> $ docker login
> $ docker login registry.scout.docker.com
> $ cosign verify \
>     registry.scout.docker.com/docker/dhi-python@sha256:b5418da893ada6272add2268573a3d5f595b5c486fb7ec58370a93217a9785ae \
>     --key https://registry.scout.docker.com/keyring/dhi/latest.pub --experimental-oci11
> ```

## Verify Helm chart attestations with Docker Scout

You can use the Docker Scout CLI to list and retrieve attestations for DHI Helm
charts using the same commands as for images.

### List available chart attestations

To list all attestations attached to a DHI Helm chart:

```console
$ docker scout attestation list <your-namespace>/<chart-name>:<tag>
```

For example, to view attestations for the Redis HA chart for the `docs` organization:

```console
$ docker scout attestation list docs/dhi-redis-ha-chart:4.35.2
```

This command displays an overview of the chart including its provenance source
and a detailed list of all available attestations with their predicate types and
digest information.

### Retrieve a specific chart attestation

To retrieve the full content of a specific attestation, use:

```console
$ docker scout attestation get <your-namespace>/<chart-name>:<tag> --predicate-type <predicate-uri>
```

For example, to view the SLSA provenance:

```console
$ docker scout attestation get docs/dhi-redis-ha-chart:4.35.2 \
  --predicate-type https://slsa.dev/provenance/v0.2
```
