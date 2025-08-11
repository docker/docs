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

## Step 1: Sign up and subscribe to DHI for access

To access Docker Hardened Images, your organization must [sign
up](https://www.docker.com/products/hardened-images/#getstarted) and subscribe.

## Step 2: Find an image to use

Once subscribed, Docker Hardened Images will appear under your organization's
namespace on Docker Hub.

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub** in the top navigation.
3. In the left sidebar, choose your organization that has DHI access.
4. In the left sidebar, select **Hardened Images** > **Catalog**.

   ![Docker Hub sidebar showing DHI catalog](./images/dhi-catalog.png)

5. Use the search bar or filters to find an image (e.g., `python`, `node`,
   `golang`). For this guide, use the Python image as an example.

    ![DHI catalog with Python repository shown](./images/dhi-python-search.png)

6. Select the Python repository to view its details.

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
  
2. Follow the on-screen instructions to choose a name. For this guide, the
   example uses the name `dhi-python`. Note that the name must start with
   `dhi-`.

   ![Mirror a repository page](./images/dhi-mirror-screen.png)

3. Select **Create repository** to start the mirroring process.

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

To dive deeper into using images see [Use a Docker Hardened Image](./how-to/use.md).

## What's next

You've pulled and run your first Docker Hardened Image. Here are a few ways to keep going:

- [Migrate existing applications to DHIs](./how-to/migrate.md): Learn how to
  update your Dockerfiles to use Docker Hardened Images as the base.

- [Verify DHIs](./how-to/verify.md): Use tools like [Docker Scout](/scout/) or
  Cosign to inspect and verify signed attestations, like SBOMs and provenance.

- [Scan DHIs](./how-to/scan.md): Analyze the image with Docker
  Scout or other scanners to identify known CVEs.