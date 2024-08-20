---
title: Get started with Policy Evaluation in Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves, policy
description: |
  Policies in Docker Scout let you define supply chain rules and thresholds
  for your artifacts, and track how your artifacts perform against those
  requirements over time
---

{{% restricted title="Early Access" %}}
Policy Evaluation is an [Early Access](/release-lifecycle/#early-access-ea)
feature of Docker Scout.
{{% /restricted %}}

In software supply chain management, maintaining the security and reliability
of artifacts is a top priority. Policy Evaluation in Docker Scout introduces a
layer of control, on top of existing analysis capabilities. It lets you define
supply chain rules for your artifacts, and helps you track how your artifacts
perform, relative to your rules and thresholds, over time.

Learn how you can use Policy Evaluation to ensure that your artifacts align
with established best practices.

## How Policy Evaluation works

When you activate Docker Scout for a repository, images that you push are
[automatically analyzed](/scout/explore/analysis.md). The analysis gives you insights
about the composition of your images, including what packages they contain and
what vulnerabilities they're exposed to. Policy Evaluation builds on top of the
image analysis feature, interpreting the analysis results against the rules
defined by policies.

A policy defines image quality criteria that your artifacts should fulfill.
For example, the **No AGPL v3 licenses** policy flags any image containing packages distributed under the AGPL v3 license.
If an image contains such a package, that image is non-compliant with this policy.
Some policies, such as the **No AGPL v3 licenses** policy, are configurable.
Configurable policies let you adjust the criteria to better match your organization's needs.

In Docker Scout, policies are designed to help you ratchet forward your
security and supply chain stature. Where other tools focus on providing a pass
or fail status, Docker Scout policies visualizes how small, incremental changes
affect policy status, even when your artifacts don't meet the policy
requirements (yet). By tracking how the fail gap changes over time, you more
easily see whether your artifact is improving or deteriorating relative to
policy.

Policies don't necessarily have to be related to application security and
vulnerabilities. You can use policies to measure and track other aspects of
supply chain management as well, such as open-source license usage and base
image up-to-dateness.

## Out-of-the-box policies

Docker Scout ships the following out-of-the-box policies:

- [No fixable critical or high vulnerabilities](#no-fixable-critical-or-high-vulnerabilities)
- [No AGPL v3 licenses](#no-agpl-v3-licenses)
- [No outdated base images](#no-outdated-base-images)
- [No high-profile vulnerabilities](#no-high-profile-vulnerabilities)
- [Supply chain attestations](#supply-chain-attestations)
- [Default non-root user](#default-non-root-user)
- [No unapproved base images](#no-unapproved-base-images)

To give you a head start, Scout enables several policies by default for your
Scout-enabled repositories. You can customize the default configurations to
reflect internal requirements and standards. You can also disable a policy
altogether if it isn't relevant to you. For more information, see [Configure
policies](./configure.md).

There's also a set of [additional policies](#additional-policies) that can be
optionally enabled for repositories.

### No fixable critical or high vulnerabilities

The **No fixable critical or high vulnerabilities** policy requires that your
artifacts aren't exposed to known vulnerabilities where there's a fix version
available. Essentially, this means that there's an easy fix that you can deploy
for images that fail this policy: upgrade the vulnerable package to a version
containing a fix for the vulnerability.

By default, this policy only flags critical and high severity vulnerabilities.

This policy is violated if an artifact is affected by one or more critical-
or high-severity vulnerability, where a fix version is available.

You can configure the parameters of this policy by creating a custom version of the policy.
The following policy parameters are configurable in a custom version:

- **Age**: The minimum number of days since the vulnerability was first published

  The rationale for only flagging vulnerabilities of a certain minimum age is
  that newly discovered vulnerabilities shouldn't cause your evaluations to
  fail until you've had a chance to address them.

<!-- vale Vale.Spelling = NO -->
- **Severities**: Severity levels to consider (default: `Critical, High`)
<!-- vale Vale.Spelling = YES -->

- **Fixable vulnerabilities only**: Whether or not to only report
  vulnerabilities with a fix version available (enabled by default).

For more information about configuring policies, see [Configure policies](./configure.md).

### No AGPL v3 licenses

The **No AGPL v3 licenses** policy requires that your artifacts don't contain
packages distributed under an AGPLv3 license. This policy is violated if
your artifacts contain one or more packages with this license.

You can configure the list of licenses that this policy should look out for,
and add exceptions by specifying an allow-list (in the form of PURLs).
See [Configure policies](./configure.md).

### No outdated base images

The **No outdated base images** policy requires that the base images you use are
up-to-date.

It's violated when the tag you used to build your image points to a
different digest than what you're using. If there's a mismatch in digests, that
means the base image you're using is out of date.

Your images need provenance attestations for this policy to successfully
evaluate. For more information, see [No base image data](#no-base-image-data).

### No high-profile vulnerabilities

The **No high-profile vulnerabilities** policy requires that your artifacts don't
contain vulnerabilities from Docker Scout’s curated list. This list is kept
up-to-date with newly disclosed vulnerabilities that are widely recognized to
be risky.

The list includes the following vulnerabilities:

- [CVE-2014-0160 (OpenSSL Heartbleed)](https://scout.docker.com/v/CVE-2014-0160)
- [CVE-2021-44228 (Log4Shell)](https://scout.docker.com/v/CVE-2021-44228)
- [CVE-2023-38545 (cURL SOCKS5 heap buffer overflow)](https://scout.docker.com/v/CVE-2023-38545)
- [CVE-2023-44487 (HTTP/2 Rapid Reset)](https://scout.docker.com/v/CVE-2023-44487)
- [CVE-2024-3094 (XZ backdoor)](https://scout.docker.com/v/CVE-2024-3094)

You can configure the CVEs included in this list by creating a custom policy.
Custom configuration options include:

- **CVEs to avoid**: Specify the CVEs that you want to avoid in your artifacts.

  Default: `CVE-2014-0160`, `CVE-2021-44228`, `CVE-2023-38545`, `CVE-2023-44487`, `CVE-2024-3094`

- **CISA KEV**: Enable tracking of vulnerabilities from CISA's Known Exploited Vulnerabilities (KEV) catalog

  The [CISA KEV catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)
  includes vulnerabilities that are actively exploited in the wild. When enabled,
  the policy flags images that contain vulnerabilities from the CISA KEV catalog.

  Enabled by default.

For more information on policy configuration, see [Configure policies](./configure.md).

### Supply chain attestations

The **Supply chain attestations** policy requires that your artifacts have
[SBOM](../../build/metadata/attestations/sbom.md) and
[provenance](../../build/metadata/attestations/slsa-provenance.md) attestations.

This policy is violated if an artifact lacks either an SBOM attestation or a
provenance attestation with max mode. To ensure compliance,
update your build command to attach these attestations at build-time:

```console
$ docker buildx build --provenance=true --sbom=true -t <IMAGE> --push .
```

For more information about building with attestations, see
[Attestations](../../build/metadata/attestations/_index.md).

If you're using GitHub Actions to build and push your images,
learn how you can [configure the action](/build/ci/github-actions/attestations.md)
to apply SBOM and provenance attestations.

### Default non-root user

By default, containers run as the `root` superuser with full system
administration privileges inside the container, unless the Dockerfile specifies
a different default user. Running containers as a privileged user weakens their
runtime security, as it means any code that runs in the container can perform
administrative actions.

The **Default non-root user** policy detects images that are set to run as the
default `root` user. To comply with this policy, images must specify a non-root
user in the image configuration. Images violate this policy if they don't
specify a non-root default user for the runtime stage.

For non-compliant images, evaluation results show whether or not the `root`
user was set explicitly for the image. This helps you distinguish between
policy violations caused by images where the `root` user is implicit, and
images where `root` is set on purpose.

The following Dockerfile runs as `root` by default despite not being explicitly set:

```Dockerfile
FROM alpine
RUN echo "Hi"
```

Whereas in the following case, the `root` user is explicitly set:

```Dockerfile
FROM alpine
USER root
RUN echo "Hi"
```

> [!NOTE]
>
> This policy only checks for the default user of the image, as set in the
> image configuration blob. Even if you do specify a non-root default user,
> it's still possible to override the default user at runtime, for example by
> using the `--user` flag for the `docker run` command.

To make your images compliant with this policy, use the
[`USER`](../../reference/dockerfile.md#user) Dockerfile instruction to set
a default user that doesn't have root privileges for the runtime stage.

The following Dockerfile snippets shows the difference between a compliant and
non-compliant image.

{{< tabs >}}
{{< tab name="Non-compliant" >}}

```dockerfile
FROM alpine AS builder
COPY Makefile ./src /
RUN make build

FROM alpine AS runtime
COPY --from=builder bin/production /app
ENTRYPOINT ["/app/production"]
```

{{< /tab >}}
{{< tab name="Compliant" >}}

```dockerfile {hl_lines=7}
FROM alpine AS builder
COPY Makefile ./src /
RUN make build

FROM alpine AS runtime
COPY --from=builder bin/production /app
USER nonroot
ENTRYPOINT ["/app/production"]
```

{{< /tab >}}
{{< /tabs >}}

### No unapproved base images

The **No unapproved base images** policy ensures that the base images you use
in your builds are maintained and secure.

This policy checks whether the base images used in your builds match any of the
patterns specified in the policy configuration. The following table shows a few
example patterns for this policy.

| Use case                                                        | Pattern                          |
| --------------------------------------------------------------- | -------------------------------- |
| Allow all images from Docker Hub                                | `docker.io/*`                    |
| Allow all Docker Official Images                                | `docker.io/library/*`            |
| Allow images from a specific organization                       | `docker.io/orgname/*`            |
| Allow tags of a specific repository                             | `docker.io/orgname/repository:*` |
| Allow images on a registry with hostname `registry.example.com` | `registry.example.com/*`         |
| Allow slim tags of NodeJS images                                | `docker.io/library/node:*-slim`  |

An asterisk (`*`) matches up until the character that follows, or until the end
of the image reference. Note that the `docker.io` prefix is required in order
to match Docker Hub images. This is the registry hostname of Docker Hub.

This policy is configurable with the following options:

- **Approved base image sources**

  Specify the image reference patterns that you want to allow. The policy
  evaluates the base image references against these patterns.

  Default: `[*]` (any reference is an allowed base image)

- **Only supported tags**

  Allow only supported tags when using Docker Official Images.

  When this option is enabled, images using unsupported tags of official images
  as their base image trigger a policy violation. Supported tags for official
  images are listed in the **Supported tags** section of the repository
  overview on Docker Hub.

  Enabled by default.

- **Only supported OS distributions**

  Allow only Docker Official Images of supported Linux distribution versions.

  When this option is enabled, images using unsupported Linux distributions
  that have reached end of life (such as `ubuntu:18.04`) trigger a policy violation.

  Enabling this option may cause the policy to report no data
  if the operating system version cannot be determined.

  Enabled by default.

Your images need provenance attestations for this policy to successfully
evaluate. For more information, see [No base image data](#no-base-image-data).

## Additional policies

In addition to the [out-of-the-box policies](#out-of-the-box-policies) enabled
by default, Docker Scout supports the following optional policies. Before you
can enable these policies, you need to either configure the policies, or
configure the integration that the policy requires.

- [SonarQube quality gates passed](#sonarqube-quality-gates-passed)

### SonarQube quality gates passed

The **SonarQube quality gates passed** policy builds on the [SonarQube
integration](../integrations/code-quality/sonarqube.md) to assess the quality
of your source code. This policy works by ingesting the SonarQube code analysis
results into Docker Scout.

You define the criteria for this policy using SonarQube's [quality
gates](https://docs.sonarsource.com/sonarqube/latest/user-guide/quality-gates/).
SonarQube evaluates your source code against the quality gates you've defined
in SonarQube. Docker Scout surfaces the SonarQube assessment as a Docker Scout
policy.

Docker Scout uses [provenance](../../build/metadata/attestations/slsa-provenance.md)
attestations or the `org.opencontainers.image.revision` OCI annotation to link
SonarQube analysis results with container images. In addition to enabling the
SonarQube integration, you must also make sure that your images have either the
attestation or the label.

![Git commit SHA links image with SonarQube analysis](../images/scout-sq-commit-sha.webp)

Once you push an image and policy evaluation completes, the results from the
SonarQube quality gates display as a policy in the Docker Scout Dashboard, and
in the CLI.

> [!NOTE]
>
> Docker Scout can only access SonarQube analyses created after the integration
> is enabled. Docker Scout doesn't have access to historic evaluations. Trigger
> a SonarQube analysis and policy evaluation after enabling the integration to
> view the results in Docker Scout.

## No base image data

There are cases when it's not possible to determine information about the base
images used in your builds. In such cases, the **No outdated base images** and
**No unapproved base images** policies get flagged as having **No data**.

This "no data" state occurs when:

- Docker Scout doesn't know what base image tag you used
- The base image version you used has multiple tags, but not all tags are out
  of date

To make sure that Docker Scout always knows about your base image, you can
attach [provenance attestations](../../build/metadata/attestations/slsa-provenance.md)
at build-time. Docker Scout uses provenance attestations to find out the base
image version.
