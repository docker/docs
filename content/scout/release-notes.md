---
title: Docker Scout release notes
description: Learn about the latest features of Docker Scout
keywords: docker scout, release notes, changelog, features, changes, delta, new, releases
---

This page contains information about the new features, improvements, known
issues, and bug fixes in Docker Scout releases. These release notes cover the
Docker Scout platform, including the Dashboard. For CLI release notes, refer to
the `docker/scout-cli` [GitHub repository](https://github.com/docker/scout-cli/releases).

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1)
for what's coming next.

## 2023-10-04

This marks the General Availability (GA) release of Docker Scout.

### New

The following new features are included in this release:

- [Policy Evaluation](#policy-evaluation) (Early Access)
- [Amazon ECR integration](#amazon-ecr-integration)
- [Sysdig integration](#sysdig-integration)
- [JFrog Artifactory integration](#jfrog-artifactory-integration)

#### Policy evaluation

Policy Evaluation is an early access feature that helps you ensure software
integrity and track how your artifacts are doing over time. This release ships
with four out-of-the-box policies, enabled by default for all organizations.

![Policy overview in Dashboard](./images/release-notes/policy-ea.webp)

- **Base images not up-to-date** evaluates whether the base images are out of
  date, and require updating. Up-to-date base images help you ensure that your
  environments are reliable and secure.
- **Critical and high vulnerabilities with fixes** reports if there are
  vulnerabilities with critical or high severity in your images, and where
  there's a fix version available that you can upgrade to.
- **All critical vulnerabilities** looks out for any vulnerabilities of
  critical severity found in your images.
- **Packages with AGPLv3, GPLv3 license** helps you catch possibly unwanted
  copyleft licenses used in your images.

You can view and evaluate policy status for images using the Docker Scout
Dashboard and the `docker scout policy` CLI command. For more information,
refer to the [Policy Evaluation documentation](./policy/_index.md).

#### Amazon ECR integration

The new Amazon Elastic Container Registry (ECR) integration enables image
analysis for images hosted in ECR repositories.

You set up the integration using a pre-configured CloudFormation stack template
that bootstraps the necessary AWS resources in your account. Docker Scout
automatically analyzes images that you push to your registry, storing only the
metadata about the image contents, and not the container images themselves.

The integration offers a straightforward process for adding additional
repositories, activating Docker Scout for specific repositories, and removing
the integration if needed. To learn more, refer to the [Amazon ECR integration
documentation](./integrations/registry/ecr.md).

#### Sysdig integration

The new Sysdig integration gives you real-time security insights for your
Kubernetes runtime environments.

Enabling this integration helps you address and prioritize risks for images
used to run your production workloads. It also helps reduce monitoring noise,
by automatically excluding vulnerabilities in programs that are never loaded
into memory, using VEX documents.

For more information and getting started, see [Sysdig integration
documentation](./integrations/environment/sysdig.md).

#### JFrog Artifactory integration

The new JFrog Artifactory integration enables automatic image analysis on
Artifactory registries.

![Animation of how to integrate Artifactory](./images/release-notes/artifactory-agent.gif)

The integration involves deploying a Docker Scout Artifactory agent that polls
for new images, performs analysis, and uploads results to Docker Scout, all
while preserving the integrity of image data. Learn more in the [Artifactory
integration documentation](./integrations/registry/artifactory.md)

### Known limitations

- Image analysis only works for Linux images
- Docker Scout can't process images larger than 12GB in compressed size
- Creating an image SBOM (part of image analysis) has a timeout limit of 4 minutes
