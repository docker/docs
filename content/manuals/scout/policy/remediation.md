---
title: Remediation with Docker Scout
description: Learn how Docker Scout can help you improve your software quality automatically, using remediation
keywords: scout, supply chain, security, remediation, automation
---

{{< summary-bar feature_name="Remediation with Docker Scout" >}}

Docker Scout helps you remediate supply chain or security issues by providing
recommendations based on policy evaluation results. Recommendations are
suggested actions you can take that improve policy compliance, or that add
metadata to images which enables Docker Scout to provide better evaluation
results and recommendations.

Docker Scout provides remediation advice for the default policies of the
following policy types:

- [Up-to-Date Base Images](#up-to-date-base-images-remediation)
- [Supply Chain Attestations](#supply-chain-attestations-remediation)

<!-- TODO(dvdksn): verify the following -->
> [!NOTE]
> Guided remediation is not supported for custom policies.

For images that violate policy, the recommendations focus on addressing
compliance issues and fixing violations. For images where Docker Scout is
unable to determine compliance, recommendations show you how to meet the
prerequisites that ensure Docker Scout can successfully evaluate the policy.

## View recommendations

Recommendations appear on the policy details pages of the Docker Scout
Dashboard. To get to this page:

1. Go to the [Policies page](https://scout.docker.com/reports/policy) in the Docker Scout Dashboard.
2. Select a policy in the list.

The policy details page groups evaluation results into two different tabs
depending on the policy status:

- Violations
- Compliance unknown

The **Violations** tab lists images that don't comply with the selected policy.
The **Compliance unknown** tab lists images that Docker Scout is unable to
determine the compliance status for. When compliance is unknown, Docker Scout
needs more information about the image.

To view recommended actions for an image, hover over one of the images in the
list to reveal a **View fixes** button.

![Remediation for policy violations](../images/remediation.png)

Select the **View fixes** button to opens the remediation side panel containing
recommended actions for your image.

If there are more than one recommendations available, the primary
recommendation displays as the **Recommended fix**. Additional recommendations
are listed as **Quick fixes**. Quick fixes are usually actions that provide a
temporary solution.

The side panel may also contain one or more help sections related to the
available recommendations.

## Up-to-Date Base Images remediation

The **Up-to-Date Base Images** policy checks whether the base image you use is
up-to-date. The recommended actions displayed in the remediation side panel
depend on how much information Docker Scout has about your image. The more
information that's available, the better the recommendations.

The following scenarios outline the different recommendations depending on the
information available about the image.

### No provenance attestations

For Docker Scout to be able to evaluate this policy, you must add [provenance
attestations](/manuals/build/metadata/attestations/slsa-provenance.md) to your image. If
your image doesn't have provenance attestations, compliance is undeterminable.

<!--
  TODO(dvdksn): no support for the following, yet

  When provenance attestations are unavailable, Docker Scout provides generic,
  best-effort recommendations in the remediation side panel. These
  recommendations estimate your base using information from image analysis
  results. The base image version is unknown, but you can manually select the
  version you use in the remediation side panel. This lets Docker Scout evaluate
  whether the base image detected in the image analysis is up-to-date with the
  version you selected.

  https://github.com/docker/docs/pull/18961#discussion_r1447186845
-->

### Provenance attestations available

With provenance attestations added, Docker Scout can correctly detect the base
image version that you're using. The version found in the attestations is
cross-referenced against the current version of the corresponding tag to
determine if it's up-to-date.

If there's a policy violation, the recommended actions show how to update your
base image version to the latest version, while also pinning the base image
version to a specific digest. For more information, see [Pin base image
versions](/manuals/build/building/best-practices.md#pin-base-image-versions).

### GitHub integration enabled

If you're hosting the source code for your image on GitHub, you can enable the
[GitHub integration](../integrations/source-code-management/github.md). This
integration enables Docker Scout to provide even more useful remediation
advice, and lets you initiate remediation for violations directly from the
Docker Scout Dashboard.

With the GitHub integration enabled, you can use the remediation side panel to
raise a pull request on the GitHub repository of the image. The pull request
automatically updates the base image version in your Dockerfile to the
up-to-date version.

This automated remediation pins your base image to a specific digest, while
helping you stay up-to-date as new versions become available. Pinning the base
image to a digest is important for reproducibility, and helps avoid unwanted
changes from making their way into your supply chain.

For more information about base image pinning, see [Pin base image
versions](/manuals/build/building/best-practices.md#pin-base-image-versions).

<!--
  TODO(dvdksn): no support for the following, yet

  Enabling the GitHub integration also allows Docker Scout to visualize the
  remediation workflow in the Docker Scout Dashboard. Each step, from the pull
  request being raised to the image being deployed to an environment, is
  displayed in the remediation sidebar when inspecting the image.

  https://github.com/docker/docs/pull/18961#discussion_r1447189475
-->

## Supply Chain Attestations remediation

The default **Supply Chain Attestations** policy requires full provenance and
SBOM attestations on images. If your image is missing an attestation, or if an
attestation doesn't contain enough information, the policy is violated.

The recommendations available in the remediation side panel helps guide you to
what action you need to take to address the issues. For example, if your image
has a provenance attestation, but the attestation doesn't contain enough
information, you're recommended to re-build your image with
[`mode=max`](/manuals/build/metadata/attestations/slsa-provenance.md#max) provenance.
