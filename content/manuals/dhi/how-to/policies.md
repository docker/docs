---
title: Enforce Docker Hardened Image usage with policies
linktitle: Enforce image usage
description: Learn how to use image policies with Docker Scout for Docker Hardened Images.
weight: 50
keywords: docker scout policies, enforce image compliance, container security policy, image provenance, vulnerability policy check
---

When you have a Docker Hardened Images Enterprise subscription, mirroring a
Docker Hardened Image (DHI) repository automatically enables [Docker
Scout](/scout/), allowing you to start enforcing security and compliance
policies for your images without additional setup. Using Docker Scout policies,
you can define and apply rules that ensure only approved and secure images, such
as those based on DHIs, are used across your environments.

Docker Scout includes a dedicated [**Valid Docker Hardened Image (DHI) or DHI
base
image**](../../scout/policy/_index.md#valid-docker-hardened-image-dhi-or-dhi-base-image)
policy type that validates whether your images are Docker Hardened Images or are
built using a DHI as the base image. This policy checks for valid Docker signed
verification summary attestations.

With policy evaluation built into Docker Scout, you can monitor image compliance
in real time, integrate checks into your CI/CD workflows, and maintain
consistent standards for image security and provenance.

## View existing policies

To see the current policies applied to a mirrored DHI repository:

1. Go to the mirrored DHI repository in [Docker Hub](https://hub.docker.com).
2. Select **View on Scout**.

   This opens the [Docker Scout dashboard](https://scout.docker.com), where you
   can see which policies are currently active and whether your images meet the
   policy criteria.

Docker Scout automatically evaluates policy compliance when new images are
pushed. Each policy includes a compliance result and a link to the affected
images and layers.

## Evaluate DHI policy compliance for your images

When you enable Docker Scout for your repositories, you can configure the
[**Valid Docker Hardened Image (DHI) or DHI base
image**](../../scout/policy/_index.md#valid-docker-hardened-image-dhi-or-dhi-base-image)
policy. This optional policy validates whether your images are DHIs or built with DHI
base images by checking for Docker signed verification summary attestations.

The following example shows how to build an image using a DHI base image and
evaluate its compliance with the DHI policy.

### Example: Build and evaluate a DHI-based image

#### Step 1: Use a DHI base image in your Dockerfile

Create a Dockerfile that uses a Docker Hardened Image mirrored repository as the
base. For example:

```dockerfile
# Dockerfile
FROM <your-namespace>/dhi-python:3.13-alpine3.21

ENTRYPOINT ["python", "-c", "print('Hello from a DHI-based image')"]
```

#### Step 2: Build and push the image

Open a terminal and navigate to the directory containing your Dockerfile. Then,
build and push the image to your Docker Hub repository:

```console
$ docker build \
  --push \
  -t <your-namespace>/my-dhi-app:v1 .
```

#### Step 3: Enable Docker Scout

To enable Docker Scout for your organization and the repository, run the
following commands in your terminal:

```console
$ docker login
$ docker scout enroll <your-namespace>
$ docker scout repo enable --org <your-namespace> <your-namespace>/my-dhi-app
```

#### Step 4: Configure the DHI policy

Once Docker Scout is enabled, you can configure the **Valid Docker Hardened
Image (DHI) or DHI base image** policy for your organization:

1. Go to the [Docker Scout dashboard](https://scout.docker.com).
2. Select your organization and navigate to **Policies**.
3. Configure the **Valid Docker Hardened Image (DHI) or DHI base image** policy
   to enable it for your repositories.

For more information on configuring policies, see
[Configure policies](../../scout/policy/configure.md).

#### Step 5: View policy compliance

Once the DHI policy is configured and active, you can view compliance results:

1. Go to the [Docker Scout dashboard](https://scout.docker.com).
2. Select your organization and navigate to **Images**.
3. Find your image, `<your-namespace>/my-dhi-app:v1`, and select the link in the **Compliance** column.

This shows the policy compliance results for your image. The **Valid Docker
Hardened Image (DHI) or DHI base image** policy evaluates whether your image has
a valid Docker signed verification summary attestation or if its base image has
such an attestation.

You can now [evaluate policy compliance in your CI](/scout/policy/ci/).