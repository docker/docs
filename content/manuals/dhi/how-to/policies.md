---
title: Enforce Docker Hardened Image usage with policies
linktitle: Enforce image usage
description: Learn how to use image policies with Docker Scout for Docker Hardened Images.
weight: 50
keywords: docker scout policies, enforce image compliance, container security policy, image provenance, vulnerability policy check
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

Mirroring a Docker Hardened Image (DHI) repository automatically enables [Docker
Scout](/scout/), allowing you to start enforcing security and compliance policies for your
images without additional setup. Using Docker Scout policies, you can define and
apply rules that ensure only approved and secure images, such as those based on
DHIs, are used across your environments.

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

## Create policies for your DHI-based images

To ensure that the images you build using Docker Hardened Images remain secure,
you can create Docker Scout policies tailored to your requirements for your own
repositories. These policies help enforce security standards such as preventing
high-severity vulnerabilities, requiring up-to-date base images, or validating
the presence of key metadata.

Policies evaluate images when they are pushed to a repository, allowing you to
track compliance, get notified of deviations, and integrate policy checks into
your CI/CD pipeline.

### Example: Create a policy for DHI-based images

This example shows how to create a policy that requires all images in your
organization to use Docker Hardened Images as their base. This ensures that
your applications are built on secure, minimal, and production-ready images.

#### Step 1: Use a DHI base image in your Dockerfile

Create a Dockerfile that uses a Docker Hardened Image mirrored repository as the
base. For example:

```dockerfile
# Dockerfile
FROM ORG_NAME/dhi-python:3.13-alpine3.21

ENTRYPOINT ["python", "-c", "print('Hello from a DHI-based image')"]
```

#### Step 2: Build and push the image

Open a terminal and navigate to the directory containing your Dockerfile. Then,
build and push the image to your Docker Hub repository:

```console
$ docker build \
  --push \
  -t YOUR_ORG/my-dhi-app:v1 .
```

#### Step 3: Enable Docker Scout

To enable Docker Scout for your organization and the repository, run the
following commands in your terminal:

```console
$ docker login
$ docker scout enroll YOUR_ORG
$ docker scout repo enable --org YOUR_ORG YOUR_ORG/my-dhi-app
```

#### Step 4: Create a policy

1. Go to the [Docker Scout dashboard](https://scout.docker.com).
2. Select your organization and navigate to **Policies**.
3. Select **Add policy**.
4. Select **Configure** for **Approved Base Images Policy**.
5. Give the policy a compliant name, such as **Approved DHI Base Images**.
6. In **Approved base image sources**, delete the default item.
7. In **Approved base image sources**, add approved base image sources. For this
   example, use the wildcard (`*`) to allow all mirrored DHI repositories,
   `docker.io/ORG_NAME/dhi-*`. Replace `ORG_NAME` with your organization name.
8. Select **Save policy**.

#### Step 5: Evaluate policy compliance

1. Go to the [Docker Scout dashboard](https://scout.docker.com).
2. Select your organization and navigate to **Images**.
3. Find your image, `YOUR_ORG/my-dhi-app:v1`, and select the link in the **Compliance** column.

This shows the policy compliance results for your image, including whether it
meets the requirements of the **Approved DHI Base Images** policy.

You can now [evaluate policy compliance in your CI](/scout/policy/ci/).
