---
title: Configure Docker Cloud
linktitle: Configure
weight: 20
description: Learn how to configure build settings for Docker Cloud.
keywords: cloud, configuration, settings, cloud builder, GPU, disk allocation, private resources, firewall
---

To use Docker Cloud, you must configure your environment based on how you're using it:

- If you're using Cloud mode to build and run containers through Docker Desktop,
  follow the [Docker Cloud quickstart](/cloud/quickstart/).
- If you're not using Cloud mode and only building with Docker Cloud, follow the
  steps in [Build with Docker Cloud](/cloud/build/) or [Use Docker Cloud to build images in CI](/cloud/ci-build/).

Both ways of using Docker Cloud can be further customized through **Cloud
settings** in the Docker Cloud dashboard.

## Cloud settings

The **Cloud settings** page in Docker Cloud dashboard lets you configure Docker Cloud and
GPU access, disk allocation, private resource access, and firewall settings for
your cloud builders in your organization.

To view the **Cloud settings** page:

1. Go to [Docker Home](https://app.docker.com/).
2. Select the account for which you want to manage Docker Cloud.
3. Select **Go to Docker Cloud**
4. Select **Cloud settings**.

The following sections describe the available settings.

### Cloud feature availability

The **Allow Docker Cloud usage** option lets you control whether your
organization can use Docker Cloud features through Docker Desktop for hybrid
development.

The **Allow GPU access** option lets you control whether your organization can
utilize GPU-accelerated containers when using Docker Cloud in Docker Desktop.

### Lock Docker Cloud

The **Lock access to Docker Cloud** setting removes the ability for anybody in
your organization to utilize any cloud resources or consume build or run
minutes.

### Disk allocation

The **Disk allocation** setting lets you control how much of the available
storage is dedicated to the build cache. A lower allocation increases storage
available for active builds.

To make disk allocation changes, navigate to **Cloud settings** in Docker
Cloud and then adjust the **Disk allocation** slider to specify the
percentage of storage used for build caching.

Any changes take effect immediately.

> [!TIP]
> 
> If you build very large images, consider allocating less storage for caching.

### Build cache space

Your subscription includes the following Build cache space:

| Subscription | Build cache space |
|--------------|-------------------|
| Personal     | N/A               |
| Pro          | 50GB              |
| Team         | 100GB             |
| Business     | 200GB             |

To get more Build cache space, [upgrade your subscription](/manuals/subscription/change.md).

> [!TIP]
>
> If you build large images, consider allocating less storage for caching.

### Private resource access

Private resource access lets cloud builders pull images and packages from
private resources. This feature is useful when builds rely on self-hosted
artifact repositories or private OCI registries.

For example, if your organization hosts a private [PyPI](https://pypi.org/)
repository on a private network, Docker Build Cloud would not be able to access
it by default, since the cloud builder is not connected to your private network.

To enable your cloud builders to access your private resources, enter the host
name and port of your private resource and then select **Add**.

#### Authentication

If your internal artifacts require authentication, make sure that you
authenticate with the repository either before or during the build. For internal
package repositories for npm or PyPI, use [build
secrets](/manuals/build/building/secrets.md) to authenticate during the build.
For internal OCI registries, use `docker login` to authenticate before building.

Note that if you use a private registry that requires authentication, you will
need to authenticate with `docker login` twice before building. This is because
the cloud builder needs to authenticate with Docker to use the cloud builder,
and then again to authenticate with the private registry.

```console
$ echo $DOCKER_PAT | docker login docker.io -u <username> --password-stdin
$ echo $REGISTRY_PASSWORD | docker login registry.example.com -u <username> --password-stdin
$ docker build --builder <cloud-builder> --tag registry.example.com/<image> --push .
```

### Firewall

Firewall settings let you restrict cloud builder egress traffic to specific IP
addresses. This helps enhance security by limiting external network egress from
the builder.

1. Select **Enable firewall: Restrict cloud builder egress to specific public IP address**.

2. Enter the IP address you want to allow.

3. Select **Add** to apply the restriction.
