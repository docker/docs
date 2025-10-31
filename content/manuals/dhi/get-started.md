---
linktitle: Quickstart
title: Docker Hardened Images quickstart
description: Follow a quickstart guide to explore, mirror, and run a Docker Hardened Image.
weight: 2
keywords: docker hardened images quickstart, mirror container image, run secure image
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

This guide shows you how to go from zero to running a Docker Hardened Image
(DHI) using a real example. While the steps use a specific image as an
example, they can be applied to any DHI.

> [!TIP]
>
> You can keep using the same tools and workflows you already know when moving
> to DHI from other images on Docker Hub, such as Bitnami public catalog images.
> Note that [Bitnami announced](https://github.com/bitnami/charts/issues/35164)
> that its public catalog images will no longer be available after September 29,
> 2025.
>
> In most cases, migrating is as simple as updating the image reference in your
> configuration or commands. Start with this guide, then see the [migration
> guide](./how-to/migrate.md) for more details and examples.

## Step 1: Start a free trial to access DHI

You can browse the Docker Hardened Images catalog without a subscription, but to
use an image, you must either [contact sales to
subscribe](https://www.docker.com/products/hardened-images/#getstarted) or start
a free trial for an [organization](/admin/organization/). This guide walks you
through starting a free trial.

To start a free trial:

1. Go to the Hardened Images catalog in [Docker
   Hub](https://hub.docker.com/hardened-images/catalog) and sign in.
2. Select **Start trial** and follow the on-screen instructions.

## Step 2: Find an image to use

1. Go to the Hardened Images catalog in [Docker
   Hub](https://hub.docker.com/hardened-images/catalog) and sign in.
2. In the left sidebar, choose your organization that has DHI access.
3. In the left sidebar, select **Hardened Images** > **Catalog**.

   ![Docker Hub sidebar showing DHI catalog](./images/dhi-catalog.png)

4. Use the search bar or filters to find an image (e.g., `python`, `node`,
   `golang`). For this guide, use the Python image as an example.

    ![DHI catalog with Python repository shown](./images/dhi-python-search.png)

5. Select the Python repository to view its details.

Continue to the next step to mirror the image. To dive deeper into exploring
images see [Explore Docker Hardened Images](./how-to/explore.md).

## Step 3: Mirror the image

To use a Docker Hardened Image, you must mirror it to your organization. Only
organization owners can perform this action. Mirroring creates a copy of the
image in your organization's namespace, allowing team members to pull and use
it.

1. In the image repository page, select **Mirror to repository**.
   
   ![An image of the Python page with the Mirror to repository button showing](./images/dhi-mirror-button.png)

   > [!NOTE]
   >
   > If you don't see the **Mirror to repository** button, the repository may
   > already be mirrored to your organization. In this case, you can select
   > **View in repository** to see the mirrored image's location or mirror it to
   > another repository.
  
2. Follow the on-screen instructions to mirror the repository.

It may take a few minutes for all the tags to finish mirroring. Once
mirrored, the image repository appears in your organization's namespace. For
example, in [Docker Hub](https://hub.docker.com), go to **My Hub** > ***YOUR_ORG*** > **Repositories**,
and you should see `dhi-python` listed. You can now pull it
like any other image.

![Repository list with mirrored repository showing](./images/dhi-python-mirror.png)

Continue to the next step to pull and run the image. To dive deeper into
mirroring images see [Mirror a Docker Hardened Image
repository](./how-to/mirror.md).

## Step 4: Pull and run the image

Once you've mirrored the image to your organization, you can pull and run it
like any other Docker image. Note that Docker Hardened Images are designed to be
minimal and secure, so they may not include all the tools or libraries you
expect in a typical image. You can view the typical differences in
[Considerations when adopting
DHIs](./how-to/use.md#considerations-when-adopting-dhis).

The following example demonstrates that you can run the Python image and execute
a simple Python command just like you would with any other Docker image:

1. Pull the mirrored image. Open a terminal and run the following command,
   replacing `<your-namespace>` with your organization's namespace:

   ```console
   $ docker pull <your-namespace>/dhi-python:3.13
   ```

2. Run the image to confirm everything works:

    ```console
    $ docker run --rm <your-namespace>/dhi-python:3.13 python -c "print('Hello from DHI')"
    ```
  
    This starts a container from the `dhi-python:3.13` image and runs a simple
    Python script that prints `Hello from DHI`.

To dive deeper into using images, see:

- [Use a Docker Hardened Image](./how-to/use.md) for general usage
- [Use in Kubernetes](./how-to/k8s.md) for Kubernetes deployments
- [Use a Helm chart](./how-to/helm.md) for deploying with Helm

## Step 5: Compare with the other images

You can quickly compare DHIs with other images to see the security
improvements and differences. This comparison helps you understand the value of
using hardened images.

Run the following command to see a summary comparison, replacing
`<your-namespace>` with your organization's namespace:

```console
$ docker scout compare <your-namespace>/dhi-python:3.13 --to python:3.13 --platform linux/amd64 --ignore-unchanged 2>/dev/null | sed -n '/## Overview/,/^  ## /p' | head -n -1
```

Example output:

```plaintext
  ## Overview

                      │                    Analyzed Image                     │               Comparison Image
  ────────────────────┼───────────────────────────────────────────────────────┼───────────────────────────────────────────────
    Target            │  docker/dhi-python:3.13                               │  python:3.13
      digest          │  c215e9da9f84                                         │  7f48e892134c
      tag             │  3.13                                                 │  3.13
      platform        │ linux/amd64                                           │ linux/amd64
      provenance      │ https://github.com/docker-hardened-images/definitions │ https://github.com/docker-library/python.git
                      │  77a629b3d0db035700206c2a4e7ed904e5902ea8             │  3f2d7e4c339ab883455b81a873519f1d0f2cd80a
      vulnerabilities │    0C     0H     0M     0L                            │    0C     1H     5M   141L     2?
                      │           -1     -5   -141     -2                     │
      size            │ 35 MB (-377 MB)                                       │ 412 MB
      packages        │ 80 (-530)                                             │ 610
                      │                                                       │
```

This comparison shows that the Docker Hardened Image:

- Removes vulnerabilities: 1 HIGH, 5 MEDIUM, 141 LOW, and 2 unspecified severity CVEs removed
- Reduces size: From 412 MB down to 35 MB (91% reduction)
- Minimizes packages: From 610 packages down to 80 (87% reduction)

To dive deeper into comparing images see [Compare Docker Hardened Images](./how-to/compare.md).

## What's next

You've pulled and run your first Docker Hardened Image. Here are a few ways to keep going:

- [Migrate existing applications to DHIs](./how-to/migrate.md): Learn how to
  update your Dockerfiles to use Docker Hardened Images as the base.

- [Verify DHIs](./how-to/verify.md): Use tools like [Docker Scout](/scout/) or
  Cosign to inspect and verify signed attestations, like SBOMs and provenance.

- [Scan DHIs](./how-to/scan.md): Analyze the image with Docker
  Scout or other scanners to identify known CVEs.