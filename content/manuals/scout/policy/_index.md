---
title: Policy Evaluation
linkTitle: Policy Evaluation
weight: 70
keywords: scout, supply chain, vulnerabilities, packages, cves, policy
description: |
  Policy Evaluation in Docker Scout lets you define supply chain rules for your
  artifacts and evaluate image compliance
---

Policy Evaluation in Docker Scout lets you define supply chain rules for your
artifacts and evaluate compliance using the `docker scout policy` command. Run
evaluations locally, in CI pipelines, with custom Rego policies, or using OCI
bundles. See [Evaluate policies](./local.md).

## How Policy Evaluation works

When you run `docker scout policy`, the CLI indexes the image into an SBOM and
enriches it with CVE and VEX data. It then evaluates each configured policy
in-process against that data. No data is sent to the Scout service, and an
organization is not required for most use cases.

A policy defines image quality criteria your artifacts should meet. For
example, the **No copyleft licenses** policy flags any image containing
packages distributed under a copyleft license. If an image contains such a
package, it's non-compliant with that policy.

## Policy types

Docker Scout includes the following built-in policy types:

- [Severity-Based Vulnerability](#severity-based-vulnerability)
- [Compliant Licenses](#compliant-licenses)
- [Up-to-Date Base Images](#up-to-date-base-images)
- [High-Profile Vulnerabilities](#high-profile-vulnerabilities)
- [Supply Chain Attestations](#supply-chain-attestations)
- [Default Non-Root User](#default-non-root-user)
- [Approved Base Images](#approved-base-images)
- [Valid Docker Hardened Image (DHI) or DHI base image](#valid-docker-hardened-image-dhi-or-dhi-base-image)

For configuration options for each policy type, see
[Evaluate policies](./local.md#configure-built-in-policies).

<!-- vale Docker.HeadingSentenceCase = NO -->

### Severity-Based Vulnerability

The **Severity-Based Vulnerability** policy type checks whether your artifacts
are exposed to known vulnerabilities. By default, it flags critical and high
severity vulnerabilities where a fix version is available.

Configurable parameters include severity levels, a grace period for newly
disclosed CVEs, fixable-only filtering, and package type filtering.

### Compliant Licenses

The **Compliant Licenses** policy type checks whether your images contain
packages distributed under an inappropriate license. You can configure the
list of licenses to flag and add package-level exceptions.

### Up-to-Date Base Images

The **Up-to-Date Base Images** policy type checks whether the base images you
use are current. Images are non-compliant if the tag you built from points to
a different digest than what you're using.

Your images need provenance attestations for this policy to evaluate
successfully. For more information, see [No base image data](#no-base-image-data).

### High-Profile Vulnerabilities

The **High-Profile Vulnerabilities** policy type checks whether your images
contain vulnerabilities from Docker Scout's curated list of widely recognized,
high-impact CVEs. The list includes Log4Shell, Spring4Shell, XZ backdoor, and
others, and is updated as new high-profile vulnerabilities are disclosed.

You can configure which CVEs are considered high-profile and enable tracking
of CISA's Known Exploited Vulnerabilities catalog.

### Supply Chain Attestations

The **Supply Chain Attestations** policy type checks whether your images have
[SBOM](/manuals/build/metadata/attestations/sbom.md) and
[provenance](/manuals/build/metadata/attestations/slsa-provenance.md)
attestations. Images are non-compliant if they lack either attestation type.

To ensure compliance, build with attestations:

```console
$ docker buildx build --provenance=true --sbom=true -t <IMAGE> --push .
```

### Default Non-Root User

The **Default Non-Root User** policy type detects images configured to run as
the `root` user. Use the
[`USER`](/reference/dockerfile.md#user) Dockerfile instruction to set a
non-root default user for the runtime stage.

### Approved Base Images

The **Approved Base Images** policy type ensures the base images you use match
a configurable allowlist of glob patterns. Images are non-compliant if the
base image reference doesn't match any of the allowed patterns.

Your images need provenance attestations for this policy to evaluate
successfully. For more information, see [No base image data](#no-base-image-data).

### Valid Docker Hardened Image (DHI) or DHI base image

The **Valid Docker Hardened Image (DHI) or DHI base image** policy type
ensures your images are either Docker Hardened Images or are built using a DHI
as the base image. It validates compliance by checking for a valid Docker
signed verification summary attestation.

<!-- vale Docker.HeadingSentenceCase = YES -->

## No base image data

The **Up-to-Date Base Images** and **Approved Base Images** policies require
provenance attestations to determine the base image used in your build. Without
them, these policies report **No data**.

To ensure Docker Scout always has base image information, attach provenance
attestations at build time:

```console
$ docker buildx build --provenance=true -t <IMAGE> --push .
```

## Policies page in the Dashboard

> [!IMPORTANT]
>
> The `docker scout policy` command brings policy evaluation directly to your
> CLI so you can evaluate any image locally, in CI, or with custom policies
> without needing the Dashboard. The Policies page in the Dashboard is
> deprecated and will be retired on September 1, 2026. See
> [Evaluate policies](./local.md).

The Docker Scout Dashboard previously provided a visual interface for tracking
policy compliance across your organization's images. See
[Use the Policies page in the Dashboard](./dashboard.md).
