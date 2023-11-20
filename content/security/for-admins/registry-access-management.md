---
description: Registry Access Management
keywords: registry, access, management, permissions, Docker Business feature
title: Registry Access Management
aliases:
- /desktop/hardened-desktop/registry-access-management/
- /admin/organization/registry-access/
- /docker-hub/registry-access-management/
---

> Note
>
> Registry Access Management is available to [Docker Business](../../subscription/details.md) customers only. 

With Registry Access Management (RAM), administrators can ensure that their developers using Docker Desktop only access allowed registries. This is done through the Registry Access Management dashboard on Docker Hub. 

Registry Access Management supports both cloud and on-prem registries. Example registries administrators can allow include: 
 - Docker Hub. This is enabled by default.
 - Amazon ECR
 - GitHub Container Registry
 - Google Container Registry
 - Nexus
 - Artifactory

## Prerequisites 

You need to [configure a registry.json to enforce sign-in](/docker-hub/configure-sign-in/). For Registry Access Management to take effect, Docker Desktop users must authenticate to your organization.

## Configure Registry Access Management permissions

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-registry-access product="hub" %}}

{{< /tab >}}
{{< tab name="Docker Admin" >}}

{{< include "admin-early-access.md" >}}

{{% admin-registry-access product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## Verify the restrictions

The new Registry Access Management policy takes effect after the developer successfully authenticates to Docker Desktop using their organization credentials. If a developer attempts to pull an image from a disallowed registry via the Docker CLI, they receive an error message that the organization has disallowed this registry.

## Caveats

There are certain limitations when using Registry Access Management:

- Windows image pulls and image builds are not restricted by default. For Registry Access Management to take effect on Windows Container mode, you must allow the Windows Docker daemon to use Docker Desktop's internal proxy by selecting the [Use proxy for Windows Docker daemon](../../desktop/settings/windows.md/#proxies) setting.
- Builds such as `docker buildx` using a Kubernetes driver are not restricted
- Builds such as `docker buildx` using a custom docker-container driver are not restricted
- Blocking is DNS-based; you must use a registry's access control mechanisms to distinguish between “push” and “pull”
- WSL 2 requires at least a 5.4 series Linux kernel (this does not apply to earlier Linux kernel series)
- Under the WSL 2 network, traffic from all Linux distributions is restricted (this will be resolved in the updated 5.15 series Linux kernel)

Also, Registry Access Management operates on the level of hosts, not IP addresses. Developers can bypass this restriction within their domain resolution, for example by running Docker against a local proxy or modifying their operating system's `sts` file. Blocking these forms of manipulation is outside the remit of Docker Desktop.
