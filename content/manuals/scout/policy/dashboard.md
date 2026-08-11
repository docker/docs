---
title: Use the Policies page in the Dashboard
linkTitle: Use the Dashboard
description: View policy status, configure policies, and get remediation advice using the Docker Scout Dashboard
keywords: scout, policy, dashboard, configure, remediation, status
params:
  sidebar:
    badge:
      color: gray
      text: Deprecated
aliases:
  - /scout/policy/ci/
  - /scout/policy/view/
  - /scout/policy/configure/
  - /scout/policy/remediation/
---

> [!IMPORTANT]
>
> The Policies page is deprecated and will be retired on September 1, 2026.
> The `docker scout policy` command replaces this with more options for
> evaluating policies. You can run evaluations locally, in CI, against custom
> Rego policies, or using OCI bundles. See
> [Evaluate policies](./local.md).

## View policy status

The **Overview** tab of the [Docker Scout Dashboard](https://scout.docker.com/)
displays a summary of recent changes in policy for your repositories.
This summary shows images that have seen the most change in their policy
evaluation between the most recent image and the previous image.

### Policy status per repository

The **Images** tab shows the current policy status, and recent policy trend,
for all images in the selected environment. The **Policy status** column shows:

- Number of fulfilled policies versus the total number of policies
- Recent policy trends

The policy trend, denoted by the directional arrows, indicates whether an image
is better, worse, or unchanged compared to the previous image in the same
environment.

- The green arrow pointing upwards shows the number of policies that improved.
- The red arrow pointing downwards shows the number of policies that worsened.
- The bidirectional gray arrow shows the number of policies that were unchanged.

### Detailed results

To view the full evaluation results for an image, navigate to the image tag in
the Docker Scout Dashboard and open the **Policy** tab.

For vulnerability-related policies, the details view shows the fix version when
one is available. For licensing-related policies, the list shows all packages
whose license doesn't meet the policy criteria.

## Configure policies

Some policy types are configurable. You can create customized versions with
your own parameters, disable a policy, or delete it.

> [!NOTE]
> Historic evaluation results for the default policy configuration are removed
> if you delete or customize a policy.

### Add a policy

1. Go to the [Policies page](https://scout.docker.com/reports/policy) in the Docker Scout Dashboard.
2. Select **Add policy**.
3. Locate the policy type you want to configure and select **Configure**.

   - If **Configure** is grayed out, the policy has no configurable parameters.
   - If the button reads **Integrate**, setup is required before the policy can be enabled.

4. Update the policy parameters.
5. Select **Save policy** to enable, or **Save and disable** to save without enabling.

### Edit a policy

1. Go to the [Policies page](https://scout.docker.com/reports/policy).
2. Select the policy and then select **Edit**.
3. Update the parameters and save.

### Disable a policy

Disabling a policy hides its results but doesn't delete historic data.

1. Go to the [Policies page](https://scout.docker.com/reports/policy).
2. Select the policy and then select **Disable**.

### Delete a policy

Deleting a policy removes its evaluation results.

1. Go to the [Policies page](https://scout.docker.com/reports/policy).
2. Select the policy and then select **Delete**.

To recreate a deleted policy, follow [Add a policy](#add-a-policy) and select
**Configure** on the deleted policy type.

## Remediation

Docker Scout provides remediation recommendations based on policy evaluation
results. Recommendations are available for the following policy types:

- [Up-to-Date Base Images](#up-to-date-base-images)
- [Supply Chain Attestations](#supply-chain-attestations)

To view recommendations:

1. Go to the [Policies page](https://scout.docker.com/reports/policy).
2. Select a policy.
3. Hover over an image in the list and select **View fixes**.

If more than one recommendation is available, the primary recommendation
displays as the **Recommended fix**. Additional recommendations are listed as
**Quick fixes**.

### Up-to-Date Base Images

Without provenance attestations, compliance is undeterminable. Add
[provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md)
so Docker Scout can detect the base image version you're using.

With provenance attestations, the recommended action shows how to update your
base image to the latest version and pin it to a specific digest. See
[Pin base image versions](/manuals/build/building/best-practices.md#pin-base-image-versions).

With the GitHub integration enabled, you can raise a pull request directly from
the remediation panel to update the base image version in your Dockerfile.

### Supply Chain Attestations

The **Supply Chain Attestations** policy requires SBOM and provenance
attestations. The remediation panel shows what's missing. For example, if your
image has a provenance attestation without enough information, rebuild with
[`mode=max`](/manuals/build/metadata/attestations/slsa-provenance.md#max) provenance.

## Evaluate policy compliance in CI

Adding policy evaluation to your CI pipelines helps you detect and prevent
cases where a change would cause policy compliance to worsen compared to your
baseline.

The recommended strategy involves evaluating a local image and comparing the
results to a baseline using an [environment](../integrations/environment/_index.md).
If policy compliance for the new image is worse than the baseline, the CI run
fails. If compliance is better or unchanged, the run succeeds.

The following GitHub Actions example uses the [Docker Scout GitHub Action](https://github.com/marketplace/actions/docker-scout)
to compare a pull request image against the `production` environment. The
`exit-on` input is set to `policy`, so the step fails only if policy compliance
has worsened.

> [!NOTE]
>
> Due to a limitation in Docker Engine, loading multi-platform images or images
> with attestations to the image store isn't supported. Build a single-platform
> image without attestations and load it for the policy evaluation to work.

```yaml
name: Docker

on:
  push:
    tags: ["*"]
    branches:
      - "main"
  pull_request:
    branches: ["**"]

env:
  REGISTRY: docker.io
  IMAGE_NAME: <IMAGE_NAME>
  DOCKER_ORG: <ORG>

jobs:
  build:
    permissions:
      pull-requests: write

    runs-on: ubuntu-latest
    steps:
      - name: Log into registry ${{ env.REGISTRY }}
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ secrets.REGISTRY_USER }}
          password: ${{ secrets.REGISTRY_TOKEN }}

      - name: Setup Docker buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@{{% param "metadata_action_version" %}}
        with:
          images: ${{ env.IMAGE_NAME }}

      - name: Build image
        id: build-and-push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          sbom: ${{ github.event_name != 'pull_request' }}
          provenance: ${{ github.event_name != 'pull_request' }}
          push: ${{ github.event_name != 'pull_request' }}
          load: ${{ github.event_name == 'pull_request' }}

      - name: Authenticate with Docker
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ secrets.DOCKER_USER }}
          password: ${{ secrets.DOCKER_PAT }}

      - name: Compare
        if: ${{ github.event_name == 'pull_request' }}
        uses: docker/scout-action@v1
        with:
          command: compare
          image: ${{ steps.meta.outputs.tags }}
          to-env: production
          platform: "linux/amd64"
          ignore-unchanged: true
          only-severities: critical,high
          organization: ${{ env.DOCKER_ORG }}
          exit-on: policy
```

For other CI platforms, see
[Docker Scout CI integrations](../integrations/_index.md#continuous-integration).
