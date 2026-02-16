---
title: Verify a Docker Hardened Image or chart
linktitle: Verify an image or chart
description: Use Docker Scout or cosign to verify signed attestations like SBOMs, provenance, and vulnerability data for Docker Hardened Images and charts.
weight: 40
keywords: verify container image, docker scout attest, cosign verify, sbom validation, signed container attestations, helm chart verification
---

Docker Hardened Images (DHI) and charts include signed attestations that verify
the build process, contents, and security posture.

Docker's public key for DHI images and charts is published at:

- https://registry.scout.docker.com/keyring/dhi/latest.pub
- https://github.com/docker-hardened-images/keyring

Docker recommends using [Docker Scout](/scout/), but you can use
[`regctl`](https://github.com/regclient/regclient) and
[`cosign`](https://docs.sigstore.dev/) to retrieve and verify attestations.
Docker Scout offers several key advantages: it understands DHI attestation
structures, automatically resolves platforms, provides human-readable summaries,
validates in one step with `--verify`, and integrates tightly with Docker's
attestation infrastructure.

> [!IMPORTANT]
>
> You must authenticate to the Docker Hardened Images registry (`dhi.io`) to
> pull images. Use your Docker ID credentials (the same username and password
> you use for Docker Hub) when signing in. If you don't have a Docker account,
> [create one](../../accounts/create-account.md) for free.
>
> Run `docker login dhi.io` to authenticate.

## Verify image attestations

> [!NOTE]
>
> Before you run `docker scout attest` commands, ensure any image that you have
> pulled locally is up to date with the remote image. You can do this by running
> `docker pull`. If you don't do this, you may see `No attestation found`.

### List available attestations

To list attestations for a mirrored DHI image:

{{< tabs group="tool" >}}
{{< tab name="Docker Scout" >}}

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://dhi.io/python:3.13` instead of `dhi.io/python:3.13`.

```console
$ docker scout attest list dhi.io/<image>:<tag>
```

This command shows all available attestations, including SBOMs, provenance, vulnerability reports, and more.

{{< /tab >}}
{{< tab name="regctl" >}}

First, authenticate to both registries. Prepare a [personal access token
(PAT)](../../security/access-tokens.md) for your user with `read only` access:

```console
$ export DOCKER_USERNAME="YOUR_DOCKER_USERNAME"
$ export DOCKER_PAT="YOUR_DOCKER_PAT"
$ export DOCKER_ORG="YOUR_DOCKER_ORG"
$ echo $DOCKER_PAT | regctl registry login -u "$DOCKER_USERNAME" --pass-stdin docker.io
$ echo $DOCKER_PAT | regctl registry login -u "$DOCKER_USERNAME" --pass-stdin registry.scout.docker.com
```

Then list attestations using the `--external` flag. DHI repositories store image
layers on `dhi.io` (or `docker.io` for mirrored images) and signed attestations
in `registry.scout.docker.com`:

```console
$ regctl artifact list docker.io/${DOCKER_ORG}/<image>:<tag> \
  --external registry.scout.docker.com/${DOCKER_ORG}/<image> \
  --platform linux/amd64
```

For example:

```console
$ regctl artifact list docker.io/${DOCKER_ORG}/dhi-node:22 \
  --external registry.scout.docker.com/${DOCKER_ORG}/dhi-node \
  --platform linux/amd64
```


{{< /tab >}}
{{< /tabs >}}

### Retrieve a specific attestation

{{< tabs group="tool" >}}
{{< tab name="Docker Scout" >}}

To retrieve a specific attestation, use the `--predicate-type` flag with the full predicate type URI:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  dhi.io/<image>:<tag>
```

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://dhi.io/python:3.13` instead of `dhi.io/python:3.13`.

For example:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  dhi.io/python:3.13
```

To retrieve only the predicate body:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --predicate \
  dhi.io/<image>:<tag>
```

{{< /tab >}}
{{< tab name="regctl" >}}

Once you've listed attestations, download the full attestation artifact using the digest from the `Name` field:

```console
$ regctl artifact get <attestation-digest> > attestation.json
```

For example, to save a SLSA provenance attestation:

```console
$ regctl artifact get registry.scout.docker.com/${DOCKER_ORG}/dhi-node@sha256:6cbf803796e281e535f2681de7cd33a1012202610322a50ee745d1bb02ac3c18 > slsa_provenance.json
```

{{< /tab >}}
{{< /tabs >}}

### Validate the attestation

{{< tabs >}}
{{< tab name="Docker Scout" >}}

To validate the attestation using Docker Scout, you can use the `--verify` flag:

```console
$ docker scout attest get dhi.io/<image>:<tag> \
   --predicate-type https://scout.docker.com/sbom/v0.1 --verify
```

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name
> with `registry://`. For example, use `registry://dhi.io/node:20.19-debian12`
> instead of `dhi.io/node:20.19-debian12`.


For example, to verify the SBOM attestation for the `dhi.io/node:20.19-debian12` image:

```console
$ docker scout attest get dhi.io/node:20.19-debian12 \
   --predicate-type https://scout.docker.com/sbom/v0.1 --verify
```

{{< /tab >}}
{{< tab name="cosign" >}}

Once you've listed the attestations and obtained the digest from the `Name` field, verify them using cosign:

```console
$ cosign verify \
  <attestation-digest-from-name-field> \
  --key https://registry.scout.docker.com/keyring/dhi/latest.pub \
  --insecure-ignore-tlog=true
```

For example:

```console
$ cosign verify \
  registry.scout.docker.com/${DOCKER_ORG}/dhi-node@sha256:6cbf803796e281e535f2681de7cd33a1012202610322a50ee745d1bb02ac3c18 \
  --key https://registry.scout.docker.com/keyring/dhi/latest.pub \
  --insecure-ignore-tlog=true
```

> [!NOTE]
>
> The `--insecure-ignore-tlog=true` flag is needed because DHI attestations
> may not be recorded in the public Rekor transparency log to protect private
> customer information. The attestation signature is still verified against
> Docker's public key.

{{< /tab >}}
{{< /tabs >}}

#### Handle missing transparency log entries

When using `--verify` with Docker Scout or `cosign verify`, you may sometimes
see an error like:

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
  dhi.io/<image>:<tag> \
  --verify --skip-tlog
```

> [!NOTE]
>
> The `--skip-tlog` flag is only available in Docker Scout CLI version 1.18.2 and
> later.
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://dhi.io/python:3.13` instead of `dhi.io/python:3.13`.


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
  dhi.io/<image>:<tag>
```

> [!NOTE]
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://dhi.io/python:3.13` instead of `dhi.io/python:3.13`.

For example:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --verify \
  dhi.io/python:3.13
```

If verification succeeds, Docker Scout prints the full `cosign verify` command.

Example output:

```console
    v SBOM obtained from attestation, 101 packages found
    v Provenance obtained from attestation
    v cosign verify ...
```

> [!IMPORTANT]
>
> When using cosign, you must first authenticate to both the DHI registry
> and the Docker Scout registry.
>
> For example:
>
> ```console
> $ docker login dhi.io
> $ docker login registry.scout.docker.com
> $ cosign verify ...
> ```

## Verify Helm chart attestations with Docker Scout

Docker Hardened Image Helm charts include the same comprehensive attestations
as container images. The verification process for charts is identical to that
for images, using the same Docker Scout CLI commands.

### List available chart attestations

To list attestations for a DHI Helm chart:

```console
$ docker scout attest list dhi.io/<chart>:<version>
```

For example, to list attestations for the external-dns chart:

```console
$ docker scout attest list dhi.io/external-dns-chart:1.20.0
```

This command shows all available chart attestations, including SBOMs, provenance, vulnerability reports, and more.

### Retrieve a specific chart attestation

To retrieve a specific attestation from a Helm chart, use the `--predicate-type` flag with the full predicate type URI:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  dhi.io/<chart>:<version>
```

For example:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  dhi.io/external-dns-chart:1.20.0
```

To retrieve only the predicate body:

```console
$ docker scout attest get \
  --predicate-type https://cyclonedx.org/bom/v1.6 \
  --predicate \
  dhi.io/<chart>:<version>
```

### Validate chart attestations with Docker Scout

To validate a chart attestation using Docker Scout, use the `--verify` flag:

```console
$ docker scout attest get dhi.io/<chart>:<version> \
   --predicate-type https://scout.docker.com/sbom/v0.1 --verify
```

For example, to verify the SBOM attestation for the external-dns chart:

```console
$ docker scout attest get dhi.io/external-dns-chart:1.20.0 \
   --predicate-type https://scout.docker.com/sbom/v0.1 --verify
```

The same `--skip-tlog` flag described in [Handle missing transparency log
entries](#handle-missing-transparency-log-entries) can also be used with chart
attestations when needed.

## Available DHI attestations

See [available
attestations](../core-concepts/attestations.md#image-attestations) for a list
of attestations available for each DHI image and [Helm chart
attestations](../core-concepts/attestations.md#helm-chart-attestations) for a
list of attestations available for each DHI chart.

## Explore attestations on Docker Hub

You can also browse attestations visually when [exploring an image
variant](./explore.md#view-image-variant-details). The **Attestations** section
lists each available attestation with its:

- Type (e.g. SBOM, VEX)
- Predicate type URI
- Digest reference for use with `cosign`

These attestations are generated and signed automatically as part of the Docker
Hardened Image or chart build process.