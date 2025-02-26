---
title: Builder settings
description: Set your builder settings relating to private registries, disk allocation .
keywords: build, cloud build, optimize, remote, local, cloud, registry, package repository, vpn
---

The **Builder settings** page in Docker Build Cloud lets you configure disk allocation, private resource access, and firewall settings for your cloud builders in your organization. These configurations help optimize storage, enable access to private registries, and secure outbound network traffic.

## Disk allocation

The **Disk allocation** setting lets you control how much of the available storage is dedicated to the build cache. A lower allocation increases storage available for active builds.

To make disk allocation changes, navigate to **Builder settings** in Docker Build Cloud and then adjust the **Disk allocation** slider to specify the percentage of storage used for build caching.

Any changes take effect immediately.

> [!TIP]
> 
> If you build very large images, consider allocating less storage for caching.

## Private resource access

Private resource access lets cloud builders pull images and packages from private resources. This feature is useful when builds rely on self-hosted artifact repositories or private OCI registries.

For example, if your organization hosts a private [PyPI](https://pypi.org/) repository on a private network, Docker Build Cloud would not be able to access it by default, since the cloud builder is not connected to your private network.

To enable your cloud builders to access your private resources, enter the host name and port of your private rescource and then select **Add**.

### Authentication 

If your internal artifacts require authentication, make sure that you
authenticate with the repository either before or during the build. For
internal package repositories for npm or PyPI, use [build secrets](/manuals/build/building/secrets.md)
to authenticate during the build. For internal OCI registries, use `docker
login` to authenticate before building.

Note that if you use a private registry that requires authentication, you will
need to authenticate with `docker login` twice before building. This is because
the cloud builder needs to authenticate with Docker to use the cloud builder,
and then again to authenticate with the private registry.

```console
$ echo $DOCKER_PAT | docker login docker.io -u <username> --password-stdin
$ echo $REGISTRY_PASSWORD | docker login registry.example.com -u <username> --password-stdin
$ docker build --builder <cloud-builder> --tag registry.example.com/<image> --push .
```

## Firewall

Firewall settings let you restrict cloud builder egress traffic to specific IP addresses. This helps enhance security by limiting external network egress from the builder.

1. Select the **Enable firewall: Restrict cloud builder egress to specific public IP address** checkbox.

2. Enter the IP address you want to allow.

3. Select **Add** to apply the restriction.
