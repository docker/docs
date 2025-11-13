---
title: Use a Docker Hardened Image chart
linktitle: Use a Helm chart
description: Learn how to use a Docker Hardened Image chart.
keywords: use hardened image, helm, k8s, kubernetes, dhi chart, chart
weight: 32
params:
  sidebar:
    badge:
      color: violet
      text: Early Access
---

{{< summary-bar feature_name="Docker Hardened Image charts" >}}

Docker Hardened Image (DHI) charts are Docker-provided [Helm charts](https://helm.sh/docs/) built from upstream and
community-maintained sources, designed for compatibility with Docker Hardened Images. These charts are available as OCI
artifacts within the DHI catalog on Docker Hub. For more details, see [Docker Hardened Image
charts](/dhi/features/helm/).

DHI charts incorporate multiple layers of supply chain security that aren't present in upstream charts:

- SLSA Level 3 compliance: Each chart is built with SLSA Build Level 3 standards, including detailed build provenance
- Software Bill of Materials (SBOMs): Comprehensive SBOMs detail all components referenced within the chart
- Cryptographic signing: All associated metadata is cryptographically signed by Docker for integrity and authenticity
- Hardened configuration: Charts automatically reference Docker Hardened Images for secure deployments
- Tested compatibility: Charts are robustly tested to work out-of-the-box with Docker Hardened Images

This guide walks you through how to use the DHI Redis chart. You can adapt the steps to other DHI charts and your own
Kubernetes workflows. DHI charts work like any other Helm chart, but you must mirror them to your own repository before
using them.

## Prerequisites

To follow along with this guide, you need:

- A Kubernetes cluster set up and [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed. To
  test locally, you can use Docker Desktop with Kubernetes enabled. For more information, see [Install Docker
  Desktop](/desktop/install/windows-install/) and [Enable Kubernetes](/desktop/use-desktop/kubernetes/).
- Helm installed. For more information, see the [Helm installation guide](https://helm.sh/docs/intro/install/).
- Access to DHI. For more information about starting a free trial, see [Get started with Docker Hardened
  Images](/dhi/get-started/).

## Step 1: Find a Docker Helm chart and request access

To find a Docker Helm chart for DHI:

1. Go to the Hardened Images catalog in [Docker Hub](https://hub.docker.com/hardened-images/catalog) and sign in.
2. In the left sidebar, select your organization that has DHI access.
3. In the left sidebar, select **Hardened Images** > **Catalog**.
4. In the search bar, search for a Helm chart. For this guide, search for `redis chart`.
5. Select the Helm chart to view its details. For this guide, select the **Redis HA Helm Chart**.

   You will see the **Overview** page with details about the chart.

6. If visible, select **Request access to Helm charts**.

   Before you can mirror the chart, you may need to request access for the Early Access program. If **Request access to
   Helm charts** is visible on the Helm chart repository details page, select it and wait for an email notifying you
   that the access has been granted by Docker.

## Step 2: Mirror the Docker Helm chart

You must mirror the Docker Helm chart to your own repository before using it.

To mirror the Docker Helm chart to your organization, in the Helm chart repository details page you opened in [step
1](#step-1-find-a-docker-helm-chart-and-request-access):

1. Select **Mirror Helm chart**.
2. Follow the on-screen instructions to mirror the Helm chart. For this guide, name the destination repository
   `dhi-redis-ha-chart`.

   When complete, you will see the details page for the mirrored Helm chart in your organization's namespace. On this
   page, you can verify that the necessary dependencies have also been mirrored.

3. If any dependencies are not mirrored, mirror them now. For this guide, select **Mirror image** if necessary for the
   Redis image, then follow the on-screen instructions.

You only need to mirror the Helm chart and its dependencies once. After they are mirrored, you can use them in any
Kubernetes cluster that can access your organization's namespace.

## Step 3: Optional. Mirror the Helm chart and/or its images to your own registry

By default, when you mirror a chart or image from the Docker Hardened Images catalog, the chart or image is mirrored to
your namespace in Docker Hub. If you want to then mirror to your own third-party registry, you can follow the
instructions in [How to mirror an image](/dhi/how-to/mirror/) for either the chart, the image, or both.

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

## Step 4: Create a Kubernetes secret for pulling images

You need to create a Kubernetes secret for pulling images from Docker Hub or your own registry. This is necessary
because Docker Hardened Images are in private repositories. If you mirror the images to your own registry, you still
need to create this secret if the registry requires authentication.

1. For Docker Hub, create a [personal access token (PAT)](/security/access-tokens/) using your Docker account or an
   [organization access token (OAT)](/enterprise/security/access-tokens/). Ensure the token has at least read-only
   access to the Docker Hardened Image repositories.
2. Create a secret in Kubernetes using the following command. Replace `<your-secret-name>`, `<your-username>`,
   `<your-personal-access-token>`, and `<your-email>` with your own values.

   > [!NOTE]
   >
   > You need to create this secret in each Kubernetes namespace that uses a DHI. If you've mirror your DHIs to another
   > registry, replace `docker.io` with your registry's hostname. Replace `<your-username>`, `<your-access-token>`, and
   > `<your-email>` with your own values. `<your-username>` is Docker ID if using a PAT or your organization name if
   > using an OAT. `<your-secret-name>` is a name you choose for the secret.

   ```console
   $ kubectl create secret docker-registry <your-secret-name> \
       --docker-server=docker.io \
       --docker-username=<your-username> \
       --docker-password=<your-access-token> \
       --docker-email=<your-email>
   ```

   For example:

    ```console
    $ kubectl create secret docker-registry dhi-pull-secret \
        --docker-server=docker.io \
        --docker-username=docs \
        --docker-password=dckr_pat_12345 \
        --docker-email=moby@example.com
   ```

## Step 5: Update the image references in the Helm chart

DHI charts reference images stored in private repositories. While many standard Helm charts use default image locations
that are accessible to everyone, DHI images must first be mirrored to your own Docker Hub namespace or private registry.
Since each organization will have their own unique repository location, the Helm chart must be updated to point to the
correct image locations specific to your organization's Docker Hub namespace or registry.

To do this, you can use one of the following approaches:

- Pre-rendering: Uses a values override file to set the image references before Helm renders the chart templates.
- Post-rendering: Uses a script that automatically rewrites image references after Helm renders the templates but
  before deploying to Kubernetes. The script is invoked by Helm during the `helm install` command using the
  `--post-renderer` flag, where you pass it the new image prefix as an argument.

{{< tabs group="rendering" >}} {{< tab name="Pre-rendering" >}}

Create a file named `dhi-images.yaml` file with the following:

```yaml
image:
  repository: <your-namespace>/dhi-redis
haproxy:
  image:
    repository: <your-namespace>/dhi-haproxy
sysctlImage:
  image:
    repository: <your-namespace>/dhi-busybox
configmapTest:
  image:
    repository: <your-namespace>/dhi-shellcheck
exporter:
  image:
    repository: <your-namespace>/dhi-redis-exporter
```

Replace `<your-namespace>` with your Docker Hub namespace or with your own namespace in your own registry.

For example, for the Redis chart:

```yaml
image:
  repository: docs/dhi-redis
haproxy:
  image:
    repository: docs/dhi-haproxy
sysctlImage:
  image:
    repository: docs/dhi-busybox
configmapTest:
  image:
    repository: docs/dhi-shellcheck
exporter:
  image:
    repository: docs/dhi-redis-exporter
```

{{< /tab >}} {{< tab name="Post-rendering" >}}

Create a script named `post-renderer.sh` using the following command:

```bash
cat > post-renderer.sh << 'EOF'
#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <new-prefix>" >&2
  exit 1
fi

# Replaces dhi/ or docker.io/dhi with the specified PREFIX
PREFIX="$1"
sed -E "s|(image: )\"?(docker\.io/)?dhi/|\1$PREFIX|g"
EOF
chmod +x post-renderer.sh
```

This script will replace all references to `dhi/` or `docker.io/dhi/` with the prefix you provide when running `helm
install`.

{{< /tab >}} {{< /tabs >}}

## Step 6: Install the Helm chart

1. If the chart is in a private repository, sign in to the registry using Helm:

   ```console
   $ echo "<your-access-token>" | helm registry login registry-1.docker.io --username <your-username> --password-stdin
   ```

   For example:

   ```console
   $ echo "dckr_pat_12345" | helm registry login registry-1.docker.io --username docs --password-stdin
   ```

2. Install the chart using `helm install`. The command differs slightly depending on whether you are using
   post-rendering or pre-rendering. Optionally, you can also use the `--dry-run` flag to test the installation without
   actually installing anything.

   {{< tabs group="rendering" >}} {{< tab name="Pre-rendering" >}}

   ```console
   $ helm install <release-name> oci://registry-1.docker.io/<your-namespace>/<helm-chart-repository> --version <chart-version> \
     --set "imagePullSecrets[0].name=<your-secret-name>" \
     -f dhi-images.yaml
   ```

   Replace `<your-namespace>` and `<chart-version>` accordingly. If the chart is in your own registry, replace
   `registry-1.docker.io/<your-namespace>` with your own registry and namespace. Replace `<your-secret-name>` with the
   name of the image pull secret you created earlier.

   For example, for the Redis chart:

   ```console
   $ helm install my-redis-ha oci://registry-1.docker.io/docs/dhi-redis-ha-chart --version 0.1.0 \
     --set "imagePullSecrets[0].name=dhi-pull-secret" \
     -f dhi-images.yaml
   ```

   {{< /tab >}} {{< tab name="Post-rendering" >}}

   ```console
   $ helm install <release-name> oci://registry-1.docker.io/<your-namespace>/<helm-chart-repository> --version <chart-version> \
     --set "imagePullSecrets[0].name=<your-secret-name>" \
     --post-renderer ./post-renderer.sh --post-renderer-args "<your-registry-and-repository>"
   ```

   Replace `<your-namespace>` and `<chart-version>` accordingly. If the chart is in your own registry, replace
   `registry-1.docker.io/<your-namespace>` with your own registry and namespace. Replace
   `<your-registry-and-repository>` with the registry and repository prefix you want to use for the images, for example,
   `gcr.io/my-project/dhi-`, or `your-namespace/` if you are using Docker Hub. Replace `<your-secret-name>` with the
   name of the image pull secret you created earlier.

   For example, for the Redis chart:

   ```console
   $ helm install my-redis-ha oci://registry-1.docker.io/docs/dhi-redis-ha-chart --version 0.1.0 \
     --set "imagePullSecrets[0].name=dhi-pull-secret" \
     --post-renderer ./post-renderer.sh --post-renderer-args "docs/"
   ```

   {{< /tab >}} {{< /tabs >}}

## Step 7: Verify the installation

After a few seconds all the pods should be up and running.

```console
$ kubectl get pods
NAME                                  READY   STATUS    RESTARTS   AGE
<release-name>-<chart-name>-server-0   3/3     Running   0          33s
```

For example, for the Redis chart:

```console
$ kubectl get pods
NAME                                  READY   STATUS    RESTARTS   AGE
my-redis-ha-redis-ha-chart-server-0   3/3     Running   0          33s
```

## Step 8: Uninstall the Helm chart

To uninstall the Helm chart, run:

```console
$ helm uninstall <release-name>
```

For example, for the Redis chart:

```console
$ helm uninstall my-redis-ha
```