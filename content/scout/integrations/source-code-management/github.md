---
title: Integrate Docker Scout with GitHub
description: Integrate Docker Scout using the GitHub app to get remediation advice directly in your repositories
keywords: scout, github, integration, image analysis, supply chain, remediation, source code
---

> **Beta feature**
>
> The GitHub integration is currently in [Beta](../../../release-lifecycle.md#Beta).
{ .experimental }

The GitHub app integration for Docker Scout grants Docker Scout access to your
source code repository on GitHub. This improved visibility into how your image
gets created means Docker Scout can give you automated and contextual
remediation advice.

## How it works

When you enable the GitHub integration, Docker Scout can make a direct link
between the image analysis results and the source.

When analyzing your image, Docker Scout checks for [provenance
attestations](../../../build/attestations/slsa-provenance.md) to detect the
location of the source code repository for the image. If the source location is
found, and you've enabled the GitHub app, Docker Scout parses the Dockerfile
used to create the image.

Parsing the Dockerfile reveals the base image tag used to build the image. By
knowing the base image tags used, Docker Scout can detect whether the tag is
outdated, meaning it's been changed to a different image digest. For example,
say you're using `alpine:3.18` as your base image, and at a later point in
time, the image maintainers release a patch version for version `3.18`,
containing security fixes. The `alpine:3.18` tag you've been using becomes
out-of-date; the `alpine:3.18` you're using is no longer the latest.

When this happens, Docker Scout detects the discrepancy and surfaces it through
the [Outdated base images](../../policy/_index.md#outdated-base-images) policy.
When the GitHub integration's enabled, you'll also get automated suggestions on
how to update your base image. For more information about how Docker Scout can
help you automatically improve your supply chain conduct and security posture,
see [Remediation](../../policy/remediation.md).

## Setup

To integrate Docker Scout with your GitHub organization:

1. Go to [Integrations](https://scout.docker.com/settings/integrations/) on the
   Docker Scout Dashboard.
2. Find the **GitHub** integration under the **Source code management**
   section, and select **Integrate**.
3. On the integration page that opens, select the
   **Integrate GitHub app** button to open GitHub.
4. Select the organization that you want to integrate.
5. Select whether you want to integrate all repositories in the GitHub
   organization or a manual selection of repositories.
6. Select **Install & Authorize** to add the Docker Scout app to the
   organization.

   This redirects you back to the Docker Scout Dashboard, which lists your
   active GitHub integrations.

The GitHub integration is now active.
