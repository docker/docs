---
title: Use a Docker Hardened Image chart
linktitle: Use a Helm chart
description: Learn how to use a Docker Hardened Image chart.
keywords: use hardened image, helm, k8s, kubernetes, dhi chart, chart
weight: 32
---

Docker Hardened Image (DHI) charts are Docker-provided [Helm
charts](https://helm.sh/docs/) built from upstream sources, designed for
compatibility with Docker Hardened Images. These charts are available as OCI
artifacts within the DHI catalog on Docker Hub. For more details, see [Docker
Hardened Image charts](/dhi/features/helm/).

DHI charts incorporate multiple layers of supply chain security that aren't present in upstream charts:

- SLSA Level 3 compliance: Each chart is built with SLSA Build Level 3 standards, including detailed build provenance
- Software Bill of Materials (SBOMs): Comprehensive SBOMs detail all components referenced within the chart
- Cryptographic signing: All associated metadata is cryptographically signed by Docker for integrity and authenticity
- Hardened configuration: Charts automatically reference Docker Hardened Images for secure deployments
- Tested compatibility: Charts are robustly tested to work out-of-the-box with Docker Hardened Images

You can use a DHI chart like any other Helm chart stored in an OCI registry.
When you have a Docker Hardened Images subscription, you can also customize DHI
charts to reference customized images and mirrored repositories. The customized
chart build pipeline ensures that your customizations are built securely, use
the latest base charts, and include attestations.

## Find a Docker Helm chart

To find a Docker Helm chart for DHI:

1. Go to the Hardened Images catalog in [Docker Hub](https://hub.docker.com/hardened-images/catalog) and sign in.
2. In the left sidebar, select **Hardened Images** > **Catalog**.
3. Select **Filter by** for **Helm Charts**.
4. Select a Helm chart repository to view its details.

## Mirror a Helm chart and/or its images to a third-party registry

If you want to mirror to your own third-party registry, you can follow the
instructions in [Mirror a Docker Hardened Image repository](/dhi/how-to/mirror/) for either the
chart, the image, or both.

The same `regctl` tool that is used for mirroring container images can also be used for mirroring Helm charts, as Helm
charts are OCI artifacts.

For example:

```console
regctl image copy \
    "${SRC_CHART_REPO}:${TAG}" \
    "${DEST_REG}/${DEST_CHART_REPO}:${TAG}" \
    --referrers \
    --referrers-src "${SRC_ATT_REPO}" \
    --referrers-tgt "${DEST_REG}/${DEST_CHART_REPO}" \
    --force-recursive
```

## Create a Kubernetes secret for pulling images

You need to create a Kubernetes secret for pulling images from `dhi.io`, Docker
Hub, or your own registry. This is necessary because Docker Hardened Image
repositories require authentication. If you mirror the images to your own
registry, you still need to create this secret if the registry requires
authentication.

1. For `dhi.io` or Docker Hub, create a [personal access token
   (PAT)](/security/access-tokens/) using your Docker account or an
   [organization access token (OAT)](/enterprise/security/access-tokens/).
   Ensure the token has at least read-only access to the Docker Hardened Image
   repositories.
2. Create a secret in Kubernetes using the following command. Replace `<your-secret-name>`, `<your-username>`,
   `<your-personal-access-token>`, and `<your-email>` with your own values.

   > [!NOTE]
   >
   > You need to create this secret in each Kubernetes namespace that uses a
   > DHI. If you've mirror your DHIs to another registry, replace
   > `dhi.io` with your registry's hostname. Replace
   > `<your-username>`, `<your-access-token>`, and `<your-email>` with your own
   > values. `<your-username>` is Docker ID if using a PAT or your organization
   > name if using an OAT. `<your-secret-name>` is a name you choose for the
   > secret.

   ```console
   $ kubectl create secret docker-registry <your-secret-name> \
       --docker-server=dhi.io \
       --docker-username=<your-username> \
       --docker-password=<your-access-token> \
       --docker-email=<your-email>
   ```

   For example:

    ```console
    $ kubectl create secret docker-registry dhi-pull-secret \
        --docker-server=dhi.io \
        --docker-username=docs \
        --docker-password=dckr_pat_12345 \
        --docker-email=moby@example.com
   ```

## Install a Helm chart

To install a Helm chart from Docker Hardened Images:

1. Sign in to the registry using Helm:

   ```console
   $ echo $ACCESS_TOKEN | helm registry login dhi.io --username <your-username> --password-stdin
   ```

    Replace `<your-username>` and set `$ACCESS_TOKEN`.

2. Install the chart using `helm install`. Optionally, you can also use the `--dry-run` flag to test the installation without
   actually installing anything.


   ```console
   $ helm install <release-name> oci://dhi.io/<helm-chart-repository> --version <chart-version> \
     --set "imagePullSecrets[0].name=<your-secret-name>"
   ```

   Replace `<helm-chart-repository>` and `<chart-version>` accordingly. If the
   chart is in your own registry or another repository, replace
   `dhi.io/<helm-chart-repository>` with your own location. Replace
   `<your-secret-name>` with the name of the image pull secret created
   from [Create a Kubernetes secret for pulling images](#create-a-kubernetes-secret-for-pulling-images).

## Customize a Helm chart

You can customize Docker Hardened Image Helm charts to reference customized
images and mirrored repositories. For more details, see [Customize Docker
Hardened Images and charts](./customize.md).

## Verify a Helm chart and view its attestations

You can verify Helm charts. For more details, see [Verify Helm chart attestations](./verify.md#verify-helm-chart-attestations-with-docker-scout).