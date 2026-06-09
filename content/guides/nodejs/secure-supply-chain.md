---
title: Secure your Node.js image supply chain
linkTitle: Secure your supply chain
weight: 45
keywords: node.js, node, sbom, provenance, attestations, docker scout, supply chain, security
description: Learn how to inspect, generate, and verify supply chain attestations for your Node.js container image.
---

## Prerequisites

Complete [Automate your builds with GitHub Actions](configure-github-actions.md).

## Overview

When you ship a container image, what's inside it and where it came from
matters. Supply chain attestations are signed records that answer questions
like which packages are in the image, what vulnerabilities affect them, how
the image was built, and what security checks it passed.

In this section, you'll inspect the attestations that ship with your Docker
Hardened Image base, generate your own SBOM and provenance attestations
during CI, and pin the base image by digest so your builds are reproducible.

The inspection commands in this topic are shown manually so you can see what
each one returns. In a real workflow you'd automate these checks with
[Docker Scout](/scout/), which runs the same scans on every push,
enforces policies in CI, and surfaces results in your registry and pull
requests.

## Inspect the base image attestations

Docker Hardened Images are built to SLSA Build Level 3 and ship with a set of
signed attestations covering bill-of-materials, vulnerabilities, build
provenance, and security scans. See
[DHI attestations](/manuals/dhi/core-concepts/attestations.md) for the full
list of types and how to verify their signatures with Cosign.

List all the attestations available on the Node.js DHI:

```console
$ docker scout attest list registry://dhi.io/node:24-alpine3.23-dev
```

View the SBOM:

```console
$ docker scout sbom registry://dhi.io/node:24-alpine3.23-dev
```

Check known vulnerabilities:

```console
$ docker scout cves registry://dhi.io/node:24-alpine3.23-dev
```

> [!NOTE]
>
> The `registry://` prefix forces `docker scout` to fetch the image and its
> attestations from the registry instead of reading a locally pulled copy. If
> you've already pulled or built against the base image, the local copy
> doesn't have the attached attestations, so the prefix is required to see
> them.

When you base your own image on a DHI image, these attestations stay attached to the base layer in the registry. Tools that inspect your image can follow the chain back to the DHI source.

## Generate attestations for your image

Update your GitHub Actions workflow to attach SBOM and provenance attestations to the image you push.

Edit `.github/workflows/build.yml` and update the build-and-push step:

```yaml {hl_lines="6-7"}
- name: Build and push Docker image
  uses: docker/build-push-action@{{% param "build_push_action_version" %}}
  with:
    context: .
    push: true
    sbom: true
    provenance: mode=max
    tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
```

- `sbom: true` tells BuildKit to scan the built image and attach an SBOM attestation.
- `provenance: mode=max` records detailed build provenance, including the source repository, commit, and build parameters.

The next time your workflow runs, the pushed image will carry these attestations alongside the image manifest in the registry.

## Inspect your pushed image's attestations

After your workflow pushes the image, inspect it the same way you inspected the base image:

```console
$ docker scout attest list registry://DOCKER_USERNAME/REPO_NAME:latest
$ docker scout sbom registry://DOCKER_USERNAME/REPO_NAME:latest
```

The SBOM includes packages from every layer, including those inherited from `dhi.io/node:24-alpine3.23-dev`. The provenance record references the DHI base image by digest, so consumers of your image can trace the build chain back to the DHI source.

## Pin the base image by digest

Image tags like `dhi.io/node:24-alpine3.23-dev` move over time as new patches land. For reproducible builds, pin to an immutable digest.

Look up the digest for each image:

```console
$ docker buildx imagetools inspect dhi.io/node:24-alpine3.23-dev --format "{{ .Manifest.Digest }}"
sha256:2bf01111c7dfe429362f64b3977f0cd6e63ff39023012f88487dec7e83aa26ca
$ docker buildx imagetools inspect dhi.io/node:24-alpine3.23 --format "{{ .Manifest.Digest }}"
sha256:868827fd45c6a01f7f3337ba7ff3f48ebb14da10d8cf3d347f98ded5481317a5
```

Each digest is a 64-character hex string. Update your `Dockerfile` to reference each digest on its corresponding `FROM` line:

```dockerfile
FROM dhi.io/node:24-alpine3.23-dev@sha256:2bf01111c7dfe429362f64b3977f0cd6e63ff39023012f88487dec7e83aa26ca AS dev
# ...
FROM dhi.io/node:24-alpine3.23@sha256:868827fd45c6a01f7f3337ba7ff3f48ebb14da10d8cf3d347f98ded5481317a5 AS runner
```

> [!TIP]
>
> Pinning by digest also pins you to that image's vulnerabilities. Use [Dependabot](https://docs.github.com/en/code-security/dependabot) or [Renovate](https://docs.renovatebot.com/) to automate digest updates so you get a PR when a new patched image is available, with a changelog to review before merging.

## Summary

In this section, you learned how to:

- Inspect the supply chain attestations that ship with the DHI base image, including SBOMs, CVE reports, and build provenance
- Generate SBOM and provenance attestations for your own image in CI
- Pin base images by digest for reproducible builds

Related information:

- [DHI attestations](/manuals/dhi/core-concepts/attestations.md)
- [Verify a Docker Hardened Image](/manuals/dhi/how-to/verify.md)
- [Docker Scout](/scout/)
- [Build attestations](/manuals/build/metadata/attestations/_index.md)

## Next steps

In the next section, you'll deploy your application to Kubernetes.
